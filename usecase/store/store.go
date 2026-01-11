package store

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/i18n"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StoreInteractor = (*StoreInteractor)(nil)

type StoreInteractor struct {
	DS domain.DataStore
}

func (interactor *StoreInteractor) CreateStore(ctx context.Context, domainCStore *domain.CreateStoreParams) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.CreateStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainStore, err := interactor.CheckCreateStoreFields(ctx, domainCStore)
	if err != nil {
		return
	}

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.StoreRepo().ExistsStore(ctx, &domain.ExistsStoreParams{
			MerchantID: domainStore.MerchantID,
			StoreName:  domainStore.StoreName,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrStoreNameExists
		}

		userExists, err := ds.StoreUserRepo().Exists(ctx, domain.StoreUserExistsParams{
			Username: domainStore.LoginAccount,
		})
		if err != nil {
			return err
		}
		if userExists {
			return domain.ErrUsernameExist
		}
		storeID := uuid.New()
		err = ds.StoreUserRepo().Create(ctx, &domain.StoreUser{
			ID:             uuid.New(),
			Username:       domainStore.LoginAccount,
			HashedPassword: domainStore.LoginPassword,
			Nickname:       "admin",
			MerchantID:     domainStore.MerchantID,
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

		domainStore.ID = storeID
		err = ds.StoreRepo().Create(ctx, domainStore)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return
	}
	return
}

func (interactor *StoreInteractor) UpdateStore(ctx context.Context, domainUStore *domain.UpdateStoreParams) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.UpdateStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainStore, err := interactor.CheckUpdateStoreFields(ctx, domainUStore)
	if err != nil {
		return
	}

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.StoreRepo().ExistsStore(ctx, &domain.ExistsStoreParams{
			MerchantID: domainStore.MerchantID,
			StoreName:  domainStore.StoreName,
			ExcludeID:  domainStore.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrStoreNameExists
		}
		err = ds.StoreRepo().Update(ctx, domainStore)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

func (interactor *StoreInteractor) DeleteStore(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.DeleteStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return interactor.DS.StoreRepo().Delete(ctx, id)
}

func (interactor *StoreInteractor) GetStore(ctx context.Context, id uuid.UUID) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.GetStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	domainStore, err = interactor.DS.StoreRepo().FindByID(ctx, id)
	if err != nil {
		return
	}
	if msgID, ok := domain.BusinessTypeI18NMap[domainStore.BusinessTypeCode]; ok {
		name := i18n.Translate(ctx, msgID, nil)
		domainStore.BusinessTypeName = name
	}
	return
}

func (interactor *StoreInteractor) GetStoreByMerchantID(ctx context.Context, merchantID uuid.UUID) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.GetStoreByMerchantID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainStore, err = interactor.DS.StoreRepo().FindStoreMerchant(ctx, merchantID)
	if err != nil {
		return
	}
	if msgID, ok := domain.BusinessTypeI18NMap[domainStore.BusinessTypeCode]; ok {
		name := i18n.Translate(ctx, msgID, nil)
		domainStore.BusinessTypeName = name
	}

	return
}

func (interactor *StoreInteractor) GetStores(ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.StoreListFilter,
	orderBys ...domain.StoreListOrderBy,
) (domainStores []*domain.Store, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.GetStores")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	domainStores, total, err = interactor.DS.StoreRepo().GetStores(ctx, pager, filter, orderBys...)
	if err != nil {
		return
	}
	return
}

func (interactor *StoreInteractor) StoreSimpleUpdate(ctx context.Context,
	updateField domain.StoreSimpleUpdateField,
	domainUStoreParams *domain.UpdateStoreParams,
) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.StoreSimpleUpdate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		domainStore, err := ds.StoreRepo().FindByID(ctx, domainUStoreParams.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrStoreNotExists
			}
			return err
		}
		switch updateField {
		case domain.StoreSimpleUpdateFieldStatus:
			if domainStore.Status == domainUStoreParams.Status {
				return nil
			}
			domainStore.Status = domainUStoreParams.Status
		default:
			return fmt.Errorf("unsupported update field: %v", updateField)
		}
		err = ds.StoreRepo().Update(ctx, domainStore)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

func (interactor *StoreInteractor) CheckCreateStoreFields(ctx context.Context,
	domainCStore *domain.CreateStoreParams,
) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.CheckCreateStoreFields")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainStore = &domain.Store{
		MerchantID:              domainCStore.MerchantID,
		AdminPhoneNumber:        domainCStore.AdminPhoneNumber,
		StoreName:               domainCStore.StoreName,
		StoreShortName:          domainCStore.StoreShortName,
		StoreCode:               domainCStore.StoreCode,
		Status:                  domainCStore.Status,
		BusinessModel:           domainCStore.BusinessModel,
		BusinessTypeCode:        domainCStore.BusinessTypeCode,
		LocationNumber:          domainCStore.LocationNumber,
		ContactName:             domainCStore.ContactName,
		ContactPhone:            domainCStore.ContactPhone,
		UnifiedSocialCreditCode: domainCStore.UnifiedSocialCreditCode,
		StoreLogo:               domainCStore.StoreLogo,
		BusinessLicenseURL:      domainCStore.BusinessLicenseURL,
		StorefrontURL:           domainCStore.StorefrontURL,
		CashierDeskURL:          domainCStore.CashierDeskURL,
		DiningEnvironmentURL:    domainCStore.DiningEnvironmentURL,
		FoodOperationLicenseURL: domainCStore.FoodOperationLicenseURL,
		BusinessHours:           domainCStore.BusinessHours,
		DiningPeriods:           domainCStore.DiningPeriods,
		ShiftTimes:              domainCStore.ShiftTimes,
		Address:                 domainCStore.Address,
		LoginAccount:            domainCStore.LoginAccount,
		LoginPassword:           domainCStore.LoginPassword,
	}

	if err = interactor.validateTimeConfigs(domainStore); err != nil {
		return nil, err
	}

	return
}

