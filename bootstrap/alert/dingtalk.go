package alert

import (
	"context"
	"fmt"
	"time"

	"github.com/CatchZeng/dingtalk/pkg/dingtalk"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
)

type alertDingTalk struct {
	appName string
	ats     []string
	client  *dingtalk.Client
}

func (a *alertDingTalk) Notify(ctx context.Context, errMsg string) {
	logger := logging.FromContext(ctx).Named("alert_dingtalk")

	content := fmt.Sprintf(tmpl, a.appName, time.Now().Format(time.DateTime), errMsg)
	msg := dingtalk.NewMarkdownMessage().
		SetMarkdown("服务告警", content).
		SetAt(a.ats, false)

	if _, _, err := a.client.Send(msg); err != nil {
		logger.Errorf("Failed to send alert to dingtalk: %v", err)
	}
}
