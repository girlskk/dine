package usecase

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.OrderCartInteractor = (*OrderCartInteractor)(nil)

type OrderCartInteractor struct {
	ds domain.DataStore
}

func NewOrderCartInteractor(dataStore domain.DataStore) *OrderCartInteractor {
	return &OrderCartInteractor{
		ds: dataStore,
	}
}

func (i *OrderCartInteractor) ListByTable(ctx context.Context, tableID int) (res domain.OrderCarts, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderCartInteractor.ListByTable")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return i.ds.OrderCartRepo().ListByTable(ctx, tableID, true)
}

func (i *OrderCartInteractor) AddItem(ctx context.Context,
	params domain.OrderCartAddParams,
) (items domain.OrderCarts, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderCartInteractor.AddItem")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 检查商品
		err = i.checkProduct(ctx, ds, params)
		if err != nil {
			return err
		}
		// 查找是否存在完全匹配的项目（商品ID、规格ID、属性ID和做法ID都相同）
		existsItem, err := ds.OrderCartRepo().FindByUniqueKey(ctx, domain.OrderCartItemUniqueKey(params))
		if err != nil && !domain.IsNotFound(err) {
			return err
		}
		// 如果找到完全匹配的项目，则更新数量
		if existsItem != nil {
			err = ds.OrderCartRepo().IncrementQuantity(ctx, existsItem.ID)
			if err != nil {
				return err
			}
		} else {
			newItem := &domain.OrderCart{
				TableID:       params.TableID,
				ProductID:     params.ProductID,
				ProductSpecID: params.ProductSpecID,
				AttrID:        params.AttrID,
				RecipeID:      params.RecipeID,
				Quantity:      decimal.NewFromInt(1),
			}
			err = ds.OrderCartRepo().Create(ctx, newItem)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return i.ds.OrderCartRepo().ListByTable(ctx, params.TableID, true)
}

func (i *OrderCartInteractor) checkProduct(ctx context.Context,
	ds domain.DataStore, params domain.OrderCartAddParams,
) (err error) {
	// 检查商品是否ok
	product, err := ds.ProductRepo().FindByID(ctx, params.ProductID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrProductNotExists)
		}
		return err
	}
	if product.Status != domain.ProductStatusApproved {
		return domain.ParamsError(domain.ErrProductStatus)
	}
	if product.SaleStatus == domain.ProductSaleStatusOff {
		return domain.ParamsError(domain.ErrProductNotOnSale)
	}

	// 检查商品规格是否ok
	if product.Type == domain.ProductTypeMulti {
		if params.ProductSpecID == 0 {
			return domain.ParamsError(domain.ErrProductTypeMulti)
		}
		productSpec, err := ds.ProductSpecRelRepo().FindByID(ctx, params.ProductSpecID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}
		if productSpec.ProductID != product.ID {
			return domain.ParamsError(domain.ErrSpecNotExists)
		}
		if productSpec.SaleStatus == domain.ProductSaleStatusOff {
			return domain.ParamsError(domain.ErrProductNotOnSale)
		}
	}
	return nil
}

func (i *OrderCartInteractor) RemoveItem(ctx context.Context, id int, tableID int) (items domain.OrderCarts, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderCartInteractor.RemoveItem")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		item, err := ds.OrderCartRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
			return err
		}
		if item.TableID != tableID {
			return domain.ParamsError(domain.ErrOrderCartNotFound)
		}
		return ds.OrderCartRepo().DecrementQuantity(ctx, id)
	})

	if err != nil {
		return nil, err
	}

	return i.ds.OrderCartRepo().ListByTable(ctx, tableID, true)
}
