package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StoreWithdrawInteractor = (*StoreWithdrawInteractor)(nil)

type StoreWithdrawInteractor struct {
	ds            domain.DataStore
	dailySequence domain.DailySequence
}

func NewStoreWithdrawInteractor(dataStore domain.DataStore, dailySequence domain.DailySequence) *StoreWithdrawInteractor {
	return &StoreWithdrawInteractor{
		ds:            dataStore,
		dailySequence: dailySequence,
	}
}

func (i *StoreWithdrawInteractor) Apply(ctx context.Context, withdraw *domain.StoreWithdraw) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawInteractor.Apply")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	seqNo, err := i.genSeqNo(ctx, withdraw.StoreID)
	if err != nil {
		return err
	}
	withdraw.No = seqNo

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 获取门店账户
		account, err := ds.StoreAccountRepo().FindByStore(ctx, withdraw.StoreID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStoreAccountNotExists)
			}
			return err
		}
		// 2. 检查余额是否足够
		if account.Balance.LessThan(withdraw.Amount) {
			return domain.ParamsError(domain.ErrStoreWithdrawAmountNotEnough)
		}
		// 3. 创建提现记录
		if err = ds.StoreWithdrawRepo().Create(ctx, withdraw); err != nil {
			return err
		}
		return nil
	})
}

func (i *StoreWithdrawInteractor) Update(ctx context.Context, update *domain.StoreWithdraw) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 获取提现单并加锁
		withdraw, err := ds.StoreWithdrawRepo().FindByIDForUpdate(ctx, update.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStoreWithdrawNotExists)
			}
			return err
		}

		// 2. 检查是否为该门店的提现单
		if withdraw.StoreID != update.StoreID {
			return domain.ParamsError(domain.ErrStoreWithdrawNotExists)
		}

		// 3. 检查状态是否为待提交
		if withdraw.Status != domain.StoreWithdrawStatusUncommitted {
			return domain.ParamsError(domain.ErrStoreWithdrawStatusInvalid)
		}

		// 4. 获取门店账户
		account, err := ds.StoreAccountRepo().FindByStore(ctx, withdraw.StoreID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStoreAccountNotExists)
			}
			return err
		}
		// 5. 检查余额是否足够
		if account.Balance.LessThan(update.Amount) {
			return domain.ParamsError(domain.ErrStoreWithdrawAmountNotEnough)
		}

		// 6. 更新提现记录
		if err = ds.StoreWithdrawRepo().Update(ctx, update); err != nil {
			return err
		}
		return nil
	})
}

func (i *StoreWithdrawInteractor) Commit(ctx context.Context, id int, storeID int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawInteractor.Commit")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 获取提现单并加锁
		withdraw, err := ds.StoreWithdrawRepo().FindByIDForUpdate(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStoreWithdrawNotExists)
			}
			return err
		}
		// 检查提现单状态是否为待提交
		if withdraw.Status != domain.StoreWithdrawStatusUncommitted {
			return domain.ParamsError(domain.ErrStoreWithdrawStatusInvalid)
		}
		// 检查提现单门店是否为当前门店
		if withdraw.StoreID != storeID {
			return domain.ParamsError(domain.ErrStoreWithdrawNotExists)
		}
		// 获取门店账户并加锁
		account, err := ds.StoreAccountRepo().FindByStoreForUpdate(ctx, withdraw.StoreID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStoreAccountNotExists)
			}
			return err
		}
		// 检查余额是否足够
		if account.Balance.LessThan(withdraw.Amount) {
			return domain.ParamsError(domain.ErrStoreWithdrawAmountNotEnough)
		}

		//  更新账户金额：减少可用余额，增加待提现金额
		if err = ds.StoreAccountRepo().AdjustAmount(ctx, withdraw.StoreID, domain.StoreAccountAdjustments{
			BalanceDelta:         withdraw.Amount.Neg(),
			PendingWithdrawDelta: withdraw.Amount,
		}); err != nil {
			return err
		}

		// 更新提现单状态为"待审核"
		if err = ds.StoreWithdrawRepo().UpdateStatus(ctx, id, domain.StoreWithdrawStatusPending); err != nil {
			return err
		}
		// 记录账户金额变更日志
		transaction := &domain.StoreAccountTransaction{
			StoreID: withdraw.StoreID,
			No:      withdraw.No,
			Amount:  withdraw.Amount.Neg(),
			After:   account.Balance.Sub(withdraw.Amount),
			Type:    domain.TransactionTypeWithdrawApply,
		}
		if err = ds.StoreAccountRepo().RecordTransaction(ctx, transaction); err != nil {
			return err
		}
		return nil
	})
}

