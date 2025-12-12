package usecasefx

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"usecase",
	fx.Provide(),
)
