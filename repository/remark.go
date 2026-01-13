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

func (repo *RemarkRepository) FindByID(ctx context.Context, id uuid.UUID) (domainRemark *domain.Remark, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	er, err := repo.Client.Remark.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrRemarkNotExists)
			return
		}
		err = fmt.Errorf("failed to query Remark: %w", err)
		return
	}
	domainRemark = convertRemarkToDomain(er)
	return domainRemark, nil
}

func (repo *RemarkRepository) Create(ctx context.Context, remark *domain.Remark) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if remark == nil {
		return fmt.Errorf("remark is nil")
	}

	builder := repo.Client.Remark.Create().
		SetID(remark.ID).
		SetName(remark.Name).
		SetRemarkType(remark.RemarkType).
		SetEnabled(remark.Enabled).
		SetSortOrder(remark.SortOrder).
		SetRemarkScene(remark.RemarkScene)

	if remark.MerchantID == uuid.Nil && remark.RemarkType == domain.RemarkTypeBrand {
		return fmt.Errorf("merchant ID is required for brand remark")
	}
	if remark.MerchantID != uuid.Nil {
		builder = builder.SetMerchantID(remark.MerchantID)
	}
	if remark.StoreID != uuid.Nil {
		builder = builder.SetStoreID(remark.StoreID)
	}
	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create remark: %w", err)
		return
	}
	remark.ID = created.ID
	remark.CreatedAt = created.CreatedAt
	return
}

func (repo *RemarkRepository) Update(ctx context.Context, remark *domain.Remark) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if remark == nil {
		return fmt.Errorf("remark is nil")
	}

	builder := repo.Client.Remark.UpdateOneID(remark.ID).
		SetName(remark.Name).
		SetEnabled(remark.Enabled).
		SetSortOrder(remark.SortOrder)
	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to update remark: %w", err)
		return err
	}

	remark.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *RemarkRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.Remark.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to delete remark: %w", err)
		return
	}
	return
}

