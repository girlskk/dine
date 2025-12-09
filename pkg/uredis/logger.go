package uredis

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
)

const slowThreshold = time.Millisecond * 100 // 设置慢查询阈值为100ms

type RedisLogger struct{}

// DialHook 实现连接钩子
func (h *RedisLogger) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		startTime := time.Now()
		conn, err := next(ctx, network, addr)
		duration := time.Since(startTime)
		if err != nil && !errors.Is(err, redis.Nil) {
			logger := logging.FromContext(ctx).Named("Redis.DialHook")
			logger.Errorf("Redis连接失败: network=%s, addr=%s, 耗时: %v, 错误: %v",
				network, addr, duration, err)
			return conn, err
		}
		return conn, err
	}
}

// ProcessHook 实现命令处理钩子
func (h *RedisLogger) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		startTime := time.Now()
		err := next(ctx, cmd)
		duration := time.Since(startTime)

		logger := logging.FromContext(ctx).Named("Redis.ProcessHook")
		if err != nil && !errors.Is(err, redis.Nil) {
			logger.Errorf("Redis命令执行失败: %s, 耗时: %v, 错误: %v",
				cmd.String(), duration, err)
			return err
		}

		// 只记录慢查询
		if duration > slowThreshold {
			logger.Warnf("Redis慢命令: %s, 耗时: %v", cmd.String(), duration)
		}
		return err
	}
}

// ProcessPipelineHook 实现管道命令处理钩子
func (h *RedisLogger) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		startTime := time.Now()
		err := next(ctx, cmds)
		duration := time.Since(startTime)

		logger := logging.FromContext(ctx).Named("Redis.ProcessPipelineHook")
		if err != nil && !errors.Is(err, redis.Nil) {
			logger.Errorf("Redis管道命令执行失败: 命令数量=%d, 耗时: %v, 错误: %v",
				len(cmds), duration, err)
			return err
		}

		// 检查每个命令的执行结果
		for _, cmd := range cmds {
			if cmdErr := cmd.Err(); cmdErr != nil && !errors.Is(cmdErr, redis.Nil) {
				logger.Errorf("Redis管道命令部分失败: %s, 耗时: %v, 错误: %v",
					cmd.String(), duration, cmdErr)
			}
		}

		if duration > slowThreshold {
			logger.Warnf("Redis慢管道命令: 命令数量=%d, 耗时: %v", len(cmds), duration)
		}
		return err
	}
}
