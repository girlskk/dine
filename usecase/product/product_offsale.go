package product

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *ProductInteractor) OffSale(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductInteractor.OffSale")
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

		// 2. 验证商品当前状态是否为"在售"
		if product.SaleStatus != domain.ProductSaleStatusOnSale {
			return domain.ParamsError(domain.ErrProductNotExists)
		}

		// 3. 更新商品售卖状态为"停售"
		product.SaleStatus = domain.ProductSaleStatusOffSale
		err = ds.ProductRepo().Update(ctx, product)
		if err != nil {
			return err
		}

		return nil
	})
}
