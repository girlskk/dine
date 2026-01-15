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

type ProductTagTestSuite struct {
	RepositoryTestSuite
	repo *ProductTagRepository
	ctx  context.Context
}

func TestProductTagTestSuite(t *testing.T) {
	suite.Run(t, new(ProductTagTestSuite))
}

func (s *ProductTagTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &ProductTagRepository{
		Client: s.client,
	}
	s.ctx = context.Background()
}

func (s *ProductTagTestSuite) createTestProductTag(
	merchantID, storeID uuid.UUID,
	name string,
) *domain.ProductTag {
	tag := &domain.ProductTag{
		ID:           uuid.New(),
		Name:         name,
		MerchantID:   merchantID,
		StoreID:      storeID,
		ProductCount: 0,
	}
	err := s.repo.Create(s.ctx, tag)
	require.NoError(s.T(), err)
	return tag
}

func (s *ProductTagTestSuite) TestProductTag_Create() {
	s.T().Run("创建商品标签成功", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()

		tag := &domain.ProductTag{
			ID:           uuid.New(),
			Name:         "热销",
			MerchantID:   merchantID,
			StoreID:      storeID,
			ProductCount: 0,
		}

		err := s.repo.Create(s.ctx, tag)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, tag.ID)
		require.NotZero(t, tag.CreatedAt)

		// 验证数据库记录
		dbTag := s.client.ProductTag.GetX(s.ctx, tag.ID)
		require.Equal(t, "热销", dbTag.Name)
		require.Equal(t, merchantID, dbTag.MerchantID)
		require.Equal(t, storeID, dbTag.StoreID)
		require.Equal(t, 0, dbTag.ProductCount)
	})
}

func (s *ProductTagTestSuite) TestProductTag_FindByID() {
	s.T().Run("查找存在的商品标签", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()
		tag := s.createTestProductTag(merchantID, storeID, "热销")

		found, err := s.repo.FindByID(s.ctx, tag.ID)
		require.NoError(t, err)
		require.NotNil(t, found)
		require.Equal(t, tag.ID, found.ID)
		require.Equal(t, "热销", found.Name)
	})

	s.T().Run("查找不存在的商品标签", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *ProductTagTestSuite) TestProductTag_Update() {
	s.T().Run("更新商品标签成功", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()
		tag := s.createTestProductTag(merchantID, storeID, "热销")

		tag.Name = "新品"

		err := s.repo.Update(s.ctx, tag)
		require.NoError(t, err)

		// 验证更新
		updated := s.client.ProductTag.GetX(s.ctx, tag.ID)
		require.Equal(t, "新品", updated.Name)
	})
}

func (s *ProductTagTestSuite) TestProductTag_Delete() {
	s.T().Run("删除商品标签成功", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()
		tag := s.createTestProductTag(merchantID, storeID, "热销")

		err := s.repo.Delete(s.ctx, tag.ID)
		require.NoError(t, err)

		// 验证已删除
		_, err = s.client.ProductTag.Get(s.ctx, tag.ID)
		require.Error(t, err)
	})
}

func (s *ProductTagTestSuite) TestProductTag_Exists() {
	s.T().Run("检查名称是否存在", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()
		tag := s.createTestProductTag(merchantID, storeID, "热销")

		// 检查相同门店下的相同名称
		exists, err := s.repo.Exists(s.ctx, domain.ProductTagExistsParams{
			MerchantID: merchantID,
			Name:       "热销",
		})
		require.NoError(t, err)
		require.True(t, exists)

		// 检查不同门店下的相同名称
		exists, err = s.repo.Exists(s.ctx, domain.ProductTagExistsParams{
			MerchantID: uuid.New(),
			Name:       "热销",
		})
		require.NoError(t, err)
		require.False(t, exists)

		// 更新时排除自身ID
		exists, err = s.repo.Exists(s.ctx, domain.ProductTagExistsParams{
			MerchantID: merchantID,
			Name:       "热销",
			ExcludeID:  tag.ID,
		})
		require.NoError(t, err)
		require.False(t, exists)
	})
}

func (s *ProductTagTestSuite) TestProductTag_PagedListBySearch() {
	s.T().Run("查询商品标签列表", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()

		// 创建多个商品标签
		s.createTestProductTag(merchantID, storeID, "热销")
		s.createTestProductTag(merchantID, storeID, "新品")
		s.createTestProductTag(merchantID, storeID, "推荐")

		page := upagination.New(1, 10)

		// 查询列表
		res, err := s.repo.PagedListBySearch(s.ctx, page, domain.ProductTagSearchParams{
			MerchantID: merchantID,
		})
		require.NoError(t, err)
		require.Equal(t, 3, len(res.Items))
		require.Equal(t, 3, res.Total)

		// 验证按创建时间倒序排列（后创建的在前）
		require.Equal(t, "推荐", res.Items[0].Name)
		require.Equal(t, "新品", res.Items[1].Name)
		require.Equal(t, "热销", res.Items[2].Name)
	})

	s.T().Run("分页查询", func(t *testing.T) {
		merchantID := uuid.New()
		storeID := uuid.New()

		// 创建5个商品标签
		for i := 0; i < 5; i++ {
			s.createTestProductTag(merchantID, storeID, "标签"+string(rune('A'+i)))
		}

		page := upagination.New(1, 2)

		// 第一页
		res, err := s.repo.PagedListBySearch(s.ctx, page, domain.ProductTagSearchParams{
			MerchantID: merchantID,
		})
		require.NoError(t, err)
		require.Equal(t, 5, res.Total)
		require.Len(t, res.Items, 2)

		// 第二页
		page = upagination.New(2, 2)
		res, err = s.repo.PagedListBySearch(s.ctx, page, domain.ProductTagSearchParams{
			MerchantID: merchantID,
		})
		require.NoError(t, err)
		require.Equal(t, 5, res.Total)
		require.Len(t, res.Items, 2)
	})
}
