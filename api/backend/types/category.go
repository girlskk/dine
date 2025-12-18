package types

import "github.com/google/uuid"

// CreateRootCategoryReq 创建一级商品分类请求
type CreateRootCategoryReq struct {
	Name          string     `json:"name" binding:"required,max=255"` // 分类名称
	TaxRateID     *uuid.UUID `json:"tax_rate_id"`                     // 税率ID
	StallID       *uuid.UUID `json:"stall_id"`                        // 出品部门ID
	ChildrenNames []string   `json:"children_names"`                  // 子分类名称列表
}

// CreateChildCategoryReq 创建二级商品分类请求
type CreateChildCategoryReq struct {
	Name           string     `json:"name" binding:"required,max=255"` // 分类名称
	InheritTaxRate bool       `json:"inherit_tax_rate"`                // 是否继承父分类的税率ID
	TaxRateID      *uuid.UUID `json:"tax_rate_id"`                     // 税率ID
	InheritStall   bool       `json:"inherit_stall"`                   // 是否继承父分类的出品部门ID
	StallID        *uuid.UUID `json:"stall_id"`                        // 出品部门ID
}

// // UpdateCategoryReq 更新商品分类请求
// type UpdateCategoryReq struct {
// 	Name         string    `json:"name" binding:"required,max=255"`  // 分类名称
// 	TaxRateID    uuid.UUID `json:"tax_rate_id" binding:"required"`   // 税率ID
// 	DepartmentID uuid.UUID `json:"department_id" binding:"required"` // 出品部门ID
// }

// // ListCategoryReq 查询商品分类列表请求
// type ListCategoryReq struct {
// 	ParentID *uuid.UUID `form:"parent_id"` // 父分类ID，为空表示查询一级分类
// }

// // CategoryResp 商品分类响应
// type CategoryResp struct {
// 	ID           uuid.UUID      `json:"id"`                 // 分类ID
// 	Name         string         `json:"name"`               // 分类名称
// 	StoreID      uuid.UUID      `json:"store_id"`           // 门店ID
// 	ParentID     *uuid.UUID     `json:"parent_id"`          // 父分类ID，nil表示一级分类
// 	TaxRateID    uuid.UUID      `json:"tax_rate_id"`        // 税率ID
// 	DepartmentID uuid.UUID      `json:"department_id"`      // 出品部门ID
// 	ProductCount int            `json:"product_count"`      // 关联的商品数量
// 	CreatedAt    string         `json:"created_at"`         // 创建时间
// 	UpdatedAt    string         `json:"updated_at"`         // 更新时间
// 	Parent       *CategoryResp  `json:"parent,omitempty"`   // 父分类
// 	Children     []CategoryResp `json:"children,omitempty"` // 子分类列表
// }
