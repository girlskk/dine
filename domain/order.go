package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OrderType 订单类型
type OrderType string

const (
	OrderTypeSale          OrderType = "SALE"           // 销售单
	OrderTypeRefund        OrderType = "REFUND"         // 退单
	OrderTypePartialRefund OrderType = "PARTIAL_REFUND" // 部分退款单
)

func (OrderType) Values() []string {
	return []string{
		string(OrderTypeSale),
		string(OrderTypeRefund),
		string(OrderTypePartialRefund),
	}
}

// DiningMode 就餐模式
type DiningMode string

const (
	DiningModeDineIn DiningMode = "DINE_IN" // 堂食
)

func (DiningMode) Values() []string {
	return []string{
		string(DiningModeDineIn),
	}
}

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPlaced    OrderStatus = "PLACED"    // 已下单
	OrderStatusCompleted OrderStatus = "COMPLETED" // 已完成
	OrderStatusCancelled OrderStatus = "CANCELLED" // 已取消
)

func (OrderStatus) Values() []string {
	return []string{
		string(OrderStatusPlaced),
		string(OrderStatusCompleted),
		string(OrderStatusCancelled),
	}
}

// PaymentStatus 支付状态
type PaymentStatus string

const (
	PaymentStatusUnpaid   PaymentStatus = "UNPAID"   // 未支付
	PaymentStatusPaying   PaymentStatus = "PAYING"   // 支付中
	PaymentStatusPaid     PaymentStatus = "PAID"     // 已支付
	PaymentStatusRefunded PaymentStatus = "REFUNDED" // 全额退款
)

func (PaymentStatus) Values() []string {
	return []string{
		string(PaymentStatusUnpaid),
		string(PaymentStatusPaying),
		string(PaymentStatusPaid),
		string(PaymentStatusRefunded),
	}
}

// Channel 下单渠道/操作来源
type Channel string

const (
	ChannelPOS Channel = "POS" // POS收银
	ChannelH5  Channel = "H5"  // 扫码点餐H5
	ChannelApp Channel = "APP" // 移动点餐APP
)

func (Channel) Values() []string {
	return []string{
		string(ChannelPOS),
		string(ChannelH5),
		string(ChannelApp),
	}
}

func (c Channel) Label() string {
	switch c {
	case ChannelPOS:
		return "POS收银"
	case ChannelH5:
		return "扫码点餐H5"
	case ChannelApp:
		return "移动点餐APP"
	default:
		return string(c)
	}
}

// FeeType 费用类型
type FeeType string

const (
	FeeTypeService   FeeType = "SERVICE"   // 服务费
	FeeTypePackaging FeeType = "PACKAGING" // 打包费
)

func (FeeType) Values() []string {
	return []string{
		string(FeeTypeService),
		string(FeeTypePackaging),
	}
}

type OrderRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Order, error)
	Create(ctx context.Context, order *Order) error
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params OrderListParams) ([]*Order, int, error)
	SalesReport(ctx context.Context, params OrderSalesReportParams) ([]*OrderSalesReportItem, int, error)
	ProductSalesSummary(ctx context.Context, params ProductSalesSummaryParams) ([]*ProductSalesSummaryItem, int, error)
	ProductSalesDetail(ctx context.Context, params ProductSalesDetailParams) ([]*ProductSalesDetailItem, int, error)
}

type OrderInteractor interface {
	Create(ctx context.Context, order *Order) error
	Get(ctx context.Context, id uuid.UUID) (*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params OrderListParams) ([]*Order, int, error)
	SalesReport(ctx context.Context, params OrderSalesReportParams) ([]*OrderSalesReportItem, int, error)
	ProductSalesSummary(ctx context.Context, params ProductSalesSummaryParams) ([]*ProductSalesSummaryItem, int, error)
	ProductSalesDetail(ctx context.Context, params ProductSalesDetailParams) ([]*ProductSalesDetailItem, int, error)
}

// OrderCashier 收银员信息
type OrderCashier struct {
	CashierID   uuid.UUID `json:"cashier_id,omitempty"`   // 收银员ID
	CashierName string    `json:"cashier_name,omitempty"` // 收银员名称
}

// OrderPOS POS 终端信息 对应device表
type OrderPOS struct {
	ID   uuid.UUID `json:"id"`   // POS 设备id
	Name string    `json:"name"` // POS 设备名称
}

