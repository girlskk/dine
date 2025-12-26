package menu

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.MenuInteractor = (*MenuInteractor)(nil)

type MenuInteractor struct {
	DS domain.DataStore
}

func NewMenuInteractor(ds domain.DataStore) *MenuInteractor {
	return &MenuInteractor{
		DS: ds,
	}
}

func (i *MenuInteractor) Create(ctx context.Context, menu *domain.Menu) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MenuInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 业务规则校验
		if err = validateMenuBusinessRules(ctx, ds, menu, uuid.Nil); err != nil {
			return err
		}

		// 创建菜单
		return ds.MenuRepo().Create(ctx, menu)
	})
}

func (i *MenuInteractor) Update(ctx context.Context, menu *domain.Menu, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MenuInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证菜单存在
		menu, err := ds.MenuRepo().FindByID(ctx, menu.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrMenuNotExists)
			}
			return err
		}
		if menu.MerchantID != user.GetMerchantID() {
			return domain.ParamsError(domain.ErrMenuNotExists)
		}

		// 业务规则校验（排除自身）
		if err = validateMenuBusinessRules(ctx, ds, menu, menu.ID); err != nil {
			return err
		}
		// 更新菜单
		return ds.MenuRepo().Update(ctx, menu)
	})
}

func (i *MenuInteractor) Delete(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MenuInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证菜单存在
		menu, err := ds.MenuRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrMenuNotExists)
			}
			return err
		}
		if menu.MerchantID != user.GetMerchantID() {
			return domain.ParamsError(domain.ErrMenuNotExists)
		}
		if menu.StoreCount > 0 {
			return domain.ParamsError(domain.ErrMenuHasStores)
		}

		// 删除菜单
		return ds.MenuRepo().Delete(ctx, id)
	})
}

func (i *MenuInteractor) GetDetail(ctx context.Context, id uuid.UUID, user domain.User) (res *domain.Menu, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MenuInteractor.GetDetail")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	menu, err := i.DS.MenuRepo().GetDetail(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ParamsError(domain.ErrMenuNotExists)
		}
		return nil, err
	}
	if menu.MerchantID != user.GetMerchantID() {
		return nil, domain.ParamsError(domain.ErrMenuNotExists)
	}
	return menu, nil
}

func (i *MenuInteractor) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.MenuSearchParams,
) (res *domain.MenuSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MenuInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.MenuRepo().PagedListBySearch(ctx, page, params)
}

// ============================================
// 校验函数
// ============================================

func validateMenuBusinessRules(ctx context.Context, ds domain.DataStore, menu *domain.Menu, excludeMenuID uuid.UUID) error {
	// 1. 检查菜单名称是否唯一
	exists, err := ds.MenuRepo().Exists(ctx, domain.MenuExistsParams{
		MerchantID: menu.MerchantID,
		Name:       menu.Name,
		ExcludeID:  excludeMenuID,
	})
	if err != nil {
		return err
	}
	if exists {
		return domain.ParamsError(domain.ErrMenuNameExists)
	}

	storeIDs := lo.Map(menu.Stores, func(store *domain.StoreSimple, _ int) uuid.UUID {
		return store.ID
	})

	// 2. 检查门店是否有效且属于当前品牌商
	// @TODO

	// 3. 检查门店是否已绑定其他菜单
	hasBound, err := ds.MenuRepo().CheckStoreBound(ctx, storeIDs, excludeMenuID)
	if err != nil {
		return err
	}
	if hasBound {
		return domain.ParamsError(domain.ErrMenuStoreBound)
	}

	// 4. 检查菜品是否有效
	productIDs := lo.Map(menu.Items, func(item *domain.MenuItem, _ int) uuid.UUID {
		return item.ProductID
	})
	products, err := ds.ProductRepo().ListByIDs(ctx, productIDs)
	if err != nil {
		return err
	}
	productMap := lo.SliceToMap(products, func(product *domain.Product) (uuid.UUID, *domain.Product) {
		return product.ID, product
	})

	for _, item := range menu.Items {
		product, ok := productMap[item.ProductID]
		if !ok {
			return domain.ParamsError(fmt.Errorf("菜品ID %s 不存在", item.ProductID))
		}

		// 检查菜品是否属于当前品牌商
		if product.MerchantID != menu.MerchantID {
			return domain.ParamsError(domain.ErrMenuItemProductInvalid)
		}
	}

	return nil
}
