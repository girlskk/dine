package merchant

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
	"gitlab.jiguang.dev/pos-dine/dine/usecase/store"
)

var _ domain.MerchantInteractor = (*MerchantInteractor)(nil)

type MerchantInteractor struct {
	DataStore domain.DataStore
}

func (interactor *MerchantInteractor) CreateMerchant(ctx context.Context, domainCMerchant *domain.CreateMerchantParams) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.CreateMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if domainCMerchant == nil {
		return domain.ParamsError(errors.New("domainCMerchant is required"))
	}

	domainMerchant, err := interactor.CheckCreateMerchantFields(ctx, domainCMerchant)
	if err != nil {
		return err
	}

	expireTime := domain.CalculateExpireTime(time.Now().UTC(), domainCMerchant.PurchaseDuration, domainCMerchant.PurchaseDurationUnit)
	domainMerchant.ExpireUTC = expireTime
	err = interactor.createMerchant(ctx, domainMerchant, nil)
	if err != nil {
		err = fmt.Errorf("failed to create merchant: %w", err)
		return
	}
	return
}

func (interactor *MerchantInteractor) CreateMerchantAndStore(ctx context.Context, domainCMerchant *domain.CreateMerchantParams, domainCStore *domain.CreateStoreParams) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.CreateMerchantAndStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	if domainCMerchant == nil {
		return domain.ParamsError(errors.New("domainCMerchant is required"))
	}

	if domainCStore == nil {
		return domain.ParamsError(errors.New("domainCStore is required"))
	}

	domainMerchant, err := interactor.CheckCreateMerchantFields(ctx, domainCMerchant)
	if err != nil {
		return err
	}

	storeInteractor := store.NewStoreInteractor(interactor.DataStore)
	domainStore, err := storeInteractor.CheckCreateStoreFields(ctx, domainCStore)
	if err != nil {
		return err
	}

	err = interactor.checkFields(ctx, domainMerchant)
	if err != nil {
		return
	}

	expireTime := domain.CalculateExpireTime(time.Now().UTC(), domainCMerchant.PurchaseDuration, domainCMerchant.PurchaseDurationUnit)
	domainMerchant.ExpireUTC = expireTime
	err = interactor.createMerchant(ctx, domainMerchant, domainStore)
	if err != nil {
		return err
	}
	return
}

func (interactor *MerchantInteractor) createMerchant(ctx context.Context, domainMerchant *domain.Merchant, domainStore *domain.Store) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.createMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		var err error
		adminUserID := uuid.New()
		err = ds.AdminUserRepo().Create(ctx, &domain.AdminUser{
			ID:             adminUserID,
			Username:       domainMerchant.LoginAccount,
			HashedPassword: domainMerchant.LoginPassword,
			AccountType:    domain.AdminUserAccountTypeSuperAdmin,
		})
		if err != nil {
			return err
		}
		domainMerchant.AdminUserID = adminUserID
		merchantID, err := ds.MerchantRepo().Create(ctx, domainMerchant)
		if err != nil {
			return err
		}

		if domainStore != nil {
			domainStore.AdminUserID = adminUserID
			domainStore.MerchantID = merchantID
			err = ds.StoreRepo().Create(ctx, domainStore)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return
}

func (interactor *MerchantInteractor) UpdateMerchant(ctx context.Context, domainUMerchant *domain.UpdateMerchantParams) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.UpdateMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if domainUMerchant == nil {
		return domain.ParamsError(errors.New("domainUMerchant is required"))
	}

	domainMerchant, err := interactor.CheckUpdateMerchantFields(ctx, domainUMerchant)
	if err != nil {
		return
	}

	err = interactor.updateMerchant(ctx, domainMerchant, nil)
	if err != nil {
		err = fmt.Errorf("failed to update merchant: %w", err)
		return err
	}
	return
}

