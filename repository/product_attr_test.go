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

type ProductAttrTestSuite struct {
	RepositoryTestSuite
	repo *ProductAttrRepository
	ctx  context.Context
}

func TestProductAttrTestSuite(t *testing.T) {
	suite.Run(t, new(ProductAttrTestSuite))
}

func (s *ProductAttrTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &ProductAttrRepository{
		Client: s.client,
	}
	s.ctx = context.Background()
}

func (s *ProductAttrTestSuite) TestProductAttr_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		attr := &domain.ProductAttr{
			Name:    "微微辣",
			StoreID: 1,
		}
		err := s.repo.Create(s.ctx, attr)
		require.NoError(t, err)
		require.NotZero(t, attr.ID)

		// 验证数据库记录
		dbAttr := s.client.Attr.GetX(s.ctx, attr.ID)
		require.Equal(t, "微微辣", dbAttr.Name)
	})
}

func (s *ProductAttrTestSuite) createTestAttr() *ent.Attr {
	return s.client.Attr.Create().
		SetName("测试属性").
		SetStoreID(1).
		SaveX(s.ctx)
}

func (s *ProductAttrTestSuite) TestProductAttr_FindByID() {
	pu := s.createTestAttr()

	s.T().Run("正常查询", func(t *testing.T) {
		Attr, err := s.repo.FindByID(s.ctx, 1)
		require.NoError(t, err)
		require.Equal(t, pu.ID, Attr.ID)
		require.Equal(t, "测试属性", Attr.Name)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, 9999)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *ProductAttrTestSuite) TestProductAttr_Update() {
	pu := s.createTestAttr()
	s.T().Run("正常更新", func(t *testing.T) {
		Attr := &domain.ProductAttr{
			ID:   pu.ID,
			Name: "更新后的名称",
		}

		err := s.repo.Update(s.ctx, Attr)
		require.NoError(t, err)

		// 验证更新结果
		updated := s.client.Attr.GetX(s.ctx, pu.ID)
		require.Equal(t, "更新后的名称", updated.Name)
	})

	s.T().Run("更新不存在的ID", func(t *testing.T) {
		Attr := &domain.ProductAttr{
			ID:   9999,
			Name: "无效属性",
		}

		err := s.repo.Update(s.ctx, Attr)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *ProductAttrTestSuite) TestProductAttr_Delete() {
	pu := s.createTestAttr()

	s.T().Run("正常删除", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, pu.ID)
		require.NoError(t, err)

		deleted, err := s.client.Attr.Get(s.ctx, pu.ID)
		require.Nil(t, deleted)
		require.Error(t, err)
		require.True(t, ent.IsNotFound(err))
	})

	s.T().Run("删除不存在的ID", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, 9999)
		require.Error(t, err)
	})
}

func (s *ProductAttrTestSuite) TestProductAttr_Exists() {
	pu := s.createTestAttr()
	// 测试用例
	s.T().Run("存在记录", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.AttrExistsParams{
			StoreID: pu.StoreID,
			Name:    pu.Name,
		})
		s.NoError(err)
		s.True(exists)
	})
}

func (s *ProductAttrTestSuite) TestProductAttr_PagedListBySearch() {
	// 创建测试数据
	testAttrs := []struct {
		name    string
		storeID int
	}{
		{"属性A", 1},
		{"属性B", 1},
		{"属性C", 1},
		{"属性D", 2},
		{"属性E", 2},
	}

	for _, u := range testAttrs {
		s.client.Attr.Create().
			SetName(u.name).
			SetStoreID(u.storeID).
			SaveX(s.ctx)
	}

	s.T().Run("默认ID倒序", func(t *testing.T) {
		page := upagination.New(1, 10)
		res, err := s.repo.PagedListBySearch(s.ctx, page, domain.AttrSearchParams{
			StoreID: 1,
		})
		s.NoError(err)
		s.Equal(3, res.Total)
		s.Equal(3, len(res.Items))
		util.PrettyJson(res)
	})
}

func (s *ProductAttrTestSuite) TestProductAttr_ListByIDs() {
	// 创建测试数据
	testAttrs := []struct {
		name    string
		storeID int
	}{
		{"属性A", 1},
		{"属性B", 1},
		{"属性C", 1},
		{"属性D", 2},
		{"属性E", 2},
	}

	for _, u := range testAttrs {
		s.client.Attr.Create().
			SetName(u.name).
			SetStoreID(u.storeID).
			SaveX(s.ctx)
	}

	s.T().Run("根据IDs查询", func(t *testing.T) {
		res, err := s.repo.ListByIDs(s.ctx, []int{1, 2, 3})
		s.NoError(err)
		s.Equal(3, len(res))
		util.PrettyJson(res)
	})
}
