package order

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

type DomainService struct {
	OrderEventTrigger domain.OrderEventTrigger
	ObjectStorage     domain.ObjectStorage
	DataStore         domain.DataStore
}

func NewDomainService(trigger domain.OrderEventTrigger, objectStorage domain.ObjectStorage, dataStore domain.DataStore) *DomainService {
	return &DomainService{
		OrderEventTrigger: trigger,
		ObjectStorage:     objectStorage,
		DataStore:         dataStore,
	}
}

type PaidParams struct {
	Order       *domain.Order
	Amount      decimal.Decimal
	Channel     domain.OrderPaidChannel
	Operator    *domain.FrontendUser
	MemberID    int
	MemberName  string
	MemberPhone string
	SeqNo       string
}

// 订单支付逻辑，必须在事务中调用
func (s *DomainService) Paid(ctx context.Context, ds domain.DataStore, params *PaidParams) (updateOrder *domain.Order, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderDomainService.Paid")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if params.Amount.IsNegative() {
		err = domain.ParamsErrorf("支付金额必须大于0")
		return
	}

	if !ds.IsTransactionActive() {
		panic("transaction not started")
	}

	od := params.Order

	if err = od.CanPaid(params.Amount, params.Channel == domain.OrderPaidChannelPoint); err != nil {
		return
	}

	if params.Amount.IsZero() && (!od.RealPrice.IsZero() || params.Channel != domain.OrderPaidChannelCash) {
		err = domain.ParamsErrorf("支付金额必须大于0")
		return
	}

	od.Paid = od.Paid.Add(params.Amount)
	od.PaidChannels = od.PaidChannels.Add(params.Channel)
	switch params.Channel {
	case domain.OrderPaidChannelCash:
		od.CashPaid = od.CashPaid.Add(params.Amount)
	case domain.OrderPaidChannelWechatPay:
		od.WechatPaid = od.WechatPaid.Add(params.Amount)
	case domain.OrderPaidChannelAlipay:
		od.AlipayPaid = od.AlipayPaid.Add(params.Amount)
	case domain.OrderPaidChannelPoint:
		od.PointsPaid = od.PointsPaid.Add(params.Amount)
	case domain.OrderPaidChannelPointWallet:
		od.PointsWalletPaid = od.PointsWalletPaid.Add(params.Amount)
	}

	od.Status = lo.Ternary(od.Paid.GreaterThanOrEqual(od.RealPrice), domain.OrderStatusPaid, domain.OrderStatusPartPaid)
	od.LastPaidAt = lo.ToPtr(time.Now())
	od.MemberName = lo.Ternary(params.MemberName != "", params.MemberName, od.MemberName)
	od.MemberPhone = lo.Ternary(params.MemberPhone != "", params.MemberPhone, od.MemberPhone)
	od.MemberID = lo.Ternary(params.MemberID != 0, params.MemberID, od.MemberID)

	if od.Status == domain.OrderStatusPaid {
		od.FinishedAt = od.LastPaidAt
	}

	updateOrder, err = ds.OrderRepo().Update(ctx, od)
	if err != nil {
		err = fmt.Errorf("failed to update order: %w", err)
		return
	}

	baseEventParams := domain.OrderEventBaseParams{
		DataStore: ds,
		Order:     updateOrder,
		Operator:  params.Operator,
	}

	if err = s.OrderEventTrigger.FirePaid(ctx, &domain.OrderEventPaidParams{
		OrderEventBaseParams: baseEventParams,
		Amount:               params.Amount,
		Channel:              params.Channel,
		SeqNo:                params.SeqNo,
	}); err != nil {
		err = fmt.Errorf("failed to fire order paid event: %w", err)
		return
	}

	if updateOrder.Status == domain.OrderStatusPaid {
		if err = s.OrderEventTrigger.FireFinish(ctx, &baseEventParams); err != nil {
			err = fmt.Errorf("failed to fire order finish event: %w", err)
			return
		}
	}

	return
}

