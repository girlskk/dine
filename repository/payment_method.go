package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/paymentmethod"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.PaymentMethodRepository = (*PaymentMethodRepository)(nil)

type PaymentMethodRepository struct {
	Client *ent.Client
}

func NewPaymentMethodRepository(client *ent.Client) *PaymentMethodRepository {
	return &PaymentMethodRepository{
		Client: client,
	}
}

func (repo *PaymentMethodRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.PaymentMethod, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentMethodRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	pm, err := repo.Client.PaymentMethod.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrMenuNotExists)
		}
		return nil, err
	}

	res = convertPaymentMethodToDomain(pm)
	return res, nil
}

func (repo *PaymentMethodRepository) GetDetail(ctx context.Context, id uuid.UUID) (res *domain.PaymentMethod, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentMethodRepository.GetDetail")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return repo.FindByID(ctx, id)
}

func (repo *PaymentMethodRepository) Create(ctx context.Context, p *domain.PaymentMethod) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentMethodRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	builder := repo.Client.PaymentMethod.Create().
		SetID(p.ID).
		SetMerchantID(p.MerchantID).
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

func (repo *PaymentMethodRepository) Update(ctx context.Context, p *domain.PaymentMethod) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentMethodRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.PaymentMethod.UpdateOneID(p.ID).
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
	_, err = builder.Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PaymentMethodRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentMethodRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.PaymentMethod.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PaymentMethodRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.PaymentMethodSearchParams,
) (res *domain.PaymentMethodSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PaymentMethodRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.PaymentMethod.Query()

	if params.MerchantID != uuid.Nil {
		query.Where(paymentmethod.MerchantID(params.MerchantID))
	}
	if params.Name != "" {
		query.Where(paymentmethod.NameContains(params.Name))
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
	list, err := query.Order(ent.Desc(paymentmethod.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.PaymentMethods, 0, len(list))
	for _, m := range list {
		items = append(items, convertPaymentMethodToDomain(m))
	}

	page.SetTotal(total)

	return &domain.PaymentMethodSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

// ============================================
// 转换函数
// ============================================

func convertPaymentMethodToDomain(pm *ent.PaymentMethod) *domain.PaymentMethod {
	if pm == nil {
		return nil
	}
	m := &domain.PaymentMethod{
		ID:               pm.ID,
		MerchantID:       pm.MerchantID,
		Name:             pm.Name,
		AccountingRule:   pm.AccountingRule,
		PaymentType:      pm.PaymentType,
		InvoiceRule:      pm.InvoiceRule,
		CashDrawerStatus: pm.CashDrawerStatus,
		DisplayChannels:  pm.DisplayChannels,
		Status:           pm.Status,
		CreatedAt:        pm.CreatedAt,
	}
	return m
}
