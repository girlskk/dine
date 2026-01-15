package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// StallListReq 出品部门列表查询
type StallListReq struct {
	Enabled   *bool                 `form:"enabled"`    // 启用状态
	PrintType domain.StallPrintType `form:"print_type"` // 打印类型
}

type StallListResp struct {
	Stalls []*domain.Stall `json:"stalls"` // 出品部门列表
	Total  int             `json:"total"`  // 出品部门总数
}
