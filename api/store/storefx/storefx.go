package storefx

import (
	"net/http"

	"gitlab.jiguang.dev/pos-dine/dine/api/store"
	"gitlab.jiguang.dev/pos-dine/dine/api/store/handler"
	mid "gitlab.jiguang.dev/pos-dine/dine/api/store/middleware"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/httpserver"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/middleware"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"storefx",
	fx.Provide(
		fx.Annotate(
			store.New,
			fx.As(new(http.Handler)),
		)),

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
		asHandler(handler.NewCategoryHandler),
		asHandler(handler.NewProductTagHandler),
		asHandler(handler.NewProductUnitHandler),
		asHandler(handler.NewProductSpecHandler),
		asHandler(handler.NewProductAttrHandler),
		asHandler(handler.NewProductHandler),
		asHandler(handler.NewMenuHandler),
		asHandler(handler.NewRegionHandler),
		asHandler(handler.NewStoreHandler),
		asHandler(handler.NewDeviceHandler),
		asHandler(handler.NewProfitDistributionBillHandler),
		asHandler(handler.NewRoleHandler),
		asHandler(handler.NewDepartmentHandler),
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
