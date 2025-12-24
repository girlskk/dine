package types

// ProductUnitCreateReq 创建商品单位请求
type ProductUnitCreateReq struct {
	Name string `json:"name" binding:"required,max=255"`               // 单位名称
	Type string `json:"type" binding:"required,oneof=quantity weight"` // 单位类型：quantity（数量单位）、weight（重量单位）
}

// ProductUnitUpdateReq 更新商品单位请求
type ProductUnitUpdateReq struct {
	Name string `json:"name" binding:"required,max=255"`               // 单位名称
	Type string `json:"type" binding:"required,oneof=quantity weight"` // 单位类型：quantity（数量单位）、weight（重量单位）
}

// ProductUnitListReq 查询商品单位列表请求
type ProductUnitListReq struct {
	Name string `form:"name" binding:"omitempty,max=255"`               // 单位名称
	Type string `form:"type" binding:"omitempty,oneof=quantity weight"` // 单位类型：quantity（数量单位）、weight（重量单位）
	Page int    `form:"page" binding:"omitempty,min=1"`                 // 页码
	Size int    `form:"size" binding:"omitempty,min=1"`                 // 每页数量
}
