package domainservicefx

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"domainservice",
	fx.Provide(),
)
