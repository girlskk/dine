package types

// ProductTagCreateReq 创建商品标签请求
type ProductTagCreateReq struct {
	Name string `json:"name" binding:"required,max=255"` // 标签名称
}

// ProductTagUpdateReq 更新商品标签请求
type ProductTagUpdateReq struct {
	Name string `json:"name" binding:"required,max=255"` // 标签名称
}

// ProductTagListReq 查询商品标签列表请求
type ProductTagListReq struct {
	Name string `form:"name" binding:"omitempty,max=255"` // 标签名称
	Page int    `form:"page" binding:"omitempty,min=1"`   // 页码
	Size int    `form:"size" binding:"omitempty,min=1"`   // 每页数量
}
