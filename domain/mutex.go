package domain

import (
	"context"
	"fmt"
	"time"
)

const (
	mutexKeyOrderPrefix      = "mutex:order"
	mutexKeyPaymentPrefix    = "mutex:payment"
	mutexKeyDataExportPrefix = "mutex:data_export"
	mutexKeyCartPrefix       = "mutex:cart"
)

type MutexOption func(*MutexConfig)

type MutexConfig struct {
	Expiry time.Duration
	Wait   time.Duration
}

// MutexWithExpiry 锁过期时间（默认 8s）
func MutexWithExpiry(expiry time.Duration) MutexOption {
	return func(opts *MutexConfig) {
		opts.Expiry = expiry
	}
}

// MutexWithWait 保证的最小等待时间（默认 1.6s）
func MutexWithWait(wait time.Duration) MutexOption {
	return func(opts *MutexConfig) {
		opts.Wait = wait
	}
}

type Mutex interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) (bool, error)
	TryLock(ctx context.Context) error
}

type MutexManager interface {
	NewMutex(key string, options ...MutexOption) Mutex
}

func NewMutexOrderKey(orderNo string) string {
	return fmt.Sprintf("%s:%s", mutexKeyOrderPrefix, orderNo)
}

func NewMutexPaymentKey(seqNo string) string {
	return fmt.Sprintf("%s:%s", mutexKeyPaymentPrefix, seqNo)
}

func NewMutexDataExportKey(id int) string {
	return fmt.Sprintf("%s:%d", mutexKeyDataExportPrefix, id)
}

func NewMutexCartKey(tableID int) string {
	return fmt.Sprintf("%s:%d", mutexKeyCartPrefix, tableID)
}
