package interceptor

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/metrics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryHandlerFunc(alert alert.Alert) recovery.RecoveryHandlerFuncContext {
	return func(ctx context.Context, p any) (err error) {
		logger := logging.FromContext(ctx).Named("interceptor.Recovery")
		logger.Errorw(
			"grpc handle panic",
			"error", p,
		)
		metrics.RecoverCounter.Inc()

		go alert.Notify(ctx, fmt.Sprintf("panic: %v", p))

		return status.Error(codes.Internal, "internal server error")
	}
}
