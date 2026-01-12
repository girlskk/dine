package businessconfig

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.BusinessConfigInteractor = (*BusinessConfigInteractor)(nil)

type BusinessConfigInteractor struct {
	DS domain.DataStore
}

func NewBusinessConfigInteractor(ds domain.DataStore) *BusinessConfigInteractor {
	return &BusinessConfigInteractor{
		DS: ds,
	}
}

func (i *BusinessConfigInteractor) ListBySearch(
	ctx context.Context,
	params domain.BusinessConfigSearchParams,
) (res *domain.BusinessConfigSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BusinessConfigInteractor.ListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.BusinessConfigRepo().ListBySearch(ctx, params)
}
