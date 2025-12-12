package util

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func SpanErrFinish(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}

// StartSpan 创建并开始一个新的 span
// tracerName: tracer 名称
// spanName: span 名称
// attrs: 可选的属性
func StartSpan(ctx context.Context, tracerName, spanName string, attrs ...attribute.KeyValue) (trace.Span, context.Context) {
	tr := otel.Tracer(tracerName)
	ctx, span := tr.Start(ctx, spanName)
	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}
	return span, ctx
}
