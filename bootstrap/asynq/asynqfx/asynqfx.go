package asynqfx

import (
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/asynq"
	"go.uber.org/fx"
)

var ServerModule = fx.Module(
	"asynq_server",
	fx.Provide(
		asynq.NewRedisConnOpt,
		fx.Annotate(
			asynq.NewServeMux,
			fx.ParamTags(`group:"handlers"`),
		),
		asynq.NewServer,
		asynq.NewScheduler,
		asynq.NewClient,
	),
)

var ClientModule = fx.Module(
	"asynq_client",
	fx.Provide(
		asynq.NewRedisConnOpt,
		asynq.NewClient,
	),
)
