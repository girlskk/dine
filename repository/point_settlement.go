package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/pointsettlement"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.PointSettlementRepository = (*PointSettlementRepository)(nil)

type PointSettlementRepository struct {
	Client *ent.Client
}

func NewPointSettlementRepository(client *ent.Client) *PointSettlementRepository {
	return &PointSettlementRepository{
		Client: client,
	}
}

func (r *PointSettlementRepository) BatchCreate(ctx context.Context, settlements domain.PointSettlements) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementRepository.BatchCreate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	settlementCreates := lo.Map(settlements, func(item *domain.PointSettlement, _ int) *ent.PointSettlementCreate {
		return r.Client.PointSettlement.Create().
			SetNo(item.No).
			SetStoreID(item.StoreID).
			SetStoreName(item.StoreName).
			SetOrderCount(item.OrderCount).
			SetAmount(item.Amount).
			SetTotalPoints(item.TotalPoints).
			SetDate(item.Date).
			SetStatus(int(item.Status)).
			SetPointSettlementRate(item.PointSettlementRate)
	})

	_, err = r.Client.PointSettlement.CreateBulk(settlementCreates...).Save(ctx)
	return
}

func (r *PointSettlementRepository) PagedListBySearch(ctx context.Context,
	page *upagination.Pagination, params domain.PointSettlementSearchParams,
) (res *domain.PointSettlementSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := r.filterBuildQuery(params)

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	settlements, err := query.Order(ent.Desc(pointsettlement.FieldID)).
		Offset(page.Offset()).
		Limit(page.Size).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.PointSettlements, 0, len(settlements))
	for _, c := range settlements {
		items = append(items, convertToPointSettlement(c))
	}

	return &domain.PointSettlementSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

func convertToPointSettlement(eu *ent.PointSettlement) *domain.PointSettlement {
	return &domain.PointSettlement{
		ID:                  eu.ID,
		No:                  eu.No,
		StoreID:             eu.StoreID,
		StoreName:           eu.StoreName,
		OrderCount:          eu.OrderCount,
		Amount:              eu.Amount,
		TotalPoints:         eu.TotalPoints,
		Date:                eu.Date,
		Status:              domain.PointSettlementStatus(eu.Status),
		PointSettlementRate: eu.PointSettlementRate,
		ApprovedAt:          eu.ApprovedAt,
		ApproverID:          eu.ApproverID,
		CreatedAt:           eu.CreatedAt,
		UpdatedAt:           eu.UpdatedAt,
	}
}

func (r *PointSettlementRepository) FindByID(ctx context.Context, id int) (res *domain.PointSettlement, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	settlement, err := r.Client.PointSettlement.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrPointSettlementNotExists)
		}
		return nil, err
	}
	return convertToPointSettlement(settlement), nil
}

func (r *PointSettlementRepository) FindByIDForUpdate(ctx context.Context, id int) (res *domain.PointSettlement, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementRepository.FindByIDForUpdate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	settlement, err := r.Client.PointSettlement.Query().
		Where(pointsettlement.ID(id)).
		ForUpdate().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrPointSettlementNotExists)
		}
		return nil, err
	}
	return convertToPointSettlement(settlement), nil
}

func (r *PointSettlementRepository) UpdateStatus(ctx context.Context, id int, status domain.PointSettlementStatus, approverID *int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementRepository.UpdateStatus")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	updateQuery := r.Client.PointSettlement.
		UpdateOneID(id).
		SetStatus(int(status))

	// 如果是审批通过，设置审批人和审批时间
	if approverID != nil {
		updateQuery.SetApproverID(*approverID).
			SetApprovedAt(time.Now())
	} else {
		// 反审批，清除审批人和审批时间
		updateQuery.ClearApproverID().
			ClearApprovedAt()
	}
	_, err = updateQuery.Save(ctx)
	return err
}

func (r *PointSettlementRepository) GetPointSettlementRange(ctx context.Context, params domain.PointSettlementSearchParams,
) (res domain.PointSettlementRange, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementRepository.GetPointSettlementRange")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := r.filterBuildQuery(params)

	var v []struct {
		Min, Max, Count int
	}

	err = query.Aggregate(
		ent.Min(pointsettlement.FieldID),
		ent.Max(pointsettlement.FieldID),
		ent.Count(),
	).Scan(ctx, &v)

	if err != nil {
		err = fmt.Errorf("failed to get max point settlement id: %w", err)
		return
	}

	if len(v) > 0 {
		res = domain.PointSettlementRange{
			MinID: v[0].Min,
			MaxID: v[0].Max,
			Count: v[0].Count,
		}
	}

	return
}

func (r *PointSettlementRepository) filterBuildQuery(params domain.PointSettlementSearchParams) *ent.PointSettlementQuery {
	query := r.Client.PointSettlement.Query()

	if params.StoreID > 0 {
		query.Where(pointsettlement.StoreID(params.StoreID))
	}

	if params.StartAt != nil {
		query.Where(pointsettlement.DateGTE(*params.StartAt))
	}
	if params.EndAt != nil {
		query.Where(pointsettlement.DateLTE(*params.EndAt))
	}

	if params.IDGte > 0 {
		query = query.Where(pointsettlement.IDGTE(params.IDGte))
	}
	if params.IDLte > 0 {
		query = query.Where(pointsettlement.IDLTE(params.IDLte))
	}

	return query
}
