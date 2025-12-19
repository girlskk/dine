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
	ErrProductAttrNotExists             = errors.New("商品口味做法不存在")
	ErrProductAttrNameExists            = errors.New("商品口味做法名称已存在")
	ErrProductAttrDeleteHasItems        = errors.New("商品口味做法下有口味做法项，不能删除")
	ErrProductAttrItemNotExists         = errors.New("商品口味做法项不存在")
	ErrProductAttrItemNameExists        = errors.New("商品口味做法项名称已存在")
	ErrProductAttrItemDeleteHasProducts = errors.New("商品口味做法项下有商品，不能删除")
	ErrProductAttrItemBasePriceInvalid  = errors.New("基础加价必须为非负数")
)

// ProductAttrRepository 商品口味做法仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_attr_repository.go -package=mock . ProductAttrRepository
type ProductAttrRepository interface {
	// ProductAttr 相关操作
	FindByID(ctx context.Context, id uuid.UUID) (*ProductAttr, error)
	Create(ctx context.Context, attr *ProductAttr) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params ProductAttrExistsParams) (bool, error)
	ListBySearch(ctx context.Context, params ProductAttrSearchParams) (ProductAttrs, error)
	GetDetail(ctx context.Context, id uuid.UUID) (*ProductAttr, error)

	// ProductAttrItem 相关操作（作为 ProductAttr 的一部分）
	FindItemByID(ctx context.Context, id uuid.UUID) (*ProductAttrItem, error)
	CreateItems(ctx context.Context, items []*ProductAttrItem) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
}

// ProductAttrInteractor 商品口味做法用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_attr_interactor.go -package=mock . ProductAttrInteractor
type ProductAttrInteractor interface {
	Create(ctx context.Context, attr *ProductAttr) error
	Update(ctx context.Context, attr *ProductAttr) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
	ListBySearch(ctx context.Context, params ProductAttrSearchParams) (ProductAttrs, error)
}

// 售卖渠道枚举
type SaleChannel string

const (
	SaleChannelPOS                SaleChannel = "POS"         // POS
	SaleChannelMobileOrdering     SaleChannel = "Mobile"      // 移动点餐
	SaleChannelScanOrdering       SaleChannel = "Scan"        // 扫码点餐
	SaleChannelSelfService        SaleChannel = "SelfService" // 自助点餐
	SaleChannelThirdPartyDelivery SaleChannel = "ThirdParty"  // 三方外卖
)

func (SaleChannel) Values() []string {
	return []string{
		string(SaleChannelPOS),
		string(SaleChannelMobileOrdering),
		string(SaleChannelScanOrdering),
		string(SaleChannelSelfService),
		string(SaleChannelThirdPartyDelivery),
	}
}

// ProductAttr 商品口味做法实体
type ProductAttr struct {
	ID           uuid.UUID     `json:"id"`            // 口味做法ID
	Name         string        `json:"name"`          // 口味做法名称
	Channels     []SaleChannel `json:"channels"`      // 售卖渠道列表
	MerchantID   uuid.UUID     `json:"merchant_id"`   // 品牌商ID
	StoreID      uuid.UUID     `json:"store_id"`      // 门店ID
	ProductCount int           `json:"product_count"` // 关联的商品数量（所有项的关联商品数量累加）
	CreatedAt    time.Time     `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time     `json:"updated_at"`    // 更新时间

	// 关联信息
	Items []*ProductAttrItem `json:"items,omitempty"` // 口味做法项列表
}

// ProductAttrs 商品口味做法集合
type ProductAttrs []*ProductAttr

// ProductAttrItem 商品口味做法项实体
type ProductAttrItem struct {
	ID           uuid.UUID       `json:"id"`            // 口味做法项ID
	AttrID       uuid.UUID       `json:"attr_id"`       // 口味做法ID（外键）
	Name         string          `json:"name"`          // 口味做法项名称
	Image        string          `json:"image"`         // 图片URL（可选）
	BasePrice    decimal.Decimal `json:"base_price"`    // 基础加价（单位：分）
	ProductCount int             `json:"product_count"` // 关联的商品数量
	CreatedAt    time.Time       `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time       `json:"updated_at"`    // 更新时间
}

// ProductAttrItems 商品口味做法项集合
type ProductAttrItems []*ProductAttrItem

// ------------------------------------------------------------
// 参数定义（DTO）
// ------------------------------------------------------------

// ProductAttrExistsParams 存在性检查参数
type ProductAttrExistsParams struct {
	MerchantID uuid.UUID
	Name       string
	ExcludeID  uuid.UUID // 排除的ID（用于更新时检查名称唯一性）
}

// ProductAttrItemExistsParams 存在性检查参数
type ProductAttrItemExistsParams struct {
	AttrID    uuid.UUID
	Name      string
	ExcludeID uuid.UUID // 排除的ID（用于更新时检查名称唯一性）
}

// ProductAttrSearchParams 查询参数
type ProductAttrSearchParams struct {
	MerchantID uuid.UUID
}

type ProductAttrSearchRes struct {
	*upagination.Pagination
	Items ProductAttrs `json:"items"`
}
