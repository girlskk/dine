package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
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

func (s *ProductUnitTestSuite) createTestProductUnit(
	merchantID, storeID uuid.UUID,
	name string, unitType domain.ProductUnitType,
) *domain.ProductUnit {
	unit := &domain.ProductUnit{
		ID:           uuid.New(),
		Name:         name,
		Type:         unitType,
		MerchantID:   merchantID,
		StoreID:      storeID,
		ProductCount: 0,
	}
	err := s.repo.Create(s.ctx, unit)
	require.NoError(s.T(), err)
	return unit
}

func (s *ProductUnitTestSuite) TestProductUnit_Create() {
	s.T().Run("创建商品单位成功", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()

		unit := &domain.ProductUnit{
			ID:           uuid.New(),
			Name:         "个",
			Type:         domain.ProductUnitTypeQuantity,
			MerchantID:   merchantID,
			StoreID:      storeID,
			ProductCount: 0,
		}

		err := s.repo.Create(s.ctx, unit)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, unit.ID)
		require.NotZero(t, unit.CreatedAt)

		// 验证数据库记录
		dbUnit := s.client.ProductUnit.GetX(s.ctx, unit.ID)
		require.Equal(t, "个", dbUnit.Name)
		require.Equal(t, domain.ProductUnitTypeQuantity, dbUnit.Type)
		require.Equal(t, merchantID, dbUnit.MerchantID)
		require.Equal(t, storeID, dbUnit.StoreID)
		require.Equal(t, 0, dbUnit.ProductCount)
	})
}

func (s *ProductUnitTestSuite) TestProductUnit_FindByID() {
	s.T().Run("查找存在的商品单位", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()
		unit := s.createTestProductUnit(merchantID, storeID, "个", domain.ProductUnitTypeQuantity)

		found, err := s.repo.FindByID(s.ctx, unit.ID)
		require.NoError(t, err)
		require.NotNil(t, found)
		require.Equal(t, unit.ID, found.ID)
		require.Equal(t, "个", found.Name)
		require.Equal(t, domain.ProductUnitTypeQuantity, found.Type)
	})

	s.T().Run("查找不存在的商品单位", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *ProductUnitTestSuite) TestProductUnit_Update() {
	s.T().Run("更新商品单位成功", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()
		unit := s.createTestProductUnit(merchantID, storeID, "个", domain.ProductUnitTypeQuantity)

		unit.Name = "件"
		unit.Type = domain.ProductUnitTypeWeight

		err := s.repo.Update(s.ctx, unit)
		require.NoError(t, err)

		// 验证更新
		updated := s.client.ProductUnit.GetX(s.ctx, unit.ID)
		require.Equal(t, "件", updated.Name)
		require.Equal(t, domain.ProductUnitTypeWeight, updated.Type)
	})
}

func (s *ProductUnitTestSuite) TestProductUnit_Delete() {
	s.T().Run("删除商品单位成功", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()
		unit := s.createTestProductUnit(merchantID, storeID, "个", domain.ProductUnitTypeQuantity)

		err := s.repo.Delete(s.ctx, unit.ID)
		require.NoError(t, err)

		// 验证已删除
		_, err = s.client.ProductUnit.Get(s.ctx, unit.ID)
		require.Error(t, err)
	})
}

func (s *ProductUnitTestSuite) TestProductUnit_Exists() {
	s.T().Run("检查名称是否存在", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()
		unit := s.createTestProductUnit(merchantID, storeID, "个", domain.ProductUnitTypeQuantity)

		// 检查相同品牌商下的相同名称
		exists, err := s.repo.Exists(s.ctx, domain.ProductUnitExistsParams{
			MerchantID: merchantID,
			Name:       "个",
		})
		require.NoError(t, err)
		require.True(t, exists)

		// 检查不同品牌商下的相同名称（应该不存在）
		otherMerchantID := uuid.New()
		exists, err = s.repo.Exists(s.ctx, domain.ProductUnitExistsParams{
			MerchantID: otherMerchantID,
			Name:       "个",
		})
		require.NoError(t, err)
		require.False(t, exists)

		// 更新时排除自身ID（应该不存在）
		exists, err = s.repo.Exists(s.ctx, domain.ProductUnitExistsParams{
			MerchantID: merchantID,
			Name:       "个",
			ExcludeID:  unit.ID,
		})
		require.NoError(t, err)
		require.False(t, exists)
	})
}

