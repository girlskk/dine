package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gitlab.jiguang.dev/pos-dine/dine/buildinfo"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/utracing"
)

type Observability struct {
	skippers []SkipperFunc
}

func NewObservability(skippers ...SkipperFunc) *Observability {
	return &Observability{skippers: skippers}
}

func (o *Observability) Name() string {
	return "Observability"
}

func (o *Observability) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, o.skippers...) {
			c.Next()
			return
		}

		r := c.Request
		ctx := r.Context()

		tracer := opentracing.GlobalTracer()
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan(r.URL.Path, ext.RPCServerOption(spanCtx))
		defer span.Finish()

		if id := utracing.RequestIDFromContext(ctx); id != "" {
			span.SetTag("request_id", id)
		}
		span.SetTag("build_version", buildinfo.Version)
		span.SetTag("build_at", buildinfo.BuildAt)

		ctx = opentracing.ContextWithSpan(ctx, span)
		c.Request = r.Clone(ctx)

		c.Next()
	}
}
