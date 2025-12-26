package productattr

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductAttrInteractor = (*ProductAttrInteractor)(nil)

type ProductAttrInteractor struct {
	DS domain.DataStore
}

func NewProductAttrInteractor(ds domain.DataStore) *ProductAttrInteractor {
	return &ProductAttrInteractor{
		DS: ds,
	}
}

func (i *ProductAttrInteractor) Create(ctx context.Context, attr *domain.ProductAttr) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductAttrInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证名称在当前门店下是否唯一
		exists, err := ds.ProductAttrRepo().Exists(ctx, domain.ProductAttrExistsParams{
			MerchantID: attr.MerchantID,
			StoreID:    attr.StoreID,
			Name:       attr.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrProductAttrNameExists
		}

		// 3. 创建商品口味做法
		err = ds.ProductAttrRepo().Create(ctx, attr)
		if err != nil {
			return err
		}

		// 4. 创建口味做法项
		if len(attr.Items) > 0 {
			err = ds.ProductAttrRepo().CreateItems(ctx, attr.Items)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (i *ProductAttrInteractor) Update(ctx context.Context, attr *domain.ProductAttr, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductAttrInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 使用 GetDetail 加载所有子项
		existingAttr, err := ds.ProductAttrRepo().GetDetail(ctx, attr.ID)
		if err != nil {
			return err
		}

		if err := verifyProductAttrOwnership(user, existingAttr); err != nil {
			return err
		}

		// 2. 验证更新后的名称在当前门店下是否唯一（排除自身）
		if attr.Name != existingAttr.Name {
			exists, err := ds.ProductAttrRepo().Exists(ctx, domain.ProductAttrExistsParams{
				MerchantID: existingAttr.MerchantID,
				StoreID:    existingAttr.StoreID,
				Name:       attr.Name,
				ExcludeID:  attr.ID,
			})
			if err != nil {
				return err
			}
			if exists {
				return domain.ErrProductAttrNameExists
			}
		}

		// 3. 处理口味做法项
		existingItemMap := make(map[uuid.UUID]*domain.ProductAttrItem)
		for _, item := range existingAttr.Items {
			existingItemMap[item.ID] = item
		}

		// 构建传入项的ID映射
		requestedItemMap := make(map[uuid.UUID]*domain.ProductAttrItem)
		for _, item := range attr.Items {
			// 如果传入了ID，检查是否存在于现有项中
			if item.ID != uuid.Nil {
				if _, exists := existingItemMap[item.ID]; !exists {
					return domain.ParamsError(domain.ErrProductAttrItemNotExists)
				}
			} else {
				// 新增项，生成ID
				item.ID = uuid.New()
			}
			requestedItemMap[item.ID] = item
		}

		// 4. 找出需要删除的项（存在于现有项但不在传入项中）
		var itemsToDelete []uuid.UUID
		for existingID, existingItem := range existingItemMap {
			if _, exists := requestedItemMap[existingID]; !exists {
				itemsToDelete = append(itemsToDelete, existingID)
				// 检查是否有关联商品
				if existingItem.ProductCount > 0 {
					return domain.ParamsError(domain.ErrProductAttrItemDeleteHasProducts)
				}
			}
		}

		// 5. 执行删除
		if len(itemsToDelete) > 0 {
			err = ds.ProductAttrRepo().DeleteItems(ctx, itemsToDelete)
			if err != nil {
				return err
			}
		}

		// 6. 更新口味做法基本信息
		existingAttr.Name = attr.Name
		existingAttr.Channels = attr.Channels
		err = ds.ProductAttrRepo().Update(ctx, existingAttr)
		if err != nil {
			return err
		}

		// 7. 批量保存口味做法项（新增和编辑统一处理）
		if len(attr.Items) > 0 {
			err = ds.ProductAttrRepo().SaveItems(ctx, attr.Items)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (i *ProductAttrInteractor) Delete(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductAttrInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 查找口味做法，验证口味做法存在
		attr, err := ds.ProductAttrRepo().GetDetail(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductAttrNotExists)
			}
			return err
		}
		if err := verifyProductAttrOwnership(user, attr); err != nil {
			return err
		}

		if len(attr.Items) > 0 {
			return domain.ErrProductAttrDeleteHasItems
		}

		// 2. 删除口味做法
		err = ds.ProductAttrRepo().Delete(ctx, id)

		return err
	})
}

func (i *ProductAttrInteractor) DeleteItem(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductAttrInteractor.DeleteItem")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 查找口味做法项，验证项存在
		item, err := ds.ProductAttrRepo().FindItemByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductAttrItemNotExists)
			}
			return err
		}

		// 2. 查找口味做法，验证口味做法存在，当前用户是否拥有该口味做法操作权限
		attr, err := ds.ProductAttrRepo().FindByID(ctx, item.AttrID)

		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductAttrNotExists)
			}
			return err
		}
		if err := verifyProductAttrOwnership(user, attr); err != nil {
			return err
		}

		// 3. 如果口味做法项下有关联商品，不能删除
		if item.ProductCount > 0 {
			return domain.ErrProductAttrItemDeleteHasProducts
		}

		// 4. 删除口味做法项
		err = ds.ProductAttrRepo().DeleteItem(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (i *ProductAttrInteractor) ListBySearch(
	ctx context.Context,
	params domain.ProductAttrSearchParams,
) (res domain.ProductAttrs, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductAttrInteractor.ListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return i.DS.ProductAttrRepo().ListBySearch(ctx, params)
}

func verifyProductAttrOwnership(user domain.User, attr *domain.ProductAttr) error {
	if user.GetMerchantID() != attr.MerchantID || user.GetStoreID() != attr.StoreID {
		return domain.ParamsError(domain.ErrProductAttrNotExists)
	}
	return nil
}
