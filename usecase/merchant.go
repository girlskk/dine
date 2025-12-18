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

func (interactor *MerchantInteractor) CreateMerchant(ctx context.Context, domainMerchant *domain.Merchant) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.CreateMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = interactor.checkFields(ctx, domainMerchant)
	if err != nil {
		return
	}

	err = interactor.DataStore.MerchantRepo().Create(ctx, domainMerchant)
	if err != nil {
		return err
	}
	return
}

func (interactor *MerchantInteractor) checkFields(ctx context.Context, domainMerchant *domain.Merchant) (err error) {
	exists, err := interactor.DataStore.MerchantRepo().ExistMerchant(ctx, &domain.MerchantExistsParams{
		MerchantName: domainMerchant.MerchantName,
		NotID:        domainMerchant.ID,
	})
	if err != nil {
		return err
	}
	if exists {
		return domain.ParamsError(domain.ErrMerchantNameExists)
	}
	return
}

func (interactor *MerchantInteractor) CreateMerchantAndStore(ctx context.Context, domainMerchant *domain.Merchant, domainStore *domain.Store) (err error) {
	//TODO implement me
	panic("implement me")
}

func (interactor *MerchantInteractor) UpdateMerchant(ctx context.Context, domainMerchant *domain.Merchant) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.UpdateMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = interactor.checkFields(ctx, domainMerchant)
	if err != nil {
		return
	}

	err = interactor.DataStore.MerchantRepo().Update(ctx, domainMerchant)
	if err != nil {
		return err
	}
	return
}

func (interactor *MerchantInteractor) UpdateMerchantAndStore(ctx context.Context, domainMerchant *domain.Merchant, domainStore *domain.Store) (err error) {
	//TODO implement me
	panic("implement me")
}

func (interactor *MerchantInteractor) DeleteMerchant(ctx context.Context, id int) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.DeleteMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return interactor.DataStore.MerchantRepo().Delete(ctx, id)
}

func (interactor *MerchantInteractor) GetMerchant(ctx context.Context, id int) (domainMerchant *domain.Merchant, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.GetMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainMerchant, err = interactor.DataStore.MerchantRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ParamsError(domain.ErrMerchantNotExists)
		}
		return nil, err
	}
	return
}

func (interactor *MerchantInteractor) GetMerchants(ctx context.Context, pager *upagination.Pagination, filter *domain.MerchantListFilter, orderBys ...domain.MerchantListOrderBy) (domainMerchants []*domain.Merchant, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.GetMerchants")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return interactor.DataStore.MerchantRepo().GetMerchants(ctx, pager, filter, orderBys...)
}

func (interactor *MerchantInteractor) CountMerchant(ctx context.Context) (merchantCount *domain.MerchantCount, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.CountMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return interactor.DataStore.MerchantRepo().CountMerchant(ctx)
}

func (interactor *MerchantInteractor) MerchantRenewal(ctx context.Context, merchantRenewal *domain.MerchantRenewal) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.MerchantRenewal")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return interactor.DataStore.MerchantRepo().MerchantRenewal(ctx, merchantRenewal)
}

func NewMerchantInteractor(dataStore domain.DataStore) *MerchantInteractor {
	return &MerchantInteractor{
		DataStore: dataStore,
	}
}
