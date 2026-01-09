package region

import (
	"context"
	"fmt"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.CountryInteractor = (*CountryInteractor)(nil)

type CountryInteractor struct {
	DS domain.DataStore
}

func NewCountryInteractor(ds domain.DataStore) *CountryInteractor {
	return &CountryInteractor{DS: ds}
}

func (interactor *CountryInteractor) GetAllCountries(ctx context.Context) (countryList []*domain.Country, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "CountryInteractor.GetAllCountries")
	defer func() { util.SpanErrFinish(span, err) }()

	countryList, err = interactor.DS.CountryRepo().GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all countries: %w", err)
	}
	return
}
