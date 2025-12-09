package repository

import (
	"context"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
)

type ProductTestSuite struct {
	RepositoryTestSuite
	repo *ProductRepository
	ctx  context.Context
}

func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}

func (s *ProductTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &ProductRepository{
		Client: s.client,
	}
	s.ctx = context.Background()
}

func (s *ProductTestSuite) TestProduct_FindByID() {
	s.T().Run("正常查询", func(t *testing.T) {
		product, err := s.repo.FindByID(s.ctx, 16)
		require.NoError(t, err)
		util.PrettyJson(product)
	})
}

func (s *ProductTestSuite) createTestProduct(name string) *ent.Product {
	storeID := 1
	// 1. 创建商品分类
	category := s.client.Category.Create().
		SetName("测试分类").
		SetStoreID(storeID).
		SaveX(s.ctx)

	// 3. 创建单位
	unit := s.client.Unit.Create().
		SetName("个").
		SetStoreID(storeID).
		SaveX(s.ctx)

	return s.client.Product.Create().
		SetName(name).
		SetStoreID(storeID).
		SetCategoryID(category.ID).
		SetUnitID(unit.ID).
		SetPrice(decimal.NewFromInt(100)).
		SaveX(s.ctx)
}
