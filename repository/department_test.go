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

type DepartmentRepositoryTestSuite struct {
	RepositoryTestSuite
	repo *DepartmentRepository
	ctx  context.Context
}

func TestDepartmentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DepartmentRepositoryTestSuite))
}

func (s *DepartmentRepositoryTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &DepartmentRepository{Client: s.client}
	s.ctx = context.Background()
}

func (s *DepartmentRepositoryTestSuite) createTestDepartment(tag string, deptType domain.DepartmentType) *ent.Department {
	return s.client.Department.Create().
		SetID(uuid.New()).
		SetName(tag + "-部门").
		SetCode(tag + "-CODE").
		SetDepartmentType(deptType).
		SetEnable(true).
		SetMerchantID(uuid.New()).
		SetStoreID(uuid.New()).
		SaveX(s.ctx)
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_Create() {
	s.T().Run("正常创建", func(t *testing.T) {
		dept := &domain.Department{
			ID:             uuid.New(),
			Name:           "新部门",
			Code:           "DEPT-001",
			DepartmentType: domain.DepartmentBackend,
			Enable:         true,
		}
		require.NoError(t, s.repo.Create(s.ctx, dept))
		require.NotEqual(t, uuid.Nil, dept.ID)

		db := s.client.Department.GetX(s.ctx, dept.ID)
		require.Equal(t, "新部门", db.Name)
		require.Equal(t, "DEPT-001", db.Code)
		require.Equal(t, domain.DepartmentBackend, db.DepartmentType)
		require.True(t, db.Enable)
	})
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_FindByID() {
	entity := s.createTestDepartment("find", domain.DepartmentAdmin)

	s.T().Run("存在记录", func(t *testing.T) {
		dept, err := s.repo.FindByID(s.ctx, entity.ID)
		require.NoError(t, err)
		require.Equal(t, entity.ID, dept.ID)
		require.Equal(t, "find-部门", dept.Name)
	})

	s.T().Run("不存在记录", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
	})
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_Update() {
	entity := s.createTestDepartment("update", domain.DepartmentStore)

	s.T().Run("更新成功", func(t *testing.T) {
		payload := &domain.Department{
			ID:     entity.ID,
			Name:   "更新名称",
			Enable: false,
		}
		require.NoError(t, s.repo.Update(s.ctx, payload))

		db := s.client.Department.GetX(s.ctx, entity.ID)
		require.Equal(t, "更新名称", db.Name)
		require.False(t, db.Enable)
	})
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_Delete() {
	entity := s.createTestDepartment("delete", domain.DepartmentBackend)

	s.T().Run("删除成功", func(t *testing.T) {
		require.NoError(t, s.repo.Delete(s.ctx, entity.ID))
		_, err := s.client.Department.Get(s.ctx, entity.ID)
		require.Error(t, err)
		require.True(t, ent.IsNotFound(err))
	})

	s.T().Run("删除不存在", func(t *testing.T) {
		require.Error(t, s.repo.Delete(s.ctx, uuid.New()))
	})
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_Exists() {
	entity := s.createTestDepartment("exists", domain.DepartmentAdmin)

	s.T().Run("返回真", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.DepartmentExistsParams{Name: entity.Name, MerchantID: entity.MerchantID})
		require.NoError(t, err)
		require.True(t, exists)
	})

	s.T().Run("排除自身后假", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.DepartmentExistsParams{Name: entity.Name, ExcludeID: entity.ID})
		require.NoError(t, err)
		require.False(t, exists)
	})
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_CheckUserInDepartment() {
	d := s.createTestDepartment("checkadmin", domain.DepartmentAdmin)
	_, err := s.client.AdminUser.Create().
		SetID(uuid.New()).
		SetUsername("dept-admin").
		SetHashedPassword("pass").
		SetNickname("dept admin").
		SetDepartmentID(d.ID).
		SetCode("CODE-DEP").
		SetRealName("dept").
		SetGender(domain.GenderMale).
		SetEnabled(true).
		SetIsSuperadmin(false).
		Save(s.ctx)
	require.NoError(s.T(), err)

	s.T().Run("存在用户", func(t *testing.T) {
		exists, err := s.repo.CheckUserInDepartment(s.ctx, d.ID)
		require.NoError(t, err)
		require.True(t, exists)
	})

	s.T().Run("不存在部门报错", func(t *testing.T) {
		exists, err := s.repo.CheckUserInDepartment(s.ctx, uuid.New())
		require.Error(t, err)
		require.False(t, exists)
	})
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_Integration() {
	s.T().Run("CRUD流程", func(t *testing.T) {
		d := &domain.Department{
			ID:             uuid.New(),
			Name:           "integration",
			Code:           "DEPT-INT",
			DepartmentType: domain.DepartmentBackend,
			Enable:         true,
		}
		require.NoError(t, s.repo.Create(s.ctx, d))

		fetched, err := s.repo.FindByID(s.ctx, d.ID)
		require.NoError(t, err)
		require.Equal(t, d.Name, fetched.Name)

		d.Name = "integration-2"
		require.NoError(t, s.repo.Update(s.ctx, d))

		updated, err := s.repo.FindByID(s.ctx, d.ID)
		require.NoError(t, err)
		require.Equal(t, "integration-2", updated.Name)

		require.NoError(t, s.repo.Delete(s.ctx, d.ID))
		_, err = s.repo.FindByID(s.ctx, d.ID)
		require.Error(t, err)
	})
}
