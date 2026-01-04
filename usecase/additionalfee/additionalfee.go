package additionalfee

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

// Ensure implementation satisfies the interface.
var _ domain.AdditionalFeeInteractor = (*AdditionalFeeInteractor)(nil)

type AdditionalFeeInteractor struct {
	ds domain.DataStore
}

func NewAdditionalFeeInteractor(ds domain.DataStore) *AdditionalFeeInteractor {
	return &AdditionalFeeInteractor{ds: ds}
}

func (interactor *AdditionalFeeInteractor) Create(ctx context.Context, fee *domain.AdditionalFee) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()

	if fee == nil {
		return fmt.Errorf("additional fee is nil")
	}

	if err = interactor.checkExists(ctx, fee); err != nil {
		return err
	}

	fee.ID = uuid.New()
	err = interactor.ds.AdditionalFeeRepo().Create(ctx, fee)
	if err != nil {
		return fmt.Errorf("failed to create additional fee: %w", err)
	}

	return nil
}

func (interactor *AdditionalFeeInteractor) Update(ctx context.Context, fee *domain.AdditionalFee) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()

	if fee == nil {
		return fmt.Errorf("additional fee is nil")
	}

	oldFee, err := interactor.ds.AdditionalFeeRepo().FindByID(ctx, fee.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrAdditionalFeeNotExists)
		}
		return fmt.Errorf("failed to fetch additional fee: %w", err)
	}

	updatedFee := &domain.AdditionalFee{
		ID:                  fee.ID,
		Name:                fee.Name,
		FeeType:             oldFee.FeeType,
		FeeCategory:         fee.FeeCategory,
		ChargeMode:          fee.ChargeMode,
		FeeValue:            fee.FeeValue,
		IncludeInReceivable: fee.IncludeInReceivable,
		Taxable:             fee.Taxable,
		DiscountScope:       fee.DiscountScope,
		OrderChannels:       fee.OrderChannels,
		DiningWays:          fee.DiningWays,
		Enabled:             fee.Enabled,
		SortOrder:           fee.SortOrder,
		MerchantID:          oldFee.MerchantID,
		StoreID:             oldFee.StoreID,
	}

	if err = interactor.checkExists(ctx, updatedFee); err != nil {
		return err
	}

	if err = interactor.ds.AdditionalFeeRepo().Update(ctx, updatedFee); err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrAdditionalFeeNotExists)
		}
		return fmt.Errorf("failed to update additional fee: %w", err)
	}

	return nil
}

func (interactor *AdditionalFeeInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()
	err = interactor.ds.AdditionalFeeRepo().Delete(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrAdditionalFeeNotExists)
		}
		return fmt.Errorf("failed to delete additional fee: %w", err)
	}

	return
}

func (interactor *AdditionalFeeInteractor) GetAdditionalFee(ctx context.Context, id uuid.UUID) (fee *domain.AdditionalFee, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeInteractor.GetAdditionalFee")
	defer func() { util.SpanErrFinish(span, err) }()

	fee, err = interactor.ds.AdditionalFeeRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrAdditionalFeeNotExists)
			return
		}
		err = fmt.Errorf("failed to fetch additional fee: %w", err)
		return
	}

	return
}

func (interactor *AdditionalFeeInteractor) GetAdditionalFees(ctx context.Context, pager *upagination.Pagination, filter *domain.AdditionalFeeListFilter, orderBys ...domain.AdditionalFeeOrderBy) (fees []*domain.AdditionalFee, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeInteractor.GetAdditionalFees")
	defer func() { util.SpanErrFinish(span, err) }()

	if pager == nil {
		err = domain.ParamsError(errors.New("pager is required"))
		return
	}
	if filter == nil {
		err = domain.ParamsError(errors.New("filter is required"))
		return
	}

	fees, total, err = interactor.ds.AdditionalFeeRepo().GetAdditionalFees(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get additional fees: %w", err)
		return
	}

	return
}

func (interactor *AdditionalFeeInteractor) AdditionalFeeSimpleUpdate(ctx context.Context, updateField domain.AdditionalFeeSimpleUpdateType, fee *domain.AdditionalFee) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeInteractor.AdditionalFeeSimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	if fee == nil {
		return fmt.Errorf("additional fee is nil")
	}

	oldFee, err := interactor.ds.AdditionalFeeRepo().FindByID(ctx, fee.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrAdditionalFeeNotExists)
		}
		return fmt.Errorf("failed to fetch additional fee: %w", err)
	}

	switch updateField {
	case domain.AdditionalFeeSimpleUpdateTypeEnabled:
		if oldFee.Enabled == fee.Enabled {
			return
		}
		oldFee.Enabled = fee.Enabled
	default:
		return domain.ParamsError(errors.New("unsupported update field"))
	}

	err = interactor.ds.AdditionalFeeRepo().Update(ctx, oldFee)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrAdditionalFeeNotExists)
		}
		return fmt.Errorf("failed to update additional fee: %w", err)
	}

	return nil
}

func (interactor *AdditionalFeeInteractor) checkExists(ctx context.Context, fee *domain.AdditionalFee) (err error) {
	exists, err := interactor.ds.AdditionalFeeRepo().Exists(ctx, domain.AdditionalFeeExistsParams{
		MerchantID: fee.MerchantID,
		StoreID:    fee.StoreID,
		Name:       fee.Name,
		ExcludeID:  fee.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to check additional fee exists: %w", err)
	}
	if exists {
		return domain.ParamsError(domain.ErrRemarkNameExists)
	}
	return nil
}
