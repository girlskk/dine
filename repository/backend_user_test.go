package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

type BackendUserTestSuite struct {
	RepositoryTestSuite
	repo *BackendUserRepository
	ctx  context.Context
}

func TestBackendUserTestSuite(t *testing.T) {
	suite.Run(t, new(BackendUserTestSuite))
}

func (s *BackendUserTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &BackendUserRepository{
		Client: s.client,
	}
	s.ctx = context.Background()
}

func (s *BackendUserTestSuite) createTestBackendUser() *ent.BackendUser {
	userID := uuid.New()
	merchantID := uuid.New()
	hashedPassword, err := util.HashPassword("123456")
	require.NoError(s.T(), err)

	return s.client.BackendUser.Create().
		SetID(userID).
		SetMerchantID(merchantID).
		SetUsername("admin").
		SetHashedPassword(hashedPassword).
		SetNickname("测试用户").
		SaveX(s.ctx)
}

func (s *BackendUserTestSuite) TestBackendUser_Create() {
	user := s.createTestBackendUser()
	require.NotNil(s.T(), user)
	require.NotEqual(s.T(), uuid.Nil, user.ID)
	require.Equal(s.T(), "admin", user.Username)
	require.Equal(s.T(), "测试用户", user.Nickname)
}
