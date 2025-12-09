package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"gitlab.jiguang.dev/pos-dine/dine/ent/storeaccount"
	"gitlab.jiguang.dev/pos-dine/dine/ent/storeaccounttransaction"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StoreAccountRepository = (*StoreAccountRepository)(nil)

type StoreAccountRepository struct {
	Client *ent.Client
}

func NewStoreAccountRepository(client *ent.Client) *StoreAccountRepository {
	return &StoreAccountRepository{
		Client: client,
	}
}

func (s *StoreAccountRepository) Create(ctx context.Context, storeAccount *domain.StoreAccount) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreAccountRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	_, err = s.Client.StoreAccount.Create().
		SetStoreID(storeAccount.StoreID).
		SetBalance(storeAccount.Balance).
		SetPendingWithdraw(storeAccount.PendingWithdraw).
		SetTotalAmount(storeAccount.TotalAmount).
		SetWithdrawn(storeAccount.Withdrawn).
		Save(ctx)
	return
}

func (s *StoreAccountRepository) FindByStoreForUpdate(ctx context.Context, storeID int) (res *domain.StoreAccount, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreAccountRepository.FindByStoreForUpdate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	account, err := s.Client.StoreAccount.Query().
		Where(storeaccount.StoreID(storeID)).
		ForUpdate().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrStoreAccountNotExists)
		}
		return nil, err
	}
	return convertToStoreAccount(account), nil
}

func (s *StoreAccountRepository) FindByStore(ctx context.Context, storeID int) (res *domain.StoreAccount, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreAccountRepository.FindByStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	account, err := s.Client.StoreAccount.Query().
		Where(storeaccount.StoreID(storeID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrStoreAccountNotExists)
		}
		return nil, err
	}
	return convertToStoreAccount(account), nil
}

func convertToStoreAccount(account *ent.StoreAccount) *domain.StoreAccount {
	return &domain.StoreAccount{
		ID:              account.ID,
		StoreID:         account.StoreID,
		Balance:         account.Balance,
		PendingWithdraw: account.PendingWithdraw,
		Withdrawn:       account.Withdrawn,
		TotalAmount:     account.TotalAmount,
		CreatedAt:       account.CreatedAt,
		UpdatedAt:       account.UpdatedAt,
	}
}

func (s *StoreAccountRepository) AdjustAmount(ctx context.Context, storeID int, adj domain.StoreAccountAdjustments) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreAccountRepository.AdjustAmount")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = s.Client.StoreAccount.Update().
		Where(storeaccount.StoreID(storeID)).
		Modify(func(u *sql.UpdateBuilder) {
			u.Set(storeaccount.FieldBalance, sql.ExprFunc(func(b *sql.Builder) {
				b.Ident(storeaccount.FieldBalance).WriteOp(sql.OpAdd).Arg(adj.BalanceDelta)
			}))
			u.Set(storeaccount.FieldPendingWithdraw, sql.ExprFunc(func(b *sql.Builder) {
				b.Ident(storeaccount.FieldPendingWithdraw).WriteOp(sql.OpAdd).Arg(adj.PendingWithdrawDelta)
			}))
			u.Set(storeaccount.FieldWithdrawn, sql.ExprFunc(func(b *sql.Builder) {
				b.Ident(storeaccount.FieldWithdrawn).WriteOp(sql.OpAdd).Arg(adj.WithdrawnDelta)
			}))
			u.Set(storeaccount.FieldTotalAmount, sql.ExprFunc(func(b *sql.Builder) {
				b.Ident(storeaccount.FieldTotalAmount).WriteOp(sql.OpAdd).Arg(adj.TotalDelta)
			}))
		}).Exec(ctx)

	if err != nil {
		return fmt.Errorf("更新账户金额失败：%w", err)
	}
	return nil
}

func (s *StoreAccountRepository) RecordTransaction(ctx context.Context, tx *domain.StoreAccountTransaction) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreAccountRepository.RecordTransaction")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = s.Client.StoreAccountTransaction.Create().
		SetStoreID(tx.StoreID).
		SetAmount(tx.Amount).
		SetType(tx.Type).
		SetNo(tx.No).
		SetAfter(tx.After).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("记录账户流水失败：%w", err)
	}
	return nil
}

func (s *StoreAccountRepository) PagedListTransactions(ctx context.Context,
	page *upagination.Pagination, params domain.StoreAccountTransactionSearchParams,
) (res *domain.StoreAccountTransactionSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreAccountRepository.PagedListTransactions")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := s.Client.StoreAccountTransaction.Query()

	if params.StoreID > 0 {
		query.Where(storeaccounttransaction.StoreID(params.StoreID))
	}

	if params.StartAt != nil {
		query.Where(storeaccounttransaction.CreatedAtGTE(util.DayStart(*params.StartAt)))
	}
	if params.EndAt != nil {
		query.Where(storeaccounttransaction.CreatedAtLTE(util.DayEnd(*params.EndAt)))
	}

	if params.Type > 0 {
		query.Where(storeaccounttransaction.Type(params.Type))
	}

	// 计算总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	// 查询数据
	transactions, err := query.
		Order(ent.Desc(storeaccounttransaction.FieldID)).
		Offset(page.Offset()).
		Limit(page.Size).
		All(ctx)

	if err != nil {
		return nil, err
	}

	// 转换数据
	items := make(domain.StoreAccountTransactions, 0, len(transactions))
	for _, t := range transactions {
		items = append(items, convertToStoreAccountTransaction(t))
	}

	return &domain.StoreAccountTransactionSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

func convertToStoreAccountTransaction(t *ent.StoreAccountTransaction) *domain.StoreAccountTransaction {
	return &domain.StoreAccountTransaction{
		ID:        t.ID,
		StoreID:   t.StoreID,
		No:        t.No,
		Amount:    t.Amount,
		After:     t.After,
		Type:      t.Type,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
