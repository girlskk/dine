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
)

// Store repository tests

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

func (s *StoreRepositoryTestSuite) createMerchant(tag string) *ent.Merchant {
	return s.client.Merchant.Create().
		SetID(uuid.New()).
		SetMerchantCode("MC-" + tag).
		SetMerchantName("商户-" + tag).
		SetMerchantShortName("简称-" + tag).
		SetMerchantType(domain.MerchantTypeBrand).
		SetBrandName("品牌-" + tag).
		SetAdminPhoneNumber("13800000000").
		SetBusinessTypeCode(domain.BusinessTypeBakery).
		SetMerchantLogo("logo-" + tag).
		SetDescription("描述-" + tag).
		SetStatus(domain.MerchantStatusActive).
		SetSuperAccount("account-" + tag).
		SetCountry(domain.CountryMY).
		SetProvince(domain.ProvinceMY01).
		SetAddress("地址-" + tag).
		SetLng("120.00").
		SetLat("30.00").
		SaveX(s.ctx)
}

func (s *StoreRepositoryTestSuite) newDomainStore(tag string, merchant *ent.Merchant) *domain.Store {
	return &domain.Store{
		ID:                      uuid.New(),
		MerchantID:              merchant.ID,
		AdminPhoneNumber:        "1390000" + tag[len(tag)-4:],
		StoreName:               "门店-" + tag,
		StoreShortName:          "简称-" + tag,
		StoreCode:               "SC-" + tag,
		Status:                  domain.StoreStatusOpen,
		BusinessModel:           domain.BusinessModelDirect,
		BusinessTypeCode:        domain.BusinessTypeBakery,
		LocationNumber:          "L-" + tag,
		ContactName:             "联系人-" + tag,
		ContactPhone:            "1370000" + tag[len(tag)-4:],
		UnifiedSocialCreditCode: "USC-" + tag,
		StoreLogo:               "logo-" + tag,
		BusinessLicenseURL:      "license-" + tag,
		StorefrontURL:           "front-" + tag,
		CashierDeskURL:          "cashier-" + tag,
		DiningEnvironmentURL:    "environment-" + tag,
		FoodOperationLicenseURL: "food-" + tag,
		LoginAccount:            "login-" + tag,
		BusinessHours:           []domain.BusinessHours{{Weekdays: []time.Weekday{time.Monday}, BusinessHours: []domain.BusinessHour{{StartTime: "09:00:00", EndTime: "18:00:00"}}}},
		DiningPeriods:           []domain.DiningPeriod{{Name: "午餐", StartTime: "11:00:00", EndTime: "13:00:00"}},
		ShiftTimes:              []domain.ShiftTime{{Name: "早班", StartTime: "09:00:00", EndTime: "15:00:00"}},
		Address: &domain.Address{
			Country:  domain.CountryMY,
			Province: domain.ProvinceMY01,
			Address:  "门店地址-" + tag,
			Lng:      "121.00",
			Lat:      "31.00",
		},
	}
}

func (s *StoreRepositoryTestSuite) TestStore_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		tag := "create-test"

		merchant := s.createMerchant(tag)
		store := s.newDomainStore(tag, merchant)

		err := s.repo.Create(s.ctx, store)
		require.NoError(t, err)

		saved := s.client.Store.GetX(s.ctx, store.ID)
		require.Equal(t, store.StoreName, saved.StoreName)
		require.Equal(t, store.Address.Address, saved.Address)
		require.Equal(t, store.BusinessTypeCode, saved.BusinessTypeCode)
		require.Equal(t, store.LoginAccount, saved.SuperAccount)
		require.False(t, saved.CreatedAt.IsZero())
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Create(s.ctx, nil)
		require.Error(t, err)
	})

	s.T().Run("地址为空", func(t *testing.T) {
		tag := "create-no-address"
		merchant := s.createMerchant(tag)
		store := s.newDomainStore(tag, merchant)
		store.Address = nil

		err := s.repo.Create(s.ctx, store)
		require.Error(t, err)
	})
}

