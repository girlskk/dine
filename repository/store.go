package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/store"
	"gitlab.jiguang.dev/pos-dine/dine/ent/storefinance"
	"gitlab.jiguang.dev/pos-dine/dine/ent/storeinfo"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StoreRepository = (*StoreRepository)(nil)

type StoreRepository struct {
	Client *ent.Client
}

func NewStoreRepository(client *ent.Client) *StoreRepository {
	return &StoreRepository{
		Client: client,
	}
}

func (repo *StoreRepository) Find(ctx context.Context, id int) (res *domain.Store, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreRepository.Find")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	store, err := repo.Client.Store.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrStoreNotExists)
		}
		return nil, err
	}
	return convertStore(store), nil
}

func (repo *StoreRepository) ListAll(ctx context.Context) (res domain.Stores, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreRepository.ListAll")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	entStores, err := repo.Client.Store.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	stores := make(domain.Stores, len(entStores))
	for i, t := range entStores {
		stores[i] = convertStore(t)
	}
	return stores, nil
}

func (repo *StoreRepository) Exists(ctx context.Context, params domain.StoreExistsParams) (exists bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := repo.Client.Store.Query()

	if params.Name != "" {
		query.Where(store.Name(params.Name))
	}
	return query.Exist(ctx)
}

func (repo *StoreRepository) Create(ctx context.Context, store *domain.Store) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	create := repo.Client.Store.Create().
		SetName(store.Name).
		SetType(store.Type).
		SetCooperationType(store.CooperationType).
		SetNeedAudit(store.NeedAudit).
		SetEnabled(store.Enabled).
		SetPointSettlementRate(store.PointSettlementRate).
		SetPointWithdrawalRate(store.PointWithdrawalRate).
		SetHuifuID(store.HuifuID).
		SetZxhID(store.ZxhID).
		SetZxhSecret(store.ZxhSecret)

	created, err := create.Save(ctx)
	if err != nil {
		return err
	}
	store.ID = created.ID

	if store.Info != nil {
		info := repo.Client.StoreInfo.Create().
			SetCity(store.Info.City).
			SetAddress(store.Info.Address).
			SetContactName(store.Info.ContactName).
			SetContactPhone(store.Info.ContactPhone).
			SetImages(store.Info.Images).
			SetStoreID(store.ID)
		_, err = info.Save(ctx)
		if err != nil {
			return err
		}
	}

	if store.Finance != nil {
		finance := repo.Client.StoreFinance.Create().
			SetBankAccount(store.Finance.BankAccount).
			SetBankCardName(store.Finance.BankCardName).
			SetBankName(store.Finance.BankName).
			SetBranchName(store.Finance.BranchName).
			SetPublicAccount(store.Finance.PublicAccount).
			SetCompanyName(store.Finance.CompanyName).
			SetPublicBankName(store.Finance.PublicBankName).
			SetPublicBranchName(store.Finance.PublicBranchName).
			SetCreditCode(store.Finance.CreditCode).
			SetStoreID(store.ID)
		_, err = finance.Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *StoreRepository) Update(ctx context.Context, dstore *domain.Store) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	// 更新主表
	_, err = repo.Client.Store.Update().
		SetName(dstore.Name).
		SetType(dstore.Type).
		SetCooperationType(dstore.CooperationType).
		SetNeedAudit(dstore.NeedAudit).
		SetEnabled(dstore.Enabled).
		SetPointSettlementRate(dstore.PointSettlementRate).
		SetPointWithdrawalRate(dstore.PointWithdrawalRate).
		SetHuifuID(dstore.HuifuID).
		SetZxhID(dstore.ZxhID).
		SetZxhSecret(dstore.ZxhSecret).
		Where(store.ID(dstore.ID)).
		Save(ctx)

	if err != nil {
		return err
	}

	if dstore.Info != nil {
		_, err = repo.Client.StoreInfo.Update().
			SetCity(dstore.Info.City).
			SetAddress(dstore.Info.Address).
			SetContactName(dstore.Info.ContactName).
			SetContactPhone(dstore.Info.ContactPhone).
			SetImages(dstore.Info.Images).
			Where(storeinfo.StoreID(dstore.ID)).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	if dstore.Finance != nil {
		_, err = repo.Client.StoreFinance.Update().
			SetBankAccount(dstore.Finance.BankAccount).
			SetBankCardName(dstore.Finance.BankCardName).
			SetBankName(dstore.Finance.BankName).
			SetBranchName(dstore.Finance.BranchName).
			SetPublicAccount(dstore.Finance.PublicAccount).
			SetCompanyName(dstore.Finance.CompanyName).
			SetPublicBankName(dstore.Finance.PublicBankName).
			SetPublicBranchName(dstore.Finance.PublicBranchName).
			SetCreditCode(dstore.Finance.CreditCode).
			Where(storefinance.StoreID(dstore.ID)).
			Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *StoreRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.StoreSearchParams,
) (res *domain.StoreSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 构建基础查询
	query := repo.Client.Store.Query()

	// 应用搜索条件
	if params.Name != "" {
		query.Where(store.NameContains(params.Name))
	}

	if params.City != "" {
		query.Where(store.HasStoreInfoWith(storeinfo.CityContains(params.City)))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	// 应用排序分页
	query.Order(ent.Desc(store.FieldID)).
		Offset(page.Offset()).
		Limit(page.Size)

	// 执行查询
	entStores, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为领域对象
	stores := make(domain.Stores, len(entStores))
	for i, t := range entStores {
		stores[i] = convertStore(t)
	}
	return &domain.StoreSearchRes{
		Pagination: page,
		Items:      stores,
	}, nil
}

func (repo *StoreRepository) GetDetail(ctx context.Context, id int) (res *domain.Store, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreRepository.GetDetail")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	p, err := repo.Client.Store.Query().
		Where(store.ID(id)).
		WithStoreInfo().
		WithStoreFinance().
		WithBackendUser().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrStoreNameExists)
		}
		return nil, err
	}
	return convertStore(p), nil
}

