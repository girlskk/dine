package category

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *CategoryInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 查找分类，验证分类存在
		category, err := ds.CategoryRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrCategoryNotExists)
			}
			return err
		}

		// 2. 如果分类下有关联商品，不能删除
		if category.ProductCount > 0 {
			return domain.ErrCategoryDeleteHasProducts
		}

		// 3. 如果是一级分类，需要检查是否有子分类
		if category.IsRoot() {
			childrenCount, err := ds.CategoryRepo().CountChildrenByParentID(ctx, id)
			if err != nil {
				return err
			}
			if childrenCount > 0 {
				return domain.ErrCategoryDeleteHasChildren
			}
		}

		// 4. 删除分类
		err = ds.CategoryRepo().Delete(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}
