package category

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.CategoryInteractor = (*CategoryInteractor)(nil)

type CategoryInteractor struct {
	DS domain.DataStore
}

func NewCategoryInteractor(ds domain.DataStore) *CategoryInteractor {
	return &CategoryInteractor{
		DS: ds,
	}
}

func (i *CategoryInteractor) ListBySearch(ctx context.Context, params domain.CategorySearchParams) (res domain.Categories, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.ListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.CategoryRepo().ListBySearch(ctx, params)
}
