package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func TimeLimiter(timeout int) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
		return handler(ctx, req)
	}
}
