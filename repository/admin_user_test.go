package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/e"
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

func (s *AdminUserTestSuite) createTestAdminUser() *ent.AdminUser {
	userID := uuid.New()
	hashedPassword, err := util.HashPassword("123456")
	require.NoError(s.T(), err)

	return s.client.AdminUser.Create().
		SetID(userID).
		SetUsername("testuser").
		SetHashedPassword(hashedPassword).
		SetNickname("测试用户").
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
		}
		err = s.repo.Create(s.ctx, user)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, user.ID)

		// 验证数据库记录
		dbUser := s.client.AdminUser.GetX(s.ctx, user.ID)
		require.Equal(t, "newuser", dbUser.Username)
		require.Equal(t, "新用户", dbUser.Nickname)
		require.Equal(t, hashedPassword, dbUser.HashedPassword)
	})

	s.T().Run("创建重复用户名", func(t *testing.T) {
		// 先创建一个用户
		user1 := &domain.AdminUser{
			ID:             uuid.New(),
			Username:       "duplicate",
			HashedPassword: "hashed_password",
			Nickname:       "用户1",
		}
		err := s.repo.Create(s.ctx, user1)
		require.NoError(t, err)

		// 尝试创建相同用户名的用户（唯一约束冲突）
		user2 := &domain.AdminUser{
			ID:             uuid.New(),
			Username:       "duplicate",
			HashedPassword: "hashed_password",
			Nickname:       "用户2",
		}
		err = s.repo.Create(s.ctx, user2)
		require.Error(t, err)

		// 验证是 Conflict 错误
		var apiErr *errorx.Error
		require.True(t, errors.As(err, &apiErr))
		require.Equal(t, e.Conflict, apiErr.Code)
	})
}

func (s *AdminUserTestSuite) TestAdminUser_Find() {
	au := s.createTestAdminUser()

	s.T().Run("正常查询", func(t *testing.T) {
		user, err := s.repo.Find(s.ctx, au.ID)
		require.NoError(t, err)
		require.Equal(t, au.ID, user.ID)
		require.Equal(t, "testuser", user.Username)
		require.Equal(t, "测试用户", user.Nickname)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		nonExistentID := uuid.New()
		_, err := s.repo.Find(s.ctx, nonExistentID)
		require.Error(t, err)

		// 验证是 NotFound 错误
		var apiErr *errorx.Error
		require.True(t, errors.As(err, &apiErr))
		require.Equal(t, e.NotFound, apiErr.Code)
	})
}

func (s *AdminUserTestSuite) TestAdminUser_FindByUsername() {
	au := s.createTestAdminUser()

	s.T().Run("正常查询", func(t *testing.T) {
		user, err := s.repo.FindByUsername(s.ctx, "testuser")
		require.NoError(t, err)
		require.Equal(t, au.ID, user.ID)
		require.Equal(t, "testuser", user.Username)
		require.Equal(t, "测试用户", user.Nickname)
	})

	s.T().Run("不存在的用户名", func(t *testing.T) {
		_, err := s.repo.FindByUsername(s.ctx, "nonexistent")
		require.Error(t, err)

		// 验证是 NotFound 错误
		var apiErr *errorx.Error
		require.True(t, errors.As(err, &apiErr))
		require.Equal(t, e.NotFound, apiErr.Code)
	})
}

func (s *AdminUserTestSuite) TestAdminUser_Update() {
	au := s.createTestAdminUser()

	s.T().Run("正常更新", func(t *testing.T) {
		user := &domain.AdminUser{
			ID:             au.ID,
			Username:       "updateduser",
			HashedPassword: "updated_password",
			Nickname:       "更新后的昵称",
		}

		err := s.repo.Update(s.ctx, user)
		require.NoError(t, err)

		// 验证更新结果
		updated := s.client.AdminUser.GetX(s.ctx, au.ID)
		require.Equal(t, "updateduser", updated.Username)
		require.Equal(t, "更新后的昵称", updated.Nickname)
		require.Equal(t, "updated_password", updated.HashedPassword)
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

		// 验证是 NotFound 错误
		var apiErr *errorx.Error
		require.True(t, errors.As(err, &apiErr))
		require.Equal(t, e.NotFound, apiErr.Code)
	})
}

func (s *AdminUserTestSuite) TestAdminUser_Delete() {
	au := s.createTestAdminUser()

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

		// 验证是 NotFound 错误
		var apiErr *errorx.Error
		require.True(t, errors.As(err, &apiErr))
		require.Equal(t, e.NotFound, apiErr.Code)
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

		// Delete
		err = s.repo.Delete(s.ctx, user.ID)
		require.NoError(t, err)

		// Verify delete
		_, err = s.repo.Find(s.ctx, user.ID)
		require.Error(t, err)

		// 验证是 NotFound 错误
		var apiErr *errorx.Error
		require.True(t, errors.As(err, &apiErr))
		require.Equal(t, e.NotFound, apiErr.Code)
	})
}
