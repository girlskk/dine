package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/menu"
	"gitlab.jiguang.dev/pos-dine/dine/ent/menuitem"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
	"gitlab.jiguang.dev/pos-dine/dine/ent/store"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.MenuRepository = (*MenuRepository)(nil)

type MenuRepository struct {
	Client *ent.Client
}

func NewMenuRepository(client *ent.Client) *MenuRepository {
	return &MenuRepository{
		Client: client,
	}
}

func (repo *MenuRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.Menu, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MenuRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	em, err := repo.Client.Menu.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrMenuNotExists)
		}
		return nil, err
	}

	res = convertMenuToDomain(em)
	return res, nil
}

func (repo *MenuRepository) GetDetail(ctx context.Context, id uuid.UUID) (res *domain.Menu, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MenuRepository.GetDetail")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	em, err := repo.Client.Menu.Query().
		Where(menu.IDEQ(id)).
		WithItems(func(query *ent.MenuItemQuery) {
			query.WithProduct()
		}).
		WithStores().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrMenuNotExists)
		}
		return nil, err
	}

	res = convertMenuToDomain(em)
	return res, nil
}

func (repo *MenuRepository) Create(ctx context.Context, m *domain.Menu) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MenuRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 创建菜单
	menuBuilder := repo.Client.Menu.Create().
		SetID(m.ID).
		SetMerchantID(m.MerchantID).
		SetName(m.Name).
		SetDistributionRule(m.DistributionRule).
		SetStoreCount(m.StoreCount).
		SetItemCount(m.ItemCount)

	// 设置关联门店（Many2Many）
	if len(m.Stores) > 0 {
		storeIDs := make([]uuid.UUID, 0, len(m.Stores))
		for _, store := range m.Stores {
			storeIDs = append(storeIDs, store.ID)
		}
		menuBuilder.AddStoreIDs(storeIDs...)
	}

	em, err := menuBuilder.Save(ctx)
	if err != nil {
		return err
	}

	// 创建菜单项
	if len(m.Items) > 0 {
		itemBuilders := make([]*ent.MenuItemCreate, 0, len(m.Items))
		for _, item := range m.Items {
			builder := repo.Client.MenuItem.Create().
				SetID(item.ID).
				SetMenuID(em.ID).
				SetProductID(item.ProductID).
				SetSaleRule(item.SaleRule)

			if item.BasePrice != nil {
				builder.SetBasePrice(*item.BasePrice)
			}
			if item.MemberPrice != nil {
				builder.SetMemberPrice(*item.MemberPrice)
			}
			itemBuilders = append(itemBuilders, builder)
		}
		_, err = repo.Client.MenuItem.CreateBulk(itemBuilders...).Save(ctx)
		if err != nil {
			return err
		}
	}

	m.CreatedAt = em.CreatedAt
	m.UpdatedAt = em.UpdatedAt
	return nil
}