func (interactor *MerchantInteractor) UpdateMerchantAndStore(ctx context.Context, domainUMerchant *domain.UpdateMerchantParams, domainUStore *domain.UpdateStoreParams) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.UpdateMerchantAndStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if domainUMerchant == nil {
		return domain.ParamsError(errors.New("domainUMerchant is required"))
	}

	if domainUStore == nil {
		return domain.ParamsError(errors.New("domainUStore is required"))
	}

	domainMerchant, err := interactor.CheckUpdateMerchantFields(ctx, domainUMerchant)
	if err != nil {
		return err
	}

	storeInteractor := store.NewStoreInteractor(interactor.DataStore)
	domainStore, err := storeInteractor.CheckUpdateStoreFields(ctx, domainUStore)
	if err != nil {
		return err
	}

	err = interactor.updateMerchant(ctx, domainMerchant, domainStore)
	if err != nil {
		return fmt.Errorf("failed to update merchant and store: %w", err)
	}
	return
}

func (interactor *MerchantInteractor) updateMerchant(ctx context.Context, domainMerchant *domain.Merchant, domainStore *domain.Store) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.UpdateMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		var err error
		err = ds.MerchantRepo().Update(ctx, domainMerchant)
		if err != nil {
			return err
		}

		adminUser := &domain.AdminUser{
			ID:             domainMerchant.AdminUserID,
			Username:       domainMerchant.LoginAccount,
			HashedPassword: domainMerchant.LoginPassword,
		}
		err = ds.AdminUserRepo().Update(ctx, adminUser)
		if err != nil {
			return err
		}
		if domainStore != nil {
			err = ds.StoreRepo().Update(ctx, domainStore)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (interactor *MerchantInteractor) DeleteMerchant(ctx context.Context, id int) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.DeleteMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	err = interactor.DataStore.MerchantRepo().Delete(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to delete merchant: %w", err)
		return
	}
	return
}

func (interactor *MerchantInteractor) GetMerchant(ctx context.Context, id int) (domainMerchant *domain.Merchant, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.GetMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainMerchant, err = interactor.DataStore.MerchantRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ParamsError(domain.ErrMerchantNotExists)
		}
		return nil, err
	}
	return
}

func (interactor *MerchantInteractor) GetMerchants(ctx context.Context, pager *upagination.Pagination, filter *domain.MerchantListFilter, orderBys ...domain.MerchantListOrderBy) (domainMerchants []*domain.Merchant, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.GetMerchants")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainMerchants, total, err = interactor.DataStore.MerchantRepo().GetMerchants(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get merchants: %w", err)
		return
	}

	return
}

func (interactor *MerchantInteractor) CountMerchant(ctx context.Context) (merchantCount *domain.MerchantCount, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.CountMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	merchantCount, err = interactor.DataStore.MerchantRepo().CountMerchant(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count merchants: %w", err)
		return
	}

	return
}

func (interactor *MerchantInteractor) MerchantRenewal(ctx context.Context, merchantRenewal *domain.MerchantRenewal) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantInteractor.MerchantRenewal")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		err := ds.MerchantRenewalRepo().Create(ctx, merchantRenewal)
		if err != nil {
			return err
		}
		m, err := ds.MerchantRepo().FindByID(ctx, merchantRenewal.MerchantID)
		if err != nil {
			return err
		}
		m.ExpireUTC = domain.CalculateExpireTime(*m.ExpireUTC, merchantRenewal.PurchaseDuration, merchantRenewal.PurchaseDurationUnit)
		err = ds.MerchantRepo().Update(ctx, m)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		err = fmt.Errorf("failed to renew merchant: %w", err)
		return
	}

	return
}

