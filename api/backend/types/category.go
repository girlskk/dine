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

// UpdateCategoryReq 更新商品分类请求
type UpdateCategoryReq struct {
	Name           string     `json:"name" binding:"omitempty,max=255"` // 分类名称
	InheritTaxRate bool       `json:"inherit_tax_rate"`                 // 是否继承父分类的税率ID（仅子分类有效）
	TaxRateID      *uuid.UUID `json:"tax_rate_id"`                      // 税率ID
	InheritStall   bool       `json:"inherit_stall"`                    // 是否继承父分类的出品部门ID（仅子分类有效）
	StallID        *uuid.UUID `json:"stall_id"`                         // 出品部门ID
}
