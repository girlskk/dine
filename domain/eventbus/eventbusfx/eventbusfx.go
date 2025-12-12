package eventbusfx

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"eventbus",
	fx.Provide(),
)
