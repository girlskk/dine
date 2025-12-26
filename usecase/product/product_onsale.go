package product

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *ProductInteractor) OnSale(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductInteractor.OnSale")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证商品存在
		product, err := ds.ProductRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductNotExists)
			}
			return err
		}

		// 2. 验证商品当前状态是否为"停售"
		if product.SaleStatus != domain.ProductSaleStatusOffSale {
			return domain.ParamsError(domain.ErrProductNotExists)
		}

		// 3. 更新商品售卖状态为"在售"
		product.SaleStatus = domain.ProductSaleStatusOnSale
		err = ds.ProductRepo().Update(ctx, product)
		if err != nil {
			return err
		}

		return nil
	})
}