// 生成订单列表导出
func (s *DomainService) GenerateOrderListExport(ctx context.Context, name string, orders []*domain.Order) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderDomainService.GenerateOrderListExport")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	headers := []string{
		"订单号", "订单状态", "门店名称", "下单时间", "最后支付时间", "会员姓名", "会员电话", "订单总额", "优惠总额", "实付金额",
		"已付总额", "支付渠道", "支付金额", "商品名称", "商品数量", "商品单价",
	}

	var contents [][]string

	for _, order := range orders {
		maxRow := lo.Max([]int{len(order.Items), len(order.PaidChannels)})
		orderNo := order.No
		orderStatus := order.Status.ToString()
		storeName := order.StoreName
		createdAt := order.CreatedAt.Format(time.DateTime)
		lastPaidAt := "-"
		if order.LastPaidAt != nil {
			lastPaidAt = order.LastPaidAt.Format(time.DateTime)
		}
		memberName := order.MemberName
		memberPhone := order.MemberPhone
		totalPrice := order.TotalPrice.StringFixed(2)
		discount := order.Discount.StringFixed(2)
		realPrice := order.RealPrice.StringFixed(2)
		paid := order.Paid.StringFixed(2)

		for i := range maxRow {
			var line [16]string
			if i == 0 {
				line[0] = lo.Ternary(orderNo != "", orderNo, "-")
				line[1] = lo.Ternary(orderStatus != "", orderStatus, "-")
				line[2] = lo.Ternary(storeName != "", storeName, "-")
				line[3] = lo.Ternary(createdAt != "", createdAt, "-")
				line[4] = lo.Ternary(lastPaidAt != "", lastPaidAt, "-")
				line[5] = lo.Ternary(memberName != "", memberName, "-")
				line[6] = lo.Ternary(memberPhone != "", memberPhone, "-")
				line[7] = lo.Ternary(totalPrice != "", totalPrice, "-")
				line[8] = lo.Ternary(discount != "", discount, "-")
				line[9] = lo.Ternary(realPrice != "", realPrice, "-")
				line[10] = lo.Ternary(paid != "", paid, "-")

				if len(order.PaidChannels) == 0 {
					line[11] = "-"
					line[12] = "-"
				}
			}

			if i < len(order.PaidChannels) {
				channelName := order.PaidChannels[i].ToString()
				line[11] = lo.Ternary(channelName != "", channelName, "-")
				line[12] = order.PaidChannels[i].Paid(order).StringFixed(2)
			}
			if i < len(order.Items) {
				line[13] = lo.Ternary(order.Items[i].Name != "", order.Items[i].Name, "-")
				line[14] = order.Items[i].Quantity.StringFixed(3)
				line[15] = order.Items[i].Price.StringFixed(2)
			}
			contents = append(contents, line[:])
		}
	}

	url, err = s.ObjectStorage.ExportExcelWithBlankMerge(ctx, domain.ObjectStorageSceneOrderListExport, name, headers, contents)
	if err != nil {
		err = fmt.Errorf("failed to export order list: %w", err)
		return
	}

	return
}

// 订单列表导出
func (s *DomainService) OrderListExport(ctx context.Context, filename string, params *domain.OrderListExportParams) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderDomainService.OrderListExport")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	orders, _, err := s.DataStore.OrderRepo().GetOrdersWithItems(ctx, &params.Pager, &params.Filter, domain.NewOrderListOrderByCreatedAt(false))
	if err != nil {
		err = fmt.Errorf("failed to get orders: %w", err)
		return
	}

	name, _ := util.GetFileNameAndExt(filename)

	url, err = s.GenerateOrderListExport(ctx, name, orders)
	if err != nil {
		err = fmt.Errorf("failed to export order list: %w", err)
		return
	}

	return
}

var _ domain.DataExporter = (*OrderListExporter)(nil)

type OrderListExporter struct {
	DomainService *DomainService
}

func NewOrderListExporter(domainService *DomainService) *OrderListExporter {
	return &OrderListExporter{
		DomainService: domainService,
	}
}

func (s *OrderListExporter) NewParams() any {
	return new(domain.OrderListExportParams)
}

func (s *OrderListExporter) Export(ctx context.Context, filename string, params any) (url string, err error) {
	return s.DomainService.OrderListExport(ctx, filename, params.(*domain.OrderListExportParams))
}
