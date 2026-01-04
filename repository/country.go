package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/country"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.CountryRepository = (*CountryRepository)(nil)

type CountryRepository struct {
	Client *ent.Client
}

func NewCountryRepository(client *ent.Client) *CountryRepository {
	return &CountryRepository{Client: client}
}

func (repo *CountryRepository) GetAll(ctx context.Context) (domainCountries []*domain.Country, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CountryRepository.GetAll")
	defer func() { util.SpanErrFinish(span, err) }()

	ecs, err := repo.Client.Country.Query().Order(country.BySort()).All(ctx)
	if err != nil {
		return nil, err
	}

	domainCountries = lo.Map(ecs, func(ec *ent.Country, _ int) *domain.Country {
		return convertCountryToDomain(ec)
	})
	return
}

func (repo *CountryRepository) FindByID(ctx context.Context, id uuid.UUID) (domainCountry *domain.Country, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CountryRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	ec, err := repo.Client.Country.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(err)
		}
		return nil, err
	}
	domainCountry = convertCountryToDomain(ec)
	return
}

func convertCountryToDomain(ec *ent.Country) *domain.Country {
	if ec == nil {
		return nil
	}
	return &domain.Country{
		ID:   ec.ID,
		Name: ec.Name,
		Sort: ec.Sort,
	}
}
