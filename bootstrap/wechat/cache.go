package wechat

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache .redis cache
type RedisCache struct {
	conn redis.UniversalClient
}

// NewRedisCache 实例化
func NewRedisCache(redisClient redis.UniversalClient) *RedisCache {
	return &RedisCache{conn: redisClient}
}

// SetConn 设置conn
func (r *RedisCache) SetConn(conn redis.UniversalClient) {
	r.conn = conn
}

// Get 获取一个值
func (r *RedisCache) Get(key string) any {
	return r.GetContext(context.Background(), key)
}

// GetContext 获取一个值
func (r *RedisCache) GetContext(ctx context.Context, key string) any {
	result, err := r.conn.Do(ctx, "GET", key).Result()
	if err != nil {
		return nil
	}
	return result
}

// Set 设置一个值
func (r *RedisCache) Set(key string, val any, timeout time.Duration) error {
	return r.SetContext(context.Background(), key, val, timeout)
}

// SetContext 设置一个值
func (r *RedisCache) SetContext(ctx context.Context, key string, val any, timeout time.Duration) error {
	return r.conn.SetEx(ctx, key, val, timeout).Err()
}

// IsExist 判断key是否存在
func (r *RedisCache) IsExist(key string) bool {
	return r.IsExistContext(context.Background(), key)
}

// IsExistContext 判断key是否存在
func (r *RedisCache) IsExistContext(ctx context.Context, key string) bool {
	result, _ := r.conn.Exists(ctx, key).Result()

	return result > 0
}

// Delete 删除
func (r *RedisCache) Delete(key string) error {
	return r.DeleteContext(context.Background(), key)
}

// DeleteContext 删除
func (r *RedisCache) DeleteContext(ctx context.Context, key string) error {
	return r.conn.Del(ctx, key).Err()
}
