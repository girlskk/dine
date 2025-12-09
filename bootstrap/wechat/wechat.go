package wechat

import (
	"github.com/redis/go-redis/v9"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
)

func NewMiniProgram(cfg Config, redisClient redis.UniversalClient) *miniprogram.MiniProgram {
	wc := wechat.NewWechat()
	che := NewRedisCache(redisClient)
	mcfg := &miniConfig.Config{
		AppID:     cfg.AppID,
		AppSecret: cfg.AppSecret,
		Cache:     che,
	}
	return wc.GetMiniProgram(mcfg)
}
