package metrics

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/robfig/cron/v3"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func StartPush(ctx context.Context, c *Config, scheduler *cron.Cron, g prometheus.Gatherer) (func(), error) {
	nodeAddr, err := util.NodeETH0Addr()
	if err != nil {
		return nil, err
	}

	logger := logging.FromContext(ctx).Named("prom_pusher")

	id, err := scheduler.AddFunc(fmt.Sprintf("@every %ds", c.Interval), func() {
		if err := push.New(c.Addr, c.NodeName).
			Gatherer(g).
			Grouping("instance", nodeAddr).
			Add(); err != nil {
			logger.Errorf("Failed to push gateway: %v", err)
			return
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add func: %w", err)
	}

	return func() {
		scheduler.Remove(id)
		if err := push.New(c.Addr, c.NodeName).
			Grouping("instance", nodeAddr).Delete(); err != nil {
			logger.Errorf("Failed to delete node: %v", err)
		}
	}, nil
}
