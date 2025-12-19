package usecasefx

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/category"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/merchant"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/store"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/userauth"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"usecase",
	fx.Provide(
		fx.Annotate(
			userauth.NewAdminUserInteractor,
			fx.As(new(domain.AdminUserInteractor)),
		),
		fx.Annotate(
			category.NewCategoryInteractor,
			fx.As(new(domain.CategoryInteractor)),
		),
		fx.Annotate(
			userauth.NewBackendUserInteractor,
			fx.As(new(domain.BackendUserInteractor)),
		),
		fx.Annotate(
			merchant.NewMerchantInteractor,
			fx.As(new(domain.MerchantInteractor)),
		),
		fx.Annotate(
			store.NewStoreInteractor,
			fx.As(new(domain.StoreInteractor)),
		),
	),
)
