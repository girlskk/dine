package repository

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (suite *RepositoryTestSuite) TestTableArea_Create() {
	ctx := context.Background()
	store := suite.createTestStore(ctx)
	_ = suite.createTestTableArea(ctx, store)
}

func (suite *RepositoryTestSuite) createTestTableArea(ctx context.Context, store *ent.Store) *ent.TableArea {
	return suite.client.TableArea.Create().
		SetName(util.RandomString(10)).
		SetStoreID(store.ID).
		SaveX(ctx)
}
