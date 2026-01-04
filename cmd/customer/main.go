package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gookit/event"
	"gitlab.jiguang.dev/pos-dine/dine/adapter/adapterfx"
	"gitlab.jiguang.dev/pos-dine/dine/api/customer/customerfx"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/alert/alertfx"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/asynq/asynqfx"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/db/dbfx"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/httpserver/httpserverfx"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/rdb/rdbfx"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/tracing/tracingfx"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/wechat"
	"gitlab.jiguang.dev/pos-dine/dine/buildinfo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/domain/domainservicefx"
	"gitlab.jiguang.dev/pos-dine/dine/domain/eventbus/eventbusfx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ali/oss"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/sequence"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
	"gitlab.jiguang.dev/pos-dine/dine/repository/repositoryfx"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/usecasefx"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var configFiles util.FlagStringArray

func init() {
	flag.Var(&configFiles, "conf", "App configuration files(.json,.yaml,.toml), multiple files are separated by ','")

	var err error
	time.Local, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	logger := logging.NewLoggerFromEnv().
		Named("customer").
		With("app", "customer").
		With("build_version", buildinfo.Version).
		With("build_at", buildinfo.BuildAt)
	defer logger.Sync()

	fx.New(
		fx.Supply(logger, event.Std()),
		fx.WithLogger(func(log *zap.SugaredLogger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Desugar()}
		}),
		fx.Provide(
			fx.Annotate(
				bootstrap.NewCustomerConfig,
				fx.ParamTags(`name:"config_files"`),
			),
			fx.Annotate(
				func() []string { return []string(configFiles) },
				fx.ResultTags(`name:"config_files"`),
			),
			fx.Annotate(
				sequence.NewDailySequence,
				fx.As(new(domain.DailySequence)),
			),
			oss.New,
			wechat.NewMiniProgram,
		),
		dbfx.Module,
		tracingfx.Module,
		repositoryfx.Module,
		usecasefx.Module,
		alertfx.Module,
		httpserverfx.Module,
		eventbusfx.Module,
		adapterfx.Module,
		customerfx.Module,
		domainservicefx.Module,
		rdbfx.Module,
		asynqfx.ClientModule,
		fx.Invoke(func(*http.Server) {}),
		fx.Invoke(func(tp trace.TracerProvider) {}),
	).Run()
}
