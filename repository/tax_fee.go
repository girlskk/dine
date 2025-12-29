package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/taxfee"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

// TaxFeeRepository implements domain.TaxFeeRepository
var _ domain.TaxFeeRepository = (*TaxFeeRepository)(nil)

type TaxFeeRepository struct {
	Client *ent.Client
}

func NewTaxFeeRepository(client *ent.Client) *TaxFeeRepository {
	return &TaxFeeRepository{Client: client}
}

func (repo *TaxFeeRepository) FindByID(ctx context.Context, id uuid.UUID) (fee *domain.TaxFee, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	tf, err := repo.Client.TaxFee.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrTaxFeeNotExists)
			return
		}
		err = fmt.Errorf("failed to query tax fee: %w", err)
		return
	}
	fee = convertTaxFeeToDomain(tf)
	return
}

func (repo *TaxFeeRepository) Create(ctx context.Context, fee *domain.TaxFee) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if fee == nil {
		return fmt.Errorf("tax fee is nil")
	}

	builder := repo.Client.TaxFee.Create().
		SetID(fee.ID).
		SetName(fee.Name).
		SetTaxFeeType(fee.TaxFeeType).
		SetTaxCode(fee.TaxCode).
		SetTaxRateType(fee.TaxRateType).
		SetTaxRate(fee.TaxRate).
		SetDefaultTax(fee.DefaultTax)
	if fee.MerchantID != uuid.Nil {
		builder = builder.SetMerchantID(fee.MerchantID)
	}
	if fee.StoreID != uuid.Nil {
		builder = builder.SetStoreID(fee.StoreID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create tax fee: %w", err)
		return
	}
	fee.ID = created.ID
	fee.CreatedAt = created.CreatedAt
	return
}

func (repo *TaxFeeRepository) Update(ctx context.Context, fee *domain.TaxFee) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeRepository.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if fee == nil {
		return fmt.Errorf("tax fee is nil")
	}

	builder := repo.Client.TaxFee.UpdateOneID(fee.ID).
		SetName(fee.Name).
		SetTaxRateType(fee.TaxRateType).
		SetTaxRate(fee.TaxRate).
		SetDefaultTax(fee.DefaultTax)

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to update tax fee: %w", err)
		return
	}
	fee.UpdatedAt = updated.UpdatedAt
	return
}

func (repo *TaxFeeRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.TaxFee.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to delete tax fee: %w", err)
		return
	}
	return
}

func (repo *TaxFeeRepository) Exists(ctx context.Context, params domain.TaxFeeExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeRepository.Exists")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.TaxFee.Query().
		Where(taxfee.NameEQ(params.Name))
	if params.MerchantID != uuid.Nil {
		query = query.Where(taxfee.MerchantIDEQ(params.MerchantID))
	}
	if params.StoreID != uuid.Nil {
		query = query.Where(taxfee.StoreIDEQ(params.StoreID))
	}
	if params.ExcludeID != uuid.Nil {
		query = query.Where(taxfee.IDNEQ(params.ExcludeID))
	}

	exists, err = query.Exist(ctx)
	if err != nil {
		err = fmt.Errorf("failed to check tax fee exists: %w", err)
		return
	}
	return
}

func (repo *TaxFeeRepository) GetTaxFees(ctx context.Context, pager *upagination.Pagination, filter *domain.TaxFeeListFilter, orderBys ...domain.TaxFeeOrderBy) (fees []*domain.TaxFee, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeRepository.GetTaxFees")
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
		err = fmt.Errorf("failed to count tax fees: %w", err)
		return
	}

	tf, err := query.Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query tax fees: %w", err)
		return
	}
	fees = lo.Map(tf, func(item *ent.TaxFee, _ int) *domain.TaxFee {
		return convertTaxFeeToDomain(item)
	})
	return
}

func (repo *TaxFeeRepository) filterBuildQuery(filter *domain.TaxFeeListFilter) *ent.TaxFeeQuery {
	query := repo.Client.TaxFee.Query()
	if filter.MerchantID != uuid.Nil {
		query = query.Where(taxfee.MerchantIDEQ(filter.MerchantID))
	}
	if filter.StoreID != uuid.Nil {
		query = query.Where(taxfee.StoreIDEQ(filter.StoreID))
	}
	if filter.Name != "" {
		query = query.Where(taxfee.NameContains(filter.Name))
	}
	if filter.TaxFeeType != "" {
		query = query.Where(taxfee.TaxFeeTypeEQ(filter.TaxFeeType))
	}
	return query
}

func (repo *TaxFeeRepository) orderBy(orderBys ...domain.TaxFeeOrderBy) []taxfee.OrderOption {
	var opts []taxfee.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.TaxFeeOrderByID:
			opts = append(opts, taxfee.ByID(rule))
		case domain.TaxFeeOrderByCreatedAt:
			opts = append(opts, taxfee.ByCreatedAt(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, taxfee.ByCreatedAt(sql.OrderDesc()))
	}
	return opts
}

func convertTaxFeeToDomain(tf *ent.TaxFee) *domain.TaxFee {
	return &domain.TaxFee{
		ID:          tf.ID,
		Name:        tf.Name,
		TaxFeeType:  tf.TaxFeeType,
		TaxCode:     tf.TaxCode,
		TaxRateType: tf.TaxRateType,
		TaxRate:     tf.TaxRate,
		DefaultTax:  tf.DefaultTax,
		MerchantID:  tf.MerchantID,
		StoreID:     tf.StoreID,
		CreatedAt:   tf.CreatedAt,
		UpdatedAt:   tf.UpdatedAt,
	}
}
