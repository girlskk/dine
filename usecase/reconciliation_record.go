package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ReconciliationRecordInteractor = (*ReconciliationRecordInteractor)(nil)

type ReconciliationRecordInteractor struct {
	ds domain.DataStore
}

func NewReconciliationRecordInteractor(dataStore domain.DataStore) *ReconciliationRecordInteractor {
	return &ReconciliationRecordInteractor{
		ds: dataStore,
	}
}

func (r *ReconciliationRecordInteractor) GenerateDailyRecords(ctx context.Context) (err error) {
	// 获取昨天的已完成订单列表
	yesterday := util.DayStart(time.Now()).AddDate(0, 0, -1)
	return r.GenerateDailyRecordsByDate(ctx, yesterday)
}

// GenerateDailyRecordsByDate 生成指定日期的财务对账单
func (r *ReconciliationRecordInteractor) GenerateDailyRecordsByDate(ctx context.Context, recordDate time.Time) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationRecordInteractor.GenerateDailyRecordsByDate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	dayStart := util.DayStart(recordDate)
	dayEnd := util.DayEnd(recordDate)

	orders, total, err := r.ds.OrderRepo().GetOrders(ctx, &upagination.Pagination{Page: 1, Size: upagination.MaxSize},
		&domain.OrderListFilter{
			Status:        domain.OrderStatusPaid,
			FinishedAtGte: &dayStart,
			FinishedAtLte: &dayEnd,
		})
	if err != nil {
		return err
	}
	if total == 0 {
		return nil
	}

	// 获取所有门店的费率信息
	stores, err := r.ds.StoreRepo().ListAll(ctx)
	if err != nil {
		return err
	}
	if len(stores) == 0 {
		return nil
	}
	storeIDMap := lo.SliceToMap(stores, func(item *domain.Store) (int, *domain.Store) { return item.ID, item })

	// 根据 “门店ID+支付渠道” 进行分组
	storeChannelMap := make(map[string]*domain.ReconciliationRecord)
	storePointSettlementMap := make(map[int]*domain.PointSettlement)
	for _, order := range orders {
		processPaymentChannelReconciliation(order, storeChannelMap, domain.OrderPaidChannelCash, order.CashPaid, recordDate)
		processPaymentChannelReconciliation(order, storeChannelMap, domain.OrderPaidChannelWechatPay, order.WechatPaid, recordDate)
		processPaymentChannelReconciliation(order, storeChannelMap, domain.OrderPaidChannelAlipay, order.AlipayPaid, recordDate)
		processPaymentChannelReconciliation(order, storeChannelMap, domain.OrderPaidChannelPoint, order.PointsPaid, recordDate)
		processPaymentChannelReconciliation(order, storeChannelMap, domain.OrderPaidChannelPointWallet, order.PointsWalletPaid, recordDate)
		// 统计积分结算账单
		if err = processPointSettlement(order, storePointSettlementMap, storeIDMap, recordDate); err != nil {
			return err
		}
	}
	// 财务对账单
	reconciliationRecords := lo.Values(storeChannelMap)

	storePointSettlementMap = lo.MapValues(storePointSettlementMap, func(item *domain.PointSettlement, _ int) *domain.PointSettlement {
		item.Amount = item.TotalPoints.Mul(decimal.NewFromFloat(1).Sub(item.PointSettlementRate))
		return item
	})

	pointSettlements := lo.Values(storePointSettlementMap)

	return r.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 创建财务对账单
		if len(reconciliationRecords) > 0 {
			if err := ds.ReconciliationRecordRepo().BatchCreate(ctx, reconciliationRecords); err != nil {
				return fmt.Errorf("保存财务对账单失败：%v", err)
			}
		}
		if len(pointSettlements) > 0 {
			if err := ds.PointSettlementRepo().BatchCreate(ctx, pointSettlements); err != nil {
				return fmt.Errorf("保存积分结算账单失败：%v", err)
			}
		}
		return nil
	})
}

// 不同支付渠道的财务对账处理
func processPaymentChannelReconciliation(
	order *domain.Order,
	storeChannelMap map[string]*domain.ReconciliationRecord,
	channel domain.OrderPaidChannel,
	channelAmount decimal.Decimal,
	recordDate time.Time,
) {
	if channelAmount.LessThanOrEqual(decimal.Zero) {
		return
	}
	key := fmt.Sprintf("%d:%s", order.StoreID, channel)
	if value, exists := storeChannelMap[key]; exists {
		value.OrderCount++
		value.Amount = value.Amount.Add(channelAmount)
	} else {
		storeChannelMap[key] = &domain.ReconciliationRecord{
			// PAY YYMMDD SSSS (PAY 支付方式，SSSS 是门店ID)
			No: fmt.Sprintf("%s%s%04d", channel.Prefix(),
				recordDate.AddDate(0, 0, 1).Format("060102"), order.StoreID),
			StoreID:    order.StoreID,
			StoreName:  order.StoreName,
			Channel:    channel,
			OrderCount: 1,
			Amount:     channelAmount,
			Date:       recordDate,
		}
	}
}

