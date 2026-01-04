package types

import (
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

type ProvinceListReq struct {
	CountryID uuid.UUID `form:"country_id" binding:"omitempty"`
}

type CountryListResp struct {
	Countries []*domain.Country `json:"countries"`
}
type ProvinceListResp struct {
	Provinces []*domain.Province `json:"provinces"`
}
