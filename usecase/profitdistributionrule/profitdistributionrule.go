package profitdistributionrule

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProfitDistributionRuleInteractor = (*ProfitDistributionRuleInteractor)(nil)

type ProfitDistributionRuleInteractor struct {
	DS     domain.DataStore
	Config domain.ProfitDistributionConfig
}

func NewProfitDistributionRuleInteractor(ds domain.DataStore, profitDistributionConfig domain.ProfitDistributionConfig) *ProfitDistributionRuleInteractor {

	fmt.Println(profitDistributionConfig)

	return &ProfitDistributionRuleInteractor{
		DS:     ds,
		Config: profitDistributionConfig,
	}
}

func (i *ProfitDistributionRuleInteractor) Create(ctx context.Context, rule *domain.ProfitDistributionRule, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProfitDistributionRuleInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if err = validateProfitDistributionRuleParams(rule); err != nil {
		return err
	}

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 业务规则校验
		if err = validateProfitDistributionRuleBusinessRules(ctx, ds, rule, uuid.Nil); err != nil {
			return err
		}
		// 创建分账方案
		return ds.ProfitDistributionRuleRepo().Create(ctx, rule)
	})
}

func (i *ProfitDistributionRuleInteractor) Update(ctx context.Context, rule *domain.ProfitDistributionRule, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProfitDistributionRuleInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	// 检查时间窗口：修改操作必须在定时任务执行后
	if err = i.checkModifyTimeWindow(); err != nil {
		return err
	}
	if err = validateProfitDistributionRuleParams(rule); err != nil {
		return err
	}

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证分账方案存在
		existingRule, err := ds.ProfitDistributionRuleRepo().FindByID(ctx, rule.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProfitDistributionRuleNotExists)
			}
			return err
		}

		// 验证是否属于当前品牌商
		if existingRule.MerchantID != user.GetMerchantID() {
			return domain.ParamsError(domain.ErrProfitDistributionRuleNotExists)
		}

		rule.MerchantID = existingRule.MerchantID
		rule.Status = existingRule.Status

		// 业务规则校验（排除自身）
		if err = validateProfitDistributionRuleBusinessRules(ctx, ds, rule, rule.ID); err != nil {
			return err
		}

		// 更新分账方案
		return ds.ProfitDistributionRuleRepo().Update(ctx, rule)
	})
}

func (i *ProfitDistributionRuleInteractor) Delete(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProfitDistributionRuleInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 检查时间窗口：删除操作必须在定时任务执行后
	if err = i.checkModifyTimeWindow(); err != nil {
		return err
	}

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证分账方案存在
		rule, err := ds.ProfitDistributionRuleRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProfitDistributionRuleNotExists)
			}
			return err
		}

		// 验证是否属于当前品牌商
		if rule.MerchantID != user.GetMerchantID() {
			return domain.ParamsError(domain.ErrProfitDistributionRuleNotExists)
		}

		// 删除分账方案
		return ds.ProfitDistributionRuleRepo().Delete(ctx, id)
	})
}

func (i *ProfitDistributionRuleInteractor) Enable(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProfitDistributionRuleInteractor.Enable")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 检查时间窗口：启用操作必须在定时任务执行后
	if err = i.checkModifyTimeWindow(); err != nil {
		return err
	}

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证分账方案存在
		rule, err := ds.ProfitDistributionRuleRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProfitDistributionRuleNotExists)
			}
			return err
		}

		// 验证是否属于当前品牌商
		if rule.MerchantID != user.GetMerchantID() {
			return domain.ParamsError(domain.ErrProfitDistributionRuleNotExists)
		}

		// 验证状态：只有禁用状态才能启用
		if rule.Status != domain.ProfitDistributionRuleStatusDisabled {
			return domain.ParamsError(domain.ErrProfitDistributionRuleStatusInvalid)
		}

		// 更新状态为启用
		rule.Status = domain.ProfitDistributionRuleStatusEnabled
		return ds.ProfitDistributionRuleRepo().Update(ctx, rule)
	})
}

