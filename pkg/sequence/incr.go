package sequence

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// IncrSequence 生成一直递增的编号
type IncrSequence struct {
	rdb    redis.UniversalClient
	key    string
	prefix string
	width  int
}

// NewSequenceWithConfig 支持自定义 redis key、前缀和位数
func NewSequenceWithConfig(rdb redis.UniversalClient, key, prefix string, width int) (*IncrSequence, error) {
	if key == "" {
		return nil, errors.New("incr sequence key cannot be empty")
	}
	if prefix == "" {
		return nil, errors.New("incr sequence prefix cannot be empty")
	}
	if width <= 0 {
		width = 6
	}
	return &IncrSequence{
		rdb:    rdb,
		key:    key,
		prefix: prefix,
		width:  width,
	}, nil
}

// Next 返回下一个编号
func (s *IncrSequence) Next(ctx context.Context) (string, error) {
	val, err := s.rdb.Incr(ctx, s.key).Result()
	if err != nil {
		return "", fmt.Errorf("increment incr sequence: %w。key:%v,prefix:%v", err, s.key, s.prefix)
	}
	return s.format(val), nil
}

// Current 返回当前编号（还没生成则返回空串）
func (s *IncrSequence) Current(ctx context.Context) (string, error) {
	val, err := s.rdb.Get(ctx, s.key).Int64()
	if errors.Is(err, redis.Nil) {
		return s.format(0), nil
	}
	if err != nil {
		return "", fmt.Errorf("get incr sequence: %w。key:%v,prefix:%v", err, s.key, s.prefix)
	}
	return s.format(val), nil
}

func (s *IncrSequence) format(val int64) string {
	return fmt.Sprintf("%s%0*d", s.prefix, s.width, val)
}

// NewAdminTaxSequence 返回默认前缀 OSST，4 位数字的税率序列 运营后台使用
func NewAdminTaxSequence(rdb redis.UniversalClient) *IncrSequence {
	c, _ := NewSequenceWithConfig(rdb, domain.AdminTaxSequenceKey, "OSST", 4)
	return c
}

// NewBackendTaxSequence 返回默认前缀 BMPT，4 位数字的税率序列 品牌后台使用
func NewBackendTaxSequence(rdb redis.UniversalClient) *IncrSequence {
	c, _ := NewSequenceWithConfig(rdb, domain.BackendTaxSequenceKey, "BMPT", 4)
	return c
}

// NewStoreTaxSequence 返回默认前缀 SMPT，4 位数字的税率序列 门店后台使用
func NewStoreTaxSequence(rdb redis.UniversalClient) *IncrSequence {
	c, _ := NewSequenceWithConfig(rdb, domain.StoreTaxSequenceKey, "SMPT", 4)
	return c
}

// NewAdminDepartmentSequence 返回默认前缀 OSSD，4 位数字的部门序列 运营后台使用
func NewAdminDepartmentSequence(rdb redis.UniversalClient) *IncrSequence {
	c, _ := NewSequenceWithConfig(rdb, domain.AdminDepartmentSequenceKey, "OSSD", 4)
	return c
}

// NewBackendDepartmentSequence 返回默认前缀 BMPD，4 位数字的部门序列 品牌后台使用
func NewBackendDepartmentSequence(rdb redis.UniversalClient) *IncrSequence {
	c, _ := NewSequenceWithConfig(rdb, domain.BackendDepartmentSequenceKey, "BMPD", 4)
	return c
}

// NewStoreDepartmentSequence 返回默认前缀 SMPD，4 位数字的部门序列 门店后台使用
func NewStoreDepartmentSequence(rdb redis.UniversalClient) *IncrSequence {
	c, _ := NewSequenceWithConfig(rdb, domain.StoreDepartmentSequenceKey, "SMPD", 4)
	return c
}

// NewAdminUserSequence 返回默认前缀 OSSU，6 位数字的用户序列 运营后台使用
func NewAdminUserSequence(rdb redis.UniversalClient) *IncrSequence {
	c, _ := NewSequenceWithConfig(rdb, domain.AdminUserSequenceKey, "OSSU", 6)
	return c
}

// NewBackendUserSequence 返回默认前缀 BMPU，6 位数字的用户序列 品牌后台使用
func NewBackendUserSequence(rdb redis.UniversalClient) *IncrSequence {
	c, _ := NewSequenceWithConfig(rdb, domain.BackendUserSequenceKey, "BMPU", 6)
	return c
}

// NewStoreUserSequence 返回默认前缀 SMPU，6 位数字的用户序列 门店后台使用
func NewStoreUserSequence(rdb redis.UniversalClient) *IncrSequence {
	c, _ := NewSequenceWithConfig(rdb, domain.StoreUserSequenceKey, "SMPU", 6)
	return c
}

// NewAdminRoleSequence 返回默认前缀 OSSR，4 位数字的权限序列 运营后台使用
func NewAdminRoleSequence(rdb redis.UniversalClient) *IncrSequence {
	c, _ := NewSequenceWithConfig(rdb, domain.AdminRoleSequenceKey, "OSSR", 4)
	return c
}

// NewBackendRoleSequence 返回默认前缀 BMPR，4 位数字的权限序列 品牌后台使用
func NewBackendRoleSequence(rdb redis.UniversalClient) *IncrSequence {
	c, _ := NewSequenceWithConfig(rdb, domain.BackendRoleSequenceKey, "BMPR", 4)
	return c
}

// NewStoreRoleSequence 返回默认前缀 SMPR，4 位数字的权限序列 门店后台使用
func NewStoreRoleSequence(rdb redis.UniversalClient) *IncrSequence {
	c, _ := NewSequenceWithConfig(rdb, domain.StoreRoleSequenceKey, "SMPR", 4)
	return c
}
