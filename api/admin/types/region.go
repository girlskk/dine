package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

type CountryListResp struct {
	Countries []*domain.Country `json:"countries"`
}
type ProvinceListResp struct {
	Provinces []*domain.Province `json:"provinces"`
}
