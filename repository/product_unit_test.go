package repository

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
	"testing"
)

type ProductUnitTestSuite struct {
	RepositoryTestSuite
	repo *ProductUnitRepository
	ctx  context.Context
}

func TestProductUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ProductUnitTestSuite))
}

func (s *ProductUnitTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &ProductUnitRepository{
		Client: s.client,
	}
	s.ctx = context.Background()
}

func (s *ProductUnitTestSuite) TestProductUnit_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		unit := &domain.ProductUnit{
			Name:    "盘",
			StoreID: 1,
		}
		err := s.repo.Create(s.ctx, unit)
		require.NoError(t, err)
		require.NotZero(t, unit.ID)

		// 验证数据库记录
		dbUnit := s.client.Unit.GetX(s.ctx, unit.ID)
		require.Equal(t, "盘", dbUnit.Name)
	})
}

func (s *ProductUnitTestSuite) createTestUnit() *ent.Unit {
	return s.client.Unit.Create().
		SetName("测试单位").
		SetStoreID(1).
		SaveX(s.ctx)
}

func (s *ProductUnitTestSuite) TestProductUnit_FindByID() {
	pu := s.createTestUnit()

	s.T().Run("正常查询", func(t *testing.T) {
		unit, err := s.repo.FindByID(s.ctx, 1)
		require.NoError(t, err)
		require.Equal(t, pu.ID, unit.ID)
		require.Equal(t, "测试单位", unit.Name)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, 9999)
		require.Error(t, err)
	})
}

func (s *ProductUnitTestSuite) TestProductUnit_Update() {
	pu := s.createTestUnit()
	s.T().Run("正常更新", func(t *testing.T) {
		unit := &domain.ProductUnit{
			ID:   pu.ID,
			Name: "更新后的名称",
		}

		err := s.repo.Update(s.ctx, unit)
		require.NoError(t, err)

		// 验证更新结果
		updated := s.client.Unit.GetX(s.ctx, pu.ID)
		require.Equal(t, "更新后的名称", updated.Name)
	})

	s.T().Run("更新不存在的ID", func(t *testing.T) {
		unit := &domain.ProductUnit{
			ID:   9999,
			Name: "无效单位",
		}

		err := s.repo.Update(s.ctx, unit)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *ProductUnitTestSuite) TestProductUnit_Delete() {
	pu := s.createTestUnit()

	s.T().Run("正常删除", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, pu.ID)
		require.NoError(t, err)

		deleted, err := s.client.Unit.Get(s.ctx, pu.ID)
		require.Nil(t, deleted)
		require.Error(t, err)
		require.True(t, ent.IsNotFound(err))
	})

	s.T().Run("删除不存在的ID", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, 9999)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *ProductUnitTestSuite) TestProductUnit_Exists() {
	pu := s.createTestUnit()
	// 测试用例
	s.T().Run("存在记录", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.UnitExistsParams{
			StoreID: pu.StoreID,
			Name:    pu.Name,
		})
		s.NoError(err)
		s.True(exists)
	})
}

func (s *ProductUnitTestSuite) TestProductUnit_PagedListBySearch() {
	// 创建测试数据
	testUnits := []struct {
		name    string
		storeID int
	}{
		{"单位A", 1},
		{"单位B", 1},
		{"单位C", 1},
		{"单位D", 2},
		{"单位E", 2},
	}

	for _, u := range testUnits {
		s.client.Unit.Create().
			SetName(u.name).
			SetStoreID(u.storeID).
			SaveX(s.ctx)
	}

	s.T().Run("默认ID倒序", func(t *testing.T) {
		page := upagination.New(1, 10)
		res, err := s.repo.PagedListBySearch(s.ctx, page, domain.UnitSearchParams{
			StoreID: 1,
		})
		s.NoError(err)
		s.Equal(3, res.Total)
		s.Equal(3, len(res.Items))
		util.PrettyJson(res)
	})
}
