package bootstrap

import (
	"github.com/jinzhu/configor"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/alert"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/db"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/rdb"
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/zxh"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ali/oss"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/huifu"
	"gitlab.jiguang.dev/pos-dine/dine/scheduler/periodic"
	"go.uber.org/fx"
)

type SchedulerConfig struct {
	fx.Out

	Database db.Config
	Alert    alert.Config
	Redis    rdb.Config
	Zxh      zxh.Config
	Oss      oss.Config
	Huifu    huifu.MerchSysConfig

	FinanceBill periodic.FinanceBillConfig
}

func NewSchedulerConfig(files []string) (cfg SchedulerConfig, err error) {
	err = configor.Load(&cfg, files...)
	return
}
