package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/domain/order"
	"gitlab.jiguang.dev/pos-dine/dine/domain/point_settlement"
	"gitlab.jiguang.dev/pos-dine/dine/domain/reconciliation"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.DataExportInteractor = (*DataExportInteractor)(nil)

const (
	DataExportTimeout = 10 * time.Minute
)

type DataExportInteractor struct {
	DataStore                     domain.DataStore
	MutexManager                  domain.MutexManager
	asynqClient                   *asynq.Client
	OrderListExporter             *order.OrderListExporter
	ReconciliationListExporter    *reconciliation.ReconciliationListExporter
	ReconciliationDetailExporter  *reconciliation.ReconciliationDetailExporter
	PointSettlementListExporter   *point_settlement.PointSettlementListExporter
	PointSettlementDetailExporter *point_settlement.PointSettlementDetailExporter
}

func NewDataExportInteractor(dataStore domain.DataStore, mutexManager domain.MutexManager, asynqClient *asynq.Client,
	orderListExporter *order.OrderListExporter,
	reconciliationListExporter *reconciliation.ReconciliationListExporter,
	reconciliationDetailExporter *reconciliation.ReconciliationDetailExporter,
	pointSettlementListExporter *point_settlement.PointSettlementListExporter,
	pointSettlementDetailExporter *point_settlement.PointSettlementDetailExporter,
) *DataExportInteractor {
	return &DataExportInteractor{
		DataStore:                     dataStore,
		MutexManager:                  mutexManager,
		asynqClient:                   asynqClient,
		OrderListExporter:             orderListExporter,
		ReconciliationListExporter:    reconciliationListExporter,
		ReconciliationDetailExporter:  reconciliationDetailExporter,
		PointSettlementListExporter:   pointSettlementListExporter,
		PointSettlementDetailExporter: pointSettlementDetailExporter,
	}
}

func (interactor *DataExportInteractor) Submit(ctx context.Context, params ...*domain.SubmitDataExportParams) (dataExports []*domain.DataExport, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DataExportInteractor.Submit")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	dataExports = make([]*domain.DataExport, 0, len(params))
	for _, param := range params {
		operatorInfo := domain.ExtractOperatorInfo(param.Operator)
		dataExports = append(dataExports, &domain.DataExport{
			StoreID:      param.StoreID,
			Type:         param.Type,
			Status:       domain.DataExportStatusPending,
			Params:       param.Params,
			FileName:     param.FileName,
			OperatorType: operatorInfo.Type,
			OperatorID:   operatorInfo.ID,
			OperatorName: operatorInfo.Name,
		})
	}

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		dataExports, err = ds.DataExportRepo().CreateBulk(ctx, dataExports)
		if err != nil {
			return fmt.Errorf("failed to create data export bulk: %w", err)
		}

		for _, dataExport := range dataExports {
			payload := domain.DataExportTaskPayload{
				ID: dataExport.ID,
			}

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("failed to marshal payload: %w", err)
			}

			t := asynq.NewTask(domain.TaskTypeDataExport, payloadBytes)
			if _, err = interactor.asynqClient.Enqueue(t, asynq.Timeout(DataExportTimeout)); err != nil {
				return fmt.Errorf("failed to enqueue data export task: %w", err)
			}
		}

		return nil
	})

	return
}

