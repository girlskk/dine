package category

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *CategoryInteractor) Update(ctx context.Context, category *domain.Category) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证分类存在
		existingCategory, err := ds.CategoryRepo().FindByID(ctx, category.ID)
		if err != nil {
			return err
		}

		// 2. 验证更新后的分类名称在同一层级下唯一（排除自身）
		if category.Name != existingCategory.Name {
			exists, err := ds.CategoryRepo().Exists(ctx, domain.CategoryExistsParams{
				MerchantID: existingCategory.MerchantID,
				Name:       category.Name,
				ParentID:   existingCategory.ParentID,
				IsRoot:     existingCategory.IsRoot(),
				ExcludeID:  category.ID,
			})
			if err != nil {
				return err
			}
			if exists {
				return domain.ErrCategoryNameExists
			}
		}

		// 3. 将请求数据合并到 existingCategory
		existingCategory.Name = category.Name
		existingCategory.InheritTaxRate = category.InheritTaxRate
		existingCategory.InheritStall = category.InheritStall
		existingCategory.TaxRateID = category.TaxRateID
		existingCategory.StallID = category.StallID

		// 4. 处理继承逻辑（如果是子分类）
		if !existingCategory.IsRoot() {
			if category.InheritTaxRate {
				existingCategory.TaxRateID = uuid.Nil
			}
			if category.InheritStall {
				existingCategory.StallID = uuid.Nil
			}
		} else {
			// 一级分类不能有继承逻辑
			existingCategory.InheritTaxRate = false
			existingCategory.InheritStall = false
		}

		// 5. 验证税率ID和出品部门ID的有效性（如果提供了）
		if existingCategory.TaxRateID != uuid.Nil {
			// @TODO: 验证税率ID是否存在且可用
			// exists, err := ds.TaxRateRepo().Exists(ctx, domain.TaxRateExistsParams{
			// 	MerchantID: existingCategory.MerchantID,
			// 	ID:         existingCategory.TaxRateID,
			// })
			// if err != nil {
			// 	return err
			// }
			// if !exists {
			// 	return domain.ParamsError(domain.ErrTaxRateNotExists)
			// }
		}

		if existingCategory.StallID != uuid.Nil {
			// @TODO: 验证出品部门ID是否存在且可用
			// exists, err := ds.StallRepo().Exists(ctx, domain.StallExistsParams{
			// 	MerchantID: existingCategory.MerchantID,
			// 	ID:         existingCategory.StallID,
			// })
			// if err != nil {
			// 	return err
			// }
			// if !exists {
			// 	return domain.ParamsError(domain.ErrStallNotExists)
			// }
		}

		// 6. 执行更新操作
		err = ds.CategoryRepo().Update(ctx, existingCategory)
		if err != nil {
			return err
		}

		return nil
	})
}
