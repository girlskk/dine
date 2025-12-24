package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
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
	DiningModeDineIn   DiningMode = "DINE_IN"  // 堂食
	DiningModeTakeaway DiningMode = "TAKEAWAY" // 外卖（自取/配送）
)

func (DiningMode) Values() []string {
	return []string{
		string(DiningModeDineIn),
		string(DiningModeTakeaway),
	}
}

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusDraft      OrderStatus = "DRAFT"       // 草稿/购物车
	OrderStatusPlaced     OrderStatus = "PLACED"      // 已下单
	OrderStatusInProgress OrderStatus = "IN_PROGRESS" // 制作中
	OrderStatusReady      OrderStatus = "READY"       // 可取餐
	OrderStatusCompleted  OrderStatus = "COMPLETED"   // 已完成
	OrderStatusCancelled  OrderStatus = "CANCELLED"   // 已取消
	OrderStatusVoided     OrderStatus = "VOIDED"      // 已作废
	OrderStatusMerged     OrderStatus = "MERGED"      // 已合并
)

func (OrderStatus) Values() []string {
	return []string{
		string(OrderStatusDraft),
		string(OrderStatusPlaced),
		string(OrderStatusInProgress),
		string(OrderStatusReady),
		string(OrderStatusCompleted),
		string(OrderStatusCancelled),
		string(OrderStatusVoided),
		string(OrderStatusMerged),
	}
}

// PaymentStatus 支付状态
type PaymentStatus string

const (
	PaymentStatusUnpaid            PaymentStatus = "UNPAID"             // 未支付
	PaymentStatusPaying            PaymentStatus = "PAYING"             // 支付中
	PaymentStatusPartiallyPaid     PaymentStatus = "PARTIALLY_PAID"     // 部分支付
	PaymentStatusPaid              PaymentStatus = "PAID"               // 已支付
	PaymentStatusPartiallyRefunded PaymentStatus = "PARTIALLY_REFUNDED" // 部分退款
	PaymentStatusRefunded          PaymentStatus = "REFUNDED"           // 全额退款
)

func (PaymentStatus) Values() []string {
	return []string{
		string(PaymentStatusUnpaid),
		string(PaymentStatusPaying),
		string(PaymentStatusPartiallyPaid),
		string(PaymentStatusPaid),
		string(PaymentStatusPartiallyRefunded),
		string(PaymentStatusRefunded),
	}
}

// FulfillmentStatus 交付状态
type FulfillmentStatus string

const (
	FulfillmentStatusNone          FulfillmentStatus = "NONE"           // 无
	FulfillmentStatusInRestaurant  FulfillmentStatus = "IN_RESTAURANT"  // 店内用餐
	FulfillmentStatusServed        FulfillmentStatus = "SERVED"         // 已上齐
	FulfillmentStatusPickupPending FulfillmentStatus = "PICKUP_PENDING" // 待取餐
	FulfillmentStatusPickedUp      FulfillmentStatus = "PICKED_UP"      // 已取餐
	FulfillmentStatusDelivering    FulfillmentStatus = "DELIVERING"     // 配送中
	FulfillmentStatusDelivered     FulfillmentStatus = "DELIVERED"      // 已送达
)

func (FulfillmentStatus) Values() []string {
	return []string{
		string(FulfillmentStatusNone),
		string(FulfillmentStatusInRestaurant),
		string(FulfillmentStatusServed),
		string(FulfillmentStatusPickupPending),
		string(FulfillmentStatusPickedUp),
		string(FulfillmentStatusDelivering),
		string(FulfillmentStatusDelivered),
	}
}

// TableStatus 桌位状态
type TableStatus string

const (
	TableStatusOpened      TableStatus = "OPENED"      // 已开台
	TableStatusTransferred TableStatus = "TRANSFERRED" // 已转台
	TableStatusReleased    TableStatus = "RELEASED"    // 已释放
)

