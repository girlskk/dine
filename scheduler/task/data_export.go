package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

type DataExport struct {
	DataExportInteractor domain.DataExportInteractor
	Alert                alert.Alert
}

func NewDataExportTask(dataExportInteractor domain.DataExportInteractor, alert alert.Alert) *DataExport {
	return &DataExport{
		DataExportInteractor: dataExportInteractor,
		Alert:                alert,
	}
}

func (task *DataExport) Type() string {
	return domain.TaskTypeDataExport
}

func (task *DataExport) ProcessTask(ctx context.Context, t *asynq.Task) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DataExport.ProcessTask")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	logger := logging.FromContext(ctx).Named("DataExport.ProcessTask")
	ctx = logging.NewContext(ctx, logger)

	var payload *domain.DataExportTaskPayload
	if err = json.Unmarshal(t.Payload(), &payload); err != nil {
		err = fmt.Errorf("json.Unmarshal payload: %w", err)
		return
	}

	defer func() {
		if err != nil {
			logger.Errorf("数据导出失败[%d]: %v", payload.ID, err)
			task.Alert.Notify(ctx, fmt.Sprintf("数据导出失败[%d]: %v", payload.ID, err))
		}
	}()

	logger.Infof("处理数据导出[%d]", payload.ID)
	if err = task.DataExportInteractor.Run(ctx, payload.ID); err != nil {
		if domain.IsNotFound(err) {
			logger.Warnf("data export %d not found", payload.ID)
			err = nil
			return
		} else if errors.Is(err, domain.DataExportStatusError) {
			logger.Warnf("data export %d status error", payload.ID)
			err = nil
			return
		}
		err = fmt.Errorf("DataExportInteractor.Run: %w", err)
		return
	}

	return
}
