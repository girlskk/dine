package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

type StoreUserTestSuite struct {
	RepositoryTestSuite
	repo *StoreUserRepository
	ctx  context.Context
}

func TestStoreUserTestSuite(t *testing.T) {
	suite.Run(t, new(StoreUserTestSuite))
}

func (s *StoreUserTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.ctx = context.Background()
	s.repo = &StoreUserRepository{
		Client: s.client,
	}
}

func (s *StoreUserTestSuite) createTestStoreUser(tag string) *ent.StoreUser {
	hashedPassword, err := util.HashPassword("123456")
	require.NoError(s.T(), err)
	merchantID := uuid.New()
	storeID := uuid.New()
	departmentID := uuid.New()
	userID := uuid.New()

	return s.client.StoreUser.Create().
		SetID(userID).
		SetUsername(tag + "-newuser").
		SetHashedPassword(hashedPassword).
		SetNickname(tag + "-新用户").
		SetMerchantID(merchantID).
		SetStoreID(storeID).
		SetDepartmentID(departmentID).
		SetCode(tag + "-CODE-NEWUSER").
		SetRealName(tag + "-新用户真实姓名").
		SetGender(domain.GenderFemale).
		SetEmail(tag + "-newuser@dine.test").
		SetPhoneNumber("17700000000").
		SetEnabled(true).
		SetIsSuperadmin(false).
		SaveX(s.ctx)
}

func (s *StoreUserTestSuite) TestStoreUser_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		hashedPassword, err := util.HashPassword("123456")
		require.NoError(t, err)

		user := &domain.StoreUser{
			ID:             uuid.New(),
			Username:       "newuser",
			HashedPassword: hashedPassword,
			Nickname:       "新用户",
			MerchantID:     uuid.New(),
			StoreID:        uuid.New(),
			DepartmentID:   uuid.New(),
			Code:           "CODE-NEWUSER",
			RealName:       "新用户真实姓名",
			Gender:         domain.GenderOther,
			Email:          "newuser@dine.test",
			PhoneNumber:    "17700000000",
			Enabled:        true,
			IsSuperAdmin:   true,
		}
		require.NoError(t, s.repo.Create(s.ctx, user))
		require.NotEqual(t, uuid.Nil, user.ID)

		dbUser := s.client.StoreUser.GetX(s.ctx, user.ID)
		require.Equal(t, "newuser", dbUser.Username)
		require.Equal(t, "新用户", dbUser.Nickname)
		require.Equal(t, hashedPassword, dbUser.HashedPassword)
		require.Equal(t, "CODE-NEWUSER", dbUser.Code)
		require.Equal(t, "新用户真实姓名", dbUser.RealName)
		require.Equal(t, domain.GenderOther, dbUser.Gender)
		require.Equal(t, "newuser@dine.test", dbUser.Email)
		require.Equal(t, "17700000000", dbUser.PhoneNumber)
		require.True(t, dbUser.Enabled)
		require.True(t, dbUser.IsSuperadmin)
	})

	s.T().Run("创建重复用户名", func(t *testing.T) {
		hashedPassword, err := util.HashPassword("123456")
		require.NoError(t, err)

		user1 := &domain.StoreUser{
			ID:             uuid.New(),
			Username:       "duplicate",
			HashedPassword: hashedPassword,
			Nickname:       "用户1",
			MerchantID:     uuid.New(),
			StoreID:        uuid.New(),
			DepartmentID:   uuid.New(),
			Code:           "CODE-DUP",
			RealName:       "duplicate",
			Gender:         domain.GenderOther,
			Enabled:        true,
			IsSuperAdmin:   true,
		}
		require.NoError(t, s.repo.Create(s.ctx, user1))

		user2 := &domain.StoreUser{
			ID:             uuid.New(),
			Username:       "duplicate",
			HashedPassword: hashedPassword,
			Nickname:       "用户2",
			MerchantID:     uuid.New(),
			StoreID:        uuid.New(),
			DepartmentID:   uuid.New(),
			Code:           "CODE-DUP-2",
			RealName:       "duplicate",
			Gender:         domain.GenderOther,
			Email:          "duplicate2@dine.test",
			PhoneNumber:    "17700000002",
			Enabled:        true,
			IsSuperAdmin:   true,
		}
		require.Error(t, s.repo.Create(s.ctx, user2))
	})
}

func (s *StoreUserTestSuite) TestStoreUser_Find() {
	su := s.createTestStoreUser("find")

	s.T().Run("正常查询", func(t *testing.T) {
		user, err := s.repo.Find(s.ctx, su.ID)
		require.NoError(t, err)
		require.Equal(t, su.ID, user.ID)
		require.Equal(t, "find-newuser", user.Username)
		require.Equal(t, "find-新用户", user.Nickname)
		require.Equal(t, "find-新用户真实姓名", user.RealName)
		require.Equal(t, domain.GenderFemale, user.Gender)
		require.Equal(t, "find-newuser@dine.test", user.Email)
		require.Equal(t, "17700000000", user.PhoneNumber)
		require.True(t, user.Enabled)
		require.False(t, user.IsSuperAdmin)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		_, err := s.repo.Find(s.ctx, uuid.New())
		require.Error(t, err)
	})
}

