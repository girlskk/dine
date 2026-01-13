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
	DS domain.DataStore
}

func NewMerchantInteractor(dataStore domain.DataStore) *MerchantInteractor {
	return &MerchantInteractor{
		DS: dataStore,
	}
}

func (interactor *MerchantInteractor) CreateMerchant(ctx context.Context, domainCMerchant *domain.CreateMerchantParams, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.CreateMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainMerchant := &domain.Merchant{
		MerchantCode:         domainCMerchant.MerchantCode,
		MerchantName:         domainCMerchant.MerchantName,
		MerchantShortName:    domainCMerchant.MerchantShortName,
		MerchantType:         domainCMerchant.MerchantType,
		BrandName:            domainCMerchant.BrandName,
		AdminPhoneNumber:     domainCMerchant.AdminPhoneNumber,
		BusinessTypeCode:     domainCMerchant.BusinessTypeCode,
		MerchantLogo:         domainCMerchant.MerchantLogo,
		Description:          domainCMerchant.Description,
		LoginAccount:         domainCMerchant.LoginAccount,
		LoginPassword:        domainCMerchant.LoginPassword,
		Address:              domainCMerchant.Address,
		PurchaseDuration:     domainCMerchant.PurchaseDuration,
		PurchaseDurationUnit: domainCMerchant.PurchaseDurationUnit,
	}

	expireTime := domain.CalculateExpireTime(time.Now().UTC(), domainCMerchant.PurchaseDuration, domainCMerchant.PurchaseDurationUnit)
	domainMerchant.ExpireUTC = expireTime
	domainMerchant.Status = domain.MerchantStatusActive
	err = interactor.createMerchant(ctx, domainMerchant, nil)
	if err != nil {
		return err
	}
	return
}

func (interactor *MerchantInteractor) CreateMerchantAndStore(ctx context.Context,
	domainCMerchant *domain.CreateMerchantParams,
	domainCStore *domain.CreateStoreParams,
	user domain.User,
) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.CreateMerchantAndStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainMerchant := &domain.Merchant{
		MerchantCode:         domainCMerchant.MerchantCode,
		MerchantName:         domainCMerchant.MerchantName,
		MerchantShortName:    domainCMerchant.MerchantShortName,
		MerchantType:         domainCMerchant.MerchantType,
		BrandName:            domainCMerchant.BrandName,
		AdminPhoneNumber:     domainCMerchant.AdminPhoneNumber,
		BusinessTypeCode:     domainCMerchant.BusinessTypeCode,
		MerchantLogo:         domainCMerchant.MerchantLogo,
		Description:          domainCMerchant.Description,
		LoginAccount:         domainCMerchant.LoginAccount,
		LoginPassword:        domainCMerchant.LoginPassword,
		Address:              domainCMerchant.Address,
		PurchaseDuration:     domainCMerchant.PurchaseDuration,
		PurchaseDurationUnit: domainCMerchant.PurchaseDurationUnit,
	}

	storeInteractor := store.NewStoreInteractor(interactor.DS)
	domainStore, err := storeInteractor.CheckCreateStoreFields(ctx, domainCStore)
	if err != nil {
		return err
	}

	expireTime := domain.CalculateExpireTime(time.Now().UTC(), domainCMerchant.PurchaseDuration, domainCMerchant.PurchaseDurationUnit)
	domainMerchant.ExpireUTC = expireTime
	domainMerchant.Status = domain.MerchantStatusActive
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

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.MerchantRepo().ExistMerchant(ctx, &domain.MerchantExistsParams{
			MerchantName: domainMerchant.MerchantName,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrMerchantNameExists
		}
		userExists, err := ds.BackendUserRepo().Exists(ctx, domain.BackendUserExistsParams{
			Username: domainMerchant.LoginAccount,
		})
		if err != nil {
			return err
		}
		if userExists {
			return domain.ErrUsernameExist
		}
		merchantID := uuid.New()
		err = ds.BackendUserRepo().Create(ctx, &domain.BackendUser{
			ID:             uuid.New(),
			Username:       domainMerchant.LoginAccount,
			HashedPassword: domainMerchant.LoginPassword,
			Nickname:       "admin",
			MerchantID:     merchantID,
			Code:           merchantID.String(),
			RealName:       "admin",
			Gender:         domain.GenderOther,
			Email:          "",
			PhoneNumber:    "",
			Enabled:        true,
			IsSuperAdmin:   true,
		})
		if err != nil {
			return err
		}
		domainMerchant.ID = merchantID
		err = ds.MerchantRepo().Create(ctx, domainMerchant)
		if err != nil {
			return err
		}
		if domainStore != nil {
			storeExists, err := ds.StoreRepo().ExistsStore(ctx, &domain.ExistsStoreParams{
				MerchantID: domainStore.MerchantID,
				StoreName:  domainStore.StoreName,
			})
			if err != nil {
				return err
			}
			if storeExists {
				return domain.ErrStoreNameExists
			}

			storeUserExists, err := ds.StoreUserRepo().Exists(ctx, domain.StoreUserExistsParams{
				Username: domainStore.LoginAccount,
			})
			if err != nil {
				return err
			}
			if storeUserExists {
				return domain.ErrUsernameExist
			}
			storeID := uuid.New()
			err = ds.StoreUserRepo().Create(ctx, &domain.StoreUser{
				ID:             uuid.New(),
				Username:       domainMerchant.LoginAccount,
				HashedPassword: domainMerchant.LoginPassword,
				Nickname:       "admin",
				MerchantID:     merchantID,
				StoreID:        storeID,
				Code:           storeID.String(),
				RealName:       "admin",
				Gender:         domain.GenderOther,
				Email:          "",
				PhoneNumber:    "",
				Enabled:        true,
				IsSuperAdmin:   true,
			})
			if err != nil {
				return err
			}
			domainStore.MerchantID = merchantID
			domainStore.ID = storeID
			err = ds.StoreRepo().Create(ctx, domainStore)
			if err != nil {
				return err
			}

		}
		err = ds.MerchantRenewalRepo().Create(ctx, &domain.MerchantRenewal{
			ID:                   uuid.New(),
			MerchantID:           merchantID,
			PurchaseDuration:     domainMerchant.PurchaseDuration,
			PurchaseDurationUnit: domainMerchant.PurchaseDurationUnit,
			OperatorName:         "",
			OperatorAccount:      "",
		})
		if err != nil {
			return err
		}
		return nil
	})

	return
}

