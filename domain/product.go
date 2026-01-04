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
	// 商品基础错误
	ErrProductNotExists  = errors.New("商品不存在")
	ErrProductNameExists = errors.New("商品名称已存在")

	// 商品分类相关错误
	ErrProductCategoryNotExists = errors.New("商品分类不存在")
	ErrProductCategoryInvalid   = errors.New("商品分类无效")

	// 商品单位相关错误
	ErrProductUnitInvalid = errors.New("商品单位无效")

	// 商品规格相关错误
	ErrProductSpecInvalid = errors.New("商品规格无效")

	ErrProductSpecRelationNoDefault = errors.New("商品规格必须且只有一个默认项")

	// 商品口味做法相关错误
	ErrProductAttrInvalid           = errors.New("商品口味做法无效")
	ErrProductAttrRelationNoDefault = errors.New("每个口味做法分组必须且只有一个默认项")

	// 商品标签相关错误
	ErrProductTagInvalid = errors.New("商品标签无效")

	// 商品验证错误
	ErrProductEffectiveDateInvalid = errors.New("生效日期无效，开始时间必须早于结束时间")
	ErrProductTaxRateNotExists     = errors.New("指定税率不存在")
	ErrProductStallNotExists       = errors.New("指定出品部门不存在")
	ErrProductPackingFeeNotExists  = errors.New("打包费配置不存在")

	// 商品套餐相关错误
	ErrProductBelongToSetMeal = errors.New("商品属于套餐组，不能停售或删除")

	// 商品下发相关错误
	ErrProductDistributeStoreInvalid = errors.New("门店无效，必须属于当前品牌商")
)

// ------------------------------------------------------------
// 仓储和用例接口
// ------------------------------------------------------------

// ProductRepository 商品仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_repository.go -package=mock . ProductRepository
type ProductRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetDetail(ctx context.Context, id uuid.UUID) (*Product, error)
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params ProductExistsParams) (bool, error)
	ListByIDs(ctx context.Context, ids []uuid.UUID) (Products, error)
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProductSearchParams) (*ProductSearchRes, error)
	FindByNameInStore(ctx context.Context, storeID uuid.UUID, name string) (*Product, error)
}

// ProductInteractor 商品用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_interactor.go -package=mock . ProductInteractor
type ProductInteractor interface {
	Create(ctx context.Context, product *Product, user User) error
	CreateSetMeal(ctx context.Context, product *Product, user User) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProductSearchParams) (*ProductSearchRes, error)
	Update(ctx context.Context, product *Product, user User) error
	UpdateSetMeal(ctx context.Context, product *Product, user User) error
	Delete(ctx context.Context, id uuid.UUID, user User) error
	OffSale(ctx context.Context, id uuid.UUID, user User) error
	OnSale(ctx context.Context, id uuid.UUID, user User) error
	GetDetail(ctx context.Context, id uuid.UUID, user User) (*Product, error)
	Distribute(ctx context.Context, params ProductDistributeParams, user User) error
}

// ------------------------------------------------------------
// 枚举定义
// ------------------------------------------------------------

// ProductSupportType 支持类型
type ProductSupportType string

const (
	ProductSupportTypeDine     ProductSupportType = "dine"     // 堂食
	ProductSupportTypeTakeaway ProductSupportType = "takeaway" // 外带
	ProductSupportTypeDelivery ProductSupportType = "delivery" // 外卖
)

func (ProductSupportType) Values() []string {
	return []string{
		string(ProductSupportTypeDine),
		string(ProductSupportTypeTakeaway),
		string(ProductSupportTypeDelivery),
	}
}

// ProductAttrSelectionType 口味做法点单限制
type ProductAttrSelectionType string

const (
	ProductAttrSelectionTypeRequiredOne ProductAttrSelectionType = "required_one" // 必选一项
	ProductAttrSelectionTypeMultiple    ProductAttrSelectionType = "multiple"     // 可多选
)

func (ProductAttrSelectionType) Values() []string {
	return []string{
		string(ProductAttrSelectionTypeRequiredOne),
		string(ProductAttrSelectionTypeMultiple),
	}
}

// ProductSaleStatus 售卖状态
type ProductSaleStatus string

const (
	ProductSaleStatusOnSale  ProductSaleStatus = "on_sale"  // 在售
	ProductSaleStatusOffSale ProductSaleStatus = "off_sale" // 停售
)

func (ProductSaleStatus) Values() []string {
	return []string{
		string(ProductSaleStatusOnSale),
		string(ProductSaleStatusOffSale),
	}
}

// EffectiveDateType 生效日期类型
type EffectiveDateType string

const (
	EffectiveDateTypeDaily  EffectiveDateType = "daily"  // 按天
	EffectiveDateTypeCustom EffectiveDateType = "custom" // 自定义
)

func (EffectiveDateType) Values() []string {
	return []string{
		string(EffectiveDateTypeDaily),
		string(EffectiveDateTypeCustom),
	}
}

// ProductType 商品类型
type ProductType string

const (
	ProductTypeNormal  ProductType = "normal"   // 普通商品
	ProductTypeSetMeal ProductType = "set_meal" // 套餐商品
)

func (ProductType) Values() []string {
	return []string{
		string(ProductTypeNormal),
		string(ProductTypeSetMeal),
	}
}

// ------------------------------------------------------------
// 实体定义
// ------------------------------------------------------------

