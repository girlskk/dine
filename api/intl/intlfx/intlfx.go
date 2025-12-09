package intlfx

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"gitlab.jiguang.dev/pos-dine/dine/api/intl"
	"gitlab.jiguang.dev/pos-dine/dine/api/intl/pb"
	"gitlab.jiguang.dev/pos-dine/dine/api/intl/service"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/rpcserver"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugrpc/interceptor"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var Module = fx.Module(
	"intl",
	fx.Provide(
		intl.New,
	),
	fx.Provide(
		fx.Annotate(
			func(alert alert.Alert, cfg rpcserver.Config, logger *zap.SugaredLogger) grpc.ServerOption {
				return grpc.ChainUnaryInterceptor(
					recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(interceptor.RecoveryHandlerFunc(alert))),
					interceptor.TimeLimiter(cfg.RequestTimeout),
					interceptor.RequestID(),
					interceptor.Observability(),
					interceptor.PopulateLogger(logger),
					logging.UnaryServerInterceptor(
						interceptor.InterceptorLogger(),
						logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
					),
					interceptor.ErrorHandling(alert),
					interceptor.Validator(),
				)
			},
			fx.As(new(grpc.ServerOption)),
			fx.ResultTags(`group:"server_options"`),
		),
	),
	fx.Provide(
		fx.Annotate(
			service.NewIntlService,
			fx.As(new(pb.IntlServer)),
		),
	),
)
