package eventbus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gookit/event"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.OrderEventTrigger = (*OrderEventBus)(nil)

type orderEventContext struct {
	baseEventContext
	order    *domain.Order
	operator any
}

func newOrderEventContext(ctx context.Context, name string, base *domain.OrderEventBaseParams) *orderEventContext {
	c := &orderEventContext{
		baseEventContext: baseEventContext{
			ds:  base.DataStore,
			ctx: ctx,
		},
		order:    base.Order,
		operator: base.Operator,
	}
	c.SetName(name)
	return c
}

const (
	orderEventCreate      = "order.create"       // 创建
	orderEventModifyPrice = "order.modify_price" // 改价
	orderEventAppendItem  = "order.append_item"  // 添加商品
	orderEventRemoveItem  = "order.remove_item"  // 删除商品
	orderEventTurnTable   = "order.turn_table"   // 转台
	orderEventPaid        = "order.paid"         // 支付
	orderEventCancel      = "order.cancel"       // 取消
	orderEventDiscount    = "order.discount"     // 折扣
	orderEventFinish      = "order.finish"       // 完成
)

// Redis Stream 相关常量
const (
	// 订单事件流
	OrderEventStreamKey = "event:order"
)

type OrderEventBus struct {
	em  *event.Manager
	rdb redis.UniversalClient
}

func NewOrderEventBus(em *event.Manager, rdb redis.UniversalClient) *OrderEventBus {
	bus := &OrderEventBus{em: em, rdb: rdb}
	em.On("order.*", event.ListenerFunc(bus.logHandler))                    // 记录订单日志
	em.On(orderEventCreate, event.ListenerFunc(bus.occupyTableHandler))     // 创建订单时占用台桌
	em.On(orderEventCancel, event.ListenerFunc(bus.releaseTableHandler))    // 取消订单时释放台桌
	em.On(orderEventTurnTable, event.ListenerFunc(bus.releaseTableHandler)) // 转台时释放老台桌
	em.On(orderEventTurnTable, event.ListenerFunc(bus.occupyTableHandler))  // 转台时占用新台桌
	em.On(orderEventPaid, event.ListenerFunc(bus.financeLogHandler))        // 支付时记录财务日志
	em.On(orderEventFinish, event.ListenerFunc(bus.releaseTableHandler))    // 完成时释放台桌

	// 添加 Redis Stream 通知处理器
	em.On(orderEventCreate, event.ListenerFunc(bus.redisStreamNotifyHandler))     // 创建订单时发送通知
	em.On(orderEventAppendItem, event.ListenerFunc(bus.redisStreamNotifyHandler)) // 添加商品时发送通知
	em.On(orderEventRemoveItem, event.ListenerFunc(bus.redisStreamNotifyHandler)) // 删除商品时发送通知
	em.On(orderEventFinish, event.ListenerFunc(bus.redisStreamNotifyHandler))     // 支付完成时发送通知
	return bus
}

func (bus *OrderEventBus) FireCreateOrder(ctx context.Context, params *domain.OrderEventBaseParams) error {
	c := newOrderEventContext(ctx, orderEventCreate, params)
	c.Set("occupy_table_id", params.Order.TableID)
	return bus.em.FireEvent(c)
}

func (bus *OrderEventBus) FireModifyPrice(ctx context.Context, params *domain.OrderEventBaseParams) error {
	c := newOrderEventContext(ctx, orderEventModifyPrice, params)
	return bus.em.FireEvent(c)
}

func (bus *OrderEventBus) FireAppendItem(ctx context.Context, params *domain.OrderEventBaseParams) error {
	c := newOrderEventContext(ctx, orderEventAppendItem, params)
	c.Set("operated_items", params.OperatedItems)
	return bus.em.FireEvent(c)
}

func (bus *OrderEventBus) FireRemoveItem(ctx context.Context, params *domain.OrderEventBaseParams) error {
	c := newOrderEventContext(ctx, orderEventRemoveItem, params)
	c.Set("operated_items", params.OperatedItems)
	return bus.em.FireEvent(c)
}

func (bus *OrderEventBus) FireTurnTable(ctx context.Context, params *domain.OrderEventTurnTableParams) error {
	c := newOrderEventContext(ctx, orderEventTurnTable, &params.OrderEventBaseParams)
	c.Set("release_table_id", params.OldTableID)
	c.Set("occupy_table_id", params.Order.TableID)
	return bus.em.FireEvent(c)
}

func (bus *OrderEventBus) FirePaid(ctx context.Context, params *domain.OrderEventPaidParams) error {
	c := newOrderEventContext(ctx, orderEventPaid, &params.OrderEventBaseParams)
	c.Set("amount", params.Amount)
	c.Set("channel", params.Channel)
	c.Set("seq_no", params.SeqNo)
	c.Set("log_type", domain.OrderFinanceLogTypePaid)
	return bus.em.FireEvent(c)
}

func (bus *OrderEventBus) FireCancel(ctx context.Context, params *domain.OrderEventBaseParams) error {
	c := newOrderEventContext(ctx, orderEventCancel, params)
	c.Set("release_table_id", params.Order.TableID)
	return bus.em.FireEvent(c)
}

func (bus *OrderEventBus) FireDiscount(ctx context.Context, params *domain.OrderEventBaseParams) error {
	c := newOrderEventContext(ctx, orderEventDiscount, params)
	return bus.em.FireEvent(c)
}

