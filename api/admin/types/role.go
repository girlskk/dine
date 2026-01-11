package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// RoleCreateReq 创建角色请求
type RoleCreateReq struct {
	Name   string `json:"name" binding:"required,max=50"`
	Enable bool   `json:"enable"`
}

// RoleUpdateReq 更新角色请求
type RoleUpdateReq struct {
	Name   string `json:"name" binding:"required,max=50"`
	Enable bool   `json:"enable"`
}

// RoleListReq 角色列表查询请求
type RoleListReq struct {
	upagination.RequestPagination
	Name   string `form:"name"`
	Enable *bool  `form:"enable"`
}

// RoleListResp 角色列表响应
type RoleListResp struct {
	Roles []*domain.Role `json:"roles"`
	Total int            `json:"total"`
}

type SetMenusReq struct {
	Paths []string `json:"paths" binding:"required"`
}

type RoleMenusResp struct {
	Paths []string `json:"paths"`
}
