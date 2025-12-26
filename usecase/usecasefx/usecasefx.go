package usecasefx

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/additionalfee"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/category"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/merchant"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/remark"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/stall"
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
		fx.Annotate(
			remark.NewRemarkInteractor,
			fx.As(new(domain.RemarkInteractor)),
		),
		fx.Annotate(
			remark.NewRemarkCategoryInteractor,
			fx.As(new(domain.RemarkCategoryInteractor)),
		),
		fx.Annotate(
			stall.NewStallInteractor,
			fx.As(new(domain.StallInteractor)),
		),
		fx.Annotate(
			additionalfee.NewAdditionalFeeInteractor,
			fx.As(new(domain.AdditionalFeeInteractor)),
		),
	),
)
