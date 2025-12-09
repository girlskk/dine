package adapterfx

import (
	"gitlab.jiguang.dev/pos-dine/dine/adapter/mutex"
	"gitlab.jiguang.dev/pos-dine/dine/adapter/objectstorage"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"adapter",
	fx.Provide(
		fx.Annotate(
			mutex.NewRedlockMutexManager,
			fx.As(new(domain.MutexManager)),
		),
		fx.Annotate(
			objectstorage.NewStorage,
			fx.As(new(domain.ObjectStorage)),
		),
	),
)
