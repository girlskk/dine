package interceptor

import (
	"context"
	"fmt"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorHandling(alert alert.Alert) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		res, err := handler(ctx, req)
		if err != nil {
			status := status.Convert(err)
			if status != nil && (status.Code() == codes.Internal || status.Code() == codes.Unknown) {
				go alert.Notify(ctx, fmt.Sprintf("grpc server error: %v", err))
			}
		}
		return res, err
	}
}
