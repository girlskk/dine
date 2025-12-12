package repository

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/enttest"
	_ "gitlab.jiguang.dev/pos-dine/dine/ent/runtime"
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
	opts := []enttest.Option{
		enttest.WithOptions(ent.Log(t.Log)),
	}
	suite.client = connectDB(t, opts)
}

func connectDB(t *testing.T, opts []enttest.Option) *ent.Client {
	// 使用本地真实数据库测试
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
	// 	"root", "pass", "127.0.0.1", "33061", "dine")
	// client, err := ent.Open("mysql", dsn)
	// if err != nil {
	// 	t.Fatalf("failed opening connection to mysql: %v", err)
	// }
	// return client.Debug()

	// 使用内存 SQLite 数据库进行测试
	return enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", opts...)
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
}
