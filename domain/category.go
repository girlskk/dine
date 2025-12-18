package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCategoryNotExists         = errors.New("商品分类不存在")
	ErrCategoryNameExists        = errors.New("商品分类名称已存在")
	ErrCategoryParentNotExists   = errors.New("父分类不存在")
	ErrCategoryParentHasProducts = errors.New("父分类下有商品，不能创建子分类")
	ErrCategoryInvalidLevel      = errors.New("分类级别无效，只支持两级分类")
	ErrCategoryDeleteHasChildren = errors.New("商品分类下有子分类，不能删除")
	ErrCategoryDeleteHasProducts = errors.New("商品分类下有商品，不能删除")
)

// CategoryRepository 商品分类仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/category_repository.go -package=mock . CategoryRepository
type CategoryRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Category, error)
	Create(ctx context.Context, category *Category) error
	CreateBulk(ctx context.Context, categories []*Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params CategoryExistsParams) (bool, error)
	CountChildrenByParentID(ctx context.Context, parentID uuid.UUID) (int, error)
}

// CategoryInteractor 商品分类用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/category_interactor.go -package=mock . CategoryInteractor
type CategoryInteractor interface {
	CreateRoot(ctx context.Context, category *Category) error
	CreateChild(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	// Update(ctx context.Context, category *Category) (*Category, error)
}

// Category 商品分类实体
type Category struct {
	ID             uuid.UUID `json:"id"`               // 分类ID
	Name           string    `json:"name"`             // 分类名称
	MerchantID     uuid.UUID `json:"merchant_id"`      // 品牌商ID
	StoreID        uuid.UUID `json:"store_id"`         // 门店ID
	ParentID       uuid.UUID `json:"parent_id"`        // 父分类ID
	InheritTaxRate bool      `json:"inherit_tax_rate"` // 是否继承父分类的税率ID
	TaxRateID      uuid.UUID `json:"tax_rate_id"`      // 税率ID
	InheritStall   bool      `json:"inherit_stall"`    // 是否继承父分类的出品部门ID
	StallID        uuid.UUID `json:"stall_id"`         // 出品部门ID
	ProductCount   int       `json:"product_count"`    // 关联的商品数量
	SortOrder      int       `json:"sort_order"`       // 排序，值越小越靠前
	CreatedAt      time.Time `json:"created_at"`       // 创建时间
	UpdatedAt      time.Time `json:"updated_at"`       // 更新时间

	// @TODO 关联信息
	Childrens []*Category `json:"children,omitempty"` // 子分类列表
	// TaxRate *TaxRate `json:"tax_rate,omitempty"` // 税率
	// Stall   *Stall   `json:"stall,omitempty"`    // 出品部门
}

// Categories 商品分类集合
type Categories []*Category

// IsRoot 判断是否为一级分类
func (c *Category) IsRoot() bool {
	return c.ParentID == uuid.Nil
}

// ------------------------------------------------------------
// 参数定义（DTO）
// ------------------------------------------------------------

// CategoryExistsParams 存在性检查参数
type CategoryExistsParams struct {
	MerchantID uuid.UUID
	Name       string
	ParentID   uuid.UUID
	IsRoot     bool
}
