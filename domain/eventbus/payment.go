package eventbus

import (
	"context"
	"fmt"

	"github.com/gookit/event"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/domain/order"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.PaymentEventTrigger = (*PaymentEventBus)(nil)

const (
	paymentEventSuccess = "payment.success" // 支付成功
)

type paymentEventContext struct {
	baseEventContext
	payment  *domain.Payment
	operator any
}

type PaymentEventBus struct {
	em                 *event.Manager
	OrderDomainService *order.DomainService
}

func newPaymentEventContext(ctx context.Context, name string, base *domain.PaymentEventBaseParams) *paymentEventContext {
	c := &paymentEventContext{
		baseEventContext: baseEventContext{
			ds:  base.DataStore,
			ctx: ctx,
		},
		payment:  base.Payment,
		operator: base.Operator,
	}
	c.SetName(name)
	return c
}

func NewPaymentEventBus(em *event.Manager, orderDomainService *order.DomainService) *PaymentEventBus {
	bus := &PaymentEventBus{em: em, OrderDomainService: orderDomainService}
	em.On(paymentEventSuccess, event.ListenerFunc(bus.successHandler)) // 支付成功处理关联的业务
	return bus
}

func (bus *PaymentEventBus) FireSuccess(ctx context.Context, params *domain.PaymentEventSuccessParams) error {
	c := newPaymentEventContext(ctx, paymentEventSuccess, &params.PaymentEventBaseParams)
	c.Set("member_info", params.MemberInfo)
	return bus.em.FireEvent(c)
}

// 支付成功处理订单
func (bus PaymentEventBus) payOrder(ctx context.Context, ds domain.DataStore, payment *domain.Payment, operator *domain.FrontendUser, memberInfo domain.PaymentMemberInfo) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentEventBus.payOrder")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	logger := logging.FromContext(ctx).Named("PaymentEventBus.payOrder")

	var od *domain.Order
	od, err = ds.OrderRepo().Find(ctx, payment.BizID)
	if err != nil {
		err = fmt.Errorf("failed to get order: %w", err)
		return
	}
	logger = logger.With("order_no", od.No)
	ctx = logging.NewContext(ctx, logger)

	var orderPaych domain.OrderPaidChannel
	switch payment.Channel {
	case domain.PayChannelWechatPay:
		orderPaych = domain.OrderPaidChannelWechatPay
	case domain.PayChannelAlipay:
		orderPaych = domain.OrderPaidChannelAlipay
	case domain.PayChannelPoint:
		orderPaych = domain.OrderPaidChannelPoint
	case domain.PayChannelPointWallet:
		orderPaych = domain.OrderPaidChannelPointWallet
	default:
		return fmt.Errorf("unsupported pay channel: %s", payment.Channel)
	}

	paidParams := &order.PaidParams{
		Order:       od,
		Amount:      payment.Amount,
		Channel:     orderPaych,
		Operator:    operator,
		MemberID:    memberInfo.ID,
		MemberName:  memberInfo.Name,
		MemberPhone: memberInfo.Phone,
		SeqNo:       payment.SeqNo,
	}

	if _, err = bus.OrderDomainService.Paid(ctx, ds, paidParams); err != nil {
		err = fmt.Errorf("failed to paid order: %w", err)
		return
	}

	return
}

// 支付成功处理关联的业务
func (bus *PaymentEventBus) successHandler(evt event.Event) (err error) {
	c := evt.(*paymentEventContext)
	ctx := c.ctx

	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentEventBus.successHandler")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	payment := c.payment
	ds := c.ds
	operator := c.operator.(*domain.FrontendUser)

	logger := logging.FromContext(ctx).Named("PaymentEventBus.successHandler")
	logger = logger.With("biz_id", payment.BizID, "biz_type", payment.PayBizType)
	ctx = logging.NewContext(ctx, logger)

	switch payment.PayBizType {
	default:
		return fmt.Errorf("unsupported pay biz type: %s", payment.PayBizType)
	case domain.PayBizTypeOrder:
		return bus.payOrder(ctx, ds, payment, operator, c.Get("member_info").(domain.PaymentMemberInfo))
	}
}
