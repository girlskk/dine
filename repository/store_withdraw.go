package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/storewithdraw"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StoreWithdrawRepository = (*StoreWithdrawRepository)(nil)

type StoreWithdrawRepository struct {
	Client *ent.Client
}

func NewStoreWithdrawRepository(client *ent.Client) *StoreWithdrawRepository {
	return &StoreWithdrawRepository{
		Client: client,
	}
}

func (s *StoreWithdrawRepository) Create(ctx context.Context, withdraw *domain.StoreWithdraw) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	_, err = s.Client.StoreWithdraw.Create().
		SetStoreID(withdraw.StoreID).
		SetStoreName(withdraw.StoreName).
		SetNo(withdraw.No).
		SetAmount(withdraw.Amount).
		SetPointWithdrawalRate(withdraw.PointWithdrawalRate).
		SetActualAmount(withdraw.ActualAmount).
		SetAccountType(withdraw.AccountType).
		SetBankAccount(withdraw.BankAccount).
		SetBankCardName(withdraw.BankCardName).
		SetBankName(withdraw.BankName).
		SetBankBranch(withdraw.BankBranch).
		SetInvoiceAmount(withdraw.InvoiceAmount).
		SetStatus(withdraw.Status).
		Save(ctx)
	return
}

func (s *StoreWithdrawRepository) Update(ctx context.Context, withdraw *domain.StoreWithdraw) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	_, err = s.Client.StoreWithdraw.UpdateOneID(withdraw.ID).
		SetAmount(withdraw.Amount).
		SetPointWithdrawalRate(withdraw.PointWithdrawalRate).
		SetActualAmount(withdraw.ActualAmount).
		SetAccountType(withdraw.AccountType).
		SetBankAccount(withdraw.BankAccount).
		SetBankCardName(withdraw.BankCardName).
		SetBankName(withdraw.BankName).
		SetBankBranch(withdraw.BankBranch).
		SetInvoiceAmount(withdraw.InvoiceAmount).
		Save(ctx)
	return
}

func (s *StoreWithdrawRepository) PagedListBySearch(ctx context.Context, page *upagination.Pagination,
	params domain.StoreWithdrawSearchParams,
) (res *domain.StoreWithdrawSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := s.Client.StoreWithdraw.Query()

	if params.StoreID > 0 {
		query.Where(storewithdraw.StoreID(params.StoreID))
	}

	if params.Status > 0 {
		query.Where(storewithdraw.Status(params.Status))
	}

	if params.StartAt != nil {
		query.Where(storewithdraw.CreatedAtGTE(util.DayStart(*params.StartAt)))
	}
	if params.EndAt != nil {
		query.Where(storewithdraw.CreatedAtLTE(util.DayEnd(*params.EndAt)))
	}

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}

	page.Total = total

	withdraws, err := query.Order(ent.Desc(storewithdraw.FieldID)).
		Offset(page.Offset()).
		Limit(page.Size).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*domain.StoreWithdraw, 0, len(withdraws))
	for _, w := range withdraws {
		items = append(items, convertToStoreWithdraw(w))
	}

	return &domain.StoreWithdrawSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

// convertToStoreWithdraw 将ent对象转换为domain对象
func convertToStoreWithdraw(w *ent.StoreWithdraw) *domain.StoreWithdraw {
	return &domain.StoreWithdraw{
		ID:                  w.ID,
		StoreID:             w.StoreID,
		StoreName:           w.StoreName,
		No:                  w.No,
		Amount:              w.Amount,
		PointWithdrawalRate: w.PointWithdrawalRate,
		ActualAmount:        w.ActualAmount,
		AccountType:         w.AccountType,
		BankAccount:         w.BankAccount,
		BankCardName:        w.BankCardName,
		BankName:            w.BankName,
		BankBranch:          w.BankBranch,
		InvoiceAmount:       w.InvoiceAmount,
		Status:              w.Status,
		CreatedAt:           w.CreatedAt,
		UpdatedAt:           w.UpdatedAt,
	}
}

func (s *StoreWithdrawRepository) FindByIDForUpdate(ctx context.Context, id int) (res *domain.StoreWithdraw, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawRepository.FindByIDForUpdate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	withdraw, err := s.Client.StoreWithdraw.Query().
		Where(storewithdraw.ID(id)).
		ForUpdate().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrStoreWithdrawNotExists)
		}
		return nil, err
	}
	return convertToStoreWithdraw(withdraw), nil
}

func (s *StoreWithdrawRepository) FindByID(ctx context.Context, id int) (res *domain.StoreWithdraw, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	withdraw, err := s.Client.StoreWithdraw.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrStoreWithdrawNotExists)
		}
		return nil, err
	}
	return convertToStoreWithdraw(withdraw), nil
}

func (s *StoreWithdrawRepository) UpdateStatus(ctx context.Context, id int, status domain.StoreWithdrawStatus) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawRepository.UpdateStatus")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	_, err = s.Client.StoreWithdraw.UpdateOneID(id).
		SetStatus(status).
		Save(ctx)
	return
}

func (s *StoreWithdrawRepository) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	err = s.Client.StoreWithdraw.DeleteOneID(id).Exec(ctx)
	return
}
