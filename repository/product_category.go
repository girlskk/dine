package repository

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/category"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.CategoryRepository = (*CategoryRepository)(nil)

type CategoryRepository struct {
	Client *ent.Client
}

func NewCategoryRepository(client *ent.Client) *CategoryRepository {
	return &CategoryRepository{
		Client: client,
	}
}

func (repo *CategoryRepository) FindByID(ctx context.Context, id int) (res *domain.Category, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	entCat, err := repo.Client.Category.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrCategoryNotExists)
		}
		return nil, err
	}
	return repo.convertToDomain(entCat), nil
}

func (repo *CategoryRepository) Exists(ctx context.Context, params domain.CategoryExistsParams) (exists bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Category.Query()
	if params.StoreID != 0 {
		query.Where(category.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(category.Name(params.Name))
	}
	return query.Exist(ctx)
}

func (repo *CategoryRepository) Create(ctx context.Context, cat *domain.Category) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	created, err := repo.Client.Category.Create().
		SetName(cat.Name).
		SetStoreID(cat.StoreID).
		Save(ctx)

	if err != nil {
		return err
	}

	cat.ID = created.ID
	cat.CreatedAt = created.CreatedAt
	return nil
}

func (repo *CategoryRepository) Update(ctx context.Context, cat *domain.Category) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	updated, err := repo.Client.Category.UpdateOneID(cat.ID).
		SetName(cat.Name).
		Save(ctx)

	if err != nil {
		return err
	}
	cat.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *CategoryRepository) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.Category.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *CategoryRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.CategorySearchParams,
) (res *domain.CategorySearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Category.Query()

	if params.StoreID != 0 {
		query.Where(category.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(category.NameContains(params.Name))
	}

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	entCats, err := query.Order(ent.Desc(category.FieldID)).
		Offset(page.Offset()).
		Limit(page.Size).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.Categories, 0, len(entCats))
	for _, c := range entCats {
		items = append(items, repo.convertToDomain(c))
	}

	return &domain.CategorySearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

func (repo *CategoryRepository) convertToDomain(c *ent.Category) *domain.Category {
	return &domain.Category{
		ID:        c.ID,
		Name:      c.Name,
		StoreID:   c.StoreID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
