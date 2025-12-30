package types

// ProductSpecCreateReq 创建商品规格请求
type ProductSpecCreateReq struct {
	Name string `json:"name" binding:"required,max=255"` // 规格名称
}

// ProductSpecUpdateReq 更新商品规格请求
type ProductSpecUpdateReq struct {
	Name string `json:"name" binding:"required,max=255"` // 规格名称
}

// ProductSpecListReq 查询商品规格列表请求
type ProductSpecListReq struct {
	Name string `form:"name" binding:"omitempty,max=255"` // 规格名称
	Page int    `form:"page" binding:"omitempty,min=1"`   // 页码
	Size int    `form:"size" binding:"omitempty,min=1"`   // 每页数量
}
