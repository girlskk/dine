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
	CashierID   string `json:"cashier_id,omitempty"`   // 收银员ID
	CashierName string `json:"cashier_name,omitempty"` // 收银员名称
}

// OrderPOS POS 终端信息
type OrderPOS struct {
	PosID    string `json:"pos_id,omitempty"`    // POS ID
	PosCode  string `json:"pos_code,omitempty"`  // POS 编码
	DeviceID string `json:"device_id,omitempty"` // 设备ID
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
	FeeID   string `json:"fee_id,omitempty"`   // 费用ID
	FeeName string `json:"fee_name,omitempty"` // 费用名称
	FeeType string `json:"fee_type,omitempty"` // 费用类型（TIP/PACKAGING/SURCHARGE/OTHER 等）

	Amount decimal.Decimal `json:"amount"`         // 费用金额
	Meta   interface{}     `json:"meta,omitempty"` // 扩展信息
}

// OrderTaxRate 税率信息
type OrderTaxRate struct {
	TaxRateID   string `json:"tax_rate_id,omitempty"`   // 税率ID
	TaxRateName string `json:"tax_rate_name,omitempty"` // 税率名称

	Rate decimal.Decimal `json:"rate"` // 税率（百分比）

	TaxableAmount decimal.Decimal `json:"taxable_amount"` // 计税金额
	TaxAmount     decimal.Decimal `json:"tax_amount"`     // 税额
	Meta          interface{}     `json:"meta,omitempty"` // 扩展信息
}

// OrderRefund 退款单信息（order_type=REFUND/PARTIAL_REFUND）
type OrderRefund struct {
	OriginOrderID string `json:"origin_order_id,omitempty"` // 原正单订单ID
	OriginOrderNo string `json:"origin_order_no,omitempty"` // 原正单订单号
	Reason        string `json:"reason,omitempty"`          // 退款说明
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
	ChangeAmount    decimal.Decimal `json:"change_amount"`     // 找零
	AmountRefunded  decimal.Decimal `json:"amount_refunded"`   // 已退款
}

// OrderPayment 支付记录
type OrderPayment struct {
	PaymentNo     string          `json:"payment_no"`     // 支付号（第三方/外部交易号）
	PaymentMethod string          `json:"payment_method"` // 支付方式
	PaymentAmount decimal.Decimal `json:"payment_amount"` // 支付金额

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

	OrderType OrderType   `json:"order_type"` // 订单类型
	Refund    OrderRefund `json:"refund"`     // 退款单信息

	PlacedAt    time.Time `json:"placed_at"`    // 下单时间
	PaidAt      time.Time `json:"paid_at"`      // 支付完成时间
	CompletedAt time.Time `json:"completed_at"` // 完成时间

	PlacedBy string `json:"placed_by"` // 下单人

	DiningMode    DiningMode    `json:"dining_mode"`    // 堂食/外卖
	OrderStatus   OrderStatus   `json:"order_status"`   // 订单状态
	PaymentStatus PaymentStatus `json:"payment_status"` // 支付状态

	TableID    string `json:"table_id"`    // 桌位ID
	TableName  string `json:"table_name"`  // 桌位名称
	GuestCount int    `json:"guest_count"` // 用餐人数

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
