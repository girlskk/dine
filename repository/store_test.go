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
	storeent "gitlab.jiguang.dev/pos-dine/dine/ent/store"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type StoreRepositoryTestSuite struct {
	RepositoryTestSuite
	repo *StoreRepository
	ctx  context.Context
}

func TestStoreRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(StoreRepositoryTestSuite))
}

func (s *StoreRepositoryTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &StoreRepository{Client: s.client}
	s.ctx = context.Background()
}

type storeLocation struct {
	countryID    uuid.UUID
	provinceID   uuid.UUID
	cityID       uuid.UUID
	districtID   uuid.UUID
	countryName  string
	provinceName string
	cityName     string
	districtName string
}

func (s *StoreRepositoryTestSuite) createAdminUser(username string) *ent.AdminUser {
	return s.client.AdminUser.Create().
		SetUsername(username).
		SetHashedPassword("hashed").
		SetNickname("store-admin").
		SaveX(s.ctx)
}

func (s *StoreRepositoryTestSuite) createBusinessType(code, name string) *ent.MerchantBusinessType {
	return s.client.MerchantBusinessType.Create().
		SetTypeCode(code).
		SetTypeName(name).
		SaveX(s.ctx)
}

func (s *StoreRepositoryTestSuite) createLocation(tag string) storeLocation {
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

	return storeLocation{
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

func (s *StoreRepositoryTestSuite) createMerchant(tag string, businessType *ent.MerchantBusinessType, admin *ent.AdminUser, loc storeLocation) *ent.Merchant {
	return s.client.Merchant.Create().
		SetMerchantCode("MC-" + tag).
		SetMerchantName("商户-" + tag).
		SetMerchantShortName("简称-" + tag).
		SetMerchantType(domain.MerchantTypeBrand).
		SetBrandName("品牌-" + tag).
		SetAdminPhoneNumber("13800000000").
		SetMerchantLogo("logo-" + tag).
		SetDescription("描述-" + tag).
		SetStatus(domain.MerchantStatusActive).
		SetBusinessTypeID(businessType.ID).
		SetMerchantBusinessType(businessType).
		SetAdminUser(admin).
		SetCountryID(loc.countryID).
		SetProvinceID(loc.provinceID).
		SetCityID(loc.cityID).
		SetDistrictID(loc.districtID).
		SetAddress("地址-" + tag).
		SetLng("120.00").
		SetLat("30.00").
		SaveX(s.ctx)
}

func (s *StoreRepositoryTestSuite) newStore(tag string) (*domain.Store, storeLocation, *ent.MerchantBusinessType, *ent.AdminUser, *ent.Merchant) {
	loc := s.createLocation(tag)
	admin := s.createAdminUser("admin-" + tag)
	businessType := s.createBusinessType("bt-"+tag, "业态-"+tag)
	merchant := s.createMerchant(tag, businessType, admin, loc)

	shortTag := tag
	if len(shortTag) > 8 {
		shortTag = shortTag[:8]
	}

	businessHours := []domain.BusinessHours{
		{Weekdays: []time.Weekday{time.Monday, time.Tuesday}, StartTime: "09:00:00", EndTime: "18:00:00"},
	}
	diningPeriods := []domain.DiningPeriod{
		{Name: "午餐", StartTime: "11:00:00", EndTime: "14:00:00"},
	}
	shiftTimes := []domain.ShiftTime{
		{Name: "白班", StartTime: "09:00:00", EndTime: "17:00:00"},
	}

	store := &domain.Store{
		ID:                      uuid.New(),
		MerchantID:              merchant.ID,
		AdminPhoneNumber:        "13900000000",
		StoreName:               "门店-" + shortTag,
		StoreShortName:          "简称-" + shortTag,
		StoreCode:               "SC-" + tag,
		Status:                  domain.StoreStatusOpen,
		BusinessModel:           domain.BusinessModelDirect,
		BusinessTypeID:          businessType.ID,
		LocationNumber:          "位置编号-" + tag,
		ContactName:             "联系人-" + shortTag,
		ContactPhone:            "13700000000",
		UnifiedSocialCreditCode: "USCC-" + tag,
		StoreLogo:               "store-logo-" + tag,
		BusinessLicenseURL:      "license-" + tag,
		StorefrontURL:           "storefront-" + tag,
		CashierDeskURL:          "cashier-" + tag,
		DiningEnvironmentURL:    "env-" + tag,
		FoodOperationLicenseURL: "food-" + tag,
		AdminUserID:             admin.ID,
		BusinessHours:           businessHours,
		DiningPeriods:           diningPeriods,
		ShiftTimes:              shiftTimes,
		Address: &domain.Address{
			CountryID:  loc.countryID,
			ProvinceID: loc.provinceID,
			CityID:     loc.cityID,
			DistrictID: loc.districtID,
			Address:    "地址-" + tag,
			Lng:        "121.00",
			Lat:        "31.00",
		},
	}

	return store, loc, businessType, admin, merchant
}

func (s *StoreRepositoryTestSuite) TestStore_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		store, loc, businessType, admin, merchant := s.newStore(uuid.NewString())

		err := s.repo.Create(s.ctx, store)
		require.NoError(t, err)

		saved := s.client.Store.Query().
			Where(storeent.IDEQ(store.ID)).
			WithAdminUser().
			WithMerchantBusinessType().
			WithMerchant().
			OnlyX(s.ctx)

		require.Equal(t, merchant.ID, saved.MerchantID)
		require.Equal(t, businessType.ID, saved.BusinessTypeID)
		require.Equal(t, admin.ID, saved.AdminUserID)
		require.Equal(t, store.StoreName, saved.StoreName)
		require.Equal(t, store.StoreShortName, saved.StoreShortName)
		require.Equal(t, store.Status, saved.Status)

		found, err := s.repo.FindByID(s.ctx, store.ID)
		require.NoError(t, err)
		require.Len(t, found.BusinessHours, 1)
		require.Len(t, found.DiningPeriods, 1)
		require.Len(t, found.ShiftTimes, 1)
		require.Equal(t, loc.countryName, found.Address.CountryName)
		require.Equal(t, loc.provinceName, found.Address.ProvinceName)
		require.Equal(t, loc.cityName, found.Address.CityName)
		require.Equal(t, loc.districtName, found.Address.DistrictName)
	})

	s.T().Run("参数为空", func(t *testing.T) {
		err := s.repo.Create(s.ctx, nil)
		require.Error(t, err)
	})

	s.T().Run("地址为空", func(t *testing.T) {
		store, _, _, _, _ := s.newStore(uuid.NewString())
		store.Address = nil
		err := s.repo.Create(s.ctx, store)
		require.Error(t, err)
	})
}