// OrderStore 门店信息
type OrderStore struct {
	ID           uuid.UUID `json:"id"`
	MerchantID   uuid.UUID `json:"merchant_id"`   // 商户 ID
	StoreCode    string    `json:"store_code"`    // 门店编码(保留字段)
	MerchantName string    `json:"merchant_name"` // 商户名称
	ContactPhone string    `json:"contact_phone"` // 联系电话
	StoreName    string    `json:"store_name"`    // 门店名称
}

// OrderFee 费用明细
type OrderFee struct {
	FeeID   uuid.UUID `json:"fee_id,omitempty"`   // 费用ID
	FeeName string    `json:"fee_name,omitempty"` // 费用名称
	FeeType FeeType   `json:"fee_type,omitempty"` // 费用类型

	Amount decimal.Decimal `json:"amount"`         // 费用金额
	Meta   interface{}     `json:"meta,omitempty"` // 扩展信息
}

// OrderTaxRate 税率信息
type OrderTaxRate struct {
	TaxRateID   uuid.UUID `json:"tax_rate_id,omitempty"`   // 税率ID
	TaxRateName string    `json:"tax_rate_name,omitempty"` // 税率名称

	Rate decimal.Decimal `json:"rate"` // 税率（百分比）

	TaxableAmount decimal.Decimal `json:"taxable_amount"` // 计税金额
	TaxAmount     decimal.Decimal `json:"tax_amount"`     // 税额
	Meta          interface{}     `json:"meta,omitempty"` // 扩展信息
}

// OrderAmount 金额汇总
type OrderAmount struct {
	ItemsSubtotal   decimal.Decimal `json:"items_subtotal"`    // 商品小计
	DiscountTotal   decimal.Decimal `json:"discount_total"`    // 折扣合计
	TaxTotal        decimal.Decimal `json:"tax_total"`         // 税费合计
	ServiceFeeTotal decimal.Decimal `json:"service_fee_total"` // 服务费合计
	DeliveryFee     decimal.Decimal `json:"delivery_fee"`      // 配送费
	FeeTotal        decimal.Decimal `json:"fee_total"`         // 其他费用合计
	RoundingAmount  decimal.Decimal `json:"rounding_amount"`   // 舍入/抹零
	AmountDue       decimal.Decimal `json:"amount_due"`        // 应收
	AmountPaid      decimal.Decimal `json:"amount_paid"`       // 实收
	OverpayAmount   decimal.Decimal `json:"overpay_amount"`    // 溢收
	ChangeAmount    decimal.Decimal `json:"change_amount"`     // 找零
	AmountRefunded  decimal.Decimal `json:"amount_refunded"`   // 已退款
}

// OrderPayment 支付记录
type OrderPayment struct {
	PaymentNo     string               `json:"payment_no"`     // 支付号（第三方/外部交易号）
	PaymentMethod PaymentMethodPayType `json:"payment_method"` // 支付方式
	PaymentStatus PaymentStatus        `json:"payment_status"` // 支付状态
	PaymentAmount decimal.Decimal      `json:"payment_amount"` // 支付金额
	RefundAmount  decimal.Decimal      `json:"refund_amount"`  // 退款金额
	ChangeAmount  decimal.Decimal      `json:"change_amount"`  // 找零
	POS           OrderPOS             `json:"pos"`            // POS 终端信息
	Cashier       OrderCashier         `json:"cashier"`        // 收银员信息

	PaidAt time.Time `json:"paid_at,omitempty"` // 支付时间
}

