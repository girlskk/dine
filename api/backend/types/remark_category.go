package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// RemarkCategoryCreateReq 创建备注分类
type RemarkCategoryCreateReq struct {
	Name        string             `json:"name" binding:"required,max=50"`                                                                          // 分类名称
	RemarkScene domain.RemarkScene `json:"remark_scene" binding:"required,oneof=whole_order item cancel_reason discount gift rebill refund_reject"` // 使用场景
	Description string             `json:"description" binding:"omitempty,max=255"`                                                                 // 分类描述
	SortOrder   int                `json:"sort_order" binding:"omitempty,gte=0"`                                                                    // 排序，越小越靠前
}

// RemarkCategoryUpdateReq 更新备注分类请求
type RemarkCategoryUpdateReq struct {
	Name        string             `json:"name" binding:"required,max=50"`
	RemarkScene domain.RemarkScene `json:"remark_scene" binding:"required,oneof=whole_order item cancel_reason discount gift rebill refund_reject"`
	Description string             `json:"description" binding:"omitempty,max=255"`
	SortOrder   int                `json:"sort_order" binding:"omitempty,gte=0"`
}

// RemarkCategoryListReq 备注分类列表查询
type RemarkCategoryListReq struct {
}

// RemarkGroupListResp 备注类型列表
type RemarkGroupListResp struct {
	RemarkGroups []domain.RemarkGroup `json:"remark_groups"`
}
