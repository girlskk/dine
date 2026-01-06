package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type TaxFeeRepositoryTestSuite struct {
	RepositoryTestSuite
	repo *TaxFeeRepository
	ctx  context.Context
}

func TestTaxFeeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TaxFeeRepositoryTestSuite))
}

func (s *TaxFeeRepositoryTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &TaxFeeRepository{Client: s.client}
	s.ctx = context.Background()
}

func (s *TaxFeeRepositoryTestSuite) newTaxFee(tag string, merchantID, storeID uuid.UUID) *domain.TaxFee {
	return &domain.TaxFee{
		ID:          uuid.New(),
		Name:        "税费-" + tag,
		TaxFeeType:  domain.TaxFeeTypeStore,
		TaxCode:     "CODE-" + tag,
		TaxRateType: domain.TaxRateTypeUnified,
		TaxRate:     decimal.NewFromFloat(0.06),
		DefaultTax:  true,
		MerchantID:  merchantID,
		StoreID:     storeID,
	}
}

func (s *TaxFeeRepositoryTestSuite) TestTaxFee_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		fee := s.newTaxFee("create", uuid.New(), uuid.Nil)

		err := s.repo.Create(s.ctx, fee)
		require.NoError(t, err)

		saved := s.client.TaxFee.GetX(s.ctx, fee.ID)
		require.Equal(t, fee.Name, saved.Name)
		require.Equal(t, fee.TaxCode, saved.TaxCode)
		require.True(t, saved.DefaultTax)
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Create(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *TaxFeeRepositoryTestSuite) TestTaxFee_FindByID() {
	fee := s.newTaxFee("find", uuid.New(), uuid.Nil)
	require.NoError(s.T(), s.repo.Create(s.ctx, fee))

	s.T().Run("查询成功", func(t *testing.T) {
		found, err := s.repo.FindByID(s.ctx, fee.ID)
		require.NoError(t, err)
		require.Equal(t, fee.ID, found.ID)
		require.Equal(t, fee.Name, found.Name)
	})

	s.T().Run("不存在", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *TaxFeeRepositoryTestSuite) TestTaxFee_Update() {
	fee := s.newTaxFee("update", uuid.New(), uuid.Nil)
	require.NoError(s.T(), s.repo.Create(s.ctx, fee))

	s.T().Run("更新成功", func(t *testing.T) {
		fee.Name = "更新-" + fee.Name
		fee.TaxRate = decimal.NewFromFloat(0.08)
		fee.DefaultTax = false

		err := s.repo.Update(s.ctx, fee)
		require.NoError(t, err)

		updated := s.client.TaxFee.GetX(s.ctx, fee.ID)
		require.Equal(t, fee.Name, updated.Name)
		require.True(t, fee.TaxRate.Equal(updated.TaxRate))
		require.Equal(t, fee.DefaultTax, updated.DefaultTax)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		missing := s.newTaxFee("missing", uuid.New(), uuid.Nil)
		err := s.repo.Update(s.ctx, missing)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Update(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *TaxFeeRepositoryTestSuite) TestTaxFee_Delete() {
	fee := s.newTaxFee("delete", uuid.New(), uuid.Nil)
	require.NoError(s.T(), s.repo.Create(s.ctx, fee))

	s.T().Run("删除成功", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, fee.ID)
		require.NoError(t, err)
		_, err = s.client.TaxFee.Get(s.ctx, fee.ID)
		require.Error(t, err)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *TaxFeeRepositoryTestSuite) TestTaxFee_GetTaxFees() {
	m1 := uuid.New()
	m2 := uuid.New()
	fee1 := s.newTaxFee("001", m1, uuid.New())
	require.NoError(s.T(), s.repo.Create(s.ctx, fee1))
	time.Sleep(10 * time.Millisecond)
	fee2 := s.newTaxFee("002", m1, uuid.New())
	fee2.DefaultTax = false
	require.NoError(s.T(), s.repo.Create(s.ctx, fee2))
	time.Sleep(10 * time.Millisecond)
	fee3 := s.newTaxFee("003", m2, uuid.New())
	fee3.TaxFeeType = domain.TaxFeeTypeMerchant
	require.NoError(s.T(), s.repo.Create(s.ctx, fee3))

	pager := upagination.New(1, 10)

	s.T().Run("按商户筛选默认排序", func(t *testing.T) {
		list, total, err := s.repo.GetTaxFees(s.ctx, pager, &domain.TaxFeeListFilter{MerchantID: m1})
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, list, 2)
		require.Equal(t, fee2.ID, list[0].ID)
	})

	s.T().Run("按创建时间升序", func(t *testing.T) {
		order := domain.NewTaxFeeOrderByCreatedAt(false)
		list, total, err := s.repo.GetTaxFees(s.ctx, pager, &domain.TaxFeeListFilter{MerchantID: m1}, order)
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Equal(t, fee1.ID, list[0].ID)
	})

	s.T().Run("按类型筛选", func(t *testing.T) {
		list, total, err := s.repo.GetTaxFees(s.ctx, pager, &domain.TaxFeeListFilter{MerchantID: m2, TaxFeeType: domain.TaxFeeTypeMerchant})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, fee3.ID, list[0].ID)
	})

	s.T().Run("按名称模糊筛选", func(t *testing.T) {
		list, total, err := s.repo.GetTaxFees(s.ctx, pager, &domain.TaxFeeListFilter{Name: "002"})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, fee2.ID, list[0].ID)
	})

	s.T().Run("入参缺失", func(t *testing.T) {
		_, _, err := s.repo.GetTaxFees(s.ctx, nil, &domain.TaxFeeListFilter{})
		require.Error(t, err)
		_, _, err = s.repo.GetTaxFees(s.ctx, pager, nil)
		require.Error(t, err)
	})
}

func (s *TaxFeeRepositoryTestSuite) TestTaxFee_Exists() {
	merchantID := uuid.New()
	storeID := uuid.New()
	fee := s.newTaxFee("exists", merchantID, storeID)
	require.NoError(s.T(), s.repo.Create(s.ctx, fee))

	s.T().Run("同名存在", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.TaxFeeExistsParams{Name: fee.Name, MerchantID: merchantID, StoreID: storeID})
		require.NoError(t, err)
		require.True(t, exists)
	})

	s.T().Run("排除自身", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.TaxFeeExistsParams{Name: fee.Name, MerchantID: merchantID, StoreID: storeID, ExcludeID: fee.ID})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("不同商户不冲突", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.TaxFeeExistsParams{Name: fee.Name, MerchantID: uuid.New(), StoreID: storeID})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("不同门店不冲突", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.TaxFeeExistsParams{Name: fee.Name, MerchantID: merchantID, StoreID: uuid.New()})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("仅名称查询", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.TaxFeeExistsParams{Name: fee.Name})
		require.NoError(t, err)
		require.True(t, exists)
	})
}
