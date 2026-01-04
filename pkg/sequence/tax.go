package sequence

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// TaxSequence 生成一直递增的税费编号
type TaxSequence struct {
	rdb    redis.UniversalClient
	key    string
	prefix string
	width  int
}

// NewTaxSequence 返回默认前缀 SST，4 位数字的税费序列
func NewTaxSequence(rdb redis.UniversalClient) *TaxSequence {
	return &TaxSequence{
		rdb:    rdb,
		key:    "tax:sequence",
		prefix: "SST",
		width:  4,
	}
}

// NewTaxSequenceWithConfig 支持自定义 redis key、前缀和位数
func NewTaxSequenceWithConfig(rdb redis.UniversalClient, key, prefix string, width int) *TaxSequence {
	if key == "" {
		key = "tax:sequence"
	}
	if prefix == "" {
		prefix = "SST"
	}
	if width <= 0 {
		width = 4
	}
	return &TaxSequence{
		rdb:    rdb,
		key:    key,
		prefix: prefix,
		width:  width,
	}
}

// Next 返回下一个税费编号
func (s *TaxSequence) Next(ctx context.Context) (string, error) {
	val, err := s.rdb.Incr(ctx, s.key).Result()
	if err != nil {
		return "", fmt.Errorf("increment tax sequence: %w", err)
	}
	return s.format(val), nil
}

// Current 返回当前编号（还没生成则返回空串）
func (s *TaxSequence) Current(ctx context.Context) (string, error) {
	val, err := s.rdb.Get(ctx, s.key).Int64()
	if err == redis.Nil {
		return s.format(0), nil
	}
	if err != nil {
		return "", fmt.Errorf("get tax sequence: %w", err)
	}
	return s.format(val), nil
}

func (s *TaxSequence) format(val int64) string {
	return fmt.Sprintf("%s%0*d", s.prefix, s.width, val)
}
