package types

type CategoryListReq struct {
	StoreID int `json:"store_id"`
}

type ProductListReq struct {
	StoreID    int `json:"store_id"`    // 门店ID
	CategoryID int `json:"category_id"` // 分类ID
}

type ProductIDReq struct {
	ID int `json:"id"` // 商品ID
}
