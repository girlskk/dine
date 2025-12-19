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

func (i *ProductAttrInteractor) Update(ctx context.Context, attr *domain.ProductAttr) (err error) {
	// span, ctx := util.StartSpan(ctx, "usecase", "ProductAttrInteractor.Update")
	// defer func() {
	// 	util.SpanErrFinish(span, err)
	// }()

	// return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
	// 	// 1. 验证口味做法存在
	// 	existingAttr, err := ds.ProductAttrRepo().FindByID(ctx, attr.ID)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// 2. 验证更新后的名称在当前门店下是否唯一（排除自身）
	// 	if attr.Name != existingAttr.Name {
	// 		exists, err := ds.ProductAttrRepo().Exists(ctx, domain.ProductAttrExistsParams{
	// 			MerchantID: existingAttr.MerchantID,
	// 			Name:       attr.Name,
	// 			ExcludeID:  attr.ID,
	// 		})
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if exists {
	// 			return domain.ErrProductAttrNameExists
	// 		}
	// 	}

	// 	// 3. 更新口味做法基本信息
	// 	existingAttr.Name = attr.Name
	// 	existingAttr.Channels = attr.Channels
	// 	err = ds.ProductAttrRepo().Update(ctx, existingAttr)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// 4. 处理口味做法项的新增、修改、删除
	// 	// 获取现有的项
	// 	existingItems, err := ds.ProductAttrRepo().ListItemsByAttrID(ctx, attr.ID)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	existingItemMap := make(map[uuid.UUID]*domain.ProductAttrItem)
	// 	for _, item := range existingItems {
	// 		existingItemMap[item.ID] = item
	// 	}

	// 	// 处理传入的项
	// 	itemMap := make(map[uuid.UUID]*domain.ProductAttrItem)
	// 	var itemsToCreate []*domain.ProductAttrItem
	// 	var itemsToUpdate []*domain.ProductAttrItem

	// 	for _, item := range attr.Items {
	// 		// 验证基础加价必须为非负数
	// 		if item.BasePrice.IsNegative() {
	// 			return domain.ErrProductAttrItemBasePriceInvalid
	// 		}

	// 		item.AttrID = attr.ID
	// 		if item.ID == uuid.Nil {
	// 			// 新项，需要创建
	// 			item.ID = uuid.New()
	// 			itemsToCreate = append(itemsToCreate, item)
	// 		} else {
	// 			// 已存在的项，需要更新
	// 			if _, exists := existingItemMap[item.ID]; !exists {
	// 				return domain.ErrProductAttrItemNotExists
	// 			}
	// 			itemsToUpdate = append(itemsToUpdate, item)
	// 		}
	// 		itemMap[item.ID] = item

	// 		// 验证名称在口味做法下唯一（排除自身）
	// 		exists, err := ds.ProductAttrRepo().ItemExists(ctx, domain.ProductAttrItemExistsParams{
	// 			AttrID:    attr.ID,
	// 			Name:      item.Name,
	// 			ExcludeID: item.ID,
	// 		})
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if exists {
	// 			return domain.ErrProductAttrItemNameExists
	// 		}
	// 	}

	// 	// 找出需要删除的项
	// 	var itemsToDelete []uuid.UUID
	// 	for _, existingItem := range existingItems {
	// 		if _, exists := itemMap[existingItem.ID]; !exists {
	// 			itemsToDelete = append(itemsToDelete, existingItem.ID)
	// 		}
	// 	}

	// 	// 执行删除
	// 	for _, itemID := range itemsToDelete {
	// 		// 检查是否有关联商品
	// 		item, err := ds.ProductAttrRepo().FindItemByID(ctx, itemID)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if item.ProductCount > 0 {
	// 			return domain.ErrProductAttrItemDeleteHasProducts
	// 		}
	// 		err = ds.ProductAttrRepo().DeleteItem(ctx, itemID)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}

	// 	// 执行创建
	// 	if len(itemsToCreate) > 0 {
	// 		err = ds.ProductAttrRepo().CreateItems(ctx, itemsToCreate)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}

	// 	// 执行更新
	// 	for _, item := range itemsToUpdate {
	// 		err = ds.ProductAttrRepo().UpdateItem(ctx, item)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}

	// 	// 5. 更新返回的 attr 对象
	// 	attr.UpdatedAt = existingAttr.UpdatedAt

	// 	return nil
	// })
	return nil
}

func (i *ProductAttrInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
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

		if len(attr.Items) > 0 {
			return domain.ErrProductAttrDeleteHasItems
		}

		// 2. 删除口味做法
		err = ds.ProductAttrRepo().Delete(ctx, id)

		return err
	})
}

func (i *ProductAttrInteractor) DeleteItem(ctx context.Context, id uuid.UUID) (err error) {
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

		// 2. 如果口味做法项下有关联商品，不能删除
		if item.ProductCount > 0 {
			return domain.ErrProductAttrItemDeleteHasProducts
		}

		// 3. 删除口味做法项
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
