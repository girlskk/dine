package schedulerfx

import (
	"gitlab.jiguang.dev/pos-dine/dine/scheduler"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"scheduler",
	// provide handlers
	fx.Provide(),
	// provide periodic tasks
	fx.Provide(),
)

func asHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(scheduler.Handler)),
		fx.ResultTags(`group:"handlers"`),
	)
}

func asPeriodic(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(scheduler.Periodic)),
		fx.ResultTags(`group:"periodics"`),
	)
}
