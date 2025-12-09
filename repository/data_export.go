package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/dataexport"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.DataExportRepository = (*DataExportRepository)(nil)

type DataExportRepository struct {
	Client *ent.Client
}

func NewDataExportRepository(client *ent.Client) *DataExportRepository {
	return &DataExportRepository{
		Client: client,
	}
}

func (r *DataExportRepository) Create(ctx context.Context, dataExport *domain.DataExport) (newDataExport *domain.DataExport, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DataExportRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	de, err := r.Client.DataExport.Create().
		SetStoreID(dataExport.StoreID).
		SetType(dataExport.Type).
		SetStatus(dataExport.Status).
		SetParams(dataExport.Params).
		SetOperatorType(dataExport.OperatorType).
		SetOperatorID(dataExport.OperatorID).
		SetOperatorName(dataExport.OperatorName).
		SetFileName(dataExport.FileName).
		SetURL(dataExport.URL).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create data export: %w", err)
		return
	}

	newDataExport = convertDataExport(de)

	return
}

func (r *DataExportRepository) CreateBulk(ctx context.Context, dataExports []*domain.DataExport) (newDataExports []*domain.DataExport, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DataExportRepository.CreateBulk")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	des, err := r.Client.DataExport.MapCreateBulk(dataExports, func(c *ent.DataExportCreate, idx int) {
		c.SetStoreID(dataExports[idx].StoreID).
			SetType(dataExports[idx].Type).
			SetStatus(dataExports[idx].Status).
			SetParams(dataExports[idx].Params).
			SetOperatorType(dataExports[idx].OperatorType).
			SetOperatorID(dataExports[idx].OperatorID).
			SetOperatorName(dataExports[idx].OperatorName).
			SetFileName(dataExports[idx].FileName).
			SetURL(dataExports[idx].URL)
	}).Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create data export bulk: %w", err)
		return
	}

	newDataExports = lo.Map(des, func(de *ent.DataExport, _ int) *domain.DataExport {
		return convertDataExport(de)
	})

	return
}

func convertDataExport(de *ent.DataExport) *domain.DataExport {
	return &domain.DataExport{
		ID:           de.ID,
		StoreID:      de.StoreID,
		Type:         de.Type,
		Status:       de.Status,
		Params:       de.Params,
		FailedReason: de.FailedReason,
		OperatorType: de.OperatorType,
		OperatorID:   de.OperatorID,
		OperatorName: de.OperatorName,
		FileName:     de.FileName,
		URL:          de.URL,
		CreatedAt:    de.CreatedAt,
		UpdatedAt:    de.UpdatedAt,
	}
}

func (r *DataExportRepository) Update(ctx context.Context, dataExport *domain.DataExport) (updatedDataExport *domain.DataExport, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DataExportRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	de, err := r.Client.DataExport.UpdateOneID(dataExport.ID).
		SetStatus(dataExport.Status).
		SetFailedReason(dataExport.FailedReason).
		SetURL(dataExport.URL).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to update data export: %w", err)
		return
	}

	updatedDataExport = convertDataExport(de)

	return
}

func (r *DataExportRepository) Find(ctx context.Context, id int) (dataExport *domain.DataExport, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DataExportRepository.Find")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	de, err := r.Client.DataExport.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to get data export: %w", err)
		return
	}

	dataExport = convertDataExport(de)

	return
}

func (r *DataExportRepository) List(ctx context.Context, pager *upagination.Pagination, filter *domain.DataExportFilter) (dataExports []*domain.DataExport, total int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DataExportRepository.List")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := r.Client.DataExport.Query().
		Where(dataexport.StoreID(filter.StoreID))

	if filter.Type != "" {
		query = query.Where(dataexport.TypeEQ(filter.Type))
	}
	if filter.Status != "" {
		query = query.Where(dataexport.StatusEQ(filter.Status))
	}
	if filter.CreatedAtGte != nil {
		query = query.Where(dataexport.CreatedAtGTE(*filter.CreatedAtGte))
	}
	if filter.CreatedAtLte != nil {
		query = query.Where(dataexport.CreatedAtLTE(*filter.CreatedAtLte))
	}

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count: %w", err)
		return
	}

	des, err := query.Order(dataexport.ByCreatedAt(sql.OrderDesc()), dataexport.ByID(sql.OrderDesc())).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)

	if err != nil {
		err = fmt.Errorf("failed to get data exports: %w", err)
		return
	}

	dataExports = lo.Map(des, func(de *ent.DataExport, _ int) *domain.DataExport {
		return convertDataExport(de)
	})

	return
}
