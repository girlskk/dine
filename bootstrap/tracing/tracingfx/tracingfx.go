package tracingfx

import (
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/tracing"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"tracing",
	fx.Provide(
		tracing.New,
	),
)
