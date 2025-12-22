package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

type locationInfo struct {
	countryID    uuid.UUID
	provinceID   uuid.UUID
	cityID       uuid.UUID
	districtID   uuid.UUID
	countryName  string
	provinceName string
	cityName     string
	districtName string
}

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

func (s *MerchantRepositoryTestSuite) createAdminUser(username string) *ent.AdminUser {
	hashedPassword, err := util.HashPassword("123456")
	require.NoError(s.T(), err)

	return s.client.AdminUser.Create().
		SetUsername(username).
		SetHashedPassword(hashedPassword).
		SetNickname("测试管理员").
		SaveX(s.ctx)
}

func (s *MerchantRepositoryTestSuite) createBusinessType(code, name string) *ent.MerchantBusinessType {
	return s.client.MerchantBusinessType.Create().
		SetTypeCode(code).
		SetTypeName(name).
		SaveX(s.ctx)
}

func (s *MerchantRepositoryTestSuite) createLocation(tag string) locationInfo {
	countryName := "国家-" + tag
	country := s.client.Country.Create().
		SetName(countryName).
		SaveX(s.ctx)

	provinceName := "省份-" + tag
	province := s.client.Province.Create().
		SetCountry(country).
		SetName(provinceName).
		SaveX(s.ctx)

	cityName := "城市-" + tag
	city := s.client.City.Create().
		SetCountry(country).
		SetProvince(province).
		SetName(cityName).
		SaveX(s.ctx)

	districtName := "区域-" + tag
	district := s.client.District.Create().
		SetCountry(country).
		SetProvince(province).
		SetCity(city).
		SetName(districtName).
		SaveX(s.ctx)

	return locationInfo{
		countryID:    country.ID,
		provinceID:   province.ID,
		cityID:       city.ID,
		districtID:   district.ID,
		countryName:  countryName,
		provinceName: provinceName,
		cityName:     cityName,
		districtName: districtName,
	}
}

func (s *MerchantRepositoryTestSuite) newMerchant(tag string, status domain.MerchantStatus, merchantType domain.MerchantType, expire *time.Time) (*domain.Merchant, locationInfo, *ent.MerchantBusinessType, *ent.AdminUser) {
	loc := s.createLocation(tag)
	admin := s.createAdminUser("admin-" + tag)
	businessType := s.createBusinessType("bt-code-"+tag, "业态-"+tag)

	merchant := &domain.Merchant{
		ID:                uuid.New(),
		MerchantCode:      "MC-" + tag,
		MerchantName:      "商户-" + tag,
		MerchantShortName: "简称-" + tag,
		MerchantType:      merchantType,
		BrandName:         "品牌-" + tag,
		AdminPhoneNumber:  "13800000000",
		ExpireUTC:         expire,
		BusinessTypeID:    businessType.ID,
		MerchantLogo:      "logo-" + tag,
		Description:       "描述-" + tag,
		Status:            status,
		AdminUserID:       admin.ID,
		Address: &domain.Address{
			CountryID:  loc.countryID,
			ProvinceID: loc.provinceID,
			CityID:     loc.cityID,
			DistrictID: loc.districtID,
			Address:    "地址-" + tag,
			Lng:        "120.00",
			Lat:        "30.00",
		},
	}

	return merchant, loc, businessType, admin
}

