package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OrderProduct 订单商品明细
type OrderProduct struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 订单关联
	OrderID     uuid.UUID `json:"order_id"`      // 所属订单ID
	OrderItemID string    `json:"order_item_id"` // 订单内明细ID
	Index       int       `json:"index"`         // 下单序号（同订单内第几次下单）

	// 商品基础信息
	ProductID   uuid.UUID   `json:"product_id"`   // 商品ID
	ProductName string      `json:"product_name"` // 商品名称
	ProductType ProductType `json:"product_type"` // 商品类型
	Category    Category    `json:"category"`     // 分类信息
	UnitID      uuid.UUID   `json:"unit_id"`      // 单位ID
	MainImage   string      `json:"main_image"`   // 商品主图
	Description string      `json:"description"`  // 菜品描述

	// 数量与金额
	Price           decimal.Decimal `json:"price"`             // 单价
	Qty             int             `json:"qty"`               // 数量
	IsGift          bool            `json:"is_gift"`           // 是否赠送
	GiftQty         int             `json:"gift_qty"`          // 赠送数量
	Subtotal        decimal.Decimal `json:"subtotal"`          // 小计
	DiscountAmount  decimal.Decimal `json:"discount_amount"`   // 优惠金额
	AmountBeforeTax decimal.Decimal `json:"amount_before_tax"` // 税前金额
	TaxRate         decimal.Decimal `json:"tax_rate"`          // 税率（百分比，如 6.00 表示 6%）
	Tax             decimal.Decimal `json:"tax"`               // 税额
	AmountAfterTax  decimal.Decimal `json:"amount_after_tax"`  // 税后金额
	Total           decimal.Decimal `json:"total"`             // 合计

	// 促销信息
	PromotionDiscount decimal.Decimal `json:"promotion_discount"` // 促销优惠金额

	// 做法金额与赠送金额
	AttrAmount decimal.Decimal `json:"attr_amount"` // 做法金额
	GiftAmount decimal.Decimal `json:"gift_amount"` // 赠送金额

	// 退菜信息
	VoidQty      int             `json:"void_qty"`      // 已退菜数量汇总
	VoidAmount   decimal.Decimal `json:"void_amount"`   // 已退菜金额汇总
	RefundReason string          `json:"refund_reason"` // 退菜原因
	RefundedBy   uuid.UUID       `json:"refunded_by"`   // 退菜操作人
	RefundedAt   time.Time       `json:"refunded_at"`   // 退菜时间

	// 其他信息
	Note string `json:"note"` // 备注

	// 套餐信息（仅套餐商品使用）
	Groups SetMealGroups `json:"groups"` // 套餐组信息

	// 规格信息
	SpecRelations ProductSpecRelations `json:"spec_relations"` // 商品规格关联信息

	// 口味做法信息
	AttrRelations ProductAttrRelations `json:"attr_relations"` // 商品口味做法关联信息
}
