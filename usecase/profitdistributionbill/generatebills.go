package profitdistributionbill

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

const (
	maxWorkers = 10 // 最大并发数
)

type jobResult struct {
	merchantID uuid.UUID
	err        error
}

func (i *ProfitDistributionBillInteractor) GenerateProfitDistributionBills(ctx context.Context) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProfitDistributionBillInteractor.GenerateProfitDistributionBills")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	logger := logging.FromContext(ctx).Named("ProfitDistributionBillInteractor.GenerateProfitDistributionBills")
	ctx = logging.NewContext(ctx, logger)

	// 1. 获取所有启用的分账方案
	rules, err := i.DS.ProfitDistributionRuleRepo().ListAllEnabled(ctx)
	if err != nil {
		logger.Errorf("获取启用的分账方案失败: %v", err)
		return err
	}

	if len(rules) == 0 {
		logger.Info("没有启用的分账方案，跳过生成账单")
		return nil
	}
	// 2. 过滤今天可以执行的方案
	executableRules := i.filterTodayExecutableRules(rules)
	if len(executableRules) == 0 {
		logger.Info("今天没有可执行的分账方案，跳过生成账单")
		return nil
	}
	logger.Infof("找到 %d 个今天可执行的分账方案", len(executableRules))

	// 3. 按品牌商进行分组
	rulesByMerchant := make(map[uuid.UUID]domain.ProfitDistributionRules)
	for _, rule := range rules {
		rulesByMerchant[rule.MerchantID] = append(rulesByMerchant[rule.MerchantID], rule)
	}

	// 4. 使用协程池并发处理品牌商
	merchantIDs := lo.Keys(rulesByMerchant)
	jobs := make(chan uuid.UUID, len(merchantIDs))
	results := make(chan jobResult, len(merchantIDs))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Worker: 处理单个品牌商的所有门店
	worker := func(workerID int) {
		for merchantID := range jobs {
			select {
			case <-ctx.Done():
				return
			default:
				logger.Infof("开始处理品牌商任务：worker_id=%d, merchant_id=%s", workerID, merchantID.String())
				err := i.processMerchantBills(ctx, rulesByMerchant[merchantID])
				results <- jobResult{
					merchantID: merchantID,
					err:        err,
				}
			}
		}
	}

	// 启动 worker
	workerCount := maxWorkers
	if len(merchantIDs) < maxWorkers {
		workerCount = len(merchantIDs)
	}
	for i := 0; i < workerCount; i++ {
		go worker(i)
	}

	// 投递任务
	for _, merchantID := range merchantIDs {
		jobs <- merchantID
	}
	close(jobs)

	// 收集结果
	var failed []string
	for range len(merchantIDs) {
		r := <-results
		if r.err != nil {
			failed = append(
				failed,
				fmt.Sprintf("merchant_id=%s err=%v", r.merchantID.String(), r.err),
			)
			logger.Errorf("品牌商 %s 处理失败: %v", r.merchantID.String(), r.err)
		}
	}

	if len(failed) > 0 {
		err = fmt.Errorf("部分品牌商生成分账账单失败: %s", strings.Join(failed, "; "))
		logger.Errorf("%v", err)
		return err
	}

	logger.Infof("成功生成所有品牌商的分账账单")
	return nil
}

// filterTodayExecutableRules 过滤今天可以执行的方案
func (i *ProfitDistributionBillInteractor) filterTodayExecutableRules(rules domain.ProfitDistributionRules) domain.ProfitDistributionRules {
	todayDate := util.DayStart(time.Now())
	todayDay := todayDate.Day()

	executableRules := make(domain.ProfitDistributionRules, 0)
	for _, rule := range rules {
		effectiveDate := util.DayStart(rule.EffectiveDate)
		expiryDate := util.DayStart(rule.ExpiryDate)
		// 生效日期：today >= effectiveDate（今天大于等于生效日期）
		// 失效日期：today < expiryDate（今天小于失效日期，不包含失效日期当天）
		// 例如：1-2生效，1-3过期，则1-2可以执行，1-3 00:00:00就失效了
		if todayDate.Before(effectiveDate) || !todayDate.Before(expiryDate) {
			continue
		}
		// 根据账单生成周期判断
		switch rule.BillingCycle {
		case domain.ProfitDistributionRuleBillingCycleDaily:
			// 按日生成：今天在生效日期和失效日期内即可
			executableRules = append(executableRules, rule)
		case domain.ProfitDistributionRuleBillingCycleMonthly:
			// 按月生成：需要检查账单生成日是否为今天
			if todayDay == rule.BillGenerationDay {
				executableRules = append(executableRules, rule)
			}
		}
	}

	return executableRules
}

