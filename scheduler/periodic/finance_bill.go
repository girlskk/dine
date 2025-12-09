package periodic

import (
	"time"

	"github.com/hibiken/asynq"
	"gitlab.jiguang.dev/pos-dine/dine/scheduler/task"
)

type FinanceBillConfig struct {
	baseConfig
}

type FinanceBill struct {
	config FinanceBillConfig
}

func NewFinanceBill(config FinanceBillConfig) *FinanceBill {
	return &FinanceBill{config: config}
}

func (e *FinanceBill) Register(scheduler *asynq.Scheduler) (err error) {
	opts := []asynq.Option{
		asynq.TaskID(task.TaskTypeFinanceBill),
	}

	if e.config.Timeout > 0 {
		opts = append(opts, asynq.Timeout(time.Duration(e.config.Timeout)*time.Second))
	}

	_, err = scheduler.Register(
		e.config.Cron,
		asynq.NewTask(task.TaskTypeFinanceBill, nil),
		opts...,
	)
	return
}
