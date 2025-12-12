package backendfx

import (
	"net/http"

	"gitlab.jiguang.dev/pos-dine/dine/api/backend"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/httpserver"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/middleware"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"backendfx",
	fx.Provide(
		fx.Annotate(
			backend.New,
			fx.As(new(http.Handler)),
		)),

	// base middleware
	fx.Provide(
		asMiddleware(middleware.NewRecovery),
		asMiddleware(middleware.NewErrorHandling),
		asMiddleware(func(c httpserver.Config) *middleware.TimeLimiter { return middleware.NewTimeLimiter(c.RequestTimeout) }),
		asMiddleware(middleware.NewPopulateRequestID),
		asMiddleware(middleware.NewPopulateLogger),
		fx.Annotate(
			middleware.NewObservability,
			fx.As(new(ugin.Middleware)),
			fx.ResultTags(`group:"middlewares"`),
		),
		asMiddleware(middleware.NewLogger),
	),

	// handler
	fx.Provide(),
)

func asHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(ugin.Handler)),
		fx.ResultTags(`group:"handlers"`),
	)
}

func asMiddleware(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(ugin.Middleware)),
		fx.ResultTags(`group:"middlewares"`),
	)
}
