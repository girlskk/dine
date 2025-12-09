package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	_ "gitlab.jiguang.dev/pos-dine/dine/ent/runtime"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

type RepositoryTestSuite struct {
	suite.Suite
	client *ent.Client
}

func (suite *RepositoryTestSuite) SetupTest() {
	suite.initDB()
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.NoError(suite.clearDB())
}

func (suite *RepositoryTestSuite) initDB() {
	t := suite.T()
	// opts := []enttest.Option{
	// 	enttest.WithOptions(ent.Log(t.Log)),
	// }

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
		"root", "pass", "127.0.0.1", "33061", "dine")
	// 连接到真实数据库
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed opening connection to mysql: %v", err)
	}
	suite.client = client.Debug()

	// suite.client = enttest.Open(t, "mysql", dsn, opts...).Debug()
}

func (suite *RepositoryTestSuite) clearDB() (err error) {
	if client := suite.client; client != nil {
		err = client.Close()
	}
	return
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestAtomic() {
	ctx := context.Background()
	repo := New(suite.client)

	suite.Run("测试事务正常提交", func() {
		err := repo.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
			err := ds.StoreRepo().Create(ctx, &domain.Store{
				Name:                util.RandomString(10),
				CooperationType:     domain.StoreCooperationTypeJoin,
				NeedAudit:           false,
				Enabled:             true,
				PointSettlementRate: decimal.Zero,
				PointWithdrawalRate: decimal.Zero,
			})
			suite.Require().NoError(err)

			err = ds.StoreRepo().Create(ctx, &domain.Store{
				Name:                util.RandomString(10),
				CooperationType:     domain.StoreCooperationTypeJoin,
				NeedAudit:           false,
				Enabled:             true,
				PointSettlementRate: decimal.Zero,
				PointWithdrawalRate: decimal.Zero,
			})
			suite.Require().NoError(err)
			return nil
		})
		suite.Require().NoError(err)
		defer func() {
			_, err := repo.client.Store.Delete().
				Exec(ctx)
			suite.Require().NoError(err)
		}()

		stores, err := repo.client.Store.Query().All(ctx)
		suite.Require().NoError(err)
		suite.Require().Len(stores, 2)
	})

	suite.Run("测试事务异常提交回滚", func() {
		err := repo.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
			err := ds.StoreRepo().Create(ctx, &domain.Store{
				Name:                util.RandomString(10),
				CooperationType:     domain.StoreCooperationTypeJoin,
				NeedAudit:           false,
				Enabled:             true,
				PointSettlementRate: decimal.Zero,
				PointWithdrawalRate: decimal.Zero,
			})
			suite.Require().NoError(err)

			err = ds.StoreRepo().Create(ctx, &domain.Store{
				Name:                util.RandomString(10),
				CooperationType:     domain.StoreCooperationTypeJoin,
				NeedAudit:           false,
				Enabled:             true,
				PointSettlementRate: decimal.Zero,
				PointWithdrawalRate: decimal.Zero,
			})
			suite.Require().NoError(err)

			return errors.New("test error")
		})
		suite.Require().Error(err)

		stores, err := repo.client.Store.Query().All(ctx)
		suite.Require().NoError(err)
		suite.Require().Empty(stores)
	})

	suite.Run("测试事务钩子", func() {
		trigger := false
		err := repo.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
			ds.AddHook(func() {
				trigger = true
			})

			err := ds.StoreRepo().Create(ctx, &domain.Store{
				Name:                util.RandomString(10),
				CooperationType:     domain.StoreCooperationTypeJoin,
				NeedAudit:           false,
				Enabled:             true,
				PointSettlementRate: decimal.Zero,
				PointWithdrawalRate: decimal.Zero,
			})
			suite.Require().NoError(err)

			err = ds.StoreRepo().Create(ctx, &domain.Store{
				Name:                util.RandomString(10),
				CooperationType:     domain.StoreCooperationTypeJoin,
				NeedAudit:           false,
				Enabled:             true,
				PointSettlementRate: decimal.Zero,
				PointWithdrawalRate: decimal.Zero,
			})
			suite.Require().NoError(err)

			return nil
		})
		suite.Require().NoError(err)
		defer func() {
			_, err := repo.client.Store.Delete().
				Exec(ctx)
			suite.Require().NoError(err)
		}()

		stores, err := repo.client.Store.Query().All(ctx)
		suite.Require().NoError(err)
		suite.Require().Len(stores, 2)
		suite.Require().True(trigger)
	})
}
