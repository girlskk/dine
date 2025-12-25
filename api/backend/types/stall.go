package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// StallCreateReq 创建出品部门请求
type StallCreateReq struct {
	Name      string                `json:"name" binding:"required,max=20"`                    // 出品部门名称
	PrintType domain.StallPrintType `json:"print_type" binding:"required,oneof=receipt label"` // 打印类型
	Enabled   bool                  `json:"enabled"`                                           // 是否启用
	SortOrder int                   `json:"sort_order" binding:"omitempty,gte=0"`              // 排序
}

// StallUpdateReq 更新出品部门请求
type StallUpdateReq struct {
	Name      string                `json:"name" binding:"required,max=20"`
	PrintType domain.StallPrintType `json:"print_type" binding:"required,oneof=receipt label"`
	Enabled   bool                  `json:"enabled"`
	SortOrder int                   `json:"sort_order" binding:"omitempty,gte=0"`
}

// StallListReq 出品部门列表查询
type StallListReq struct {
	upagination.RequestPagination
	Name      string                `form:"name" json:"name"`             // 名称模糊查询
	Enabled   *bool                 `form:"enabled" json:"enabled"`       // 启用状态
	PrintType domain.StallPrintType `form:"print_type" json:"print_type"` // 打印类型
}

type StallListResp struct {
	Stalls []*domain.Stall `json:"stalls"` // 出品部门列表
	Total  int             `json:"total"`  // 出品部门总数
}

// StallSimpleUpdateReq 简单更新出品部门字段
type StallSimpleUpdateReq struct {
	SimpleUpdateType domain.StallSimpleUpdateType `json:"simple_update_type" binding:"required,oneof=enabled"` // 更新的字段名称
	Enabled          bool                         `json:"enabled"`                                             // 启用状态
}