func (s *MerchantRepositoryTestSuite) TestMerchant_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		expire := time.Now().UTC().Add(48 * time.Hour)
		merchant, _, businessType, admin := s.newMerchant(uuid.NewString(), domain.MerchantStatusActive, domain.MerchantTypeBrand, &expire)

		err := s.repo.Create(s.ctx, merchant)
		require.NoError(t, err)

		saved := s.client.Merchant.Query().
			WithAdminUser().
			WithMerchantBusinessType().
			OnlyX(s.ctx)

		require.Equal(t, merchant.MerchantName, saved.MerchantName)
		require.Equal(t, merchant.MerchantShortName, saved.MerchantShortName)
		require.Equal(t, merchant.Status, saved.Status)
		require.Equal(t, businessType.ID, saved.BusinessTypeID)
		require.Equal(t, admin.ID, saved.AdminUserID)
	})

	s.T().Run("参数为空", func(t *testing.T) {
		err := s.repo.Create(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *MerchantRepositoryTestSuite) TestMerchant_Update() {
	tag := uuid.NewString()
	merchant, _, _, _ := s.newMerchant(tag, domain.MerchantStatusActive, domain.MerchantTypeBrand, nil)
	require.NoError(s.T(), s.repo.Create(s.ctx, merchant))

	newBusinessType := s.createBusinessType("bt-new-"+tag, "新业态-"+tag)

	merchant.MerchantName = "更新-" + merchant.MerchantName
	merchant.BrandName = "新品牌-" + tag
	merchant.Status = domain.MerchantStatusDisabled
	merchant.BusinessTypeID = newBusinessType.ID
	merchant.Address.Address = "更新地址-" + tag
	merchant.Address.Lng = "121.00"
	merchant.Address.Lat = "31.00"

	err := s.repo.Update(s.ctx, merchant)
	require.NoError(s.T(), err)

	updated := s.client.Merchant.GetX(s.ctx, merchant.ID)
	require.Equal(s.T(), merchant.MerchantName, updated.MerchantName)
	require.Equal(s.T(), merchant.BrandName, updated.BrandName)
	require.Equal(s.T(), merchant.Status, updated.Status)
	require.Equal(s.T(), merchant.Address.Address, updated.Address)
	require.Equal(s.T(), merchant.Address.Lng, updated.Lng)
	require.Equal(s.T(), merchant.Address.Lat, updated.Lat)
	require.Equal(s.T(), newBusinessType.ID, updated.BusinessTypeID)

	s.T().Run("不存在的ID", func(t *testing.T) {
		missingTag := uuid.NewString()
		missingMerchant, _, _, _ := s.newMerchant(missingTag, domain.MerchantStatusActive, domain.MerchantTypeStore, nil)
		err := s.repo.Update(s.ctx, missingMerchant)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *MerchantRepositoryTestSuite) TestMerchant_Delete() {
	tag := uuid.NewString()
	merchant, _, _, _ := s.newMerchant(tag, domain.MerchantStatusActive, domain.MerchantTypeBrand, nil)
	require.NoError(s.T(), s.repo.Create(s.ctx, merchant))

	s.T().Run("正常删除", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, merchant.ID)
		require.NoError(t, err)

		_, err = s.client.Merchant.Get(s.ctx, merchant.ID)
		require.Error(t, err)
		require.True(t, ent.IsNotFound(err))
	})

	s.T().Run("删除不存在的记录", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *MerchantRepositoryTestSuite) TestMerchant_FindByID() {
	tag := uuid.NewString()
	merchant, loc, businessType, admin := s.newMerchant(tag, domain.MerchantStatusActive, domain.MerchantTypeBrand, nil)
	require.NoError(s.T(), s.repo.Create(s.ctx, merchant))

	s.T().Run("查询成功", func(t *testing.T) {
		found, err := s.repo.FindByID(s.ctx, merchant.ID)
		require.NoError(t, err)

		require.Equal(t, merchant.ID, found.ID)
		require.Equal(t, merchant.MerchantName, found.MerchantName)
		require.Equal(t, businessType.TypeName, found.BusinessTypeName)
		require.Equal(t, admin.Username, found.LoginAccount)
		require.Equal(t, admin.HashedPassword, found.LoginPassword)
		require.Equal(t, loc.countryName, found.Address.CountryName)
		require.Equal(t, loc.provinceName, found.Address.ProvinceName)
		require.Equal(t, loc.cityName, found.Address.CityName)
		require.Equal(t, loc.districtName, found.Address.DistrictName)
	})

	s.T().Run("记录不存在", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *MerchantRepositoryTestSuite) TestMerchant_GetMerchants() {
	activeTag := uuid.NewString()
	activeMerchant, _, _, _ := s.newMerchant(activeTag, domain.MerchantStatusActive, domain.MerchantTypeBrand, nil)
	require.NoError(s.T(), s.repo.Create(s.ctx, activeMerchant))

	disabledTag := uuid.NewString()
	disabledMerchant, _, _, _ := s.newMerchant(disabledTag, domain.MerchantStatusDisabled, domain.MerchantTypeStore, nil)
	require.NoError(s.T(), s.repo.Create(s.ctx, disabledMerchant))

	s.T().Run("正常分页查询", func(t *testing.T) {
		pager := upagination.New(1, 10)
		filter := &domain.MerchantListFilter{Status: domain.MerchantStatusActive}

		merchants, total, err := s.repo.GetMerchants(s.ctx, pager, filter, domain.NewMerchantListOrderByCreatedAt(false))
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Len(t, merchants, 1)
		require.Equal(t, activeMerchant.MerchantName, merchants[0].MerchantName)
	})

	s.T().Run("缺少参数", func(t *testing.T) {
		_, _, err := s.repo.GetMerchants(s.ctx, nil, &domain.MerchantListFilter{})
		require.Error(t, err)

		_, _, err = s.repo.GetMerchants(s.ctx, upagination.New(1, 10), nil)
		require.Error(t, err)
	})
}

func (s *MerchantRepositoryTestSuite) TestMerchant_CountMerchant() {
	future := time.Now().UTC().Add(24 * time.Hour)
	past := time.Now().UTC().Add(-24 * time.Hour)

	brandMerchant, _, _, _ := s.newMerchant(uuid.NewString(), domain.MerchantStatusActive, domain.MerchantTypeBrand, &future)
	require.NoError(s.T(), s.repo.Create(s.ctx, brandMerchant))

	storeMerchant, _, _, _ := s.newMerchant(uuid.NewString(), domain.MerchantStatusActive, domain.MerchantTypeStore, &past)
	require.NoError(s.T(), s.repo.Create(s.ctx, storeMerchant))

	count, err := s.repo.CountMerchant(s.ctx)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, count.MerchantTypeBrand)
	require.Equal(s.T(), 1, count.MerchantTypeStore)
	require.Equal(s.T(), 1, count.Expired)
}

func (s *MerchantRepositoryTestSuite) TestMerchant_ExistMerchant() {
	tag := uuid.NewString()
	merchant, _, _, _ := s.newMerchant(tag, domain.MerchantStatusActive, domain.MerchantTypeBrand, nil)
	require.NoError(s.T(), s.repo.Create(s.ctx, merchant))

	s.T().Run("已存在", func(t *testing.T) {
		exist, err := s.repo.ExistMerchant(s.ctx, &domain.MerchantExistsParams{MerchantName: merchant.MerchantName})
		require.NoError(t, err)
		require.True(t, exist)
	})

	s.T().Run("排除当前ID", func(t *testing.T) {
		exist, err := s.repo.ExistMerchant(s.ctx, &domain.MerchantExistsParams{MerchantName: merchant.MerchantName, ExcludeID: merchant.ID})
		require.NoError(t, err)
		require.False(t, exist)
	})

	s.T().Run("参数为空", func(t *testing.T) {
		_, err := s.repo.ExistMerchant(s.ctx, nil)
		require.Error(t, err)
	})
}
