package point_settlement

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

// 积分结算单列表导出
func (s *DomainService) PointSettlementListExport(ctx context.Context, filename string,
	params *domain.PointSettlementListExportParams,
) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementDomainService.PointSettlementListExport")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	pointSettlements, err := s.DataStore.PointSettlementRepo().PagedListBySearch(ctx, &params.Pager, params.Filter)
	if err != nil {
		err = fmt.Errorf("failed to get point settlements: %w", err)
		return
	}

	name, _ := util.GetFileNameAndExt(filename)

	url, err = s.GeneratePointSettlementListExport(ctx, name, pointSettlements.Items)
	if err != nil {
		err = fmt.Errorf("failed to export point settlement list: %w", err)
		return
	}

	return
}

// 生成积分结算单列表Excel
func (s *DomainService) GeneratePointSettlementListExport(ctx context.Context, name string,
	records []*domain.PointSettlement,
) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementDomainService.GeneratePointSettlementListExport")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	headers := []string{
		"序号", "单据编号", "账单日期", "门店名称", "入账笔数", "入账金额",
	}

	var contents [][]string

	for index, record := range records {
		contents = append(contents, []string{
			fmt.Sprintf("%d", index+1),
			record.No,
			record.Date.Format(time.DateOnly),
			record.StoreName,
			fmt.Sprintf("%d", record.OrderCount),
			record.Amount.String(),
		})
	}

	url, err = s.ObjectStorage.ExportExcel(ctx, domain.ObjectStorageScenePointSettlementListExport, name, headers, contents)
	if err != nil {
		err = fmt.Errorf("failed to export point settlement list: %w", err)
		return
	}

	return
}

// 积分结算单明细导出
func (s *DomainService) PointSettlementDetailExport(ctx context.Context, filename string,
	params *domain.PointSettlementListExportParams,
) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementDomainService.PointSettlementDetailExport")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	pointSettlements, err := s.DataStore.PointSettlementRepo().PagedListBySearch(ctx, &params.Pager, params.Filter)
	if err != nil {
		err = fmt.Errorf("failed to get point settlements: %w", err)
		return
	}
	// 门店+账单日期 为key
	pointSettlementMap := make(map[string]*domain.PointSettlement)
	for _, record := range pointSettlements.Items {
		key := fmt.Sprintf("%d-%s", record.StoreID, record.Date.Format(time.DateOnly))
		pointSettlementMap[key] = record
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

	details := make([]*domain.PointSettlementDetailExportColumn, 0)
	for _, order := range orders {
		for _, channel := range order.PaidChannels {
			if channel != domain.OrderPaidChannelPoint {
				continue
			}
			key := fmt.Sprintf("%d-%s", order.StoreID, order.FinishedAt.Format(time.DateOnly))
			pointSettlement, ok := pointSettlementMap[key]
			if !ok {
				continue
			}
			details = append(details, &domain.PointSettlementDetailExportColumn{
				RecordNo:        pointSettlement.No,
				RecordDate:      pointSettlement.Date,
				StoreName:       pointSettlement.StoreName,
				OrderNo:         order.No,
				OrderFinishedAt: *order.FinishedAt,
				OrderAmount:     order.TotalPrice,
			})
		}
	}
	name, _ := util.GetFileNameAndExt(filename)
	url, err = s.GeneratePointSettlementDetailExport(ctx, name, details)
	if err != nil {
		err = fmt.Errorf("failed to export point settlement detail: %w", err)
		return
	}
	return
}

// 生成积分结算单明细Excel
func (s *DomainService) GeneratePointSettlementDetailExport(ctx context.Context, name string,
	details []*domain.PointSettlementDetailExportColumn,
) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementDomainService.GeneratePointSettlementDetailExport")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	headers := []string{
		"单据编号", "账单日期", "门店名称", "订单交易完成时间", "订单号", "订单金额",
	}

	var contents [][]string

	for _, detail := range details {
		contents = append(contents, []string{
			detail.RecordNo,
			detail.RecordDate.Format(time.DateOnly),
			detail.StoreName,
			detail.OrderFinishedAt.Format(time.DateTime),
			detail.OrderNo,
			detail.OrderAmount.String(),
		})
	}

	url, err = s.ObjectStorage.ExportExcel(ctx, domain.ObjectStorageScenePointSettlementDetailExport, name, headers, contents)
	if err != nil {
		err = fmt.Errorf("failed to export point settlement detail: %w", err)
		return
	}

	return
}

// ----------------------------------------------------------------------------------------
// 积分结算单列表导出器
// ----------------------------------------------------------------------------------------
var _ domain.DataExporter = (*PointSettlementListExporter)(nil)

type PointSettlementListExporter struct {
	DomainService *DomainService
}

func NewPointSettlementListExporter(domainService *DomainService) *PointSettlementListExporter {
	return &PointSettlementListExporter{
		DomainService: domainService,
	}
}

func (s *PointSettlementListExporter) NewParams() any {
	return new(domain.PointSettlementListExportParams)
}

func (s *PointSettlementListExporter) Export(ctx context.Context, filename string, params any) (url string, err error) {
	return s.DomainService.PointSettlementListExport(ctx, filename, params.(*domain.PointSettlementListExportParams))
}

// ----------------------------------------------------------------------------------------
// 积分结算单明细导出器
// ----------------------------------------------------------------------------------------
var _ domain.DataExporter = (*PointSettlementDetailExporter)(nil)

type PointSettlementDetailExporter struct {
	DomainService *DomainService
}

func NewPointSettlementDetailExporter(domainService *DomainService) *PointSettlementDetailExporter {
	return &PointSettlementDetailExporter{
		DomainService: domainService,
	}
}

func (s *PointSettlementDetailExporter) NewParams() any {
	return new(domain.PointSettlementListExportParams)
}

func (s *PointSettlementDetailExporter) Export(ctx context.Context, filename string, params any) (url string, err error) {
	return s.DomainService.PointSettlementDetailExport(ctx, filename, params.(*domain.PointSettlementListExportParams))
}
