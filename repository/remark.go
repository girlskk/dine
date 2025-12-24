package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/remark"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RemarkRepository = (*RemarkRepository)(nil)

type RemarkRepository struct {
	Client *ent.Client
}

func NewRemarkRepository(client *ent.Client) *RemarkRepository {
	return &RemarkRepository{Client: client}
}

func (repo *RemarkRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.Remark, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	er, err := repo.Client.Remark.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrRemarkNotExists)
		}
		return nil, err
	}
	res = convertRemarkToDomain(er)
	return res, nil
}

func (repo *RemarkRepository) Create(ctx context.Context, remark *domain.Remark) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if remark == nil {
		return fmt.Errorf("remark is nil")
	}

	_, err = repo.Client.Remark.Create().
		SetID(remark.ID).
		SetName(remark.Name).
		SetRemarkType(remark.RemarkType).
		SetEnabled(remark.Enabled).
		SetSortOrder(remark.SortOrder).
		SetCategoryID(remark.CategoryID).
		SetMerchantID(remark.MerchantID).
		SetStoreID(remark.StoreID).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create remark: %w", err)
		return err
	}
	return
}

func (repo *RemarkRepository) Update(ctx context.Context, remark *domain.Remark) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if remark == nil {
		return fmt.Errorf("remark is nil")
	}

	_, err = repo.Client.Remark.UpdateOneID(remark.ID).
		SetName(remark.Name).
		SetRemarkType(remark.RemarkType).
		SetEnabled(remark.Enabled).
		SetSortOrder(remark.SortOrder).
		SetCategoryID(remark.CategoryID).
		SetMerchantID(remark.MerchantID).
		SetStoreID(remark.StoreID).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to update remark: %w", err)
		return err
	}

	return nil
}

func (repo *RemarkRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.Remark.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to delete remark: %w", err)
		return
	}
	return
}

func (repo *RemarkRepository) GetRemarks(ctx context.Context, pager *upagination.Pagination, filter *domain.RemarkListFilter, orderBys ...domain.RemarkOrderBy) (domainRemarks domain.Remarks, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.GetRemarks")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.filterBuildQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count: %w", err)
		return
	}
	remarks, err := query.Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query remark: %w", err)
		return
	}
	domainRemarks = lo.Map(remarks, func(item *ent.Remark, _ int) *domain.Remark {
		return convertRemarkToDomain(item)
	})
	return
}

func (repo *RemarkRepository) Exists(ctx context.Context, params domain.RemarkExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.Exists")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.Remark.Query()
	if params.CategoryID != uuid.Nil {
		query = query.Where(remark.CategoryID(params.CategoryID))
	}
	if params.MerchantID != uuid.Nil {
		query = query.Where(remark.MerchantID(params.MerchantID))
	}
	if params.StoreID != uuid.Nil {
		query = query.Where(remark.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query = query.Where(remark.Name(params.Name))
	}
	if params.ExcludeID != uuid.Nil {
		query = query.Where(remark.IDNEQ(params.ExcludeID))
	}
	exists, err = query.Exist(ctx)
	if err != nil {
		err = fmt.Errorf("failed to check remark existence: %w", err)
		return
	}
	return
}

func convertRemarkToDomain(er *ent.Remark) *domain.Remark {
	if er == nil {
		return nil
	}
	return &domain.Remark{
		ID:         er.ID,
		Name:       er.Name,
		RemarkType: er.RemarkType,
		Enabled:    er.Enabled,
		SortOrder:  er.SortOrder,
		CategoryID: er.CategoryID,
		MerchantID: er.MerchantID,
		StoreID:    er.StoreID,
		CreatedAt:  er.CreatedAt,
		UpdatedAt:  er.UpdatedAt,
	}
}

func (repo *RemarkRepository) filterBuildQuery(filter *domain.RemarkListFilter) *ent.RemarkQuery {
	query := repo.Client.Remark.Query()

	if filter.CategoryID != uuid.Nil {
		query = query.Where(remark.CategoryID(filter.CategoryID))
	}
	if filter.MerchantID != uuid.Nil {
		query = query.Where(remark.MerchantIDIn(filter.MerchantID, uuid.Nil)) // include both brand and system remarks
	}
	if filter.StoreID != uuid.Nil {
		query = query.Where(remark.StoreID(filter.StoreID))
	}
	if filter.Enabled != nil {
		query = query.Where(remark.EnabledEQ(*filter.Enabled))
	}
	if filter.RemarkType != "" {
		query = query.Where(remark.RemarkTypeEQ(filter.RemarkType))
	}
	return query
}

func (repo *RemarkRepository) orderBy(orderBys ...domain.RemarkOrderBy) []remark.OrderOption {
	var opts []remark.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.RemarkOrderByID:
			opts = append(opts, remark.ByID(rule))
		case domain.RemarkOrderByCreatedAt:
			opts = append(opts, remark.ByCreatedAt(rule))
		case domain.RemarkOrderBySortOrder:
			opts = append(opts, remark.BySortOrder(rule))

		}
	}

	if len(opts) == 0 {
		opts = append(opts, remark.ByCreatedAt(sql.OrderDesc()))
	}

	return opts
}
