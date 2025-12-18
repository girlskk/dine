package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/store"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StoreRepository = (*StoreRepository)(nil)

type StoreRepository struct {
	Client *ent.Client
}

func (repo *StoreRepository) Create(ctx context.Context, domainStore *domain.Store) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	if domainStore == nil {
		err = fmt.Errorf("domainStore is nil")
		return
	}
	_, err = repo.Client.Store.Create().
		SetMerchantID(domainStore.MerchantID).
		SetAdminPhoneNumber(domainStore.AdminPhoneNumber).
		SetStoreName(domainStore.StoreName).
		SetStoreShortName(domainStore.StoreShortName).
		SetStoreCode(domainStore.StoreCode).
		SetStatus(domainStore.Status).
		SetBusinessModel(domainStore.BusinessModel).
		SetBusinessTypeID(domainStore.BusinessTypeID).
		SetContactName(domainStore.ContactName).
		SetContactPhone(domainStore.ContactPhone).
		SetUnifiedSocialCreditCode(domainStore.UnifiedSocialCreditCode).
		SetStoreLogo(domainStore.StoreLogo).
		SetBusinessLicenseURL(domainStore.BusinessLicenseURL).
		SetStorefrontURL(domainStore.StorefrontURL).
		SetCashierDeskURL(domainStore.CashierDeskURL).
		SetDiningEnvironmentURL(domainStore.DiningEnvironmentURL).
		SetFoodOperationLicenseURL(domainStore.FoodOperationLicenseURL).
		SetCountryID(domainStore.CountryID).
		SetProvinceID(domainStore.ProvinceID).
		SetCityID(domainStore.CityID).
		SetDistrictID(domainStore.DistrictID).
		SetCountryName(domainStore.CountryName).
		SetProvinceName(domainStore.ProvinceName).
		SetCityName(domainStore.CityName).
		SetDistrictName(domainStore.DistrictName).
		SetAddress(domainStore.Address).
		SetLng(domainStore.Lng).
		SetLat(domainStore.Lat).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create store: %w", err)
		return
	}
	return
}

func (repo *StoreRepository) Update(ctx context.Context, domainStore *domain.Store) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	if domainStore == nil {
		err = fmt.Errorf("domainStore is nil")
		return
	}
	_, err = repo.Client.Store.UpdateOneID(domainStore.ID).
		SetAdminPhoneNumber(domainStore.AdminPhoneNumber).
		SetStoreName(domainStore.StoreName).
		SetStoreShortName(domainStore.StoreShortName).
		SetStoreCode(domainStore.StoreCode).
		SetStatus(domainStore.Status).
		SetBusinessModel(domainStore.BusinessModel).
		SetBusinessTypeID(domainStore.BusinessTypeID).
		SetContactName(domainStore.ContactName).
		SetContactPhone(domainStore.ContactPhone).
		SetUnifiedSocialCreditCode(domainStore.UnifiedSocialCreditCode).
		SetStoreLogo(domainStore.StoreLogo).
		SetBusinessLicenseURL(domainStore.BusinessLicenseURL).
		SetStorefrontURL(domainStore.StorefrontURL).
		SetCashierDeskURL(domainStore.CashierDeskURL).
		SetDiningEnvironmentURL(domainStore.DiningEnvironmentURL).
		SetFoodOperationLicenseURL(domainStore.FoodOperationLicenseURL).
		SetCountryID(domainStore.CountryID).
		SetProvinceID(domainStore.ProvinceID).
		SetCityID(domainStore.CityID).
		SetDistrictID(domainStore.DistrictID).
		SetCountryName(domainStore.CountryName).
		SetProvinceName(domainStore.ProvinceName).
		SetCityName(domainStore.CityName).
		SetDistrictName(domainStore.DistrictName).
		SetAddress(domainStore.Address).
		SetLng(domainStore.Lng).
		SetLat(domainStore.Lat).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to update store: %w", err)
		return
	}
	return
}

func (repo *StoreRepository) Delete(ctx context.Context, id int) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.Store.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to delete merchant: %w", err)
		return
	}

	return
}

func (repo *StoreRepository) FindByID(ctx context.Context, id int) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	em, err := repo.Client.Store.Query().
		Where(store.ID(id)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, domain.NotFoundError(domain.ErrStoreNotExists)
	}
	return convertStore(em), nil
}

