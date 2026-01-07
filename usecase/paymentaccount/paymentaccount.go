package paymentaccount

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.PaymentAccountInteractor = (*PaymentAccountInteractor)(nil)

type PaymentAccountInteractor struct {
	DS domain.DataStore
}

func NewPaymentAccountInteractor(ds domain.DataStore) *PaymentAccountInteractor {
	return &PaymentAccountInteractor{
		DS: ds,
	}
}

func (i *PaymentAccountInteractor) Create(ctx context.Context, account *domain.PaymentAccount) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PaymentAccountInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 检查支付商户号是否已存在
		exists, err := ds.PaymentAccountRepo().Exists(ctx, domain.PaymentAccountExistsParams{
			MerchantID:     account.MerchantID,
			MerchantNumber: account.MerchantNumber,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrPaymentAccountMerchantNumberExist)
		}
		return ds.PaymentAccountRepo().Create(ctx, account)
	})
}

func (i *PaymentAccountInteractor) Update(ctx context.Context, account *domain.PaymentAccount, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PaymentAccountInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证收款账户存在
		existing, err := ds.PaymentAccountRepo().FindByID(ctx, account.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrPaymentAccountNotExists)
			}
			return err
		}

		// 验证收款账户是否属于当前品牌商
		if existing.MerchantID != user.GetMerchantID() {
			return domain.ParamsError(domain.ErrPaymentAccountNotExists)
		}

		// 检查支付商户号是否已存在（排除自身）
		exists, err := ds.PaymentAccountRepo().Exists(ctx, domain.PaymentAccountExistsParams{
			MerchantID:     existing.MerchantID,
			MerchantNumber: account.MerchantNumber,
			ExcludeID:      account.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrPaymentAccountMerchantNumberExist)
		}

		// 更新收款账户
		return ds.PaymentAccountRepo().Update(ctx, account)
	})
}

func (i *PaymentAccountInteractor) Delete(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PaymentAccountInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证收款账户存在
		existing, err := ds.PaymentAccountRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrPaymentAccountNotExists)
			}
			return err
		}

		// 验证收款账户是否属于当前品牌商
		if existing.MerchantID != user.GetMerchantID() {
			return domain.ParamsError(domain.ErrPaymentAccountNotExists)
		}

		// 检查是否有绑定的门店收款账户
		storeAccountCount, err := ds.PaymentAccountRepo().CountStoreAccounts(ctx, id)
		if err != nil {
			return err
		}
		if storeAccountCount > 0 {
			return domain.ParamsError(domain.ErrPaymentAccountHasStoreAccounts)
		}

		// 删除收款账户
		return ds.PaymentAccountRepo().Delete(ctx, id)
	})
}

func (i *PaymentAccountInteractor) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.PaymentAccountSearchParams,
) (res *domain.PaymentAccountSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PaymentAccountInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return i.DS.PaymentAccountRepo().PagedListBySearch(ctx, page, params)
}
