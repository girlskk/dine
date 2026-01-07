package backendfx

import (
	"net/http"

	"gitlab.jiguang.dev/pos-dine/dine/api/backend"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/handler"
	mid "gitlab.jiguang.dev/pos-dine/dine/api/backend/middleware"
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
		asHandler(handler.NewProductUnitHandler),
		asHandler(handler.NewProductSpecHandler),
		asHandler(handler.NewProductTagHandler),
		asHandler(handler.NewProductAttrHandler),
		asHandler(handler.NewProductHandler),
		asHandler(handler.NewRemarkHandler),
		asHandler(handler.NewMenuHandler),
		asHandler(handler.NewProfitDistributionRuleHandler),
		asHandler(handler.NewProfitDistributionBillHandler),
		asHandler(handler.NewRegionHandler),
		asHandler(handler.NewStoreHandler),
		asHandler(handler.NewMerchantHandler),
		asHandler(handler.NewAdditionalFeeHandler),
		asHandler(handler.NewDeviceHandler),
		asHandler(handler.NewRemarkCategoryHandler),
		asHandler(handler.NewStallHandler),
		asHandler(handler.NewTaxFeeHandler),
		asHandler(handler.NewPaymentMethodHandler),
		asHandler(handler.NewPaymentAccountHandler),
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
