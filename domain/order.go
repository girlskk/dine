package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

const (
	OrderListExportSingleMaxSize int = 3000
)

type OrderListOrderBy int

const (
	_ OrderListOrderBy = iota
	OrderListOrderByID
	OrderListOrderByCreatedAt
)

type OrderListOrder struct {
	OrderBy OrderListOrderBy
	Desc    bool
}

func NewOrderListOrderByCreatedAt(desc bool) OrderListOrder {
	return OrderListOrder{
		OrderBy: OrderListOrderByCreatedAt,
		Desc:    desc,
	}
}

func NewOrderListOrderByID(desc bool) OrderListOrder {
	return OrderListOrder{
		OrderBy: OrderListOrderByID,
		Desc:    desc,
	}
}

type (
	OrderPaidChannel  string
	OrderPaidChannels []OrderPaidChannel
)

func (channel OrderPaidChannel) ToString() string {
	switch channel {
	case OrderPaidChannelCash:
		return "现金"
	case OrderPaidChannelWechatPay:
		return "微信"
	case OrderPaidChannelAlipay:
		return "支付宝"
	case OrderPaidChannelPoint:
		return "餐饮积分"
	case OrderPaidChannelPointWallet:
		return "知心话钱包"
	default:
		return ""
	}
}

func (channel OrderPaidChannel) Paid(od *Order) decimal.Decimal {
	switch channel {
	case OrderPaidChannelCash:
		return od.CashPaid
	case OrderPaidChannelWechatPay:
		return od.WechatPaid
	case OrderPaidChannelAlipay:
		return od.AlipayPaid
	case OrderPaidChannelPoint:
		return od.PointsPaid
	case OrderPaidChannelPointWallet:
		return od.PointsWalletPaid
	default:
		return decimal.Zero
	}
}

func (channels OrderPaidChannels) Add(channel OrderPaidChannel) OrderPaidChannels {
	channels = append(channels, channel)
	return lo.Uniq(channels)
}

func (channels OrderPaidChannels) Contains(channel OrderPaidChannel) bool {
	return lo.Contains(channels, channel)
}

const (
	OrderPaidChannelCash        OrderPaidChannel = "cash"         // 现金支付
	OrderPaidChannelWechatPay   OrderPaidChannel = "wechat"       // 微信支付
	OrderPaidChannelAlipay      OrderPaidChannel = "alipay"       // 支付宝支付
	OrderPaidChannelPoint       OrderPaidChannel = "point"        // 积分支付
	OrderPaidChannelPointWallet OrderPaidChannel = "point_wallet" // 知心话钱包支付
)

// 订单支付渠道编号前缀
const (
	OrderPaidChannelCashPrefix        string = "CASH" // 现金支付
	OrderPaidChannelWechatPayPrefix   string = "WX"   // 微信支付
	OrderPaidChannelAlipayPrefix      string = "ALI"  // 支付宝支付
	OrderPaidChannelPointPrefix       string = "PI"   // 积分支付
	OrderPaidChannelPointWalletPrefix string = "QB"   // 知心话钱包支付
)

func (OrderPaidChannel) Values() []string {
	return []string{
		string(OrderPaidChannelCash),
		string(OrderPaidChannelWechatPay),
		string(OrderPaidChannelAlipay),
		string(OrderPaidChannelPoint),
		string(OrderPaidChannelPointWallet),
	}
}

func (c OrderPaidChannel) Prefix() string {
	switch c {
	case OrderPaidChannelCash:
		return OrderPaidChannelCashPrefix
	case OrderPaidChannelWechatPay:
		return OrderPaidChannelWechatPayPrefix
	case OrderPaidChannelAlipay:
		return OrderPaidChannelAlipayPrefix
	case OrderPaidChannelPoint:
		return OrderPaidChannelPointPrefix
	case OrderPaidChannelPointWallet:
		return OrderPaidChannelPointWalletPrefix
	default:
		return ""
	}
}

