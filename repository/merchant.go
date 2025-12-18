package repository

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/merchant"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.MerchantRepository = (*MerchantRepository)(nil)

type MerchantRepository struct {
	Client *ent.Client
}

func NewMerchantRepository(client *ent.Client) *MerchantRepository {
	return &MerchantRepository{
		Client: client,
	}
}

func (repo *MerchantRepository) FindByID(ctx context.Context, id int) (domainMerchant *domain.Merchant, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	em, err := repo.Client.Merchant.Query().
		Where(merchant.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return convertMerchant(em), nil
}

func (repo *MerchantRepository) Create(ctx context.Context, domainMerchant *domain.Merchant) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.Merchant.Create().
		SetMerchantCode(domainMerchant.MerchantCode).
		SetMerchantName(domainMerchant.MerchantName).
		SetMerchantShortName(domainMerchant.MerchantShortName).
		SetMerchantType(domainMerchant.MerchantType).
		SetBrandName(domainMerchant.BrandName).
		SetAdminPhoneNumber(domainMerchant.AdminPhoneNumber).
		SetNillableExpireUtc(domainMerchant.ExpireUTC).
		SetBusinessTypeID(domainMerchant.BusinessTypeID).
		SetMerchantLogo(domainMerchant.MerchantLogo).
		SetDescription(domainMerchant.Description).
		SetStatus(domainMerchant.Status).
		SetLoginAccount(domainMerchant.LoginAccount).
		SetLoginPassword(domainMerchant.LoginPassword).
		SetCountryID(domainMerchant.CountryID).
		SetProvinceID(domainMerchant.ProvinceID).
		SetCityID(domainMerchant.CityID).
		SetDistrictID(domainMerchant.DistrictID).
		SetCountryName(domainMerchant.CountryName).
		SetProvinceName(domainMerchant.ProvinceName).
		SetCityName(domainMerchant.CityName).
		SetDistrictName(domainMerchant.DistrictName).
		SetAddress(domainMerchant.Address).
		SetLng(domainMerchant.Lng).
		SetLat(domainMerchant.Lat).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create merchant: %w", err)
		return
	}
	return
}

func (repo *MerchantRepository) Update(ctx context.Context, domainMerchant *domain.Merchant) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.Merchant.UpdateOneID(domainMerchant.ID).
		SetMerchantCode(domainMerchant.MerchantCode).
		SetMerchantName(domainMerchant.MerchantName).
		SetMerchantShortName(domainMerchant.MerchantShortName).
		SetMerchantType(domainMerchant.MerchantType).
		SetBrandName(domainMerchant.BrandName).
		SetAdminPhoneNumber(domainMerchant.AdminPhoneNumber).
		SetNillableExpireUtc(domainMerchant.ExpireUTC).
		SetBusinessTypeID(domainMerchant.BusinessTypeID).
		SetMerchantLogo(domainMerchant.MerchantLogo).
		SetDescription(domainMerchant.Description).
		SetStatus(domainMerchant.Status).
		SetLoginAccount(domainMerchant.LoginAccount).
		SetLoginPassword(domainMerchant.LoginPassword).
		SetCountryID(domainMerchant.CountryID).
		SetProvinceID(domainMerchant.ProvinceID).
		SetCityID(domainMerchant.CityID).
		SetDistrictID(domainMerchant.DistrictID).
		SetCountryName(domainMerchant.CountryName).
		SetProvinceName(domainMerchant.ProvinceName).
		SetCityName(domainMerchant.CityName).
		SetDistrictName(domainMerchant.DistrictName).
		SetAddress(domainMerchant.Address).
		SetLng(domainMerchant.Lng).
		SetLat(domainMerchant.Lat).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to update merchant: %w", err)
		return
	}
	return
}

func (repo *MerchantRepository) Delete(ctx context.Context, id int) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.Merchant.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to delete merchant: %w", err)
		return
	}

	return
}

func (repo *MerchantRepository) GetMerchants(ctx context.Context, pager *upagination.Pagination, filter *domain.MerchantListFilter, orderBys ...domain.MerchantListOrderBy) (domainMerchants []*domain.Merchant, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.GetMerchants")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.filterBuildQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count: %w", err)
		return
	}

	merchants, err := query.Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query merchants: %w", err)
		return nil, 0, err
	}

	domainMerchants = lo.Map(merchants, func(merchant *ent.Merchant, _ int) *domain.Merchant {
		return convertMerchant(merchant)
	})
	return
}

func (repo *MerchantRepository) CountMerchant(ctx context.Context, condition map[string]string) (merchantCount *domain.MerchantCount, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.CountMerchantByCondition")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	now := time.Now().UTC()
	var counts struct {
		BrandCount   int `sql:"brand_count"`
		StoreCount   int `sql:"store_count"`
		ExpiredCount int `sql:"expired_count"`
	}

	err = repo.Client.Merchant.Query().
		Aggregate(
			func(s *sql.Selector) string {
				// count brand merchants
				return sql.As(
					fmt.Sprintf("SUM(CASE WHEN %s = '%s' THEN 1 ELSE 0 END)", s.C(merchant.FieldMerchantType), domain.MerchantTypeBrand),
					"brand_count",
				)
			},
			func(s *sql.Selector) string {
				// count store merchants
				return sql.As(
					fmt.Sprintf("SUM(CASE WHEN %s = '%s' THEN 1 ELSE 0 END)", s.C(merchant.FieldMerchantType), domain.MerchantTypeStore),
					"store_count",
				)
			},
			func(s *sql.Selector) string {
				// count expired merchants
				return sql.As(
					fmt.Sprintf("SUM(CASE WHEN %s < '%s' THEN 1 ELSE 0 END)", s.C(merchant.FieldExpireUtc), now.Format("2006-01-02 15:04:05")),
					"expired_count",
				)
			},
		).
		Scan(ctx, &counts)
	if err != nil {
		err = fmt.Errorf("failed to count merchant: %w", err)
		return
	}

	merchantCount = &domain.MerchantCount{
		MerchantTypeBrand: counts.BrandCount,
		MerchantTypeStore: counts.StoreCount,
		Expired:           counts.ExpiredCount,
	}
	return
}

