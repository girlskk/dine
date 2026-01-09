package periodic

import (
	"time"

	"github.com/hibiken/asynq"
	"gitlab.jiguang.dev/pos-dine/dine/scheduler/task"
)

type ProfitDistributionConfig struct {
	baseConfig
}

type ProfitDistribution struct {
	config ProfitDistributionConfig
}

func NewProfitDistribution(config ProfitDistributionConfig) *ProfitDistribution {
	return &ProfitDistribution{config: config}
}

func (p *ProfitDistribution) Register(scheduler *asynq.Scheduler) (err error) {
	opts := []asynq.Option{
		asynq.TaskID(task.TaskTypeProfitDistribution),
	}

	if p.config.Timeout > 0 {
		opts = append(opts, asynq.Timeout(time.Duration(p.config.Timeout)*time.Second))
	}

	_, err = scheduler.Register(
		p.config.Cron,
		asynq.NewTask(task.TaskTypeProfitDistribution, nil),
		opts...,
	)
	return
}
