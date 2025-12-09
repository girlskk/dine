package scheduler

import (
	"github.com/hibiken/asynq"
)

type Handler interface {
	asynq.Handler
	Type() string
}

type Periodic interface {
	Register(scheduler *asynq.Scheduler) (err error)
}