// Product 商品实体
type Product struct {
	ID         uuid.UUID `json:"id"`          // 商品ID
	Name       string    `json:"name"`        // 商品名称
	MerchantID uuid.UUID `json:"merchant_id"` // 品牌商ID
	StoreID    uuid.UUID `json:"store_id"`    // 门店ID

	// 商品类型
	Type ProductType `json:"type"` // 商品类型：normal（普通商品）、set_meal（套餐商品）

	// 基础信息
	CategoryID   uuid.UUID            `json:"category_id"`   // 分类ID
	MenuID       uuid.UUID            `json:"menu_id"`       // 菜单ID
	Mnemonic     string               `json:"mnemonic"`      // 助记词
	ShelfLife    int                  `json:"shelf_life"`    // 保质期（单位：天）
	SupportTypes []ProductSupportType `json:"support_types"` // 支持类型（堂食、外带）

	// 属性关联
	UnitID uuid.UUID `json:"unit_id"` // 单位ID

	// 售卖信息
	SaleStatus         ProductSaleStatus `json:"sale_status"`          // 售卖状态
	SaleChannels       []SaleChannel     `json:"sale_channels"`        // 售卖渠道
	EffectiveDateType  EffectiveDateType `json:"effective_date_type"`  // 生效日期类型
	EffectiveStartTime *time.Time        `json:"effective_start_time"` // 生效开始时间
	EffectiveEndTime   *time.Time        `json:"effective_end_time"`   // 生效结束时间
	MinSaleQuantity    int               `json:"min_sale_quantity"`    // 起售份数
	AddSaleQuantity    int               `json:"add_sale_quantity"`    // 加售份数

	// 其他信息
	InheritTaxRate bool      `json:"inherit_tax_rate"` // 是否继承原分类税率
	TaxRateID      uuid.UUID `json:"tax_rate_id"`      // 指定税率ID
	InheritStall   bool      `json:"inherit_stall"`    // 是否继承原出品部门
	StallID        uuid.UUID `json:"stall_id"`         // 指定出品部门ID

	// 展示信息
	MainImage    string   `json:"main_image"`    // 主图
	DetailImages []string `json:"detail_images"` // 详情图片
	Description  string   `json:"description"`   // 菜品描述

	// 套餐属性（仅套餐商品使用）
	EstimatedCostPrice *decimal.Decimal `json:"estimated_cost_price,omitempty"` // 预估成本价（可选，单位：分，仅套餐商品使用）
	DeliveryCostPrice  *decimal.Decimal `json:"delivery_cost_price,omitempty"`  // 外卖成本价（可选，单位：分，仅套餐商品使用）

	// 时间戳
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间

	// 关联信息
	SpecRelations ProductSpecRelations `json:"spec_relations,omitempty"` // 商品规格关联列表
	AttrRelations ProductAttrRelations `json:"attr_relations,omitempty"` // 商品口味做法关联列表
	Tags          ProductTags          `json:"tags,omitempty"`           // 商品标签列表
	Groups        SetMealGroups        `json:"groups,omitempty"`         // 套餐组列表
	Category      *Category            `json:"category,omitempty"`       // 分类
	Unit          *ProductUnit         `json:"unit,omitempty"`           // 单位
	TaxRate       *TaxFee              `json:"tax_rate,omitempty"`       // 税率
	Stall         *Stall               `json:"stall,omitempty"`          // 出品部门
}

// Products 商品集合
type Products []*Product

// ------------------------------------------------------------
// 参数定义（DTO）
// ------------------------------------------------------------

// ProductExistsParams 存在性检查参数
type ProductExistsParams struct {
	MerchantID uuid.UUID
	StoreID    uuid.UUID
	Name       string
	ExcludeID  uuid.UUID // 排除的ID（用于更新时检查名称唯一性）
}

// ProductSearchParams 查询参数
type ProductSearchParams struct {
	MerchantID   uuid.UUID         // 品牌商ID（必填）
	StoreID      uuid.UUID         // 门店ID（可选）
	OnlyMerchant bool              // 是否只查询品牌商ID（可选）
	Name         string            // 商品名称（可选，模糊匹配）
	SaleStatus   ProductSaleStatus // 售卖状态（可选，空字符串表示全部）
	Type         ProductType       // 商品类型（可选，空字符串表示全部）
	CategoryID   uuid.UUID         // 分类ID（可选，支持一级分类和二级分类）
	StallID      uuid.UUID         // 出品部门ID（可选）
	StartAt      *time.Time        // 创建开始时间（可选）
	EndAt        *time.Time        // 创建结束时间（可选，最长支持1年跨度）
}

// ProductSearchRes 查询结果
type ProductSearchRes struct {
	*upagination.Pagination
	Items Products `json:"items"`
}

// ProductDistributeParams 商品下发参数
type ProductDistributeParams struct {
	ProductID        uuid.UUID            // 商品ID（必选）
	MerchantID       uuid.UUID            // 品牌商ID
	StoreIDs         []uuid.UUID          // 门店ID列表（必选，多选）
	DistributionRule MenuDistributionRule // 下发规则（必选）：override（新增并覆盖同名菜品）、keep（对同名菜品不做修改）
	SaleRule         MenuItemSaleRule     // 下发售卖规则（可选，仅当下发规则为override时使用）：keep_brand_status（保留品牌状态）、keep_store_status（保留门店状态）
}
