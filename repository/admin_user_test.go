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

type AdminUserTestSuite struct {
	RepositoryTestSuite
	repo *AdminUserRepository
	ctx  context.Context
}

func TestAdminUserTestSuite(t *testing.T) {
	suite.Run(t, new(AdminUserTestSuite))
}

func (s *AdminUserTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &AdminUserRepository{
		Client: s.client,
	}
	s.ctx = context.Background()
}

func (s *AdminUserTestSuite) createTestAdminUser(tag string) *ent.AdminUser {
	userID := uuid.New()
	hashedPassword, err := util.HashPassword("123456")
	require.NoError(s.T(), err)
	user := &domain.AdminUser{
		ID:             uuid.New(),
		Username:       tag + "-newuser",
		HashedPassword: hashedPassword,
		Nickname:       tag + "-新用户",
		DepartmentID:   uuid.New(),
		Code:           tag + "-CODE-NEWUSER",
		RealName:       tag + "-新用户真实姓名",
		Gender:         domain.GenderFemale,
		Email:          tag + "-newuser@dine.test",
		PhoneNumber:    "17700000000",
		Enabled:        true,
		IsSuperAdmin:   false,
	}
	return s.client.AdminUser.Create().
		SetID(userID).
		SetUsername(user.Username).
		SetHashedPassword(user.HashedPassword).
		SetNickname(user.Nickname).
		SetDepartmentID(uuid.New()).
		SetCode(user.Code).
		SetRealName(user.RealName).
		SetGender(domain.GenderFemale).
		SetEmail(user.Email).
		SetPhoneNumber(user.PhoneNumber).
		SetEnabled(user.Enabled).
		SetIsSuperadmin(user.IsSuperAdmin).
		SaveX(s.ctx)
}

func (s *AdminUserTestSuite) TestAdminUser_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		hashedPassword, err := util.HashPassword("123456")
		require.NoError(s.T(), err)

		user := &domain.AdminUser{
			ID:             uuid.New(),
			Username:       "newuser",
			HashedPassword: hashedPassword,
			Nickname:       "新用户",
			DepartmentID:   uuid.New(),
			Code:           "CODE-NEWUSER",
			RealName:       "新用户真实姓名",
			Gender:         domain.GenderOther,
			Email:          "newuser@dine.test",
			PhoneNumber:    "17700000000",
			Enabled:        true,
			IsSuperAdmin:   true,
		}
		err = s.repo.Create(s.ctx, user)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, user.ID)

		// 验证数据库记录
		dbUser := s.client.AdminUser.GetX(s.ctx, user.ID)
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
		// 先创建一个用户
		user1 := &domain.AdminUser{
			ID:             uuid.New(),
			Username:       "duplicate",
			HashedPassword: "hashed_password",
			Nickname:       "用户1",
			DepartmentID:   uuid.New(),
			Code:           "CODE-DUP",
			RealName:       "duplicate",
			Gender:         domain.GenderOther,
			Enabled:        true,
			IsSuperAdmin:   true,
		}
		err := s.repo.Create(s.ctx, user1)
		require.NoError(t, err)

		// 尝试创建相同用户名的用户（唯一约束冲突）
		user2 := &domain.AdminUser{
			ID:             uuid.New(),
			Username:       "duplicate",
			HashedPassword: "hashed_password",
			Nickname:       "用户2",
			DepartmentID:   uuid.New(),
			Code:           "CODE-DUP-2",
			RealName:       "duplicate",
			Gender:         domain.GenderOther,
			Email:          "duplicate2@dine.test",
			PhoneNumber:    "17700000002",
			Enabled:        true,
			IsSuperAdmin:   true,
		}
		err = s.repo.Create(s.ctx, user2)
		require.Error(t, err)
	})
}

func (s *AdminUserTestSuite) TestAdminUser_Find() {
	au := s.createTestAdminUser("find")

	s.T().Run("正常查询", func(t *testing.T) {
		user, err := s.repo.Find(s.ctx, au.ID)
		require.NoError(t, err)
		require.Equal(t, au.ID, user.ID)
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
		nonExistentID := uuid.New()
		_, err := s.repo.Find(s.ctx, nonExistentID)
		require.Error(t, err)
	})
}

func (s *AdminUserTestSuite) TestAdminUser_FindByUsername() {
	au := s.createTestAdminUser("findbyusername")

	s.T().Run("正常查询", func(t *testing.T) {
		user, err := s.repo.FindByUsername(s.ctx, "findbyusername-newuser")
		require.NoError(t, err)
		require.Equal(t, au.ID, user.ID)
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

func (s *AdminUserTestSuite) TestAdminUser_Update() {
	au := s.createTestAdminUser("update")

	s.T().Run("正常更新", func(t *testing.T) {
		user := &domain.AdminUser{
			ID:             au.ID,
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

		err := s.repo.Update(s.ctx, user)
		require.NoError(t, err)

		// 验证更新结果
		updated := s.client.AdminUser.GetX(s.ctx, au.ID)
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
		nonExistentID := uuid.New()
		user := &domain.AdminUser{
			ID:             nonExistentID,
			Username:       "invalid",
			HashedPassword: "password",
			Nickname:       "无效用户",
		}

		err := s.repo.Update(s.ctx, user)
		require.Error(t, err)
	})
}

func (s *AdminUserTestSuite) TestAdminUser_Delete() {
	au := s.createTestAdminUser("delete")

	s.T().Run("正常删除", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, au.ID)
		require.NoError(t, err)

		// 验证删除结果
		deleted, err := s.client.AdminUser.Get(s.ctx, au.ID)
		require.Nil(t, deleted)
		require.Error(t, err)
		require.True(t, ent.IsNotFound(err))
	})

	s.T().Run("删除不存在的ID", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := s.repo.Delete(s.ctx, nonExistentID)
		require.Error(t, err)

	})
}

func (s *AdminUserTestSuite) TestAdminUser_Integration() {
	s.T().Run("完整的CRUD流程", func(t *testing.T) {
		// Create
		user := &domain.AdminUser{
			ID:             uuid.New(),
			Username:       "integration",
			HashedPassword: "password",
			Nickname:       "集成测试用户",
			RealName:       "集成真实姓名",
			Code:           "Integration-Code",
			Gender:         domain.GenderMale,
			Email:          "integration@dine.test",
			PhoneNumber:    "16600001111",
			Enabled:        true,
			IsSuperAdmin:   true,
			DepartmentID:   uuid.New(),
		}
		err := s.repo.Create(s.ctx, user)
		require.NoError(t, err)

		// Read by ID
		found, err := s.repo.Find(s.ctx, user.ID)
		require.NoError(t, err)
		require.Equal(t, user.Username, found.Username)

		// Read by Username
		foundByUsername, err := s.repo.FindByUsername(s.ctx, user.Username)
		require.NoError(t, err)
		require.Equal(t, user.ID, foundByUsername.ID)

		// Update
		user.Nickname = "更新后的昵称"
		err = s.repo.Update(s.ctx, user)
		require.NoError(t, err)

		// Verify update
		updated, err := s.repo.Find(s.ctx, user.ID)
		require.NoError(t, err)
		require.Equal(t, "更新后的昵称", updated.Nickname)
		require.Equal(t, "集成真实姓名", updated.RealName)
		require.Equal(t, "integration", updated.Username)

		// Delete
		err = s.repo.Delete(s.ctx, user.ID)
		require.NoError(t, err)

		// Verify delete
		_, err = s.repo.Find(s.ctx, user.ID)
		require.Error(t, err)
	})
}
