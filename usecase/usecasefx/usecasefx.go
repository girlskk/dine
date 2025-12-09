package usecasefx

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/usecase"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/product"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/table"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"usecase",
	fx.Provide(
		fx.Annotate(
			usecase.NewFrontendUserInteractor,
			fx.As(new(domain.FrontendUserInteractor)),
		),
		fx.Annotate(
			usecase.NewBackendUserInteractor,
			fx.As(new(domain.BackendUserInteractor)),
		),
		fx.Annotate(
			usecase.NewAdminUserInteractor,
			fx.As(new(domain.AdminUserInteractor)),
		),

		fx.Annotate(
			table.NewTableAreaInteractor,
			fx.As(new(domain.TableAreaInteractor)),
		),

		fx.Annotate(
			table.NewTableInteractor,
			fx.As(new(domain.TableInteractor)),
		),

		fx.Annotate(
			product.NewCategoryInteractor,
			fx.As(new(domain.CategoryInteractor)),
		),

		fx.Annotate(
			product.NewProductInteractor,
			fx.As(new(domain.ProductInteractor)),
		),

		fx.Annotate(
			product.NewProductUnitInteractor,
			fx.As(new(domain.ProductUnitInteractor)),
		),

		fx.Annotate(
			product.NewProductAttrInteractor,
			fx.As(new(domain.ProductAttrInteractor)),
		),

		fx.Annotate(
			product.NewProductRecipeInteractor,
			fx.As(new(domain.ProductRecipeInteractor)),
		),

		fx.Annotate(
			product.NewProductSpecInteractor,
			fx.As(new(domain.ProductSpecInteractor)),
		),

		fx.Annotate(
			usecase.NewOrderInteractor,
			fx.As(new(domain.OrderInteractor)),
		),

		fx.Annotate(
			usecase.NewStoreInteractor,
			fx.As(new(domain.StoreInteractor)),
		),

		fx.Annotate(
			usecase.NewPaymentInteractor,
			fx.As(new(domain.PaymentInteractor)),
		),

		fx.Annotate(
			usecase.NewReconciliationRecordInteractor,
			fx.As(new(domain.ReconciliationRecordInteractor)),
		),
		fx.Annotate(
			usecase.NewPointSettlementInteractor,
			fx.As(new(domain.PointSettlementInteractor)),
		),
		fx.Annotate(
			usecase.NewDataExportInteractor,
			fx.As(new(domain.DataExportInteractor)),
		),
		fx.Annotate(
			usecase.NewStoreAccountInteractor,
			fx.As(new(domain.StoreAccountInteractor)),
		),
		fx.Annotate(
			usecase.NewStoreWithdrawInteractor,
			fx.As(new(domain.StoreWithdrawInteractor)),
		),
		fx.Annotate(
			usecase.NewOrderCartInteractor,
			fx.As(new(domain.OrderCartInteractor)),
		),
		fx.Annotate(
			usecase.NewCustomerInteractor,
			fx.As(new(domain.CustomerInteractor)),
		),
	),
)
