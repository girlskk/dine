package httpserverfx

import (
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/httpserver"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"httpserver",
	fx.Provide(
		httpserver.NewServer,
	),
)
