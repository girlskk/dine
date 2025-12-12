package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/tracing"
	"gitlab.jiguang.dev/pos-dine/dine/buildinfo"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/utracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Observability struct {
	skippers    []SkipperFunc
	serviceName string
	otelHandler gin.HandlerFunc // 缓存 otelgin 中间件
}

func NewObservability(conf tracing.Config, skippers ...SkipperFunc) *Observability {
	return &Observability{
		serviceName: conf.ServiceName,
		skippers:    skippers,
		otelHandler: otelgin.Middleware(conf.ServiceName),
	}
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

		o.otelHandler(c)

		// 获取 span 并添加自定义属性
		span := trace.SpanFromContext(c.Request.Context())
		if span.IsRecording() {
			if id := utracing.RequestIDFromContext(c.Request.Context()); id != "" {
				span.SetAttributes(attribute.String("request_id", id))
			}
			span.SetAttributes(
				attribute.String("build_version", buildinfo.Version),
				attribute.String("build_at", buildinfo.BuildAt),
			)
		}

		c.Next()
	}
}
