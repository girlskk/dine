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

type StallRepositoryTestSuite struct {
	RepositoryTestSuite
	repo *StallRepository
	ctx  context.Context
}

func TestStallRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(StallRepositoryTestSuite))
}

func (s *StallRepositoryTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &StallRepository{Client: s.client}
	s.ctx = context.Background()
}

func (s *StallRepositoryTestSuite) newStall(tag string, merchantID, storeID uuid.UUID) *domain.Stall {
	return &domain.Stall{
		ID:         uuid.New(),
		Name:       "出品-" + tag,
		StallType:  domain.StallTypeStore,
		PrintType:  domain.StallPrintTypeReceipt,
		Enabled:    true,
		SortOrder:  10,
		MerchantID: merchantID,
		StoreID:    storeID,
	}
}

func (s *StallRepositoryTestSuite) TestStall_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		stall := s.newStall("create", uuid.New(), uuid.New())

		err := s.repo.Create(s.ctx, stall)
		require.NoError(t, err)

		saved := s.client.Stall.GetX(s.ctx, stall.ID)
		require.Equal(t, stall.Name, saved.Name)
		require.Equal(t, stall.StallType, saved.StallType)
		require.Equal(t, stall.MerchantID, saved.MerchantID)
		require.Equal(t, stall.StoreID, saved.StoreID)
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Create(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *StallRepositoryTestSuite) TestStall_FindByID() {
	stall := s.newStall("find", uuid.New(), uuid.New())
	require.NoError(s.T(), s.repo.Create(s.ctx, stall))

	s.T().Run("查询成功", func(t *testing.T) {
		found, err := s.repo.FindByID(s.ctx, stall.ID)
		require.NoError(t, err)
		require.Equal(t, stall.ID, found.ID)
		require.Equal(t, stall.Name, found.Name)
	})

	s.T().Run("不存在", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *StallRepositoryTestSuite) TestStall_Update() {
	stall := s.newStall("update", uuid.New(), uuid.New())
	require.NoError(s.T(), s.repo.Create(s.ctx, stall))

	s.T().Run("更新成功", func(t *testing.T) {
		stall.Name = "更新-" + stall.Name
		stall.Enabled = false
		stall.SortOrder = 1

		err := s.repo.Update(s.ctx, stall)
		require.NoError(t, err)

		updated := s.client.Stall.GetX(s.ctx, stall.ID)
		require.Equal(t, stall.Name, updated.Name)
		require.Equal(t, stall.Enabled, updated.Enabled)
		require.Equal(t, stall.SortOrder, updated.SortOrder)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		missing := s.newStall("missing", uuid.New(), uuid.New())
		err := s.repo.Update(s.ctx, missing)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Update(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *StallRepositoryTestSuite) TestStall_Delete() {
	stall := s.newStall("delete", uuid.New(), uuid.New())
	require.NoError(s.T(), s.repo.Create(s.ctx, stall))

	s.T().Run("删除成功", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, stall.ID)
		require.NoError(t, err)
		_, err = s.client.Stall.Get(s.ctx, stall.ID)
		require.Error(t, err)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *StallRepositoryTestSuite) TestStall_GetStalls() {
	m1 := uuid.New()
	m2 := uuid.New()
	stall1 := s.newStall("001", m1, uuid.New())
	require.NoError(s.T(), s.repo.Create(s.ctx, stall1))
	time.Sleep(10 * time.Millisecond)
	stall2 := s.newStall("002", m1, uuid.New())
	stall2.Enabled = false
	stall2.SortOrder = 2
	require.NoError(s.T(), s.repo.Create(s.ctx, stall2))
	time.Sleep(10 * time.Millisecond)
	stall3 := s.newStall("003", m2, uuid.New())
	stall3.PrintType = domain.StallPrintTypeLabel
	require.NoError(s.T(), s.repo.Create(s.ctx, stall3))

	pager := upagination.New(1, 10)

	s.T().Run("按商户筛选默认排序", func(t *testing.T) {
		list, total, err := s.repo.GetStalls(s.ctx, pager, &domain.StallListFilter{MerchantID: m1})
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, list, 2)
		require.Equal(t, stall2.ID, list[0].ID)
	})

	s.T().Run("按排序字段升序", func(t *testing.T) {
		order := domain.NewStallOrderBySortOrder(false)
		list, total, err := s.repo.GetStalls(s.ctx, pager, &domain.StallListFilter{MerchantID: m1}, order)
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Equal(t, stall2.ID, list[0].ID)
	})

	s.T().Run("按打印类型与启用状态筛选", func(t *testing.T) {
		enable := true
		list, total, err := s.repo.GetStalls(s.ctx, pager, &domain.StallListFilter{MerchantID: m2, PrintType: domain.StallPrintTypeLabel, Enabled: &enable})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, stall3.ID, list[0].ID)
	})

	s.T().Run("按名称模糊筛选", func(t *testing.T) {
		list, total, err := s.repo.GetStalls(s.ctx, pager, &domain.StallListFilter{Name: "003"})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, stall3.ID, list[0].ID)
	})
}

func (s *StallRepositoryTestSuite) TestStall_Exists() {
	merchantID := uuid.New()
	storeID := uuid.New()
	stall := s.newStall("exists", merchantID, storeID)
	require.NoError(s.T(), s.repo.Create(s.ctx, stall))

	s.T().Run("同名存在", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.StallExistsParams{Name: stall.Name, MerchantID: merchantID, StoreID: storeID})
		require.NoError(t, err)
		require.True(t, exists)
	})

	s.T().Run("排除自身", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.StallExistsParams{Name: stall.Name, MerchantID: merchantID, StoreID: storeID, ExcludeID: stall.ID})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("不同商户不冲突", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.StallExistsParams{Name: stall.Name, MerchantID: uuid.New(), StoreID: storeID})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("不同门店不冲突", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.StallExistsParams{Name: stall.Name, MerchantID: merchantID, StoreID: uuid.New()})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("仅名称查询", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.StallExistsParams{Name: stall.Name})
		require.NoError(t, err)
		require.True(t, exists)
	})
}
