package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type MerchantRepositoryTestSuite struct {
	RepositoryTestSuite
	repo *MerchantRepository
	ctx  context.Context
}

func TestMerchantRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MerchantRepositoryTestSuite))
}

func (s *MerchantRepositoryTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &MerchantRepository{Client: s.client}
	s.ctx = context.Background()
}

func (s *MerchantRepositoryTestSuite) newMerchant(tag string) *domain.Merchant {
	suffix := tag
	if len(tag) > 4 {
		suffix = tag[len(tag)-4:]
	}

	return &domain.Merchant{
		ID:                uuid.New(),
		MerchantCode:      "MC-" + tag,
		MerchantName:      "商户-" + tag,
		MerchantShortName: "简称-" + tag,
		MerchantType:      domain.MerchantTypeBrand,
		BrandName:         "品牌-" + tag,
		AdminPhoneNumber:  "1380000" + suffix,
		BusinessTypeCode:  domain.BusinessTypeBakery,
		MerchantLogo:      "logo-" + tag,
		Description:       "描述-" + tag,
		Status:            domain.MerchantStatusActive,
		LoginAccount:      "login-" + tag,
		Address: &domain.Address{
			Country:  domain.CountryMY,
			Province: domain.ProvinceMY01,
			Address:  "地址-" + tag,
			Lng:      "120.0",
			Lat:      "30.0",
		},
	}
}

func (s *MerchantRepositoryTestSuite) TestMerchant_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		m := s.newMerchant("create")

		err := s.repo.Create(s.ctx, m)
		require.NoError(t, err)

		saved := s.client.Merchant.GetX(s.ctx, m.ID)
		require.Equal(t, m.MerchantName, saved.MerchantName)
		require.Equal(t, m.MerchantType, saved.MerchantType)
		require.Equal(t, m.LoginAccount, saved.SuperAccount)
		require.False(t, saved.CreatedAt.IsZero())
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Create(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *MerchantRepositoryTestSuite) TestMerchant_FindByID() {
	m := s.newMerchant("find")
	require.NoError(s.T(), s.repo.Create(s.ctx, m))

	s.T().Run("查询成功", func(t *testing.T) {
		found, err := s.repo.FindByID(s.ctx, m.ID)
		require.NoError(t, err)
		require.Equal(t, m.ID, found.ID)
		require.Equal(t, m.MerchantName, found.MerchantName)
		require.Equal(t, m.Address.Address, found.Address.Address)
	})

	s.T().Run("不存在", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *MerchantRepositoryTestSuite) TestMerchant_Update() {
	m := s.newMerchant("update")
	require.NoError(s.T(), s.repo.Create(s.ctx, m))

	s.T().Run("更新成功", func(t *testing.T) {
		m.MerchantName = "更新-" + m.MerchantName
		m.Description = ""
		m.Status = domain.MerchantStatusDisabled

		err := s.repo.Update(s.ctx, m)
		require.NoError(t, err)

		updated := s.client.Merchant.GetX(s.ctx, m.ID)
		require.Equal(t, m.MerchantName, updated.MerchantName)
		require.Equal(t, m.Status, updated.Status)
		require.Empty(t, updated.Description)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		missing := s.newMerchant("missing")
		err := s.repo.Update(s.ctx, missing)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Update(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *MerchantRepositoryTestSuite) TestMerchant_Delete() {
	m := s.newMerchant("delete")
	require.NoError(s.T(), s.repo.Create(s.ctx, m))

	s.T().Run("删除成功", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, m.ID)
		require.NoError(t, err)
		_, err = s.client.Merchant.Get(s.ctx, m.ID)
		require.Error(t, err)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *MerchantRepositoryTestSuite) TestMerchant_GetMerchants() {
	m1 := s.newMerchant("001")
	require.NoError(s.T(), s.repo.Create(s.ctx, m1))
	time.Sleep(10 * time.Millisecond)
	m2 := s.newMerchant("002")
	m2.Status = domain.MerchantStatusExpired
	require.NoError(s.T(), s.repo.Create(s.ctx, m2))

	pager := upagination.New(1, 10)

	s.T().Run("按名称筛选默认排序", func(t *testing.T) {
		list, total, err := s.repo.GetMerchants(s.ctx, pager, &domain.MerchantListFilter{MerchantName: "商户-"})
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, list, 2)
		gotIDs := map[uuid.UUID]bool{list[0].ID: true, list[1].ID: true}
		require.True(t, gotIDs[m1.ID])
		require.True(t, gotIDs[m2.ID])
	})

	s.T().Run("按创建时间升序", func(t *testing.T) {
		order := domain.NewMerchantListOrderByCreatedAt(false)
		list, total, err := s.repo.GetMerchants(s.ctx, pager, &domain.MerchantListFilter{MerchantName: "商户-"}, order)
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Equal(t, m1.ID, list[0].ID)
	})

	s.T().Run("按状态筛选", func(t *testing.T) {
		list, total, err := s.repo.GetMerchants(s.ctx, pager, &domain.MerchantListFilter{Status: domain.MerchantStatusExpired})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, m2.ID, list[0].ID)
	})

	s.T().Run("缺少入参", func(t *testing.T) {
		_, _, err := s.repo.GetMerchants(s.ctx, nil, &domain.MerchantListFilter{})
		require.Error(t, err)
		_, _, err = s.repo.GetMerchants(s.ctx, pager, nil)
		require.Error(t, err)
	})
}

func (s *MerchantRepositoryTestSuite) TestMerchant_CountMerchant() {
	m1 := s.newMerchant("count1")
	m1.MerchantType = domain.MerchantTypeBrand
	m1.ExpireUTC = func() *time.Time { t := time.Now().UTC().Add(-24 * time.Hour); return &t }()
	require.NoError(s.T(), s.repo.Create(s.ctx, m1))

	m2 := s.newMerchant("count2")
	m2.MerchantType = domain.MerchantTypeStore
	require.NoError(s.T(), s.repo.Create(s.ctx, m2))

	count, err := s.repo.CountMerchant(s.ctx)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), count)
	require.Equal(s.T(), 1, count.MerchantTypeBrand)
	require.Equal(s.T(), 1, count.MerchantTypeStore)
	require.Equal(s.T(), 1, count.Expired)
}

func (s *MerchantRepositoryTestSuite) TestMerchant_ExistMerchant() {
	m := s.newMerchant("exists")
	require.NoError(s.T(), s.repo.Create(s.ctx, m))

	s.T().Run("名称存在", func(t *testing.T) {
		exist, err := s.repo.ExistMerchant(s.ctx, &domain.MerchantExistsParams{MerchantName: m.MerchantName})
		require.NoError(t, err)
		require.True(t, exist)
	})

	s.T().Run("排除自身", func(t *testing.T) {
		exist, err := s.repo.ExistMerchant(s.ctx, &domain.MerchantExistsParams{MerchantName: m.MerchantName, ExcludeID: m.ID})
		require.NoError(t, err)
		require.False(t, exist)
	})

	s.T().Run("不存在的名称", func(t *testing.T) {
		exist, err := s.repo.ExistMerchant(s.ctx, &domain.MerchantExistsParams{MerchantName: "other"})
		require.NoError(t, err)
		require.False(t, exist)
	})

	s.T().Run("空入参", func(t *testing.T) {
		_, err := s.repo.ExistMerchant(s.ctx, nil)
		require.Error(t, err)
	})
}
