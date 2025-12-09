package domainservicefx

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain/order"
	"gitlab.jiguang.dev/pos-dine/dine/domain/payment"
	"gitlab.jiguang.dev/pos-dine/dine/domain/point_settlement"
	"gitlab.jiguang.dev/pos-dine/dine/domain/reconciliation"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"domainservice",
	fx.Provide(
		order.NewDomainService,
		order.NewOrderListExporter,
		payment.NewDomainService,
		reconciliation.NewDomainService,
		reconciliation.NewReconciliationListExporter,
		reconciliation.NewReconciliationDetailExporter,
		point_settlement.NewDomainService,
		point_settlement.NewPointSettlementListExporter,
		point_settlement.NewPointSettlementDetailExporter,
	),
)
