package repository

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)
ßß

var _ domain.PaymentMethodRepository = (*PaymentMethodRepository)(nil)

type PaymentMethodRepository struct {
	Client *ent.Client
}

func NewPaymentMethodRepository(client *ent.Client) *PaymentMethodRepository {
	return &PaymentMethodRepository{
		Client: client,
	}
}

//func (repo *PaymentMethodRepository) GetDetail(ctx context.Context, id uuid.UUID) (res *domain.Menu, err error) {
//	span, ctx := util.StartSpan(ctx, "repository", "PaymentMethodRepository.GetDetail")
//	defer func() {
//		util.SpanErrFinish(span, err)
//	}()
//
//	em, err := repo.Client.PaymentMethod.Query().
//		Where(paymentmethod.IDEQ(id)).
//		WithItems(func(query *ent.MenuItemQuery) {
//			query.WithProduct()
//		}).
//		WithStores().
//		Only(ctx)
//	if err != nil {
//		if ent.IsNotFound(err) {
//			return nil, domain.NotFoundError(domain.ErrMenuNotExists)
//		}
//		return nil, err
//	}
//
//	res = convertMenuToDomain(em)
//	return res, nil
//}

func (repo *PaymentMethodRepository) Create(ctx context.Context, p *domain.PaymentMethod) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentMethodRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	builder := repo.Client.PaymentMethod.Create().
		SetID(p.ID).
		SetName(p.Name).
		SetAccountingRule(p.AccountingRule).
		SetPaymentType(p.PaymentType).
		SetInvoiceRule(p.InvoiceRule).
		SetCashDrawerStatus(p.CashDrawerStatus).
		SetDisplayChannels(p.DisplayChannels).
		SetStatus(p.Status)
	if p.FeeRate != nil {
		builder = builder.SetFeeRate(*p.FeeRate)
	}
	created, err := builder.Save(ctx)
	if err != nil {
		return err
	}
	p.ID = created.ID
	p.CreatedAt = created.CreatedAt
	p.UpdatedAt = created.UpdatedAt
	return nil
}

//func (repo *PaymentMethodRepository) Update(ctx context.Context, m *domain.Menu) (err error) {
//	span, ctx := util.StartSpan(ctx, "repository", "PaymentMethodRepository.Update")
//	defer func() {
//		util.SpanErrFinish(span, err)
//	}()
//
//	// 先删除旧的菜单项记录（使用 DELETE，而不是 SET NULL）
//	skipSoftDeleteCtx := schematype.SkipSoftDelete(ctx)
//	_, err = repo.Client.PaymentMethodItem.Delete().
//		Where(menuitem.MenuIDEQ(m.ID)).
//		Exec(skipSoftDeleteCtx)
//	if err != nil {
//		return err
//	}
//
//	// 更新菜单基本信息
//	builder := repo.Client.PaymentMethod.UpdateOneID(m.ID).
//		SetName(m.Name).
//		SetDistributionRule(m.DistributionRule).
//		SetStoreCount(m.StoreCount).
//		SetItemCount(m.ItemCount)
//
//	// 更新关联门店（Many2Many）
//	if len(m.Stores) > 0 {
//		storeIDs := make([]uuid.UUID, 0, len(m.Stores))
//		for _, store := range m.Stores {
//			storeIDs = append(storeIDs, store.ID)
//		}
//		builder = builder.AddStoreIDs(storeIDs...)
//	} else {
//		builder = builder.ClearStores()
//	}
//
//	_, err = builder.Save(ctx)
//	if err != nil {
//		return err
//	}
//
//	// 更新菜单项
//	if len(m.Items) > 0 {
//		itemBuilders := make([]*ent.MenuItemCreate, 0, len(m.Items))
//		for _, item := range m.Items {
//			builder := repo.Client.PaymentMethodItem.Create().
//				SetID(item.ID).
//				SetMenuID(m.ID).
//				SetProductID(item.ProductID).
//				SetSaleRule(item.SaleRule)
//
//			if item.BasePrice != nil {
//				builder.SetBasePrice(*item.BasePrice)
//			}
//			if item.MemberPrice != nil {
//				builder.SetMemberPrice(*item.MemberPrice)
//			}
//			itemBuilders = append(itemBuilders, builder)
//		}
//		_, err = repo.Client.PaymentMethodItem.CreateBulk(itemBuilders...).Save(ctx)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func (repo *PaymentMethodRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
//	span, ctx := util.StartSpan(ctx, "repository", "PaymentMethodRepository.Delete")
//	defer func() {
//		util.SpanErrFinish(span, err)
//	}()
//
//	// 删除菜单项
//	_, err = repo.Client.PaymentMethodItem.Delete().Where(menuitem.MenuID(id)).Exec(ctx)
//	if err != nil {
//		return err
//	}
//
//	// 删除菜单
//	err = repo.Client.PaymentMethod.DeleteOneID(id).Exec(ctx)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (repo *PaymentMethodRepository) PagedListBySearch(
//	ctx context.Context,
//	page *upagination.Pagination,
//	params domain.MenuSearchParams,
//) (res *domain.MenuSearchRes, err error) {
//	span, ctx := util.StartSpan(ctx, "repository", "PaymentMethodRepository.PagedListBySearch")
//	defer func() {
//		util.SpanErrFinish(span, err)
//	}()
//
//	query := repo.Client.PaymentMethod.Query()
//
//	if params.MerchantID != uuid.Nil {
//		query.Where(menu.MerchantID(params.MerchantID))
//	}
//
//	if params.StoreID != uuid.Nil {
//		query.Where(menu.HasStoresWith(store.IDEQ(params.StoreID)))
//	}
//
//	if params.Name != "" {
//		query.Where(menu.NameContains(params.Name))
//	}
//
//	// 获取总数
//	total, err := query.Clone().Count(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	// 分页处理
//	query = query.
//		Offset(page.Offset()).
//		Limit(page.Size)
//
//	// 按创建时间倒序排列
//	entMenus, err := query.Order(ent.Desc(menu.FieldCreatedAt)).All(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	items := make(domain.Menus, 0, len(entMenus))
//	for _, m := range entMenus {
//		items = append(items, convertMenuToDomain(m))
//	}
//
//	page.SetTotal(total)
//
//	return &domain.MenuSearchRes{
//		Pagination: page,
//		Items:      items,
//	}, nil
//}