// processMerchantBills 处理单个品牌商的所有门店账单
func (i *ProfitDistributionBillInteractor) processMerchantBills(
	ctx context.Context,
	rules domain.ProfitDistributionRules,
) error {
	// logger := logging.FromContext(ctx).Named("processMerchantBills")
	merchantID := rules[0].MerchantID
	today := util.DayStart(time.Now())

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 收集所有门店ID
		storeIDSet := make(map[uuid.UUID]*domain.ProfitDistributionRule)
		for _, rule := range rules {
			for _, store := range rule.Stores {
				// 如果门店已存在，报错
				if _, exists := storeIDSet[store.ID]; exists {
					return fmt.Errorf("门店 [%s] 规则重复", store.StoreName)
				}
				storeIDSet[store.ID] = rule
			}
		}
		if len(storeIDSet) == 0 {
			return nil
		}

		// 按账单周期分组处理
		dailyBills := make(domain.ProfitDistributionBills, 0)
		monthlyBills := make(domain.ProfitDistributionBills, 0)

		dailyStoreIDs := make([]uuid.UUID, 0)
		monthlyStoreIDs := make([]uuid.UUID, 0)

		for storeID, rule := range storeIDSet {
			if rule.BillingCycle == domain.ProfitDistributionRuleBillingCycleDaily {
				dailyStoreIDs = append(dailyStoreIDs, storeID)
			} else {
				monthlyStoreIDs = append(monthlyStoreIDs, storeID)
			}
		}

		// 处理按日生成的账单
		if len(dailyStoreIDs) > 0 {
			bills, err := i.generateDailyBills(ctx, ds, merchantID, dailyStoreIDs, storeIDSet, today)
			if err != nil {
				return fmt.Errorf("生成按日账单失败: %w", err)
			}
			dailyBills = bills
		}

		// 处理按月生成的账单
		if len(monthlyStoreIDs) > 0 {
			bills, err := i.generateMonthlyBills(ctx, ds, merchantID, monthlyStoreIDs, storeIDSet, today)
			if err != nil {
				return fmt.Errorf("生成按月账单失败: %w", err)
			}
			monthlyBills = bills
		}

		// 批量创建账单
		allBills := append(dailyBills, monthlyBills...)

		util.PrettyJson(allBills)

		if len(allBills) > 0 {
			// if err := ds.ProfitDistributionBillRepo().CreateBulk(ctx, allBills); err != nil {
			// 	return fmt.Errorf("批量创建账单失败: %w", err)
			// }
			// logger.Infof("品牌商 %s 成功创建 %d 条账单", merchantID.String(), len(allBills))
		}
		return nil
	})
}

// @TODO 获取营业额
type StoreDailyRevenue struct {
	ID      uuid.UUID       `json:"id"`
	StoreID uuid.UUID       `json:"store_id"`
	Date    time.Time       `json:"date"`
	Amount  decimal.Decimal `json:"amount"`
}

type StoreDailyRevenues []*StoreDailyRevenue

// generateDailyBills 生成按日账单
func (i *ProfitDistributionBillInteractor) generateDailyBills(
	ctx context.Context,
	ds domain.DataStore,
	merchantID uuid.UUID,
	storeIDs []uuid.UUID,
	storeIDToRule map[uuid.UUID]*domain.ProfitDistributionRule,
	today time.Time,
) (domain.ProfitDistributionBills, error) {
	// 获取昨日营业额（按日账单）
	yesterday := util.DayStart(today.AddDate(0, 0, -1))
	// @TODO 获取营业额
	// revenues, err := ds.StoreDailyRevenueRepo().FindByStoreIDsAndDateRange(ctx, storeIDs, billDate)
	// if err != nil {
	// 	return nil, err
	// }

	// mock 一个营业额
	revenues := make(StoreDailyRevenues, 0)
	for _, storeID := range storeIDs {
		revenues = append(revenues, &StoreDailyRevenue{
			ID:      uuid.New(),
			StoreID: storeID,
			Date:    yesterday,
			Amount:  decimal.NewFromInt(100),
		})
	}

	// 建立门店ID到营业额的映射
	revenueMap := lo.SliceToMap(revenues, func(r *StoreDailyRevenue) (uuid.UUID, *StoreDailyRevenue) {
		return r.StoreID, r
	})

	bills := make(domain.ProfitDistributionBills, 0)
	for _, storeID := range storeIDs {
		rule := storeIDToRule[storeID]
		revenue, exists := revenueMap[storeID]
		// 如果没有营业额记录，跳过
		if !exists || revenue == nil {
			continue
		}
		// 计算分账金额
		receivableAmount := revenue.Amount.Mul(rule.SplitRatio)

		// // 生成账单编号
		// billNo, err := i.generateBillNo(ctx, billDate)
		// if err != nil {
		// 	return nil, fmt.Errorf("生成账单编号失败: %w", err)
		// }

		// 创建账单
		bill := &domain.ProfitDistributionBill{
			ID:               uuid.New(),
			No:               util.RandomString(10),
			MerchantID:       merchantID,
			StoreID:          storeID,
			RevenueID:        revenue.ID,
			ReceivableAmount: receivableAmount,
			PaymentAmount:    decimal.Zero,
			Status:           domain.ProfitDistributionBillStatusUnpaid,
			BillDate:         today,
			StartDate:        yesterday,
			EndDate:          yesterday,
			RuleSnapshot: &domain.ProfitDistributionRuleSnapshot{
				RuleID:     rule.ID,
				RuleName:   rule.Name,
				SplitRatio: rule.SplitRatio,
			},
		}

		bills = append(bills, bill)
	}

	return bills, nil
}

