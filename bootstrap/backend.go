package bootstrap

import (
	"github.com/jinzhu/configor"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/alert"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/db"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/httpserver"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/rdb"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/zxh"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ali/oss"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/huifu"
	"go.uber.org/fx"
)

type BackendConfig struct {
	fx.Out

	App      domain.AppConfig
	HTTP     httpserver.Config
	Database db.Config
	Redis    rdb.Config
	Alert    alert.Config
	Auth     domain.AuthConfig
	Huifu    huifu.MerchSysConfig
	Zxh      zxh.Config
	Oss      oss.Config
}

func NewBackendConfig(files []string) (cfg BackendConfig, err error) {
	err = configor.Load(&cfg, files...)
	return
}
