package intl

import (
	"gitlab.jiguang.dev/pos-dine/dine/api/intl/pb"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Params struct {
	fx.In

	AppConfig domain.AppConfig

	ServerOptions []grpc.ServerOption `group:"server_options"`
	IntlService   pb.IntlServer
}

func New(p Params) *grpc.Server {
	srv := grpc.NewServer(p.ServerOptions...)
	pb.RegisterIntlServer(srv, p.IntlService)

	if p.AppConfig.RunMode == domain.RunModeDev {
		reflection.Register(srv)
	}

	return srv
}
