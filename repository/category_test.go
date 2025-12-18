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

type CategoryTestSuite struct {
	RepositoryTestSuite
	repo *CategoryRepository
	ctx  context.Context
}

func TestCategoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryTestSuite))
}

func (s *CategoryTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &CategoryRepository{
		Client: s.client,
	}
	s.ctx = context.Background()
}

func (s *CategoryTestSuite) createTestCategory(parentID uuid.UUID) *ent.Category {
	storeID := uuid.New()
	taxRateID := uuid.New()
	stallID := uuid.New()

	builder := s.client.Category.Create().
		SetID(uuid.New()).
		SetName("测试分类").
		SetStoreID(storeID).
		SetInheritTaxRate(false).
		SetInheritStall(false).
		SetSortOrder(0)

	if parentID != uuid.Nil {
		builder = builder.SetParentID(parentID)
	}

	if taxRateID != uuid.Nil {
		builder = builder.SetTaxRateID(taxRateID)
	}

	if stallID != uuid.Nil {
		builder = builder.SetStallID(stallID)
	}

	return builder.SaveX(s.ctx)
}

func (s *CategoryTestSuite) createTestRootCategory() *ent.Category {
	return s.createTestCategory(uuid.Nil)
}

func (s *CategoryTestSuite) TestCategory_Create() {
	s.T().Run("创建一级分类成功", func(t *testing.T) {
		storeID := uuid.New()
		taxRateID := uuid.New()
		stallID := uuid.New()

		category := &domain.Category{
			ID:             uuid.New(),
			Name:           "一级分类",
			StoreID:        storeID,
			ParentID:       uuid.Nil, // 一级分类
			InheritTaxRate: false,
			TaxRateID:      taxRateID,
			InheritStall:   false,
			StallID:        stallID,
			SortOrder:      0,
		}

		err := s.repo.Create(s.ctx, category)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, category.ID)
		require.True(t, category.IsRoot())

		// 验证数据库记录
		dbCategory := s.client.Category.GetX(s.ctx, category.ID)
		require.Equal(t, "一级分类", dbCategory.Name)
		require.Equal(t, storeID, dbCategory.StoreID)
		require.Equal(t, uuid.Nil, dbCategory.ParentID)
		require.Equal(t, taxRateID, dbCategory.TaxRateID)
		require.Equal(t, stallID, dbCategory.StallID)
	})

	s.T().Run("创建二级分类成功", func(t *testing.T) {
		// 先创建一级分类
		rootCategory := s.createTestRootCategory()

		storeID := uuid.New()

		category := &domain.Category{
			ID:             uuid.New(),
			Name:           "二级分类",
			StoreID:        storeID,
			ParentID:       rootCategory.ID, // 二级分类
			InheritTaxRate: true,            // 继承父分类税率
			TaxRateID:      uuid.Nil,
			InheritStall:   true, // 继承父分类档口
			StallID:        uuid.Nil,
			SortOrder:      1,
		}

		err := s.repo.Create(s.ctx, category)
		require.NoError(t, err)
		require.False(t, category.IsRoot())

		// 验证数据库记录
		dbCategory := s.client.Category.GetX(s.ctx, category.ID)
		require.Equal(t, "二级分类", dbCategory.Name)
		require.Equal(t, rootCategory.ID, dbCategory.ParentID)
		require.True(t, dbCategory.InheritTaxRate)
		require.True(t, dbCategory.InheritStall)
	})
}

func (s *CategoryTestSuite) TestCategory_FindByID() {
	category := s.createTestRootCategory()

	s.T().Run("正常查询", func(t *testing.T) {
		found, err := s.repo.FindByID(s.ctx, category.ID)

		util.PrettyJson(found)
		require.NoError(t, err)
		require.Equal(t, category.ID, found.ID)
		require.Equal(t, "测试分类", found.Name)
		require.True(t, found.IsRoot())
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		nonExistentID := uuid.New()
		_, err := s.repo.FindByID(s.ctx, nonExistentID)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *CategoryTestSuite) TestCategory_Delete() {
	category := s.createTestRootCategory()

	s.T().Run("正常删除", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, category.ID)
		require.NoError(t, err)

		// 验证删除结果
		deleted, err := s.client.Category.Get(s.ctx, category.ID)
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

func (s *CategoryTestSuite) TestCategory_Update() {
	category := s.createTestRootCategory()

	s.T().Run("正常更新", func(t *testing.T) {
		newTaxRateID := uuid.New()
		newStallID := uuid.New()

		cat := &domain.Category{
			ID:             category.ID,
			Name:           "更新后的分类",
			StoreID:        category.StoreID,
			ParentID:       category.ParentID,
			InheritTaxRate: false,
			TaxRateID:      newTaxRateID,
			InheritStall:   false,
			StallID:        newStallID,
			SortOrder:      10,
			ProductCount:   5,
		}

		err := s.repo.Update(s.ctx, cat)
		require.NoError(t, err)

		// 验证更新结果
		updated := s.client.Category.GetX(s.ctx, category.ID)
		require.Equal(t, "更新后的分类", updated.Name)
		require.Equal(t, newTaxRateID, updated.TaxRateID)
		require.Equal(t, newStallID, updated.StallID)
		require.Equal(t, 10, updated.SortOrder)
		require.Equal(t, 5, updated.ProductCount)
	})

	s.T().Run("更新继承字段", func(t *testing.T) {
		// 创建二级分类
		rootCategory := s.createTestRootCategory()
		childCategory := s.createTestCategory(rootCategory.ID)

		cat := &domain.Category{
			ID:             childCategory.ID,
			Name:           childCategory.Name,
			StoreID:        childCategory.StoreID,
			ParentID:       childCategory.ParentID,
			InheritTaxRate: true, // 改为继承
			TaxRateID:      uuid.Nil,
			InheritStall:   true, // 改为继承
			StallID:        uuid.Nil,
			SortOrder:      childCategory.SortOrder,
			ProductCount:   childCategory.ProductCount,
		}

		err := s.repo.Update(s.ctx, cat)
		require.NoError(t, err)

		// 验证更新结果
		updated := s.client.Category.GetX(s.ctx, childCategory.ID)
		require.True(t, updated.InheritTaxRate)
		require.True(t, updated.InheritStall)
		require.Equal(t, uuid.Nil, updated.TaxRateID)
		require.Equal(t, uuid.Nil, updated.StallID)
	})

	s.T().Run("更新不存在的ID", func(t *testing.T) {
		nonExistentID := uuid.New()
		cat := &domain.Category{
			ID:       nonExistentID,
			Name:     "无效分类",
			StoreID:  uuid.New(),
			ParentID: uuid.Nil,
		}

		err := s.repo.Update(s.ctx, cat)
		require.Error(t, err)
	})
}
