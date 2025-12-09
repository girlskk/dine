package product

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductRecipeInteractor = (*ProductRecipeInteractor)(nil)

type ProductRecipeInteractor struct {
	ds domain.DataStore
}

func NewProductRecipeInteractor(ds domain.DataStore) *ProductRecipeInteractor {
	return &ProductRecipeInteractor{
		ds: ds,
	}
}

func (i *ProductRecipeInteractor) Create(ctx context.Context, recipe *domain.ProductRecipe) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	recipe.StoreID = user.Store.ID

	exists, err := i.ds.ProductRecipeRepo().Exists(ctx, domain.RecipeExistsParams{
		StoreID: recipe.StoreID,
		Name:    recipe.Name,
	})
	if err != nil {
		return nil
	}
	if exists {
		return domain.ParamsError(domain.ErrRecipeNameExists)
	}

	return i.ds.ProductRecipeRepo().Create(ctx, recipe)
}

func (i *ProductRecipeInteractor) Update(ctx context.Context, recipe *domain.ProductRecipe) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	recipe.StoreID = user.Store.ID

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		existing, err := ds.ProductRecipeRepo().FindByID(ctx, recipe.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrRecipeNotExists)
			}
			return err
		}

		if existing.StoreID != recipe.StoreID {
			return domain.ParamsError(domain.ErrRecipeNotExists)
		}

		if existing.Name == recipe.Name {
			return nil
		}

		exists, err := ds.ProductRecipeRepo().Exists(ctx, domain.RecipeExistsParams{
			StoreID: recipe.StoreID,
			Name:    recipe.Name,
		})
		if err != nil {
			return nil
		}
		if exists {
			return domain.ParamsError(domain.ErrRecipeNameExists)
		}

		existing.Name = recipe.Name
		return ds.ProductRecipeRepo().Update(ctx, existing)
	})
}

func (i *ProductRecipeInteractor) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	user := domain.FromBackendUserContext(ctx)

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		existing, err := ds.ProductRecipeRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrRecipeNotExists)
			}
			return err
		}

		if existing.StoreID != user.Store.ID {
			return domain.ParamsError(domain.ErrRecipeNotExists)
		}

		// 检查是否被商品关联
		used, err := ds.ProductRecipeRepo().IsUsedByProduct(ctx, id)
		if err != nil {
			return err
		}
		if used {
			return domain.ParamsError(domain.ErrRecipeUsing)
		}

		return ds.ProductRecipeRepo().Delete(ctx, id)
	})
}

func (i *ProductRecipeInteractor) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.RecipeSearchParams,
) (res *domain.RecipeSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	params.StoreID = user.Store.ID

	return i.ds.ProductRecipeRepo().PagedListBySearch(ctx, page, params)
}
