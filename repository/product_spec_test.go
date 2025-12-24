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

type ProductSpecTestSuite struct {
	RepositoryTestSuite
	repo *ProductSpecRepository
	ctx  context.Context
}

func TestProductSpecTestSuite(t *testing.T) {
	suite.Run(t, new(ProductSpecTestSuite))
}

func (s *ProductSpecTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &ProductSpecRepository{
		Client: s.client,
	}
	s.ctx = context.Background()
}

func (s *ProductSpecTestSuite) createTestProductSpec(
	merchantID, storeID uuid.UUID,
	name string,
) *domain.ProductSpec {
	spec := &domain.ProductSpec{
		ID:           uuid.New(),
		Name:         name,
		MerchantID:   merchantID,
		StoreID:      storeID,
		ProductCount: 0,
	}
	err := s.repo.Create(s.ctx, spec)
	require.NoError(s.T(), err)
	return spec
}

func (s *ProductSpecTestSuite) TestProductSpec_Create() {
	s.T().Run("创建商品规格成功", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()

		spec := &domain.ProductSpec{
			ID:           uuid.New(),
			Name:         "大",
			MerchantID:   merchantID,
			StoreID:      storeID,
			ProductCount: 0,
		}

		err := s.repo.Create(s.ctx, spec)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, spec.ID)
		require.NotZero(t, spec.CreatedAt)

		// 验证数据库记录
		dbSpec := s.client.ProductSpec.GetX(s.ctx, spec.ID)
		require.Equal(t, "大", dbSpec.Name)
		require.Equal(t, merchantID, dbSpec.MerchantID)
		require.Equal(t, storeID, dbSpec.StoreID)
		require.Equal(t, 0, dbSpec.ProductCount)
	})
}

func (s *ProductSpecTestSuite) TestProductSpec_FindByID() {
	s.T().Run("查找存在的商品规格", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()
		spec := s.createTestProductSpec(merchantID, storeID, "大")

		found, err := s.repo.FindByID(s.ctx, spec.ID)
		require.NoError(t, err)
		require.NotNil(t, found)
		require.Equal(t, spec.ID, found.ID)
		require.Equal(t, "大", found.Name)
	})

	s.T().Run("查找不存在的商品规格", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *ProductSpecTestSuite) TestProductSpec_Update() {
	s.T().Run("更新商品规格成功", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()
		spec := s.createTestProductSpec(merchantID, storeID, "大")

		spec.Name = "中"

		err := s.repo.Update(s.ctx, spec)
		require.NoError(t, err)

		// 验证更新
		updated := s.client.ProductSpec.GetX(s.ctx, spec.ID)
		require.Equal(t, "中", updated.Name)
	})
}

func (s *ProductSpecTestSuite) TestProductSpec_Delete() {
	s.T().Run("删除商品规格成功", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()
		spec := s.createTestProductSpec(merchantID, storeID, "大")

		err := s.repo.Delete(s.ctx, spec.ID)
		require.NoError(t, err)

		// 验证已删除
		_, err = s.client.ProductSpec.Get(s.ctx, spec.ID)
		require.Error(t, err)
	})
}

func (s *ProductSpecTestSuite) TestProductSpec_Exists() {
	s.T().Run("检查名称是否存在", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()
		spec := s.createTestProductSpec(merchantID, storeID, "大")

		// 检查相同门店下的相同名称
		exists, err := s.repo.Exists(s.ctx, domain.ProductSpecExistsParams{
			MerchantID: merchantID,
			Name:       "大",
		})
		require.NoError(t, err)
		require.True(t, exists)

		// 检查不同门店下的相同名称
		exists, err = s.repo.Exists(s.ctx, domain.ProductSpecExistsParams{
			MerchantID: uuid.New(),
			Name:       "大",
		})
		require.NoError(t, err)
		require.False(t, exists)

		// 更新时排除自身ID
		exists, err = s.repo.Exists(s.ctx, domain.ProductSpecExistsParams{
			MerchantID: merchantID,
			Name:       "大",
			ExcludeID:  spec.ID,
		})
		require.NoError(t, err)
		require.False(t, exists)
	})
}

func (s *ProductSpecTestSuite) TestProductSpec_PagedListBySearch() {
	s.T().Run("查询商品规格列表", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()

		// 创建多个商品规格
		s.createTestProductSpec(merchantID, storeID, "大")
		s.createTestProductSpec(merchantID, storeID, "中")
		s.createTestProductSpec(merchantID, storeID, "小")

		page := upagination.New(1, 10)

		// 查询列表
		res, err := s.repo.PagedListBySearch(s.ctx, page, domain.ProductSpecSearchParams{
			MerchantID: merchantID,
		})
		require.NoError(t, err)
		require.Equal(t, 3, len(res.Items))
		require.Equal(t, 3, res.Total)

		// 验证按创建时间倒序排列（后创建的在前）
		require.Equal(t, "小", res.Items[0].Name)
		require.Equal(t, "中", res.Items[1].Name)
		require.Equal(t, "大", res.Items[2].Name)
	})

	s.T().Run("分页查询", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()

		// 创建5个商品规格
		for i := 0; i < 5; i++ {
			s.createTestProductSpec(merchantID, storeID, "规格"+string(rune('A'+i)))
		}

		page := upagination.New(1, 2)

		// 第一页
		res, err := s.repo.PagedListBySearch(s.ctx, page, domain.ProductSpecSearchParams{
			MerchantID: merchantID,
		})
		require.NoError(t, err)
		require.Equal(t, 5, res.Total)
		require.Len(t, res.Items, 2)

		// 第二页
		page = upagination.New(2, 2)
		res, err = s.repo.PagedListBySearch(s.ctx, page, domain.ProductSpecSearchParams{
			MerchantID: merchantID,
		})
		require.NoError(t, err)
		require.Equal(t, 5, res.Total)
		require.Len(t, res.Items, 2)
	})
}
