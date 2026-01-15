package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
)

type RoleRepositoryTestSuite struct {
	RepositoryTestSuite
	repo *RoleRepository
	ctx  context.Context
}

func TestRoleRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RoleRepositoryTestSuite))
}

func (s *RoleRepositoryTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.ctx = context.Background()
	s.repo = &RoleRepository{Client: s.client}
}

func (s *RoleRepositoryTestSuite) createTestRole(tag string, roleType domain.RoleType) *ent.Role {
	merchantID := uuid.New()
	storeID := uuid.New()
	return s.client.Role.Create().
		SetID(uuid.New()).
		SetName(tag + "-角色").
		SetCode(tag + "-CODE").
		SetRoleType(roleType).
		SetEnabled(true).
		SetDataScope(domain.RoleDataScopeDepartment).
		SetMerchantID(merchantID).
		SetStoreID(storeID).
		SetLoginChannels([]domain.LoginChannel{domain.LoginChannelPos}).
		SaveX(s.ctx)
}

func (s *RoleRepositoryTestSuite) TestRole_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		role := &domain.Role{
			ID:            uuid.New(),
			Name:          "新角色",
			Code:          "ROLE-NEW",
			RoleType:      domain.RoleTypeBackend,
			DataScope:     domain.RoleDataScopeCustom,
			Enabled:       true,
			MerchantID:    uuid.New(),
			StoreID:       uuid.New(),
			LoginChannels: []domain.LoginChannel{domain.LoginChannelStore},
		}
		require.NoError(t, s.repo.Create(s.ctx, role))
		require.NotEqual(t, uuid.Nil, role.ID)

		dbRole := s.client.Role.GetX(s.ctx, role.ID)
		require.Equal(t, "新角色", dbRole.Name)
		require.Equal(t, "ROLE-NEW", dbRole.Code)
		require.Equal(t, domain.RoleTypeBackend, dbRole.RoleType)
		require.Equal(t, domain.RoleDataScopeCustom, dbRole.DataScope)
		require.True(t, dbRole.Enabled)
		require.ElementsMatch(t, []domain.LoginChannel{domain.LoginChannelStore}, dbRole.LoginChannels)
	})

	s.T().Run("重复编码", func(t *testing.T) {
		role1 := &domain.Role{
			ID:            uuid.New(),
			Name:          "角色A",
			Code:          "ROLE-DUP",
			RoleType:      domain.RoleTypeAdmin,
			DataScope:     domain.RoleDataScopeAll,
			Enabled:       true,
			MerchantID:    uuid.New(),
			StoreID:       uuid.New(),
			LoginChannels: []domain.LoginChannel{domain.LoginChannelPos},
		}
		require.NoError(t, s.repo.Create(s.ctx, role1))

		role2 := &domain.Role{
			ID:            uuid.New(),
			Name:          "角色B",
			Code:          "ROLE-DUP",
			RoleType:      domain.RoleTypeAdmin,
			DataScope:     domain.RoleDataScopeAll,
			Enabled:       true,
			MerchantID:    uuid.New(),
			StoreID:       uuid.New(),
			LoginChannels: []domain.LoginChannel{domain.LoginChannelPos},
		}
		require.Error(t, s.repo.Create(s.ctx, role2))
	})
}

func (s *RoleRepositoryTestSuite) TestRole_FindByID() {
	entity := s.createTestRole("find", domain.RoleTypeStore)

	s.T().Run("存在的角色", func(t *testing.T) {
		role, err := s.repo.FindByID(s.ctx, entity.ID)
		require.NoError(t, err)
		require.Equal(t, entity.ID, role.ID)
		require.Equal(t, "find-角色", role.Name)
		require.Equal(t, "find-CODE", role.Code)
		require.Equal(t, domain.RoleTypeStore, role.RoleType)
		require.Equal(t, domain.RoleDataScopeDepartment, role.DataScope)
		require.True(t, role.Enabled)
	})

	s.T().Run("不存在的角色", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
	})
}

