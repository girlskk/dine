package rpcserverfx

import (
	"net"

	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/rpcserver"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"rpcserver",
	fx.Provide(
		fx.Annotate(
			rpcserver.New,
			fx.As(new(net.Listener)),
		),
	),
)
