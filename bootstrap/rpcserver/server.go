package rpcserver

import (
	"context"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func New(lc fx.Lifecycle, cfg Config, srv *grpc.Server) (net.Listener, error) {
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.Serve(lis); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.GracefulStop()
			return nil
		},
	})

	return lis, nil
}
