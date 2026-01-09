package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
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

func (s *DepartmentRepositoryTestSuite) newDepartment(tag string, merchantID, storeID uuid.UUID) *domain.Department {
	return &domain.Department{
		ID:             uuid.New(),
		Name:           "Dept-" + tag,
		Code:           "CODE-" + tag,
		DepartmentType: domain.DepartmentBackend,
		Enable:         true,
		MerchantID:     merchantID,
		StoreID:        storeID,
	}
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		dept := s.newDepartment("create", uuid.New(), uuid.Nil)

		err := s.repo.Create(s.ctx, dept)
		require.NoError(t, err)

		saved := s.client.Department.GetX(s.ctx, dept.ID)
		require.Equal(t, dept.Name, saved.Name)
		require.Equal(t, dept.Code, saved.Code)
		require.Equal(t, dept.DepartmentType, saved.DepartmentType)
		require.Equal(t, dept.Enable, saved.Enable)
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Create(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_FindByID() {
	dept := s.newDepartment("find", uuid.New(), uuid.Nil)
	require.NoError(s.T(), s.repo.Create(s.ctx, dept))

	s.T().Run("查询成功", func(t *testing.T) {
		found, err := s.repo.FindByID(s.ctx, dept.ID)
		require.NoError(t, err)
		require.Equal(t, dept.ID, found.ID)
		require.Equal(t, dept.Name, found.Name)
		require.Equal(t, dept.Code, found.Code)
	})

	s.T().Run("不存在", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_Update() {
	dept := s.newDepartment("update", uuid.New(), uuid.Nil)
	require.NoError(s.T(), s.repo.Create(s.ctx, dept))

	s.T().Run("更新成功", func(t *testing.T) {
		dept.Name = "Updated-" + dept.Name
		dept.Enable = false

		err := s.repo.Update(s.ctx, dept)
		require.NoError(t, err)

		updated := s.client.Department.GetX(s.ctx, dept.ID)
		require.Equal(t, dept.Name, updated.Name)
		require.Equal(t, dept.Enable, updated.Enable)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		missing := s.newDepartment("missing", uuid.New(), uuid.Nil)
		err := s.repo.Update(s.ctx, missing)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Update(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_Delete() {
	dept := s.newDepartment("delete", uuid.New(), uuid.Nil)
	require.NoError(s.T(), s.repo.Create(s.ctx, dept))

	s.T().Run("删除成功", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, dept.ID)
		require.NoError(t, err)

		_, err = s.client.Department.Get(s.ctx, dept.ID)
		require.Error(t, err)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_GetDepartments() {
	m1 := uuid.New()
	m2 := uuid.New()
	s1 := uuid.New()
	s2 := uuid.New()

	dept1 := s.newDepartment("001", m1, s1)
	require.NoError(s.T(), s.repo.Create(s.ctx, dept1))
	time.Sleep(10 * time.Millisecond)
	dept2 := s.newDepartment("002", m1, s1)
	dept2.DepartmentType = domain.DepartmentAdmin
	dept2.Enable = false
	require.NoError(s.T(), s.repo.Create(s.ctx, dept2))
	time.Sleep(10 * time.Millisecond)
	dept3 := s.newDepartment("003", m2, s2)
	dept3.DepartmentType = domain.DepartmentStore
	require.NoError(s.T(), s.repo.Create(s.ctx, dept3))

	pager := upagination.New(1, 10)

	s.T().Run("按商户筛选默认排序", func(t *testing.T) {
		list, total, err := s.repo.GetDepartments(s.ctx, pager, &domain.DepartmentListFilter{MerchantID: m1})
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, list, 2)
		require.Equal(t, dept2.ID, list[0].ID)
	})

	s.T().Run("按创建时间升序", func(t *testing.T) {
		order := domain.NewDepartmentListOrderByCreatedAt(false)
		list, total, err := s.repo.GetDepartments(s.ctx, pager, &domain.DepartmentListFilter{MerchantID: m1}, order)
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Equal(t, dept1.ID, list[0].ID)
	})

	s.T().Run("按部门类型与启用状态筛选", func(t *testing.T) {
		enable := true
		list, total, err := s.repo.GetDepartments(s.ctx, pager, &domain.DepartmentListFilter{MerchantID: m2, DepartmentType: domain.DepartmentStore, Enable: &enable})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, dept3.ID, list[0].ID)
	})

	s.T().Run("按名称模糊筛选", func(t *testing.T) {
		list, total, err := s.repo.GetDepartments(s.ctx, pager, &domain.DepartmentListFilter{Name: dept1.Name})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, dept1.ID, list[0].ID)
	})

	s.T().Run("按编码模糊筛选", func(t *testing.T) {
		list, total, err := s.repo.GetDepartments(s.ctx, pager, &domain.DepartmentListFilter{Code: "002"})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, dept2.ID, list[0].ID)
	})

	s.T().Run("无筛选返回全部", func(t *testing.T) {
		list, total, err := s.repo.GetDepartments(s.ctx, pager, nil)
		require.NoError(t, err)
		require.Equal(t, 3, total)
		require.Len(t, list, 3)
	})
}

func (s *DepartmentRepositoryTestSuite) TestDepartment_Exists() {
	merchantID := uuid.New()
	storeID := uuid.New()
	dept := s.newDepartment("exists", merchantID, storeID)
	require.NoError(s.T(), s.repo.Create(s.ctx, dept))

	s.T().Run("同名存在", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.DepartmentExistsParams{Name: dept.Name, MerchantID: merchantID, StoreID: storeID})
		require.NoError(t, err)
		require.True(t, exists)
	})

	s.T().Run("排除自身", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.DepartmentExistsParams{Name: dept.Name, MerchantID: merchantID, StoreID: storeID, ExcludeID: dept.ID})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("不同商户不冲突", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.DepartmentExistsParams{Name: dept.Name, MerchantID: uuid.New(), StoreID: storeID})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("不同门店不冲突", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.DepartmentExistsParams{Name: dept.Name, MerchantID: merchantID, StoreID: uuid.New()})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("仅名称查询", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.DepartmentExistsParams{Name: dept.Name})
		require.NoError(t, err)
		require.True(t, exists)
	})
}
