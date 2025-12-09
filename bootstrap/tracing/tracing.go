package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/fx"
)

// Config tracing config
type Config struct {
	ServiceName              string
	Disabled                 bool `default:"true"`
	AgentHostPort            string
	ReporterLogSpans         bool
	SamplerType              string
	SamplerParam             float64
	SamplerSamplingServerURL string
	AppVersion               string
	AppBuildAt               string
}

func New(lc fx.Lifecycle, conf Config) (opentracing.Tracer, error) {
	cfg := config.Configuration{
		ServiceName: conf.ServiceName,
		Disabled:    conf.Disabled,
		Sampler: &config.SamplerConfig{
			Type:              conf.SamplerType,
			Param:             conf.SamplerParam,
			SamplingServerURL: conf.SamplerSamplingServerURL,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: conf.AgentHostPort,
			LogSpans:           conf.ReporterLogSpans,
		},
		Tags: []opentracing.Tag{
			{Key: "app.version", Value: conf.AppVersion},
			{Key: "app.build_at", Value: conf.AppBuildAt},
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return closer.Close()
		},
	})

	return tracer, nil
}
