// filepath: /Users/rrr/go/src/dine-api/repository/remark_test.go
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

type RemarkTestSuite struct {
	RepositoryTestSuite
	repo *RemarkRepository
	ctx  context.Context
}

func TestRemarkTestSuite(t *testing.T) {
	suite.Run(t, new(RemarkTestSuite))
}

func (s *RemarkTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &RemarkRepository{Client: s.client}
	s.ctx = context.Background()
}

func (s *RemarkTestSuite) createRemarkCategory(name string, merchantID uuid.UUID) *ent.RemarkCategory {
	if merchantID != uuid.Nil {
		s.ensureMerchant(merchantID)
	}
	builder := s.client.RemarkCategory.Create().
		SetID(uuid.New()).
		SetName(name).
		SetRemarkScene(domain.RemarkSceneWholeOrder).
		SetSortOrder(0).
		SetDescription("desc")

	if merchantID != uuid.Nil {
		builder = builder.SetMerchantID(merchantID)
	}

	return builder.SaveX(s.ctx)
}

func (s *RemarkTestSuite) createRemark(name string, remarkType domain.RemarkType, merchantID, storeID uuid.UUID, sort int) *domain.Remark {
	return s.createRemarkWithCategory(name, remarkType, merchantID, storeID, sort, uuid.Nil)
}

func (s *RemarkTestSuite) createRemarkWithCategory(name string, remarkType domain.RemarkType, merchantID, storeID uuid.UUID, sort int, categoryID uuid.UUID) *domain.Remark {
	if merchantID != uuid.Nil {
		s.ensureMerchant(merchantID)
	}
	if storeID != uuid.Nil {
		s.ensureStore(storeID, merchantID)
	}

	catID := categoryID
	if catID == uuid.Nil {
		catID = s.createRemarkCategory(name+"-cat", merchantID).ID
	}
	remark := &domain.Remark{
		ID:         uuid.New(),
		Name:       name,
		RemarkType: remarkType,
		Enabled:    true,
		SortOrder:  sort,
		CategoryID: catID,
		MerchantID: merchantID,
		StoreID:    storeID,
	}
	require.NoError(s.T(), s.repo.Create(s.ctx, remark))
	return remark
}

// ensureMerchant creates a minimal merchant with required dependencies.
func (s *RemarkTestSuite) ensureMerchant(id uuid.UUID) *ent.Merchant {
	if id == uuid.Nil {
		return nil
	}
	if m, err := s.client.Merchant.Get(s.ctx, id); err == nil {
		return m
	}

	storeUser := s.client.StoreUser.Create().
		SetUsername("storeUser-" + id.String()).
		SetHashedPassword("pwd").
		SetNickname("storeUser").
		SaveX(s.ctx)

	bt := s.client.MerchantBusinessType.Create().
		SetTypeCode("code-" + id.String()).
		SetTypeName("name-" + id.String()).
		SaveX(s.ctx)

	return s.client.Merchant.Create().
		SetID(id).
		SetMerchantCode("MC-" + id.String()).
		SetMerchantName("Merchant-" + id.String()).
		SetMerchantShortName("M-" + id.String()).
		SetMerchantType(domain.MerchantTypeBrand).
		SetBrandName("Brand-" + id.String()).
		SetAdminPhoneNumber("13800000000").
		SetBusinessTypeID(bt.ID).
		SetMerchantLogo("logo").
		SetDescription("desc").
		SetStatus(domain.MerchantStatusActive).
		SetSuperAccount(storeUser.Username).
		SetAddress("addr").
		SetLng("0").
		SetLat("0").
		SaveX(s.ctx)
}

