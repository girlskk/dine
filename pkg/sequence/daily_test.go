package sequence

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/jonboulle/clockwork"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRedis struct {
	s   *miniredis.Miniredis
	rdb redis.UniversalClient
}

func newTestRedis(t *testing.T) *testRedis {
	s, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(func() { s.Close() })

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	t.Cleanup(func() { rdb.Close() })

	return &testRedis{s: s, rdb: rdb}
}

func TestDailySequence_Next(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T, tr *testRedis, seq *DailySequence)
	}{
		{
			name: "第一次调用返回1",
			validate: func(t *testing.T, tr *testRedis, seq *DailySequence) {
				val, err := seq.Next(context.Background(), "test")
				require.NoError(t, err)
				assert.Equal(t, int64(1), val)
			},
		},
		{
			name: "连续调用序号递增",
			validate: func(t *testing.T, tr *testRedis, seq *DailySequence) {
				ctx := context.Background()
				prefix := "test"

				val1, err := seq.Next(ctx, prefix)
				require.NoError(t, err)
				assert.Equal(t, int64(1), val1)

				val2, err := seq.Next(ctx, prefix)
				require.NoError(t, err)
				assert.Equal(t, int64(2), val2)

				val3, err := seq.Next(ctx, prefix)
				require.NoError(t, err)
				assert.Equal(t, int64(3), val3)
			},
		},
		{
			name: "不同日期使用不同序列",
			validate: func(t *testing.T, tr *testRedis, seq *DailySequence) {
				ctx := context.Background()
				now := time.Date(2025, 3, 7, 15, 30, 0, 0, time.Local)
				c := clockwork.NewFakeClockAt(now)
				seq.clock = c
				prefix := "test"

				tr.s.SetTime(now)

				// 获取今天的序列号
				val1, err := seq.Next(ctx, prefix)

				require.NoError(t, err)
				assert.Equal(t, int64(1), val1)

				// 验证今天的 key 存在
				todayKey := seq.key(now, prefix)
				assert.True(t, tr.s.Exists(todayKey))

				// 切换到明天
				tomorrow := now.Add(24 * time.Hour)
				tr.s.SetTime(tomorrow)
				tr.s.FastForward(24 * time.Hour)

				// 验证今天的 key 已过期
				assert.False(t, tr.s.Exists(todayKey))

				c.Advance(24 * time.Hour)
				// 获取明天的序列号
				val2, err := seq.Next(ctx, prefix)

				require.NoError(t, err)
				assert.Equal(t, int64(1), val2)

				// 验证明天的 key 存在
				tomorrowKey := seq.key(tomorrow, prefix)
				assert.True(t, tr.s.Exists(tomorrowKey))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := newTestRedis(t)
			seq := NewDailySequence(tr.rdb)

			tt.validate(t, tr, seq)
		})
	}
}

func TestDailySequence_Current(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, tr *testRedis, seq *DailySequence)
		wantVal int64
		wantErr bool
	}{
		{
			name:    "key不存在时返回0",
			wantVal: 0,
		},
		{
			name: "返回当前值",
			setup: func(t *testing.T, tr *testRedis, seq *DailySequence) {
				ctx := context.Background()
				val, err := seq.Next(ctx, "test")
				require.NoError(t, err)
				assert.Equal(t, int64(1), val)
			},
			wantVal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := newTestRedis(t)
			seq := NewDailySequence(tr.rdb)

			prefix := "test"
			if tt.setup != nil {
				tt.setup(t, tr, seq)
			}

			val, err := seq.Current(context.Background(), prefix)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantVal, val)
		})
	}
}