func (i *StoreWithdrawInteractor) genSeqNo(ctx context.Context, storeID int) (string, error) {
	no, err := i.dailySequence.Next(ctx, domain.DailySequencePrefixStoreWithdrawNo)
	if err != nil {
		return "", fmt.Errorf("failed to generate seq no: %w", err)
	}
	datePart := time.Now().Format("060102")
	return fmt.Sprintf("%s%s%03d%03d", domain.StoreWithdrawNoPrefix, datePart, storeID, no), nil
}

func (i *StoreWithdrawInteractor) PagedListBySearch(ctx context.Context, page *upagination.Pagination,
	params domain.StoreWithdrawSearchParams,
) (res *domain.StoreWithdrawSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return i.ds.StoreWithdrawRepo().PagedListBySearch(ctx, page, params)
}

func (i *StoreWithdrawInteractor) Approve(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawInteractor.Approve")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 获取提现单并加锁
		withdraw, err := ds.StoreWithdrawRepo().FindByIDForUpdate(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStoreWithdrawNotExists)
			}
			return err
		}

		// 2. 检查状态是否为待审核
		if withdraw.Status != domain.StoreWithdrawStatusPending {
			return domain.ParamsError(domain.ErrStoreWithdrawStatusInvalid)
		}

		// 3. 获取门店账户并加锁
		account, err := ds.StoreAccountRepo().FindByStoreForUpdate(ctx, withdraw.StoreID)
		if err != nil {
			return err
		}

		// 4. 更新提现单状态为已审核
		if err = ds.StoreWithdrawRepo().UpdateStatus(ctx, id, domain.StoreWithdrawStatusApproved); err != nil {
			return err
		}

		// 5. 更新账户：-待提现金额 + 已提现金额
		if err = ds.StoreAccountRepo().AdjustAmount(ctx, withdraw.StoreID, domain.StoreAccountAdjustments{
			PendingWithdrawDelta: withdraw.Amount.Neg(),
			WithdrawnDelta:       withdraw.Amount,
		}); err != nil {
			return err
		}

		// 6. 记录账户金额变更日志
		transaction := &domain.StoreAccountTransaction{
			StoreID: account.StoreID,
			No:      withdraw.No,
			Amount:  decimal.Zero, // 余额不变，只是从待提现转为已提现
			After:   account.Balance,
			Type:    domain.TransactionTypeWithdrawApprove,
		}
		if err = ds.StoreAccountRepo().RecordTransaction(ctx, transaction); err != nil {
			return err
		}

		return nil
	})
}

