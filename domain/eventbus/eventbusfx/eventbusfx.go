package eventbusfx

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/domain/eventbus"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"eventbus",
	fx.Provide(
		fx.Annotate(
			eventbus.NewOrderEventBus,
			fx.As(new(domain.OrderEventTrigger)),
		),
		fx.Annotate(
			eventbus.NewPaymentEventBus,
			fx.As(new(domain.PaymentEventTrigger)),
		),
	),
)
