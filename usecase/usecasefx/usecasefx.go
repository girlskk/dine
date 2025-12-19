package usecasefx

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/category"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/productattr"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/productspec"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/producttag"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/productunit"
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
			productunit.NewProductUnitInteractor,
			fx.As(new(domain.ProductUnitInteractor)),
		),
		fx.Annotate(
			productspec.NewProductSpecInteractor,
			fx.As(new(domain.ProductSpecInteractor)),
		),
		fx.Annotate(
			producttag.NewProductTagInteractor,
			fx.As(new(domain.ProductTagInteractor)),
		),
		fx.Annotate(
			productattr.NewProductAttrInteractor,
			fx.As(new(domain.ProductAttrInteractor)),
		),
	),
)
