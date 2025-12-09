package domain

import (
	"time"
)

type OrderLogEvent string

const (
	OrderLogEventCreate     OrderLogEvent = "create"      // 创建
	OrderLogEventAppendItem OrderLogEvent = "append_item" // 加菜
	OrderLogEventRemoveItem OrderLogEvent = "remove_item" // 退菜
	OrderLogEventTurnTable  OrderLogEvent = "turn_table"  // 转台
	OrderLogEventTurnItem   OrderLogEvent = "turn_item"   // 转菜
	OrderLogEventPaid       OrderLogEvent = "paid"        // 支付
	OrderLogEventCancel     OrderLogEvent = "cancel"      // 取消
	OrderLogEventFinish     OrderLogEvent = "finish"      // 完成
)

func (OrderLogEvent) Values() []string {
	return []string{
		string(OrderLogEventCreate),
		string(OrderLogEventAppendItem),
		string(OrderLogEventRemoveItem),
		string(OrderLogEventTurnTable),
		string(OrderLogEventTurnItem),
		string(OrderLogEventPaid),
		string(OrderLogEventCancel),
		string(OrderLogEventFinish),
	}
}

// 订单操作日志
type OrderLog struct {
	ID           int           `json:"id"`
	OrderID      int           `json:"order_id"`      // 订单ID
	Event        OrderLogEvent `json:"event"`         // 事件 create: 创建订单, append_item: 加菜, remove_item: 退菜, turn_table: 转台, turn_item: 转菜, paid: 支付, cancel: 取消订单, finish: 订单完成
	OperatorType OperatorType  `json:"operator_type"` // 操作人类型 frontend: 前台用户, backend: 后台用户, admin: 管理员, system: 系统
	OperatorID   int           `json:"operator_id"`   // 操作人ID
	OperatorName string        `json:"operator_name"` // 操作人名称
	CreatedAt    time.Time     `json:"created_at"`    // 操作时间
}