func (repo *MerchantRepository) CreateMerchantAndStore(ctx context.Context, domainMerchant *domain.Merchant, domainStore *domain.Store) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.CreateMerchantAndStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	repoTx := New(repo.Client)
	err = repoTx.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		err := ds.MerchantRepo().Create(ctx, domainMerchant)
		if err != nil {
			return err
		}
		err = ds.StoreRepo().Create(ctx, domainStore)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (repo *MerchantRepository) MerchantRenewal(ctx context.Context, merchantRenewal *domain.MerchantRenewal) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.MerchantRenewal")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	repoTx := New(repo.Client)
	err = repoTx.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		err := ds.MerchantRenewalRepo().Create(ctx, merchantRenewal)
		if err != nil {
			return err
		}
		m, err := ds.MerchantRepo().FindByID(ctx, merchantRenewal.MerchantID)
		if err != nil {
			return err
		}
		m.ExpireUTC = sumRenewalDuration(*m.ExpireUTC, merchantRenewal.PurchaseDuration, merchantRenewal.PurchaseDurationUnit)
		err = ds.MerchantRepo().Update(ctx, m)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (repo *MerchantRepository) filterBuildQuery(filter *domain.MerchantListFilter) *ent.MerchantQuery {
	query := repo.Client.Merchant.Query()

	if filter.Status != "" {
		query = query.Where(merchant.StatusEQ(filter.Status))
	}
	if filter.MerchantName != "" {
		query = query.Where(merchant.MerchantNameContains(filter.MerchantName))
	}
	if filter.AdminPhoneNumber != "" {
		query = query.Where(merchant.AdminPhoneNumberEQ(filter.AdminPhoneNumber))
	}
	if filter.MerchantType != "" {
		query = query.Where(merchant.MerchantTypeEQ(filter.MerchantType))
	}
	if filter.CreatedAtGte != nil {
		query = query.Where(merchant.CreatedAtGTE(*filter.CreatedAtGte))
	}
	if filter.CreatedAtLte != nil {
		query = query.Where(merchant.CreatedAtLTE(*filter.CreatedAtLte))
	}
	if filter.ProvinceID > 0 {
		query = query.Where(merchant.ProvinceIDEQ(filter.ProvinceID))
	}
	return query
}

func (repo *MerchantRepository) orderBy(orderBys ...domain.MerchantListOrderBy) []merchant.OrderOption {
	var opts []merchant.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.MerchantListOrderByID:
			opts = append(opts, merchant.ByID(rule))
		case domain.MerchantListOrderByCreatedAt:
			opts = append(opts, merchant.ByCreatedAt(rule))
		}
	}

	if len(opts) == 0 {
		opts = append(opts, merchant.ByCreatedAt(sql.OrderDesc()))
	}

	return opts
}

func convertMerchant(em *ent.Merchant) *domain.Merchant {
	return &domain.Merchant{
		ID:                em.ID,
		MerchantCode:      em.MerchantCode,
		MerchantName:      em.MerchantName,
		MerchantShortName: em.MerchantShortName,
		MerchantType:      em.MerchantType,
		BrandName:         em.BrandName,
		AdminPhoneNumber:  em.AdminPhoneNumber,
		ExpireUTC:         em.ExpireUtc,
		BusinessTypeID:    em.BusinessTypeID,
		MerchantLogo:      em.MerchantLogo,
		Description:       em.Description,
		Status:            em.Status,
		LoginAccount:      em.LoginAccount,
		LoginPassword:     em.LoginPassword,
		CountryID:         em.CountryID,
		ProvinceID:        em.ProvinceID,
		CityID:            em.CityID,
		DistrictID:        em.DistrictID,
		CountryName:       em.CountryName,
		ProvinceName:      em.ProvinceName,
		CityName:          em.CityName,
		DistrictName:      em.DistrictName,
		Address:           em.Address,
		Lng:               em.Lng,
		Lat:               em.Lat,
		CreatedAt:         em.CreatedAt,
		UpdatedAt:         em.UpdatedAt,
	}
}

func sumRenewalDuration(oldTime time.Time, d int, durationUnit domain.PurchaseDurationUnit) *time.Time {
	newTime := oldTime
	switch durationUnit {
	case domain.PurchaseDurationUnitDay:
		newTime = oldTime.AddDate(0, 0, d)
	case domain.PurchaseDurationUnitMonth:
		newTime = oldTime.AddDate(0, d, 0)
	case domain.PurchaseDurationUnitYear:
		newTime = oldTime.AddDate(d, 0, 0)
	case domain.PurchaseDurationUnitWeek:
		newTime = oldTime.AddDate(0, 0, d*7)
	default:
	}
	return &newTime
}
