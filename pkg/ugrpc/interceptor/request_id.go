package interceptor

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/utracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const MetaDataXRequestID = "x-request-id"

func RequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		requestIDs := md.Get(MetaDataXRequestID)

		var requestID string
		if len(requestIDs) > 0 {
			requestID = requestIDs[0]
		} else {
			requestID = uuid.New().String()
		}

		ctx = utracing.NewRequestIDContext(ctx, requestID)

		return handler(ctx, req)
	}
}
