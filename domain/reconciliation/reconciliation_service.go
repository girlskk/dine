package reconciliation

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

type DomainService struct {
	DataStore     domain.DataStore
	ObjectStorage domain.ObjectStorage
}

func NewDomainService(dataStore domain.DataStore, objectStorage domain.ObjectStorage) *DomainService {
	return &DomainService{
		DataStore:     dataStore,
		ObjectStorage: objectStorage,
	}
}

// 财务对账单列表导出
func (s *DomainService) ReconciliationListExport(ctx context.Context, filename string,
	params *domain.ReconciliationRecordListExportParams,
) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationDomainService.ReconciliationListExport")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	reconciliationRecords, err := s.DataStore.ReconciliationRecordRepo().PagedListBySearch(ctx, &params.Pager, params.Filter)
	if err != nil {
		err = fmt.Errorf("failed to get reconciliation records: %w", err)
		return
	}

	name, _ := util.GetFileNameAndExt(filename)

	url, err = s.GenerateReconciliationListExport(ctx, name, reconciliationRecords.Items)
	if err != nil {
		err = fmt.Errorf("failed to export reconciliation record list: %w", err)
		return
	}

	return
}

// 生成财务对账单列表Excel
func (s *DomainService) GenerateReconciliationListExport(ctx context.Context, name string,
	records []*domain.ReconciliationRecord,
) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationDomainService.GenerateReconciliationListExport")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	headers := []string{
		"序号", "单据编号", "账单日期", "门店名称", "支付渠道", "入账笔数", "入账金额",
	}

	var contents [][]string

	for index, record := range records {
		contents = append(contents, []string{
			fmt.Sprintf("%d", index+1),
			record.No,
			record.Date.Format(time.DateOnly),
			record.StoreName,
			record.Channel.ToString(),
			fmt.Sprintf("%d", record.OrderCount),
			record.Amount.String(),
		})
	}

	url, err = s.ObjectStorage.ExportExcel(ctx, domain.ObjectStorageSceneReconciliationListExport, name, headers, contents)
	if err != nil {
		err = fmt.Errorf("failed to export reconciliation list: %w", err)
		return
	}

	return
}

// 财务对账单明细导出
func (s *DomainService) ReconciliationDetailExport(ctx context.Context, filename string,
	params *domain.ReconciliationRecordListExportParams,
) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationDomainService.ReconciliationDetailExport")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	reconciliationRecords, err := s.DataStore.ReconciliationRecordRepo().PagedListBySearch(ctx, &params.Pager, params.Filter)
	if err != nil {
		err = fmt.Errorf("failed to get reconciliation records: %w", err)
		return
	}
	// 门店+账单日期+渠道 为key
	reconciliationRecordMap := make(map[string]*domain.ReconciliationRecord)
	for _, record := range reconciliationRecords.Items {
		key := fmt.Sprintf("%d-%s-%s", record.StoreID, record.Date.Format(time.DateOnly), record.Channel)
		reconciliationRecordMap[key] = record
	}

	// 获取订单详情
	dayStart := util.DayStart(*params.Filter.StartAt)
	dayEnd := util.DayEnd(*params.Filter.EndAt)
	filter := &domain.OrderListFilter{
		Status:        domain.OrderStatusPaid,
		FinishedAtGte: &dayStart,
		FinishedAtLte: &dayEnd,
		StoreID:       params.Filter.StoreID,
	}

	orders, _, err := s.DataStore.OrderRepo().GetOrders(ctx, &upagination.Pagination{
		Page: 1,
		Size: upagination.MaxSize,
	}, filter)
	if err != nil {
		err = fmt.Errorf("failed to get orders: %w", err)
		return
	}

	details := make([]*domain.ReconciliationDetailExportColumn, 0)
	for _, order := range orders {
		for _, channel := range order.PaidChannels {
			key := fmt.Sprintf("%d-%s-%s", order.StoreID, order.FinishedAt.Format(time.DateOnly), channel)
			reconciliationRecord, ok := reconciliationRecordMap[key]
			if !ok {
				continue
			}
			details = append(details, &domain.ReconciliationDetailExportColumn{
				RecordNo:        reconciliationRecord.No,
				RecordDate:      reconciliationRecord.Date,
				StoreName:       reconciliationRecord.StoreName,
				Channel:         channel,
				OrderNo:         order.No,
				OrderFinishedAt: *order.FinishedAt,
				OrderAmount:     order.TotalPrice,
			})
		}
	}
	name, _ := util.GetFileNameAndExt(filename)
	url, err = s.GenerateReconciliationDetailExport(ctx, name, details)
	if err != nil {
		err = fmt.Errorf("failed to export reconciliation detail: %w", err)
		return
	}
	return
}

// 生成财务对账单明细Excel
func (s *DomainService) GenerateReconciliationDetailExport(ctx context.Context, name string,
	details []*domain.ReconciliationDetailExportColumn,
) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationDomainService.GenerateReconciliationDetailExport")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	headers := []string{
		"单据编号", "账单日期", "门店名称", "支付渠道", "订单交易完成时间", "订单号", "订单金额",
	}

	var contents [][]string

	for _, detail := range details {
		contents = append(contents, []string{
			detail.RecordNo,
			detail.RecordDate.Format(time.DateOnly),
			detail.StoreName,
			detail.Channel.ToString(),
			detail.OrderFinishedAt.Format(time.DateTime),
			detail.OrderNo,
			detail.OrderAmount.String(),
		})
	}

	url, err = s.ObjectStorage.ExportExcel(ctx, domain.ObjectStorageSceneReconciliationDetailExport, name, headers, contents)
	if err != nil {
		err = fmt.Errorf("failed to export reconciliation detail: %w", err)
		return
	}

	return
}

// ----------------------------------------------------------------------------------------
// 财务对账单列表导出器
// ----------------------------------------------------------------------------------------
var _ domain.DataExporter = (*ReconciliationListExporter)(nil)

type ReconciliationListExporter struct {
	DomainService *DomainService
}

func NewReconciliationListExporter(domainService *DomainService) *ReconciliationListExporter {
	return &ReconciliationListExporter{
		DomainService: domainService,
	}
}

func (s *ReconciliationListExporter) NewParams() any {
	return new(domain.ReconciliationRecordListExportParams)
}

func (s *ReconciliationListExporter) Export(ctx context.Context, filename string, params any) (url string, err error) {
	return s.DomainService.ReconciliationListExport(ctx, filename, params.(*domain.ReconciliationRecordListExportParams))
}

// ----------------------------------------------------------------------------------------
// 财务对账单明细导出器
// ----------------------------------------------------------------------------------------
var _ domain.DataExporter = (*ReconciliationDetailExporter)(nil)

type ReconciliationDetailExporter struct {
	DomainService *DomainService
}

func NewReconciliationDetailExporter(domainService *DomainService) *ReconciliationDetailExporter {
	return &ReconciliationDetailExporter{
		DomainService: domainService,
	}
}

func (s *ReconciliationDetailExporter) NewParams() any {
	return new(domain.ReconciliationRecordListExportParams)
}

func (s *ReconciliationDetailExporter) Export(ctx context.Context, filename string, params any) (url string, err error) {
	return s.DomainService.ReconciliationDetailExport(ctx, filename, params.(*domain.ReconciliationRecordListExportParams))
}
