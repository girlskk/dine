package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// CreateRefundOrderReq 创建退款订单请求
type CreateRefundOrderReq struct {
	ID           uuid.UUID `json:"id" binding:"required"`            // 退款订单ID
	StoreID      uuid.UUID `json:"store_id" binding:"required"`      // 门店ID
	BusinessDate string    `json:"business_date" binding:"required"` // 营业日
	ShiftNo      string    `json:"shift_no"`                         // 班次号
	RefundNo     string    `json:"refund_no"`                        // 退款单号（可选，不传则自动生成）

	OriginOrderID uuid.UUID `json:"origin_order_id" binding:"required"` // 原订单ID

	RefundType       domain.RefundType       `json:"refund_type" binding:"required"` // 退款类型
	RefundReasonCode domain.RefundReasonCode `json:"refund_reason_code"`             // 退款原因代码
	RefundReason     string                  `json:"refund_reason"`                  // 退款原因描述

	Store   domain.OrderStore   `json:"store" binding:"required"`   // 门店信息
	Pos     domain.OrderPOS     `json:"pos" binding:"required"`     // POS终端信息
	Cashier domain.OrderCashier `json:"cashier" binding:"required"` // 收银员信息

	RefundAmount   domain.RefundAmount    `json:"refund_amount" binding:"required"` // 退款金额明细
	RefundPayments []domain.RefundPayment `json:"refund_payments"`                  // 退款支付记录

	RefundProducts []CreateRefundProductReq `json:"refund_products" binding:"required"` // 退款商品明细

	Remark string `json:"remark"` // 备注
}

// CreateRefundProductReq 创建退款商品请求
type CreateRefundProductReq struct {
	OriginOrderProductID uuid.UUID `json:"origin_order_product_id" binding:"required"` // 原订单商品明细ID
	OriginOrderItemID    string    `json:"origin_order_item_id"`                       // 原订单内明细ID

	ProductID   uuid.UUID          `json:"product_id" binding:"required"`   // 商品ID
	ProductName string             `json:"product_name" binding:"required"` // 商品名称
	ProductType domain.ProductType `json:"product_type"`                    // 商品类型
	Category    domain.Category    `json:"category"`                        // 分类信息
	ProductUnit domain.ProductUnit `json:"product_unit"`                    // 商品单位信息
	MainImage   string             `json:"main_image"`                      // 商品主图
	Description string             `json:"description"`                     // 菜品描述

	OriginQty      int             `json:"origin_qty" binding:"required"` // 原购买数量
	OriginPrice    decimal.Decimal `json:"origin_price"`                  // 原单价
	OriginSubtotal decimal.Decimal `json:"origin_subtotal"`               // 原小计
	OriginDiscount decimal.Decimal `json:"origin_discount"`               // 原优惠金额
	OriginTax      decimal.Decimal `json:"origin_tax"`                    // 原税额
	OriginTotal    decimal.Decimal `json:"origin_total"`                  // 原合计

	RefundQty      int             `json:"refund_qty" binding:"required"` // 退款数量
	RefundSubtotal decimal.Decimal `json:"refund_subtotal"`               // 退款小计
	RefundDiscount decimal.Decimal `json:"refund_discount"`               // 退款优惠分摊
	RefundTax      decimal.Decimal `json:"refund_tax"`                    // 退款税额
	RefundTotal    decimal.Decimal `json:"refund_total"`                  // 退款合计

	Groups        domain.SetMealGroups        `json:"groups"`         // 套餐组信息
	SpecRelations domain.ProductSpecRelations `json:"spec_relations"` // 规格信息
	AttrRelations domain.ProductAttrRelations `json:"attr_relations"` // 口味做法

	RefundReason string `json:"refund_reason"` // 单品退款原因
}

// UpdateRefundOrderReq 更新退款订单请求
type UpdateRefundOrderReq struct {
	RefundStatus     domain.RefundStatus     `json:"refund_status"`      // 退款状态
	RefundReasonCode domain.RefundReasonCode `json:"refund_reason_code"` // 退款原因代码
	RefundReason     string                  `json:"refund_reason"`      // 退款原因描述

	ApprovedBy     uuid.UUID `json:"approved_by"`      // 审批人ID
	ApprovedByName string    `json:"approved_by_name"` // 审批人名称
	ApprovedAt     time.Time `json:"approved_at"`      // 审批时间

	RefundedAt time.Time `json:"refunded_at"` // 退款完成时间

	RefundPayments []domain.RefundPayment `json:"refund_payments"` // 退款支付记录

	Remark string `json:"remark"` // 备注
}

// RefundOrderListReq 退款订单列表请求
type RefundOrderListReq struct {
	StoreID           uuid.UUID `form:"store_id" binding:"required"` // 门店ID
	OriginOrderID     uuid.UUID `form:"origin_order_id"`             // 原订单ID
	BusinessDateStart string    `form:"business_date_start"`         // 营业日开始
	BusinessDateEnd   string    `form:"business_date_end"`           // 营业日结束
	RefundNo          string    `form:"refund_no"`                   // 退款单号
	RefundType        string    `form:"refund_type"`                 // 退款类型
	RefundStatus      string    `form:"refund_status"`               // 退款状态

	Page int `form:"page" binding:"omitempty,min=1"` // 页码
	Size int `form:"size" binding:"omitempty,min=1"` // 每页数量
}

// RefundOrderListResp 退款订单列表响应
type RefundOrderListResp struct {
	Items []*domain.RefundOrder `json:"items"` // 退款订单列表
	Total int                   `json:"total"` // 总数
}
