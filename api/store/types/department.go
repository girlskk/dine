package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// DepartmentCreateReq 创建部门请求
type DepartmentCreateReq struct {
	Name   string `json:"name" binding:"required"`
	Enable bool   `json:"enable"`
}

// DepartmentUpdateReq 更新部门请求
type DepartmentUpdateReq struct {
	Name   string `json:"name" binding:"required"`
	Enable bool   `json:"enable"`
}

// DepartmentListReq 部门列表查询请求
type DepartmentListReq struct {
	upagination.RequestPagination
	Name   string `form:"name"`
	Code   string `form:"code"`
	Enable *bool  `form:"enable"`
}

// DepartmentListResp 部门列表响应
type DepartmentListResp struct {
	Departments []*domain.Department `json:"departments"`
	Total       int                  `json:"total"`
}
