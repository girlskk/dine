package mutex

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/e"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

const (
	minRetryDelayMilliSec = 50
	maxRetryDelayMilliSec = 250
)

var (
	_ domain.MutexManager = (*RedlockMutexManager)(nil)
	_ domain.Mutex        = (*redlockMutex)(nil)
)

type RedlockMutexManager struct {
	redsync *redsync.Redsync
}

func NewRedlockMutexManager(client goredislib.UniversalClient) *RedlockMutexManager {
	return &RedlockMutexManager{
		redsync: redsync.New(goredis.NewPool(client)),
	}
}

func (mgr *RedlockMutexManager) NewMutex(key string, options ...domain.MutexOption) domain.Mutex {
	opt := new(domain.MutexConfig)
	for _, option := range options {
		option(opt)
	}

	redsyncOpts := []redsync.Option{
		redsync.WithRetryDelayFunc(func(tries int) time.Duration {
			return time.Duration(rand.Intn(maxRetryDelayMilliSec-minRetryDelayMilliSec)+minRetryDelayMilliSec) * time.Millisecond
		}),
	}

	if opt.Expiry > time.Second {
		redsyncOpts = append(redsyncOpts, redsync.WithExpiry(opt.Expiry))
	} else {
		redsyncOpts = append(redsyncOpts, redsync.WithExpiry(8*time.Second))
	}

	if opt.Wait > time.Second {
		tries := int(opt.Wait / (minRetryDelayMilliSec * time.Millisecond))
		redsyncOpts = append(redsyncOpts, redsync.WithTries(tries))
	} else {
		redsyncOpts = append(redsyncOpts, redsync.WithTries(32))
	}

	return &redlockMutex{mgr.redsync.NewMutex(key, redsyncOpts...)}
}

type redlockMutex struct {
	*redsync.Mutex
}

func (m *redlockMutex) Lock(ctx context.Context) (err error) {
	span, ctx := util.StartSpan(ctx, "adapter", "RedlockMutex.Lock")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if err = m.Mutex.LockContext(ctx); err != nil && (isErrTaken(err) || isErrNodeTaken(err)) {
		return errorx.Fail(e.Conflict, err)
	}
	return
}

func (m *redlockMutex) TryLock(ctx context.Context) (err error) {
	span, ctx := util.StartSpan(ctx, "adapter", "RedlockMutex.TryLock")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if err := m.Mutex.TryLockContext(ctx); err != nil && (isErrTaken(err) || isErrNodeTaken(err)) {
		return errorx.Fail(e.Conflict, err)
	}
	return
}

func (m *redlockMutex) Unlock(ctx context.Context) (ok bool, err error) {
	span, _ := util.StartSpan(ctx, "adapter", "RedlockMutex.Unlock")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.Mutex.UnlockContext(ctx)
}

func isErrNodeTaken(err error) bool {
	if err == nil {
		return false
	}
	var e *redsync.ErrNodeTaken
	return errors.As(err, &e)
}

func isErrTaken(err error) bool {
	if err == nil {
		return false
	}
	var e *redsync.ErrTaken
	return errors.As(err, &e)
}
