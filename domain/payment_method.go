package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// ------------------------------------------------------------
// 错误定义
// ------------------------------------------------------------

var (
	ErrPaymentMethodNotExists = errors.New("结算方式不存在")
)

// ------------------------------------------------------------
// 枚举定义ßß
// ------------------------------------------------------------

// PaymentMethodAccountingRule 计入规则
type PaymentMethodAccountingRule string

const (
	PaymentMethodAccountingRuleIncome   PaymentMethodAccountingRule = "income"   // 计入实收
	PaymentMethodAccountingRuleDiscount PaymentMethodAccountingRule = "discount" // 计入优惠
)

func (PaymentMethodAccountingRule) Values() []string {
	return []string{
		string(PaymentMethodAccountingRuleIncome),
		string(PaymentMethodAccountingRuleDiscount),
	}
}

// PaymentMethodPayType 结算类型
type PaymentMethodPayType string

const (
	PaymentMethodPayTypeOther         PaymentMethodPayType = "other"          // 其他
	PaymentMethodPayTypeCash          PaymentMethodPayType = "cash"           // 现金
	PaymentMethodPayTypeOfflineCard   PaymentMethodPayType = "offline_card"   // 线下刷卡
	PaymentMethodPayTypeCustomCoupon  PaymentMethodPayType = "custom_coupon"  // 自定义券
	PaymentMethodPayTypePartnerCoupon PaymentMethodPayType = "partner_coupon" // 三方合作券
)

func (PaymentMethodPayType) Values() []string {
	return []string{
		string(PaymentMethodPayTypeOther),
		string(PaymentMethodPayTypeCash),
		string(PaymentMethodPayTypeOfflineCard),
		string(PaymentMethodPayTypeCustomCoupon),
		string(PaymentMethodPayTypePartnerCoupon),
	}
}

// PaymentMethodInvoiceRule 实收部分开票规则
type PaymentMethodInvoiceRule string

const (
	PaymentMethodInvoiceRuleNotInvoice   PaymentMethodInvoiceRule = "no_invoice"    // 不开发票
	PaymentMethodInvoiceRuleActualAmount PaymentMethodInvoiceRule = "actual_amount" // 按实收金额
)

func (PaymentMethodInvoiceRule) Values() []string {
	return []string{
		string(PaymentMethodInvoiceRuleNotInvoice),
		string(PaymentMethodInvoiceRuleActualAmount),
	}
}

// PaymentMethodDisplayChannel 收银终端显示渠道枚举
type PaymentMethodDisplayChannel string

const (
	PaymentMethodDisplayChannelPOS                SaleChannel = "POS"         // POS
	PaymentMethodDisplayChannelMobileOrdering     SaleChannel = "Mobile"      // 移动点餐
	PaymentMethodDisplayChannelScanOrdering       SaleChannel = "Scan"        // 扫码点餐
	PaymentMethodDisplayChannelSelfService        SaleChannel = "SelfService" // 自助点餐
	PaymentMethodDisplayChannelThirdPartyDelivery SaleChannel = "ThirdParty"  // 三方外卖
)

func (PaymentMethodDisplayChannel) Values() []string {
	return []string{
		string(PaymentMethodDisplayChannelPOS),
		string(PaymentMethodDisplayChannelMobileOrdering),
		string(PaymentMethodDisplayChannelScanOrdering),
		string(PaymentMethodDisplayChannelSelfService),
		string(PaymentMethodDisplayChannelThirdPartyDelivery),
	}
}

type PaymentMethodSource string

const (
	PaymentMethodSourceBrand  PaymentMethodSource = "brand"  // 品牌
	PaymentMethodSourceStore  PaymentMethodSource = "store"  // 门店
	PaymentMethodSourceSystem PaymentMethodSource = "system" // 系统
)

func (PaymentMethodSource) Values() []string {
	return []string{
		string(PaymentMethodSourceBrand),
		string(PaymentMethodSourceStore),
		string(PaymentMethodSourceSystem),
	}
}

type PaymentMethod struct {
	ID               uuid.UUID                     `json:"id"`
	MerchantID       uuid.UUID                     `json:"merchant_id"`        // 品牌商ID
	Name             string                        `json:"name"`               // 结算方式名称
	AccountingRule   PaymentMethodAccountingRule   `json:"accounting_rule"`    // 计入规则:income-计入实收,discount-计入优惠
	PaymentType      PaymentMethodPayType          `json:"payment_type"`       // 结算类型:other-其他,cash-现金,offline_card-线下刷卡,custom_coupon-自定义券,partner_coupon-三方合作券
	FeeRate          *decimal.Decimal              `json:"fee_rate"`           // 手续费率,百分比
	InvoiceRule      PaymentMethodInvoiceRule      `json:"invoice_rule"`       // 实收部分开票规则:no_invoice-不开发票,actual_amount-按实收金额
	CashDrawerStatus bool                          `json:"cash_drawer_status"` // 开钱箱状态:false-不开钱箱, true-开钱箱（必选）
	DisplayChannels  []PaymentMethodDisplayChannel `json:"display_channels"`   // 收银终端显示渠道（可选，可多选）：POS、移动点餐、扫码点餐、自助点餐、三方外卖
	Status           bool                          `json:"status"`             // 启用/停用状态: true-启用, false-停用（必选）
	CreatedAt        time.Time                     `json:"created_at"`         // 创建时间
	UpdatedAt        time.Time                     `json:"updated_at"`         // 更新时间
}

type PaymentMethods []*PaymentMethod

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/payment_method_repository.go -package=mock . PaymentMethodRepository
type PaymentMethodRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*PaymentMethod, error)
	GetDetail(ctx context.Context, id uuid.UUID) (*PaymentMethod, error)
	Create(ctx context.Context, menu *PaymentMethod) error
	Update(ctx context.Context, menu *PaymentMethod) error
	Delete(ctx context.Context, id uuid.UUID) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params PaymentMethodSearchParams) (*PaymentMethodSearchRes, error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/payment_method_interactor.go -package=mock . PaymentMethodInteractor
type PaymentMethodInteractor interface {
	Create(ctx context.Context, menu *PaymentMethod) error
	Update(ctx context.Context, menu *PaymentMethod, user User) error
	Delete(ctx context.Context, id uuid.UUID, user User) error
	GetDetail(ctx context.Context, id uuid.UUID, user User) (*PaymentMethod, error)
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params PaymentMethodSearchParams) (*PaymentMethodSearchRes, error)
}
type PaymentMethodSearchParams struct {
	MerchantID uuid.UUID
	Name       string // 菜单名称（模糊匹配）
}

// PaymentMethodSearchRes 查询结果
type PaymentMethodSearchRes struct {
	*upagination.Pagination
	Items PaymentMethods `json:"items"`
}
