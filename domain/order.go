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

// Channel 下单渠道
type Channel string

const (
	ChannelPOS Channel = "POS" // POS终端
)

func (Channel) Values() []string {
	return []string{
		string(ChannelPOS),
	}
}

func (c Channel) GetName() string {
	switch c {
	case ChannelPOS:
		return "POS"
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

// PaymentMethod 支付方式/结算分类
type PaymentMethod string

const (
	PaymentMethodCash          PaymentMethod = "CASH"           // 现金
	PaymentMethodOnlinePayment PaymentMethod = "ONLINE_PAYMENT" // 在线支付
	PaymentMethodMemberCard    PaymentMethod = "MEMBER_CARD"    // 会员卡
	PaymentMethodCustomCoupon  PaymentMethod = "CUSTOM_COUPON"  // 自定义券
	PaymentMethodPartnerCoupon PaymentMethod = "PARTNER_COUPON" // 三方合作券
	PaymentMethodBankCard      PaymentMethod = "BANK_CARD"      // 银行卡
)

func (PaymentMethod) Values() []string {
	return []string{
		string(PaymentMethodCash),
		string(PaymentMethodOnlinePayment),
		string(PaymentMethodMemberCard),
		string(PaymentMethodCustomCoupon),
		string(PaymentMethodPartnerCoupon),
		string(PaymentMethodBankCard),
	}
}

type OrderRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Order, error)
	Create(ctx context.Context, order *Order) error
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params OrderListParams) ([]*Order, int, error)
}

type OrderInteractor interface {
	Create(ctx context.Context, order *Order) error
	Get(ctx context.Context, id uuid.UUID) (*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params OrderListParams) ([]*Order, int, error)
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
	PaymentNo     string          `json:"payment_no"`     // 支付号（第三方/外部交易号）
	PaymentMethod PaymentMethod   `json:"payment_method"` // 支付方式
	PaymentStatus PaymentStatus   `json:"payment_status"` // 支付状态
	PaymentAmount decimal.Decimal `json:"payment_amount"` // 支付金额
	RefundAmount  decimal.Decimal `json:"refund_amount"`  // 退款金额

	POS     OrderPOS     `json:"pos"`     // POS 终端信息
	Cashier OrderCashier `json:"cashier"` // 收银员信息

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

	OrderProducts []OrderProduct `json:"order_products"` // 订单商品明细
}

type OrderListParams struct {
	MerchantID uuid.UUID
	StoreID    uuid.UUID

	BusinessDate string
	OrderNo      string
	OrderType    OrderType

	OrderStatus   OrderStatus
	PaymentStatus PaymentStatus

	Page int
	Size int
}