func convertStore(s *ent.Store) *domain.Store {
	if s == nil {
		return nil
	}
	store := &domain.Store{
		ID:                  s.ID,
		Name:                s.Name,
		Type:                s.Type,
		CooperationType:     s.CooperationType,
		NeedAudit:           s.NeedAudit,
		Enabled:             s.Enabled,
		PointSettlementRate: s.PointSettlementRate,
		PointWithdrawalRate: s.PointWithdrawalRate,
		HuifuID:             s.HuifuID,
		ZxhID:               s.ZxhID,
		ZxhSecret:           s.ZxhSecret,
		CreatedAt:           s.CreatedAt,
		UpdatedAt:           s.UpdatedAt,
		BackendUser:         convertBackendUser(s.Edges.BackendUser),
	}

	if s.Edges.StoreInfo != nil {
		store.Info = &domain.StoreInfo{
			City:         s.Edges.StoreInfo.City,
			Address:      s.Edges.StoreInfo.Address,
			ContactName:  s.Edges.StoreInfo.ContactName,
			ContactPhone: s.Edges.StoreInfo.ContactPhone,
			Images:       s.Edges.StoreInfo.Images,
		}
	}

	if s.Edges.StoreFinance != nil {
		store.Finance = &domain.StoreFinance{
			BankAccount:      s.Edges.StoreFinance.BankAccount,
			BankCardName:     s.Edges.StoreFinance.BankCardName,
			BankName:         s.Edges.StoreFinance.BankName,
			BranchName:       s.Edges.StoreFinance.BranchName,
			PublicAccount:    s.Edges.StoreFinance.PublicAccount,
			CompanyName:      s.Edges.StoreFinance.CompanyName,
			PublicBankName:   s.Edges.StoreFinance.PublicBankName,
			PublicBranchName: s.Edges.StoreFinance.PublicBranchName,
			CreditCode:       s.Edges.StoreFinance.CreditCode,
		}
	}
	return store
}
