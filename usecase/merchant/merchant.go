package merchant

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
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
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.CreateMerchant")
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
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.CreateMerchantAndStore")
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
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.createMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		var err error
		adminUserID := uuid.New()
		merchantID := uuid.New()

		err = ds.AdminUserRepo().Create(ctx, &domain.AdminUser{
			ID:             adminUserID,
			Username:       domainMerchant.LoginAccount,
			HashedPassword: domainMerchant.LoginPassword,
			AccountType:    domain.AdminUserAccountTypeSuperAdmin,
		})
		if err != nil {
			return err
		}
		domainMerchant.ID = merchantID
		domainMerchant.AdminUserID = adminUserID
		err = ds.MerchantRepo().Create(ctx, domainMerchant)
		if err != nil {
			return err
		}

		if domainStore != nil {
			domainStore.AdminUserID = adminUserID
			domainStore.MerchantID = merchantID
			domainStore.ID = uuid.New()
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
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.UpdateMerchant")
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
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.UpdateMerchantAndStore")
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

func (interactor *MerchantInteractor) MerchantSimpleUpdate(ctx context.Context, updateField domain.MerchantSimpleUpdateType, domainUMerchant *domain.UpdateMerchantParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.MerchantSimpleUpdate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if domainUMerchant == nil {
		return domain.ParamsError(errors.New("domainUMerchant is required"))
	}
	merchant, err := interactor.DataStore.MerchantRepo().FindByID(ctx, domainUMerchant.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrMerchantNotExists)
			return
		}
		return
	}
	switch updateField {
	case domain.MerchantSimpleUpdateTypeStatus:
		if merchant.Status == domainUMerchant.Status {
			return
		}
		merchant.Status = domainUMerchant.Status

	default:
		return domain.ParamsError(fmt.Errorf("unsupported update field: %v", updateField))
	}

	return interactor.DataStore.MerchantRepo().Update(ctx, merchant)
}

func (interactor *MerchantInteractor) updateMerchant(ctx context.Context, domainMerchant *domain.Merchant, domainStore *domain.Store) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.updateMerchant")
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

func (interactor *MerchantInteractor) DeleteMerchant(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.DeleteMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	_, err = interactor.DataStore.MerchantRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrMerchantNotExists)
		}
		return
	}
	err = interactor.DataStore.MerchantRepo().Delete(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to delete merchant: %w", err)
		return
	}
	return
}

func (interactor *MerchantInteractor) GetMerchant(ctx context.Context, id uuid.UUID) (domainMerchant *domain.Merchant, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.GetMerchant")
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
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.GetMerchants")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainMerchants, total, err = interactor.DataStore.MerchantRepo().GetMerchants(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get merchants: %w", err)
		return
	}
	merchantIds := lo.Map(domainMerchants, func(item *domain.Merchant, _ int) uuid.UUID {
		return item.ID
	})

	merchantStoreCounts, err := interactor.DataStore.StoreRepo().CountStoresByMerchantID(ctx, merchantIds)
	if err != nil {
		err = fmt.Errorf("failed to count stores by merchant ids: %w", err)
		return
	}

	storeCountMap := lo.SliceToMap(merchantStoreCounts, func(item *domain.MerchantStoreCount) (uuid.UUID, int) {
		return item.MerchantID, item.StoreCount
	})

	for _, m := range domainMerchants {
		if count, ok := storeCountMap[m.ID]; ok {
			m.StoreCount = count
		}
	}
	return
}

func (interactor *MerchantInteractor) CountMerchant(ctx context.Context) (merchantCount *domain.MerchantCount, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.CountMerchant")
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
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.MerchantRenewal")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		var err error
		merchantRenewal.ID = uuid.New()
		err = ds.MerchantRenewalRepo().Create(ctx, merchantRenewal)
		if err != nil {
			return err
		}
		m, err := ds.MerchantRepo().FindByID(ctx, merchantRenewal.MerchantID)
		if err != nil {
			return err
		}
		oldExpireTime := time.Now().UTC()
		if m.ExpireUTC != nil {
			oldExpireTime = *m.ExpireUTC
		}
		m.ExpireUTC = domain.CalculateExpireTime(oldExpireTime, merchantRenewal.PurchaseDuration, merchantRenewal.PurchaseDurationUnit)
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
		LoginAccount:      domainCMerchant.LoginAccount,
		LoginPassword:     domainCMerchant.LoginPassword,
		Address:           domainCMerchant.Address,
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
		Address:           domainUMerchant.Address,
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
		ExcludeID:    domainMerchant.ID,
	})
	if err != nil {
		err = fmt.Errorf("failed to query merchant name existence: %w", err)
		return err
	}
	if exists {
		return domain.ConflictError(domain.ErrMerchantNameExists)
	}

	return
}

func NewMerchantInteractor(dataStore domain.DataStore) *MerchantInteractor {
	return &MerchantInteractor{
		DataStore: dataStore,
	}
}
