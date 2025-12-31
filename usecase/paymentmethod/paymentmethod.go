package paymentmethod

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
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
		// 创建菜单
		return ds.PaymentMethodRepo().Create(ctx, p)
	})
}ß
