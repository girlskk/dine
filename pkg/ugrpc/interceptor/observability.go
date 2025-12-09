package interceptor

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gitlab.jiguang.dev/pos-dine/dine/buildinfo"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/utracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Observability() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		tracer := opentracing.GlobalTracer()

		spanCtx, _ := tracer.Extract(opentracing.TextMap, metadataTextMap(md))
		span := tracer.StartSpan(
			info.FullMethod,
			opentracing.ChildOf(spanCtx),
			opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
			opentracing.Tag{Key: "grpc.method", Value: info.FullMethod},
			ext.SpanKindRPCServer,
			ext.RPCServerOption(spanCtx),
		)
		defer span.Finish()

		if id := utracing.RequestIDFromContext(ctx); id != "" {
			span.SetTag("request_id", id)
		}

		span.SetTag("build_version", buildinfo.Version)
		span.SetTag("build_at", buildinfo.BuildAt)

		ctx = opentracing.ContextWithSpan(ctx, span)

		return handler(ctx, req)
	}
}

// metadataTextMap 适配器
type metadataTextMap metadata.MD

func (m metadataTextMap) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range m {
		lowerKey := strings.ToLower(k)
		for _, v := range vals {
			if err := handler(lowerKey, v); err != nil {
				return err
			}
		}
	}
	return nil
}