func (s *StoreRepositoryTestSuite) TestStore_Update() {
	tag := "update-test"

	merchant := s.createMerchant(tag)
	store := s.newDomainStore(tag, merchant)
	require.NoError(s.T(), s.repo.Create(s.ctx, store))

	store.StoreName = "更新-" + store.StoreName
	store.BusinessTypeCode = domain.BusinessTypeChineseFood
	store.Address.Address = "新地址-" + tag
	store.Address.Lng = "122.00"
	store.Address.Lat = "32.00"

	err := s.repo.Update(s.ctx, store)
	require.NoError(s.T(), err)

	updated := s.client.Store.GetX(s.ctx, store.ID)
	require.Equal(s.T(), store.StoreName, updated.StoreName)
	require.Equal(s.T(), store.Address.Address, updated.Address)
	require.Equal(s.T(), domain.BusinessTypeChineseFood, updated.BusinessTypeCode)

	s.T().Run("不存在的ID", func(t *testing.T) {
		missing := s.newDomainStore("missing"+tag, merchant)
		err := s.repo.Update(s.ctx, missing)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})

	s.T().Run("地址为空", func(t *testing.T) {
		store.Address = nil
		err := s.repo.Update(s.ctx, store)
		require.Error(t, err)
	})

	s.T().Run("入参为空", func(t *testing.T) {
		err := s.repo.Update(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *StoreRepositoryTestSuite) TestStore_Delete() {
	tag := "delete-test"
	merchant := s.createMerchant(tag)
	store := s.newDomainStore(tag, merchant)
	require.NoError(s.T(), s.repo.Create(s.ctx, store))

	err := s.repo.Delete(s.ctx, store.ID)
	require.NoError(s.T(), err)

	_, err = s.client.Store.Get(s.ctx, store.ID)
	require.Error(s.T(), err)

	err = s.repo.Delete(s.ctx, uuid.New())
	require.Error(s.T(), err)
	require.True(s.T(), domain.IsNotFound(err))
}

func (s *StoreRepositoryTestSuite) TestStore_FindByID() {
	tag := "find-by-id-test"
	merchant := s.createMerchant(tag)
	store := s.newDomainStore(tag, merchant)
	require.NoError(s.T(), s.repo.Create(s.ctx, store))

	found, err := s.repo.FindByID(s.ctx, store.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), store.StoreName, found.StoreName)
	require.Equal(s.T(), store.Address.Address, found.Address.Address)

	_, err = s.repo.FindByID(s.ctx, uuid.New())
	require.Error(s.T(), err)
	require.True(s.T(), domain.IsNotFound(err))
}

func (s *StoreRepositoryTestSuite) TestStore_FindStoreMerchant() {
	tag := "find-store-merchant-test"
	merchant := s.createMerchant(tag)
	store := s.newDomainStore(tag, merchant)
	require.NoError(s.T(), s.repo.Create(s.ctx, store))

	found, err := s.repo.FindStoreMerchant(s.ctx, merchant.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), store.StoreName, found.StoreName)

	_, err = s.repo.FindStoreMerchant(s.ctx, uuid.Nil)
	require.Error(s.T(), err)

	_, err = s.repo.FindStoreMerchant(s.ctx, uuid.New())
	require.Error(s.T(), err)
	require.True(s.T(), domain.IsNotFound(err))
}

func (s *StoreRepositoryTestSuite) TestStore_GetStores() {
	tag := "get-stores-test"

	merchant := s.createMerchant(tag)
	store1 := s.newDomainStore(tag+"-1", merchant)
	require.NoError(s.T(), s.repo.Create(s.ctx, store1))
	time.Sleep(10 * time.Millisecond)
	store2 := s.newDomainStore(tag+"-2", merchant)
	require.NoError(s.T(), s.repo.Create(s.ctx, store2))

	pager := upagination.New(1, 10)
	filter := &domain.StoreListFilter{MerchantID: merchant.ID}
	list, total, err := s.repo.GetStores(s.ctx, pager, filter)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, total)
	require.Len(s.T(), list, 2)

	orderByIDAsc := domain.NewStoreListOrderByCreatedAt(false)
	list, total, err = s.repo.GetStores(s.ctx, pager, filter, orderByIDAsc)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, total)
	require.Equal(s.T(), store1.ID, list[0].ID)

	nameFilter := &domain.StoreListFilter{StoreName: store1.StoreName[6:], MerchantID: merchant.ID}
	list, total, err = s.repo.GetStores(s.ctx, pager, nameFilter)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, total)
	require.Equal(s.T(), store1.ID, list[0].ID)

	timeBound := time.Now()
	filterTime := &domain.StoreListFilter{MerchantID: merchant.ID, CreatedAtGte: &timeBound}
	_, _, err = s.repo.GetStores(s.ctx, pager, filterTime)
	require.NoError(s.T(), err)
}

