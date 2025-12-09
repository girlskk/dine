package types

import (
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

type CategoryListReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type CategoryCreateReq struct {
	Name string `json:"name" binding:"required"` // 分类名称
}

type CategoryUpdateReq struct {
	ID   int    `json:"id" binding:"required"`   // 分类ID
	Name string `json:"name" binding:"required"` // 分类名称
}

type CategoryDeleteReq struct {
	ID int `json:"id" binding:"required"` // 分类ID
}

type UnitCreateReq struct {
	Name string `json:"name" binding:"required"` // 单位名称
}

type UnitUpdateReq struct {
	ID   int    `json:"id" binding:"required"`   // 单位ID
	Name string `json:"name" binding:"required"` // 单位名称
}

type UnitDeleteReq struct {
	ID int `json:"id" binding:"required"` // 单位ID
}

type UnitListReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type AttrCreateReq struct {
	Name string `json:"name" binding:"required"` // 属性名称
}

type AttrUpdateReq struct {
	ID   int    `json:"id" binding:"required"`   // 属性ID
	Name string `json:"name" binding:"required"` // 属性名称
}

type AttrDeleteReq struct {
	ID int `json:"id" binding:"required"` // 属性ID
}

type AttrListReq struct {
	Page int `json:"page"` // 页码
	Size int `json:"size"` // 每页数量
}

type RecipeCreateReq struct {
	Name string `json:"name" binding:"required"` // 做法名称
}

type RecipeUpdateReq struct {
	ID   int    `json:"id" binding:"required"`   // 做法ID
	Name string `json:"name" binding:"required"` // 做法名称
}

type RecipeDeleteReq struct {
	ID int `json:"id" binding:"required"` // 做法ID
}

type RecipeListReq struct {
	Page int `json:"page"` // 页码
	Size int `json:"size"` // 每页数量
}

type SpecCreateReq struct {
	Name string `json:"name" binding:"required"` // 规格名称
}

type SpecUpdateReq struct {
	ID   int    `json:"id" binding:"required"`   // 规格ID
	Name string `json:"name" binding:"required"` // 规格名称
}

type SpecDeleteReq struct {
	ID int `json:"id" binding:"required"` // 规格ID
}

type SpecListReq struct {
	Page int `json:"page"` // 页码
	Size int `json:"size"` // 每页数量
}

type ProductCreateReq struct {
	Name         string             `json:"name" binding:"required"`             // 商品名称
	Type         domain.ProductType `json:"type" binding:"required,oneof=1 2 3"` // 商品类型 1-单规格 2-多规格 3-套餐
	Price        decimal.Decimal    `json:"price"`                               // 单规格价格
	Images       []string           `json:"images" binding:"required"`           // 商品图片
	CategoryID   int                `json:"category_id" binding:"required"`      // 商品分类ID
	UnitID       int                `json:"unit_id"`                             // 商品单位ID
	AttrIDs      []int              `json:"attr_ids"`                            // 商品属性ID
	RecipeIDs    []int              `json:"recipe_ids"`                          // 商品做法ID
	Specs        []ProductSpecReq   `json:"specs"`                               // 商品规格（多规格时必填）
	SetMealItems []SetMealItemReq   `json:"setmeal_items"`                       // 套餐商品详情（套餐商品时必填）
}

type ProductSpecReq struct {
	ID    int             `json:"id"`    // 规格ID
	Price decimal.Decimal `json:"price"` // 规格价格
}

type SetMealItemReq struct {
	ProductID     int             `json:"product_id" binding:"required"` // 子商品ID
	ProductSpecID int             `json:"product_spec_id"`               // 商品-规格ID（多规格商品必填）
	Quantity      decimal.Decimal `json:"quantity" binding:"required"`   // 数量
	Price         decimal.Decimal `json:"price" binding:"required"`      // 套餐内单价
}

type ProductUpdateReq struct {
	ID           int              `json:"id" binding:"required"`
	Name         string           `json:"name" binding:"required"`        // 商品名称
	Price        decimal.Decimal  `json:"price"`                          // 单规格价格
	Images       []string         `json:"images" binding:"required"`      // 商品图片
	CategoryID   int              `json:"category_id" binding:"required"` // 商品分类ID
	UnitID       int              `json:"unit_id"`                        // 商品单位ID
	AttrIDs      []int            `json:"attr_ids"`                       // 商品属性ID
	RecipeIDs    []int            `json:"recipe_ids"`                     // 商品做法ID
	Specs        []ProductSpecReq `json:"specs"`                          // 商品规格（多规格时必填）
	SetMealItems []SetMealItemReq `json:"setmeal_items"`                  // 套餐商品详情（套餐商品时必填）
}

type ProductDeleteReq struct {
	ID int `json:"id" binding:"required"`
}

type ProductDetailReq struct {
	ID int `json:"id" binding:"required"`
}

type ProductListReq struct {
	Page       int                        `json:"page"`
	Size       int                        `json:"size"`
	Name       string                     `json:"name"`                                             // 名称或编号
	CategoryID int                        `json:"category_id"`                                      // 分类ID
	Status     domain.ProductStatus       `json:"status" binding:"omitempty,oneof=1 2"`             // 商品状态：1-待审核 2-审核通过
	SaleStatus []domain.ProductSaleStatus `json:"sale_status" binding:"omitempty,dive,oneof=1 2 3"` // 商品售卖状态： 1-在售 2-售罄 3-部分规格售罄
	Type       domain.ProductType         `json:"type" binding:"omitempty,oneof=1 2 3"`             // 商品类型：1-单规格 2-多规格 3-套餐
}

// 创建/编辑商品转换为 usecase 请求参数
func ToProductUpsetReq(req ProductCreateReq, productID int) domain.ProductUpsetParams {
	product := &domain.Product{
		ID:            productID,
		Name:          req.Name,
		Type:          req.Type,
		Price:         req.Price,
		Images:        req.Images,
		AllowPointPay: true,
		CategoryID:    req.CategoryID,
		UnitID:        req.UnitID,
		Specs:         make(domain.ProductSpecRels, 0, len(req.Specs)),
	}
	// 商品规格
	for _, item := range req.Specs {
		product.Specs = append(product.Specs, &domain.ProductSpecRel{
			SpecID:     item.ID,
			Price:      item.Price,
			SaleStatus: domain.ProductSaleStatusOn,
		})
	}

	setmealDetails := make(domain.SetMealDetails, 0, len(req.SetMealItems))
	for _, item := range req.SetMealItems {
		setmealDetails = append(setmealDetails, &domain.SetMealDetail{
			SetMealPrice: item.Price,
			Quantity:     item.Quantity,
			ProductID:    item.ProductID,
			Spec: &domain.ProductSpecRel{
				ID: item.ProductSpecID,
			},
		})
	}

	return domain.ProductUpsetParams{
		Product:        product,
		AttrIDs:        req.AttrIDs,
		RecipeIDs:      req.RecipeIDs,
		SetMealDetails: setmealDetails,
	}
}

func ToProductUpdateReq(req ProductUpdateReq) domain.ProductUpsetParams {
	product := &domain.Product{
		ID:            req.ID,
		Name:          req.Name,
		Price:         req.Price,
		Images:        req.Images,
		AllowPointPay: true,
		CategoryID:    req.CategoryID,
		UnitID:        req.UnitID,
		Specs:         make(domain.ProductSpecRels, 0, len(req.Specs)),
	}
	// 商品规格
	for _, item := range req.Specs {
		product.Specs = append(product.Specs, &domain.ProductSpecRel{
			SpecID:     item.ID,
			Price:      item.Price,
			SaleStatus: domain.ProductSaleStatusOn,
		})
	}

	setmealDetails := make(domain.SetMealDetails, 0, len(req.SetMealItems))
	for _, item := range req.SetMealItems {
		setmealDetails = append(setmealDetails, &domain.SetMealDetail{
			SetMealPrice: item.Price,
			Quantity:     item.Quantity,
			ProductID:    item.ProductID,
			Spec: &domain.ProductSpecRel{
				ID: item.ProductSpecID,
			},
		})
	}

	return domain.ProductUpsetParams{
		Product:        product,
		AttrIDs:        req.AttrIDs,
		RecipeIDs:      req.RecipeIDs,
		SetMealDetails: setmealDetails,
	}
}
