package alertfx

import (
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/alert"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"alert",
	fx.Provide(
		alert.New,
	),
)
