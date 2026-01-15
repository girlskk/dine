package product

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *ProductInteractor) Create(ctx context.Context, product *domain.Product, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if err = validateProductParams(product); err != nil {
		return err
	}

	err = i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		if err = i.validateProductBusinessRules(ctx, ds, product, user, uuid.Nil); err != nil {
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

		// 创建口味做法关联
		if len(product.AttrRelations) > 0 {
			err = ds.ProductAttrRelRepo().CreateBulk(ctx, product.AttrRelations)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