// ensureStore creates a minimal store under the given merchant.
func (s *RemarkTestSuite) ensureStore(id, merchantID uuid.UUID) *ent.Store {
	if id == uuid.Nil {
		return nil
	}
	if st, err := s.client.Store.Get(s.ctx, id); err == nil {
		return st
	}
	m := s.ensureMerchant(merchantID)

	storeUser := s.client.StoreUser.Create().
		SetUsername("store-storeUser-" + id.String()).
		SetHashedPassword("pwd").
		SetNickname("store-storeUser").
		SaveX(s.ctx)

	shortName := "S-" + id.String()
	if len(shortName) > 30 {
		shortName = shortName[:30]
	}
	name := "Store-" + id.String()
	if len(name) > 30 {
		name = name[:30]
	}

	return s.client.Store.Create().
		SetID(id).
		SetMerchantID(m.ID).
		SetAdminPhoneNumber("13800000001").
		SetStoreName(name).
		SetStoreShortName(shortName).
		SetStoreCode("SC-" + id.String()).
		SetStatus(domain.StoreStatusOpen).
		SetBusinessModel(domain.BusinessModelDirect).
		SetBusinessTypeID(m.BusinessTypeID).
		SetLocationNumber("loc").
		SetContactName("contact").
		SetContactPhone("13800000002").
		SetUnifiedSocialCreditCode("USC-" + id.String()).
		SetStoreLogo("logo").
		SetBusinessLicenseURL("bl").
		SetStorefrontURL("sf").
		SetCashierDeskURL("cd").
		SetDiningEnvironmentURL("de").
		SetFoodOperationLicenseURL("fo").
		SetBusinessHours([]domain.BusinessHours{}).
		SetDiningPeriods([]domain.DiningPeriod{}).
		SetShiftTimes([]domain.ShiftTime{}).
		SetAddress("addr").
		SetLng("0").
		SetLat("0").
		SetSuperAccount(storeUser.Username).
		SaveX(s.ctx)
}

func (s *RemarkTestSuite) TestRemark_CreateAndFind() {
	remark := s.createRemark("system-remark", domain.RemarkTypeSystem, uuid.Nil, uuid.Nil, 1)

	found, err := s.repo.FindByID(s.ctx, remark.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), remark.ID, found.ID)
	require.Equal(s.T(), remark.Name, found.Name)
	require.Equal(s.T(), remark.RemarkType, found.RemarkType)
	require.True(s.T(), found.Enabled)
}

func (s *RemarkTestSuite) TestRemark_CreateNilRemark() {
	err := s.repo.Create(s.ctx, nil)
	require.Error(s.T(), err)
}

func (s *RemarkTestSuite) TestRemark_Update() {
	remark := s.createRemark("brand-remark", domain.RemarkTypeBrand, uuid.New(), uuid.Nil, 2)

	remark.Name = "brand-remark-updated"
	remark.Enabled = false
	remark.SortOrder = 5
	err := s.repo.Update(s.ctx, remark)
	require.NoError(s.T(), err)

	updated := s.client.Remark.GetX(s.ctx, remark.ID)
	require.Equal(s.T(), "brand-remark-updated", updated.Name)
	require.False(s.T(), updated.Enabled)
	require.Equal(s.T(), 5, updated.SortOrder)
}

func (s *RemarkTestSuite) TestRemark_UpdateNotFound() {
	remark := &domain.Remark{ID: uuid.New(), Name: "missing"}
	err := s.repo.Update(s.ctx, remark)
	require.Error(s.T(), err)
	require.True(s.T(), domain.IsNotFound(err))
}

func (s *RemarkTestSuite) TestRemark_Delete() {
	remark := s.createRemark("to-delete", domain.RemarkTypeSystem, uuid.Nil, uuid.Nil, 3)

	require.NoError(s.T(), s.repo.Delete(s.ctx, remark.ID))

	_, err := s.client.Remark.Get(s.ctx, remark.ID)
	require.Error(s.T(), err)
	require.True(s.T(), ent.IsNotFound(err))
}

func (s *RemarkTestSuite) TestRemark_DeleteNotFound() {
	err := s.repo.Delete(s.ctx, uuid.New())
	require.Error(s.T(), err)
	require.True(s.T(), domain.IsNotFound(err))
}

