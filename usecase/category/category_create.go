package category

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *CategoryInteractor) CreateRoot(ctx context.Context, category *domain.Category) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.CreateRoot")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// Check if the category name already exists
		exists, err := ds.CategoryRepo().Exists(ctx, domain.CategoryExistsParams{
			MerchantID: category.MerchantID,
			Name:       category.Name,
			IsRoot:     true,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ConflictError(domain.ErrCategoryNameExists)
		}

		// check tax rate id
		if category.TaxRateID != uuid.Nil {
			// @TODO
			// exists, err := ds.TaxRateRepo().Exists(ctx, domain.TaxRateExistsParams{
			// 	MerchantID: category.MerchantID,
			// 	ID:         category.TaxRateID,
			// })
			// if err != nil {
			// 	return err
			// }
			// if !exists {
			// 	return domain.NotFoundError(domain.ErrTaxRateNotExists)
			// }
		}

		if category.StallID != uuid.Nil {
			// @TODO
			// exists, err := ds.StallRepo().Exists(ctx, domain.StallExistsParams{
			// 	MerchantID: category.MerchantID,
			// 	ID:         category.StallID,
			// })
			// if err != nil {
			// 	return err
			// }
			// if !exists {
			// 	return domain.NotFoundError(domain.ErrStallNotExists)
			// }
		}

		// create root category
		err = ds.CategoryRepo().Create(ctx, category)
		if err != nil {
			return err
		}

		// create children categories
		if len(category.Childrens) > 0 {
			err = ds.CategoryRepo().CreateBulk(ctx, category.Childrens)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (i *CategoryInteractor) CreateChild(ctx context.Context, category *domain.Category) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.CreateChild")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证父分类存在
		parentCategory, err := ds.CategoryRepo().FindByID(ctx, category.ParentID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.NotFoundError(domain.ErrCategoryParentNotExists)
			}
			return err
		}

		// 2. 验证父分类是一级分类（不能是二级分类）
		if !parentCategory.IsRoot() {
			return domain.ParamsError(domain.ErrCategoryInvalidLevel)
		}

		// 3. 验证父分类下没有关联商品
		if parentCategory.ProductCount > 0 {
			return domain.ParamsError(domain.ErrCategoryParentHasProducts)
		}

		// 4. 验证同一父分类下名称唯一性
		exists, err := ds.CategoryRepo().Exists(ctx, domain.CategoryExistsParams{
			MerchantID: category.MerchantID,
			Name:       category.Name,
			ParentID:   category.ParentID,
			IsRoot:     false,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ConflictError(domain.ErrCategoryNameExists)
		}

		// 5. 处理继承逻辑
		// 如果继承税率，则 TaxRateID 保持为 uuid.Nil，表示继承父分类的值
		if category.InheritTaxRate {
			category.TaxRateID = uuid.Nil
		}

		// 如果继承出品部门，则 StallID 保持为 uuid.Nil，表示继承父分类的值
		if category.InheritStall {
			category.StallID = uuid.Nil
		}

		// 6. 验证税率ID和出品部门ID的有效性（如果提供了）
		if category.TaxRateID != uuid.Nil {
			// @TODO: 验证税率ID是否存在且可用
			// exists, err := ds.TaxRateRepo().Exists(ctx, domain.TaxRateExistsParams{
			// 	MerchantID: category.MerchantID,
			// 	ID:         category.TaxRateID,
			// })
			// if err != nil {
			// 	return err
			// }
			// if !exists {
			// 	return domain.NotFoundError(domain.ErrTaxRateNotExists)
			// }
		}

		if category.StallID != uuid.Nil {
			// @TODO: 验证出品部门ID是否存在且可用
			// exists, err := ds.StallRepo().Exists(ctx, domain.StallExistsParams{
			// 	MerchantID: category.MerchantID,
			// 	ID:         category.StallID,
			// })
			// if err != nil {
			// 	return err
			// }
			// if !exists {
			// 	return domain.NotFoundError(domain.ErrStallNotExists)
			// }
		}

		// 创建二级分类
		err = ds.CategoryRepo().Create(ctx, category)
		if err != nil {
			return err
		}

		return nil
	})
}
