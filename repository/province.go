package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/province"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProvinceRepository = (*ProvinceRepository)(nil)

type ProvinceRepository struct {
	Client *ent.Client
}

func NewProvinceRepository(client *ent.Client) *ProvinceRepository {
	return &ProvinceRepository{Client: client}
}

func (repo *ProvinceRepository) GetAll(ctx context.Context, countryID uuid.UUID) (provinces []*domain.Province, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProvinceRepository.GetAll")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.Province.Query().Order(province.BySort())
	if countryID != uuid.Nil {
		query = query.Where(province.CountryID(countryID))
	}

	eps, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	provinces = make([]*domain.Province, 0, len(eps))
	for _, p := range eps {
		provinces = append(provinces, convertProvinceToDomain(p))
	}
	return provinces, nil
}

func (repo *ProvinceRepository) FindByID(ctx context.Context, id uuid.UUID) (pv *domain.Province, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProvinceRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	ep, err := repo.Client.Province.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(err)
		}
		return nil, err
	}
	return convertProvinceToDomain(ep), nil
}

func (repo *ProvinceRepository) Create(ctx context.Context, province *domain.Province) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProvinceRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()

	builder := repo.Client.Province.Create().
		SetID(province.ID).
		SetCountryID(province.CountryID).
		SetName(province.Name).
		SetSort(province.Sort)

	created, err := builder.Save(ctx)
	if err != nil {
		return err
	}

	province.ID = created.ID
	province.CountryID = created.CountryID
	return nil
}

func (repo *ProvinceRepository) Update(ctx context.Context, province *domain.Province) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProvinceRepository.Update")
	defer func() { util.SpanErrFinish(span, err) }()

	updated, err := repo.Client.Province.UpdateOneID(province.ID).
		SetName(province.Name).
		SetSort(province.Sort).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.NotFoundError(err)
		}
		return err
	}

	province.CountryID = updated.CountryID
	return nil
}

func (repo *ProvinceRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProvinceRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.Province.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.NotFoundError(err)
		}
		return err
	}
	return nil
}

func (repo *ProvinceRepository) GetByFilter(ctx context.Context, filter *domain.ProvinceListFilter) (domainProvinces []*domain.Province, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProvinceRepository.GetByFilter")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.Province.Query()
	if filter != nil {
		if filter.CountryID != uuid.Nil {
			query = query.Where(province.CountryID(filter.CountryID))
		}
		if filter.Name != "" {
			query = query.Where(province.Name(filter.Name))
		}
	}

	eps, err := query.Order(province.BySort()).All(ctx)
	if err != nil {
		return nil, err
	}

	domainProvinces = lo.Map(eps, func(ep *ent.Province, _ int) *domain.Province {
		return convertProvinceToDomain(ep)
	})
	return
}

func convertProvinceToDomain(ep *ent.Province) *domain.Province {
	if ep == nil {
		return nil
	}
	return &domain.Province{
		ID:        ep.ID,
		CountryID: ep.CountryID,
		Name:      ep.Name,
		Sort:      ep.Sort,
	}
}