// Order 订单
type Order struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	MerchantID uuid.UUID `json:"merchant_id"` // 品牌商ID
	StoreID    uuid.UUID `json:"store_id"`    // 门店ID

	BusinessDate string `json:"business_date"` // 营业日
	ShiftNo      string `json:"shift_no"`      // 班次号
	OrderNo      string `json:"order_no"`      // 订单号

	OrderType OrderType `json:"order_type"` // 订单类型

	PlacedAt    time.Time `json:"placed_at"`    // 下单时间
	PaidAt      time.Time `json:"paid_at"`      // 支付完成时间
	CompletedAt time.Time `json:"completed_at"` // 完成时间

	PlacedBy     uuid.UUID `json:"placed_by"`      // 下单人ID
	PlacedByName string    `json:"placed_by_name"` // 下单人名称

	DiningMode    DiningMode    `json:"dining_mode"`    // 堂食/外卖
	OrderStatus   OrderStatus   `json:"order_status"`   // 订单状态
	PaymentStatus PaymentStatus `json:"payment_status"` // 支付状态

	TableID    uuid.UUID `json:"table_id"`    // 桌位ID
	TableName  string    `json:"table_name"`  // 桌位名称
	GuestCount int       `json:"guest_count"` // 用餐人数

	Store   OrderStore   `json:"store"`   // 门店信息
	Channel Channel      `json:"channel"` // 下单渠道
	Pos     OrderPOS     `json:"pos"`     // POS 终端信息
	Cashier OrderCashier `json:"cashier"` // 收银员信息

	TaxRates []OrderTaxRate `json:"tax_rates"` // 税率
	Fees     []OrderFee     `json:"fees"`      // 费用
	Payments []OrderPayment `json:"payments"`  // 支付记录
	Amount   OrderAmount    `json:"amount"`    // 金额汇总

	Remark string `json:"remark"` // 整单备注

	OperationLogs []OrderOperationLog `json:"operation_logs"` // 操作日志

	OrderProducts []OrderProduct `json:"order_products"` // 订单商品明细
}

type OrderListParams struct {
	MerchantID uuid.UUID
	StoreID    uuid.UUID

	BusinessDateStart string
	BusinessDateEnd   string
	OrderNo           string
	OrderType         OrderType

	OrderStatus   OrderStatus
	PaymentStatus PaymentStatus

	Page int
	Size int
}

// OrderSalesReportParams 销售报表查询参数
type OrderSalesReportParams struct {
	MerchantID        uuid.UUID   // 品牌商ID
	StoreIDs          []uuid.UUID // 门店ID列表（可多选）
	BusinessDateStart string      // 营业日开始
	BusinessDateEnd   string      // 营业日结束

	Page int
	Size int
}

// OrderSalesReportItem 销售报表单条记录
type OrderSalesReportItem struct {
	BusinessDate       string          `json:"business_date"`         // 营业日
	StoreID            uuid.UUID       `json:"store_id"`              // 门店ID
	StoreName          string          `json:"store_name"`            // 门店名称
	OrderCount         int             `json:"order_count"`           // 单量
	GuestCount         int             `json:"guest_count"`           // 用餐人数
	AmountDue          decimal.Decimal `json:"amount_due"`            // 应收金额
	DiscountTotal      decimal.Decimal `json:"discount_total"`        // 优惠金额
	FeeTotal           decimal.Decimal `json:"fee_total"`             // 附加费金额
	AmountPaid         decimal.Decimal `json:"amount_paid"`           // 实收金额
	CashAmount         decimal.Decimal `json:"cash_amount"`           // 现金金额
	ThirdPartyAmount   decimal.Decimal `json:"third_party_amount"`    // 三方支付金额
	ChangeAmount       decimal.Decimal `json:"change_amount"`         // 零钱实收
	AmountPaidPerGuest decimal.Decimal `json:"amount_paid_per_guest"` // 人均实收
	RefundCount        int             `json:"refund_count"`          // 退款单数
	RefundAmount       decimal.Decimal `json:"refund_amount"`         // 退款金额
	NetAmount          decimal.Decimal `json:"net_amount"`            // 净收金额（实收-退款）
}

// ProductSalesSummaryParams 商品销售汇总查询参数
type ProductSalesSummaryParams struct {
	MerchantID        uuid.UUID   // 品牌商ID
	StoreIDs          []uuid.UUID // 门店ID列表
	BusinessDateStart string      // 营业日开始
	BusinessDateEnd   string      // 营业日结束
	OrderChannel      Channel     // 订单来源
	CategoryID        uuid.UUID   // 商品分类ID
	ProductName       string      // 商品名称（模糊搜索）
	ProductType       ProductType // 商品类型

	Page int
	Size int
}

