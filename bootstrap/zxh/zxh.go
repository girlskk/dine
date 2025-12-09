package zxh

import "gitlab.jiguang.dev/pos-dine/dine/pkg/zxh"

func New(cfg Config) *zxh.Manager {
	return zxh.NewManager(cfg.BaseURL)
}
