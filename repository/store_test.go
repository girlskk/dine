package repository

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (suite *RepositoryTestSuite) TestStore_Create() {
	ctx := context.Background()
	_ = suite.createTestStore(ctx)
}

func (suite *RepositoryTestSuite) createTestStore(ctx context.Context) *ent.Store {
	return suite.client.Store.Create().
		SetName(util.RandomString(10)).
		SetCooperationType(domain.StoreCooperationTypeJoin).
		SetNeedAudit(false).
		SetEnabled(true).
		SetHuifuID(util.RandomString(10)).
		SetZxhID(util.RandomString(10)).
		SetZxhSecret(util.RandomString(10)).
		SaveX(ctx)
}
