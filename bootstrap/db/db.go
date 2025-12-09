package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"go.uber.org/fx"

	_ "github.com/go-sql-driver/mysql"
	"gitlab.jiguang.dev/pos-dine/dine/ent/migrate"
	_ "gitlab.jiguang.dev/pos-dine/dine/ent/runtime"
)

func NewDB(lc fx.Lifecycle, c Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", c.DSN())
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return db.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

func NewClient(lc fx.Lifecycle, c Config, db *sql.DB) (*ent.Client, error) {
	drv := entsql.OpenDB(dialect.MySQL, db)
	client := ent.NewClient(ent.Driver(drv))
	if c.Debug {
		client = client.Debug()
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if c.AutoMigrate {
				return client.Schema.Create(
					ctx,
					migrate.WithDropIndex(true),
					migrate.WithDropColumn(true),
				)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}