// 积分结算账单处理
func processPointSettlement(
	order *domain.Order,
	storePointSettlementMap map[int]*domain.PointSettlement,
	storeIDMap map[int]*domain.Store,
	recordDate time.Time,
) error {
	if order.PointsPaid.LessThanOrEqual(decimal.Zero) {
		return nil
	}
	if value, exists := storePointSettlementMap[order.StoreID]; exists {
		value.OrderCount++
		value.TotalPoints = value.TotalPoints.Add(order.PointsPaid)
	} else {
		store, ok := storeIDMap[order.StoreID]
		if !ok {
			return fmt.Errorf("门店不存在")
		}
		storePointSettlementMap[order.StoreID] = &domain.PointSettlement{
			// PAY YYMMDD SSSS (PAY 支付方式，SSSS 是门店ID)
			No:                  fmt.Sprintf("%s%s%04d", domain.OrderPaidChannelPoint.Prefix(), time.Now().Format("060102"), order.StoreID),
			StoreID:             order.StoreID,
			StoreName:           order.StoreName,
			OrderCount:          1,
			TotalPoints:         order.PointsPaid,
			Date:                recordDate,
			Status:              domain.PointSettlementStatusPending,
			PointSettlementRate: store.PointSettlementRate,
		}
	}
	return nil
}

func (r *ReconciliationRecordInteractor) PagedListBySearch(ctx context.Context, page *upagination.Pagination,
	params domain.ReconciliationSearchParams,
) (res *domain.ReconciliationSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationRecordInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return r.ds.ReconciliationRecordRepo().PagedListBySearch(ctx, page, params)
}

func (r *ReconciliationRecordInteractor) ListDetails(ctx context.Context, id, storeID int,
) (res domain.ReconciliationDetails, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationRecordInteractor.ListDetails")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	record, err := r.ds.ReconciliationRecordRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ParamsError(domain.ErrReconciliationRecordNotExists)
		}
		return nil, err
	}
	if storeID > 0 && record.StoreID != storeID {
		return nil, domain.ParamsError(domain.ErrReconciliationRecordNotExists)
	}

	dayStart := util.DayStart(record.Date)
	dayEnd := util.DayEnd(record.Date)
	filter := &domain.OrderListFilter{
		Status:        domain.OrderStatusPaid,
		FinishedAtGte: &dayStart,
		FinishedAtLte: &dayEnd,
	}
	orders, total, err := r.ds.OrderRepo().GetOrders(ctx, &upagination.Pagination{
		Page: 1,
		Size: upagination.MaxSize,
	}, filter)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return res, nil
	}
	var orderIDs []int
	for _, order := range orders {
		if order.PaidChannels.Contains(record.Channel) {
			item := &domain.ReconciliationDetail{
				OrderID:         order.ID,
				OrderNo:         order.No,
				OrderAmount:     order.TotalPrice,
				OrderFinishedAt: *order.FinishedAt,
			}
			switch record.Channel {
			case domain.OrderPaidChannelCash:
				item.Amount = order.CashPaid
			case domain.OrderPaidChannelPoint:
				item.Amount = order.PointsPaid
			case domain.OrderPaidChannelAlipay:
				item.Amount = order.AlipayPaid
			case domain.OrderPaidChannelWechatPay:
				item.Amount = order.WechatPaid
			case domain.OrderPaidChannelPointWallet:
				item.Amount = order.PointsWalletPaid
			}
			res = append(res, item)
			orderIDs = append(orderIDs, order.ID)
		}
	}
	if len(orderIDs) == 0 {
		return res, nil
	}
	orderItemNamesMap, err := r.ds.OrderRepo().ListItemNamesByOrders(ctx, orderIDs)
	if err != nil {
		return nil, err
	}
	for _, item := range res {
		item.ProductInfo = orderItemNamesMap[item.OrderID]
	}
	return res, nil
}

func (r *ReconciliationRecordInteractor) Summary(ctx context.Context,
	params domain.ReconciliationSearchParams,
) (res *domain.ReconciliationSummaryRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationRecordInteractor.Summary")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	page := &upagination.Pagination{
		Page: 1,
		Size: upagination.MaxSize,
	}
	records, err := r.ds.ReconciliationRecordRepo().PagedListBySearch(ctx, page, params)
	if err != nil {
		return nil, err
	}

	res = &domain.ReconciliationSummaryRes{}

	for _, record := range records.Items {
		res.TotalCount += record.OrderCount
		res.TotalAmount = res.TotalAmount.Add(record.Amount)
		switch record.Channel {
		case domain.OrderPaidChannelCash:
			res.CashAmount = res.CashAmount.Add(record.Amount)
		case domain.OrderPaidChannelWechatPay:
			res.WechatAmount = res.WechatAmount.Add(record.Amount)
		case domain.OrderPaidChannelAlipay:
			res.AlipayAmount = res.AlipayAmount.Add(record.Amount)
		case domain.OrderPaidChannelPoint:
			res.PointAmount = res.PointAmount.Add(record.Amount)
		case domain.OrderPaidChannelPointWallet:
			res.PointWalletAmount = res.PointWalletAmount.Add(record.Amount)
		}
	}
	return res, nil
}

func (r *ReconciliationRecordInteractor) GetReconciliationRange(ctx context.Context, params domain.ReconciliationSearchParams,
) (res domain.ReconciliationRange, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationRecordInteractor.GetReconciliationRange")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	res, err = r.ds.ReconciliationRecordRepo().GetReconciliationRange(ctx, params)
	return
}
