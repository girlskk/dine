package usecasefx

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
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
	),
)
