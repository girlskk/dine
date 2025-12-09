package types

import "gitlab.jiguang.dev/pos-dine/dine/domain"

type CreateOrderReq struct {
	TableID      int `json:"table_id" binding:"required,gt=0"`      // 台桌ID
	PeopleNumber int `json:"people_number" binding:"required,gt=0"` // 就餐人数
}

// CreateOrderResp 创建订单响应
type CreateOrderResp struct {
	No string `json:"no"` // 订单号
}

type OrderAppendItemsReq struct {
	No      string `json:"no" binding:"required,max=50"`     // 订单号
	TableID int    `json:"table_id" binding:"required,gt=0"` // 台桌ID
}

type OrderListResp struct {
	Orders []*domain.Order `json:"orders"` // 订单列表
	Total  int             `json:"total"`  // 总数
}

type OrderDetailReq struct {
	No string `json:"no" binding:"required"` // 订单号
}
