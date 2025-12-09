package product

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductAttrInteractor = (*ProductAttrInteractor)(nil)

type ProductAttrInteractor struct {
	ds domain.DataStore
}

func NewProductAttrInteractor(dataStore domain.DataStore) *ProductAttrInteractor {
	return &ProductAttrInteractor{
		ds: dataStore,
	}
}

func (i *ProductAttrInteractor) Create(ctx context.Context, attr *domain.ProductAttr) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	attr.StoreID = user.Store.ID

	exists, err := i.ds.ProductAttrRepo().Exists(ctx, domain.AttrExistsParams{
		StoreID: attr.StoreID,
		Name:    attr.Name,
	})
	if err != nil {
		return err
	}
	if exists {
		return domain.ParamsError(domain.ErrAttrNameExists)
	}
	return i.ds.ProductAttrRepo().Create(ctx, attr)
}

func (i *ProductAttrInteractor) Update(ctx context.Context, attr *domain.ProductAttr) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	attr.StoreID = user.Store.ID

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 检查属性是否存在
		existingAttr, err := ds.ProductAttrRepo().FindByID(ctx, attr.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}
		if existingAttr.StoreID != attr.StoreID {
			return domain.ParamsError(domain.ErrAttrNotExists)
		}
		if existingAttr.Name == attr.Name {
			return nil
		}
		// 检查名称是否重复
		exists, err := ds.ProductAttrRepo().Exists(ctx, domain.AttrExistsParams{
			StoreID: attr.StoreID,
			Name:    attr.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrAttrNameExists)
		}

		// 更新属性
		existingAttr.Name = attr.Name
		return ds.ProductAttrRepo().Update(ctx, existingAttr)
	})
}

func (i *ProductAttrInteractor) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 检查属性是否存在
		existingAttr, err := ds.ProductAttrRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}
		if existingAttr.StoreID != user.Store.ID {
			return domain.ParamsError(domain.ErrAttrNotExists)
		}
		// 2. 检查是否被商品关联
		used, err := ds.ProductAttrRepo().IsUsedByProduct(ctx, id)
		if err != nil {
			return err
		}
		if used {
			return domain.ParamsError(domain.ErrAttrUsing)
		}
		// 3. 执行删除
		return ds.ProductAttrRepo().Delete(ctx, id)
	})
}

func (i *ProductAttrInteractor) PagedListBySearch(ctx context.Context,
	page *upagination.Pagination, params domain.AttrSearchParams,
) (res *domain.AttrSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	params.StoreID = user.Store.ID

	return i.ds.ProductAttrRepo().PagedListBySearch(ctx, page, params)
}
