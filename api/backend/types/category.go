package types

import "github.com/google/uuid"

// CategoryCreateRootReq 创建一级商品分类请求
type CategoryCreateRootReq struct {
	Name          string     `json:"name" binding:"required,max=255"` // 分类名称
	TaxRateID     *uuid.UUID `json:"tax_rate_id"`                     // 税率ID
	StallID       *uuid.UUID `json:"stall_id"`                        // 出品部门ID
	ChildrenNames []string   `json:"children_names"`                  // 子分类名称列表
}

// CategoryCreateChildReq 创建二级商品分类请求
type CategoryCreateChildReq struct {
	Name           string     `json:"name" binding:"required,max=255"` // 分类名称
	InheritTaxRate bool       `json:"inherit_tax_rate"`                // 是否继承父分类的税率ID
	TaxRateID      *uuid.UUID `json:"tax_rate_id"`                     // 税率ID
	InheritStall   bool       `json:"inherit_stall"`                   // 是否继承父分类的出品部门ID
	StallID        *uuid.UUID `json:"stall_id"`                        // 出品部门ID
}

// CategoryPagedListReq 分页查询商品分类列表请求
type CategoryPagedListReq struct {
	Page int       `json:"page"`
	Size int       `json:"size"`
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}
