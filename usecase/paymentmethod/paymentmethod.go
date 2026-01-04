package paymentmethod

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.PaymentMethodInteractor = (*PaymentMethodInteractor)(nil)

type PaymentMethodInteractor struct {
	DS domain.DataStore
}

func NewPaymentMethodInteractor(ds domain.DataStore) *PaymentMethodInteractor {
	return &PaymentMethodInteractor{
		DS: ds,
	}
}

func (i *PaymentMethodInteractor) Create(ctx context.Context, p *domain.PaymentMethod) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PaymentMethodInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		return ds.PaymentMethodRepo().Create(ctx, p)
	})
}

func (i *PaymentMethodInteractor) Update(ctx context.Context, p *domain.PaymentMethod, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PaymentMethodInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证结算方式存在
		paymentMethod, err := ds.PaymentMethodRepo().FindByID(ctx, p.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrPaymentMethodNotExists)
			}
			return err
		}
		if paymentMethod.MerchantID != user.GetMerchantID() {
			return domain.ParamsError(domain.ErrPaymentMethodNotExists)
		}
		return ds.PaymentMethodRepo().Update(ctx, p)
	})
}

func (i *PaymentMethodInteractor) Delete(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PaymentMethodInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证结算方式存在
		paymentMethod, err := ds.PaymentMethodRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrPaymentMethodNotExists)
			}
			return err
		}
		if paymentMethod.MerchantID != user.GetMerchantID() {
			return domain.ParamsError(domain.ErrPaymentMethodNotExists)
		}
		return ds.PaymentMethodRepo().Delete(ctx, id)
	})
}

func (i *PaymentMethodInteractor) GetDetail(ctx context.Context, id uuid.UUID, user domain.User) (res *domain.PaymentMethod, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PaymentMethodInteractor.GetDetail")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	paymentMethod, err := i.DS.PaymentMethodRepo().GetDetail(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ParamsError(domain.ErrPaymentMethodNotExists)
		}
		return nil, err
	}
	if paymentMethod.MerchantID != user.GetMerchantID() {
		return nil, domain.ParamsError(domain.ErrPaymentMethodNotExists)
	}
	return paymentMethod, nil
}

func (i *PaymentMethodInteractor) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.PaymentMethodSearchParams,
) (res *domain.PaymentMethodSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PaymentMethodInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.PaymentMethodRepo().PagedListBySearch(ctx, page, params)
}
