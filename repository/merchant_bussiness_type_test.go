package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

type MerchantBusinessTypeRepositoryTestSuite struct {
	RepositoryTestSuite
	repo *MerchantBusinessTypeRepository
	ctx  context.Context
}

func TestMerchantBusinessTypeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MerchantBusinessTypeRepositoryTestSuite))
}

func (s *MerchantBusinessTypeRepositoryTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &MerchantBusinessTypeRepository{Client: s.client}
	s.ctx = context.Background()
}

func (s *MerchantBusinessTypeRepositoryTestSuite) createType(code, name string) uuid.UUID {
	mbt := s.client.MerchantBusinessType.Create().
		SetTypeCode(code).
		SetTypeName(name).
		SaveX(s.ctx)
	return mbt.ID
}

func (s *MerchantBusinessTypeRepositoryTestSuite) TestFindById() {
	typeID := s.createType("code-1", "业态1")

	s.T().Run("查询成功", func(t *testing.T) {
		bt, err := s.repo.FindById(s.ctx, typeID)
		require.NoError(t, err)
		require.Equal(t, typeID, bt.ID)
		require.Equal(t, "code-1", bt.TypeCode)
		require.Equal(t, "业态1", bt.TypeName)
	})

	s.T().Run("未找到", func(t *testing.T) {
		_, err := s.repo.FindById(s.ctx, uuid.New())
		require.Error(t, err)
	})
}

func (s *MerchantBusinessTypeRepositoryTestSuite) TestGetAll() {
	id1 := s.createType("code-a", "业态A")
	id2 := s.createType("code-b", "业态B")

	types, err := s.repo.GetAll(s.ctx)
	require.NoError(s.T(), err)
	require.GreaterOrEqual(s.T(), len(types), 2)

	m := make(map[uuid.UUID]*domain.MerchantBusinessType)
	for _, t := range types {
		m[t.ID] = t
	}

	require.Contains(s.T(), m, id1)
	require.Contains(s.T(), m, id2)
	require.Equal(s.T(), "code-a", m[id1].TypeCode)
	require.Equal(s.T(), "业态A", m[id1].TypeName)
	require.Equal(s.T(), "code-b", m[id2].TypeCode)
	require.Equal(s.T(), "业态B", m[id2].TypeName)
}