func (s *RemarkTestSuite) TestRemark_GetRemarks_BrandTypeIncludesSystem() {
	merchantID := uuid.New()
	otherMerchantID := uuid.New()

	rBrand := s.createRemark("brand", domain.RemarkTypeBrand, merchantID, uuid.Nil, 1)
	rSystem := s.createRemark("system", domain.RemarkTypeSystem, uuid.Nil, uuid.Nil, 2)
	s.createRemark("other", domain.RemarkTypeBrand, otherMerchantID, uuid.Nil, 3)

	pager := upagination.New(1, 3)
	filter := &domain.RemarkListFilter{MerchantID: merchantID, RemarkType: domain.RemarkTypeBrand}
	remarks, total, err := s.repo.GetRemarks(s.ctx, pager, filter, domain.NewRemarkOrderBySortOrder(false))

	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, total)
	require.Len(s.T(), remarks, 2)
	require.Equal(s.T(), rBrand.ID, remarks[0].ID)
	require.Equal(s.T(), rSystem.ID, remarks[1].ID)
}

func (s *RemarkTestSuite) TestRemark_GetRemarks_MerchantOnlyExcludesSystemWhenTypeEmpty() {
	merchantID := uuid.New()

	rBrand := s.createRemark("brand", domain.RemarkTypeBrand, merchantID, uuid.Nil, 1)
	s.createRemark("system", domain.RemarkTypeSystem, uuid.Nil, uuid.Nil, 2)

	pager := upagination.New(1, 3)
	filter := &domain.RemarkListFilter{MerchantID: merchantID}
	remarks, total, err := s.repo.GetRemarks(s.ctx, pager, filter, domain.NewRemarkOrderBySortOrder(false))

	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, total)
	require.Len(s.T(), remarks, 1)
	require.Equal(s.T(), rBrand.ID, remarks[0].ID)
}

func (s *RemarkTestSuite) TestRemark_GetRemarks_StoreTypeIncludesStoreBrandAndSystem() {
	merchantID := uuid.New()
	storeID := uuid.New()

	rStore := s.createRemark("store", domain.RemarkTypeStore, merchantID, storeID, 1)
	rBrand := s.createRemark("brand", domain.RemarkTypeBrand, merchantID, uuid.Nil, 2)
	rSystem := s.createRemark("system", domain.RemarkTypeSystem, uuid.Nil, uuid.Nil, 3)

	pager := upagination.New(1, 5)
	filter := &domain.RemarkListFilter{MerchantID: merchantID, StoreID: storeID, RemarkType: domain.RemarkTypeStore}
	remarks, total, err := s.repo.GetRemarks(s.ctx, pager, filter, domain.NewRemarkOrderBySortOrder(false))

	require.NoError(s.T(), err)
	require.Equal(s.T(), 3, total)
	require.Len(s.T(), remarks, 3)
	require.Equal(s.T(), rStore.ID, remarks[0].ID)
	require.Equal(s.T(), rBrand.ID, remarks[1].ID)
	require.Equal(s.T(), rSystem.ID, remarks[2].ID)
}

func (s *RemarkTestSuite) TestRemark_GetRemarks_WithEnabledAndTypeFilter() {
	merchantID := uuid.New()
	enabled := true

	s.createRemark("brand-enabled", domain.RemarkTypeBrand, merchantID, uuid.Nil, 1)
	disabledRemark := s.createRemark("brand-disabled", domain.RemarkTypeBrand, merchantID, uuid.Nil, 2)
	require.NoError(s.T(), s.repo.Update(s.ctx, &domain.Remark{
		ID:         disabledRemark.ID,
		Name:       disabledRemark.Name,
		RemarkType: disabledRemark.RemarkType,
		Enabled:    false,
		SortOrder:  disabledRemark.SortOrder,
		CategoryID: disabledRemark.CategoryID,
		MerchantID: disabledRemark.MerchantID,
		StoreID:    disabledRemark.StoreID,
	}))

	filter := &domain.RemarkListFilter{MerchantID: merchantID, Enabled: &enabled, RemarkType: domain.RemarkTypeBrand}
	pager := upagination.New(1, 10)

	remarks, total, err := s.repo.GetRemarks(s.ctx, pager, filter)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, total)
	require.Len(s.T(), remarks, 1)
	require.Equal(s.T(), "brand-enabled", remarks[0].Name)
}