func (interactor *DataExportInteractor) Retry(ctx context.Context, storeID int, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DataExportInteractor.Retry")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	logger := logging.FromContext(ctx).Named("DataExportInteractor.Retry")
	logger = logger.With("id", id)
	ctx = logging.NewContext(ctx, logger)

	mu := interactor.MutexManager.NewMutex(domain.NewMutexDataExportKey(id))
	if err = mu.Lock(ctx); err != nil {
		if domain.IsAlreadyTakenError(err) {
			err = domain.ParamsErrorf("数据导出 %d 正在被其他操作修改", id)
		}
		err = fmt.Errorf("failed to lock data export: %w", err)
		return
	}
	defer func() {
		if _, err := mu.Unlock(ctx); err != nil {
			logger.Errorf("failed to unlock data export: %s", err)
		}
	}()

	dataExport, err := interactor.DataStore.DataExportRepo().Find(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsErrorf("数据导出记录不存在")
			return
		}
		err = fmt.Errorf("failed to find data export: %w", err)
		return
	}

	if dataExport.StoreID != storeID {
		err = domain.ParamsErrorf("数据导出记录不存在")
		return
	}

	if dataExport.Status != domain.DataExportStatusFailed {
		err = domain.DataExportStatusError
		return
	}

	payload, err := json.Marshal(domain.DataExportTaskPayload{
		ID: id,
	})
	if err != nil {
		err = fmt.Errorf("failed to marshal payload: %w", err)
		return
	}

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) (err error) {
		dataExport.Status = domain.DataExportStatusPending
		dataExport.FailedReason = ""
		if _, err = ds.DataExportRepo().Update(ctx, dataExport); err != nil {
			err = fmt.Errorf("failed to update data export: %w", err)
			return
		}

		t := asynq.NewTask(domain.TaskTypeDataExport, payload)
		if _, err = interactor.asynqClient.Enqueue(t, asynq.Timeout(DataExportTimeout)); err != nil {
			err = fmt.Errorf("failed to enqueue data export task: %w", err)
			return
		}
		return
	})

	return
}

func (interactor *DataExportInteractor) Run(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DataExportInteractor.Run")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	logger := logging.FromContext(ctx).Named("DataExportInteractor.Run")
	logger = logger.With("id", id)
	ctx = logging.NewContext(ctx, logger)

	mu := interactor.MutexManager.NewMutex(domain.NewMutexDataExportKey(id), domain.MutexWithExpiry(DataExportTimeout))
	if err = mu.Lock(ctx); err != nil {
		if domain.IsAlreadyTakenError(err) {
			err = domain.ParamsErrorf("数据导出 %d 正在被其他操作修改", id)
		}
		err = fmt.Errorf("failed to lock data export: %w", err)
		return
	}
	defer func() {
		if _, err := mu.Unlock(ctx); err != nil {
			logger.Errorf("failed to unlock data export: %s", err)
		}
	}()

	dataExport, err := interactor.DataStore.DataExportRepo().Find(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to find data export: %w", err)
		return
	}

	if dataExport.Status != domain.DataExportStatusPending {
		err = domain.DataExportStatusError
		return
	}

	exporter := interactor.getExporter(dataExport.Type)

	params := exporter.NewParams()
	if err = json.Unmarshal(dataExport.Params, params); err != nil {
		err = fmt.Errorf("failed to unmarshal params: %w", err)
		return
	}

	url, err := exporter.Export(ctx, dataExport.FileName, params)

	if err != nil && !domain.IsDataExportBizError(err) {
		err = fmt.Errorf("failed to export order list: %w", err)
		return
	}

	if err != nil {
		dataExport.Status = domain.DataExportStatusFailed
		dataExport.FailedReason = err.Error()
	} else {
		dataExport.Status = domain.DataExportStatusSuccess
		dataExport.URL = url
	}

	_, err = interactor.DataStore.DataExportRepo().Update(ctx, dataExport)
	if err != nil {
		err = fmt.Errorf("failed to update data export: %w", err)
		return
	}

	return
}

func (interactor *DataExportInteractor) getExporter(typ domain.DataExportType) domain.DataExporter {
	switch typ {
	default:
		panic("not implemented")
	case domain.DataExportTypeOrderListExport:
		return interactor.OrderListExporter
	case domain.DataExportTypeReconciliationRecordListExport:
		return interactor.ReconciliationListExporter
	case domain.DataExportTypeReconciliationRecordDetailsExport:
		return interactor.ReconciliationDetailExporter
	case domain.DataExportTypePointSettlementListExport:
		return interactor.PointSettlementListExporter
	case domain.DataExportTypePointSettlementDetailsExport:
		return interactor.PointSettlementDetailExporter
	}
}

func (interactor *DataExportInteractor) List(ctx context.Context, pager *upagination.Pagination, filter *domain.DataExportFilter) (dataExports []*domain.DataExport, total int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DataExportInteractor.List")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	dataExports, total, err = interactor.DataStore.DataExportRepo().List(ctx, pager, filter)
	if err != nil {
		err = fmt.Errorf("failed to list data exports: %w", err)
		return
	}

	return
}
