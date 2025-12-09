package repository

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/reconciliationrecord"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ReconciliationRecordRepository = (*ReconciliationRecordRepository)(nil)

type ReconciliationRecordRepository struct {
	Client *ent.Client
}

func NewReconciliationRecordRepository(client *ent.Client) *ReconciliationRecordRepository {
	return &ReconciliationRecordRepository{
		Client: client,
	}
}

func (r *ReconciliationRecordRepository) BatchCreate(ctx context.Context, records domain.ReconciliationRecords) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationRecordRepository.BatchCreate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	recordCreates := lo.Map(records, func(item *domain.ReconciliationRecord, _ int) *ent.ReconciliationRecordCreate {
		return r.Client.ReconciliationRecord.Create().
			SetNo(item.No).
			SetStoreID(item.StoreID).
			SetStoreName(item.StoreName).
			SetOrderCount(item.OrderCount).
			SetAmount(item.Amount).
			SetChannel(item.Channel).
			SetDate(item.Date)
	})

	_, err = r.Client.ReconciliationRecord.CreateBulk(recordCreates...).Save(ctx)
	return
}

func (r *ReconciliationRecordRepository) PagedListBySearch(ctx context.Context,
	page *upagination.Pagination, params domain.ReconciliationSearchParams,
) (res *domain.ReconciliationSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationRecordRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := r.filterBuildQuery(params)

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	records, err := query.Order(ent.Desc(reconciliationrecord.FieldID)).
		Offset(page.Offset()).
		Limit(page.Size).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.ReconciliationRecords, 0, len(records))
	for _, c := range records {
		items = append(items, convertToReconciliationRecord(c))
	}

	return &domain.ReconciliationSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

func convertToReconciliationRecord(eu *ent.ReconciliationRecord) *domain.ReconciliationRecord {
	return &domain.ReconciliationRecord{
		ID:         eu.ID,
		No:         eu.No,
		StoreID:    eu.StoreID,
		StoreName:  eu.StoreName,
		OrderCount: eu.OrderCount,
		Amount:     eu.Amount,
		Channel:    eu.Channel,
		Date:       eu.Date,
		CreatedAt:  eu.CreatedAt,
		UpdatedAt:  eu.UpdatedAt,
	}
}

func (r *ReconciliationRecordRepository) FindByID(ctx context.Context, id int) (res *domain.ReconciliationRecord, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationRecordRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	record, err := r.Client.ReconciliationRecord.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrReconciliationRecordNotExists)
		}
		return nil, err
	}
	return convertToReconciliationRecord(record), nil
}

func (r *ReconciliationRecordRepository) GetReconciliationRange(ctx context.Context, params domain.ReconciliationSearchParams,
) (res domain.ReconciliationRange, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationRecordRepository.GetReconciliationRange")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := r.filterBuildQuery(params)

	var v []struct {
		Min, Max, Count int
	}

	err = query.Aggregate(
		ent.Min(reconciliationrecord.FieldID),
		ent.Max(reconciliationrecord.FieldID),
		ent.Count(),
	).Scan(ctx, &v)

	if err != nil {
		err = fmt.Errorf("failed to get max reconciliation record id: %w", err)
		return
	}

	if len(v) > 0 {
		res = domain.ReconciliationRange{
			MinID: v[0].Min,
			MaxID: v[0].Max,
			Count: v[0].Count,
		}
	}

	return
}

func (r *ReconciliationRecordRepository) filterBuildQuery(params domain.ReconciliationSearchParams) *ent.ReconciliationRecordQuery {
	query := r.Client.ReconciliationRecord.Query()

	if params.StoreID > 0 {
		query.Where(reconciliationrecord.StoreID(params.StoreID))
	}

	if params.StartAt != nil {
		query.Where(reconciliationrecord.DateGTE(*params.StartAt))
	}
	if params.EndAt != nil {
		query.Where(reconciliationrecord.DateLTE(*params.EndAt))
	}

	if params.Channel != "" {
		query.Where(reconciliationrecord.ChannelEQ(params.Channel))
	}

	if params.IDGte > 0 {
		query = query.Where(reconciliationrecord.IDGTE(params.IDGte))
	}
	if params.IDLte > 0 {
		query = query.Where(reconciliationrecord.IDLTE(params.IDLte))
	}

	return query
}
