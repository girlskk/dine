package product

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

var _ domain.ProductInteractor = (*ProductInteractor)(nil)

type ProductInteractor struct {
	DS domain.DataStore
}

func NewProductInteractor(ds domain.DataStore) *ProductInteractor {
	return &ProductInteractor{
		DS: ds,
	}
}
