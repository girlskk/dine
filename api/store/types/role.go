package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// RoleCreateReq 创建角色请求
type RoleCreateReq struct {
	Name          string                `json:"name" binding:"required,max=50"`                                 // 角色名称
	Enabled       bool                  `json:"enabled"`                                                        // 是否启用
	LoginChannels []domain.LoginChannel `json:"login_channels" binding:"omitempty,dive,oneof=pos mobile store"` // 角色登录渠道pos/mobile/store
}

// RoleUpdateReq 更新角色请求
type RoleUpdateReq struct {
	Name          string                `json:"name" binding:"required,max=50"`                                 // 角色名称
	Enabled       bool                  `json:"enabled"`                                                        // 是否启用
	LoginChannels []domain.LoginChannel `json:"login_channels" binding:"omitempty,dive,oneof=pos mobile store"` // 角色登录渠道pos/mobile/store
}

// RoleListReq 角色列表查询请求
type RoleListReq struct {
	upagination.RequestPagination
	Name    string `form:"name"`    // 角色名称
	Enabled *bool  `form:"enabled"` // 启用状态
}

// RoleListResp 角色列表响应
type RoleListResp struct {
	Roles []*domain.Role `json:"roles"`
	Total int            `json:"total"`
}
type SetMenusReq struct {
	Paths []string `json:"paths" binding:"required"` // 菜单路径列表
}

type RoleMenusResp struct {
	Paths []string `json:"paths"` // 菜单路径列表
}
