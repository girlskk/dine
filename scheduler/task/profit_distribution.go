package task

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

const (
	TaskTypeProfitDistribution = "profit_distribution"
)

type ProfitDistributionTask struct {
	ProfitDistributionBillInteractor domain.ProfitDistributionBillInteractor
	Alert                            alert.Alert
}

func NewProfitDistributionTask(interactor domain.ProfitDistributionBillInteractor, alert alert.Alert) *ProfitDistributionTask {
	return &ProfitDistributionTask{
		ProfitDistributionBillInteractor: interactor,
		Alert:                            alert,
	}
}

func (task *ProfitDistributionTask) Type() string {
	return TaskTypeProfitDistribution
}

// 每日定时生成分账账单
func (task *ProfitDistributionTask) ProcessTask(ctx context.Context, t *asynq.Task) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProfitDistributionTask.ProcessTask")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	logger := logging.FromContext(ctx).Named("ProfitDistributionTask.ProcessTask")
	ctx = logging.NewContext(ctx, logger)

	err = task.ProfitDistributionBillInteractor.GenerateProfitDistributionBills(ctx)
	if err != nil {
		logger.Errorf("生成分账账单失败: %v", err)
		task.Alert.Notify(ctx, fmt.Sprintf("生成分账账单失败: %v", err))
	}
	return
}
