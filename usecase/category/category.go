package category

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
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

func (i *CategoryInteractor) ListBySearch(ctx context.Context, params domain.CategorySearchParams) (res domain.Categories, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.ListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.CategoryRepo().ListBySearch(ctx, params)
}

// verifyCategoryOwnership 验证分类是否属于当前用户可操作
func verifyCategoryOwnership(user domain.User, category *domain.Category) error {
	if !domain.VerifyOwnerShip(user, category.MerchantID, category.StoreID) {
		return domain.ParamsError(domain.ErrCategoryNotExists)
	}
	return nil
}

func (i *CategoryInteractor) Reorder(ctx context.Context, parentID *uuid.UUID, categoryIDs []uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.Reorder")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	if len(categoryIDs) == 0 {
		return domain.ParamsError(fmt.Errorf("分类ID列表不能为空"))
	}
	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		merchantID := user.GetMerchantID()
		storeID := user.GetStoreID()
		actualParentID := uuid.Nil

		if parentID != nil {
			actualParentID = *parentID
			// 验证父分类存在且属于当前用户
			parent, err := ds.CategoryRepo().FindByID(ctx, actualParentID)
			if err != nil {
				return err
			}

			if err := verifyCategoryOwnership(user, parent); err != nil {
				return err
			}
		}
		// 1. 查询同级分类
		sameLevelCategories, err := ds.CategoryRepo().ListByParentID(ctx, merchantID, storeID, actualParentID)
		if err != nil {
			return err
		}
		// 3. 验证所有传入的 categoryIDs 都存在且属于同级
		if len(sameLevelCategories) != len(categoryIDs) {
			return domain.ParamsError(fmt.Errorf("部分分类ID无效或不属于指定的父分类"))
		}

		// 2. 建立ID到分类的映射
		categoryIDMap := lo.SliceToMap(sameLevelCategories, func(cat *domain.Category) (uuid.UUID, *domain.Category) {
			return cat.ID, cat
		})

		for _, id := range categoryIDs {
			if _, exists := categoryIDMap[id]; !exists {
				return domain.ParamsError(fmt.Errorf("分类ID %s 不存在或不属于指定的父分类", id))
			}
		}

		// 4. 计算新的 sort_order 值
		newSortMap := make(map[uuid.UUID]int)
		for i, id := range categoryIDs {
			newSortMap[id] = i + 1
		}

		// 5. 对比现有顺序，找出需要更新的分类
		updates := make(map[uuid.UUID]int)
		for _, id := range categoryIDs {
			cat := categoryIDMap[id]
			newSortOrder := newSortMap[id]

			if cat.SortOrder != newSortOrder {
				updates[id] = newSortOrder
			}
		}

		// 6. 如果没有需要更新的分类，直接返回
		if len(updates) == 0 {
			return nil
		}

		// 7. 调用 Repository 批量更新
		return ds.CategoryRepo().UpdateSortOrders(ctx, updates)
	})
}

// 检查税率是否有效
func (i *CategoryInteractor) checkTaxRate(ctx context.Context, ds domain.DataStore,
	category *domain.Category, user domain.User,
) error {
	if category.TaxRateID == uuid.Nil {
		return nil
	}
	taxRate, err := ds.TaxFeeRepo().FindByID(ctx, category.TaxRateID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrTaxFeeNotExists)
		}
		return err
	}
	if !domain.VerifyOwnerShip(user, taxRate.MerchantID, taxRate.StoreID) {
		return domain.ParamsError(domain.ErrTaxFeeNotExists)
	}
	return nil
}

// 检查出品部门是否有效
func (i *CategoryInteractor) checkStall(ctx context.Context, ds domain.DataStore,
	category *domain.Category, user domain.User,
) error {
	if category.StallID == uuid.Nil {
		return nil
	}
	stall, err := ds.StallRepo().FindByID(ctx, category.StallID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrStallNotExists)
		}
		return err
	}

	if !domain.VerifyOwnerShip(user, stall.MerchantID, stall.StoreID) {
		return domain.ParamsError(domain.ErrStallNotExists)
	}
	return nil
}