// generateMonthlyBills 生成按月账单
func (i *ProfitDistributionBillInteractor) generateMonthlyBills(
	ctx context.Context,
	ds domain.DataStore,
	merchantID uuid.UUID,
	storeIDs []uuid.UUID,
	storeIDToRule map[uuid.UUID]*domain.ProfitDistributionRule,
	today time.Time,
) (domain.ProfitDistributionBills, error) {
	return nil, nil
	// // 获取上月营业额（按月账单）
	// lastMonth := today.AddDate(0, -1, 0)
	// startDate := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, lastMonth.Location())

	// // 计算上月的最后一天
	// nextMonth := startDate.AddDate(0, 1, 0)
	// endDate := nextMonth.AddDate(0, 0, -1)

	// revenues, err := ds.StoreDailyRevenueRepo().FindByStoreIDsAndDateRange(ctx, storeIDs, startDate, endDate)
	// if err != nil {
	// 	return nil, err
	// }

	// // 按门店汇总营业额
	// storeRevenueMap := make(map[uuid.UUID]*domain.StoreDailyRevenue)
	// for _, revenue := range revenues {
	// 	if existing, ok := storeRevenueMap[revenue.StoreID]; ok {
	// 		existing.Revenue = existing.Revenue.Add(revenue.Revenue)
	// 	} else {
	// 		storeRevenueMap[revenue.StoreID] = &domain.StoreDailyRevenue{
	// 			ID:      revenue.ID, // 使用第一条记录的ID
	// 			StoreID: revenue.StoreID,
	// 			Revenue: revenue.Revenue,
	// 			Date:    startDate,
	// 		}
	// 	}
	// }

	// bills := make(domain.ProfitDistributionBills, 0)
	// for _, storeID := range storeIDs {
	// 	rule := storeIDToRule[storeID]
	// 	revenue, exists := storeRevenueMap[storeID]

	// 	// 如果没有营业额记录，跳过
	// 	if !exists || revenue == nil {
	// 		continue
	// 	}

	// 	// 检查是否已存在账单
	// 	billDate := today
	// 	billDateOnly := time.Date(billDate.Year(), billDate.Month(), billDate.Day(), 0, 0, 0, 0, billDate.Location())
	// 	existingBill, err := ds.ProfitDistributionBillRepo().FindByStoreIDAndBillDate(ctx, storeID, billDateOnly)
	// 	if err != nil && !domain.IsNotFound(err) {
	// 		return nil, err
	// 	}
	// 	if existingBill != nil {
	// 		continue // 已存在账单，跳过
	// 	}

	// 	// 计算分账金额
	// 	receivableAmount := revenue.Revenue.Mul(rule.SplitRatio)
	// 	paymentAmount := receivableAmount

	// 	// 生成账单编号
	// 	billNo, err := i.generateBillNo(ctx, billDateOnly)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("生成账单编号失败: %w", err)
	// 	}

	// 	// 创建账单
	// 	bill := &domain.ProfitDistributionBill{
	// 		ID:               uuid.New(),
	// 		No:               billNo,
	// 		MerchantID:       merchantID,
	// 		StoreID:          storeID,
	// 		RevenueID:        revenue.ID,
	// 		ReceivableAmount: receivableAmount,
	// 		PaymentAmount:    paymentAmount,
	// 		Status:           domain.ProfitDistributionBillStatusUnpaid,
	// 		BillDate:         billDateOnly,
	// 		StartDate:        startDate,
	// 		EndDate:          endDate,
	// 		RuleSnapshot: &domain.ProfitDistributionRuleSnapshot{
	// 			RuleID:     rule.ID,
	// 			RuleName:   rule.Name,
	// 			SplitRatio: rule.SplitRatio,
	// 		},
	// 	}

	// 	bills = append(bills, bill)
	// }

	// return bills, nil
}
