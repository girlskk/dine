package rdbfx

import (
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/rdb"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"redis",
	fx.Provide(
		rdb.New,
	),
)
