package usecase

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StoreAccountInteractor = (*StoreAccountInteractor)(nil)

type StoreAccountInteractor struct {
	ds domain.DataStore
}

func NewStoreAccountInteractor(dataStore domain.DataStore) *StoreAccountInteractor {
	return &StoreAccountInteractor{
		ds: dataStore,
	}
}

func (s *StoreAccountInteractor) GetDetail(ctx context.Context, storeID int) (res *domain.StoreAccount, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreAccountInteractor.GetDetail")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return s.ds.StoreAccountRepo().FindByStore(ctx, storeID)
}

func (s *StoreAccountInteractor) PagedListTransactions(ctx context.Context, page *upagination.Pagination,
	params domain.StoreAccountTransactionSearchParams,
) (res *domain.StoreAccountTransactionSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreAccountInteractor.PagedListTransactions")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return s.ds.StoreAccountRepo().PagedListTransactions(ctx, page, params)
}
