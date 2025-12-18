package usecase

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.MerchantInteractor = (*MerchantInteractor)(nil)

type MerchantInteractor struct {
	DataStore domain.DataStore
}

func (interactor *MerchantInteractor) CreateMerchant(ctx context.Context, domainMerchant *domain.Merchant, domainStore *domain.Store) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.CreateMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return
}

func (interactor *MerchantInteractor) UpdateMerchant(ctx context.Context, domainMerchant *domain.Merchant, domainStore *domain.Store) (err error) {
	//TODO implement me
	panic("implement me")
}

func (interactor *MerchantInteractor) DeleteMerchant(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (interactor *MerchantInteractor) GetMerchant(ctx context.Context, id int) (domainMerchant *domain.Merchant, err error) {
	//TODO implement me
	panic("implement me")
}

func (interactor *MerchantInteractor) GetMerchants(ctx context.Context, pager *upagination.Pagination, filter *domain.MerchantListFilter, orderBys ...domain.MerchantListOrderBy) (domainMerchants []*domain.Merchant, total int, err error) {
	//TODO implement me
	panic("implement me")
}

func (interactor *MerchantInteractor) CountMerchant(ctx context.Context, condition map[string]string) (merchantCount *domain.MerchantCount, err error) {
	//TODO implement me
	panic("implement me")
}

func (interactor *MerchantInteractor) MerchantRenewal(ctx context.Context, merchantRenewal *domain.MerchantRenewal) (err error) {
	//TODO implement me
	panic("implement me")
}

func NewMerchantInteractor(dataStore domain.DataStore) *MerchantInteractor {
	return &MerchantInteractor{
		DataStore: dataStore,
	}
}
