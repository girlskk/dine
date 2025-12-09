package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

// OrderListReq 订单列表请求
type OrderListReq struct {
	upagination.RequestPagination
	Status            domain.OrderStatus `json:"status" binding:"omitempty,oneof=unpaid part_paid paid canceled"` // 订单状态
	HasItemName       string             `json:"has_item_name"`                                                   // 包含的商品名称
	MemberNameOrPhone string             `json:"member_name_or_phone"`                                            // 会员名称或手机号
	CreatedAtGte      util.RequestDate   `json:"created_at_gte"`                                                  // 创建时间大于等于（格式：YYYY-MM-DD）
	CreatedAtLte      util.RequestDate   `json:"created_at_lte"`                                                  // 创建时间小于等于（格式：YYYY-MM-DD）
}

// OrderListResp 订单列表响应
type OrderListResp struct {
	Orders []*domain.Order `json:"orders"` // 订单列表
	Total  int             `json:"total"`  // 总数
}

// OrderDetailReq 订单详情请求
type OrderDetailReq struct {
	No string `json:"no"` // 订单号
}

// OrderListExportReq 订单导出请求
type OrderListExportReq struct {
	Status            domain.OrderStatus `json:"status" binding:"omitempty,oneof=unpaid part_paid paid canceled"` // 订单状态
	HasItemName       string             `json:"has_item_name"`                                                   // 包含的商品名称
	MemberNameOrPhone string             `json:"member_name_or_phone"`                                            // 会员名称或手机号
	CreatedAtGte      util.RequestDate   `json:"created_at_gte"`                                                  // 创建时间大于等于（格式：YYYY-MM-DD）
	CreatedAtLte      util.RequestDate   `json:"created_at_lte"`                                                  // 创建时间小于等于（格式：YYYY-MM-DD）
}
