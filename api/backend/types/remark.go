package types

import (
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// RemarkCreateReq 创建备注请求
// 名称、分类、类型为必填；排序为可选；Enabled 默认 true
// MerchantID 从登录用户上下文取；StoreID 由前端传入（按需）
type RemarkCreateReq struct {
	Name       string    `json:"name" binding:"required,max=255"`      // 备注名称
	Enabled    bool      `json:"enabled"`                              // 是否启用
	SortOrder  int       `json:"sort_order" binding:"omitempty,gte=0"` // 排序，越小越靠前
	CategoryID uuid.UUID `json:"category_id" binding:"required"`       // 备注类型
	StoreID    uuid.UUID `json:"store_id"`                             // 可选，品牌级可为空
}

// RemarkUpdateReq 更新备注请求
type RemarkUpdateReq struct {
	Name      string `json:"name" binding:"required,max=255"`      // 备注名称
	Enabled   bool   `json:"enabled"`                              // 是否启用
	SortOrder int    `json:"sort_order" binding:"omitempty,gte=0"` // 排序，越小越靠前
}

// RemarkListReq 备注列表查询
type RemarkListReq struct {
	upagination.RequestPagination
	Name       string    `form:"name" json:"name"`               // 备注名称
	CategoryID uuid.UUID `form:"category_id" json:"category_id"` // 备注类型
	StoreID    uuid.UUID `form:"store_id" json:"store_id"`       // 门店 ID
	Enabled    *bool     `form:"enabled" json:"enabled"`         // 启用状态
}

type RemarkListResp struct {
	Remarks []*domain.Remark `json:"remarks"` // 备注列表
	Total   int              `json:"total"`   // 备注总数
}

type RemarkSimpleUpdateReq struct {
	SimpleUpdateType domain.RemarkSimpleUpdateType `json:"simple_update_type" binding:"required,oneof=status"` // 简单更新类型
	Enabled          bool                          `json:"enabled"`                                            // 启用状态
}