func (TableStatus) Values() []string {
	return []string{
		string(TableStatusOpened),
		string(TableStatusTransferred),
		string(TableStatusReleased),
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
	Create(ctx context.Context, order *Order) (*Order, error)
	Get(ctx context.Context, id uuid.UUID) (*Order, error)
	Update(ctx context.Context, order *Order) (*Order, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params OrderListParams) ([]*Order, int, error)
}

// OrderMember 会员信息
type OrderMember struct {
	MemberID        string `json:"member_id,omitempty"`         // 会员ID
	MemberNo        string `json:"member_no,omitempty"`         // 会员号
	MemberName      string `json:"member_name,omitempty"`       // 会员姓名/昵称
	MemberPhone     string `json:"member_phone,omitempty"`      // 会员手机号
	MemberLevelName string `json:"member_level_name,omitempty"` // 会员等级名称
}

// OrderStore 门店信息
type OrderStore struct {
	StoreID      uuid.UUID `json:"store_id"`                // 门店ID
	StoreNo      string    `json:"store_no,omitempty"`      // 门店编号
	StoreName    string    `json:"store_name,omitempty"`    // 门店名称
	StorePhone   string    `json:"store_phone,omitempty"`   // 门店电话
	StoreAddress string    `json:"store_address,omitempty"` // 门店地址
}

// OrderChannel 下单渠道信息
type OrderChannel struct {
	Code string `json:"code,omitempty"` // 渠道编码（POS/MINI_PROGRAM）
	Name string `json:"name,omitempty"` // 渠道名称
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

// OrderTakeaway 外卖信息
type OrderTakeaway struct {
	TakeawayType string `json:"takeaway_type,omitempty"` // 外卖类型（PICKUP/DELIVERY）

	ContactName  string `json:"contact_name,omitempty"`  // 联系人/收货人
	ContactPhone string `json:"contact_phone,omitempty"` // 联系电话

	PickupNo    string `json:"pickup_no,omitempty"`     // 取餐号（自取）
	PickupEtaAt string `json:"pickup_eta_at,omitempty"` // 预计取餐时间

	DeliveryAddress string `json:"delivery_address,omitempty"` // 配送地址（配送）
	DeliveryFee     int64  `json:"delivery_fee,omitempty"`     // 配送费（分）

	DeliveryPlatform   string `json:"delivery_platform,omitempty"`    // 配送平台（如 MEITUAN/ELEME/SELF）
	DeliveryOrderNo    string `json:"delivery_order_no,omitempty"`    // 平台配送单号
	DeliveryTrackingNo string `json:"delivery_tracking_no,omitempty"` // 运单号/骑手单号
	DeliveryStatus     string `json:"delivery_status,omitempty"`      // 配送状态（平台侧）

	DeliveryRiderName  string `json:"delivery_rider_name,omitempty"`  // 骑手姓名
	DeliveryRiderPhone string `json:"delivery_rider_phone,omitempty"` // 骑手电话

	DeliveryRemark string `json:"delivery_remark,omitempty"` // 配送备注
}

// OrderPromotion 促销信息
type OrderPromotion struct {
	PromotionID   string `json:"promotion_id,omitempty"`   // 促销ID
	PromotionName string `json:"promotion_name,omitempty"` // 促销名称
	PromotionType string `json:"promotion_type,omitempty"` // 促销类型

	DiscountAmount int64       `json:"discount_amount"` // 促销优惠金额（分）
	Meta           interface{} `json:"meta,omitempty"`  // 促销扩展信息
}

// OrderFee 费用明细
type OrderFee struct {
	FeeID   string `json:"fee_id,omitempty"`   // 费用ID
	FeeName string `json:"fee_name,omitempty"` // 费用名称
	FeeType string `json:"fee_type,omitempty"` // 费用类型（TIP/PACKAGING/SURCHARGE/OTHER 等）

	Amount int64       `json:"amount"`         // 费用金额（分）
	Meta   interface{} `json:"meta,omitempty"` // 扩展信息
}

// OrderCoupon 卡券信息
type OrderCoupon struct {
	CouponID   string `json:"coupon_id,omitempty"`   // 卡券ID
	CouponName string `json:"coupon_name,omitempty"` // 卡券名称
	CouponType string `json:"coupon_type,omitempty"` // 卡券类型
	CouponCode string `json:"coupon_code,omitempty"` // 卡券码

	DiscountAmount int64       `json:"discount_amount"` // 卡券优惠金额（分）
	Meta           interface{} `json:"meta,omitempty"`  // 卡券扩展信息
}

// OrderTaxRate 税率信息
type OrderTaxRate struct {
	TaxRateID   string `json:"tax_rate_id,omitempty"`   // 税率ID
	TaxRateName string `json:"tax_rate_name,omitempty"` // 税率名称

	Rate int64 `json:"rate"` // 税率（万分比，如 600 表示 6%）

	TaxableAmount int64       `json:"taxable_amount"` // 计税金额（分）
	TaxAmount     int64       `json:"tax_amount"`     // 税额（分）
	Meta          interface{} `json:"meta,omitempty"` // 扩展信息
}

// OrderProduct 商品
type OrderProduct struct {
	OrderItemID string `json:"order_item_id"`   // 订单内明细ID
	Index       int    `json:"index,omitempty"` // 下单序号（同订单内第几次下单）

	RefundReason string    `json:"refund_reason,omitempty"` // 退菜原因
	RefundedBy   string    `json:"refunded_by,omitempty"`   // 退菜操作人
	RefundedAt   time.Time `json:"refunded_at,omitempty"`   // 退菜时间

	Promotions        []OrderPromotion `json:"promotions,omitempty"` // 促销明细（可叠加）
	PromotionDiscount int64            `json:"promotion_discount"`   // 促销优惠金额（分）

	ProductID   string `json:"product_id,omitempty"` // 商品ID
	ProductName string `json:"product_name"`         // 商品名
	SkuID       string `json:"sku_id,omitempty"`     // SKU ID
	SkuName     string `json:"sku_name,omitempty"`   // SKU 名称

	Qty int `json:"qty"` // 数量

	Price           int64 `json:"price"`             // 单价（分）
	Subtotal        int64 `json:"subtotal"`          // 小计（分）
	DiscountAmount  int64 `json:"discount_amount"`   // 优惠金额（分）
	AmountBeforeTax int64 `json:"amount_before_tax"` // 税前金额（分）
	TaxRate         int64 `json:"tax_rate"`          // 税率（万分比，如 600 表示 6%）
	Tax             int64 `json:"tax"`               // 税额（分）
	AmountAfterTax  int64 `json:"amount_after_tax"`  // 税后金额（分）
	Total           int64 `json:"total"`             // 合计（分）

	VoidQty    int   `json:"void_qty"`    // 已退菜数量汇总
	VoidAmount int64 `json:"void_amount"` // 已退菜金额汇总（分）

	Note    string      `json:"note,omitempty"`    // 备注
	Options interface{} `json:"options,omitempty"` // 做法/加料
}

// OrderRefund 退款单信息（order_type=REFUND/PARTIAL_REFUND）
type OrderRefund struct {
	OriginOrderID uuid.UUID `json:"origin_order_id,omitempty"` // 原正单订单ID
	OriginOrderNo string    `json:"origin_order_no,omitempty"` // 原正单订单号
	Reason        string    `json:"reason,omitempty"`          // 退款说明
}

// OrderAmount 金额汇总（分）
type OrderAmount struct {
	ItemsSubtotal          int64 `json:"items_subtotal"`           // 商品小计（分）
	DiscountTotal          int64 `json:"discount_total"`           // 折扣合计（分）
	PromotionDiscountTotal int64 `json:"promotion_discount_total"` // 促销优惠合计（分）
	VoucherDiscountTotal   int64 `json:"voucher_discount_total"`   // 卡券优惠合计（分）
	TaxTotal               int64 `json:"tax_total"`                // 税费合计（分）
	ServiceFeeTotal        int64 `json:"service_fee_total"`        // 服务费合计（分）
	DeliveryFee            int64 `json:"delivery_fee"`             // 配送费（分）
	FeeTotal               int64 `json:"fee_total"`                // 其他费用合计（分，小费/附加费等）
	RoundingAmount         int64 `json:"rounding_amount"`          // 舍入/抹零（分）
	AmountDue              int64 `json:"amount_due"`               // 应收（分）
	AmountPaid             int64 `json:"amount_paid"`              // 实收（分）
	ChangeAmount           int64 `json:"change_amount"`            // 找零（分）
	AmountRefunded         int64 `json:"amount_refunded"`          // 已退款（分）
}

// OrderPayment 支付记录
type OrderPayment struct {
	PaymentNo     string `json:"payment_no"`     // 支付号（第三方/外部交易号）
	PaymentMethod string `json:"payment_method"` // 支付方式
	PaymentAmount int64  `json:"payment_amount"` // 支付金额（分）

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

	OrderType     OrderType    `json:"order_type"`      // 订单类型
	OriginOrderID string       `json:"origin_order_id"` // 原正单订单ID（退款/部分退款单使用）
	Refund        *OrderRefund `json:"refund"`

	OpenedAt    *time.Time `json:"opened_at"`    // 开台时间
	PlacedAt    *time.Time `json:"placed_at"`    // 下单时间
	PaidAt      *time.Time `json:"paid_at"`      // 支付完成时间
	CompletedAt *time.Time `json:"completed_at"` // 完成时间

	OpenedBy string `json:"opened_by"` // 开台人
	PlacedBy string `json:"placed_by"` // 下单人
	PaidBy   string `json:"paid_by"`   // 收款人

	DiningMode        DiningMode        `json:"dining_mode"`        // 堂食/外卖
	OrderStatus       OrderStatus       `json:"order_status"`       // 订单状态
	PaymentStatus     PaymentStatus     `json:"payment_status"`     // 支付状态
	FulfillmentStatus FulfillmentStatus `json:"fulfillment_status"` // 交付状态
	TableStatus       TableStatus       `json:"table_status"`       // 桌位状态

	TableID       string `json:"table_id"`       // 桌位ID
	TableName     string `json:"table_name"`     // 桌位名称
	TableCapacity int    `json:"table_capacity"` // 桌位容量
	GuestCount    int    `json:"guest_count"`    // 用餐人数

	MergedToOrderID string     `json:"merged_to_order_id"` // 合并到的目标订单ID
	MergedAt        *time.Time `json:"merged_at"`          // 合并时间

	Store   *OrderStore   `json:"store"`   // 门店信息
	Channel *OrderChannel `json:"channel"` // 下单渠道
	Pos     *OrderPOS     `json:"pos"`     // POS 终端信息
	Cashier *OrderCashier `json:"cashier"` // 收银员信息

	Member   *OrderMember   `json:"member"`   // 会员信息
	Takeaway *OrderTakeaway `json:"takeaway"` // 外卖信息

	Cart            *[]OrderProduct   `json:"cart"`             // 购物车
	Products        *[]OrderProduct   `json:"products"`         // 下单商品
	Promotions      *[]OrderPromotion `json:"promotions"`       // 促销
	Coupons         *[]OrderCoupon    `json:"coupons"`          // 卡券
	TaxRates        *[]OrderTaxRate   `json:"tax_rates"`        // 税率
	Fees            *[]OrderFee       `json:"fees"`             // 费用
	Payments        *[]OrderPayment   `json:"payments"`         // 支付记录
	RefundsProducts *[]OrderProduct   `json:"refunds_products"` // 退菜记录
	Amount          *OrderAmount      `json:"amount"`           // 金额汇总
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
