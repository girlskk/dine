package product

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductSpecInteractor = (*ProductSpecInteractor)(nil)

type ProductSpecInteractor struct {
	ds domain.DataStore
}

func NewProductSpecInteractor(ds domain.DataStore) *ProductSpecInteractor {
	return &ProductSpecInteractor{
		ds: ds,
	}
}

func (i *ProductSpecInteractor) Create(ctx context.Context, spec *domain.ProductSpec) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	spec.StoreID = user.Store.ID

	exists, err := i.ds.ProductSpecRepo().Exists(ctx, domain.SpecExistsParams{
		StoreID: spec.StoreID,
		Name:    spec.Name,
	})
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(err)
		}
		return err
	}
	if exists {
		return domain.ParamsError(domain.ErrSpecNameExists)
	}

	return i.ds.ProductSpecRepo().Create(ctx, spec)
}

func (i *ProductSpecInteractor) Update(ctx context.Context, spec *domain.ProductSpec) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	currentStoreID := user.Store.ID

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		existing, err := ds.ProductSpecRepo().FindByID(ctx, spec.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrSpecNotExists)
			}
			return err
		}

		if existing.StoreID != currentStoreID {
			return domain.ParamsError(domain.ErrSpecNotExists)
		}

		if existing.Name == spec.Name {
			return nil
		}

		exists, err := ds.ProductSpecRepo().Exists(ctx, domain.SpecExistsParams{
			StoreID: currentStoreID,
			Name:    spec.Name,
		})
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrSpecNameExists)
		}

		existing.Name = spec.Name
		return ds.ProductSpecRepo().Update(ctx, existing)
	})
}

func (i *ProductSpecInteractor) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		existing, err := ds.ProductSpecRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrSpecNotExists)
			}
			return err
		}

		if existing.StoreID != user.Store.ID {
			return domain.ParamsError(domain.ErrSpecNotExists)
		}

		// 检查是否被商品关联
		used, err := ds.ProductSpecRelRepo().Exists(ctx, domain.ProductSpecRelExistsParams{SpecID: id})
		if err != nil {
			return err
		}
		if used {
			return domain.ParamsError(domain.ErrSpecUsing)
		}
		return ds.ProductSpecRepo().Delete(ctx, id)
	})
}

func (i *ProductSpecInteractor) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.SpecSearchParams,
) (res *domain.SpecSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	params.StoreID = user.Store.ID

	return i.ds.ProductSpecRepo().PagedListBySearch(ctx, page, params)
}
