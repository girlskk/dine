package category

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

var _ domain.CategoryInteractor = (*CategoryInteractor)(nil)

type CategoryInteractor struct {
	DS domain.DataStore
}

func NewCategoryInteractor(ds domain.DataStore) *CategoryInteractor {
	return &CategoryInteractor{
		DS: ds,
	}
}

// func (interactor *CategoryInteractor) Update(ctx context.Context, params domain.CategoryUpdateParams) (res *domain.Category, err error) {
// 	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.Update")
// 	defer func() {
// 		util.SpanErrFinish(span, err)
// 	}()

// 	// 查找分类
// 	category, err := interactor.DS.CategoryRepo().FindByID(ctx, params.ID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to find category: %w", err)
// 	}

// 	// 如果名称有变化，需要验证新名称是否已存在
// 	if category.Name != params.Name {
// 		exists, err := interactor.DS.CategoryRepo().Exists(ctx, domain.CategoryExistsParams{
// 			StoreID:  category.StoreID,
// 			Name:     params.Name,
// 			ParentID: category.ParentID,
// 		})
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to check category exists: %w", err)
// 		}
// 		if exists {
// 			return nil, domain.ConflictError(domain.ErrCategoryNameExists)
// 		}
// 	}

// 	// 更新分类
// 	category.Name = params.Name
// 	category.TaxRateID = params.TaxRateID
// 	category.DepartmentID = params.DepartmentID

// 	err = interactor.DS.CategoryRepo().Update(ctx, category)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update category: %w", err)
// 	}

// 	// 重新查询以获取完整信息
// 	res, err = interactor.DS.CategoryRepo().FindByID(ctx, params.ID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to find updated category: %w", err)
// 	}

// 	return res, nil
// }

// func (interactor *CategoryInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
// 	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.Delete")
// 	defer func() {
// 		util.SpanErrFinish(span, err)
// 	}()

// 	// 查找分类
// 	category, err := interactor.DS.CategoryRepo().FindByID(ctx, id)
// 	if err != nil {
// 		return fmt.Errorf("failed to find category: %w", err)
// 	}

// 	// 如果分类下有关联商品，不能删除
// 	if category.ProductCount > 0 {
// 		return domain.ParamsError(fmt.Errorf("分类下有关联商品，不能删除"))
// 	}

// 	// 如果是一级分类，需要检查是否有子分类
// 	if category.IsRoot() && len(category.Children) > 0 {
// 		return domain.ParamsError(fmt.Errorf("一级分类下有子分类，不能删除"))
// 	}

// 	// 删除分类
// 	err = interactor.DS.CategoryRepo().Delete(ctx, id)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete category: %w", err)
// 	}

// 	return nil
// }

// func (interactor *CategoryInteractor) GetByID(ctx context.Context, id uuid.UUID) (res *domain.Category, err error) {
// 	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.GetByID")
// 	defer func() {
// 		util.SpanErrFinish(span, err)
// 	}()

// 	res, err = interactor.DS.CategoryRepo().FindByID(ctx, id)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to find category: %w", err)
// 	}

// 	return res, nil
// }

// func (interactor *CategoryInteractor) ListByStoreID(ctx context.Context, storeID uuid.UUID, parentID *uuid.UUID) (res domain.Categories, err error) {
// 	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.ListByStoreID")
// 	defer func() {
// 		util.SpanErrFinish(span, err)
// 	}()

// 	res, err = interactor.DS.CategoryRepo().ListByStoreID(ctx, storeID, parentID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list categories: %w", err)
// 	}

// 	return res, nil
// }
