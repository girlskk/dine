package types

import (
	"github.com/google/uuid"
	"time"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// Member 会员信息
type Member struct {
	MemberID        string `json:"member_id,omitempty"`         // 会员ID
	MemberNo        string `json:"member_no,omitempty"`         // 会员号
	MemberName      string `json:"member_name,omitempty"`       // 会员姓名/昵称
	MemberPhone     string `json:"member_phone,omitempty"`      // 会员手机号
	MemberLevelName string `json:"member_level_name,omitempty"` // 会员等级名称
}

// Store 门店信息
type Store struct {
	StoreID      uuid.UUID `json:"store_id"`                // 门店ID
	StoreNo      string    `json:"store_no,omitempty"`      // 门店编号
	StoreName    string    `json:"store_name,omitempty"`    // 门店名称
	StorePhone   string    `json:"store_phone,omitempty"`   // 门店电话
	StoreAddress string    `json:"store_address,omitempty"` // 门店地址
}

// Channel 下单渠道信息
type Channel struct {
	Code string `json:"code,omitempty"` // 渠道编码（POS/MINI_PROGRAM）
	Name string `json:"name,omitempty"` // 渠道名称
}

// Cashier 收银员信息
type Cashier struct {
	CashierID   string `json:"cashier_id,omitempty"`   // 收银员ID
	CashierName string `json:"cashier_name,omitempty"` // 收银员名称
}

// POS POS 终端信息
type POS struct {
	PosID    string `json:"pos_id,omitempty"`    // POS ID
	PosCode  string `json:"pos_code,omitempty"`  // POS 编码
	DeviceID string `json:"device_id,omitempty"` // 设备ID
}

// Takeaway 外卖信息
type Takeaway struct {
	TakeawayType string `json:"takeaway_type,omitempty"` // 外卖类型（PICKUP/DELIVERY）

	ContactName  string `json:"contact_name,omitempty"`  // 联系人/收货人
	ContactPhone string `json:"contact_phone,omitempty"` // 联系电话

	PickupNo    string `json:"pickup_no,omitempty"`     // 取餐号（自取）
	PickupEtaAt string `json:"pickup_eta_at,omitempty"` // 预计取餐时间（RFC3339 字符串）

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

// Promotion 促销信息
type Promotion struct {
	PromotionID   string `json:"promotion_id,omitempty"`   // 促销ID
	PromotionName string `json:"promotion_name,omitempty"` // 促销名称
	PromotionType string `json:"promotion_type,omitempty"` // 促销类型

	DiscountAmount int64       `json:"discount_amount"` // 促销优惠金额（分）
	Meta           interface{} `json:"meta,omitempty"`  // 促销扩展信息
}

// Fee 费用明细
type Fee struct {
	FeeID   string `json:"fee_id,omitempty"`   // 费用ID
	FeeName string `json:"fee_name,omitempty"` // 费用名称
	FeeType string `json:"fee_type,omitempty"` // 费用类型（TIP/PACKAGING/SURCHARGE/OTHER 等）

	Amount int64       `json:"amount"`         // 费用金额（分）
	Meta   interface{} `json:"meta,omitempty"` // 扩展信息
}

// Coupon 卡券信息
type Coupon struct {
	CouponID   string `json:"coupon_id,omitempty"`   // 卡券ID
	CouponName string `json:"coupon_name,omitempty"` // 卡券名称
	CouponType string `json:"coupon_type,omitempty"` // 卡券类型
	CouponCode string `json:"coupon_code,omitempty"` // 卡券码

	DiscountAmount int64       `json:"discount_amount"` // 卡券优惠金额（分）
	Meta           interface{} `json:"meta,omitempty"`  // 卡券扩展信息
}

// TaxRate 税率信息
type TaxRate struct {
	TaxRateID   string `json:"tax_rate_id,omitempty"`   // 税率ID
	TaxRateName string `json:"tax_rate_name,omitempty"` // 税率名称

	Rate int64 `json:"rate"` // 税率（万分比，如 600 表示 6%）

	TaxableAmount int64       `json:"taxable_amount"` // 计税金额（分）
	TaxAmount     int64       `json:"tax_amount"`     // 税额（分）
	Meta          interface{} `json:"meta,omitempty"` // 扩展信息
}

// Product 商品
type Product struct {
	OrderItemID string `json:"order_item_id"`   // 订单内明细ID
	Index       int    `json:"index,omitempty"` // 下单序号（同订单内第几次下单）

	RefundReason string    `json:"refund_reason,omitempty"` // 退菜原因
	RefundedBy   string    `json:"refunded_by,omitempty"`   // 退菜操作人
	RefundedAt   time.Time `json:"refunded_at,omitempty"`   // 退菜时间

	Promotions        []Promotion `json:"promotions,omitempty"` // 促销明细（可叠加）
	PromotionDiscount int64       `json:"promotion_discount"`   // 促销优惠金额（分）

	ProductID   string `json:"product_id,omitempty"` // 商品ID
	ProductName string `json:"product_name"`         // 商品名（快照）
	SkuID       string `json:"sku_id,omitempty"`     // SKU ID
	SkuName     string `json:"sku_name,omitempty"`   // SKU 名称（快照）

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
	Options interface{} `json:"options,omitempty"` // 做法/加料（结构化）
}

// Refund 退款单信息（order_type=REFUND/PARTIAL_REFUND）
type Refund struct {
	OriginOrderID uuid.UUID `json:"origin_order_id,omitempty"` // 原正单订单ID
	OriginOrderNo string    `json:"origin_order_no,omitempty"` // 原正单订单号
	Reason        string    `json:"reason,omitempty"`          // 退款说明
}

// Amount 金额汇总（分）
type Amount struct {
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

// Payment 支付记录
type Payment struct {
	PaymentNo     string `json:"payment_no"`     // 支付号（第三方/外部交易号）
	PaymentMethod string `json:"payment_method"` // 支付方式
	PaymentAmount int64  `json:"payment_amount"` // 支付金额（分）

	POS     POS     `json:"pos"`     // POS 终端信息
	Cashier Cashier `json:"cashier"` // 收银员信息

	PaidAt time.Time `json:"paid_at,omitempty"` // 支付时间
}

// Receipt 结账后小票
type Receipt struct {
	OrderID    uuid.UUID `json:"order_id"`    // 关联订单ID
	MerchantID uuid.UUID `json:"merchant_id"` // 品牌商ID
	Store      Store     `json:"store"`       // 门店信息
	Channel    Channel   `json:"channel"`     // 下单渠道
	POS        POS       `json:"pos"`         // POS 终端信息
	Cashier    Cashier   `json:"cashier"`     // 收银员信息

	BusinessDate string `json:"business_date"` // 营业日（字符串，便于快照展示）

	OrderNo    string `json:"order_no"`    // 订单号
	OrderType  string `json:"order_type"`  // 订单类型（SALE/REFUND/PARTIAL_REFUND）
	DiningMode string `json:"dining_mode"` // 堂食/外卖（DINE_IN/TAKEAWAY）

	Member   Member   `json:"member"`   // 会员信息
	Takeaway Takeaway `json:"takeaway"` // 外卖信息（dining_mode=TAKEAWAY）

	Refund Refund `json:"refund"` // 退款单信息（退单/部分退款单）

	Products   []Product   `json:"products"`             // 下单商品
	Promotions []Promotion `json:"promotions,omitempty"` // 促销
	Coupons    []Coupon    `json:"coupons,omitempty"`    // 卡券
	TaxRates   []TaxRate   `json:"tax_rates,omitempty"`  // 税率
	Fees       []Fee       `json:"fees,omitempty"`       // 费用明细
	Payments   []Payment   `json:"payments,omitempty"`   // 支付记录

	Amount Amount `json:"amount"` // 金额汇总（分）
}

type CreateOrderReq struct {
	MerchantID uuid.UUID `json:"merchant_id" binding:"required"` // 品牌商ID
	StoreID    uuid.UUID `json:"store_id" binding:"required"`    // 门店ID

	BusinessDate string `json:"business_date" binding:"required"` // 营业日
	ShiftNo      string `json:"shift_no,omitempty"`               // 班次号

	OrderNo       string    `json:"order_no,omitempty"`                                      // 订单号（可选，未传则后端生成）
	OrderType     string    `json:"order_type,omitempty" enums:"SALE,REFUND,PARTIAL_REFUND"` // 订单类型（SALE/REFUND/PARTIAL_REFUND）
	OriginOrderID uuid.UUID `json:"origin_order_id,omitempty"`                               // 原正单订单ID（退款/部分退款单使用）
	Refund        *Refund   `json:"refund,omitempty"`                                        // 退款单信息

	DiningMode string   `json:"dining_mode" binding:"required" enums:"DINE_IN,TAKEAWAY"` // 堂食/外卖（DINE_IN/TAKEAWAY）
	Store      *Store   `json:"store,omitempty"`
	Channel    *Channel `json:"channel,omitempty"`
	POS        *POS     `json:"pos,omitempty"`
	Cashier    *Cashier `json:"cashier,omitempty"`

	OrderStatus       string `json:"order_status,omitempty" enums:"DRAFT,PLACED,IN_PROGRESS,READY,COMPLETED,CANCELLED,VOIDED,MERGED"`
	PaymentStatus     string `json:"payment_status,omitempty" enums:"UNPAID,PAYING,PARTIALLY_PAID,PAID,PARTIALLY_REFUNDED,REFUNDED"`
	FulfillmentStatus string `json:"fulfillment_status,omitempty" enums:"NONE,IN_RESTAURANT,SERVED,PICKUP_PENDING,PICKED_UP,DELIVERING,DELIVERED"`
	TableStatus       string `json:"table_status,omitempty" enums:"OPENED,TRANSFERRED,RELEASED"`
	// 桌位状态

	TableID       string `json:"table_id,omitempty"`
	TableName     string `json:"table_name,omitempty"`
	TableCapacity int    `json:"table_capacity,omitempty"`
	GuestCount    int    `json:"guest_count,omitempty"`

	OpenedAt    *time.Time `json:"opened_at,omitempty"`
	OpenedBy    string     `json:"opened_by,omitempty"`
	PlacedAt    *time.Time `json:"placed_at,omitempty"`
	PlacedBy    string     `json:"placed_by,omitempty"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	PaidBy      string     `json:"paid_by,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	Member   *Member   `json:"member,omitempty"`
	Takeaway *Takeaway `json:"takeaway,omitempty"`

	Cart            *[]Product   `json:"cart,omitempty"`
	Products        *[]Product   `json:"products,omitempty"`
	Promotions      *[]Promotion `json:"promotions,omitempty"`
	Coupons         *[]Coupon    `json:"coupons,omitempty"`
	TaxRates        *[]TaxRate   `json:"tax_rates,omitempty"`
	Fees            *[]Fee       `json:"fees,omitempty"`
	Payments        *[]Payment   `json:"payments,omitempty"`
	RefundsProducts *[]Product   `json:"refunds_products,omitempty"`

	Amount *Amount `json:"amount,omitempty"`
}

type UpdateOrderReq struct {
	BusinessDate *string `json:"business_date,omitempty"`
	ShiftNo      *string `json:"shift_no,omitempty"`

	OrderNo       *string    `json:"order_no,omitempty"`
	OrderType     *string    `json:"order_type,omitempty" enums:"SALE,REFUND,PARTIAL_REFUND"` // 订单类型（SALE/REFUND/PARTIAL_REFUND）
	OriginOrderID *uuid.UUID `json:"origin_order_id,omitempty"`
	Refund        *Refund    `json:"refund,omitempty"`

	DiningMode *string  `json:"dining_mode,omitempty" enums:"DINE_IN,TAKEAWAY"` // 堂食/外卖（DINE_IN/TAKEAWAY）
	Store      *Store   `json:"store,omitempty"`
	Channel    *Channel `json:"channel,omitempty"`
	POS        *POS     `json:"pos,omitempty"`
	Cashier    *Cashier `json:"cashier,omitempty"`

	OrderStatus       *string `json:"order_status,omitempty" enums:"DRAFT,PLACED,IN_PROGRESS,READY,COMPLETED,CANCELLED,VOIDED,MERGED"`
	PaymentStatus     *string `json:"payment_status,omitempty" enums:"UNPAID,PAYING,PARTIALLY_PAID,PAID,PARTIALLY_REFUNDED,REFUNDED"`
	FulfillmentStatus *string `json:"fulfillment_status,omitempty" enums:"NONE,IN_RESTAURANT,SERVED,PICKUP_PENDING,PICKED_UP,DELIVERING,DELIVERED"`
	TableStatus       *string `json:"table_status,omitempty" enums:"OPENED,TRANSFERRED,RELEASED"`

	TableID       *string `json:"table_id,omitempty"`
	TableName     *string `json:"table_name,omitempty"`
	TableCapacity *int    `json:"table_capacity,omitempty"`
	GuestCount    *int    `json:"guest_count,omitempty"`

	OpenedAt    *time.Time `json:"opened_at,omitempty"`
	OpenedBy    *string    `json:"opened_by,omitempty"`
	PlacedAt    *time.Time `json:"placed_at,omitempty"`
	PlacedBy    *string    `json:"placed_by,omitempty"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	PaidBy      *string    `json:"paid_by,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	Member   *Member   `json:"member,omitempty"`
	Takeaway *Takeaway `json:"takeaway,omitempty"`

	Cart            *[]Product   `json:"cart,omitempty"`
	Products        *[]Product   `json:"products,omitempty"`
	Promotions      *[]Promotion `json:"promotions,omitempty"`
	Coupons         *[]Coupon    `json:"coupons,omitempty"`
	TaxRates        *[]TaxRate   `json:"tax_rates,omitempty"`
	Fees            *[]Fee       `json:"fees,omitempty"`
	Payments        *[]Payment   `json:"payments,omitempty"`
	RefundsProducts *[]Product   `json:"refunds_products,omitempty"`

	Amount *Amount `json:"amount,omitempty"`
}

type ListOrderReq struct {
	MerchantID uuid.UUID `form:"merchant_id" binding:"required"` // 品牌商ID
	StoreID    uuid.UUID `form:"store_id" binding:"required"`    // 门店ID

	BusinessDate  string `form:"business_date"`                                                                         // 营业日
	OrderNo       string `form:"order_no"`                                                                              // 订单号
	OrderType     string `form:"order_type" enums:"SALE,REFUND,PARTIAL_REFUND"`                                         // 订单类型
	OrderStatus   string `form:"order_status" enums:"DRAFT,PLACED,IN_PROGRESS,READY,COMPLETED,CANCELLED,VOIDED,MERGED"` // 订单状态
	PaymentStatus string `form:"payment_status" enums:"UNPAID,PAYING,PARTIALLY_PAID,PAID,PARTIALLY_REFUNDED,REFUNDED"`  // 支付状态

	upagination.RequestPagination
}

type ListOrderResp struct {
	Items      []*OrderResp            `json:"items"`      // 列表数据
	Pagination *upagination.Pagination `json:"pagination"` // 分页信息
}

// OrderResp 订单返回结构
type OrderResp struct {
	OrderID uuid.UUID `json:"order_id"` // 订单ID

	MerchantID uuid.UUID `json:"merchant_id"` // 品牌商ID
	Store      Store     `json:"store"`       // 门店信息

	BusinessDate string `json:"business_date"`      // 营业日
	ShiftNo      string `json:"shift_no,omitempty"` // 班次号

	OrderNo   string `json:"order_no"`   // 订单号
	OrderType string `json:"order_type"` // 订单类型（SALE/REFUND/PARTIAL_REFUND）
	Refund    Refund `json:"refund"`     // 退款单信息

	DiningMode string  `json:"dining_mode"` // 堂食/外卖
	Channel    Channel `json:"channel"`     // 下单渠道
	POS        POS     `json:"pos"`         // POS 终端信息
	Cashier    Cashier `json:"cashier"`     // 收银员信息

	OrderStatus       string `json:"order_status"`       // 订单状态
	PaymentStatus     string `json:"payment_status"`     // 支付状态
	FulfillmentStatus string `json:"fulfillment_status"` // 交付状态
	TableStatus       string `json:"table_status"`       // 桌位状态

	TableID       string `json:"table_id,omitempty"`       // 桌位ID（堂食）
	TableName     string `json:"table_name,omitempty"`     // 桌位名称（堂食，如 A01/1号桌）
	TableCapacity int    `json:"table_capacity,omitempty"` // 桌位容量（几人桌，如 2/4/6，仅堂食）
	GuestCount    int    `json:"guest_count,omitempty"`    // 用餐人数（堂食）

	OpenedAt    time.Time `json:"opened_at,omitempty"`    // 开台时间
	OpenedBy    string    `json:"opened_by,omitempty"`    // 开台人
	PlacedAt    time.Time `json:"placed_at,omitempty"`    // 下单时间
	PlacedBy    string    `json:"placed_by,omitempty"`    // 下单人
	PaidAt      time.Time `json:"paid_at,omitempty"`      // 支付完成时间
	PaidBy      string    `json:"paid_by,omitempty"`      // 收款人
	CompletedAt time.Time `json:"completed_at,omitempty"` // 完成时间

	Member   Member   `json:"member"`   // 会员信息
	Takeaway Takeaway `json:"takeaway"` // 外卖信息（dining_mode=TAKEAWAY）

	Cart            []Product   `json:"cart,omitempty"`             // 购物车（未下单商品）
	Products        []Product   `json:"products"`                   // 下单商品
	Promotions      []Promotion `json:"promotions,omitempty"`       // 促销
	Coupons         []Coupon    `json:"coupons,omitempty"`          // 卡券
	TaxRates        []TaxRate   `json:"tax_rates,omitempty"`        // 税率
	Fees            []Fee       `json:"fees,omitempty"`             // 费用明细
	Payments        []Payment   `json:"payments,omitempty"`         // 支付记录
	RefundsProducts []Product   `json:"refunds_products,omitempty"` // 退菜记录

	Amount Amount `json:"amount"` // 金额汇总（分）
}
