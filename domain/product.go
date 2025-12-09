package domain

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// ProductRepository 仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_repository.go -package=mock . ProductRepository
type ProductRepository interface {
	FindByID(ctx context.Context, id int) (*Product, error)          // 主表查询
	GetDetail(ctx context.Context, id int) (*Product, error)         // 详细信息，包含所有的关联数据
	GetDetailWithSpec(ctx context.Context, id int) (*Product, error) // 仅加载商品和规格信息
	Exists(ctx context.Context, params ProductExistsParams) (bool, error)
	Create(ctx context.Context, product *Product, attrIDs, recipeIDs []int) error
	Update(ctx context.Context, product *Product, attrIDs, recipeIDs []int) error  // 商品整体更新
	UpdateAttr(ctx context.Context, id int, attr ProductUpdateAttrs) error         // 商品部分字段更新
	BatchUpdateAttr(ctx context.Context, ids []int, attr ProductUpdateAttrs) error // 批量更新商品部分字段
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProductSearchParams) (*ProductSearchRes, error)
	ListByIDs(ctx context.Context, ids []int) (Products, error)
	GetDetailsByIDs(ctx context.Context, ids []int) (Products, error)
}

// ProductInteractor 用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_interactor.go -package=mock . ProductInteractor
type ProductInteractor interface {
	Create(ctx context.Context, params ProductUpsetParams) error
	Update(ctx context.Context, params ProductUpsetParams) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProductSearchParams) (*ProductSearchRes, error)
	GetDetail(ctx context.Context, id int) (*Product, error)
	Approve(ctx context.Context, ids []int, allowPointPay *bool) error
	UnApprove(ctx context.Context, ids []int) error
	ListSetmealDetails(ctx context.Context, id int) (SetMealDetails, error)
	ClearStock(ctx context.Context, id int, specIDs []int) error   // 商品估清
	RestoreStock(ctx context.Context, id int, specIDs []int) error // 取消估清
}

var (
	ErrProductNameExists       = errors.New("商品名称已存在")
	ErrProductNotExists        = errors.New("商品不存在")
	ErrProductStatus           = errors.New("商品状态错误")
	ErrProductNeedUnit         = errors.New("商品必须设置单位")
	ErrProductTypeMulti        = errors.New("多规格商品必须传递规格参数")
	ErrSetMealProductNotExists = errors.New("套餐必须传递套餐明细")
	ErrSetMealProductInvalid   = errors.New("套餐包含无效商品")
	ErrProductNotOnSale        = errors.New("商品不可售")
)

type ProductStatus int // 商品状态：1-待审核 2-审核通过

const (
	_                      ProductStatus = iota
	ProductStatusUnApprove               // 待审核
	ProductStatusApproved                // 审核通过
)

type ProductSaleStatus int // 商品售卖状态： 1-在售 2-售罄 3-部分规格售罄

const (
	_                        ProductSaleStatus = iota
	ProductSaleStatusOn                        // 在售
	ProductSaleStatusOff                       // 售罄
	ProductSaleStatusPartOff                   // 部分规格售罄（对于多规格商品）
)

type ProductType int // 商品类型：1-单规格 2-多规格 3-套餐商品

const (
	_                  ProductType = iota
	ProductTypeSingle              // 单规格
	ProductTypeMulti               // 多规格
	ProductTypeSetMeal             // 套餐商品
)

func (pt ProductType) String() string {
	switch pt {
	case ProductTypeSingle:
		return "单规格"
	case ProductTypeMulti:
		return "多规格"
	case ProductTypeSetMeal:
		return "套餐商品"
	default:
		return "未知状态"
	}
}

// Product 商品实体
type Product struct {
	ID            int               `json:"id"`
	Name          string            `json:"name"`            // 商品名称
	Type          ProductType       `json:"type"`            // 商品类型：1-单规格 2-多规格 3-套餐商品
	Price         decimal.Decimal   `json:"price"`           // 商品价格
	Status        ProductStatus     `json:"status"`          // 商品状态：1-待审核 2-审核通过
	SaleStatus    ProductSaleStatus `json:"sale_status"`     // 商品售卖状态：1-在售 2-售罄 3-部分规格售罄
	StoreID       int               `json:"store_id"`        // 所属门店ID
	Images        []string          `json:"images"`          // 商品图片
	AllowPointPay bool              `json:"allow_point_pay"` // 是否支持积分支付
	CreatedAt     time.Time         `json:"created_at"`      // 创建时间
	UpdatedAt     time.Time         `json:"updated_at"`      // 更新时间
	CategoryID    int               `json:"category_id"`     // 分类ID
	UnitID        int               `json:"unit_id"`         // 单位ID

	// 关联信息
	Specs          ProductSpecRels `json:"specs"`            // 商品规格
	Attrs          ProductAttrs    `json:"attrs"`            // 商品属性
	Recipes        ProductRecipes  `json:"recipes"`          // 商品做法
	Category       *Category       `json:"category"`         // 商品分类
	Unit           *ProductUnit    `json:"unit"`             // 商品单位
	SetMealDetails SetMealDetails  `json:"set_meal_details"` // 套餐详情
}

// Products 商品集合
type Products []*Product

// ProductSearchParams 查询参数
type ProductSearchParams struct {
	StoreID    int
	CategoryID int
	Name       string
	Status     ProductStatus
	SaleStatus []ProductSaleStatus
	Type       ProductType
}

type ProductSearchRes struct {
	*upagination.Pagination
	Items Products `json:"items"`
}

// ProductExistsParams 存在性检查参数
type ProductExistsParams struct {
	StoreID    int
	Name       string
	CategoryID int
	UnitID     int
}

// 创建/编辑商品参数
type ProductUpsetParams struct {
	Product        *Product
	AttrIDs        []int
	RecipeIDs      []int
	SetMealDetails SetMealDetails
}

// 商品字段更新
type ProductUpdateAttrs struct {
	Status        ProductStatus
	SaleStatus    ProductSaleStatus
	AllowPointPay *bool
}

func (p *Product) CheckSale() bool {
	if p.Status != ProductStatusApproved {
		return false
	}

	if p.SaleStatus == ProductSaleStatusOff {
		return false
	}

	return true
}
