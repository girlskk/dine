package schedulerfx

import (
	"gitlab.jiguang.dev/pos-dine/dine/scheduler"
	"gitlab.jiguang.dev/pos-dine/dine/scheduler/periodic"
	"gitlab.jiguang.dev/pos-dine/dine/scheduler/task"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"scheduler",
	// provide handlers
	fx.Provide(
		asHandler(task.NewPaymentCallbackTask),
		asHandler(task.NewFinanceBillTask),
		asHandler(task.NewDataExportTask),
	),
	// provide periodic tasks
	fx.Provide(
		asPeriodic(periodic.NewFinanceBill),
	),
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
