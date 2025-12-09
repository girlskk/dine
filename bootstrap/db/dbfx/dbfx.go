package dbfx

import (
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/db"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"database",
	fx.Provide(
		db.NewDB,
		db.NewClient,
	),
)
