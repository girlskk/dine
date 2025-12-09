package domain

import (
	"context"
	"encoding/json"
	"regexp"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/zxh"
)

const (
	TaskTypePaymentCallback = "payment_callback"
)

// 知心话支付码前缀
var (
	regexPointCode       = regexp.MustCompile(`^z\d+$`)  // 积分支付
	regexPointWalletCode = regexp.MustCompile(`^zx\d+$`) // 钱包支付
)

type (
	PayChannel  string // 支付渠道
	PayState    string // 支付状态
	PayProvider string // 支付供应商
	PayBizType  string // 支付业务类型
)

const (
	PayChannelWechatPay   PayChannel = "wxpay"        // 微信支付
	PayChannelAlipay      PayChannel = "alipay"       // 支付宝支付
	PayChannelPoint       PayChannel = "point"        // 积分支付
	PayChannelPointWallet PayChannel = "point_wallet" // 知心话钱包支付
)

func (PayChannel) Values() []string {
	return []string{
		string(PayChannelWechatPay),
		string(PayChannelAlipay),
		string(PayChannelPoint),
		string(PayChannelPointWallet),
	}
}

const (
	PayStateUnknown    PayState = "U" // 未知
	PayStateProcessing PayState = "P" // 处理中
	PayStateSuccess    PayState = "S" // 成功
	PayStateFailure    PayState = "F" // 失败
	PayStateWaiting    PayState = "W" // 等待
)

func (PayState) Values() []string {
	return []string{
		string(PayStateUnknown),
		string(PayStateProcessing),
		string(PayStateSuccess),
		string(PayStateFailure),
		string(PayStateWaiting),
	}
}

func (state PayState) IsFinished() bool {
	return state == PayStateSuccess || state == PayStateFailure
}

const (
	PayProviderZhiXinHua       PayProvider = "zxh"        // 知心话
	PayProviderZhiXinHuaWallet PayProvider = "zxh_wallet" // 知心话钱包
	PayProviderHuifu           PayProvider = "huifu"      //汇付
)

func (PayProvider) Values() []string {
	return []string{
		string(PayProviderZhiXinHua),
		string(PayProviderZhiXinHuaWallet),
		string(PayProviderHuifu),
	}
}

const (
	PayBizTypeOrder PayBizType = "order" // 订单支付
)

func (PayBizType) Values() []string {
	return []string{
		string(PayBizTypeOrder),
	}
}

type PaymentMemberInfo struct {
	ID    int
	Name  string
	Phone string
}

type Payment struct {
	ID          int             `json:"id"`
	SeqNo       string          `json:"seq_no"`       // 流水号
	Provider    PayProvider     `json:"provider"`     // 支付供应商
	Channel     PayChannel      `json:"channel"`      // 支付渠道
	State       PayState        `json:"state"`        // 支付状态
	Amount      decimal.Decimal `json:"amount"`       // 支付金额
	GoodsDesc   string          `json:"goods_desc"`   // 商品描述
	MchID       string          `json:"mch_id"`       // 商户ID
	IPAddr      string          `json:"ip_addr"`      // IP地址
	Req         json.RawMessage `json:"req"`          // 请求参数
	Resp        json.RawMessage `json:"resp"`         // 响应参数
	Callback    json.RawMessage `json:"callback"`     // 回调
	FinishedAt  *time.Time      `json:"finished_at"`  // 完成时间
	Refunded    decimal.Decimal `json:"refunded"`     // 已退款金额
	FailReason  string          `json:"fail_reason"`  // 失败原因
	PayBizType  PayBizType      `json:"pay_biz_type"` // 支付业务类型
	BizID       int             `json:"biz_id"`       // 业务ID
	CreatorType OperatorType    `json:"creator_type"` // 创建人类型
	CreatorID   int             `json:"creator_id"`   // 创建人ID
	CreatorName string          `json:"creator_name"` // 创建人姓名
	StoreID     int             `json:"store_id"`     // 门店ID
	CreatedAt   time.Time       `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time       `json:"updated_at"`   // 更新时间
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/payment_repository.go -package=mock . PaymentRepository
type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) (*Payment, error)
	Update(ctx context.Context, payment *Payment) (*Payment, error)
	FindBySeqNo(ctx context.Context, seqNo string) (*Payment, error)
	CreateCallback(ctx context.Context, callback *PaymentCallback) (*PaymentCallback, error)
	GetCallback(ctx context.Context, callbackID int) (*PaymentCallback, error)
	RemoveCallback(ctx context.Context, callbackID int) error
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/payment_interactor.go -package=mock . PaymentInteractor
type PaymentInteractor interface {
	PayHuifuCallback(ctx context.Context, sign, respData string) (err error)
	PayZxhCallback(ctx context.Context, params zxh.ZhixinhuaPointPayCallBack) (err error)
	PayPolling(ctx context.Context, seqNo string, operator *FrontendUser) (state PayState, reason string, err error)
}

type PaymentCallbackType string

const (
	PaymentCallbackTypePay    PaymentCallbackType = "pay"
	PaymentCallbackTypeRefund PaymentCallbackType = "refund"
)

func (PaymentCallbackType) Values() []string {
	return []string{
		string(PaymentCallbackTypePay),
		string(PaymentCallbackTypeRefund),
	}
}

type PaymentCallback struct {
	ID       int                 `json:"id"`
	SeqNo    string              `json:"seq_no"`   // 流水号
	Type     PaymentCallbackType `json:"type"`     // 回调类型
	Raw      json.RawMessage     `json:"raw"`      // 原始数据
	Provider PayProvider         `json:"provider"` // 支付供应商
}

type PaymentCallbackTaskPayload struct {
	CallbackID int         `json:"callback_id"`
	SeqNo      string      `json:"seq_no"`
	Provider   PayProvider `json:"provider"`
}

func IsPointCode(code string) bool {
	return regexPointCode.MatchString(code)
}

func IsPointWalletCode(code string) bool {
	return regexPointWalletCode.MatchString(code)
}

type PaymentEventBaseParams struct {
	DataStore DataStore
	Payment   *Payment
	Operator  any
}

type PaymentEventSuccessParams struct {
	PaymentEventBaseParams
	MemberInfo PaymentMemberInfo
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/payment_event_trigger.go -package=mock . PaymentEventTrigger
type PaymentEventTrigger interface {
	FireSuccess(ctx context.Context, params *PaymentEventSuccessParams) error
}
