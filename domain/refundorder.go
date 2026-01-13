package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// RefundType 退款类型
type RefundType string

const (
	RefundTypeFull    RefundType = "FULL"    // 全额退款
	RefundTypePartial RefundType = "PARTIAL" // 部分退款
)

func (RefundType) Values() []string {
	return []string{
		string(RefundTypeFull),
		string(RefundTypePartial),
	}
}

// RefundStatus 退款状态
type RefundStatus string

const (
	RefundStatusPending    RefundStatus = "PENDING"    // 待处理
	RefundStatusProcessing RefundStatus = "PROCESSING" // 处理中
	RefundStatusCompleted  RefundStatus = "COMPLETED"  // 已完成
	RefundStatusFailed     RefundStatus = "FAILED"     // 退款失败
	RefundStatusCancelled  RefundStatus = "CANCELLED"  // 已取消
)

func (RefundStatus) Values() []string {
	return []string{
		string(RefundStatusPending),
		string(RefundStatusProcessing),
		string(RefundStatusCompleted),
		string(RefundStatusFailed),
		string(RefundStatusCancelled),
	}
}

// RefundPaymentStatus 退款支付状态
type RefundPaymentStatus string

const (
	RefundPaymentStatusPending    RefundPaymentStatus = "PENDING"    // 待退款
	RefundPaymentStatusProcessing RefundPaymentStatus = "PROCESSING" // 退款中
	RefundPaymentStatusSuccess    RefundPaymentStatus = "SUCCESS"    // 退款成功
	RefundPaymentStatusFailed     RefundPaymentStatus = "FAILED"     // 退款失败
)

func (RefundPaymentStatus) Values() []string {
	return []string{
		string(RefundPaymentStatusPending),
		string(RefundPaymentStatusProcessing),
		string(RefundPaymentStatusSuccess),
		string(RefundPaymentStatusFailed),
	}
}

// RefundChannel 退款渠道
type RefundChannel string

const (
	RefundChannelOriginal RefundChannel = "ORIGINAL" // 原路退回
	RefundChannelCash     RefundChannel = "CASH"     // 现金退款
	RefundChannelBalance  RefundChannel = "BALANCE"  // 余额退款
)

func (RefundChannel) Values() []string {
	return []string{
		string(RefundChannelOriginal),
		string(RefundChannelCash),
		string(RefundChannelBalance),
	}
}

// RefundReasonCode 退款原因代码
type RefundReasonCode string

const (
	RefundReasonCustomerRequest RefundReasonCode = "CUSTOMER_REQUEST" // 顾客要求
	RefundReasonQualityIssue    RefundReasonCode = "QUALITY_ISSUE"    // 质量问题
	RefundReasonWrongOrder      RefundReasonCode = "WRONG_ORDER"      // 下错单
	RefundReasonOutOfStock      RefundReasonCode = "OUT_OF_STOCK"     // 缺货
	RefundReasonServiceIssue    RefundReasonCode = "SERVICE_ISSUE"    // 服务问题
	RefundReasonOther           RefundReasonCode = "OTHER"            // 其他
)

func (RefundReasonCode) Values() []string {
	return []string{
		string(RefundReasonCustomerRequest),
		string(RefundReasonQualityIssue),
		string(RefundReasonWrongOrder),
		string(RefundReasonOutOfStock),
		string(RefundReasonServiceIssue),
		string(RefundReasonOther),
	}
}

// RefundAmount 退款金额明细
type RefundAmount struct {
	ItemsSubtotal   decimal.Decimal `json:"items_subtotal"`    // 商品退款小计
	DiscountTotal   decimal.Decimal `json:"discount_total"`    // 优惠退款分摊
	TaxTotal        decimal.Decimal `json:"tax_total"`         // 税费退款
	ServiceFeeTotal decimal.Decimal `json:"service_fee_total"` // 服务费退款
	FeeTotal        decimal.Decimal `json:"fee_total"`         // 其他费用退款
	RefundTotal     decimal.Decimal `json:"refund_total"`      // 退款总额
}

// RefundPayment 退款支付记录
type RefundPayment struct {
	RefundPaymentNo    string               `json:"refund_payment_no"`        // 退款流水号
	OriginPaymentNo    string               `json:"origin_payment_no"`        // 原支付流水号
	PaymentMethod      PaymentMethodPayType `json:"payment_method"`           // 支付方式
	RefundAmount       decimal.Decimal      `json:"refund_amount"`            // 退款金额
	RefundStatus       RefundPaymentStatus  `json:"refund_status"`            // 退款状态
	RefundChannel      RefundChannel        `json:"refund_channel"`           // 退款渠道
	ThirdPartyRefundNo string               `json:"third_party_refund_no"`    // 第三方退款单号
	RefundedAt         time.Time            `json:"refunded_at,omitempty"`    // 退款完成时间
	FailureReason      string               `json:"failure_reason,omitempty"` // 失败原因

	POS     OrderPOS     `json:"pos"`     // POS 终端信息
	Cashier OrderCashier `json:"cashier"` // 收银员信息
}

// RefundOrderRepository 退款订单仓储接口
type RefundOrderRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*RefundOrder, error)
	Create(ctx context.Context, refundOrder *RefundOrder) error
	Update(ctx context.Context, refundOrder *RefundOrder) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params RefundOrderListParams) ([]*RefundOrder, int, error)
	FindByOriginOrderID(ctx context.Context, originOrderID uuid.UUID) ([]*RefundOrder, error)
}