func (repo *StoreRepository) GetStores(ctx context.Context, pager *upagination.Pagination, filter *domain.StoreListFilter, orderBys ...domain.StoreListOrderBy) (domainStores []*domain.Store, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreRepository.GetStores")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if pager == nil {
		err = fmt.Errorf("pager is nil")
		return
	}
	if filter == nil {
		err = fmt.Errorf("filter is nil")
		return
	}

	query := repo.filterBuildQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count: %w", err)
		return
	}

	stores, err := query.Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query stores: %w", err)
		return nil, 0, err
	}

	domainStores = lo.Map(stores, func(store *ent.Store, _ int) *domain.Store {
		return convertStore(store)
	})
	return
}

func (repo *StoreRepository) ExistsStore(ctx context.Context, existsStoreParams *domain.ExistsStoreParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreRepository.ExistsStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if existsStoreParams == nil {
		err = fmt.Errorf("existsStoreParams is nil")
		return
	}

	query := repo.Client.Store.Query().
		Where(store.StoreNameEQ(existsStoreParams.StoreName))
	if existsStoreParams.NotID > 0 {
		query = query.Where(store.IDNEQ(existsStoreParams.NotID))
	}

	exists, err = query.Exist(ctx)
	if err != nil {
		err = fmt.Errorf("failed to check store existence: %w", err)
		return
	}

	return
}

func convertStore(es *ent.Store) *domain.Store {
	return &domain.Store{
		ID:                      es.ID,
		MerchantID:              es.MerchantID,
		AdminPhoneNumber:        es.AdminPhoneNumber,
		StoreName:               es.StoreName,
		StoreShortName:          es.StoreShortName,
		StoreCode:               es.StoreCode,
		Status:                  es.Status,
		BusinessModel:           es.BusinessModel,
		BusinessTypeID:          es.BusinessTypeID,
		ContactName:             es.ContactName,
		ContactPhone:            es.ContactPhone,
		UnifiedSocialCreditCode: es.UnifiedSocialCreditCode,
		StoreLogo:               es.StoreLogo,
		BusinessLicenseURL:      es.BusinessLicenseURL,
		StorefrontURL:           es.StorefrontURL,
		CashierDeskURL:          es.CashierDeskURL,
		DiningEnvironmentURL:    es.DiningEnvironmentURL,
		FoodOperationLicenseURL: es.FoodOperationLicenseURL,
		CountryID:               es.CountryID,
		ProvinceID:              es.ProvinceID,
		CityID:                  es.CityID,
		DistrictID:              es.DistrictID,
		CountryName:             es.CountryName,
		ProvinceName:            es.ProvinceName,
		CityName:                es.CityName,
		DistrictName:            es.DistrictName,
		Address:                 es.Address,
		Lng:                     es.Lng,
		Lat:                     es.Lat,
		CreatedAt:               es.CreatedAt,
		UpdatedAt:               es.UpdatedAt,
	}
}

func (repo *StoreRepository) filterBuildQuery(filter *domain.StoreListFilter) *ent.StoreQuery {
	query := repo.Client.Store.Query()

	if filter.StoreName != "" {
		query = query.Where(store.StoreNameContains(filter.StoreName))
	}
	if filter.MerchantID != 0 {
		query = query.Where(store.MerchantIDEQ(filter.MerchantID))
	}
	if filter.BusinessTypeID != 0 {
		query = query.Where(store.BusinessTypeIDEQ(filter.BusinessTypeID))
	}
	if filter.AdminPhoneNumber != "" {
		query = query.Where(store.AdminPhoneNumberEQ(filter.AdminPhoneNumber))
	}
	if filter.Status != "" {
		query = query.Where(store.StatusEQ(filter.Status))
	}
	if filter.BusinessModel != "" {
		query = query.Where(store.BusinessModelEQ(filter.BusinessModel))
	}
	if filter.CreatedAtGte != nil {
		query = query.Where(store.CreatedAtGTE(*filter.CreatedAtGte))
	}
	if filter.CreatedAtLte != nil {
		query = query.Where(store.CreatedAtLTE(*filter.CreatedAtLte))
	}
	if filter.ProvinceID != 0 {
		query = query.Where(store.ProvinceIDEQ(filter.ProvinceID))
	}

	return query
}

func (repo *StoreRepository) orderBy(orderBys ...domain.StoreListOrderBy) []store.OrderOption {
	var opts []store.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.StoreListOrderByID:
			opts = append(opts, store.ByID(rule))
		case domain.StoreListOrderByCreatedAt:
			opts = append(opts, store.ByCreatedAt(rule))
		}
	}

	if len(opts) == 0 {
		opts = append(opts, store.ByCreatedAt(sql.OrderDesc()))
	}

	return opts
}

func NewStoreRepository(client *ent.Client) *StoreRepository {
	return &StoreRepository{
		Client: client,
	}
}