func (s *RoleRepositoryTestSuite) TestRole_Update() {
	entity := s.createTestRole("update", domain.RoleTypeBackend)

	s.T().Run("更新成功", func(t *testing.T) {
		updated := &domain.Role{
			ID:            entity.ID,
			Name:          "更新后的角色",
			RoleType:      domain.RoleTypeBackend,
			DataScope:     domain.RoleDataScopeMerchant,
			Enabled:       false,
			LoginChannels: []domain.LoginChannel{domain.LoginChannelMobile},
		}
		require.NoError(t, s.repo.Update(s.ctx, updated))

		dbRole := s.client.Role.GetX(s.ctx, entity.ID)
		require.Equal(t, "更新后的角色", dbRole.Name)
		require.False(t, dbRole.Enabled)
		require.Equal(t, domain.RoleDataScopeMerchant, dbRole.DataScope)
		require.ElementsMatch(t, []domain.LoginChannel{domain.LoginChannelMobile}, dbRole.LoginChannels)
	})

	s.T().Run("更新不存在角色", func(t *testing.T) {
		require.Error(t, s.repo.Update(s.ctx, &domain.Role{ID: uuid.New(), Name: "无效"}))
	})
}

func (s *RoleRepositoryTestSuite) TestRole_Delete() {
	entity := s.createTestRole("delete", domain.RoleTypeAdmin)

	s.T().Run("删除成功", func(t *testing.T) {
		require.NoError(t, s.repo.Delete(s.ctx, entity.ID))

		_, err := s.client.Role.Get(s.ctx, entity.ID)
		require.Error(t, err)
		require.True(t, ent.IsNotFound(err))
	})

	s.T().Run("删除不存在角色", func(t *testing.T) {
		require.Error(t, s.repo.Delete(s.ctx, uuid.New()))
	})
}

func (s *RoleRepositoryTestSuite) TestRole_Exists() {
	entity := s.createTestRole("exists", domain.RoleTypeBackend)

	s.T().Run("存在角色", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.RoleExistsParams{
			Name:       entity.Name,
			Code:       entity.Code,
			MerchantID: entity.MerchantID,
			StoreID:    entity.StoreID,
		})
		require.NoError(t, err)
		require.True(t, exists)
	})

	s.T().Run("排除自身后不存在", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.RoleExistsParams{
			Name:      entity.Name,
			ExcludeID: entity.ID,
		})
		require.NoError(t, err)
		require.False(t, exists)
	})
}

func (s *RoleRepositoryTestSuite) TestRole_ListByIDs() {
	r1 := s.createTestRole("list1", domain.RoleTypeAdmin)
	r2 := s.createTestRole("list2", domain.RoleTypeAdmin)

	s.T().Run("返回全部", func(t *testing.T) {
		roles, err := s.repo.ListByIDs(s.ctx, r1.ID, r2.ID)
		require.NoError(t, err)
		require.Len(t, roles, 2)
	})
}

func (s *RoleRepositoryTestSuite) TestRole_Integration() {
	s.T().Run("完整流程", func(t *testing.T) {
		role := &domain.Role{
			ID:            uuid.New(),
			Name:          "integration",
			Code:          "ROLE-INT",
			RoleType:      domain.RoleTypeStore,
			DataScope:     domain.RoleDataScopeStore,
			Enabled:       true,
			LoginChannels: []domain.LoginChannel{domain.LoginChannelStore},
		}
		require.NoError(t, s.repo.Create(s.ctx, role))

		found, err := s.repo.FindByID(s.ctx, role.ID)
		require.NoError(t, err)
		require.Equal(t, role.Name, found.Name)

		role.Name = "integration-2"
		require.NoError(t, s.repo.Update(s.ctx, role))

		updated, err := s.repo.FindByID(s.ctx, role.ID)
		require.NoError(t, err)
		require.Equal(t, "integration-2", updated.Name)

		require.NoError(t, s.repo.Delete(s.ctx, role.ID))
		_, err = s.repo.FindByID(s.ctx, role.ID)
		require.Error(t, err)
	})
}
