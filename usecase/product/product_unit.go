package product

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductUnitInteractor = (*ProductUnitInteractor)(nil)

type ProductUnitInteractor struct {
	ds domain.DataStore
}

func NewProductUnitInteractor(dataStore domain.DataStore) *ProductUnitInteractor {
	return &ProductUnitInteractor{
		ds: dataStore,
	}
}

func (i *ProductUnitInteractor) Create(ctx context.Context, unit *domain.ProductUnit) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	unit.StoreID = user.Store.ID

	// 检查名称重复
	exists, err := i.ds.ProductUnitRepo().Exists(ctx, domain.UnitExistsParams{
		StoreID: unit.StoreID,
		Name:    unit.Name,
	})
	if err != nil {
		return err
	}
	if exists {
		return domain.ParamsError(domain.ErrUnitNameExists)
	}
	// 调用仓储
	return i.ds.ProductUnitRepo().Create(ctx, unit)
}

func (i *ProductUnitInteractor) Update(ctx context.Context, unit *domain.ProductUnit) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	user := domain.FromBackendUserContext(ctx)
	unit.StoreID = user.Store.ID

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		existingUnit, err := ds.ProductUnitRepo().FindByID(ctx, unit.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}
		if existingUnit.StoreID != unit.StoreID {
			return domain.ParamsError(domain.ErrUnitNotExists)
		}
		if existingUnit.Name == unit.Name {
			return nil
		}
		// 检查名称是否重复
		exists, err := ds.ProductUnitRepo().Exists(ctx, domain.UnitExistsParams{
			StoreID: unit.StoreID,
			Name:    unit.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrUnitNameExists)
		}

		// 更新属性
		existingUnit.Name = unit.Name
		return ds.ProductUnitRepo().Update(ctx, existingUnit)
	})
}

func (i *ProductUnitInteractor) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	user := domain.FromBackendUserContext(ctx)

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 检查属性是否存在
		existingUnit, err := ds.ProductUnitRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrUnitNotExists)
			}
			return err
		}
		if existingUnit.StoreID != user.Store.ID {
			return domain.ParamsError(domain.ErrUnitNotExists)
		}

		// 检查是否被商品关联
		used, err := ds.ProductRepo().Exists(ctx, domain.ProductExistsParams{UnitID: id})
		if err != nil {
			return err
		}
		if used {
			return domain.ParamsError(domain.ErrUnitUsing)
		}

		// 3. 执行删除
		return ds.ProductUnitRepo().Delete(ctx, id)
	})
}

func (i *ProductUnitInteractor) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.UnitSearchParams,
) (res *domain.UnitSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 获取当前店铺ID
	user := domain.FromBackendUserContext(ctx)
	params.StoreID = user.Store.ID

	// 3. 调用仓储层
	res, err = i.ds.ProductUnitRepo().PagedListBySearch(ctx, page, params)
	if err != nil {
		return nil, err
	}

	// 4. 返回标准化分页结构
	return &domain.UnitSearchRes{
		Pagination: page,
		Items:      res.Items,
	}, nil
}
