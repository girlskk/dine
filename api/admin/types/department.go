package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// DepartmentCreateReq 创建部门请求
type DepartmentCreateReq struct {
	Name    string `json:"name" binding:"required,max=50"` // 部门名称
	Enabled bool   `json:"enabled"`                        // 是否启用
}

// DepartmentUpdateReq 更新部门请求
type DepartmentUpdateReq struct {
	Name    string `json:"name" binding:"required,max=50"` // 部门名称
	Enabled bool   `json:"enabled"`                        // 是否启用
}

// DepartmentListReq 部门列表查询请求
type DepartmentListReq struct {
	upagination.RequestPagination
	Name    string `form:"name"`    // 部门名称
	Code    string `form:"code"`    // 部门编码
	Enabled *bool  `form:"enabled"` // 启用状态
}

// DepartmentListResp 部门列表响应
type DepartmentListResp struct {
	Departments []*domain.Department `json:"departments"` // 部门列表
	Total       int                  `json:"total"`       // 部门总数
}
