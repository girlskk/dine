package asynq

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"gitlab.jiguang.dev/pos-dine/dine/scheduler"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewServer(lc fx.Lifecycle, redisOpt asynq.RedisConnOpt, logger *zap.SugaredLogger, mux *asynq.ServeMux) *asynq.Server {
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Logger: logger,
		},
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := srv.Start(mux); err != nil {
				return fmt.Errorf("failed to start asynq server: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.Stop()
			srv.Shutdown()
			return nil
		},
	})

	return srv
}

func NewClient(redisOpt asynq.RedisConnOpt) *asynq.Client {
	return asynq.NewClient(redisOpt)
}

func NewServeMux(handlers []scheduler.Handler) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	for _, h := range handlers {
		mux.Handle(h.Type(), h)
	}
	return mux
}

type SchedulerParams struct {
	fx.In

	Lifecycle    fx.Lifecycle
	RedisConnOpt asynq.RedisConnOpt
	Logger       *zap.SugaredLogger
	Alert        alert.Alert
	Periodics    []scheduler.Periodic `group:"periodics"`
}

func NewScheduler(params SchedulerParams) (*asynq.Scheduler, error) {
	logger := params.Logger.Named("scheduler")
	errHandler := func(task *asynq.Task, opts []asynq.Option, err error) {
		logger = logger.Named("errorHandler")
		if errors.Is(err, asynq.ErrTaskIDConflict) {
			logger.Infof("任务[%s]已存在", task.Type())
			return
		}

		err = fmt.Errorf("任务[%s]调度失败: %w", task.Type(), err)
		logger.Error(err)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		go params.Alert.Notify(ctx, err.Error())
	}

	scheduler := asynq.NewScheduler(
		params.RedisConnOpt,
		&asynq.SchedulerOpts{
			Location:            time.Local,
			EnqueueErrorHandler: errHandler,
			Logger:              logger,
		},
	)

	for _, p := range params.Periodics {
		if err := p.Register(scheduler); err != nil {
			return nil, fmt.Errorf("failed to register periodic task: %w", err)
		}
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := scheduler.Start(); err != nil {
				return fmt.Errorf("failed to start scheduler: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			scheduler.Shutdown()
			return nil
		},
	})

	return scheduler, nil
}