func (s *StoreRepositoryTestSuite) TestStore_ExistsStore() {
	tag := "exists-store-test"

	merchant := s.createMerchant(tag)
	store := s.newDomainStore(tag, merchant)
	require.NoError(s.T(), s.repo.Create(s.ctx, store))

	exists, err := s.repo.ExistsStore(s.ctx, &domain.ExistsStoreParams{StoreName: store.StoreName})
	require.NoError(s.T(), err)
	require.True(s.T(), exists)

	exists, err = s.repo.ExistsStore(s.ctx, &domain.ExistsStoreParams{StoreName: store.StoreName, ExcludeID: store.ID})
	require.NoError(s.T(), err)
	require.False(s.T(), exists)

	exists, err = s.repo.ExistsStore(s.ctx, &domain.ExistsStoreParams{StoreName: "不存在"})
	require.NoError(s.T(), err)
	require.False(s.T(), exists)

	_, err = s.repo.ExistsStore(s.ctx, nil)
	require.Error(s.T(), err)
}

func (s *StoreRepositoryTestSuite) TestStore_CountStoresByMerchantID() {
	tag := "count-store-test"
	merchant1 := s.createMerchant(tag + "-1")
	merchant2 := s.createMerchant(tag + "-2")
	store1 := s.newDomainStore(tag+"-1", merchant1)
	store2 := s.newDomainStore(tag+"-2", merchant1)
	store3 := s.newDomainStore(tag+"-3", merchant2)
	require.NoError(s.T(), s.repo.Create(s.ctx, store1))
	require.NoError(s.T(), s.repo.Create(s.ctx, store2))
	require.NoError(s.T(), s.repo.Create(s.ctx, store3))

	counts, err := s.repo.CountStoresByMerchantID(s.ctx, []uuid.UUID{merchant1.ID, merchant2.ID})
	require.NoError(s.T(), err)
	require.Len(s.T(), counts, 2)

	m1 := counts[0]
	m2 := counts[1]
	if m1.MerchantID != merchant1.ID {
		m1, m2 = m2, m1
	}
	require.Equal(s.T(), merchant1.ID, m1.MerchantID)
	require.Equal(s.T(), 2, m1.StoreCount)
	require.Equal(s.T(), merchant2.ID, m2.MerchantID)
	require.Equal(s.T(), 1, m2.StoreCount)

	_, err = s.repo.CountStoresByMerchantID(s.ctx, []uuid.UUID{})
	require.Error(s.T(), err)
}

func (s *StoreRepositoryTestSuite) TestStore_ListByIDs() {
	tag := "list-by-ids-test"
	merchant := s.createMerchant(tag)

	store1 := s.newDomainStore(tag+"-1", merchant)
	store2 := s.newDomainStore(tag+"-2", merchant)
	require.NoError(s.T(), s.repo.Create(s.ctx, store1))
	require.NoError(s.T(), s.repo.Create(s.ctx, store2))

	list, err := s.repo.ListByIDs(s.ctx, []uuid.UUID{store1.ID, store2.ID})
	require.NoError(s.T(), err)
	require.Len(s.T(), list, 2)

	list, err = s.repo.ListByIDs(s.ctx, []uuid.UUID{store1.ID, uuid.New()})
	require.NoError(s.T(), err)
	require.Len(s.T(), list, 1)

	list, err = s.repo.ListByIDs(s.ctx, []uuid.UUID{})
	require.NoError(s.T(), err)
	require.Nil(s.T(), list)
}
