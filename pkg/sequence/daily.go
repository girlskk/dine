package sequence

import (
	"context"
	"fmt"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/redis/go-redis/v9"
)

var incrScript = redis.NewScript(`
	local val = redis.call('INCR', KEYS[1])
	if val == 1 then
		redis.call('EXPIREAT', KEYS[1], ARGV[1])
	end
	return val
`)

// DailySequence 基于 Redis 的每日递增序列
type DailySequence struct {
	rdb   redis.UniversalClient
	clock clockwork.Clock
}

// NewDailySequence 创建一个每日递增序列
func NewDailySequence(rdb redis.UniversalClient) *DailySequence {
	return &DailySequence{
		rdb:   rdb,
		clock: clockwork.NewRealClock(),
	}
}

// Next 获取下一个序列号，每天从 1 开始递增
func (s *DailySequence) Next(ctx context.Context, prefix string) (int64, error) {
	now := s.clock.Now()
	key := s.key(now, prefix)
	// 计算明天凌晨的时间戳
	yy, mm, dd := now.Date()
	expireAt := time.Date(yy, mm, dd, 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Unix()

	val, err := incrScript.Run(ctx, s.rdb, []string{key}, expireAt).Int64()
	if err != nil {
		return 0, fmt.Errorf("increment sequence: %w", err)
	}

	return val, nil
}

// Current 获取当前序列号
func (s *DailySequence) Current(ctx context.Context, prefix string) (int64, error) {
	now := s.clock.Now()
	key := s.key(now, prefix)
	val, err := s.rdb.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("get sequence: %w", err)
	}
	return val, nil
}

func (s *DailySequence) key(t time.Time, prefix string) string {
	return fmt.Sprintf("%s:%s", prefix, t.Format("20060102"))
}