// ProductSalesSummaryItem 商品销售汇总单条记录
type ProductSalesSummaryItem struct {
	ProductID      uuid.UUID       `json:"product_id"`      // 商品ID
	ProductName    string          `json:"product_name"`    // 商品名称
	ProductType    ProductType     `json:"product_type"`    // 商品类型
	CategoryName   string          `json:"category_name"`   // 一级分类名称
	CategoryName2  string          `json:"category_name_2"` // 二级分类名称
	SpecName       string          `json:"spec_name"`       // 规格名称
	SalesQty       int             `json:"sales_qty"`       // 销售数量
	SalesAmount    decimal.Decimal `json:"sales_amount"`    // 销售金额
	AvgPrice       decimal.Decimal `json:"avg_price"`       // 销售均价
	AmountDue      decimal.Decimal `json:"amount_due"`      // 商品应收金额
	Subtotal       decimal.Decimal `json:"subtotal"`        // 商品金额（小计）
	DiscountAmount decimal.Decimal `json:"discount_amount"` // 优惠金额
	GiftAmount     decimal.Decimal `json:"gift_amount"`     // 赠送金额
	RefundQty      int             `json:"refund_qty"`      // 退款数量
	RefundAmount   decimal.Decimal `json:"refund_amount"`   // 退款金额
	GiftQty        int             `json:"gift_qty"`        // 赠送数量
	AttrAmount     decimal.Decimal `json:"attr_amount"`     // 做法金额
}

// ProductSalesDetailParams 商品销售明细查询参数
type ProductSalesDetailParams struct {
	MerchantID        uuid.UUID   // 品牌商ID
	StoreIDs          []uuid.UUID // 门店ID列表
	BusinessDateStart string      // 营业日开始
	BusinessDateEnd   string      // 营业日结束
	OrderChannel      Channel     // 订单来源
	CategoryID        uuid.UUID   // 商品分类ID
	ProductName       string      // 商品名称（模糊搜索）
	ProductType       ProductType // 商品类型
	OrderNo           string      // 订单号

	Page int
	Size int
}

// ProductSalesDetailItem 商品销售明细单条记录
type ProductSalesDetailItem struct {
	BusinessDate   string          `json:"business_date"`   // 营业日期
	StoreName      string          `json:"store_name"`      // 门店名称
	ProductName    string          `json:"product_name"`    // 商品名称
	CategoryName   string          `json:"category_name"`   // 一级分类名称
	CategoryName2  string          `json:"category_name_2"` // 二级分类名称
	ProductType    ProductType     `json:"product_type"`    // 商品类型
	OrderNo        string          `json:"order_no"`        // 订单号
	OrderType      OrderType       `json:"order_type"`      // 订单类型
	PlacedAt       time.Time       `json:"placed_at"`       // 下单时间
	PaidAt         time.Time       `json:"paid_at"`         // 支付时间
	SalesQty       int             `json:"sales_qty"`       // 销售数量
	SalesAmount    decimal.Decimal `json:"sales_amount"`    // 销售金额
	AmountDue      decimal.Decimal `json:"amount_due"`      // 应收金额
	Subtotal       decimal.Decimal `json:"subtotal"`        // 商品金额
	DiscountAmount decimal.Decimal `json:"discount_amount"` // 优惠金额
	AttrAmount     decimal.Decimal `json:"attr_amount"`     // 做法金额
	RefundQty      int             `json:"refund_qty"`      // 退款数量
	RefundAmount   decimal.Decimal `json:"refund_amount"`   // 退款金额
}

// OrderOperationType 订单操作类型
type OrderOperationType string

const (
	OrderOperationTypePlaceOrder    OrderOperationType = "PLACE_ORDER"    // 点餐（下单）
	OrderOperationTypeGiftItem      OrderOperationType = "GIFT_ITEM"      // 点餐（赠菜）
	OrderOperationTypeCoupon        OrderOperationType = "COUPON"         // 结账（优惠券）
	OrderOperationTypeDiscount      OrderOperationType = "DISCOUNT"       // 结账（折扣）
	OrderOperationTypeCheckout      OrderOperationType = "CHECKOUT"       // 结账（支付）
	OrderOperationTypeReverseSettle OrderOperationType = "REVERSE_SETTLE" // 反结账
	OrderOperationTypeRefund        OrderOperationType = "REFUND"         // 退款
	OrderOperationTypeRefundReview  OrderOperationType = "REFUND_REVIEW"  // 退款审核
)

func (OrderOperationType) Values() []string {
	return []string{
		string(OrderOperationTypePlaceOrder),
		string(OrderOperationTypeGiftItem),
		string(OrderOperationTypeCoupon),
		string(OrderOperationTypeDiscount),
		string(OrderOperationTypeCheckout),
		string(OrderOperationTypeReverseSettle),
		string(OrderOperationTypeRefund),
		string(OrderOperationTypeRefundReview),
	}
}

