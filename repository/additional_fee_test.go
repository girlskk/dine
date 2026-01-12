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

type AdditionalFeeRepositoryTestSuite struct {
	RepositoryTestSuite
	repo *AdditionalFeeRepository
	ctx  context.Context
}

func TestAdditionalFeeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AdditionalFeeRepositoryTestSuite))
}

func (s *AdditionalFeeRepositoryTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &AdditionalFeeRepository{Client: s.client}
	s.ctx = context.Background()
}
func (s *AdditionalFeeRepositoryTestSuite) newFeeWithMerchantID(tag string, merchantID uuid.UUID) *domain.AdditionalFee {
	return &domain.AdditionalFee{
		ID:                  uuid.New(),
		Name:                "附加费-" + tag,
		FeeType:             domain.AdditionalFeeTypeStore,
		FeeCategory:         domain.AdditionalCategoryService,
		ChargeMode:          domain.AdditionalFeeChargeModeFixed,
		FeeValue:            decimal.NewFromInt(100),
		IncludeInReceivable: true,
		Taxable:             true,
		DiscountScope:       domain.AdditionalFeeDiscountScopeBefore,
		OrderChannels:       []domain.OrderChannel{domain.OrderChannelPOS, domain.OrderChannelSelfOrder},
		DiningWays:          []domain.DiningWay{domain.DiningWayDineIn, domain.DiningWayTakeOut},
		Enabled:             true,
		SortOrder:           10,
		MerchantID:          merchantID,
		StoreID:             uuid.New(),
	}
}

func (s *AdditionalFeeRepositoryTestSuite) newFee(tag string) *domain.AdditionalFee {
	return s.newFeeWithMerchantID(tag, uuid.New())
}

func (s *AdditionalFeeRepositoryTestSuite) TestAdditionalFee_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		fee := s.newFee("create")

		err := s.repo.Create(s.ctx, fee)
		require.NoError(t, err)

		saved := s.client.AdditionalFee.GetX(s.ctx, fee.ID)
		require.Equal(t, fee.Name, saved.Name)
		require.Equal(t, fee.FeeType, saved.FeeType)
		require.True(t, saved.IncludeInReceivable)
		require.Equal(t, fee.OrderChannels, saved.OrderChannels)
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Create(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *AdditionalFeeRepositoryTestSuite) TestAdditionalFee_FindByID() {
	fee := s.newFee("find")
	require.NoError(s.T(), s.repo.Create(s.ctx, fee))

	s.T().Run("查询成功", func(t *testing.T) {
		found, err := s.repo.FindByID(s.ctx, fee.ID)
		require.NoError(t, err)
		require.Equal(t, fee.ID, found.ID)
		require.Equal(t, fee.Name, found.Name)
		require.Equal(t, fee.FeeValue, found.FeeValue)
	})

	s.T().Run("不存在", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *AdditionalFeeRepositoryTestSuite) TestAdditionalFee_Update() {
	fee := s.newFee("update")
	require.NoError(s.T(), s.repo.Create(s.ctx, fee))

	s.T().Run("更新成功", func(t *testing.T) {
		fee.Name = "更新-" + fee.Name
		fee.FeeValue = decimal.NewFromInt(200)
		fee.Enabled = false
		fee.SortOrder = 1

		err := s.repo.Update(s.ctx, fee)
		require.NoError(t, err)

		updated := s.client.AdditionalFee.GetX(s.ctx, fee.ID)
		require.Equal(t, fee.Name, updated.Name)
		require.True(t, fee.FeeValue.Equal(updated.FeeValue))
		require.Equal(t, fee.SortOrder, updated.SortOrder)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		missing := s.newFee("missing")
		err := s.repo.Update(s.ctx, missing)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Update(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *AdditionalFeeRepositoryTestSuite) TestAdditionalFee_Delete() {
	fee := s.newFee("delete")
	require.NoError(s.T(), s.repo.Create(s.ctx, fee))

	s.T().Run("删除成功", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, fee.ID)
		require.NoError(t, err)
		_, err = s.client.AdditionalFee.Get(s.ctx, fee.ID)
		require.Error(t, err)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *AdditionalFeeRepositoryTestSuite) TestAdditionalFee_GetAdditionalFees() {
	m1 := uuid.New()
	m2 := uuid.New()
	fee1 := s.newFeeWithMerchantID("001", m1)
	fee1.FeeType = domain.AdditionalFeeTypeMerchant
	require.NoError(s.T(), s.repo.Create(s.ctx, fee1))
	time.Sleep(10 * time.Millisecond)
	fee2 := s.newFeeWithMerchantID("002", m1)
	fee2.Enabled = false
	fee2.SortOrder = 5
	fee2.FeeType = domain.AdditionalFeeTypeStore
	require.NoError(s.T(), s.repo.Create(s.ctx, fee2))
	time.Sleep(10 * time.Millisecond)
	fee3 := s.newFeeWithMerchantID("003", m2)
	require.NoError(s.T(), s.repo.Create(s.ctx, fee3))

	pager := upagination.New(1, 10)

	s.T().Run("按商户筛选默认排序", func(t *testing.T) {
		list, total, err := s.repo.GetAdditionalFees(s.ctx, pager, &domain.AdditionalFeeListFilter{MerchantID: m1})
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, list, 2)
		require.Equal(t, fee2.ID, list[0].ID)
	})

	s.T().Run("按排序字段升序", func(t *testing.T) {
		order := domain.NewAdditionalFeeOrderBySortOrder(false)
		list, total, err := s.repo.GetAdditionalFees(s.ctx, pager, &domain.AdditionalFeeListFilter{MerchantID: m1}, order)
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Equal(t, fee2.ID, list[0].ID)
	})

	s.T().Run("按类型与启用状态筛选", func(t *testing.T) {
		enable := true
		list, total, err := s.repo.GetAdditionalFees(s.ctx, pager, &domain.AdditionalFeeListFilter{MerchantID: m2, FeeType: domain.AdditionalFeeTypeStore, Enabled: &enable})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, fee3.ID, list[0].ID)
	})

	s.T().Run("按名称模糊筛选", func(t *testing.T) {
		list, total, err := s.repo.GetAdditionalFees(s.ctx, pager, &domain.AdditionalFeeListFilter{Name: "002"})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, fee2.ID, list[0].ID)
	})

	s.T().Run("入参缺失", func(t *testing.T) {
		_, _, err := s.repo.GetAdditionalFees(s.ctx, nil, &domain.AdditionalFeeListFilter{})
		require.Error(t, err)
		_, _, err = s.repo.GetAdditionalFees(s.ctx, pager, nil)
		require.Error(t, err)
	})
}

func (s *AdditionalFeeRepositoryTestSuite) TestAdditionalFee_Exists() {
	merchantID := uuid.New()
	fee := s.newFeeWithMerchantID("exists", merchantID)
	require.NoError(s.T(), s.repo.Create(s.ctx, fee))

	s.T().Run("同名存在", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.AdditionalFeeExistsParams{Name: fee.Name, MerchantID: merchantID})
		require.NoError(t, err)
		require.True(t, exists)
	})

	s.T().Run("排除自身", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.AdditionalFeeExistsParams{Name: fee.Name, MerchantID: merchantID, ExcludeID: fee.ID})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("不同商户不冲突", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.AdditionalFeeExistsParams{Name: fee.Name, MerchantID: uuid.New()})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("仅名称查询", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.AdditionalFeeExistsParams{Name: fee.Name})
		require.NoError(t, err)
		require.True(t, exists)
	})
}
