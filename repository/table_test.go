package repository

import (
	"context"
	"math/rand"

	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (suite *RepositoryTestSuite) TestTable_Create() {
	ctx := context.Background()
	store := suite.createTestStore(ctx)
	area := suite.createTestTableArea(ctx, store)
	_ = suite.createTestTable(ctx, area)
}

func (suite *RepositoryTestSuite) createTestTable(ctx context.Context, area *ent.TableArea) *ent.DineTable {
	return suite.client.DineTable.Create().
		SetAreaID(area.ID).
		SetName(util.RandomString(10)).
		SetStoreID(area.StoreID).
		SetSeatCount(rand.Intn(10) + 1).
		SaveX(ctx)
}