func (interactor *MerchantInteractor) UpdateMerchant(ctx context.Context, domainUMerchant *domain.UpdateMerchantParams, user domain.User) (err error) {
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

	// verify ownership of the merchant being updated
	if err = verifyMerchantOwnership(user, domainMerchant.ID); err != nil {
		return err
	}

	err = interactor.updateMerchant(ctx, domainMerchant, nil)
	if err != nil {
		return err
	}
	return
}

func (interactor *MerchantInteractor) UpdateMerchantAndStore(ctx context.Context,
	domainUMerchant *domain.UpdateMerchantParams,
	domainUStore *domain.UpdateStoreParams,
	user domain.User,
) (err error) {
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

	storeInteractor := store.NewStoreInteractor(interactor.DS)
	domainStore, err := storeInteractor.CheckUpdateStoreFields(ctx, domainUStore)
	if err != nil {
		return err
	}

	// verify ownership of the merchant being updated
	if err = verifyMerchantOwnership(user, domainMerchant.ID); err != nil {
		return err
	}

	err = interactor.updateMerchant(ctx, domainMerchant, domainStore)
	if err != nil {
		return fmt.Errorf("failed to update merchant and store: %w", err)
	}
	return
}

func (interactor *MerchantInteractor) MerchantSimpleUpdate(ctx context.Context,
	updateField domain.MerchantSimpleUpdateField,
	domainMerchant *domain.Merchant,
	user domain.User,
) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.MerchantSimpleUpdate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// verify ownership
	if err = verifyMerchantOwnership(user, domainMerchant.ID); err != nil {
		return err
	}

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		merchant, err := ds.MerchantRepo().FindByID(ctx, domainMerchant.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrMerchantNotExists
			}
			return err
		}
		switch updateField {
		case domain.MerchantSimpleUpdateTypeStatus:
			if merchant.Status == domainMerchant.Status {
				return nil
			}
			merchant.Status = domainMerchant.Status

		default:
			return fmt.Errorf("unsupported update field: %v", updateField)
		}
		err = ds.MerchantRepo().Update(ctx, merchant)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return
}