func (interactor *StoreInteractor) CheckUpdateStoreFields(ctx context.Context,
	domainUStore *domain.UpdateStoreParams,
) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.CheckUpdateStoreFields")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	oldStore, err := interactor.DS.StoreRepo().FindByID(ctx, domainUStore.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ErrStoreNotExists
			return
		}
		return
	}

	domainStore = &domain.Store{
		ID:                      domainUStore.ID,
		MerchantID:              oldStore.MerchantID,
		AdminPhoneNumber:        domainUStore.AdminPhoneNumber,
		StoreName:               domainUStore.StoreName,
		StoreShortName:          domainUStore.StoreShortName,
		StoreCode:               domainUStore.StoreCode,
		Status:                  domainUStore.Status,
		BusinessModel:           domainUStore.BusinessModel,
		BusinessTypeCode:        domainUStore.BusinessTypeCode,
		LocationNumber:          domainUStore.LocationNumber,
		ContactName:             domainUStore.ContactName,
		ContactPhone:            domainUStore.ContactPhone,
		UnifiedSocialCreditCode: domainUStore.UnifiedSocialCreditCode,
		StoreLogo:               domainUStore.StoreLogo,
		BusinessLicenseURL:      domainUStore.BusinessLicenseURL,
		StorefrontURL:           domainUStore.StorefrontURL,
		CashierDeskURL:          domainUStore.CashierDeskURL,
		DiningEnvironmentURL:    domainUStore.DiningEnvironmentURL,
		FoodOperationLicenseURL: domainUStore.FoodOperationLicenseURL,
		BusinessHours:           domainUStore.BusinessHours,
		DiningPeriods:           domainUStore.DiningPeriods,
		ShiftTimes:              domainUStore.ShiftTimes,
		Address:                 domainUStore.Address,
		LoginAccount:            oldStore.LoginAccount,
	}

	if err = interactor.validateTimeConfigs(domainStore); err != nil {
		return nil, err
	}

	return
}

func (interactor *StoreInteractor) validateTimeConfigs(store *domain.Store) error {
	if err := validateBusinessHours(store.BusinessHours); err != nil {
		return err
	}
	if err := validateDiningPeriods(store.DiningPeriods); err != nil {
		return err
	}
	if err := validateShiftTimes(store.ShiftTimes); err != nil {
		return err
	}
	return nil
}

func validateBusinessHours(hours []domain.BusinessHours) error {
	if len(hours) == 0 {
		return nil
	}
	perDay := make(map[time.Weekday][][2]string)
	for _, h := range hours {
		// h now contains multiple BusinessHour entries
		for _, bh := range h.BusinessHours {
			if bh.StartTime >= bh.EndTime {
				return domain.ParamsError(domain.ErrStoreBusinessHoursTimeInvalid)
			}
			for _, wd := range h.Weekdays {
				perDay[wd] = append(perDay[wd], [2]string{bh.StartTime, bh.EndTime})
			}
		}
	}
	for _, ranges := range perDay {
		sort.Slice(ranges, func(i, j int) bool { return ranges[i][0] < ranges[j][0] })
		for i := 1; i < len(ranges); i++ {
			prev, cur := ranges[i-1], ranges[i]
			if cur[0] < prev[1] {
				return domain.ParamsError(domain.ErrStoreBusinessHoursConflict)
			}
		}
	}
	return nil
}

func validateDiningPeriods(periods []domain.DiningPeriod) error {
	if len(periods) == 0 {
		return nil
	}
	seen := make(map[string]struct{})
	for _, p := range periods {
		if p.StartTime >= p.EndTime {
			return domain.ParamsError(domain.ErrStoreDiningPeriodTimeInvalid)
		}
		if _, ok := seen[p.Name]; ok {
			return domain.ParamsError(domain.ErrStoreDiningPeriodNameExists)
		}
		seen[p.Name] = struct{}{}
	}

	sort.Slice(periods, func(i, j int) bool { return periods[i].StartTime < periods[j].StartTime })
	for i := 1; i < len(periods); i++ {
		prev, cur := periods[i-1], periods[i]
		if cur.StartTime < prev.EndTime {
			return domain.ParamsError(domain.ErrStoreDiningPeriodConflict)
		}
	}
	return nil
}

func validateShiftTimes(shifts []domain.ShiftTime) error {
	if len(shifts) == 0 {
		return nil
	}
	seen := make(map[string]struct{})
	for _, s := range shifts {
		if s.StartTime >= s.EndTime {
			return domain.ParamsError(domain.ErrStoreShiftTimeTimeInvalid)
		}
		if _, ok := seen[s.Name]; ok {
			return domain.ParamsError(domain.ErrStoreShiftTimeNameExists)
		}
		seen[s.Name] = struct{}{}
	}

	sort.Slice(shifts, func(i, j int) bool { return shifts[i].StartTime < shifts[j].StartTime })
	for i := 1; i < len(shifts); i++ {
		prev, cur := shifts[i-1], shifts[i]
		if cur.StartTime < prev.EndTime {
			return domain.ParamsError(domain.ErrStoreShiftTimeConflict)
		}
	}
	return nil
}

func NewStoreInteractor(dataStore domain.DataStore) *StoreInteractor {
	return &StoreInteractor{
		DS: dataStore,
	}
}
