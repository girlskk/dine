package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// ListOrderReq 订单列表请求
type ListOrderReq struct {
	BusinessDate  string `form:"business_date"`                                                        // 营业日
	OrderNo       string `form:"order_no"`                                                             // 订单号
	OrderType     string `form:"order_type" binding:"omitempty,oneof=SALE REFUND PARTIAL_REFUND"`      // 订单类型
	OrderStatus   string `form:"order_status" binding:"omitempty,oneof=PLACED COMPLETED CANCELLED"`    // 订单状态
	PaymentStatus string `form:"payment_status" binding:"omitempty,oneof=UNPAID PAYING PAID REFUNDED"` // 支付状态

	upagination.RequestPagination
}

// ListOrderResp 订单列表响应
type ListOrderResp struct {
	Items      []*domain.Order         `json:"items"`      // 订单列表
	Pagination *upagination.Pagination `json:"pagination"` // 分页信息
}
