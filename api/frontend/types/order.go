package types

import (
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// CreateOrderReq 创建订单请求
type CreateOrderReq struct {
	ID           uuid.UUID `json:"id" binding:"required"`            // 订单ID
	BusinessDate string    `json:"business_date" binding:"required"` // 营业日
	ShiftNo      string    `json:"shift_no"`                         // 班次号
	OrderNo      string    `json:"order_no"`                         // 订单号

	OrderType string `json:"order_type" binding:"omitempty,oneof=SALE REFUND PARTIAL_REFUND"` // 订单类型

	DiningMode    string `json:"dining_mode" binding:"required,oneof=DINE_IN"`                         // 就餐模式
	OrderStatus   string `json:"order_status" binding:"omitempty,oneof=PLACED COMPLETED CANCELLED"`    // 订单状态
	PaymentStatus string `json:"payment_status" binding:"omitempty,oneof=UNPAID PAYING PAID REFUNDED"` // 支付状态
	Channel       string `json:"channel" binding:"omitempty,oneof=POS H5 APP"`                         // 下单渠道

	TableID    uuid.UUID `json:"table_id"`    // 桌位ID
	TableName  string    `json:"table_name"`  // 桌位名称
	GuestCount int       `json:"guest_count"` // 用餐人数

	PlacedAt time.Time `json:"placed_at"` // 下单时间
	PlacedBy uuid.UUID `json:"placed_by"` // 下单人

	Store   domain.OrderStore   `json:"store"`   // 门店信息
	Pos     domain.OrderPOS     `json:"pos"`     // POS终端信息
	Cashier domain.OrderCashier `json:"cashier"` // 收银员信息

	OrderProducts []domain.OrderProduct `json:"order_products"` // 订单商品明细
	TaxRates      []domain.OrderTaxRate `json:"tax_rates"`      // 税率明细
	Fees          []domain.OrderFee     `json:"fees"`           // 费用明细
	Payments      []domain.OrderPayment `json:"payments"`       // 支付记录
	Amount        domain.OrderAmount    `json:"amount"`         // 金额汇总

	Remark string `json:"remark"` // 整单备注

	OperationLogs []domain.OrderOperationLog `json:"operation_logs"` // 操作日志
}

// UpdateOrderReq 更新订单请求
type UpdateOrderReq struct {
	BusinessDate string `json:"business_date"` // 营业日
	ShiftNo      string `json:"shift_no"`      // 班次号
	OrderNo      string `json:"order_no"`      // 订单号

	OrderType string `json:"order_type" binding:"omitempty,oneof=SALE REFUND PARTIAL_REFUND"` // 订单类型

	DiningMode    string `json:"dining_mode" binding:"omitempty,oneof=DINE_IN"`                        // 就餐模式
	OrderStatus   string `json:"order_status" binding:"omitempty,oneof=PLACED COMPLETED CANCELLED"`    // 订单状态
	PaymentStatus string `json:"payment_status" binding:"omitempty,oneof=UNPAID PAYING PAID REFUNDED"` // 支付状态
	Channel       string `json:"channel" binding:"omitempty,oneof=POS H5 APP"`                         // 下单渠道

	TableID    uuid.UUID `json:"table_id"`    // 桌位ID
	TableName  string    `json:"table_name"`  // 桌位名称
	GuestCount int       `json:"guest_count"` // 用餐人数

	PlacedAt time.Time `json:"placed_at"` // 下单时间
	PaidAt   time.Time `json:"paid_at"`   // 支付完成时间
	PlacedBy uuid.UUID `json:"placed_by"` // 下单人

	Store   domain.OrderStore   `json:"store"`   // 门店信息
	Pos     domain.OrderPOS     `json:"pos"`     // POS终端信息
	Cashier domain.OrderCashier `json:"cashier"` // 收银员信息

	OrderProducts []domain.OrderProduct `json:"order_products"` // 订单商品明细
	TaxRates      []domain.OrderTaxRate `json:"tax_rates"`      // 税率明细
	Fees          []domain.OrderFee     `json:"fees"`           // 费用明细
	Payments      []domain.OrderPayment `json:"payments"`       // 支付记录
	Amount        domain.OrderAmount    `json:"amount"`         // 金额汇总

	Remark string `json:"remark"` // 整单备注

	OperationLogs []domain.OrderOperationLog `json:"operation_logs"` // 操作日志
}

// ListOrderReq 订单列表请求
type ListOrderReq struct {
	BusinessDateStart string `form:"business_date_start"`                                                  // 营业日开始
	BusinessDateEnd   string `form:"business_date_end"`                                                    // 营业日结束
	OrderNo           string `form:"order_no"`                                                             // 订单号
	OrderType         string `form:"order_type" binding:"omitempty,oneof=SALE REFUND PARTIAL_REFUND"`      // 订单类型
	OrderStatus       string `form:"order_status" binding:"omitempty,oneof=PLACED COMPLETED CANCELLED"`    // 订单状态
	PaymentStatus     string `form:"payment_status" binding:"omitempty,oneof=UNPAID PAYING PAID REFUNDED"` // 支付状态

	upagination.RequestPagination
}

// ListOrderResp 订单列表响应
type ListOrderResp struct {
	Items      []*domain.Order         `json:"items"`      // 订单列表
	Pagination *upagination.Pagination `json:"pagination"` // 分页信息
}
