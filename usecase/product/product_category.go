package product

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.CategoryInteractor = (*CategoryInteractor)(nil)

type CategoryInteractor struct {
	ds domain.DataStore
}

func NewCategoryInteractor(ds domain.DataStore) *CategoryInteractor {
	return &CategoryInteractor{
		ds: ds,
	}
}

func (i *CategoryInteractor) Create(ctx context.Context, cat *domain.Category) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	cat.StoreID = user.Store.ID

	exists, err := i.ds.ProductCategoryRepo().Exists(ctx, domain.CategoryExistsParams{
		StoreID: cat.StoreID,
		Name:    cat.Name,
	})
	if err != nil {
		return err
	}
	if exists {
		return domain.ParamsError(domain.ErrCategoryNameExists)
	}

	return i.ds.ProductCategoryRepo().Create(ctx, cat)
}

func (i *CategoryInteractor) Update(ctx context.Context, cat *domain.Category) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	cat.StoreID = user.Store.ID

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		existing, err := ds.ProductCategoryRepo().FindByID(ctx, cat.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}
		if existing.StoreID != cat.StoreID {
			return domain.ParamsError(domain.ErrCategoryNotExists)
		}
		if existing.Name == cat.Name {
			return nil
		}

		exists, err := ds.ProductCategoryRepo().Exists(ctx, domain.CategoryExistsParams{
			StoreID: cat.StoreID,
			Name:    cat.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrCategoryNameExists)
		}

		existing.Name = cat.Name
		return ds.ProductCategoryRepo().Update(ctx, existing)
	})
}

func (i *CategoryInteractor) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		existing, err := ds.ProductCategoryRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}

		if existing.StoreID != user.Store.ID {
			return domain.ParamsError(domain.ErrCategoryNotExists)
		}

		// 检查是否被商品关联
		used, err := ds.ProductRepo().Exists(ctx, domain.ProductExistsParams{CategoryID: id})
		if err != nil {
			return err
		}
		if used {
			return domain.ParamsError(domain.ErrCategoryUsing)
		}
		return ds.ProductCategoryRepo().Delete(ctx, id)
	})
}

func (i *CategoryInteractor) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.CategorySearchParams,
) (res *domain.CategorySearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return i.ds.ProductCategoryRepo().PagedListBySearch(ctx, page, params)
}