func (s *StoreRepositoryTestSuite) TestStore_Update() {
	tag := uuid.NewString()
	store, _, _, _, _ := s.newStore(tag)
	require.NoError(s.T(), s.repo.Create(s.ctx, store))

	store.StoreName = "更新-" + store.StoreName
	store.Status = domain.StoreStatusClosed
	store.Address.Address = "更新地址-" + tag
	store.Address.Lng = "122.00"
	store.Address.Lat = "32.00"
	store.BusinessHours = []domain.BusinessHours{{Weekdays: []time.Weekday{time.Wednesday}, StartTime: "10:00:00", EndTime: "19:00:00"}}
	store.DiningPeriods = []domain.DiningPeriod{{Name: "晚餐", StartTime: "17:00:00", EndTime: "21:00:00"}}
	store.ShiftTimes = []domain.ShiftTime{{Name: "晚班", StartTime: "16:00:00", EndTime: "23:00:00"}}

	err := s.repo.Update(s.ctx, store)
	require.NoError(s.T(), err)

	updated := s.client.Store.GetX(s.ctx, store.ID)
	require.Equal(s.T(), store.StoreName, updated.StoreName)
	require.Equal(s.T(), store.Status, updated.Status)
	require.Equal(s.T(), store.Address.Address, updated.Address)
	require.Equal(s.T(), store.Address.Lng, updated.Lng)
	require.Equal(s.T(), store.Address.Lat, updated.Lat)

	found, err := s.repo.FindByID(s.ctx, store.ID)
	require.NoError(s.T(), err)
	require.Len(s.T(), found.BusinessHours, 1)
	require.Equal(s.T(), "10:00:00", found.BusinessHours[0].StartTime)
	require.Len(s.T(), found.DiningPeriods, 1)
	require.Equal(s.T(), "晚餐", found.DiningPeriods[0].Name)
	require.Len(s.T(), found.ShiftTimes, 1)
	require.Equal(s.T(), "晚班", found.ShiftTimes[0].Name)

	s.T().Run("不存在的ID", func(t *testing.T) {
		missing, _, _, _, _ := s.newStore(uuid.NewString())
		err := s.repo.Update(s.ctx, missing)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})

	s.T().Run("地址为空", func(t *testing.T) {
		storeNilAddr, _, _, _, _ := s.newStore(uuid.NewString())
		require.NoError(t, s.repo.Create(s.ctx, storeNilAddr))
		storeNilAddr.Address = nil
		err := s.repo.Update(s.ctx, storeNilAddr)
		require.Error(t, err)
	})
}

