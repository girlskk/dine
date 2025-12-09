package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/productspec"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductSpecRelRepository = (*ProductSpecRelRepository)(nil)

type ProductSpecRelRepository struct {
	Client *ent.Client
}

func NewProductSpecRelRepository(client *ent.Client) *ProductSpecRelRepository {
	return &ProductSpecRelRepository{
		Client: client,
	}
}

func (repo *ProductSpecRelRepository) ListByIDs(ctx context.Context, ids []int) (res domain.ProductSpecRels, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecRelRepository.ListByIDs")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	productSpecs, err := repo.Client.ProductSpec.Query().
		Where(productspec.IDIn(ids...)).
		All(ctx)

	if err != nil {
		return nil, err
	}
	if len(productSpecs) == 0 {
		return nil, nil
	}
	for _, u := range productSpecs {
		res = append(res, repo.convertToDomain(u))
	}
	return res, nil
}

func (repo *ProductSpecRelRepository) convertToDomain(s *ent.ProductSpec) *domain.ProductSpecRel {
	return &domain.ProductSpecRel{
		ID:         s.ID,
		SpecID:     s.SpecID,
		ProductID:  s.ProductID,
		Price:      s.Price,
		SaleStatus: domain.ProductSaleStatus(s.SaleStatus),
	}
}

func (repo *ProductSpecRelRepository) Exists(ctx context.Context, params domain.ProductSpecRelExistsParams) (exists bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecRelRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := repo.Client.ProductSpec.Query()
	if params.ProductID > 0 {
		query.Where(productspec.ProductID(params.ProductID))
	}
	if params.SpecID > 0 {
		query.Where(productspec.SpecID(params.SpecID))
	}
	return query.Exist(ctx)
}

func (repo *ProductSpecRelRepository) UpdateSaleStatusByIDs(ctx context.Context, ids []int, status domain.ProductSaleStatus) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecRelRepository.UpdateSaleStatusByIDs")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	_, err = repo.Client.ProductSpec.Update().
		Where(productspec.IDIn(ids...)).
		SetSaleStatus(int(status)).
		Save(ctx)
	return err
}

func (repo *ProductSpecRelRepository) FindByID(ctx context.Context, id int) (res *domain.ProductSpecRel, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecRelRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	productSpec, err := repo.Client.ProductSpec.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrSpecNotExists)
		}
		return nil, err
	}
	return repo.convertToDomain(productSpec), nil
}
