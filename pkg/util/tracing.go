package util

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func SpanErrFinish(span opentracing.Span, err error) {
	if err != nil {
		ext.LogError(span, err)
	}
	span.Finish()
}
