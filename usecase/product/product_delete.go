package product

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *ProductInteractor) Delete(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductInteractor.Delete")
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

		// 2. 验证是否可以操作该商品
		if err := verifyProductOwnership(user, product); err != nil {
			return err
		}

		// 3. 删除商品
		err = ds.ProductRepo().Delete(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}