func (i *ProfitDistributionRuleInteractor) Disable(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProfitDistributionRuleInteractor.Disable")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 检查时间窗口：禁用操作必须在定时任务执行后
	if err = i.checkModifyTimeWindow(); err != nil {
		return err
	}

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证分账方案存在
		rule, err := ds.ProfitDistributionRuleRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProfitDistributionRuleNotExists)
			}
			return err
		}

		// 验证是否属于当前品牌商
		if rule.MerchantID != user.GetMerchantID() {
			return domain.ParamsError(domain.ErrProfitDistributionRuleNotExists)
		}

		// 验证状态：只有启用状态才能禁用
		if rule.Status != domain.ProfitDistributionRuleStatusEnabled {
			return domain.ParamsError(domain.ErrProfitDistributionRuleStatusInvalid)
		}

		// 更新状态为禁用
		rule.Status = domain.ProfitDistributionRuleStatusDisabled
		return ds.ProfitDistributionRuleRepo().Update(ctx, rule)
	})
}

func (i *ProfitDistributionRuleInteractor) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.ProfitDistributionRuleSearchParams,
) (res *domain.ProfitDistributionRuleSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProfitDistributionRuleInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.ProfitDistributionRuleRepo().PagedListBySearch(ctx, page, params)
}

// ============================================
// 校验函数
// ============================================

// validateProfitDistributionRuleParams 验证分账方案参数
func validateProfitDistributionRuleParams(rule *domain.ProfitDistributionRule) error {
	// 1. 验证分账比例（0-1之间）
	one := decimal.NewFromInt(1)
	if rule.SplitRatio.IsNegative() || rule.SplitRatio.GreaterThan(one) {
		return domain.ParamsError(domain.ErrProfitDistributionRuleSplitRatioInvalid)
	}

	// 2. 验证日期：生效日期必须必须晚于今天，生效日期要早于失效日期
	if rule.EffectiveDate.Before(time.Now().AddDate(0, 0, 1)) {
		return domain.ParamsError(domain.ErrProfitDistributionRuleDateInvalid)
	}
	if rule.ExpiryDate.Before(rule.EffectiveDate) {
		return domain.ParamsError(domain.ErrProfitDistributionRuleDateInvalid)
	}
	return nil
}

func validateProfitDistributionRuleBusinessRules(ctx context.Context, ds domain.DataStore, rule *domain.ProfitDistributionRule, excludeRuleID uuid.UUID) error {
	// 1. 检查分账方案名称是否唯一
	exists, err := ds.ProfitDistributionRuleRepo().Exists(ctx, domain.ProfitDistributionRuleExistsParams{
		MerchantID: rule.MerchantID,
		Name:       rule.Name,
		ExcludeID:  excludeRuleID,
	})
	if err != nil {
		return err
	}
	if exists {
		return domain.ParamsError(domain.ErrProfitDistributionRuleNameExists)
	}

	storeIDs := lo.Map(rule.Stores, func(store *domain.StoreSimple, _ int) uuid.UUID {
		return store.ID
	})

	// 2. 检查门店是否有效且属于当前品牌商
	if len(storeIDs) > 0 {
		stores, err := ds.StoreRepo().ListByIDs(ctx, storeIDs)
		if err != nil {
			return err
		}
		if len(stores) != len(storeIDs) {
			return domain.ParamsError(fmt.Errorf("部分门店不存在"))
		}
		for _, store := range stores {
			if store.MerchantID != rule.MerchantID {
				return domain.ParamsError(fmt.Errorf("门店 %s 不属于当前品牌商", store.ID))
			}
		}
	}

	// 3. 检查门店是否已绑定其他分账方案
	if len(storeIDs) > 0 {
		hasBound, err := ds.ProfitDistributionRuleRepo().CheckStoreBound(ctx, storeIDs, excludeRuleID)
		if err != nil {
			return err
		}
		if hasBound {
			return domain.ParamsError(domain.ErrProfitDistributionRuleStoreBound)
		}
	}
	return nil
}

// checkModifyTimeWindow 检查是否在允许修改的时间窗口内
// 修改操作必须在定时任务执行时间之后才能进行
func (i *ProfitDistributionRuleInteractor) checkModifyTimeWindow() error {
	now := time.Now()
	currentHour := now.Hour()
	currentMinute := now.Minute()
	currentTimeMinutes := currentHour*60 + currentMinute
	taskTimeMinutes := i.Config.TaskHour*60 + i.Config.TaskMinute

	fmt.Println(i.Config.TaskHour, i.Config.TaskMinute)
	// 如果当前时间早于定时任务执行时间，不允许修改
	if currentTimeMinutes < taskTimeMinutes {
		return domain.ParamsErrorf("修改操作必须在定时任务执行时间之后才能进行")
	}
	return nil
}
