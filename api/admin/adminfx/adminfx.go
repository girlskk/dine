package adminfx

import (
	"net/http"

	"gitlab.jiguang.dev/pos-dine/dine/api/admin"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/handler"
	mid "gitlab.jiguang.dev/pos-dine/dine/api/admin/middleware"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/httpserver"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/middleware"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"admin",
	fx.Provide(
		fx.Annotate(
			admin.New,
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
		asMiddleware(middleware.NewLocale),
		fx.Annotate(
			middleware.NewObservability,
			fx.As(new(ugin.Middleware)),
			fx.ResultTags(`group:"middlewares"`),
		),
		asMiddleware(middleware.NewLogger),
		fx.Annotate(
			mid.NewAuth,
			fx.As(new(ugin.Middleware)),
			fx.ParamTags(`group:"handlers"`),
			fx.ResultTags(`group:"middlewares"`),
		),
	),
	// handler
	fx.Provide(
		asHandler(handler.NewUserHandler),
		asHandler(handler.NewMerchantHandler),
		asHandler(handler.NewStoreHandler),
		asHandler(handler.NewRegionHandler),
	),
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
