package types

import (
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// CreateOrderReq 创建订单请求
type CreateOrderReq struct {
	TableID      int                  `json:"table_id" binding:"omitempty,gt=0"`      // 台桌ID
	PeopleNumber int                  `json:"people_number" binding:"omitempty,gt=0"` // 就餐人数
	Items        []CreateOrderReqItem `json:"items" binding:"required,min=1,dive"`    // 商品列表
}

type CreateOrderReqItem struct {
	ProductID int             `json:"product_id" binding:"required,gt=0"`         // 商品ID
	Quantity  decimal.Decimal `json:"quantity" binding:"d_positive"`              // 商品数量
	Price     decimal.Decimal `json:"price" binding:"d_nonnegative"`              // 商品价格
	SpecID    int             `json:"product_spec_id" binding:"omitempty,gt=0"`   // 商品规格ID
	AttrID    int             `json:"product_attr_id" binding:"omitempty,gt=0"`   // 商品属性ID
	RecipeID  int             `json:"product_recipe_id" binding:"omitempty,gt=0"` // 商品做法ID
	Remark    string          `json:"remark" binding:"omitempty,max=500"`         // 备注
}

// CreateOrderResp 创建订单响应
type CreateOrderResp struct {
	No string `json:"no"` // 订单号
}

// OrderDetailReq 订单详情请求
type OrderDetailReq struct {
	No string `json:"no" binding:"required"` // 订单号
}

// OrderListReq 订单列表请求
type OrderListReq struct {
	Status domain.OrderStatus `json:"status" binding:"omitempty,oneof=unpaid part_paid paid canceled"` // 订单状态
	Page   int                `json:"page" binding:"omitempty,gt=0"`                                   // 页码
	Size   int                `json:"size" binding:"omitempty,gt=0"`                                   // 每页数量
}

// OrderListResp 订单列表响应
type OrderListResp struct {
	Orders []*domain.Order `json:"orders"` // 订单列表
	Total  int             `json:"total"`  // 总数
}

// OrderModifyPriceReq 订单改价格请求
type OrderModifyPriceReq struct {
	No     string          `json:"no" binding:"required"`           // 订单号
	ItemID int             `json:"item_id" binding:"required,gt=0"` // 订单商品ID
	Price  decimal.Decimal `json:"price" binding:"d_nonnegative"`   // 价格
}

// OrderAppendItemsReq 订单添加商品请求
type OrderAppendItemsReq struct {
	No    string               `json:"no" binding:"required"`               // 订单号
	Items []CreateOrderReqItem `json:"items" binding:"required,min=1,dive"` // 商品列表
}

// OrderTurnTableReq 订单转台请求
type OrderTurnTableReq struct {
	No      string `json:"no" binding:"required"`            // 订单号
	TableID int    `json:"table_id" binding:"required,gt=0"` // 台桌ID
}

// OrderRemoveItemsReq 退菜请求
type OrderRemoveItemsReq struct {
	No       string          `json:"no" binding:"required"`           // 订单号
	ItemID   int             `json:"item_id" binding:"required,gt=0"` // 订单商品ID
	Quantity decimal.Decimal `json:"quantity" binding:"d_positive"`   // 商品数量
}

// OrderCancelReq 撤单请求
type OrderCancelReq struct {
	No string `json:"no" binding:"required"` // 订单号
}

// OrderDiscountReq 订单折扣请求
type OrderDiscountReq struct {
	No       string          `json:"no" binding:"required"`            // 订单号
	Discount decimal.Decimal `json:"discount" binding:"d_nonnegative"` // 折扣
}

// OrderCashPaidReq 现金支付请求
type OrderCashPaidReq struct {
	No     string          `json:"no" binding:"required"`          // 订单号
	Amount decimal.Decimal `json:"amount" binding:"d_nonnegative"` // 金额
}

type OrderScanPaidChannel int

const (
	_                                   OrderScanPaidChannel = iota
	OrderScanPaidChannelWechatAndAlipay                      // 微信/支付宝/知心话钱包
	OrderScanPaidChannelPoint                                // 积分
)

// OrderScanPaidReq 扫码支付请求
type OrderScanPaidReq struct {
	No       string               `json:"no" binding:"required"`                // 订单号
	Amount   decimal.Decimal      `json:"amount" binding:"d_positive"`          // 金额
	AuthCode string               `json:"auth_code" binding:"required"`         // 支付授权码
	Channel  OrderScanPaidChannel `json:"channel" binding:"required,oneof=1 2"` // 支付渠道：1-微信/支付宝/知心话钱包 2-积分
}

// OrderScanPaidResp 扫码支付响应
type OrderScanPaidResp struct {
	SeqNo string `json:"seq_no"` // 序列号
}
