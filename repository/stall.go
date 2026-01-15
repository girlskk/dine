package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/stall"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StallRepository = (*StallRepository)(nil)

// StallRepository implements Stall CRUD and pagination.
type StallRepository struct {
	Client *ent.Client
}

func NewStallRepository(client *ent.Client) *StallRepository {
	return &StallRepository{Client: client}
}

func (repo *StallRepository) FindByID(ctx context.Context, id uuid.UUID) (domainStall *domain.Stall, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StallRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	es, err := repo.Client.Stall.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrStallNotExists)
			return
		}
		return
	}
	domainStall = convertStallToDomain(es)
	return
}

func (repo *StallRepository) Create(ctx context.Context, domainStall *domain.Stall) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StallRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainStall == nil {
		return fmt.Errorf("stall is nil")
	}

	builder := repo.Client.Stall.Create().SetID(domainStall.ID).
		SetName(domainStall.Name).
		SetStallType(domainStall.StallType).
		SetPrintType(domainStall.PrintType).
		SetEnabled(domainStall.Enabled).
		SetSortOrder(domainStall.SortOrder)

	if domainStall.MerchantID != uuid.Nil {
		builder = builder.SetMerchantID(domainStall.MerchantID)
	}
	if domainStall.StoreID != uuid.Nil {
		builder = builder.SetStoreID(domainStall.StoreID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create stall: %w", err)
		return
	}
	domainStall.ID = created.ID
	domainStall.CreatedAt = created.CreatedAt
	domainStall.UpdatedAt = created.UpdatedAt
	return
}

func (repo *StallRepository) Update(ctx context.Context, domainStall *domain.Stall) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StallRepository.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainStall == nil {
		return fmt.Errorf("stall is nil")
	}

	builder := repo.Client.Stall.UpdateOneID(domainStall.ID).
		SetName(domainStall.Name).
		SetPrintType(domainStall.PrintType).
		SetEnabled(domainStall.Enabled).
		SetSortOrder(domainStall.SortOrder)

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrStallNotExists)
			return
		}
		err = fmt.Errorf("failed to update stall: %w", err)
		return
	}
	domainStall.UpdatedAt = updated.UpdatedAt
	return
}

func (repo *StallRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StallRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.Stall.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to delete stall: %w", err)
		return
	}
	return nil
}

func (repo *StallRepository) GetStalls(ctx context.Context, pager *upagination.Pagination, filter *domain.StallListFilter, orderBys ...domain.StallOrderBy) (domainStalls []*domain.Stall, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StallRepository.GetStalls")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.buildFilterQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count stall: %w", err)
		return
	}

	stalls, err := query.
		Order(repo.orderBy(orderBys...)...).
		WithCategories().
		WithProducts().
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query stall: %w", err)
		return
	}

	domainStalls = lo.Map(stalls, func(item *ent.Stall, _ int) *domain.Stall {
		return convertStallToDomain(item)
	})
	return
}

func (repo *StallRepository) Exists(ctx context.Context, params domain.StallExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StallRepository.Exists")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.Stall.Query()
	if params.MerchantID != uuid.Nil {
		query = query.Where(stall.MerchantID(params.MerchantID))
	}
	if params.StoreID != uuid.Nil {
		query = query.Where(stall.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query = query.Where(stall.Name(params.Name))
	}
	if params.ExcludeID != uuid.Nil {
		query = query.Where(stall.IDNEQ(params.ExcludeID))
	}

	exists, err = query.Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check stall existence: %w", err)
	}
	return
}

func (repo *StallRepository) buildFilterQuery(filter *domain.StallListFilter) *ent.StallQuery {
	query := repo.Client.Stall.Query()
	if filter == nil {
		return query
	}

	if filter.MerchantID != uuid.Nil {
		query = query.Where(stall.MerchantID(filter.MerchantID))
	}
	if filter.StoreID != uuid.Nil {
		query = query.Where(stall.StoreID(filter.StoreID))
	}
	if filter.StallType != "" {
		query = query.Where(stall.StallTypeEQ(filter.StallType))
	}
	if filter.PrintType != "" {
		query = query.Where(stall.PrintTypeEQ(filter.PrintType))
	}
	if filter.Enabled != nil {
		query = query.Where(stall.EnabledEQ(*filter.Enabled))
	}
	if filter.Name != "" {
		query = query.Where(stall.NameContains(filter.Name))
	}

	return query
}

func (repo *StallRepository) orderBy(orderBys ...domain.StallOrderBy) []stall.OrderOption {
	var opts []stall.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.StallOrderByID:
			opts = append(opts, stall.ByID(rule))
		case domain.StallOrderByCreatedAt:
			opts = append(opts, stall.ByCreatedAt(rule))
		case domain.StallOrderBySortOrder:
			opts = append(opts, stall.BySortOrder(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, stall.ByCreatedAt(sql.OrderDesc()))
	}
	return opts
}

func convertStallToDomain(es *ent.Stall) *domain.Stall {
	if es == nil {
		return nil
	}
	ds := &domain.Stall{
		ID:         es.ID,
		Name:       es.Name,
		StallType:  es.StallType,
		PrintType:  es.PrintType,
		Enabled:    es.Enabled,
		SortOrder:  es.SortOrder,
		MerchantID: es.MerchantID,
		StoreID:    es.StoreID,
		CreatedAt:  es.CreatedAt,
		UpdatedAt:  es.UpdatedAt,
	}
	if len(es.Edges.Products) > 0 {
		ds.RelationProductNumbers += len(es.Edges.Products)
	}
	if len(es.Edges.Categories) > 0 {
		ds.RelationProductNumbers += lo.SumBy(es.Edges.Categories, func(c *ent.Category) int {
			return c.ProductCount
		})
	}
	return ds
}
