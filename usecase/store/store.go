package store

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StoreInteractor = (*StoreInteractor)(nil)

type StoreInteractor struct {
	DataStore domain.DataStore
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

	domainStore.ID = uuid.New()
	err = interactor.DataStore.StoreRepo().Create(ctx, domainStore)
	if err != nil {
		err = fmt.Errorf("failed to create store: %w", err)
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
	err = interactor.DataStore.StoreRepo().Update(ctx, domainStore)
	if err != nil {
		err = fmt.Errorf("failed to update store: %w", err)
		return
	}

	return
}

func (interactor *StoreInteractor) DeleteStore(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.DeleteStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	_, err = interactor.DataStore.StoreRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrStoreNotExists)
			return
		}
		return
	}

	err = interactor.DataStore.StoreRepo().Delete(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to delete store: %w", err)
		return
	}

	return
}

func (interactor *StoreInteractor) GetStore(ctx context.Context, id uuid.UUID) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.GetStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	domainStore, err = interactor.DataStore.StoreRepo().FindByID(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to get store: %w", err)
		return
	}

	return
}

func (interactor *StoreInteractor) GetStoreByMerchantID(ctx context.Context, merchantID uuid.UUID) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.GetStoreByMerchantID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	domainStore, err = interactor.DataStore.StoreRepo().FindStoreMerchant(ctx, merchantID)
	if err != nil {
		err = fmt.Errorf("failed to get store by merchant id: %w", err)
		return
	}
	return
}

func (interactor *StoreInteractor) GetStores(ctx context.Context, pager *upagination.Pagination, filter *domain.StoreListFilter, orderBys ...domain.StoreListOrderBy) (domainStores []*domain.Store, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.GetStores")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	domainStores, total, err = interactor.DataStore.StoreRepo().GetStores(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get stores: %w", err)
		return
	}
	return
}

func (interactor *StoreInteractor) StoreSimpleUpdate(ctx context.Context, updateField domain.StoreSimpleUpdateType, domainUStoreParams *domain.UpdateStoreParams) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.StoreSimpleUpdate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if domainUStoreParams == nil {
		return domain.ParamsError(errors.New("domainUStoreParams is required"))
	}

	domainStore, err := interactor.DataStore.StoreRepo().FindByID(ctx, domainUStoreParams.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrStoreNotExists)
			return
		}
		return
	}

	switch updateField {
	case domain.StoreSimpleUpdateTypeStatus:
		if domainStore.Status == domainUStoreParams.Status {
			return
		}
		domainStore.Status = domainUStoreParams.Status
	default:
		return domain.ParamsError(fmt.Errorf("unsupported update field: %v", updateField))
	}

	err = interactor.DataStore.StoreRepo().Update(ctx, domainStore)
	if err != nil {
		err = fmt.Errorf("failed to update store: %w", err)
		return
	}

	return
}

func (interactor *StoreInteractor) CheckCreateStoreFields(ctx context.Context, domainCStore *domain.CreateStoreParams) (domainStore *domain.Store, err error) {
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
		BusinessTypeID:          domainCStore.BusinessTypeID,
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

	err = interactor.checkNameExists(ctx, domainStore)
	if err != nil {
		return
	}

	return
}

func (interactor *StoreInteractor) CheckUpdateStoreFields(ctx context.Context, domainUStore *domain.UpdateStoreParams) (domainStore *domain.Store, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreInteractor.CheckUpdateStoreFields")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	oldStore, err := interactor.DataStore.StoreRepo().FindByID(ctx, domainUStore.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrStoreNotExists)
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
		BusinessTypeID:          domainUStore.BusinessTypeID,
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
		LoginPassword:           domainUStore.LoginPassword,
	}

	if err = interactor.validateTimeConfigs(domainStore); err != nil {
		return nil, err
	}

	err = interactor.checkNameExists(ctx, domainStore)
	if err != nil {
		return
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
		if h.StartTime >= h.EndTime {
			return domain.ParamsError(domain.ErrStoreBusinessHoursTimeInvalid)
		}
		for _, wd := range h.Weekdays {
			perDay[wd] = append(perDay[wd], [2]string{h.StartTime, h.EndTime})
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

func (interactor *StoreInteractor) checkNameExists(ctx context.Context, domainStore *domain.Store) (err error) {
	exists, err := interactor.DataStore.StoreRepo().ExistsStore(ctx, &domain.ExistsStoreParams{
		MerchantID: domainStore.MerchantID,
		StoreName:  domainStore.StoreName,
		ExcludeID:  domainStore.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to check store name exists: %w", err)
	}
	if exists {
		return domain.ConflictError(domain.ErrStoreNameExists)
	}

	return
}

func NewStoreInteractor(dataStore domain.DataStore) *StoreInteractor {
	return &StoreInteractor{
		DataStore: dataStore,
	}
}
