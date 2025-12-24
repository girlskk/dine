package product

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *ProductInteractor) CreateSetMeal(ctx context.Context, product *domain.Product) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductInteractor.CreateSetMeal")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if err = validateProductParams(product); err != nil {
		return err
	}

	err = i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		if err = validateProductBusinessRules(ctx, ds, product, uuid.Nil); err != nil {
			return err
		}

		// 创建商品
		err := ds.ProductRepo().Create(ctx, product)
		if err != nil {
			return err
		}

		// 创建规格关联
		if len(product.SpecRelations) > 0 {
			err = ds.ProductSpecRelRepo().CreateBulk(ctx, product.SpecRelations)
			if err != nil {
				return err
			}
		}

		// 批量创建套餐组
		err = ds.SetMealGroupRepo().CreateGroups(ctx, product.Groups)
		if err != nil {
			return err
		}

		// 批量创建套餐组详情
		allDetails := make([]*domain.SetMealDetail, 0)
		for _, group := range product.Groups {
			allDetails = append(allDetails, group.Details...)
		}
		if len(allDetails) > 0 {
			err = ds.SetMealGroupRepo().CreateDetails(ctx, allDetails)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}
