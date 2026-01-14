package types

import (
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// RouterMenuCreateReq 创建菜单请求
type RouterMenuCreateReq struct {
	ParentID  uuid.UUID `json:"parent_id"`                             // 父级菜单ID
	Name      string    `json:"name" binding:"required,max=100"`       // 菜单名称
	Path      string    `json:"path" binding:"omitempty,max=255"`      // 路由路径
	Component string    `json:"component" binding:"omitempty,max=255"` // 组件路径
	Icon      string    `json:"icon" binding:"omitempty,max=500"`      // 菜单图标
	Sort      int       `json:"sort" binding:"omitempty"`              // 排序
	Enabled   bool      `json:"enabled"`                               // 是否启用
}

// RouterMenuUpdateReq 更新菜单请求
type RouterMenuUpdateReq struct {
	ParentID  uuid.UUID `json:"parent_id"`                             // 父级菜单ID
	Name      string    `json:"name" binding:"required,max=100"`       // 菜单名称
	Path      string    `json:"path" binding:"omitempty,max=255"`      // 路由路径
	Component string    `json:"component" binding:"omitempty,max=255"` // 组件路径
	Icon      string    `json:"icon" binding:"omitempty,max=500"`      // 菜单图标
	Sort      int       `json:"sort" binding:"omitempty"`              // 排序
	Enabled   bool      `json:"enabled"`                               // 是否启用
}

// RouterMenuListReq 菜单列表请求
type RouterMenuListReq struct {
}

// RouterMenuListResp 菜单列表响应
type RouterMenuListResp struct {
	Menus []*domain.RouterMenu `json:"menus"`
	Total int                  `json:"total"`
}
