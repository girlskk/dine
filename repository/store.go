package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
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
	if domainStore.Address == nil {
		err = fmt.Errorf("address is nil")
		return
	}
	builder := repo.Client.Store.Create().SetID(domainStore.ID).
		SetMerchantID(domainStore.MerchantID).
		SetAdminPhoneNumber(domainStore.AdminPhoneNumber).
		SetStoreName(domainStore.StoreName).
		SetStoreShortName(domainStore.StoreShortName).
		SetStoreCode(domainStore.StoreCode).
		SetStatus(domainStore.Status).
		SetBusinessModel(domainStore.BusinessModel).
		SetBusinessTypeID(domainStore.BusinessTypeID).
		SetLocationNumber(domainStore.LocationNumber).
		SetContactName(domainStore.ContactName).
		SetContactPhone(domainStore.ContactPhone).
		SetUnifiedSocialCreditCode(domainStore.UnifiedSocialCreditCode).
		SetStoreLogo(domainStore.StoreLogo).
		SetBusinessLicenseURL(domainStore.BusinessLicenseURL).
		SetStorefrontURL(domainStore.StorefrontURL).
		SetCashierDeskURL(domainStore.CashierDeskURL).
		SetDiningEnvironmentURL(domainStore.DiningEnvironmentURL).
		SetFoodOperationLicenseURL(domainStore.FoodOperationLicenseURL).
		SetBusinessHours(domainStore.BusinessHours).
		SetDiningPeriods(domainStore.DiningPeriods).
		SetShiftTimes(domainStore.ShiftTimes).
		SetCountryID(domainStore.Address.CountryID).
		SetProvinceID(domainStore.Address.ProvinceID).
		SetCityID(domainStore.Address.CityID).
		SetDistrictID(domainStore.Address.DistrictID).
		SetAddress(domainStore.Address.Address).
		SetLng(domainStore.Address.Lng).
		SetLat(domainStore.Address.Lat).
		SetSuperAccount(domainStore.LoginAccount)
	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create store: %w", err)
		return
	}
	domainStore.CreatedAt = created.CreatedAt
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
	if domainStore.Address == nil {
		err = fmt.Errorf("address is nil")
		return
	}

	builder := repo.Client.Store.UpdateOneID(domainStore.ID).
		SetAdminPhoneNumber(domainStore.AdminPhoneNumber).
		SetStoreName(domainStore.StoreName).
		SetStoreShortName(domainStore.StoreShortName).
		SetStoreCode(domainStore.StoreCode).
		SetStatus(domainStore.Status).
		SetBusinessModel(domainStore.BusinessModel).
		SetBusinessTypeID(domainStore.BusinessTypeID).
		SetLocationNumber(domainStore.LocationNumber).
		SetContactName(domainStore.ContactName).
		SetContactPhone(domainStore.ContactPhone).
		SetUnifiedSocialCreditCode(domainStore.UnifiedSocialCreditCode).
		SetStoreLogo(domainStore.StoreLogo).
		SetBusinessLicenseURL(domainStore.BusinessLicenseURL).
		SetStorefrontURL(domainStore.StorefrontURL).
		SetCashierDeskURL(domainStore.CashierDeskURL).
		SetDiningEnvironmentURL(domainStore.DiningEnvironmentURL).
		SetFoodOperationLicenseURL(domainStore.FoodOperationLicenseURL).
		SetBusinessHours(domainStore.BusinessHours).
		SetDiningPeriods(domainStore.DiningPeriods).
		SetShiftTimes(domainStore.ShiftTimes).
		SetCountryID(domainStore.Address.CountryID).
		SetProvinceID(domainStore.Address.ProvinceID).
		SetCityID(domainStore.Address.CityID).
		SetDistrictID(domainStore.Address.DistrictID).
		SetAddress(domainStore.Address.Address).
		SetLng(domainStore.Address.Lng).
		SetLat(domainStore.Address.Lat)
	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to update store: %w", err)
		return
	}
	domainStore.UpdatedAt = updated.UpdatedAt
	return
}

func (repo *StoreRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.Store.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to delete store: %w", err)
		return
	}

	return
}

func (repo *StoreRepository) FindByID(ctx context.Context, id uuid.UUID) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	em, err := repo.Client.Store.Query().
		Where(store.ID(id)).
		WithMerchant().
		WithCountry().
		WithProvince().
		WithCity().
		WithDistrict().
		WithMerchantBusinessType().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrStoreNotExists)
			return
		}
		err = fmt.Errorf("failed to find store by id: %w", err)
		return
	}
	domainStore = convertStore(em)
	return
}

func (repo *StoreRepository) FindStoreMerchant(ctx context.Context, merchantID uuid.UUID) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreRepository.FindStoreMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if merchantID == uuid.Nil {
		err = fmt.Errorf("merchantID is nil")
		return
	}

	em, err := repo.Client.Store.Query().
		Where(store.MerchantIDEQ(merchantID)).
		WithCountry().
		WithProvince().
		WithCity().
		WithDistrict().
		WithMerchantBusinessType().
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, domain.NotFoundError(domain.ErrStoreNotExists)
	}
	if err != nil {
		err = fmt.Errorf("failed to find store by merchant id: %w", err)
		return
	}
	domainStore = convertStore(em)
	return
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
	query.
		WithMerchant(). // 加载商户信息
		WithProvince()  // 加载省信息
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

