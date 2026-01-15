package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

type StoreListReq struct {
	Status           domain.StoreStatus   `form:"status" binding:"omitempty"`             // 营业/停业
	BusinessModel    domain.BusinessModel `form:"business_model" binding:"omitempty"`     // 直营/加盟
	BusinessTypeCode domain.BusinessType  `form:"business_type_code" binding:"omitempty"` // 业务类型
}

type StoreListResp struct {
	Stores []*domain.Store `json:"stores"` // 门店列表
	Total  int             `json:"total"`  // 总数
}