func (repo *MenuRepository) Update(ctx context.Context, m *domain.Menu) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MenuRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 先删除旧的菜单项记录（使用 DELETE，而不是 SET NULL）
	skipSoftDeleteCtx := schematype.SkipSoftDelete(ctx)
	_, err = repo.Client.MenuItem.Delete().
		Where(menuitem.MenuIDEQ(m.ID)).
		Exec(skipSoftDeleteCtx)
	if err != nil {
		return err
	}

	// 更新菜单基本信息
	builder := repo.Client.Menu.UpdateOneID(m.ID).
		SetName(m.Name).
		SetDistributionRule(m.DistributionRule).
		SetStoreCount(m.StoreCount).
		SetItemCount(m.ItemCount)

	// 更新关联门店（Many2Many）
	if len(m.Stores) > 0 {
		storeIDs := make([]uuid.UUID, 0, len(m.Stores))
		for _, store := range m.Stores {
			storeIDs = append(storeIDs, store.ID)
		}
		builder = builder.AddStoreIDs(storeIDs...)
	} else {
		builder = builder.ClearStores()
	}

	_, err = builder.Save(ctx)
	if err != nil {
		return err
	}

	// 更新菜单项
	if len(m.Items) > 0 {
		itemBuilders := make([]*ent.MenuItemCreate, 0, len(m.Items))
		for _, item := range m.Items {
			builder := repo.Client.MenuItem.Create().
				SetID(item.ID).
				SetMenuID(m.ID).
				SetProductID(item.ProductID).
				SetSaleRule(item.SaleRule)

			if item.BasePrice != nil {
				builder.SetBasePrice(*item.BasePrice)
			}
			if item.MemberPrice != nil {
				builder.SetMemberPrice(*item.MemberPrice)
			}
			itemBuilders = append(itemBuilders, builder)
		}
		_, err = repo.Client.MenuItem.CreateBulk(itemBuilders...).Save(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *MenuRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MenuRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 删除菜单项
	_, err = repo.Client.MenuItem.Delete().Where(menuitem.MenuID(id)).Exec(ctx)
	if err != nil {
		return err
	}

	// 删除菜单
	err = repo.Client.Menu.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MenuRepository) Exists(ctx context.Context, params domain.MenuExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MenuRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Menu.Query()
	if params.MerchantID != uuid.Nil {
		query.Where(menu.MerchantID(params.MerchantID))
	}
	if params.Name != "" {
		query.Where(menu.Name(params.Name))
	}
	if params.ExcludeID != uuid.Nil {
		query.Where(menu.IDNEQ(params.ExcludeID))
	}
	return query.Exist(ctx)
}

func (repo *MenuRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.MenuSearchParams,
) (res *domain.MenuSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MenuRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Menu.Query()

	if params.MerchantID != uuid.Nil {
		query.Where(menu.MerchantID(params.MerchantID))
	}

	if params.Name != "" {
		query.Where(menu.NameContains(params.Name))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}

	// 分页处理
	query = query.
		Offset(page.Offset()).
		Limit(page.Size)

	// 按创建时间倒序排列
	entMenus, err := query.Order(ent.Desc(menu.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.Menus, 0, len(entMenus))
	for _, m := range entMenus {
		items = append(items, convertMenuToDomain(m))
	}

	page.SetTotal(total)

	return &domain.MenuSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

func (repo *MenuRepository) CheckStoreBound(ctx context.Context, storeIDs []uuid.UUID, excludeMenuID uuid.UUID) (has bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MenuRepository.CheckStoreBound")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(storeIDs) == 0 {
		return false, nil
	}

	query := repo.Client.Menu.Query().
		Where(menu.HasStoresWith(store.IDIn(storeIDs...))).
		WithStores()

	if excludeMenuID != uuid.Nil {
		query.Where(menu.IDNEQ(excludeMenuID))
	}

	menus, err := query.All(ctx)
	if err != nil {
		return false, err
	}

	return len(menus) > 0, nil
}

// ============================================
// 转换函数
// ============================================

func convertMenuToDomain(em *ent.Menu) *domain.Menu {
	if em == nil {
		return nil
	}

	m := &domain.Menu{
		ID:               em.ID,
		MerchantID:       em.MerchantID,
		Name:             em.Name,
		DistributionRule: em.DistributionRule,
		StoreCount:       em.StoreCount,
		ItemCount:        em.ItemCount,
		CreatedAt:        em.CreatedAt,
		UpdatedAt:        em.UpdatedAt,
	}

	// 转换菜单项
	if len(em.Edges.Items) > 0 {
		m.Items = make(domain.MenuItems, 0, len(em.Edges.Items))
		for _, item := range em.Edges.Items {
			itemDomain := &domain.MenuItem{
				ID:        item.ID,
				MenuID:    item.MenuID,
				ProductID: item.ProductID,
				SaleRule:  item.SaleRule,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}

			if item.BasePrice != nil {
				itemDomain.BasePrice = item.BasePrice
			}
			if item.MemberPrice != nil {
				itemDomain.MemberPrice = item.MemberPrice
			}

			// 关联商品信息
			if item.Edges.Product != nil {
				itemDomain.Product = convertProductToDomain(item.Edges.Product)
			}

			m.Items = append(m.Items, itemDomain)
		}
	}

	// 转换门店ID列表
	if len(em.Edges.Stores) > 0 {
		m.Stores = lo.Map(em.Edges.Stores, func(store *ent.Store, _ int) *domain.StoreSimple {
			return &domain.StoreSimple{
				ID:        store.ID,
				StoreName: store.StoreName,
			}
		})
	}

	return m
}