func (s *StoreRepositoryTestSuite) TestStore_Delete() {
	tag := uuid.NewString()
	store, _, _, _, _ := s.newStore(tag)
	require.NoError(s.T(), s.repo.Create(s.ctx, store))

	s.T().Run("正常删除", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, store.ID)
		require.NoError(t, err)

		_, err = s.client.Store.Get(s.ctx, store.ID)
		require.Error(t, err)
		require.True(t, ent.IsNotFound(err))
	})

	s.T().Run("删除不存在的记录", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *StoreRepositoryTestSuite) TestStore_FindByID() {
	tag := uuid.NewString()
	store, loc, businessType, admin, _ := s.newStore(tag)
	require.NoError(s.T(), s.repo.Create(s.ctx, store))

	s.T().Run("查询成功", func(t *testing.T) {
		found, err := s.repo.FindByID(s.ctx, store.ID)
		require.NoError(t, err)

		require.Equal(t, store.ID, found.ID)
		require.Equal(t, store.StoreName, found.StoreName)
		require.Equal(t, businessType.TypeName, found.BusinessTypeName)
		require.Equal(t, admin.Username, found.LoginAccount)
		require.Equal(t, admin.HashedPassword, found.LoginPassword)
		require.Equal(t, loc.countryName, found.Address.CountryName)
		require.Equal(t, loc.provinceName, found.Address.ProvinceName)
		require.Equal(t, loc.cityName, found.Address.CityName)
		require.Equal(t, loc.districtName, found.Address.DistrictName)
		require.Len(t, found.BusinessHours, 1)
		require.Len(t, found.DiningPeriods, 1)
		require.Len(t, found.ShiftTimes, 1)
	})

	s.T().Run("记录不存在", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *StoreRepositoryTestSuite) TestStore_GetStores() {
	tag1 := uuid.NewString()
	store1, _, _, _, _ := s.newStore(tag1)
	store1.Status = domain.StoreStatusOpen
	require.NoError(s.T(), s.repo.Create(s.ctx, store1))

	tag2 := uuid.NewString()
	store2, _, _, _, _ := s.newStore(tag2)
	store2.Status = domain.StoreStatusClosed
	require.NoError(s.T(), s.repo.Create(s.ctx, store2))

	// 新增按商户过滤的店铺，确保过滤条件被覆盖
	extraTag := uuid.NewString()
	storeByMerchant, _, _, _, merchantFilter := s.newStore(extraTag)
	storeByMerchant.Status = domain.StoreStatusOpen
	require.NoError(s.T(), s.repo.Create(s.ctx, storeByMerchant))

	s.T().Run("正常分页查询", func(t *testing.T) {
		pager := upagination.New(1, 10)
		filter := &domain.StoreListFilter{Status: domain.StoreStatusOpen}

		stores, total, err := s.repo.GetStores(s.ctx, pager, filter, domain.NewStoreListOrderByCreatedAt(false))
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, stores, 2)

		names := map[string]bool{stores[0].StoreName: true, stores[1].StoreName: true}
		require.True(t, names[store1.StoreName])
		require.True(t, names[storeByMerchant.StoreName])
	})

	s.T().Run("按商户过滤", func(t *testing.T) {
		pager := upagination.New(1, 10)
		filter := &domain.StoreListFilter{MerchantID: merchantFilter.ID}

		stores, total, err := s.repo.GetStores(s.ctx, pager, filter, domain.NewStoreListOrderByCreatedAt(false))
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Len(t, stores, 1)
		require.Equal(t, storeByMerchant.StoreName, stores[0].StoreName)
	})

	s.T().Run("缺少参数", func(t *testing.T) {
		_, _, err := s.repo.GetStores(s.ctx, nil, &domain.StoreListFilter{})
		require.Error(t, err)

		_, _, err = s.repo.GetStores(s.ctx, upagination.New(1, 10), nil)
		require.Error(t, err)
	})
}

