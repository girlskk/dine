package repository

import (
	"context"
	"fmt"
	"sort"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
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

func (repo *MerchantRepository) Create(ctx context.Context, domainMerchant *domain.Merchant) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if domainMerchant == nil {
		err = fmt.Errorf("domainMerchant is nil")
		return
	}

	builder := repo.Client.Merchant.Create().SetID(domainMerchant.ID).
		SetMerchantCode(domainMerchant.MerchantCode).
		SetMerchantName(domainMerchant.MerchantName).
		SetMerchantShortName(domainMerchant.MerchantShortName).
		SetMerchantType(domainMerchant.MerchantType).
		SetBrandName(domainMerchant.BrandName).
		SetAdminPhoneNumber(domainMerchant.AdminPhoneNumber).
		SetNillableExpireUtc(domainMerchant.ExpireUTC).
		SetBusinessTypeCode(domainMerchant.BusinessTypeCode).
		SetMerchantLogo(domainMerchant.MerchantLogo).
		SetDescription(domainMerchant.Description).
		SetStatus(domainMerchant.Status).
		SetSuperAccount(domainMerchant.LoginAccount)

	if domainMerchant.Address != nil {
		if domainMerchant.Address.Country != "" {
			builder = builder.SetCountry(domainMerchant.Address.Country)
		}
		if domainMerchant.Address.Province != "" {
			builder = builder.SetProvince(domainMerchant.Address.Province)
		}
		builder = builder.SetAddress(domainMerchant.Address.Address).
			SetLng(domainMerchant.Address.Lng).
			SetLat(domainMerchant.Address.Lat)
	}
	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create merchant: %w", err)
		return
	}
	domainMerchant.CreatedAt = created.CreatedAt
	return
}

func (repo *MerchantRepository) Update(ctx context.Context, domainMerchant *domain.Merchant) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if domainMerchant == nil {
		err = fmt.Errorf("domainMerchant is nil")
		return
	}

	builder := repo.Client.Merchant.UpdateOneID(domainMerchant.ID).
		SetMerchantName(domainMerchant.MerchantName).
		SetMerchantShortName(domainMerchant.MerchantShortName).
		SetMerchantType(domainMerchant.MerchantType).
		SetBrandName(domainMerchant.BrandName).
		SetAdminPhoneNumber(domainMerchant.AdminPhoneNumber).
		SetNillableExpireUtc(domainMerchant.ExpireUTC).
		SetBusinessTypeCode(domainMerchant.BusinessTypeCode).
		SetMerchantLogo(domainMerchant.MerchantLogo).
		SetDescription(domainMerchant.Description).
		SetStatus(domainMerchant.Status)
	if len(domainMerchant.MerchantCode) > 0 {
		builder = builder.SetMerchantCode(domainMerchant.MerchantCode)
		builder = builder.ClearMerchantCode()
	}
	if len(domainMerchant.Description) > 0 {
		builder = builder.SetDescription(domainMerchant.Description)
	} else {
		builder = builder.ClearDescription()
	}
	if domainMerchant.Address != nil {
		if domainMerchant.Address.Country != "" {
			builder = builder.SetCountry(domainMerchant.Address.Country)
		} else {
			builder = builder.ClearCountry()
		}
		if domainMerchant.Address.Province != "" {
			builder = builder.SetProvince(domainMerchant.Address.Province)
		} else {
			builder = builder.ClearProvince()
		}
		builder = builder.SetAddress(domainMerchant.Address.Address).
			SetLng(domainMerchant.Address.Lng).
			SetLat(domainMerchant.Address.Lat)
	}
	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to update merchant: %w", err)
		return
	}
	domainMerchant.UpdatedAt = updated.UpdatedAt
	return
}

func (repo *MerchantRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.Merchant.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to delete merchant: %w", err)
		return
	}

	return
}

func (repo *MerchantRepository) FindByID(ctx context.Context, id uuid.UUID) (domainMerchant *domain.Merchant, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	em, err := repo.Client.Merchant.Query().
		Where(merchant.ID(id)).
		WithMerchantRenewals().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrMerchantNotExists)
			return
		}
		err = fmt.Errorf("failed to find merchant by id: %w", err)
		return
	}

	domainMerchant = convertMerchant(em)
	return
}

func (repo *MerchantRepository) GetMerchants(ctx context.Context, pager *upagination.Pagination, filter *domain.MerchantListFilter, orderBys ...domain.MerchantListOrderBy) (domainMerchants []*domain.Merchant, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.GetMerchants")
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

	merchants, err := query.Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query merchants: %w", err)
		return
	}

	domainMerchants = lo.Map(merchants, func(merchant *ent.Merchant, _ int) *domain.Merchant {
		return convertMerchant(merchant)
	})
	return
}

func (repo *MerchantRepository) CountMerchant(ctx context.Context) (merchantCount *domain.MerchantCount, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.CountMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	now := time.Now().UTC()
	var counts []struct {
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

	if len(counts) == 0 {
		merchantCount = &domain.MerchantCount{}
		return
	}

	merchantCount = &domain.MerchantCount{
		MerchantTypeBrand: counts[0].BrandCount,
		MerchantTypeStore: counts[0].StoreCount,
		Expired:           counts[0].ExpiredCount,
	}
	return
}

func (repo *MerchantRepository) ExistMerchant(ctx context.Context, merchantExistsParams *domain.MerchantExistsParams) (exist bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.ExistMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if merchantExistsParams == nil {
		err = fmt.Errorf("merchantExistsParams is nil")
		return
	}
	query := repo.Client.Merchant.Query().
		Where(merchant.MerchantNameEQ(merchantExistsParams.MerchantName))
	if merchantExistsParams.ExcludeID != uuid.Nil {
		query = query.Where(merchant.IDNEQ(merchantExistsParams.ExcludeID))
	}
	exist, err = query.Exist(ctx)
	if err != nil {
		err = fmt.Errorf("failed to check merchant existence: %w", err)
		return
	}

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
	if filter.Province != "" {
		query = query.Where(merchant.ProvinceEQ(filter.Province))
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
	address := &domain.Address{
		Country:  em.Country,
		Province: em.Province,
		Address:  em.Address,
		Lng:      em.Lng,
		Lat:      em.Lat,
	}

	repoMerchant := &domain.Merchant{
		ID:                em.ID,
		MerchantCode:      em.MerchantCode,
		MerchantName:      em.MerchantName,
		MerchantShortName: em.MerchantShortName,
		MerchantType:      em.MerchantType,
		BrandName:         em.BrandName,
		AdminPhoneNumber:  em.AdminPhoneNumber,
		ExpireUTC:         em.ExpireUtc,
		BusinessTypeCode:  em.BusinessTypeCode,
		MerchantLogo:      em.MerchantLogo,
		Description:       em.Description,
		Status:            em.Status,
		Address:           address,
		LoginAccount:      em.SuperAccount,
		CreatedAt:         em.CreatedAt,
		UpdatedAt:         em.UpdatedAt,
	}

	if len(em.Edges.MerchantRenewals) > 0 {
		// 倒序返回最后一条
		sort.Slice(em.Edges.MerchantRenewals, func(i, j int) bool {
			return em.Edges.MerchantRenewals[i].CreatedAt.After(em.Edges.MerchantRenewals[j].CreatedAt)
		})
		repoMerchant.PurchaseDuration = em.Edges.MerchantRenewals[0].PurchaseDuration
		repoMerchant.PurchaseDurationUnit = em.Edges.MerchantRenewals[0].PurchaseDurationUnit
	}
	return repoMerchant
}
