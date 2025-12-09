package httpserver

import (
	"context"
	"net"
	"net/http"

	"go.uber.org/fx"
)

func NewServer(lc fx.Lifecycle, cfg Config, handler http.Handler) *http.Server {
	srv := &http.Server{Addr: ":" + cfg.Port, Handler: handler}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