func (interactor *MerchantInteractor) updateMerchant(ctx context.Context, domainMerchant *domain.Merchant, domainStore *domain.Store) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.updateMerchant")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.MerchantRepo().ExistMerchant(ctx, &domain.MerchantExistsParams{
			MerchantName: domainMerchant.MerchantName,
			ExcludeID:    domainMerchant.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrMerchantNameExists
		}
		err = ds.MerchantRepo().Update(ctx, domainMerchant)
		if err != nil {
			return err
		}

		if domainStore != nil {
			storeExists, err := ds.StoreRepo().ExistsStore(ctx, &domain.ExistsStoreParams{
				MerchantID: domainStore.MerchantID,
				StoreName:  domainStore.StoreName,
				ExcludeID:  domainStore.ID,
			})
			if err != nil {
				return err
			}
			if storeExists {
				return domain.ErrStoreNameExists
			}
			err = ds.StoreRepo().Update(ctx, domainStore)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (interactor *MerchantInteractor) DeleteMerchant(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.DeleteMerchant")
	defer func() { util.SpanErrFinish(span, err) }()
	// verify ownership
	if err = verifyMerchantOwnership(user, id); err != nil {
		return err
	}
	_, err = interactor.DS.MerchantRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrMerchantNotExists)
		}
		return
	}
	err = interactor.DS.MerchantRepo().Delete(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to delete merchant: %w", err)
		return
	}
	return
}

func (interactor *MerchantInteractor) GetMerchant(ctx context.Context, id uuid.UUID, user domain.User) (domainMerchant *domain.Merchant, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.GetMerchant")
	defer func() { util.SpanErrFinish(span, err) }()

	// verify ownership
	if err = verifyMerchantOwnership(user, id); err != nil {
		return nil, err
	}
	domainMerchant, err = interactor.DS.MerchantRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ParamsError(domain.ErrMerchantNotExists)
		}
		return nil, err
	}

	if renewals, err := interactor.DS.MerchantRenewalRepo().GetByMerchant(ctx, id); err == nil && len(renewals) > 0 {
		latestRenewal := renewals[0]
		domainMerchant.PurchaseDuration = latestRenewal.PurchaseDuration
		domainMerchant.PurchaseDurationUnit = latestRenewal.PurchaseDurationUnit
	}
	return domainMerchant, nil
}

func (interactor *MerchantInteractor) GetMerchants(ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.MerchantListFilter,
	orderBys ...domain.MerchantListOrderBy,
) (domainMerchants []*domain.Merchant, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.GetMerchants")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainMerchants, total, err = interactor.DS.MerchantRepo().GetMerchants(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get merchants: %w", err)
		return
	}
	merchantIds := lo.Map(domainMerchants, func(item *domain.Merchant, _ int) uuid.UUID {
		return item.ID
	})

	if len(merchantIds) == 0 {
		return
	}
	merchantStoreCounts, err := interactor.DS.StoreRepo().CountStoresByMerchantID(ctx, merchantIds)
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
	merchantCount, err = interactor.DS.MerchantRepo().CountMerchant(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count merchants: %w", err)
		return
	}

	return
}

func (interactor *MerchantInteractor) MerchantRenewal(ctx context.Context, merchantRenewal *domain.MerchantRenewal, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.MerchantRenewal")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// verify ownership
	if err = verifyMerchantOwnership(user, merchantRenewal.MerchantID); err != nil {
		return err
	}

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
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
			if m.ExpireUTC.After(oldExpireTime) {
				oldExpireTime = *m.ExpireUTC
			}
		}
		m.Status = domain.MerchantStatusActive
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

func (interactor *MerchantInteractor) CheckUpdateMerchantFields(ctx context.Context,
	domainUMerchant *domain.UpdateMerchantParams,
) (domainMerchant *domain.Merchant, err error) {
	oldMerchant, err := interactor.DS.MerchantRepo().FindByID(ctx, domainUMerchant.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ErrMerchantNotExists
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
		BusinessTypeCode:  domainUMerchant.BusinessTypeCode,
		MerchantLogo:      domainUMerchant.MerchantLogo,
		Description:       domainUMerchant.Description,
		Status:            oldMerchant.Status,
		Address:           domainUMerchant.Address,
		LoginAccount:      oldMerchant.LoginAccount,
	}

	return
}

func verifyMerchantOwnership(user domain.User, merchantID uuid.UUID) error {
	switch user.GetUserType() {
	case domain.UserTypeAdmin:
	case domain.UserTypeBackend:
		if !domain.VerifyOwnerMerchant(user, merchantID) {
			return domain.ErrMerchantNotExists
		}
	case domain.UserTypeStore:
		if !domain.VerifyOwnerShip(user, merchantID, uuid.Nil) {
			return domain.ErrMerchantNotExists
		}
	}

	return nil
}
