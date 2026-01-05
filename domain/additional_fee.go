package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrAdditionalFeeNotExists  = errors.New("附加费不存在")
	ErrAdditionalFeeNameExists = errors.New("附加费名称已存在")
)

// AdditionalFeeRepository 附加费仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/additional_fee_repository.go -package=mock . AdditionalFeeRepository
type AdditionalFeeRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (fee *AdditionalFee, err error)
	ListByIDs(ctx context.Context, ids []uuid.UUID) (fees []*AdditionalFee, err error)
	Create(ctx context.Context, fee *AdditionalFee) (err error)
	Update(ctx context.Context, fee *AdditionalFee) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetAdditionalFees(ctx context.Context, pager *upagination.Pagination, filter *AdditionalFeeListFilter, orderBys ...AdditionalFeeOrderBy) (fees []*AdditionalFee, total int, err error)
	Exists(ctx context.Context, params AdditionalFeeExistsParams) (exists bool, err error)
}

// AdditionalFeeInteractor 附加费用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/additional_fee_interactor.go -package=mock . AdditionalFeeInteractor
type AdditionalFeeInteractor interface {
	Create(ctx context.Context, fee *AdditionalFee) (err error)
	Update(ctx context.Context, fee *AdditionalFee) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetAdditionalFee(ctx context.Context, id uuid.UUID) (fee *AdditionalFee, err error)
	GetAdditionalFees(ctx context.Context, pager *upagination.Pagination, filter *AdditionalFeeListFilter, orderBys ...AdditionalFeeOrderBy) (fees []*AdditionalFee, total int, err error)
	AdditionalFeeSimpleUpdate(ctx context.Context, updateField AdditionalFeeSimpleUpdateType, fee *AdditionalFee) (err error)
}
type DiningWay string

const (
	DiningWayDineIn   DiningWay = "dine_in"  // 堂食
	DiningWayTakeOut  DiningWay = "take_out" // 外带
	DiningWayDelivery DiningWay = "delivery" // 外送
)

func (DiningWay) Values() []string {
	return []string{string(DiningWayDineIn), string(DiningWayTakeOut), string(DiningWayDelivery)}
}

type OrderChannel string

const (
	OrderChannelPOS           OrderChannel = "pos"            // POS 端
	OrderChannelSelfOrder     OrderChannel = "self_order"     // 自助点餐
	OrderChannelMiniProgram   OrderChannel = "mini_program"   // 小程序
	OrderChannelMobileOrder   OrderChannel = "mobile_order"   // 手机点餐
	OrderChannelScanOrder     OrderChannel = "scan_order"     // 扫码点餐
	OrderChannelThirdDelivery OrderChannel = "third_delivery" // 三方外卖
)

func (OrderChannel) Values() []string {
	return []string{
		string(OrderChannelPOS),
		string(OrderChannelSelfOrder),
		string(OrderChannelMiniProgram),
		string(OrderChannelMobileOrder),
		string(OrderChannelScanOrder),
		string(OrderChannelThirdDelivery),
	}
}

type AdditionalCategory string

const (
	AdditionalCategoryService    AdditionalCategory = "service_fee"    // 服务费
	AdditionalCategoryAdditional AdditionalCategory = "additional_fee" // 附加费
	AdditionalCategoryPacking    AdditionalCategory = "packing_fee"    // 打包费
)

func (AdditionalCategory) Values() []string {
	return []string{string(AdditionalCategoryService), string(AdditionalCategoryAdditional), string(AdditionalCategoryPacking)}
}

type AdditionalFeeType string

const (
	AdditionalFeeTypeMerchant AdditionalFeeType = "merchant" // 商户
	AdditionalFeeTypeStore    AdditionalFeeType = "store"    // 门店
)

func (AdditionalFeeType) Values() []string {
	return []string{string(AdditionalFeeTypeMerchant), string(AdditionalFeeTypeStore)}
}

type AdditionalFeeChargeMode string

const (
	AdditionalFeeChargeModePercent AdditionalFeeChargeMode = "percent" // 百分比
	AdditionalFeeChargeModeFixed   AdditionalFeeChargeMode = "fixed"   // 固定金额
)