type OrderStatus string

const (
	OrderStatusUnpaid   OrderStatus = "unpaid"
	OrderStatusPartPaid OrderStatus = "part_paid"
	OrderStatusPaid     OrderStatus = "paid"
	OrderStatusCanceled OrderStatus = "canceled"
)

func (OrderStatus) Values() []string {
	return []string{
		string(OrderStatusUnpaid),
		string(OrderStatusPartPaid),
		string(OrderStatusPaid),
		string(OrderStatusCanceled),
	}
}

func (s OrderStatus) ToString() string {
	switch s {
	case OrderStatusUnpaid:
		return "未支付"
	case OrderStatusPartPaid:
		return "部分支付"
	case OrderStatusPaid:
		return "已支付"
	case OrderStatusCanceled:
		return "已取消"
	default:
		return ""
	}
}

type OrderSource string

const (
	OrderSourceOffline     OrderSource = "offline"      // 线下下单
	OrderSourceMiniProgram OrderSource = "mini_program" // 小程序下单
)

func (OrderSource) Values() []string {
	return []string{
		string(OrderSourceOffline),
		string(OrderSourceMiniProgram),
	}
}

type OrderType string

const (
	OrderTypeDineIn OrderType = "dine_in" // 堂食
)

func (OrderType) Values() []string {
	return []string{
		string(OrderTypeDineIn),
	}
}

// 订单
type Order struct {
	ID                   int
	No                   string            `json:"no"`                     // 订单号
	Type                 OrderType         `json:"type"`                   // 订单类型
	Source               OrderSource       `json:"source"`                 // 订单来源
	Status               OrderStatus       `json:"status"`                 // 订单状态
	TotalPrice           decimal.Decimal   `json:"total_price"`            // 总价（优惠前）
	Discount             decimal.Decimal   `json:"discount"`               // 优惠金额
	RealPrice            decimal.Decimal   `json:"real_price"`             // 实际总价（优惠后）
	PointsAvailable      decimal.Decimal   `json:"points_available"`       // 积分可用额度
	LastPaidAt           *time.Time        `json:"last_paid_at,omitempty"` // 最后支付时间
	FinishedAt           *time.Time        `json:"finished_at,omitempty"`  // 完成时间
	MemberID             int               `json:"member_id"`              // 会员ID
	MemberName           string            `json:"member_name"`            // 会员姓名
	MemberPhone          string            `json:"member_phone"`           // 会员手机号
	StoreID              int               `json:"store_id"`               // 门店ID
	StoreName            string            `json:"store_name"`             // 门店名称
	TableID              int               `json:"table_id"`               // 台桌ID
	TableName            string            `json:"table_name"`             // 台桌名称
	PeopleNumber         int               `json:"people_number"`          // 就餐人数
	CreatorID            int               `json:"creator_id"`             // 创建人ID
	CreatorName          string            `json:"creator_name"`           // 创建人姓名
	CreatorType          OperatorType      `json:"creator_type"`           // 创建人类型
	Paid                 decimal.Decimal   `json:"paid"`                   // 已支付金额
	Refunded             decimal.Decimal   `json:"refunded"`               // 已退款金额
	PaidChannels         OrderPaidChannels `json:"paid_channels"`          // 支付渠道
	CashPaid             decimal.Decimal   `json:"cash_paid"`              // 现金已支付金额
	WechatPaid           decimal.Decimal   `json:"wechat_paid"`            // 微信已支付金额
	WechatRefunded       decimal.Decimal   `json:"wechat_refunded"`        // 微信已退款金额
	AlipayPaid           decimal.Decimal   `json:"alipay_paid"`            // 支付宝已支付金额
	AlipayRefunded       decimal.Decimal   `json:"alipay_refunded"`        // 支付宝已退款金额
	PointsPaid           decimal.Decimal   `json:"points_paid"`            // 积分已支付金额
	PointsRefunded       decimal.Decimal   `json:"points_refunded"`        // 积分已退款金额
	PointsWalletPaid     decimal.Decimal   `json:"points_wallet_paid"`     // 知心话钱包已支付金额
	PointsWalletRefunded decimal.Decimal   `json:"points_wallet_refunded"` // 知心话钱包已退款金额
	CreatedAt            time.Time         `json:"created_at"`             // 创建时间
	UpdatedAt            time.Time         `json:"updated_at"`             // 更新时间

	Items OrderItems  `json:"items,omitempty"` // 订单商品
	Logs  []*OrderLog `json:"logs,omitempty"`  // 订单日志
}