func (t OrderOperationType) Label() string {
	switch t {
	case OrderOperationTypePlaceOrder, OrderOperationTypeGiftItem:
		return "点餐"
	case OrderOperationTypeCoupon, OrderOperationTypeDiscount, OrderOperationTypeCheckout:
		return "结账"
	case OrderOperationTypeReverseSettle:
		return "反结账"
	case OrderOperationTypeRefund:
		return "退款"
	case OrderOperationTypeRefundReview:
		return "退款审核（P1）"
	default:
		return string(t)
	}
}

// OrderOperationLog 订单操作日志
type OrderOperationLog struct {
	OperatedAt    time.Time              `json:"operated_at"`    // 操作时间
	Source        Channel                `json:"source"`         // 操作来源
	OperatorID    uuid.UUID              `json:"operator_id"`    // 操作人ID
	OperatorName  string                 `json:"operator_name"`  // 操作人名称
	OperationType OrderOperationType     `json:"operation_type"` // 操作类型
	Content       map[string]interface{} `json:"content"`        // 操作内容
}

// PlaceOrderItem 点餐商品项
type PlaceOrderItem struct {
	ProductID   uuid.UUID `json:"product_id"`           // 商品ID
	ProductName string    `json:"product_name"`         // 商品名称
	Qty         int       `json:"qty"`                  // 数量
	SpecName    string    `json:"spec_name,omitempty"`  // 规格名称
	AttrNames   []string  `json:"attr_names,omitempty"` // 做法名称列表
}

// PlaceOrderContent 点餐操作内容
type PlaceOrderContent struct {
	Items []PlaceOrderItem `json:"items"` // 商品列表
}

// GiftItem 赠菜商品项
type GiftItem struct {
	ProductID   uuid.UUID `json:"product_id"`   // 商品ID
	ProductName string    `json:"product_name"` // 商品名称
	Qty         int       `json:"qty"`          // 数量
}

// GiftItemContent 赠菜操作内容
type GiftItemContent struct {
	Items  []GiftItem `json:"items"`  // 赠菜商品列表
	Reason string     `json:"reason"` // 赠菜原因
}

// CouponContent 优惠券操作内容
type CouponContent struct {
	CouponID   uuid.UUID       `json:"coupon_id"`   // 优惠券ID
	CouponName string          `json:"coupon_name"` // 优惠券名称
	Amount     decimal.Decimal `json:"amount"`      // 优惠金额
}

// DiscountContent 折扣操作内容
type DiscountContent struct {
	DiscountRate   decimal.Decimal `json:"discount_rate"`   // 折扣率（如 0.95 表示 9.5 折）
	DiscountAmount decimal.Decimal `json:"discount_amount"` // 折扣金额
	Reason         string          `json:"reason"`          // 折扣原因
}

// PaymentDetail 支付明细
type PaymentDetail struct {
	PaymentMethod string          `json:"payment_method"` // 支付方式名称
	Amount        decimal.Decimal `json:"amount"`         // 支付金额
	Currency      string          `json:"currency"`       // 币种（如 CNY、USD）
}

// CheckoutContent 结账操作内容
type CheckoutContent struct {
	Payments []PaymentDetail `json:"payments"` // 支付明细列表
}

// ReverseSettleContent 反结账操作内容
type ReverseSettleContent struct {
	Reason string `json:"reason"` // 反结账原因
}

// RefundDetail 退款明细
type RefundDetail struct {
	PaymentMethod string          `json:"payment_method"` // 支付方式名称
	Amount        decimal.Decimal `json:"amount"`         // 退款金额
	Currency      string          `json:"currency"`       // 币种（如 CNY、USD）
}

// RefundContent 退款操作内容
type RefundContent struct {
	Refunds []RefundDetail `json:"refunds"` // 退款明细列表
}

// RefundReviewContent 退款审核操作内容
type RefundReviewContent struct {
	ReviewerID   uuid.UUID `json:"reviewer_id"`   // 审核人ID
	ReviewerName string    `json:"reviewer_name"` // 审核人名称
	Approved     bool      `json:"approved"`      // 是否通过
	Remark       string    `json:"remark"`        // 审核备注
}
