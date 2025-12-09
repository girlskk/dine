package customerfx

import (
	"net/http"

	"gitlab.jiguang.dev/pos-dine/dine/api/customer"
	"gitlab.jiguang.dev/pos-dine/dine/api/customer/handler"
	mid "gitlab.jiguang.dev/pos-dine/dine/api/customer/middleware"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/httpserver"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/middleware"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"customer",
	fx.Provide(
		fx.Annotate(
			customer.New,
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
		fx.Annotate(
			mid.NewAuth,
			fx.As(new(ugin.Middleware)),
			fx.ParamTags(`group:"handlers"`),
			fx.ResultTags(`group:"middlewares"`),
		),
	),
	// handler
	fx.Provide(
		asHandler(handler.NewCustomerHandler),
		asHandler(handler.NewProductHandler),
		asHandler(handler.NewOrderHandler),
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