func (bus *OrderEventBus) FireFinish(ctx context.Context, params *domain.OrderEventBaseParams) error {
	c := newOrderEventContext(ctx, orderEventFinish, params)
	c.Set("release_table_id", params.Order.TableID)
	return bus.em.FireEvent(c)
}

// 记录订单日志
func (bus *OrderEventBus) logHandler(evt event.Event) (err error) {
	c := evt.(*orderEventContext)
	eName := evt.Name()
	ctx := c.ctx

	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderEventBus.logHandler")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	var logEvent domain.OrderLogEvent
	switch eName {
	default:
		return nil
	case orderEventCreate:
		logEvent = domain.OrderLogEventCreate
	case orderEventAppendItem:
		logEvent = domain.OrderLogEventAppendItem
	case orderEventRemoveItem:
		logEvent = domain.OrderLogEventRemoveItem
	case orderEventTurnTable:
		logEvent = domain.OrderLogEventTurnTable
	// case orderEventTurnItem:
	// 	logEvent = domain.OrderLogEventTurnItem
	case orderEventPaid:
		logEvent = domain.OrderLogEventPaid
	case orderEventCancel:
		logEvent = domain.OrderLogEventCancel
	case orderEventFinish:
		logEvent = domain.OrderLogEventFinish
	}

	operatorInfo := domain.ExtractOperatorInfo(c.operator)

	log := &domain.OrderLog{
		OrderID:      c.order.ID,
		Event:        logEvent,
		OperatorType: operatorInfo.Type,
		OperatorID:   operatorInfo.ID,
		OperatorName: operatorInfo.Name,
	}

	if _, err := c.ds.OrderRepo().CreateLog(ctx, log); err != nil {
		return fmt.Errorf("failed to create order log: %w", err)
	}

	return
}

// 记录订单财务日志
func (bus *OrderEventBus) financeLogHandler(evt event.Event) (err error) {
	c := evt.(*orderEventContext)
	ctx := c.ctx

	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderEventBus.financeLogHandler")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	amount := c.Get("amount").(decimal.Decimal)
	channel := c.Get("channel").(domain.OrderPaidChannel)
	seqNo := c.Get("seq_no").(string)
	logType := c.Get("log_type").(domain.OrderFinanceLogType)

	creator := domain.ExtractOperatorInfo(c.operator)

	log := &domain.OrderFinanceLog{
		OrderID:     c.order.ID,
		Amount:      amount,
		Type:        logType,
		Channel:     channel,
		SeqNo:       seqNo,
		CreatorType: creator.Type,
		CreatorID:   creator.ID,
		CreatorName: creator.Name,
	}

	if _, err := c.ds.OrderRepo().CreateFinanceLog(ctx, log); err != nil {
		return fmt.Errorf("failed to create order finance log: %w", err)
	}

	return
}

// 占用台桌
func (bus *OrderEventBus) occupyTableHandler(evt event.Event) (err error) {
	c := evt.(*orderEventContext)
	ctx := c.ctx

	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderEventBus.occupyTableHandler")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	tableID := c.Get("occupy_table_id").(int)
	if tableID == 0 {
		return nil
	}

	ok, err := c.ds.TableRepo().UpdateOrderIDAndStatusFrom(ctx, tableID, c.order.ID, domain.TableStatusFree, domain.TableStatusOccupied)
	if err != nil {
		return fmt.Errorf("failed to update table status: %w", err)
	}
	if !ok {
		return domain.ParamsErrorf("台桌已被占用")
	}

	return
}

// 释放台桌
func (bus *OrderEventBus) releaseTableHandler(evt event.Event) (err error) {
	c := evt.(*orderEventContext)
	ctx := c.ctx

	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderEventBus.releaseTableHandler")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	tableID := c.Get("release_table_id").(int)

	if tableID == 0 {
		return nil
	}

	ok, err := c.ds.TableRepo().UpdateOrderIDAndStatusFrom(ctx, tableID, 0, domain.TableStatusOccupied, domain.TableStatusFree)
	if err != nil {
		return fmt.Errorf("failed to update table status: %w", err)
	}
	if !ok {
		return domain.ParamsErrorf("台桌状态异常，无法释放")
	}

	return
}

func (bus *OrderEventBus) redisStreamNotifyHandler(evt event.Event) (err error) {
	c := evt.(*orderEventContext)
	ctx := c.ctx
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderEventBus.releaseTableHandler")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 如果 order 详情为空，查询订单详情
	if c.order.Items == nil {
		order, err := c.ds.OrderRepo().Find(ctx, c.order.ID)
		if err != nil {
			return fmt.Errorf("failed to get order: %w", err)
		}
		c.order = order
	}

	data := map[string]any{
		"type":           evt.Name(),
		"order":          c.order,
		"operated_items": c.Get("operated_items"),
	}
	// 序列化为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}
	// 添加到 Redis Stream
	_, err = bus.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: OrderEventStreamKey,
		ID:     "*", // 让 Redis 自动生成 ID
		Values: map[string]any{
			"data": string(jsonData),
		},
	}).Result()

	if err != nil {
		return fmt.Errorf("failed to publish event to Redis Stream: %w", err)
	}
	return
}
