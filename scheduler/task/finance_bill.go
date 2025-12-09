package task

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

const (
	TaskTypeFinanceBill = "finance_bill"
)

type FinanceBillTask struct {
	ReconciliationRecordInteractor domain.ReconciliationRecordInteractor
	Alert                          alert.Alert
}

func NewFinanceBillTask(interactor domain.ReconciliationRecordInteractor, alert alert.Alert) *FinanceBillTask {
	return &FinanceBillTask{
		ReconciliationRecordInteractor: interactor,
	}
}

func (task *FinanceBillTask) Type() string {
	return TaskTypeFinanceBill
}

// 每日定时生成财务账单
func (task *FinanceBillTask) ProcessTask(ctx context.Context, t *asynq.Task) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FinanceBillTask.ProcessTask")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	logger := logging.FromContext(ctx).Named("FinanceBillTask.ProcessTask")
	ctx = logging.NewContext(ctx, logger)

	err = task.ReconciliationRecordInteractor.GenerateDailyRecords(ctx)
	if err != nil {
		logger.Errorf("统计每日财务账单失败: %v", err)
		task.Alert.Notify(ctx, fmt.Sprintf("统计每日财务账单失败: %v", err))
	}

	return
}
