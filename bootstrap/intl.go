package bootstrap

import (
	"github.com/jinzhu/configor"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/alert"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/db"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/rdb"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/rpcserver"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ali/oss"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/huifu"
	"go.uber.org/fx"
)

type IntlConfig struct {
	fx.Out

	App      domain.AppConfig
	RPC      rpcserver.Config
	Database db.Config
	Redis    rdb.Config
	Alert    alert.Config
	Auth     domain.AuthConfig
	Huifu    huifu.MerchSysConfig
	Oss      oss.Config
}

func NewIntlConfig(files []string) (cfg IntlConfig, err error) {
	err = configor.Load(&cfg, files...)
	return
}