func (i *StoreWithdrawInteractor) Reject(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawInteractor.Reject")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 获取提现单并加锁
		withdraw, err := ds.StoreWithdrawRepo().FindByIDForUpdate(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStoreWithdrawNotExists)
			}
			return err
		}

		// 2. 检查状态是否为待审核
		if withdraw.Status != domain.StoreWithdrawStatusPending {
			return domain.ParamsError(domain.ErrStoreWithdrawStatusInvalid)
		}

		// 3. 获取门店账户并加锁
		account, err := ds.StoreAccountRepo().FindByStoreForUpdate(ctx, withdraw.StoreID)
		if err != nil {
			return err
		}

		// 4. 更新提现单状态为已驳回
		if err = ds.StoreWithdrawRepo().UpdateStatus(ctx, id, domain.StoreWithdrawStatusRejected); err != nil {
			return err
		}

		// 5. 更新账户：-待提现金额 +可提现金额
		if err = ds.StoreAccountRepo().AdjustAmount(ctx, withdraw.StoreID, domain.StoreAccountAdjustments{
			PendingWithdrawDelta: withdraw.Amount.Neg(),
			BalanceDelta:         withdraw.Amount,
		}); err != nil {
			return err
		}

		// 6. 记录账户金额变更日志
		transaction := &domain.StoreAccountTransaction{
			StoreID: account.StoreID,
			No:      withdraw.No,
			Amount:  withdraw.Amount,
			After:   account.Balance.Add(withdraw.Amount),
			Type:    domain.TransactionTypeWithdrawReject,
		}
		if err = ds.StoreAccountRepo().RecordTransaction(ctx, transaction); err != nil {
			return err
		}

		return nil
	})
}

// Cancel 撤回提现单
func (i *StoreWithdrawInteractor) Cancel(ctx context.Context, id int, storeID int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawInteractor.Cancel")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 获取提现单并加锁
		withdraw, err := ds.StoreWithdrawRepo().FindByIDForUpdate(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStoreWithdrawNotExists)
			}
			return err
		}

		// 2. 检查是否为该门店的提现单
		if withdraw.StoreID != storeID {
			return domain.ParamsError(domain.ErrStoreWithdrawNotExists)
		}

		// 3. 检查状态是否为待审核
		if withdraw.Status != domain.StoreWithdrawStatusPending {
			return domain.ParamsError(domain.ErrStoreWithdrawStatusInvalid)
		}

		// 4. 获取门店账户并加锁
		account, err := ds.StoreAccountRepo().FindByStoreForUpdate(ctx, withdraw.StoreID)
		if err != nil {
			return err
		}

		// 5. 更新提现单状态为待提交
		if err = ds.StoreWithdrawRepo().UpdateStatus(ctx, id, domain.StoreWithdrawStatusUncommitted); err != nil {
			return err
		}

		// 6. 更新账户：减少待提现金额，增加可用余额
		if err = ds.StoreAccountRepo().AdjustAmount(ctx, withdraw.StoreID, domain.StoreAccountAdjustments{
			BalanceDelta:         withdraw.Amount,
			PendingWithdrawDelta: withdraw.Amount.Neg(),
		}); err != nil {
			return err
		}

		// 7. 记录账户金额变更日志
		transaction := &domain.StoreAccountTransaction{
			StoreID: account.StoreID,
			No:      withdraw.No,
			Amount:  withdraw.Amount,
			After:   account.Balance.Add(withdraw.Amount),
			Type:    domain.TransactionTypeWithdrawCancel,
		}
		if err = ds.StoreAccountRepo().RecordTransaction(ctx, transaction); err != nil {
			return err
		}

		return nil
	})
}

// Delete 删除提现单
func (i *StoreWithdrawInteractor) Delete(ctx context.Context, id int, storeID int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 获取提现单
		withdraw, err := ds.StoreWithdrawRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStoreWithdrawNotExists)
			}
			return err
		}
		// 2. 检查是否为该门店的提现单
		if withdraw.StoreID != storeID {
			return domain.ParamsError(domain.ErrStoreWithdrawNotExists)
		}

		// 3. 检查状态是否允许删除（待提交、已驳回的可以删除）
		if withdraw.Status != domain.StoreWithdrawStatusUncommitted &&
			withdraw.Status != domain.StoreWithdrawStatusRejected {
			return domain.ParamsError(domain.ErrStoreWithdrawStatusInvalid)
		}

		// 4. 删除提现单
		if err = ds.StoreWithdrawRepo().Delete(ctx, id); err != nil {
			return err
		}

		return nil
	})
}
