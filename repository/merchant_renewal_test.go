package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

type MerchantRenewalRepositoryTestSuite struct {
	RepositoryTestSuite
	repo *MerchantRenewalRepository
	ctx  context.Context
}

func TestMerchantRenewalRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MerchantRenewalRepositoryTestSuite))
}

func (s *MerchantRenewalRepositoryTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &MerchantRenewalRepository{Client: s.client}
	s.ctx = context.Background()
}

func (s *MerchantRenewalRepositoryTestSuite) newRenewal(merchantID uuid.UUID, tag string) *domain.MerchantRenewal {
	return &domain.MerchantRenewal{
		ID:                   uuid.New(),
		MerchantID:           merchantID,
		PurchaseDuration:     12,
		PurchaseDurationUnit: domain.PurchaseDurationUnitMonth,
		OperatorName:         "操作员-" + tag,
		OperatorAccount:      "account-" + tag,
	}
}

func (s *MerchantRenewalRepositoryTestSuite) TestMerchantRenewal_Create() {
	merchantID := uuid.New()

	s.T().Run("创建成功", func(t *testing.T) {
		renewal := s.newRenewal(merchantID, "create")

		err := s.repo.Create(s.ctx, renewal)
		require.NoError(t, err)

		saved := s.client.MerchantRenewal.GetX(s.ctx, renewal.ID)
		require.Equal(t, renewal.MerchantID, saved.MerchantID)
		require.Equal(t, renewal.PurchaseDuration, saved.PurchaseDuration)
		require.Equal(t, renewal.OperatorAccount, saved.OperatorAccount)
		require.False(t, saved.CreatedAt.IsZero())
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Create(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *MerchantRenewalRepositoryTestSuite) TestMerchantRenewal_GetByMerchant() {
	merchantID := uuid.New()
	renewal1 := s.newRenewal(merchantID, "001")
	require.NoError(s.T(), s.repo.Create(s.ctx, renewal1))
	time.Sleep(10 * time.Millisecond)
	renewal2 := s.newRenewal(merchantID, "002")
	renewal2.PurchaseDuration = 24
	require.NoError(s.T(), s.repo.Create(s.ctx, renewal2))

	s.T().Run("按创建时间倒序返回", func(t *testing.T) {
		list, err := s.repo.GetByMerchant(s.ctx, merchantID)
		require.NoError(t, err)
		require.Len(t, list, 2)
		require.Equal(t, renewal2.ID, list[0].ID)
		require.Equal(t, renewal2.PurchaseDuration, list[0].PurchaseDuration)
	})

	s.T().Run("无记录返回空列表", func(t *testing.T) {
		list, err := s.repo.GetByMerchant(s.ctx, uuid.New())
		require.NoError(t, err)
		require.Len(t, list, 0)
	})
}