func (s *StoreUserTestSuite) TestStoreUser_FindByUsername() {
	su := s.createTestStoreUser("findbyusername")

	s.T().Run("正常查询", func(t *testing.T) {
		user, err := s.repo.FindByUsername(s.ctx, "findbyusername-newuser")
		require.NoError(t, err)
		require.Equal(t, su.ID, user.ID)
		require.Equal(t, "findbyusername-newuser", user.Username)
		require.Equal(t, "findbyusername-新用户", user.Nickname)
		require.Equal(t, "findbyusername-新用户真实姓名", user.RealName)
		require.Equal(t, domain.GenderFemale, user.Gender)
		require.Equal(t, "findbyusername-newuser@dine.test", user.Email)
		require.Equal(t, "17700000000", user.PhoneNumber)
		require.True(t, user.Enabled)
		require.False(t, user.IsSuperAdmin)
	})

	s.T().Run("不存在的用户名", func(t *testing.T) {
		_, err := s.repo.FindByUsername(s.ctx, "nonexistent")
		require.Error(t, err)
	})
}

func (s *StoreUserTestSuite) TestStoreUser_Update() {
	su := s.createTestStoreUser("update")

	s.T().Run("正常更新", func(t *testing.T) {
		user := &domain.StoreUser{
			ID:             su.ID,
			Username:       "updateduser",
			HashedPassword: "updated_password",
			Nickname:       "更新后的昵称",
			RealName:       "更新后的真实姓名",
			Gender:         domain.GenderMale,
			Email:          "updated@dine.test",
			PhoneNumber:    "18800002222",
			Enabled:        false,
			IsSuperAdmin:   true,
		}
		require.NoError(t, s.repo.Update(s.ctx, user))

		updated := s.client.StoreUser.GetX(s.ctx, su.ID)
		require.Equal(t, "updateduser", updated.Username)
		require.Equal(t, "更新后的昵称", updated.Nickname)
		require.Equal(t, "updated_password", updated.HashedPassword)
		require.Equal(t, "更新后的真实姓名", updated.RealName)
		require.Equal(t, domain.GenderMale, updated.Gender)
		require.Equal(t, "updated@dine.test", updated.Email)
		require.Equal(t, "18800002222", updated.PhoneNumber)
		require.False(t, updated.Enabled)
		require.True(t, updated.IsSuperadmin)
	})

	s.T().Run("更新不存在的ID", func(t *testing.T) {
		require.Error(t, s.repo.Update(s.ctx, &domain.StoreUser{
			ID:             uuid.New(),
			Username:       "invalid",
			HashedPassword: "password",
			Nickname:       "无效用户",
		}))
	})
}

func (s *StoreUserTestSuite) TestStoreUser_Delete() {
	su := s.createTestStoreUser("delete")

	s.T().Run("正常删除", func(t *testing.T) {
		require.NoError(t, s.repo.Delete(s.ctx, su.ID))

		deleted, err := s.client.StoreUser.Get(s.ctx, su.ID)
		require.Nil(t, deleted)
		require.Error(t, err)
		require.True(t, ent.IsNotFound(err))
	})

	s.T().Run("删除不存在的ID", func(t *testing.T) {
		require.Error(t, s.repo.Delete(s.ctx, uuid.New()))
	})
}

func (s *StoreUserTestSuite) TestStoreUser_Integration() {
	s.T().Run("完整的CRUD流程", func(t *testing.T) {
		hashedPassword, err := util.HashPassword("123456")
		require.NoError(t, err)

		user := &domain.StoreUser{
			ID:             uuid.New(),
			Username:       "integration",
			HashedPassword: hashedPassword,
			Nickname:       "集成测试用户",
			RealName:       "集成真实姓名",
			Code:           "Integration-Code",
			Gender:         domain.GenderMale,
			Email:          "integration@dine.test",
			PhoneNumber:    "16600001111",
			Enabled:        true,
			IsSuperAdmin:   true,
			MerchantID:     uuid.New(),
			StoreID:        uuid.New(),
			DepartmentID:   uuid.New(),
		}
		require.NoError(t, s.repo.Create(s.ctx, user))

		found, err := s.repo.Find(s.ctx, user.ID)
		require.NoError(t, err)
		require.Equal(t, user.Username, found.Username)

		foundByUsername, err := s.repo.FindByUsername(s.ctx, user.Username)
		require.NoError(t, err)
		require.Equal(t, user.ID, foundByUsername.ID)

		user.Nickname = "更新后的昵称"
		require.NoError(t, s.repo.Update(s.ctx, user))

		updated, err := s.repo.Find(s.ctx, user.ID)
		require.NoError(t, err)
		require.Equal(t, "更新后的昵称", updated.Nickname)
		require.Equal(t, "集成真实姓名", updated.RealName)
		require.Equal(t, "integration", updated.Username)

		require.NoError(t, s.repo.Delete(s.ctx, user.ID))
		_, err = s.repo.Find(s.ctx, user.ID)
		require.Error(t, err)
	})
}
