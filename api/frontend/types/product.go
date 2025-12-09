package types

import "gitlab.jiguang.dev/pos-dine/dine/domain"

// CategoryListReq 分类列表请求
type CategoryListReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

// ProductListReq 商品列表请求
type ProductListReq struct {
	Page       int                        `json:"page"`
	Size       int                        `json:"size"`
	CategoryID int                        `json:"category_id"`                                      // 分类ID
	SaleStatus []domain.ProductSaleStatus `json:"sale_status" binding:"omitempty,dive,oneof=1 2 3"` // 商品售卖状态： 1-在售 2-售罄 3-部分规格售罄
}

// ProductIDReq 商品ID请求
type ProductIDReq struct {
	ID int `json:"id"` // 商品ID
}

// 商品估清请求参数
type ProductClearStockReq struct {
	ProductID int   `json:"product_id"` // 商品ID
	SpecIDs   []int `json:"spec_ids"`   // 商品-规格ID列表（多规格商品必传）
}

// 取消估清操作
type ProductRestoreStockReq struct {
	ProductID int   `json:"product_id"` // 商品ID
	SpecIDs   []int `json:"spec_ids"`   // 商品-规格ID列表（多规格商品必传）
}
