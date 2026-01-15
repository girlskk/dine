package category

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *CategoryInteractor) CreateRoot(ctx context.Context, category *domain.Category, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.CreateRoot")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 分类名称唯一性检查
		exists, err := ds.CategoryRepo().Exists(ctx, domain.CategoryExistsParams{
			MerchantID: category.MerchantID,
			StoreID:    category.StoreID,
			Name:       category.Name,
			IsRoot:     true,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrCategoryNameExists
		}
		// 检查税率是否有效
		err = i.checkTaxRate(ctx, ds, category, user)
		if err != nil {
			return err
		}
		// 检查出品部门是否有效
		err = i.checkStall(ctx, ds, category, user)
		if err != nil {
			return err
		}

		// 创建一级分类
		err = ds.CategoryRepo().Create(ctx, category)
		if err != nil {
			return err
		}

		// 创建子分类
		if len(category.Childrens) > 0 {
			err = ds.CategoryRepo().CreateBulk(ctx, category.Childrens)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (i *CategoryInteractor) CreateChild(ctx context.Context, category *domain.Category, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.CreateChild")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证父分类存在
		parentCategory, err := ds.CategoryRepo().FindByID(ctx, category.ParentID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrCategoryParentNotExists)
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

		// 4. 验证父分类是否属于当前用户可操作
		if !domain.VerifyOwnerShip(user, parentCategory.MerchantID, parentCategory.StoreID) {
			return domain.ParamsError(domain.ErrCategoryParentNotExists)
		}

		// 5. 验证同一父分类下名称唯一性
		exists, err := ds.CategoryRepo().Exists(ctx, domain.CategoryExistsParams{
			MerchantID: category.MerchantID,
			StoreID:    category.StoreID,
			Name:       category.Name,
			ParentID:   category.ParentID,
			IsRoot:     false,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrCategoryNameExists
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

		// 检查税率是否有效
		err = i.checkTaxRate(ctx, ds, category, user)
		if err != nil {
			return err
		}
		// 检查出品部门是否有效
		err = i.checkStall(ctx, ds, category, user)
		if err != nil {
			return err
		}

		// 创建二级分类
		err = ds.CategoryRepo().Create(ctx, category)
		if err != nil {
			return err
		}

		return nil
	})
}
