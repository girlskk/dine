package table

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.TableInteractor = (*TableInteractor)(nil)

type TableInteractor struct {
	ds domain.DataStore
}

func NewTableInteractor(dataStore domain.DataStore) *TableInteractor {
	return &TableInteractor{
		ds: dataStore,
	}
}

func (i *TableInteractor) Create(ctx context.Context, table *domain.Table) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	table.StoreID = user.Store.ID

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 检查名称是否重复
		exists, err := i.ds.TableRepo().Exists(ctx, domain.TableExistsParams{
			StoreID: table.StoreID,
			Name:    table.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrTableNameExists)
		}
		// 检查区域是否存在
		_, err = i.ds.TableAreaRepo().FindByID(ctx, table.AreaID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}

		err = i.ds.TableRepo().Create(ctx, table)
		if err != nil {
			return err
		}

		return i.ds.TableAreaRepo().IncreaseTableCount(ctx, table.AreaID)
	})
}

func (i *TableInteractor) Update(ctx context.Context, table *domain.Table) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	user := domain.FromBackendUserContext(ctx)
	table.StoreID = user.Store.ID

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		existingTable, err := ds.TableRepo().FindByID(ctx, table.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}
		if existingTable.StoreID != table.StoreID {
			return domain.ParamsError(domain.ErrTableNotExists)
		}
		if existingTable.Name != table.Name {
			// 检查名称是否重复
			exists, err := ds.TableRepo().Exists(ctx, domain.TableExistsParams{
				StoreID: table.StoreID,
				Name:    table.Name,
			})
			if err != nil {
				return err
			}
			if exists {
				return domain.ParamsError(domain.ErrTableNameExists)
			}
		}

		// 检查区域是否存在
		_, err = i.ds.TableAreaRepo().FindByID(ctx, table.AreaID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}

		// 更新台桌区域台桌数
		if existingTable.AreaID != table.AreaID {
			if err = ds.TableAreaRepo().IncreaseTableCount(ctx, table.AreaID); err != nil {
				return err
			}
			if err = ds.TableAreaRepo().DecreaseTableCount(ctx, existingTable.AreaID); err != nil {
				return err
			}
		}

		// 更新属性
		existingTable.Name = table.Name
		existingTable.AreaID = table.AreaID
		existingTable.SeatCount = table.SeatCount
		return ds.TableRepo().Update(ctx, existingTable)
	})
}

func (i *TableInteractor) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	user := domain.FromBackendUserContext(ctx)

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		existingTable, err := ds.TableRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}
		if existingTable.StoreID != user.Store.ID {
			return domain.ParamsError(domain.ErrTableNotExists)
		}
		if existingTable.Status == domain.TableStatusOccupied {
			return domain.ParamsError(domain.ErrTableIsOccupied)
		}
		err = ds.TableRepo().Delete(ctx, id)
		if err != nil {
			return err
		}
		return ds.TableAreaRepo().DecreaseTableCount(ctx, existingTable.AreaID)
	})
}

func (i *TableInteractor) PagedListBySearch(ctx context.Context, page *upagination.Pagination,
	params domain.TableSearchParams,
) (res *domain.TableSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return i.ds.TableRepo().PagedListBySearch(ctx, page, params)
}

func (i *TableInteractor) Get(ctx context.Context, id int) (t *domain.Table, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableInteractor.Get")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	t, err = i.ds.TableRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ParamsError(domain.ErrTableNotExists)
		}
		return nil, err
	}
	return t, nil
}

func (i *TableInteractor) GetWithOrder(ctx context.Context, id int) (t *domain.Table, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableInteractor.GetWithOrder")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	t, err = i.ds.TableRepo().FindWithOrder(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ParamsError(domain.ErrTableNotExists)
		}
		return nil, err
	}
	return t, nil
}
