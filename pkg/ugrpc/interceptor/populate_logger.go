package interceptor

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/utracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func PopulateLogger(originalLogger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		logger := originalLogger

		if id := utracing.RequestIDFromContext(ctx); id != "" {
			logger = logger.With("request_id", id)
		}

		ctx = logging.NewContext(ctx, logger)

		return handler(ctx, req)
	}
}
