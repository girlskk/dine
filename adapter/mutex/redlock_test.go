package mutex

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

type RedlockTestSuite struct {
	suite.Suite
	mgr    *RedlockMutexManager
	server *miniredis.Miniredis
}

func (suite *RedlockTestSuite) SetupTest() {
	server, err := miniredis.Run()
	suite.Require().NoError(err)
	suite.server = server

	client := goredislib.NewClient(&goredislib.Options{
		Addr: server.Addr(),
	})
	suite.T().Cleanup(func() { client.Close() })
	suite.mgr = NewRedlockMutexManager(client)
}

func (suite *RedlockTestSuite) TearDownTest() {
	if suite.server != nil {
		suite.server.Close()
	}
}

func TestRedlockTestSuite(t *testing.T) {
	suite.Run(t, new(RedlockTestSuite))
}

func (suite *RedlockTestSuite) TestLock() {
	ctx := context.Background()
	key := "test_lock"

	locked := make(chan struct{})
	quit := make(chan struct{})
	go func() {
		m := suite.mgr.NewMutex(key)
		suite.Require().NoError(m.Lock(ctx))
		defer m.Unlock(ctx)
		close(locked)
		<-quit
	}()

	<-locked
	m := suite.mgr.NewMutex(key)
	err := m.TryLock(ctx)
	suite.Error(err)
	close(quit)
}

func (suite *RedlockTestSuite) TestLockWait() {
	ctx := context.Background()
	key := "test_lock_wait"

	locked := make(chan struct{})
	go func() {
		m := suite.mgr.NewMutex(key)
		suite.Require().NoError(m.Lock(ctx))
		defer m.Unlock(ctx)
		close(locked)
		time.Sleep(2 * time.Second)
	}()

	<-locked
	m := suite.mgr.NewMutex(key, domain.MutexWithWait(3*time.Second))
	err := m.Lock(ctx)
	suite.NoError(err)
}

func (suite *RedlockTestSuite) TestLockExpiry() {
	ctx := context.Background()
	key := "test_lock_expiry"

	m1 := suite.mgr.NewMutex(key, domain.MutexWithExpiry(3*time.Second))
	err := m1.Lock(ctx)
	suite.NoError(err)

	suite.server.FastForward(4 * time.Second)

	m2 := suite.mgr.NewMutex(key)
	err = m2.Lock(ctx)
	suite.NoError(err)
}