func (repo *RemarkRepository) GetRemarks(ctx context.Context, pager *upagination.Pagination, filter *domain.RemarkListFilter, orderBys ...domain.RemarkOrderBy) (domainRemarks domain.Remarks, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.GetRemarks")
	defer func() { util.SpanErrFinish(span, err) }()

	if pager == nil {
		err = fmt.Errorf("pager is nil")
		return
	}
	if filter == nil {
		err = fmt.Errorf("filter is nil")
		return
	}
	if filter.StoreID != uuid.Nil && filter.MerchantID == uuid.Nil {
		err = fmt.Errorf("merchant ID is required when store ID is provided")
		return
	}
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

func (repo *RemarkRepository) CountRemarkByScene(ctx context.Context, params domain.CountRemarkParams) (countRemark map[domain.RemarkScene]int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.CountRemarkByScene")
	defer func() { util.SpanErrFinish(span, err) }()

	if params.StoreID != uuid.Nil && params.MerchantID == uuid.Nil {
		err = fmt.Errorf("merchant ID is required when store ID is provided")
		return
	}

	countRemark = make(map[domain.RemarkScene]int)
	if len(params.RemarkScenes) == 0 {
		return countRemark, nil
	}

	type result struct {
		RemarkScene domain.RemarkScene `json:"remark_scene"`
		Count       int                `json:"count"`
	}
	var results []result
	builder := repo.Client.Remark.Query().
		Where(remark.RemarkSceneIn(params.RemarkScenes...))
	if params.RemarkType != "" {
		if params.MerchantID != uuid.Nil && params.StoreID != uuid.Nil {
		}
		if params.MerchantID != uuid.Nil && params.StoreID == uuid.Nil {
		}
		if params.MerchantID == uuid.Nil && params.StoreID == uuid.Nil {
			builder = builder.Where(remark.RemarkTypeEQ(params.RemarkType))
		}
	}

	builder = repo.convertMerchantIDFilter(params.RemarkType, params.MerchantID, builder)
	builder = repo.convertStoreIDFilter(params.RemarkType, params.StoreID, builder)

	err = builder.
		GroupBy(remark.FieldRemarkScene).
		Aggregate(ent.Count()).
		Scan(ctx, &results)
	if err != nil {
		err = fmt.Errorf("failed to count remarks by scenes: %w", err)
		return
	}
	countRemark = lo.SliceToMap(results, func(item result) (domain.RemarkScene, int) {
		return item.RemarkScene, item.Count
	})
	return
}

func (repo *RemarkRepository) Exists(ctx context.Context, params domain.RemarkExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkRepository.Exists")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.Remark.Query()
	if params.RemarkType != "" {
		query = query.Where(remark.RemarkTypeEQ(params.RemarkType))
	}
	if params.RemarkScene != "" {
		query = query.Where(remark.RemarkSceneEQ(params.RemarkScene))
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
		ID:          er.ID,
		Name:        er.Name,
		RemarkType:  er.RemarkType,
		Enabled:     er.Enabled,
		SortOrder:   er.SortOrder,
		RemarkScene: er.RemarkScene,
		MerchantID:  er.MerchantID,
		StoreID:     er.StoreID,
		CreatedAt:   er.CreatedAt,
		UpdatedAt:   er.UpdatedAt,
	}
}

func (repo *RemarkRepository) filterBuildQuery(filter *domain.RemarkListFilter) *ent.RemarkQuery {
	query := repo.Client.Remark.Query()

	query = repo.convertMerchantIDFilter(filter.RemarkType, filter.MerchantID, query)
	query = repo.convertStoreIDFilter(filter.RemarkType, filter.StoreID, query)

	if filter.RemarkScene != "" {
		query = query.Where(remark.RemarkSceneEQ(filter.RemarkScene))
	}
	if filter.Enabled != nil {
		query = query.Where(remark.EnabledEQ(*filter.Enabled))
	}
	if filter.RemarkType != "" {
		if filter.MerchantID != uuid.Nil && filter.StoreID != uuid.Nil {
		}
		if filter.MerchantID != uuid.Nil && filter.StoreID == uuid.Nil {
		}
		if filter.MerchantID == uuid.Nil && filter.StoreID == uuid.Nil {
			query = query.Where(remark.RemarkTypeEQ(filter.RemarkType))
		}
	}
	return query
}

// 根据merchant ID的查询时 可查询出系统的备注
// 系统备注只查询系统级别的备注
// 品牌备注查询品牌和系统级别的备注
// 门店备注查询门店所属品牌和系统级别的备注
// Remark Type为空时只查询当前商户的备注
// convertMerchantIDFilter 系统备注只查询系统级别的备注
func (repo *RemarkRepository) convertMerchantIDFilter(remarkType domain.RemarkType, merchantID uuid.UUID, query *ent.RemarkQuery) *ent.RemarkQuery {
	if merchantID != uuid.Nil {
		switch remarkType {
		case domain.RemarkTypeSystem: // 系统备注只查询系统级别的备注
			query = query.Where(remark.Or(remark.MerchantIDIsNil(), remark.MerchantID(uuid.Nil)))
			query = query.Where(remark.RemarkTypeEQ(domain.RemarkTypeSystem))
		case domain.RemarkTypeBrand: // 品牌备注查询品牌和系统级别的备注
			query = query.Where(remark.Or(remark.MerchantID(merchantID), remark.MerchantIDIsNil(), remark.MerchantID(uuid.Nil)))
			query = query.Where(remark.RemarkTypeIn(domain.RemarkTypeSystem, domain.RemarkTypeBrand))
		case domain.RemarkTypeStore: // 门店备注查询门店所属品牌和系统级别的备注
			query = query.Where(remark.Or(remark.MerchantID(merchantID), remark.MerchantIDIsNil(), remark.MerchantID(uuid.Nil)))
		default:
			// Remark Type为空时只查询当前商户的备注
			query = query.Where(remark.MerchantID(merchantID))

		}
	}
	return query
}

// 根据store ID的查询时 可查询出系统和品牌的备注
// 系统备注只查询系统级别的备注
// 品牌备注查询品牌和系统级别的备注
// 门店备注查询门店和系统级别的备注
// Remark Type为空时只查询当前门店的备注
// convertStoreIDFilter 系统备注只查询系统级别的备注
func (repo *RemarkRepository) convertStoreIDFilter(remarkType domain.RemarkType, storeID uuid.UUID, query *ent.RemarkQuery) *ent.RemarkQuery {
	if storeID != uuid.Nil {
		switch remarkType {
		case domain.RemarkTypeSystem: // 系统备注只查询系统级别的备注
			query = query.Where(remark.Or(remark.StoreIDIsNil(), remark.StoreID(uuid.Nil)))
			query = query.Where(remark.RemarkTypeEQ(domain.RemarkTypeSystem))
		case domain.RemarkTypeBrand: // 品牌备注查询品牌和系统级别的备注
			query = query.Where(remark.Or(remark.StoreIDIsNil(), remark.StoreID(uuid.Nil)))
			query = query.Where(remark.RemarkTypeIn(domain.RemarkTypeSystem, domain.RemarkTypeBrand))
		case domain.RemarkTypeStore: // 门店备注查询门店和系统级别的备注
			query = query.Where(remark.Or(remark.StoreID(storeID), remark.StoreIDIsNil(), remark.StoreID(uuid.Nil)))
		default: // Remark Type为空时只查询当前门店的备注
			query = query.Where(remark.StoreID(storeID))
		}
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
