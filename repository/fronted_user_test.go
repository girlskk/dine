package repository

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (suite *RepositoryTestSuite) createTestFrontendUser(ctx context.Context, store *ent.Store) *ent.FrontendUser {
	return suite.client.FrontendUser.Create().
		SetUsername(util.RandomString(10)).
		SetHashedPassword(util.RandomString(10)).
		SetNickname(util.RandomString(10)).
		SetStore(store).
		SaveX(ctx)
}

func (suite *RepositoryTestSuite) TestFrontendUser_FindByUsername() {
	ctx := context.Background()
	store := suite.createTestStore(ctx)
	u := suite.createTestFrontendUser(ctx, store)

	repo := &FrontendUserRepository{
		Client: suite.client,
	}

	u2, err := repo.FindByUsername(ctx, u.Username)
	suite.Require().NoError(err)
	suite.Equal(u.ID, u2.ID)
	suite.Equal(u.Username, u2.Username)
	suite.Equal(u.HashedPassword, u2.HashedPassword)
	suite.Equal(u.Nickname, u2.Nickname)
	suite.NotNil(u2.Store)
}