func (s *RemarkTestSuite) TestRemark_Exists() {
	merchantID := uuid.New()
	storeID := uuid.New()
	remark := s.createRemark("exists", domain.RemarkTypeBrand, merchantID, storeID, 1)

	params := domain.RemarkExistsParams{
		CategoryID: remark.CategoryID,
		MerchantID: merchantID,
		StoreID:    storeID,
		Name:       remark.Name,
	}

	exists, err := s.repo.Exists(s.ctx, params)
	require.NoError(s.T(), err)
	require.True(s.T(), exists)

	params.ExcludeID = remark.ID
	exists, err = s.repo.Exists(s.ctx, params)
	require.NoError(s.T(), err)
	require.False(s.T(), exists)

	params.Name = "not-exists"
	exists, err = s.repo.Exists(s.ctx, params)
	require.NoError(s.T(), err)
	require.False(s.T(), exists)
}

func (s *RemarkTestSuite) TestRemark_FindByID_NotFound() {
	_, err := s.repo.FindByID(s.ctx, uuid.New())
	require.Error(s.T(), err)
	require.True(s.T(), domain.IsNotFound(err))
}

func (s *RemarkTestSuite) TestRemark_CreateBrandWithoutMerchant() {
	remark := &domain.Remark{
		ID:         uuid.New(),
		Name:       "brand-no-merchant",
		RemarkType: domain.RemarkTypeBrand,
		Enabled:    true,
		SortOrder:  1,
		CategoryID: s.createRemarkCategory("brand-cat", uuid.New()).ID,
	}

	err := s.repo.Create(s.ctx, remark)
	require.Error(s.T(), err)
}

func (s *RemarkTestSuite) TestRemark_CountRemarkByCategories_FiltersByRemarkTypeAndMerchant() {
	merchantID := uuid.New()
	brandCat := s.createRemarkCategory("brand-cat", merchantID)
	systemCat := s.createRemarkCategory("system-cat", uuid.Nil)

	s.createRemarkWithCategory("brand-remark", domain.RemarkTypeBrand, merchantID, uuid.Nil, 1, brandCat.ID)
	s.createRemarkWithCategory("system-remark", domain.RemarkTypeSystem, uuid.Nil, uuid.Nil, 2, systemCat.ID)

	params := domain.CountRemarkParams{CategoryIDs: []uuid.UUID{brandCat.ID, systemCat.ID}, RemarkType: domain.RemarkTypeBrand, MerchantID: merchantID}
	counts, err := s.repo.CountRemarkByCategories(s.ctx, params)

	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, counts[brandCat.ID])
	require.Equal(s.T(), 1, counts[systemCat.ID])
}

func (s *RemarkTestSuite) TestRemark_CountRemarkByCategories_EmptyCategoryIDs() {
	counts, err := s.repo.CountRemarkByCategories(s.ctx, domain.CountRemarkParams{})

	require.NoError(s.T(), err)
	require.Empty(s.T(), counts)
}

func (s *RemarkTestSuite) TestRemark_GetRemarks_DefaultOrderByCreatedAtDesc() {
	first := s.createRemark("first", domain.RemarkTypeSystem, uuid.Nil, uuid.Nil, 1)
	time.Sleep(10 * time.Millisecond)
	second := s.createRemark("second", domain.RemarkTypeSystem, uuid.Nil, uuid.Nil, 2)

	pager := upagination.New(1, 10)
	remarks, total, err := s.repo.GetRemarks(s.ctx, pager, &domain.RemarkListFilter{})

	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, total)
	require.Len(s.T(), remarks, 2)
	require.Equal(s.T(), second.ID, remarks[0].ID)
	require.Equal(s.T(), first.ID, remarks[1].ID)
}
