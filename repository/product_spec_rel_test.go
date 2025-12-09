package repository

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
	"testing"
)

type ProductSpecRelTestSuite struct {
	RepositoryTestSuite
	repo *ProductSpecRelRepository
	ctx  context.Context
}

func TestProductSpecRelTestSuite(t *testing.T) {
	suite.Run(t, new(ProductSpecRelTestSuite))
}

func (s *ProductSpecRelTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &ProductSpecRelRepository{
		Client: s.client,
	}
	s.ctx = context.Background()
}

func (s *ProductSpecRelTestSuite) TestProductAttr_ListByIDs() {
	s.T().Run("通过ID加载", func(t *testing.T) {

		res, err := s.repo.ListByIDs(s.ctx, []int{1, 2})
		require.NoError(t, err)

		util.PrettyJson(res)
	})
}
