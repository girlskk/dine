package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	entpaymentaccount "gitlab.jiguang.dev/pos-dine/dine/ent/paymentaccount"
	entstore "gitlab.jiguang.dev/pos-dine/dine/ent/store"
	entstorepaymentaccount "gitlab.jiguang.dev/pos-dine/dine/ent/storepaymentaccount"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StorePaymentAccountRepository = (*StorePaymentAccountRepository)(nil)

type StorePaymentAccountRepository struct {
	Client *ent.Client
}

func NewStorePaymentAccountRepository(client *ent.Client) *StorePaymentAccountRepository {
	return &StorePaymentAccountRepository{
		Client: client,
	}
}

func (repo *StorePaymentAccountRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.StorePaymentAccount, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StorePaymentAccountRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	espa, err := repo.Client.StorePaymentAccount.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrStorePaymentAccountNotExists)
		}
		return nil, err
	}

	return convertStorePaymentAccountToDomain(espa), nil
}

func (repo *StorePaymentAccountRepository) Create(ctx context.Context, account *domain.StorePaymentAccount) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StorePaymentAccountRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	created, err := repo.Client.StorePaymentAccount.Create().
		SetID(account.ID).
		SetMerchantID(account.MerchantID).
		SetStoreID(account.StoreID).
		SetPaymentAccountID(account.PaymentAccountID).
		SetMerchantNumber(account.MerchantNumber).
		Save(ctx)
	if err != nil {
		return err
	}

	account.CreatedAt = created.CreatedAt
	account.UpdatedAt = created.UpdatedAt
	return nil
}

func (repo *StorePaymentAccountRepository) Update(ctx context.Context, account *domain.StorePaymentAccount) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StorePaymentAccountRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	updated, err := repo.Client.StorePaymentAccount.UpdateOneID(account.ID).
		SetMerchantNumber(account.MerchantNumber).
		Save(ctx)
	if err != nil {
		return err
	}

	account.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *StorePaymentAccountRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StorePaymentAccountRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return repo.Client.StorePaymentAccount.DeleteOneID(id).Exec(ctx)
}

func (repo *StorePaymentAccountRepository) Exists(ctx context.Context, params domain.StorePaymentAccountExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StorePaymentAccountRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.StorePaymentAccount.Query().
		Where(entstorepaymentaccount.StoreID(params.StoreID)).
		Where(entstorepaymentaccount.PaymentAccountID(params.PaymentAccountID))

	// 排除指定的ID（用于更新时检查唯一性）
	if params.ExcludeID != uuid.Nil {
		query.Where(entstorepaymentaccount.IDNEQ(params.ExcludeID))
	}

	return query.Exist(ctx)
}

func (repo *StorePaymentAccountRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.StorePaymentAccountSearchParams,
) (res *domain.StorePaymentAccountSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StorePaymentAccountRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.StorePaymentAccount.Query()

	// 必填条件：品牌商ID（通过门店关联查询）
	if params.MerchantID != uuid.Nil {
		query.Where(entstorepaymentaccount.HasStoreWith(
			entstore.MerchantID(params.MerchantID),
		))
	}

	// 可选条件：门店ID列表
	if len(params.StoreIDs) > 0 {
		query.Where(entstorepaymentaccount.StoreIDIn(params.StoreIDs...))
	}

	// 可选条件：品牌商支付商户名称（模糊匹配，通过关联的品牌商收款账户查询）
	if params.MerchantName != "" {
		query.Where(entstorepaymentaccount.HasPaymentAccountWith(
			entpaymentaccount.MerchantNameContains(params.MerchantName),
		))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}

	// 分页处理
	query = query.
		WithStore().
		WithPaymentAccount().
		Offset(page.Offset()).
		Limit(page.Size)

	// 按创建时间倒序排列
	entAccounts, err := query.Order(ent.Desc(entstorepaymentaccount.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.StorePaymentAccounts, 0, len(entAccounts))
	for _, a := range entAccounts {
		items = append(items, convertStorePaymentAccountToDomain(a))
	}

	page.SetTotal(total)

	return &domain.StorePaymentAccountSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

// ============================================
// 转换函数
// ============================================

func convertStorePaymentAccountToDomain(espa *ent.StorePaymentAccount) *domain.StorePaymentAccount {
	if espa == nil {
		return nil
	}

	account := &domain.StorePaymentAccount{
		ID:               espa.ID,
		MerchantID:       espa.MerchantID,
		StoreID:          espa.StoreID,
		PaymentAccountID: espa.PaymentAccountID,
		MerchantNumber:   espa.MerchantNumber,
		CreatedAt:        espa.CreatedAt,
		UpdatedAt:        espa.UpdatedAt,
	}

	if espa.Edges.Store != nil {
		account.Store = &domain.StoreSimple{
			ID:        espa.Edges.Store.ID,
			StoreName: espa.Edges.Store.StoreName,
		}
	}

	if espa.Edges.PaymentAccount != nil {
		account.PaymentAccount = &domain.PaymentAccount{
			ID:             espa.Edges.PaymentAccount.ID,
			MerchantID:     espa.Edges.PaymentAccount.MerchantID,
			Channel:        espa.Edges.PaymentAccount.Channel,
			MerchantNumber: espa.Edges.PaymentAccount.MerchantNumber,
			MerchantName:   espa.Edges.PaymentAccount.MerchantName,
			IsDefault:      espa.Edges.PaymentAccount.IsDefault,
			CreatedAt:      espa.Edges.PaymentAccount.CreatedAt,
			UpdatedAt:      espa.Edges.PaymentAccount.UpdatedAt,
		}
	}

	return account
}
