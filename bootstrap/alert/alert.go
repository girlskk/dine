package alert

import (
	"github.com/CatchZeng/dingtalk/pkg/dingtalk"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
)

const tmpl = `
### 服务告警

**服务名称**: %s

**时间**: %s

**错误信息**: <span style="color:red; font-weight:bold;">%s</span>
`

func New(c Config) alert.Alert {
	if c.Disabled {
		return &alert.AlertNoop{}
	}

	return &alertDingTalk{
		appName: c.AppName,
		ats:     c.Ats,
		client:  dingtalk.NewClient(c.AccessToken, c.Secret),
	}
}