func (s *ProductUnitTestSuite) TestProductUnit_ListBySearch() {
	s.T().Run("查询所有商品单位列表", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()

		// 创建多个商品单位
		s.createTestProductUnit(merchantID, storeID, "个", domain.ProductUnitTypeQuantity)
		s.createTestProductUnit(merchantID, storeID, "件", domain.ProductUnitTypeQuantity)
		s.createTestProductUnit(merchantID, storeID, "公斤", domain.ProductUnitTypeWeight)

		page := upagination.New(1, 10)

		// 查询所有列表（不指定Name）
		units, err := s.repo.PagedListBySearch(s.ctx, page, domain.ProductUnitSearchParams{
			MerchantID: merchantID,
		})
		require.NoError(t, err)
		require.Equal(t, 3, page.Total)
		require.Len(t, units, 3)
	})

	s.T().Run("按名称模糊查询", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()

		// 创建多个商品单位
		s.createTestProductUnit(merchantID, storeID, "个", domain.ProductUnitTypeQuantity)
		s.createTestProductUnit(merchantID, storeID, "件", domain.ProductUnitTypeQuantity)
		s.createTestProductUnit(merchantID, storeID, "公斤", domain.ProductUnitTypeWeight)

		page := upagination.New(1, 10)

		// 查询名称包含"个"的单位（应该只返回"个"）
		units, err := s.repo.PagedListBySearch(s.ctx, page, domain.ProductUnitSearchParams{
			MerchantID: merchantID,
			Name:       "个",
		})
		require.NoError(t, err)
		require.Equal(t, 1, page.Total)
		require.Len(t, units, 1)
		require.Equal(t, "个", units.Items[0].Name)
	})

	s.T().Run("按类型查询", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()

		// 创建多个商品单位
		s.createTestProductUnit(merchantID, storeID, "个", domain.ProductUnitTypeQuantity)
		s.createTestProductUnit(merchantID, storeID, "件", domain.ProductUnitTypeQuantity)
		s.createTestProductUnit(merchantID, storeID, "公斤", domain.ProductUnitTypeWeight)

		page := upagination.New(1, 10)

		// 查询数量类型的单位
		units, err := s.repo.PagedListBySearch(s.ctx, page, domain.ProductUnitSearchParams{
			MerchantID: merchantID,
			Type:       domain.ProductUnitTypeQuantity,
		})
		require.NoError(t, err)
		require.Equal(t, 2, page.Total)
		require.Len(t, units, 2)
		for _, unit := range units.Items {
			require.Equal(t, domain.ProductUnitTypeQuantity, unit.Type)
		}
	})

	s.T().Run("分页查询", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()

		// 创建10个商品单位
		for i := 0; i < 10; i++ {
			s.createTestProductUnit(merchantID, storeID, "单位"+string(rune('A'+i)), domain.ProductUnitTypeQuantity)
		}

		// 第一页，每页3条
		page := upagination.New(1, 3)
		units, err := s.repo.PagedListBySearch(s.ctx, page, domain.ProductUnitSearchParams{
			MerchantID: merchantID,
		})
		require.NoError(t, err)
		require.Equal(t, 10, page.Total)
		require.Len(t, units, 3)

		// 第二页
		page2 := upagination.New(2, 3)
		units2, err := s.repo.PagedListBySearch(s.ctx, page2, domain.ProductUnitSearchParams{
			MerchantID: merchantID,
		})
		require.NoError(t, err)
		require.Equal(t, 10, page2.Total)
		require.Len(t, units2, 3)

		// 确保两页的数据不同
		require.NotEqual(t, units.Items[0].ID, units2.Items[0].ID)
	})
}
