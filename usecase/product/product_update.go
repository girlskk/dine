package product

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *ProductInteractor) Update(ctx context.Context, product *domain.Product) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 参数校验
	if err = validateProductParams(product); err != nil {
		return err
	}

	// 在事务中执行更新操作
	err = i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证商品是否存在
		existingProduct, err := ds.ProductRepo().FindByID(ctx, product.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductNotExists)
			}
			return err
		}

		// 验证商品是否属于当前商户
		if existingProduct.MerchantID != product.MerchantID {
			return domain.ParamsError(domain.ErrProductNotExists)
		}

		// 业务规则校验
		if err = validateProductBusinessRules(ctx, ds, product, existingProduct.ID); err != nil {
			return err
		}

		// 更新商品基本信息（同时清除旧的关联）
		if err = ds.ProductRepo().Update(ctx, product); err != nil {
			return err
		}

		// 创建新的规格关联
		if len(product.SpecRelations) > 0 {
			if err = ds.ProductSpecRelRepo().CreateBulk(ctx, product.SpecRelations); err != nil {
				return err
			}
		}

		// 创建新的口味做法关联
		if len(product.AttrRelations) > 0 {
			if err = ds.ProductAttrRelRepo().CreateBulk(ctx, product.AttrRelations); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
