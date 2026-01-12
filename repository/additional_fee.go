package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/additionalfee"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

// AdditionalFeeRepository 实现附加费仓储
var _ domain.AdditionalFeeRepository = (*AdditionalFeeRepository)(nil)

type AdditionalFeeRepository struct {
	Client *ent.Client
}

func NewAdditionalFeeRepository(client *ent.Client) *AdditionalFeeRepository {
	return &AdditionalFeeRepository{Client: client}
}

func (repo *AdditionalFeeRepository) FindByID(ctx context.Context, id uuid.UUID) (fee *domain.AdditionalFee, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	ef, err := repo.Client.AdditionalFee.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrAdditionalFeeNotExists)
			return
		}
		err = fmt.Errorf("failed to query AdditionalFee: %w", err)
		return
	}
	fee = convertAdditionalFeeToDomain(ef)
	return
}

func (repo *AdditionalFeeRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) (fees []*domain.AdditionalFee, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeRepository.ListByIDs")
	defer func() { util.SpanErrFinish(span, err) }()

	efs, err := repo.Client.AdditionalFee.Query().Where(additionalfee.IDIn(ids...)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query additional fees: %w", err)
	}
	fees = lo.Map(efs, func(item *ent.AdditionalFee, _ int) *domain.AdditionalFee {
		return convertAdditionalFeeToDomain(item)
	})
	return fees, nil
}

func (repo *AdditionalFeeRepository) Create(ctx context.Context, fee *domain.AdditionalFee) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if fee == nil {
		return fmt.Errorf("additional fee is nil")
	}

	builder := repo.Client.AdditionalFee.Create().
		SetID(fee.ID).
		SetName(fee.Name).
		SetFeeType(fee.FeeType).
		SetFeeCategory(fee.FeeCategory).
		SetChargeMode(fee.ChargeMode).
		SetFeeValue(fee.FeeValue).
		SetIncludeInReceivable(fee.IncludeInReceivable).
		SetTaxable(fee.Taxable).
		SetDiscountScope(fee.DiscountScope).
		SetOrderChannels(fee.OrderChannels).
		SetDiningWays(fee.DiningWays).
		SetEnabled(fee.Enabled).
		SetSortOrder(fee.SortOrder)
	if fee.MerchantID != uuid.Nil {
		builder = builder.SetMerchantID(fee.MerchantID)
	}
	if fee.StoreID != uuid.Nil {
		builder = builder.SetStoreID(fee.StoreID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create additional fee: %w", err)
		return
	}
	fee.ID = created.ID
	fee.CreatedAt = created.CreatedAt
	return
}

func (repo *AdditionalFeeRepository) Update(ctx context.Context, fee *domain.AdditionalFee) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeRepository.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if fee == nil {
		return fmt.Errorf("additional fee is nil")
	}

	builder := repo.Client.AdditionalFee.UpdateOneID(fee.ID).
		SetName(fee.Name).
		SetFeeCategory(fee.FeeCategory).
		SetChargeMode(fee.ChargeMode).
		SetFeeValue(fee.FeeValue).
		SetIncludeInReceivable(fee.IncludeInReceivable).
		SetTaxable(fee.Taxable).
		SetDiscountScope(fee.DiscountScope).
		SetOrderChannels(fee.OrderChannels).
		SetDiningWays(fee.DiningWays).
		SetEnabled(fee.Enabled).
		SetSortOrder(fee.SortOrder)

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to update additional fee: %w", err)
		return
	}
	fee.UpdatedAt = updated.UpdatedAt
	return
}

func (repo *AdditionalFeeRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.AdditionalFee.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to delete additional fee: %w", err)
		return
	}
	return
}

func (repo *AdditionalFeeRepository) Exists(ctx context.Context, params domain.AdditionalFeeExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeRepository.Exists")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.AdditionalFee.Query().
		Where(additionalfee.NameEQ(params.Name))
	if params.MerchantID != uuid.Nil {
		query = query.Where(additionalfee.MerchantIDEQ(params.MerchantID))
	}
	if params.StoreID != uuid.Nil {
		query = query.Where(additionalfee.StoreIDEQ(params.StoreID))
	}
	if params.ExcludeID != uuid.Nil {
		query = query.Where(additionalfee.IDNEQ(params.ExcludeID))
	}
	if params.FeeType != "" {
		query = query.Where(additionalfee.FeeTypeEQ(params.FeeType))
	}
	exists, err = query.Exist(ctx)
	if err != nil {
		err = fmt.Errorf("failed to check additional fee exists: %w", err)
		return
	}
	return
}

func (repo *AdditionalFeeRepository) GetAdditionalFees(ctx context.Context, pager *upagination.Pagination, filter *domain.AdditionalFeeListFilter, orderBys ...domain.AdditionalFeeOrderBy) (fees []*domain.AdditionalFee, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeRepository.GetAdditionalFees")
	defer func() { util.SpanErrFinish(span, err) }()

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

	ef, err := query.Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query additional fee: %w", err)
		return
	}
	fees = lo.Map(ef, func(item *ent.AdditionalFee, _ int) *domain.AdditionalFee {
		return convertAdditionalFeeToDomain(item)
	})
	return
}

func (repo *AdditionalFeeRepository) filterBuildQuery(filter *domain.AdditionalFeeListFilter) *ent.AdditionalFeeQuery {
	query := repo.Client.AdditionalFee.Query()
	if filter.MerchantID != uuid.Nil {
		query = query.Where(additionalfee.MerchantIDEQ(filter.MerchantID))
	}
	if filter.StoreID != uuid.Nil {
		query = query.Where(additionalfee.StoreIDEQ(filter.StoreID))
	}
	if filter.Name != "" {
		query = query.Where(additionalfee.NameContains(filter.Name))
	}
	if filter.FeeType != "" {
		query = query.Where(additionalfee.FeeTypeEQ(filter.FeeType))
	}
	if filter.Enabled != nil {
		query = query.Where(additionalfee.EnabledEQ(*filter.Enabled))
	}
	return query
}

func (repo *AdditionalFeeRepository) orderBy(orderBys ...domain.AdditionalFeeOrderBy) []additionalfee.OrderOption {
	var opts []additionalfee.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.AdditionalFeeOrderByID:
			opts = append(opts, additionalfee.ByID(rule))
		case domain.AdditionalFeeOrderByCreatedAt:
			opts = append(opts, additionalfee.ByCreatedAt(rule))
		case domain.AdditionalFeeOrderBySortOrder:
			opts = append(opts, additionalfee.BySortOrder(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, additionalfee.ByCreatedAt(sql.OrderDesc()))
	}
	return opts
}

func convertAdditionalFeeToDomain(ef *ent.AdditionalFee) *domain.AdditionalFee {
	return &domain.AdditionalFee{
		ID:                  ef.ID,
		Name:                ef.Name,
		FeeType:             ef.FeeType,
		ChargeMode:          ef.ChargeMode,
		FeeValue:            ef.FeeValue,
		IncludeInReceivable: ef.IncludeInReceivable,
		Taxable:             ef.Taxable,
		DiscountScope:       ef.DiscountScope,
		OrderChannels:       ef.OrderChannels,
		DiningWays:          ef.DiningWays,
		Enabled:             ef.Enabled,
		SortOrder:           ef.SortOrder,
		MerchantID:          ef.MerchantID,
		StoreID:             ef.StoreID,
		CreatedAt:           ef.CreatedAt,
		UpdatedAt:           ef.UpdatedAt,
	}
}
