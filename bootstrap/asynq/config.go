package asynq

import (
	"github.com/hibiken/asynq"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/rdb"
)

func NewRedisConnOpt(cfg rdb.Config) asynq.RedisConnOpt {
	return asynq.RedisClientOpt{
		Addr:     cfg.Addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
}