func (s *StoreRepositoryTestSuite) TestStore_ExistsStore() {
	tag := uuid.NewString()
	store, _, _, _, merchant := s.newStore(tag)
	require.NoError(s.T(), s.repo.Create(s.ctx, store))

	s.T().Run("已存在", func(t *testing.T) {
		exists, err := s.repo.ExistsStore(s.ctx, &domain.ExistsStoreParams{StoreName: store.StoreName})
		require.NoError(t, err)
		require.True(t, exists)
	})

	s.T().Run("排除当前ID", func(t *testing.T) {
		exists, err := s.repo.ExistsStore(s.ctx, &domain.ExistsStoreParams{StoreName: store.StoreName, ExcludeID: store.ID})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("该商户已存在", func(t *testing.T) {
		exists, err := s.repo.ExistsStore(s.ctx, &domain.ExistsStoreParams{MerchantID: merchant.ID, StoreName: store.StoreName})
		require.NoError(t, err)
		require.True(t, exists)
	})

	s.T().Run("按商户过滤不存在", func(t *testing.T) {
		anotherTag := uuid.NewString()
		otherStore, _, _, _, otherMerchant := s.newStore(anotherTag)
		require.NoError(t, s.repo.Create(s.ctx, otherStore))

		exists, err := s.repo.ExistsStore(s.ctx, &domain.ExistsStoreParams{StoreName: store.StoreName, MerchantID: otherMerchant.ID})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("参数为空", func(t *testing.T) {
		_, err := s.repo.ExistsStore(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *StoreRepositoryTestSuite) TestStore_FindStoreMerchant() {
	tag := uuid.NewString()
	store, _, _, _, merchant := s.newStore(tag)
	require.NoError(s.T(), s.repo.Create(s.ctx, store))

	s.T().Run("查询成功", func(t *testing.T) {
		found, err := s.repo.FindStoreMerchant(s.ctx, merchant.ID)
		require.NoError(t, err)
		require.Equal(t, store.ID, found.ID)
	})

	s.T().Run("merchantID为空", func(t *testing.T) {
		_, err := s.repo.FindStoreMerchant(s.ctx, uuid.Nil)
		require.Error(t, err)
	})

	s.T().Run("不存在的merchant", func(t *testing.T) {
		_, err := s.repo.FindStoreMerchant(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *StoreRepositoryTestSuite) TestStore_CountStoresByMerchantID() {
	tag1 := uuid.NewString()
	store1, _, businessType1, admin1, merchant1 := s.newStore(tag1)
	require.NoError(s.T(), s.repo.Create(s.ctx, store1))

	tag2 := uuid.NewString()
	store2, _, _, _, merchant2 := s.newStore(tag2)
	require.NoError(s.T(), s.repo.Create(s.ctx, store2))

	tag3 := uuid.NewString()
	store3, _, _, _, _ := s.newStore(tag3)
	store3.MerchantID = merchant1.ID
	store3.BusinessTypeID = businessType1.ID
	store3.AdminUserID = admin1.ID
	require.NoError(s.T(), s.repo.Create(s.ctx, store3))

	s.T().Run("统计成功", func(t *testing.T) {
		counts, err := s.repo.CountStoresByMerchantID(s.ctx, []uuid.UUID{merchant1.ID, merchant2.ID})
		require.NoError(t, err)
		require.Len(t, counts, 2)

		m := map[uuid.UUID]int{}
		for _, c := range counts {
			m[c.MerchantID] = c.StoreCount
		}
		require.Equal(t, 2, m[merchant1.ID])
		require.Equal(t, 1, m[merchant2.ID])
	})

	s.T().Run("空参数", func(t *testing.T) {
		_, err := s.repo.CountStoresByMerchantID(s.ctx, []uuid.UUID{})
		require.Error(t, err)
	})

	s.T().Run("包含不存在的merchantID", func(t *testing.T) {
		missing := uuid.New()
		counts, err := s.repo.CountStoresByMerchantID(s.ctx, []uuid.UUID{merchant1.ID, missing})
		require.NoError(t, err)

		m := map[uuid.UUID]int{}
		for _, c := range counts {
			m[c.MerchantID] = c.StoreCount
		}
		require.Equal(t, 2, m[merchant1.ID])
		require.NotContains(t, m, missing)
	})
}
