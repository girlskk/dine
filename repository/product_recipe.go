package repository

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/recipe"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductRecipeRepository = (*ProductRecipeRepository)(nil)

type ProductRecipeRepository struct {
	Client *ent.Client
}

func NewProductRecipeRepository(client *ent.Client) *ProductRecipeRepository {
	return &ProductRecipeRepository{
		Client: client,
	}
}

func (repo *ProductRecipeRepository) FindByID(ctx context.Context, id int) (res *domain.ProductRecipe, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	entRecipe, err := repo.Client.Recipe.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrRecipeNotExists)
		}
		return nil, err
	}
	return repo.convertToDomain(entRecipe), nil
}

func (repo *ProductRecipeRepository) Exists(ctx context.Context, params domain.RecipeExistsParams) (exists bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Recipe.Query()
	if params.StoreID != 0 {
		query.Where(recipe.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(recipe.Name(params.Name))
	}
	return query.Exist(ctx)
}

func (repo *ProductRecipeRepository) Create(ctx context.Context, recipe *domain.ProductRecipe) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	created, err := repo.Client.Recipe.Create().
		SetName(recipe.Name).
		SetStoreID(recipe.StoreID).
		Save(ctx)

	if err != nil {
		return err
	}

	recipe.ID = created.ID
	recipe.CreatedAt = created.CreatedAt
	return nil
}

func (repo *ProductRecipeRepository) Update(ctx context.Context, recipe *domain.ProductRecipe) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	updated, err := repo.Client.Recipe.UpdateOneID(recipe.ID).
		SetName(recipe.Name).
		Save(ctx)
	if err != nil {
		return err
	}
	recipe.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *ProductRecipeRepository) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.Recipe.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRecipeRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.RecipeSearchParams,
) (res *domain.RecipeSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Recipe.Query()

	if params.StoreID != 0 {
		query.Where(recipe.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(recipe.NameContains(params.Name))
	}

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	entRecipes, err := query.Order(ent.Desc(recipe.FieldID)).
		Offset(page.Offset()).
		Limit(page.Size).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.ProductRecipes, 0, len(entRecipes))
	for _, r := range entRecipes {
		items = append(items, repo.convertToDomain(r))
	}

	return &domain.RecipeSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

func (repo *ProductRecipeRepository) ListByIDs(ctx context.Context, ids []int) (res domain.ProductRecipes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeRepository.ListByIDs")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	recipes, err := repo.Client.Recipe.Query().
		Where(recipe.IDIn(ids...)).
		All(ctx)

	if err != nil {
		return nil, err
	}
	if len(recipes) == 0 {
		return nil, nil
	}
	for _, u := range recipes {
		res = append(res, repo.convertToDomain(u))
	}
	return res, nil
}

func (repo *ProductRecipeRepository) IsUsedByProduct(ctx context.Context, id int) (used bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeRepository.IsUsedByProduct")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	used, err = repo.Client.Recipe.Query().
		Where(recipe.ID(id)).
		QueryProducts().
		Exist(ctx)
	if err != nil {
		return false, err
	}
	return used, nil
}

func (repo *ProductRecipeRepository) convertToDomain(r *ent.Recipe) *domain.ProductRecipe {
	return &domain.ProductRecipe{
		ID:        r.ID,
		Name:      r.Name,
		StoreID:   r.StoreID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
