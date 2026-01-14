package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// RemarkCreateReq 创建备注请求
// 名称、分类、类型为必填；排序为可选；Enabled 默认 true
// MerchantID 从登录用户上下文取；StoreID 由前端传入（按需）
type RemarkCreateReq struct {
	Name        string             `json:"name" binding:"required,max=255"`                                                                         // 备注名称
	Enabled     bool               `json:"enabled"`                                                                                                 // 是否启用
	SortOrder   int                `json:"sort_order" binding:"omitempty,gte=0"`                                                                    // 排序，越小越靠前
	RemarkScene domain.RemarkScene `json:"remark_scene" binding:"required,oneof=whole_order item cancel_reason discount gift rebill refund_reject"` // 使用场景
}

// RemarkUpdateReq 更新备注请求
type RemarkUpdateReq struct {
	Name      string `json:"name" binding:"required,max=255"`      // 备注名称
	Enabled   bool   `json:"enabled"`                              // 是否启用
	SortOrder int    `json:"sort_order" binding:"omitempty,gte=0"` // 排序，越小越靠前
}

// RemarkListReq 备注列表查询
type RemarkListReq struct {
	upagination.RequestPagination
	Name        string             `form:"name"`         // 备注名称
	RemarkScene domain.RemarkScene `form:"remark_scene"` // 使用场景
	Enabled     *bool              `form:"enabled"`      // 启用状态
}

type RemarkListResp struct {
	Remarks []*domain.Remark `json:"remarks"` // 备注列表
	Total   int              `json:"total"`   // 备注总数
}

type RemarkCountResp struct {
	Items []RemarkCountItem `json:"items"` // 备注数量列表
}

type RemarkCountItem struct {
	RemarkScene domain.RemarkScene `json:"remark_scene"` // 使用场景
	Count       int                `json:"count"`        // 备注数量
}
