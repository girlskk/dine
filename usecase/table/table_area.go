package table

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.TableAreaInteractor = (*TableAreaInteractor)(nil)

type TableAreaInteractor struct {
	ds domain.DataStore
}

func NewTableAreaInteractor(dataStore domain.DataStore) *TableAreaInteractor {
	return &TableAreaInteractor{
		ds: dataStore,
	}
}

func (i *TableAreaInteractor) Create(ctx context.Context, area *domain.TableArea) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	area.StoreID = user.Store.ID

	exists, err := i.ds.TableAreaRepo().Exists(ctx, domain.AreaExistsParams{
		StoreID: area.StoreID,
		Name:    area.Name,
	})
	if err != nil {
		return err
	}
	if exists {
		return domain.ParamsError(domain.ErrAreaNameExists)
	}
	return i.ds.TableAreaRepo().Create(ctx, area)
}

func (i *TableAreaInteractor) Update(ctx context.Context, area *domain.TableArea) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	area.StoreID = user.Store.ID

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 检查属性是否存在
		existingArea, err := ds.TableAreaRepo().FindByID(ctx, area.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}
		if existingArea.StoreID != area.StoreID {
			return domain.ParamsError(domain.ErrUnitNotExists)
		}
		if existingArea.Name == area.Name {
			return nil
		}
		// 检查名称是否重复
		exists, err := ds.TableAreaRepo().Exists(ctx, domain.AreaExistsParams{
			StoreID: area.StoreID,
			Name:    area.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrAreaNameExists)
		}

		// 更新属性
		existingArea.Name = area.Name
		return ds.TableAreaRepo().Update(ctx, existingArea)
	})
}

func (i *TableAreaInteractor) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		existingArea, err := ds.TableAreaRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}
		if existingArea.StoreID != user.Store.ID {
			return domain.ParamsError(domain.ErrAreaNotExists)
		}
		if existingArea.TableCount > 0 {
			return domain.ParamsError(domain.ErrAreaHasTables)
		}
		return ds.TableAreaRepo().Delete(ctx, id)
	})
}

func (i *TableAreaInteractor) PagedListBySearch(ctx context.Context,
	page *upagination.Pagination, params domain.AreaSearchParams,
) (res *domain.AreaSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.ds.TableAreaRepo().PagedListBySearch(ctx, page, params)
}
