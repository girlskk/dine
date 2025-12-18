package category

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
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

func (i *CategoryInteractor) PagedListBySearch(ctx context.Context,
	page *upagination.Pagination, params domain.CategorySearchParams,
) (res *domain.CategorySearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.CategoryRepo().PagedListBySearch(ctx, page, params)
}