func (repo *StoreRepository) CountStoresByMerchantID(ctx context.Context, merchantIDs []uuid.UUID) (storeCounts []*domain.MerchantStoreCount, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreRepository.CountStoresByMerchantID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(merchantIDs) == 0 {
		err = fmt.Errorf("merchantIDs is empty")
		return
	}

	type result struct {
		MerchantID uuid.UUID `json:"merchant_id"`
		Count      int       `json:"count"`
	}

	var results []result
	err = repo.Client.Store.Query().
		Where(store.MerchantIDIn(merchantIDs...)).
		GroupBy(store.FieldMerchantID).
		Aggregate(ent.Count()).
		Scan(ctx, &results)
	if err != nil {
		err = fmt.Errorf("failed to query merchant store counts: %w", err)
		return
	}

	storeCounts = lo.Map(results, func(r result, _ int) *domain.MerchantStoreCount {
		return &domain.MerchantStoreCount{
			MerchantID: r.MerchantID,
			StoreCount: r.Count,
		}
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
	if existsStoreParams.MerchantID != uuid.Nil {
		query = query.Where(store.MerchantIDEQ(existsStoreParams.MerchantID))
	}
	if existsStoreParams.ExcludeID != uuid.Nil {
		query = query.Where(store.IDNEQ(existsStoreParams.ExcludeID))
	}

	exists, err = query.Exist(ctx)
	if err != nil {
		err = fmt.Errorf("failed to check store existence: %w", err)
		return
	}

	return
}

func convertStore(es *ent.Store) *domain.Store {
	address := &domain.Address{
		CountryID:  es.CountryID,
		ProvinceID: es.ProvinceID,
		CityID:     es.CityID,
		DistrictID: es.DistrictID,
		Address:    es.Address,
		Lng:        es.Lng,
		Lat:        es.Lat,
	}
	if es.Edges.Country != nil {
		address.CountryName = es.Edges.Country.Name
	}
	if es.Edges.Province != nil {
		address.ProvinceName = es.Edges.Province.Name
	}
	if es.Edges.City != nil {
		address.CityName = es.Edges.City.Name
	}
	if es.Edges.District != nil {
		address.DistrictName = es.Edges.District.Name
	}
	repoStore := &domain.Store{
		ID:                      es.ID,
		MerchantID:              es.MerchantID,
		AdminPhoneNumber:        es.AdminPhoneNumber,
		StoreName:               es.StoreName,
		StoreShortName:          es.StoreShortName,
		StoreCode:               es.StoreCode,
		Status:                  es.Status,
		BusinessModel:           es.BusinessModel,
		BusinessTypeID:          es.BusinessTypeID,
		LocationNumber:          es.LocationNumber,
		ContactName:             es.ContactName,
		ContactPhone:            es.ContactPhone,
		UnifiedSocialCreditCode: es.UnifiedSocialCreditCode,
		StoreLogo:               es.StoreLogo,
		BusinessLicenseURL:      es.BusinessLicenseURL,
		StorefrontURL:           es.StorefrontURL,
		CashierDeskURL:          es.CashierDeskURL,
		DiningEnvironmentURL:    es.DiningEnvironmentURL,
		FoodOperationLicenseURL: es.FoodOperationLicenseURL,
		BusinessHours:           es.BusinessHours,
		DiningPeriods:           es.DiningPeriods,
		ShiftTimes:              es.ShiftTimes,
		LoginAccount:            es.SuperAccount,
		Address:                 address,
		CreatedAt:               es.CreatedAt,
		UpdatedAt:               es.UpdatedAt,
	}

	if es.Edges.Merchant != nil {
		repoStore.MerchantName = es.Edges.Merchant.MerchantName
	}

	if es.Edges.MerchantBusinessType != nil {
		repoStore.BusinessTypeName = es.Edges.MerchantBusinessType.TypeName
	}
	return repoStore
}

func (repo *StoreRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) (domainStores []*domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreRepository.ListByIDs")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(ids) == 0 {
		return nil, nil
	}

	stores, err := repo.Client.Store.Query().Where(store.IDIn(ids...)).All(ctx)
	if err != nil {
		return nil, err
	}

	domainStores = lo.Map(stores, func(store *ent.Store, _ int) *domain.Store {
		return convertStore(store)
	})
	return domainStores, nil
}

func (repo *StoreRepository) filterBuildQuery(filter *domain.StoreListFilter) *ent.StoreQuery {
	query := repo.Client.Store.Query()

	if filter.StoreName != "" {
		query = query.Where(store.StoreNameContains(filter.StoreName))
	}
	if filter.MerchantID != uuid.Nil {
		query = query.Where(store.MerchantIDEQ(filter.MerchantID))
	}
	if filter.BusinessTypeID != uuid.Nil {
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
	if filter.ProvinceID != uuid.Nil {
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
