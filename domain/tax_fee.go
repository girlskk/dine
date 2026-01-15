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
	ErrTaxFeeNotExists          = errors.New("税费不存在")
	ErrTaxFeeNameExists         = errors.New("税费名称已存在")
	ErrTaxFeeSystemCannotUpdate = errors.New("系统内置税费不能修改")
	ErrTaxFeeSystemCannotDelete = errors.New("默认税费不能删除")
)

// TaxFeeRepository 税费仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/tax_fee_repository.go -package=mock . TaxFeeRepository
type TaxFeeRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (fee *TaxFee, err error)
	Create(ctx context.Context, fee *TaxFee) (err error)
	Update(ctx context.Context, fee *TaxFee) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetTaxFees(ctx context.Context, pager *upagination.Pagination, filter *TaxFeeListFilter, orderBys ...TaxFeeOrderBy) (fees []*TaxFee, total int, err error)
	Exists(ctx context.Context, params TaxFeeExistsParams) (exists bool, err error)
}

// TaxFeeInteractor 税费用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/tax_fee_interactor.go -package=mock . TaxFeeInteractor
type TaxFeeInteractor interface {
	Create(ctx context.Context, fee *TaxFee, user User) (err error)
	Update(ctx context.Context, fee *TaxFee, user User) (err error)
	Delete(ctx context.Context, id uuid.UUID, user User) (err error)
	GetTaxFee(ctx context.Context, id uuid.UUID, user User) (fee *TaxFee, err error)
	GetTaxFees(ctx context.Context, pager *upagination.Pagination, filter *TaxFeeListFilter, orderBys ...TaxFeeOrderBy) (fees []*TaxFee, total int, err error)
	TaxFeeSimpleUpdate(ctx context.Context, updateField TaxFeeSimpleUpdateField, fee *TaxFee, user User) (err error)
}
type TaxFeeType string

const (
	TaxFeeTypeSystem   TaxFeeType = "system"   // 系统内置
	TaxFeeTypeMerchant TaxFeeType = "merchant" // 商户
	TaxFeeTypeStore    TaxFeeType = "store"    // 门店
)

func (TaxFeeType) Values() []string {
	return []string{
		string(TaxFeeTypeSystem),
		string(TaxFeeTypeMerchant),
		string(TaxFeeTypeStore),
	}
}

type TaxRateType string

const (
	TaxRateTypeUnified TaxRateType = "unified" // 统一比例
	TaxRateTypeCustom  TaxRateType = "custom"  // 自定义比例
)

func (TaxRateType) Values() []string {
	return []string{string(TaxRateTypeUnified), string(TaxRateTypeCustom)}
}

// TaxFeeSimpleUpdateField 简单更新字段
type TaxFeeSimpleUpdateField string

const (
	TaxFeeSimpleUpdateFieldDefault TaxFeeSimpleUpdateField = "default_tax"
)

type TaxFeeOrderByType int

const (
	_ TaxFeeOrderByType = iota
	TaxFeeOrderByID
	TaxFeeOrderByCreatedAt
)

type TaxFeeOrderBy struct {
	OrderBy TaxFeeOrderByType
	Desc    bool
}

func NewTaxFeeOrderByID(desc bool) TaxFeeOrderBy {
	return TaxFeeOrderBy{OrderBy: TaxFeeOrderByID, Desc: desc}
}

func NewTaxFeeOrderByCreatedAt(desc bool) TaxFeeOrderBy {
	return TaxFeeOrderBy{OrderBy: TaxFeeOrderByCreatedAt, Desc: desc}
}

// TaxFee 税费实体
type TaxFee struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`          // 税费名称
	TaxFeeType  TaxFeeType      `json:"tax_fee_type"`  // 税费类型：商户/门店
	TaxCode     string          `json:"tax_code"`      // 税费代码
	TaxRateType TaxRateType     `json:"tax_rate_type"` // 税率类型：统一比例/自定义比例
	TaxRate     decimal.Decimal `json:"tax_rate"`      // 税率，6% -> 0.06
	DefaultTax  bool            `json:"default_tax"`   // 是否默认税费
	MerchantID  uuid.UUID       `json:"merchant_id"`   // 品牌商 ID
	StoreID     uuid.UUID       `json:"store_id"`      // 门店 ID
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// TaxFeeListFilter 查询过滤条件
type TaxFeeListFilter struct {
	MerchantID uuid.UUID  `json:"merchant_id"`
	StoreID    uuid.UUID  `json:"store_id"`
	Name       string     `json:"name"`         // 税费名称，支持模糊查询
	TaxFeeType TaxFeeType `json:"tax_fee_type"` // 税费类型：商户/门店
}

// TaxFeeExistsParams 存在性检查参数
type TaxFeeExistsParams struct {
	MerchantID uuid.UUID `json:"merchant_id,omitempty"`
	StoreID    uuid.UUID `json:"store_id,omitempty"`
	Name       string    `json:"name,omitempty"`
	ExcludeID  uuid.UUID `json:"exclude_id,omitempty"`
}