// RefundOrderInteractor 退款订单用例接口
type RefundOrderInteractor interface {
	Create(ctx context.Context, refundOrder *RefundOrder) error
	Get(ctx context.Context, id uuid.UUID) (*RefundOrder, error)
	Update(ctx context.Context, refundOrder *RefundOrder) error
	Cancel(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params RefundOrderListParams) ([]*RefundOrder, int, error)
}

// RefundOrder 退款订单
type RefundOrder struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 租户信息
	MerchantID uuid.UUID `json:"merchant_id"` // 品牌商ID
	StoreID    uuid.UUID `json:"store_id"`    // 门店ID

	// 营业信息
	BusinessDate string `json:"business_date"` // 营业日
	ShiftNo      string `json:"shift_no"`      // 班次号
	RefundNo     string `json:"refund_no"`     // 退款单号

	// 原订单关联
	OriginOrderID    uuid.UUID       `json:"origin_order_id"`    // 原订单ID
	OriginOrderNo    string          `json:"origin_order_no"`    // 原订单号
	OriginPaidAt     time.Time       `json:"origin_paid_at"`     // 原订单支付时间
	OriginAmountPaid decimal.Decimal `json:"origin_amount_paid"` // 原订单实付金额

	// 退款类型与状态
	RefundType   RefundType   `json:"refund_type"`   // 退款类型
	RefundStatus RefundStatus `json:"refund_status"` // 退款状态

	// 退款原因
	RefundReasonCode RefundReasonCode `json:"refund_reason_code"` // 退款原因代码
	RefundReason     string           `json:"refund_reason"`      // 退款原因描述

	// 操作人信息
	RefundedBy     uuid.UUID `json:"refunded_by"`      // 退款操作人ID
	RefundedByName string    `json:"refunded_by_name"` // 退款操作人名称
	ApprovedBy     uuid.UUID `json:"approved_by"`      // 审批人ID
	ApprovedByName string    `json:"approved_by_name"` // 审批人名称
	ApprovedAt     time.Time `json:"approved_at"`      // 审批时间

	// 时间节点
	RefundedAt time.Time `json:"refunded_at"` // 退款完成时间

	// 终端信息
	Store   OrderStore   `json:"store"`   // 门店信息
	Channel Channel      `json:"channel"` // 退款渠道
	Pos     OrderPOS     `json:"pos"`     // POS终端信息
	Cashier OrderCashier `json:"cashier"` // 收银员信息

	// 金额与支付
	RefundAmount   RefundAmount    `json:"refund_amount"`   // 退款金额明细
	RefundPayments []RefundPayment `json:"refund_payments"` // 退款支付记录

	// 备注
	Remark string `json:"remark"` // 备注

	// 退款商品明细
	RefundProducts []RefundOrderProduct `json:"refund_products"` // 退款商品明细
}

// RefundOrderProduct 退款商品明细
type RefundOrderProduct struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 退款单关联
	RefundOrderID uuid.UUID `json:"refund_order_id"` // 退款单ID

	// 原订单商品关联
	OriginOrderProductID uuid.UUID `json:"origin_order_product_id"` // 原订单商品明细ID
	OriginOrderItemID    string    `json:"origin_order_item_id"`    // 原订单内明细ID

	// 商品信息快照
	ProductID   uuid.UUID   `json:"product_id"`   // 商品ID
	ProductName string      `json:"product_name"` // 商品名称
	ProductType ProductType `json:"product_type"` // 商品类型
	Category    Category    `json:"category"`     // 分类信息
	ProductUnit ProductUnit `json:"product_unit"` // 商品单位信息
	MainImage   string      `json:"main_image"`   // 商品主图
	Description string      `json:"description"`  // 菜品描述

	// 原订单数量与金额
	OriginQty      int             `json:"origin_qty"`      // 原购买数量
	OriginPrice    decimal.Decimal `json:"origin_price"`    // 原单价
	OriginSubtotal decimal.Decimal `json:"origin_subtotal"` // 原小计
	OriginDiscount decimal.Decimal `json:"origin_discount"` // 原优惠金额
	OriginTax      decimal.Decimal `json:"origin_tax"`      // 原税额
	OriginTotal    decimal.Decimal `json:"origin_total"`    // 原合计

	// 退款数量与金额
	RefundQty      int             `json:"refund_qty"`      // 退款数量
	RefundSubtotal decimal.Decimal `json:"refund_subtotal"` // 退款小计
	RefundDiscount decimal.Decimal `json:"refund_discount"` // 退款优惠分摊
	RefundTax      decimal.Decimal `json:"refund_tax"`      // 退款税额
	RefundTotal    decimal.Decimal `json:"refund_total"`    // 退款合计

	// 规格/口味/套餐快照
	Groups        SetMealGroups        `json:"groups"`         // 套餐组信息
	SpecRelations ProductSpecRelations `json:"spec_relations"` // 规格信息
	AttrRelations ProductAttrRelations `json:"attr_relations"` // 口味做法

	// 退款原因
	RefundReason string `json:"refund_reason"` // 单品退款原因
}

// RefundOrderListParams 退款订单列表查询参数
type RefundOrderListParams struct {
	MerchantID        uuid.UUID
	StoreID           uuid.UUID
	OriginOrderID     uuid.UUID
	BusinessDateStart string
	BusinessDateEnd   string
	RefundNo          string
	RefundType        RefundType
	RefundStatus      RefundStatus

	Page int
	Size int
}
