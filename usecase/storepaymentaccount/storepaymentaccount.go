package storepaymentaccount

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StorePaymentAccountInteractor = (*StorePaymentAccountInteractor)(nil)

type StorePaymentAccountInteractor struct {
	DS domain.DataStore
}

func NewStorePaymentAccountInteractor(ds domain.DataStore) *StorePaymentAccountInteractor {
	return &StorePaymentAccountInteractor{
		DS: ds,
	}
}

func (i *StorePaymentAccountInteractor) Create(ctx context.Context, account *domain.StorePaymentAccount, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StorePaymentAccountInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证门店存在且属于当前品牌商
		store, err := ds.StoreRepo().FindByID(ctx, account.StoreID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStorePaymentAccountStoreNotBelongMerchant)
			}
			return err
		}
		if store.MerchantID != account.MerchantID {
			return domain.ParamsError(domain.ErrStorePaymentAccountStoreNotBelongMerchant)
		}

		// 2. 验证品牌商收款账户存在且属于当前品牌商
		paymentAccount, err := ds.PaymentAccountRepo().FindByID(ctx, account.PaymentAccountID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStorePaymentAccountPaymentAccountNotBelongMerchant)
			}
			return err
		}
		if paymentAccount.MerchantID != account.MerchantID {
			return domain.ParamsError(domain.ErrStorePaymentAccountPaymentAccountNotBelongMerchant)
		}

		// 3. 检查每个门店的品牌商收款账户ID只能对应一个门店的收款账户
		exists, err := ds.StorePaymentAccountRepo().Exists(ctx, domain.StorePaymentAccountExistsParams{
			StoreID:          account.StoreID,
			PaymentAccountID: account.PaymentAccountID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrStorePaymentAccountPaymentAccountExist)
		}
		// 5. 创建门店收款账户
		return ds.StorePaymentAccountRepo().Create(ctx, account)
	})
}

func (i *StorePaymentAccountInteractor) Update(ctx context.Context, account *domain.StorePaymentAccount, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StorePaymentAccountInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证门店收款账户存在
		existing, err := ds.StorePaymentAccountRepo().FindByID(ctx, account.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStorePaymentAccountNotExists)
			}
			return err
		}

		// 2. 验证门店收款账户是否属于当前品牌商（通过门店验证）
		store, err := ds.StoreRepo().FindByID(ctx, existing.StoreID)
		if err != nil {
			return err
		}
		if store.MerchantID != account.MerchantID {
			return domain.ParamsError(domain.ErrStorePaymentAccountNotExists)
		}

		existing.MerchantNumber = account.MerchantNumber
		// 3. 更新门店收款账户
		return ds.StorePaymentAccountRepo().Update(ctx, existing)
	})
}

func (i *StorePaymentAccountInteractor) Delete(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StorePaymentAccountInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证门店收款账户存在
		existing, err := ds.StorePaymentAccountRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrStorePaymentAccountNotExists)
			}
			return err
		}

		// 2. 验证门店收款账户是否属于当前品牌商（通过门店验证）
		store, err := ds.StoreRepo().FindByID(ctx, existing.StoreID)
		if err != nil {
			return err
		}
		if store.MerchantID != user.GetMerchantID() {
			return domain.ParamsError(domain.ErrStorePaymentAccountNotExists)
		}

		// 3. 删除门店收款账户
		return ds.StorePaymentAccountRepo().Delete(ctx, id)
	})
}

func (i *StorePaymentAccountInteractor) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.StorePaymentAccountSearchParams,
	user domain.User,
) (res *domain.StorePaymentAccountSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StorePaymentAccountInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return i.DS.StorePaymentAccountRepo().PagedListBySearch(ctx, page, params)
}
