package store

import (
	"context"
	"fmt"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StoreInteractor = (*StoreInteractor)(nil)

type StoreInteractor struct {
	DataStore domain.DataStore
}

func (interactor *StoreInteractor) CreateStore(ctx context.Context, domainCStore *domain.CreateStoreParams) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.CreateStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainStore, err := interactor.CheckCreateStoreFields(ctx, domainCStore)
	if err != nil {
		return
	}

	err = interactor.DataStore.StoreRepo().Create(ctx, domainStore)
	if err != nil {
		err = fmt.Errorf("failed to create store: %w", err)
		return
	}

	return
}

func (interactor *StoreInteractor) UpdateStore(ctx context.Context, domainUStore *domain.UpdateStoreParams) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.UpdateStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainStore, err := interactor.CheckUpdateStoreFields(ctx, domainUStore)
	if err != nil {
		return
	}
	err = interactor.DataStore.StoreRepo().Update(ctx, domainStore)
	if err != nil {
		err = fmt.Errorf("failed to update store: %w", err)
		return
	}

	return
}

func (interactor *StoreInteractor) DeleteStore(ctx context.Context, id int) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.DeleteStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = interactor.DataStore.StoreRepo().Delete(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to delete store: %w", err)
		return
	}

	return
}

func (interactor *StoreInteractor) GetStore(ctx context.Context, id int) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.GetStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	domainStore, err = interactor.DataStore.StoreRepo().FindByID(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to get store: %w", err)
		return
	}

	return
}

func (interactor *StoreInteractor) GetStores(ctx context.Context, pager *upagination.Pagination, filter *domain.StoreListFilter, orderBys ...domain.StoreListOrderBy) (domainStores []*domain.Store, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.GetStores")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	domainStores, total, err = interactor.DataStore.StoreRepo().GetStores(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get stores: %w", err)
		return
	}
	return
}

func (interactor *StoreInteractor) CheckCreateStoreFields(ctx context.Context, domainCStore *domain.CreateStoreParams) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.CheckCreateStoreFields")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainStore = &domain.Store{
		MerchantID:              domainCStore.MerchantID,
		AdminPhoneNumber:        domainCStore.AdminPhoneNumber,
		StoreName:               domainCStore.StoreName,
		StoreShortName:          domainCStore.StoreShortName,
		StoreCode:               domainCStore.StoreCode,
		Status:                  domainCStore.Status,
		BusinessModel:           domainCStore.BusinessModel,
		BusinessTypeID:          domainCStore.BusinessTypeID,
		ContactName:             domainCStore.ContactName,
		ContactPhone:            domainCStore.ContactPhone,
		UnifiedSocialCreditCode: domainCStore.UnifiedSocialCreditCode,
		StoreLogo:               domainCStore.StoreLogo,
		BusinessLicenseURL:      domainCStore.BusinessLicenseURL,
		StorefrontURL:           domainCStore.StorefrontURL,
		CashierDeskURL:          domainCStore.CashierDeskURL,
		DiningEnvironmentURL:    domainCStore.DiningEnvironmentURL,
		FoodOperationLicenseURL: domainCStore.FoodOperationLicenseURL,
		CountryID:               domainCStore.CountryID,
		ProvinceID:              domainCStore.ProvinceID,
		CityID:                  domainCStore.CityID,
		DistrictID:              domainCStore.DistrictID,
		Address:                 domainCStore.Address,
		Lng:                     domainCStore.Lng,
		Lat:                     domainCStore.Lat,
		LoginAccount:            domainCStore.LoginAccount,
		LoginPassword:           domainCStore.LoginPassword,
	}

	err = interactor.checkFields(ctx, domainStore)
	if err != nil {
		return
	}

	return
}

func (interactor *StoreInteractor) CheckUpdateStoreFields(ctx context.Context, domainUStore *domain.UpdateStoreParams) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.CheckUpdateStoreFields")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	oldStore, err := interactor.DataStore.StoreRepo().FindByID(ctx, domainUStore.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrStoreNotExists)
			return
		}
		return
	}

	domainStore = &domain.Store{
		ID:                      domainUStore.ID,
		MerchantID:              oldStore.MerchantID,
		AdminPhoneNumber:        domainUStore.AdminPhoneNumber,
		StoreName:               domainUStore.StoreName,
		StoreShortName:          domainUStore.StoreShortName,
		StoreCode:               domainUStore.StoreCode,
		Status:                  domainUStore.Status,
		BusinessModel:           domainUStore.BusinessModel,
		BusinessTypeID:          domainUStore.BusinessTypeID,
		ContactName:             domainUStore.ContactName,
		ContactPhone:            domainUStore.ContactPhone,
		UnifiedSocialCreditCode: domainUStore.UnifiedSocialCreditCode,
		StoreLogo:               domainUStore.StoreLogo,
		BusinessLicenseURL:      domainUStore.BusinessLicenseURL,
		StorefrontURL:           domainUStore.StorefrontURL,
		CashierDeskURL:          domainUStore.CashierDeskURL,
		DiningEnvironmentURL:    domainUStore.DiningEnvironmentURL,
		FoodOperationLicenseURL: domainUStore.FoodOperationLicenseURL,
		CountryID:               domainUStore.CountryID,
		ProvinceID:              domainUStore.ProvinceID,
		CityID:                  domainUStore.CityID,
		DistrictID:              domainUStore.DistrictID,
		CountryName:             oldStore.CountryName,
		ProvinceName:            oldStore.ProvinceName,
		CityName:                oldStore.CityName,
		DistrictName:            oldStore.DistrictName,
		Address:                 domainUStore.Address,
		Lng:                     domainUStore.Lng,
		Lat:                     domainUStore.Lat,
		LoginAccount:            domainUStore.LoginAccount,
		LoginPassword:           domainUStore.LoginPassword,
		AdminUserID:             oldStore.AdminUserID,
	}

	err = interactor.checkFields(ctx, domainStore)
	if err != nil {
		return
	}

	return
}

func (interactor *StoreInteractor) checkFields(ctx context.Context, domainStore *domain.Store) (err error) {
	//TODO implement me
	panic("implement me")
}

func NewStoreInteractor(dataStore domain.DataStore) *StoreInteractor {
	return &StoreInteractor{
		DataStore: dataStore,
	}
}
