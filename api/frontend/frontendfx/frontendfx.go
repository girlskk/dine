package frontendfx

import (
	"net/http"

	"gitlab.jiguang.dev/pos-dine/dine/api/frontend"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/httpserver"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/middleware"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"frontend",
	fx.Provide(
		fx.Annotate(
			frontend.New,
			fx.As(new(http.Handler)),
		),
	),
	// base middleware
	fx.Provide(
		asMiddleware(middleware.NewRecovery),
		asMiddleware(middleware.NewErrorHandling),
		asMiddleware(func(c httpserver.Config) *middleware.TimeLimiter { return middleware.NewTimeLimiter(c.RequestTimeout) }),
		asMiddleware(middleware.NewPopulateRequestID),
		asMiddleware(middleware.NewPopulateLogger),
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
