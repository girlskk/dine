package product

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *ProductInteractor) UpdateSetMeal(ctx context.Context, product *domain.Product, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductInteractor.UpdateSetMeal")
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

		// 验证是否可以操作该商品
		if err := verifyProductOwnership(user, existingProduct); err != nil {
			return err
		}

		// 业务规则校验
		if err = validateProductBusinessRules(ctx, ds, product, existingProduct.ID); err != nil {
			return err
		}

		// 删除旧的套餐组和套餐组详情（物理删除）
		if err = ds.SetMealGroupRepo().DeleteByProductID(ctx, product.ID); err != nil {
			return err
		}

		// 更新商品基本信息（同时清除旧的规格关联）
		if err = ds.ProductRepo().Update(ctx, product); err != nil {
			return err
		}

		// 创建新的规格关联
		if len(product.SpecRelations) > 0 {
			if err = ds.ProductSpecRelRepo().CreateBulk(ctx, product.SpecRelations); err != nil {
				return err
			}
		}

		// 批量创建套餐组
		if len(product.Groups) > 0 {
			if err = ds.SetMealGroupRepo().CreateGroups(ctx, product.Groups); err != nil {
				return err
			}
		}

		// 批量创建套餐组详情
		allDetails := make([]*domain.SetMealDetail, 0)
		for _, group := range product.Groups {
			allDetails = append(allDetails, group.Details...)
		}
		if len(allDetails) > 0 {
			if err = ds.SetMealGroupRepo().CreateDetails(ctx, allDetails); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