func (o *Order) CanPaid(amount decimal.Decimal, point bool) error {
	if o.Status != OrderStatusUnpaid && o.Status != OrderStatusPartPaid {
		return ParamsErrorf("订单当前状态无法支付")
	}

	if o.Paid.Add(amount).GreaterThan(o.RealPrice) {
		return ParamsErrorf("订单支付金额超出")
	}

	if point && o.PointsPaid.Add(amount).GreaterThan(o.PointsAvailable) {
		return ParamsErrorf("积分支付金额超出")
	}

	return nil
}

func (o *Order) GoodsDesc() string {
	return fmt.Sprintf("订单 %s", o.No)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/order_repository.go -package=mock . OrderRepository
type OrderRepository interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	FindByNo(ctx context.Context, no string) (*Order, error)
	Find(ctx context.Context, id int) (*Order, error)
	CreateLog(ctx context.Context, log *OrderLog) (*OrderLog, error)
	GetOrders(ctx context.Context, pager *upagination.Pagination, filter *OrderListFilter, orderBys ...OrderListOrder) (dorders []*Order, total int, err error)
	GetOrdersWithItems(ctx context.Context, pager *upagination.Pagination, filter *OrderListFilter, orderBys ...OrderListOrder) (dorders []*Order, total int, err error)
	Update(ctx context.Context, order *Order) (*Order, error)
	UpdateItem(ctx context.Context, item *OrderItem) (*OrderItem, error)
	AppendItems(ctx context.Context, orderID int, items []*OrderItem) ([]*OrderItem, error)
	RemoveItems(ctx context.Context, orderID int, itemIDs ...int) (err error)
	FindByItemID(ctx context.Context, itemID int) (*Order, error)
	CreateFinanceLog(ctx context.Context, log *OrderFinanceLog) (*OrderFinanceLog, error)
	HasIncompletePayment(ctx context.Context, orderID int) (has bool, err error)
	ListItemNamesByOrders(ctx context.Context, orderIDs []int) (map[int][]string, error)
	GetOrderRange(ctx context.Context, filter *OrderListFilter) (rg OrderRange, err error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/order_interactor.go -package=mock . OrderInteractor
type OrderInteractor interface {
	CreateOrder(ctx context.Context, params *CreateOrderParams) (*Order, error)
	CreateOrderFromCart(ctx context.Context, params *CreateOrderParams) (*Order, error)
	GetOrder(ctx context.Context, no string) (*Order, error)
	GetOrders(ctx context.Context, pager *upagination.Pagination, filter *OrderListFilter, withItems bool) (dorders []*Order, total int, err error)
	ModifyItemPrice(ctx context.Context, params *ModifyItemPriceParams) (err error)
	AppendItems(ctx context.Context, params *AppendItemParams) (err error)
	AppendItemsFromCart(ctx context.Context, params *AppendItemParams) (err error)
	TurnTable(ctx context.Context, params *TurnTableParams) (err error)
	RemoveItems(ctx context.Context, params *RemoveItemParams) (err error)
	CancelOrder(ctx context.Context, no string, operator any) (err error)
	DiscountOrder(ctx context.Context, params *DiscountOrderParams) (err error)
	CashPaid(ctx context.Context, params *OrderCashPaidParams) (err error)
	ScanPaid(ctx context.Context, params *OrderScanPaidParams) (seqNo string, err error)
	GetOrderRange(ctx context.Context, filter *OrderListFilter) (rg OrderRange, err error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/order_event_trigger.go -package=mock . OrderEventTrigger
type OrderEventTrigger interface {
	FireCreateOrder(ctx context.Context, params *OrderEventBaseParams) error
	FireModifyPrice(ctx context.Context, params *OrderEventBaseParams) error
	FireAppendItem(ctx context.Context, params *OrderEventBaseParams) error
	FireRemoveItem(ctx context.Context, params *OrderEventBaseParams) error
	FireTurnTable(ctx context.Context, params *OrderEventTurnTableParams) error
	FirePaid(ctx context.Context, params *OrderEventPaidParams) error
	FireCancel(ctx context.Context, params *OrderEventBaseParams) error
	FireDiscount(ctx context.Context, params *OrderEventBaseParams) error
	FireFinish(ctx context.Context, params *OrderEventBaseParams) error
}

type OrderEventBaseParams struct {
	DataStore     DataStore
	Order         *Order
	OperatedItems OrderItems // 当前操作的订单商品
	Operator      any
}

type OrderEventTurnTableParams struct {
	OrderEventBaseParams
	OldTableID int
}

type OrderEventPaidParams struct {
	OrderEventBaseParams
	Amount  decimal.Decimal
	Channel OrderPaidChannel
	SeqNo   string
}

type CreateOrderParams struct {
	Table        *Table
	Creator      Operator
	Store        *Store
	Items        []*CreateOrderItem
	PeopleNumber int
	Source       OrderSource
}

type CreateOrderItem struct {
	ProductID       int
	Quantity        decimal.Decimal
	Price           decimal.Decimal
	ProductSpecID   int
	ProductAttrID   int
	ProductRecipeID int
	Remark          string
}

type OrderListFilter struct {
	StoreID           int          `json:"store_id"`
	Status            OrderStatus  `json:"status"`
	FinishedAtGte     *time.Time   `json:"finished_at_gte"`
	FinishedAtLte     *time.Time   `json:"finished_at_lte"`
	HasItemName       string       `json:"has_item_name"`
	MemberNameOrPhone string       `json:"member_name_or_phone"`
	CreatedAtGte      *time.Time   `json:"created_at_gte"`
	CreatedAtLte      *time.Time   `json:"created_at_lte"`
	PointsPaidGt0     bool         `json:"points_paid_gt0"`
	IDGte             int          `json:"id_gte"`
	IDLte             int          `json:"id_lte"`
	CreatorID         int          `json:"creator_id"`
	CreatorType       OperatorType `json:"creator_type"`
}

type ModifyItemPriceParams struct {
	OrderNo  string
	ItemID   int
	Price    decimal.Decimal
	Operator *FrontendUser
}

type AppendItemParams struct {
	OrderNo  string
	Items    []*CreateOrderItem
	Operator Operator
	TableID  int
}

type DiscountOrderParams struct {
	OrderNo  string
	Discount decimal.Decimal
	Operator *FrontendUser
}

type RemoveItemParams struct {
	OrderNo  string
	ItemID   int
	Quantity decimal.Decimal
	Operator *FrontendUser
}

type TurnTableParams struct {
	OrderNo  string
	TableID  int
	Operator *FrontendUser
}

type OrderCashPaidParams struct {
	OrderNo  string
	Amount   decimal.Decimal
	Operator *FrontendUser
}

type OrderScanPaidParams struct {
	OrderNo        string
	Amount         decimal.Decimal
	AuthCode       string
	Operator       *FrontendUser
	IPAddr         string
	HuifuNotifyURL string
	ZxhNotifyURL   string
}

type OrderListExportParams struct {
	Filter OrderListFilter        `json:"filter"`
	Pager  upagination.Pagination `json:"pager"`
}

type OrderRange struct {
	MinID int
	MaxID int
	Count int
}
