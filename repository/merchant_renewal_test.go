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

type renewalLocation struct {
	countryID    uuid.UUID
	provinceID   uuid.UUID
	cityID       uuid.UUID
	districtID   uuid.UUID
	countryName  string
	provinceName string
	cityName     string
	districtName string
}

func (s *MerchantRenewalRepositoryTestSuite) createAdminUser(tag string) *ent.AdminUser {
	return s.client.AdminUser.Create().
		SetUsername("admin-" + tag).
		SetHashedPassword("hashed").
		SetNickname("续费管理员").
		SaveX(s.ctx)
}

func (s *MerchantRenewalRepositoryTestSuite) createBusinessType(tag string) *ent.MerchantBusinessType {
	short := tag
	if len(short) > 8 {
		short = short[:8]
	}
	return s.client.MerchantBusinessType.Create().
		SetTypeCode("bt-" + short).
		SetTypeName("业态-" + short).
		SaveX(s.ctx)
}

func (s *MerchantRenewalRepositoryTestSuite) createLocation(tag string) renewalLocation {
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

	return renewalLocation{
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

func (s *MerchantRenewalRepositoryTestSuite) createMerchant(tag string) (*ent.Merchant, *ent.MerchantBusinessType, *ent.AdminUser) {
	loc := s.createLocation(tag)
	admin := s.createAdminUser(tag)
	businessType := s.createBusinessType(tag)

	short := tag
	if len(short) > 8 {
		short = short[:8]
	}

	merchant := s.client.Merchant.Create().
		SetMerchantCode("MC-" + short).
		SetMerchantName("商户-" + short).
		SetMerchantShortName("简称-" + short).
		SetMerchantType(domain.MerchantTypeBrand).
		SetBrandName("品牌-" + short).
		SetAdminPhoneNumber("13800000000").
		SetMerchantLogo("logo-" + short).
		SetDescription("描述-" + short).
		SetStatus(domain.MerchantStatusActive).
		SetBusinessTypeID(businessType.ID).
		SetMerchantBusinessType(businessType).
		SetAdminUser(admin).
		SetCountryID(loc.countryID).
		SetProvinceID(loc.provinceID).
		SetCityID(loc.cityID).
		SetDistrictID(loc.districtID).
		SetAddress("地址-" + short).
		SetLng("120.00").
		SetLat("30.00").
		SaveX(s.ctx)

	return merchant, businessType, admin
}

func (s *MerchantRenewalRepositoryTestSuite) newRenewal(tag string, merchantID uuid.UUID) *domain.MerchantRenewal {
	return &domain.MerchantRenewal{
		ID:                   uuid.New(),
		MerchantID:           merchantID,
		PurchaseDuration:     12,
		PurchaseDurationUnit: domain.PurchaseDurationUnitMonth,
		OperatorName:         "操作员-" + tag,
		OperatorAccount:      "account-" + tag,
		CreatedAt:            time.Now(),
	}
}

func (s *MerchantRenewalRepositoryTestSuite) TestMerchantRenewal_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		merchant, _, _ := s.createMerchant(uuid.NewString())
		renewal := s.newRenewal("a", merchant.ID)

		err := s.repo.Create(s.ctx, renewal)
		require.NoError(t, err)

		saved := s.client.MerchantRenewal.GetX(s.ctx, renewal.ID)
		require.Equal(t, renewal.MerchantID, saved.MerchantID)
		require.Equal(t, renewal.PurchaseDuration, saved.PurchaseDuration)
		require.Equal(t, renewal.PurchaseDurationUnit, saved.PurchaseDurationUnit)
		require.Equal(t, renewal.OperatorName, saved.OperatorName)
		require.Equal(t, renewal.OperatorAccount, saved.OperatorAccount)
	})

	s.T().Run("参数为空", func(t *testing.T) {
		err := s.repo.Create(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *MerchantRenewalRepositoryTestSuite) TestMerchantRenewal_GetByMerchant() {
	merchant, _, _ := s.createMerchant(uuid.NewString())
	renewal1 := s.newRenewal("1", merchant.ID)
	renewal2 := s.newRenewal("2", merchant.ID)
	otherMerchant, _, _ := s.createMerchant(uuid.NewString())
	renewalOther := s.newRenewal("other", otherMerchant.ID)

	require.NoError(s.T(), s.repo.Create(s.ctx, renewal1))
	time.Sleep(10 * time.Millisecond)
	require.NoError(s.T(), s.repo.Create(s.ctx, renewal2))
	require.NoError(s.T(), s.repo.Create(s.ctx, renewalOther))

	renewals, err := s.repo.GetByMerchant(s.ctx, merchant.ID)
	require.NoError(s.T(), err)
	require.Len(s.T(), renewals, 2)
	require.Equal(s.T(), renewal2.ID, renewals[0].ID)
	require.Equal(s.T(), renewal1.ID, renewals[1].ID)
}
