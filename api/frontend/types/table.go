package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// TableAreaListReq 获取台桌区域列表请求
type TableAreaListReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type TableListReq struct {
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	AreaID int                 `json:"area_id"` // 区域ID
	Status *domain.TableStatus `json:"status"`  // 台桌状态：1-空闲 2-就餐中
}
