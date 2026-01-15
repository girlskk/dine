package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// AdditionalFeeListReq 列表查询
type AdditionalFeeListReq struct {
	Enabled *bool `form:"enabled"` // 是否启用
}

// AdditionalFeeListResp 附加费列表响应
type AdditionalFeeListResp struct {
	AdditionalFees []*domain.AdditionalFee `json:"additional_fees"` // 附加费列表
	Total          int                     `json:"total"`           // 总数
}
