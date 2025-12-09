package repositoryfx

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/repository"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"repository",
	fx.Provide(
		fx.Annotate(
			repository.New,
			fx.As(new(domain.DataStore)),
		),
	),
)