func (interactor *MerchantInteractor) CheckCreateMerchantFields(ctx context.Context, domainCMerchant *domain.CreateMerchantParams) (domainMerchant *domain.Merchant, err error) {
	domainMerchant = &domain.Merchant{
		MerchantCode:      domainCMerchant.MerchantCode,
		MerchantName:      domainCMerchant.MerchantName,
		MerchantShortName: domainCMerchant.MerchantShortName,
		MerchantType:      domainCMerchant.MerchantType,
		BrandName:         domainCMerchant.BrandName,
		AdminPhoneNumber:  domainCMerchant.AdminPhoneNumber,
		BusinessTypeID:    domainCMerchant.BusinessTypeID,
		MerchantLogo:      domainCMerchant.MerchantLogo,
		Description:       domainCMerchant.Description,
		Status:            domainCMerchant.Status,
		CountryID:         domainCMerchant.CountryID,
		ProvinceID:        domainCMerchant.ProvinceID,
		CityID:            domainCMerchant.CityID,
		DistrictID:        domainCMerchant.DistrictID,
		Address:           domainCMerchant.Address,
		Lng:               domainCMerchant.Lng,
		Lat:               domainCMerchant.Lat,
		LoginAccount:      domainCMerchant.LoginAccount,
		LoginPassword:     domainCMerchant.LoginPassword,
	}

	err = interactor.checkFields(ctx, domainMerchant)
	if err != nil {
		return
	}

	return
}

func (interactor *MerchantInteractor) CheckUpdateMerchantFields(ctx context.Context, domainUMerchant *domain.UpdateMerchantParams) (domainMerchant *domain.Merchant, err error) {
	oldMerchant, err := interactor.DataStore.MerchantRepo().FindByID(ctx, domainUMerchant.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrMerchantNotExists)
			return
		}
		return
	}
	domainMerchant = &domain.Merchant{
		ID:                domainUMerchant.ID,
		MerchantCode:      domainUMerchant.MerchantCode,
		MerchantName:      domainUMerchant.MerchantName,
		MerchantShortName: domainUMerchant.MerchantShortName,
		MerchantType:      oldMerchant.MerchantType,
		BrandName:         domainUMerchant.BrandName,
		AdminPhoneNumber:  domainUMerchant.AdminPhoneNumber,
		ExpireUTC:         oldMerchant.ExpireUTC,
		BusinessTypeID:    domainUMerchant.BusinessTypeID,
		MerchantLogo:      domainUMerchant.MerchantLogo,
		Description:       domainUMerchant.Description,
		Status:            domainUMerchant.Status,
		CountryID:         domainUMerchant.CountryID,
		ProvinceID:        domainUMerchant.ProvinceID,
		CityID:            domainUMerchant.CityID,
		DistrictID:        domainUMerchant.DistrictID,
		CountryName:       oldMerchant.CountryName,
		ProvinceName:      oldMerchant.ProvinceName,
		CityName:          oldMerchant.CityName,
		DistrictName:      oldMerchant.DistrictName,
		Address:           domainUMerchant.Address,
		Lng:               domainUMerchant.Lng,
		Lat:               domainUMerchant.Lat,
		LoginAccount:      domainUMerchant.LoginAccount,
		LoginPassword:     domainUMerchant.LoginPassword,
		AdminUserID:       oldMerchant.AdminUserID,
	}

	err = interactor.checkFields(ctx, domainMerchant)
	if err != nil {
		return
	}
	return
}

func (interactor *MerchantInteractor) checkFields(ctx context.Context, domainMerchant *domain.Merchant) (err error) {
	// 商户名称唯一性校验（update排除自身）
	exists, err := interactor.DataStore.MerchantRepo().ExistMerchant(ctx, &domain.MerchantExistsParams{
		MerchantName: domainMerchant.MerchantName,
		NotID:        domainMerchant.ID,
	})
	if err != nil {
		err = fmt.Errorf("failed to query merchant name existence: %w", err)
		return err
	}
	if exists {
		return domain.ParamsError(domain.ErrMerchantNameExists)
	}

	return
}

func NewMerchantInteractor(dataStore domain.DataStore) *MerchantInteractor {
	return &MerchantInteractor{
		DataStore: dataStore,
	}
}
