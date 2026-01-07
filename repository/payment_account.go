package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	entpaymentaccount "gitlab.jiguang.dev/pos-dine/dine/ent/paymentaccount"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.PaymentAccountRepository = (*PaymentAccountRepository)(nil)

type PaymentAccountRepository struct {
	Client *ent.Client
}

func NewPaymentAccountRepository(client *ent.Client) *PaymentAccountRepository {
	return &PaymentAccountRepository{
		Client: client,
	}
}

func (repo *PaymentAccountRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.PaymentAccount, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentAccountRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	epa, err := repo.Client.PaymentAccount.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrPaymentAccountNotExists)
		}
		return nil, err
	}

	return convertPaymentAccountToDomain(epa), nil
}

func (repo *PaymentAccountRepository) Create(ctx context.Context, account *domain.PaymentAccount) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentAccountRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	created, err := repo.Client.PaymentAccount.Create().
		SetID(account.ID).
		SetMerchantID(account.MerchantID).
		SetChannel(account.Channel).
		SetMerchantNumber(account.MerchantNumber).
		SetMerchantName(account.MerchantName).
		SetIsDefault(account.IsDefault).
		Save(ctx)
	if err != nil {
		return err
	}

	account.CreatedAt = created.CreatedAt
	account.UpdatedAt = created.UpdatedAt
	return nil
}

func (repo *PaymentAccountRepository) Update(ctx context.Context, account *domain.PaymentAccount) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentAccountRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	updated, err := repo.Client.PaymentAccount.UpdateOneID(account.ID).
		SetChannel(account.Channel).
		SetMerchantNumber(account.MerchantNumber).
		SetMerchantName(account.MerchantName).
		SetIsDefault(account.IsDefault).
		Save(ctx)
	if err != nil {
		return err
	}

	account.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *PaymentAccountRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentAccountRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return repo.Client.PaymentAccount.DeleteOneID(id).Exec(ctx)
}

func (repo *PaymentAccountRepository) Exists(ctx context.Context, params domain.PaymentAccountExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentAccountRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.PaymentAccount.Query().
		Where(entpaymentaccount.MerchantID(params.MerchantID)).
		Where(entpaymentaccount.ChannelEQ(params.Channel))

	// 排除指定的ID（用于更新时检查唯一性）
	if params.ExcludeID != uuid.Nil {
		query.Where(entpaymentaccount.IDNEQ(params.ExcludeID))
	}

	return query.Exist(ctx)
}

func (repo *PaymentAccountRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.PaymentAccountSearchParams,
) (res *domain.PaymentAccountSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentAccountRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.PaymentAccount.Query()

	// 必填条件：品牌商ID
	query.Where(entpaymentaccount.MerchantID(params.MerchantID))

	// 可选条件：支付渠道
	if params.Channel != "" {
		query.Where(entpaymentaccount.ChannelEQ(params.Channel))
	}

	// 可选条件：支付商户名称（模糊匹配）
	if params.MerchantName != "" {
		query.Where(entpaymentaccount.MerchantNameContains(params.MerchantName))
	}

	// 可选条件：创建时间范围
	if params.CreatedAtStart != nil {
		query.Where(entpaymentaccount.CreatedAtGTE(*params.CreatedAtStart))
	}
	if params.CreatedAtEnd != nil {
		query.Where(entpaymentaccount.CreatedAtLTE(*params.CreatedAtEnd))
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
	entAccounts, err := query.Order(ent.Desc(entpaymentaccount.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.PaymentAccounts, 0, len(entAccounts))
	for _, a := range entAccounts {
		items = append(items, convertPaymentAccountToDomain(a))
	}

	page.SetTotal(total)

	return &domain.PaymentAccountSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

// FindForUpdateByMerchantID 锁定查询该品牌商下所有账户（用于并发控制）
func (repo *PaymentAccountRepository) FindForUpdateByMerchantID(ctx context.Context, merchantID uuid.UUID) (res domain.PaymentAccounts, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentAccountRepository.FindForUpdateByMerchantID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	entAccounts, err := repo.Client.PaymentAccount.Query().
		Where(entpaymentaccount.MerchantID(merchantID)).
		ForUpdate(). // 行级锁：锁定该品牌商下所有账户
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.PaymentAccounts, 0, len(entAccounts))
	for _, a := range entAccounts {
		items = append(items, convertPaymentAccountToDomain(a))
	}

	return items, nil
}

// UpdateAllDefaultStatus 批量更新该品牌商下所有账户的默认状态
func (repo *PaymentAccountRepository) UpdateAllDefaultStatus(ctx context.Context, merchantID uuid.UUID, isDefault bool) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentAccountRepository.UpdateAllDefaultStatus")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.PaymentAccount.Update().
		Where(entpaymentaccount.MerchantID(merchantID)).
		SetIsDefault(isDefault).
		Save(ctx)
	return err
}

func (repo *PaymentAccountRepository) CountStoreAccounts(ctx context.Context, id uuid.UUID) (count int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentAccountRepository.CountStoreAccounts")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// TODO: 当门店收款账户表建立后，实现此方法
	// 目前返回0，表示没有绑定的门店收款账户
	return 0, nil
}

// ============================================
// 转换函数
// ============================================

func convertPaymentAccountToDomain(epa *ent.PaymentAccount) *domain.PaymentAccount {
	if epa == nil {
		return nil
	}

	account := &domain.PaymentAccount{
		ID:             epa.ID,
		MerchantID:     epa.MerchantID,
		Channel:        epa.Channel,
		MerchantNumber: epa.MerchantNumber,
		MerchantName:   epa.MerchantName,
		IsDefault:      epa.IsDefault,
		CreatedAt:      epa.CreatedAt,
		UpdatedAt:      epa.UpdatedAt,
	}

	return account
}
