package merchant

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/i18n"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.MerchantBusinessTypeInteractor = (*BusinessTypeInteractor)(nil)

type BusinessTypeInteractor struct {
	DS domain.DataStore
}

func (interactor *BusinessTypeInteractor) GetAll(ctx context.Context) (list []*domain.BusinessType, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BusinessTypeInteractor.GetAll")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	for _, e := range domain.BusinessTypeEntries {
		name := i18n.Translate(ctx, e.MsgID, nil)
		list = append(list, &domain.BusinessType{
			TypeCode: e.Code,
			TypeName: name,
		})
	}

	return
}

func NewMerchantBusinessTypeInteractor(ds domain.DataStore) *BusinessTypeInteractor {
	return &BusinessTypeInteractor{
		DS: ds,
	}
}
