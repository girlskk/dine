package types

import (
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// DepartmentCreateReq 创建部门请求
type DepartmentCreateReq struct {
	MerchantID     uuid.UUID             `json:"merchant_id" binding:"required"`
	StoreID        uuid.UUID             `json:"store_id" binding:"required"`
	Name           string                `json:"name" binding:"required,max=50"`
	Code           string                `json:"code" binding:"required,max=50"`
	DepartmentType domain.DepartmentType `json:"department_type" binding:"required,oneof=admin backend store"`
	Enable         bool                  `json:"enable"`
}

// DepartmentUpdateReq 更新部门请求
type DepartmentUpdateReq struct {
	Name           string                `json:"name" binding:"required,max=50"`
	Code           string                `json:"code" binding:"required,max=50"`
	DepartmentType domain.DepartmentType `json:"department_type" binding:"required,oneof=admin backend store"`
	Enable         bool                  `json:"enable"`
}

// DepartmentListReq 部门列表查询请求
type DepartmentListReq struct {
	upagination.RequestPagination
	MerchantID     uuid.UUID             `form:"merchant_id" binding:"required"`
	StoreID        uuid.UUID             `form:"store_id" binding:"required"`
	Name           string                `form:"name"`
	Code           string                `form:"code"`
	DepartmentType domain.DepartmentType `form:"department_type" binding:"omitempty,oneof=admin backend store"`
	Enable         *bool                 `form:"enable"`
}

// DepartmentListResp 部门列表响应
type DepartmentListResp struct {
	Departments []*domain.Department `json:"departments"`
	Total       int                  `json:"total"`
}
