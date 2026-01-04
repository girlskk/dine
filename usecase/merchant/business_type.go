package merchant

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.MerchantBusinessTypeInteractor = (*BusinessTypeInteractor)(nil)

type BusinessTypeInteractor struct {
	DS domain.DataStore
}

func (interactor *BusinessTypeInteractor) GetAll(ctx context.Context) (list []*domain.MerchantBusinessType, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BusinessTypeInteractor.GetAll")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return interactor.DS.MerchantBusinessTypeRepo().GetAll(ctx)
}

func NewMerchantBusinessTypeInteractor(ds domain.DataStore) *BusinessTypeInteractor {
	return &BusinessTypeInteractor{
		DS: ds,
	}
}
