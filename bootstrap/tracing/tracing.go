package tracing

import (
	"context"
	"fmt"
	"time"

	"gitlab.jiguang.dev/pos-dine/dine/buildinfo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/fx"
)

// Config tracing config
type Config struct {
	ServiceName              string
	Disabled                 bool   `default:"true"`
	OTLPEndpoint             string `default:"http://localhost:4318"`
	ReporterLogSpans         bool
	SamplerType              string
	SamplerParam             float64
	SamplerSamplingServerURL string
	AppVersion               string
	AppBuildAt               string
}

func New(lc fx.Lifecycle, conf Config) (trace.TracerProvider, error) {
	// 添加 Disabled 检查
	if conf.Disabled {
		return noop.NewTracerProvider(), nil
	}
	// 创建 OTLP HTTP exporter
	// OTLPEndpoint 可以是：
	// - "http://localhost:4318" (OTLP HTTP)
	// - "http://localhost:4317" (OTLP gRPC，需要使用 otlptracegrpc)
	ctx := context.Background()
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(conf.OTLPEndpoint),
		otlptracehttp.WithInsecure(), // 如果使用 HTTP，可以设置为 insecure
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otlp exporter: %w", err)
	}

	// 创建 resource（资源信息）
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(conf.ServiceName),
			semconv.ServiceVersionKey.String(buildinfo.Version),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 配置采样器
	var sampler tracesdk.Sampler
	switch conf.SamplerType {
	case "const":
		if conf.SamplerParam >= 1 {
			sampler = tracesdk.AlwaysSample()
		} else {
			sampler = tracesdk.NeverSample()
		}
	case "probabilistic":
		sampler = tracesdk.TraceIDRatioBased(conf.SamplerParam)
	case "ratelimiting":
		// 需要额外的配置，这里简化处理
		sampler = tracesdk.AlwaysSample()
	default:
		sampler = tracesdk.AlwaysSample()
	}

	// 创建 TracerProvider
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),     // 批量发送到 OTLP endpoint
		tracesdk.WithResource(res),    // 设置资源信息
		tracesdk.WithSampler(sampler), // 设置采样器
	)

	// 设置为全局 TracerProvider
	otel.SetTracerProvider(tp)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// 给 shutdown 设置超时
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return tp.Shutdown(shutdownCtx)
		},
	})

	return tp, nil
}
