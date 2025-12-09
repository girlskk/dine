package repository

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/spec"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductSpecRepository = (*ProductSpecRepository)(nil)

type ProductSpecRepository struct {
	Client *ent.Client
}

func NewProductSpecRepository(client *ent.Client) *ProductSpecRepository {
	return &ProductSpecRepository{
		Client: client,
	}
}

func (repo *ProductSpecRepository) FindByID(ctx context.Context, id int) (res *domain.ProductSpec, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	entSpec, err := repo.Client.Spec.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrSpecNotExists)
		}
		return nil, err
	}
	return repo.convertToDomain(entSpec), nil
}

func (repo *ProductSpecRepository) Exists(ctx context.Context, params domain.SpecExistsParams) (exists bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Spec.Query()
	if params.StoreID != 0 {
		query.Where(spec.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(spec.Name(params.Name))
	}
	return query.Exist(ctx)
}

func (repo *ProductSpecRepository) Create(ctx context.Context, spec *domain.ProductSpec) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	created, err := repo.Client.Spec.Create().
		SetName(spec.Name).
		SetStoreID(spec.StoreID).
		Save(ctx)

	if err != nil {
		return err
	}

	spec.ID = created.ID
	spec.CreatedAt = created.CreatedAt
	return nil
}

func (repo *ProductSpecRepository) Update(ctx context.Context, spec *domain.ProductSpec) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	updated, err := repo.Client.Spec.UpdateOneID(spec.ID).
		SetName(spec.Name).
		Save(ctx)

	if err != nil {
		return err
	}
	spec.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *ProductSpecRepository) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.Spec.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductSpecRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.SpecSearchParams,
) (res *domain.SpecSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Spec.Query()

	if params.StoreID != 0 {
		query.Where(spec.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(spec.NameContains(params.Name))
	}

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	entSpecs, err := query.Order(ent.Desc(spec.FieldID)).
		Offset(page.Offset()).
		Limit(page.Size).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.ProductSpecs, 0, len(entSpecs))
	for _, s := range entSpecs {
		items = append(items, repo.convertToDomain(s))
	}

	return &domain.SpecSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

func (repo *ProductSpecRepository) convertToDomain(s *ent.Spec) *domain.ProductSpec {
	return &domain.ProductSpec{
		ID:        s.ID,
		Name:      s.Name,
		StoreID:   s.StoreID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func (repo *ProductSpecRepository) ListByIDs(ctx context.Context, ids []int) (res domain.ProductSpecs, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecRepository.ListByIDs")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	specs, err := repo.Client.Spec.Query().
		Where(spec.IDIn(ids...)).
		All(ctx)

	if err != nil {
		return nil, err
	}
	if len(specs) == 0 {
		return nil, nil
	}
	for _, u := range specs {
		res = append(res, repo.convertToDomain(u))
	}
	return res, nil
}
