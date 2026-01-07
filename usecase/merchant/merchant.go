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

	userExists, err := interactor.loginUserExists(ctx, domainCMerchant.LoginAccount)
	if err != nil {
		return fmt.Errorf("failed to check login user existence: %w", err)
	}
	if userExists {
		return domain.ParamsError(domain.ErrUserExists)
	}
	expireTime := domain.CalculateExpireTime(time.Now().UTC(), domainCMerchant.PurchaseDuration, domainCMerchant.PurchaseDurationUnit)
	domainMerchant.ExpireUTC = expireTime
	domainMerchant.Status = domain.MerchantStatusActive
	err = interactor.createMerchant(ctx, domainMerchant, nil)
	if err != nil {
		err = fmt.Errorf("failed to create merchant: %w", err)
		return
	}
	return
}
func (interactor *MerchantInteractor) loginUserExists(ctx context.Context, loginUser string) (exists bool, err error) {
	return interactor.DataStore.BackendUserRepo().Exists(ctx, domain.BackendUserExistsParams{
		Username: loginUser,
	})
}
func (interactor *MerchantInteractor) CreateMerchantAndStore(ctx context.Context,
	domainCMerchant *domain.CreateMerchantParams,
	domainCStore *domain.CreateStoreParams,
) (err error) {
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

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		var err error
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

func (interactor *MerchantInteractor) UpdateMerchantAndStore(ctx context.Context,
	domainUMerchant *domain.UpdateMerchantParams,
	domainUStore *domain.UpdateStoreParams,
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

func (interactor *MerchantInteractor) MerchantSimpleUpdate(ctx context.Context,
	updateField domain.MerchantSimpleUpdateType,
	domainMerchant *domain.Merchant,
) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "MerchantInteractor.MerchantSimpleUpdate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if domainMerchant == nil {
		return domain.ParamsError(errors.New("domainMerchant is required"))
	}
	merchant, err := interactor.DataStore.MerchantRepo().FindByID(ctx, domainMerchant.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrMerchantNotExists)
			return
		}
		return
	}
	switch updateField {
	case domain.MerchantSimpleUpdateTypeStatus:
		if merchant.Status == domainMerchant.Status {
			return
		}
		merchant.Status = domainMerchant.Status

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

		// 此处不再需要修改密码
		//oldAccount, err := ds.BackendUserRepo().FindByUsername(ctx, domainMerchant.LoginAccount)
		//if err != nil {
		//	return err
		//}
		//
		//if util.CheckPassword(domainMerchant.LoginPassword, oldAccount.HashedPassword) != nil {
		//	hashPwd, err := util.HashPassword(domainMerchant.LoginPassword)
		//	if err != nil {
		//		return err
		//	}
		//	oldAccount.HashedPassword = hashPwd
		//	err = ds.BackendUserRepo().Update(ctx, oldAccount)
		//	if err != nil {
		//		return err
		//	}
		//}

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

	if bt, err := interactor.DataStore.MerchantBusinessTypeRepo().FindByCode(ctx, domainMerchant.BusinessTypeCode); err == nil {
		domainMerchant.BusinessTypeName = bt.TypeName
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

	domainMerchants, total, err = interactor.DataStore.MerchantRepo().GetMerchants(ctx, pager, filter, orderBys...)
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

func (interactor *MerchantInteractor) CheckCreateMerchantFields(ctx context.Context,
	domainCMerchant *domain.CreateMerchantParams,
) (domainMerchant *domain.Merchant, err error) {
	domainMerchant = &domain.Merchant{
		MerchantCode:      domainCMerchant.MerchantCode,
		MerchantName:      domainCMerchant.MerchantName,
		MerchantShortName: domainCMerchant.MerchantShortName,
		MerchantType:      domainCMerchant.MerchantType,
		BrandName:         domainCMerchant.BrandName,
		AdminPhoneNumber:  domainCMerchant.AdminPhoneNumber,
		BusinessTypeCode:  domainCMerchant.BusinessTypeCode,
		MerchantLogo:      domainCMerchant.MerchantLogo,
		Description:       domainCMerchant.Description,
		LoginAccount:      domainCMerchant.LoginAccount,
		LoginPassword:     domainCMerchant.LoginPassword,
		Address:           domainCMerchant.Address,
	}

	err = interactor.exists(ctx, domainMerchant)
	if err != nil {
		return
	}

	return
}

func (interactor *MerchantInteractor) CheckUpdateMerchantFields(ctx context.Context,
	domainUMerchant *domain.UpdateMerchantParams,
) (domainMerchant *domain.Merchant, err error) {
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
		BusinessTypeCode:  domainUMerchant.BusinessTypeCode,
		MerchantLogo:      domainUMerchant.MerchantLogo,
		Description:       domainUMerchant.Description,
		Status:            oldMerchant.Status,
		Address:           domainUMerchant.Address,
		LoginAccount:      oldMerchant.LoginAccount,
	}

	err = interactor.exists(ctx, domainMerchant)
	if err != nil {
		return
	}
	return
}

func (interactor *MerchantInteractor) exists(ctx context.Context, domainMerchant *domain.Merchant) (err error) {
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
