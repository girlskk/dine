package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// RemarkListReq 备注列表查询
type RemarkListReq struct {
	RemarkScene domain.RemarkScene `form:"remark_scene"` // 使用场景
	Enabled     *bool              `form:"enabled"`      // 启用状态
}

type RemarkListResp struct {
	Remarks []*domain.Remark `json:"remarks"` // 备注列表
	Total   int              `json:"total"`   // 备注总数
}