func (AdditionalFeeChargeMode) Values() []string {
	return []string{string(AdditionalFeeChargeModePercent), string(AdditionalFeeChargeModeFixed)}
}

type AdditionalFeeDiscountScope string

const (
	AdditionalFeeDiscountScopeBefore AdditionalFeeDiscountScope = "before_discount" // 折前
	AdditionalFeeDiscountScopeAfter  AdditionalFeeDiscountScope = "after_discount"  // 折后
)

func (AdditionalFeeDiscountScope) Values() []string {
	return []string{string(AdditionalFeeDiscountScopeBefore), string(AdditionalFeeDiscountScopeAfter)}
}

type AdditionalFeeSimpleUpdateType string

const (
	AdditionalFeeSimpleUpdateTypeEnabled AdditionalFeeSimpleUpdateType = "enabled"
)

type AdditionalFeeOrderByType int

const (
	_ AdditionalFeeOrderByType = iota
	AdditionalFeeOrderByID
	AdditionalFeeOrderByCreatedAt
	AdditionalFeeOrderBySortOrder
)

type AdditionalFeeOrderBy struct {
	OrderBy AdditionalFeeOrderByType
	Desc    bool
}

func NewAdditionalFeeOrderByID(desc bool) AdditionalFeeOrderBy {
	return AdditionalFeeOrderBy{OrderBy: AdditionalFeeOrderByID, Desc: desc}
}

func NewAdditionalFeeOrderByCreatedAt(desc bool) AdditionalFeeOrderBy {
	return AdditionalFeeOrderBy{OrderBy: AdditionalFeeOrderByCreatedAt, Desc: desc}
}

func NewAdditionalFeeOrderBySortOrder(desc bool) AdditionalFeeOrderBy {
	return AdditionalFeeOrderBy{OrderBy: AdditionalFeeOrderBySortOrder, Desc: desc}
}

type AdditionalFee struct {
	ID                  uuid.UUID                  `json:"id"`
	Name                string                     `json:"name"`                  // 附加费名称
	FeeType             AdditionalFeeType          `json:"fee_type"`              // 附加费类型
	FeeCategory         AdditionalCategory         `json:"fee_category"`          // 附加费类别
	ChargeMode          AdditionalFeeChargeMode    `json:"charge_mode"`           // 费用类型 percent 百分比，fixed 固定金额
	FeeValue            decimal.Decimal            `json:"fee_value"`             // fixed: 分; percent: BP
	IncludeInReceivable bool                       `json:"include_in_receivable"` // 是否计入实收
	Taxable             bool                       `json:"taxable"`               // 附加费是否收税
	DiscountScope       AdditionalFeeDiscountScope `json:"discount_scope"`        // 折扣场景
	OrderChannels       []OrderChannel             `json:"order_channels"`        // 订单渠道
	DiningWays          []DiningWay                `json:"dining_ways"`           // 就餐方式
	Enabled             bool                       `json:"enabled"`               // 是否启用
	SortOrder           int                        `json:"sort_order"`            // 排序
	MerchantID          uuid.UUID                  `json:"merchant_id"`           // 品牌商 ID
	StoreID             uuid.UUID                  `json:"store_id"`              // 门店 ID
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
}

type AdditionalFeeListFilter struct {
	MerchantID  uuid.UUID          `json:"merchant_id"`
	StoreID     uuid.UUID          `json:"store_id"`
	Name        string             `json:"name"`         // 附加费名称，支持模糊查询
	FeeType     AdditionalFeeType  `json:"fee_type"`     // 附加费类型  商户/门店
	FeeCategory AdditionalCategory `json:"fee_category"` // 附加费类别  服务费/桌台费/打包费
	Enabled     *bool              `json:"enabled"`      // 是否启用
}

type AdditionalFeeExistsParams struct {
	MerchantID uuid.UUID `json:"merchant_id,omitempty"`
	StoreID    uuid.UUID `json:"store_id,omitempty"`
	Name       string    `json:"name,omitempty"`
	ExcludeID  uuid.UUID `json:"exclude_id,omitempty"`
}
