package taxfee

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
var _ domain.TaxFeeInteractor = (*TaxFeeInteractor)(nil)

type TaxFeeInteractor struct {
	ds domain.DataStore
}

func NewTaxFeeInteractor(ds domain.DataStore) *TaxFeeInteractor {
	return &TaxFeeInteractor{ds: ds}
}

func (interactor *TaxFeeInteractor) Create(ctx context.Context, fee *domain.TaxFee) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()

	if fee == nil {
		return fmt.Errorf("tax fee is nil")
	}

	if err = interactor.checkExists(ctx, fee); err != nil {
		return err
	}

	fee.ID = uuid.New()
	if err = interactor.ds.TaxFeeRepo().Create(ctx, fee); err != nil {
		return fmt.Errorf("failed to create tax fee: %w", err)
	}
	return nil
}

func (interactor *TaxFeeInteractor) Update(ctx context.Context, fee *domain.TaxFee) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()

	if fee == nil {
		return fmt.Errorf("tax fee is nil")
	}

	oldFee, err := interactor.ds.TaxFeeRepo().FindByID(ctx, fee.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrTaxFeeNotExists)
		}
		return fmt.Errorf("failed to fetch tax fee: %w", err)
	}

	updatedFee := &domain.TaxFee{
		ID:          fee.ID,
		Name:        fee.Name,
		TaxFeeType:  oldFee.TaxFeeType,
		TaxCode:     oldFee.TaxCode,
		TaxRateType: fee.TaxRateType,
		TaxRate:     fee.TaxRate,
		DefaultTax:  fee.DefaultTax,
		MerchantID:  oldFee.MerchantID,
		StoreID:     oldFee.StoreID,
	}

	if err = interactor.checkExists(ctx, updatedFee); err != nil {
		return err
	}

	if err = interactor.ds.TaxFeeRepo().Update(ctx, updatedFee); err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrTaxFeeNotExists)
		}
		return fmt.Errorf("failed to update tax fee: %w", err)
	}
	return nil
}

func (interactor *TaxFeeInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.ds.TaxFeeRepo().Delete(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrTaxFeeNotExists)
		}
		return fmt.Errorf("failed to delete tax fee: %w", err)
	}
	return
}

func (interactor *TaxFeeInteractor) GetTaxFee(ctx context.Context, id uuid.UUID) (fee *domain.TaxFee, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeInteractor.GetTaxFee")
	defer func() { util.SpanErrFinish(span, err) }()

	fee, err = interactor.ds.TaxFeeRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrTaxFeeNotExists)
			return
		}
		err = fmt.Errorf("failed to fetch tax fee: %w", err)
		return
	}
	return
}

func (interactor *TaxFeeInteractor) GetTaxFees(ctx context.Context, pager *upagination.Pagination, filter *domain.TaxFeeListFilter, orderBys ...domain.TaxFeeOrderBy) (fees []*domain.TaxFee, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeInteractor.GetTaxFees")
	defer func() { util.SpanErrFinish(span, err) }()

	if pager == nil {
		err = domain.ParamsError(errors.New("pager is required"))
		return
	}
	if filter == nil {
		err = domain.ParamsError(errors.New("filter is required"))
		return
	}

	fees, total, err = interactor.ds.TaxFeeRepo().GetTaxFees(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get tax fees: %w", err)
		return
	}
	return
}

func (interactor *TaxFeeInteractor) TaxFeeSimpleUpdate(ctx context.Context, updateField domain.TaxFeeSimpleUpdateType, fee *domain.TaxFee) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeInteractor.TaxFeeSimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	if fee == nil {
		return fmt.Errorf("tax fee is nil")
	}

	oldFee, err := interactor.ds.TaxFeeRepo().FindByID(ctx, fee.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrTaxFeeNotExists)
		}
		return fmt.Errorf("failed to fetch tax fee: %w", err)
	}

	switch updateField {
	case domain.TaxFeeSimpleUpdateTypeDefault:
		if oldFee.DefaultTax == fee.DefaultTax {
			return nil
		}
		oldFee.DefaultTax = fee.DefaultTax
	default:
		return domain.ParamsError(errors.New("unsupported update field"))
	}

	if err = interactor.ds.TaxFeeRepo().Update(ctx, oldFee); err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrTaxFeeNotExists)
		}
		return fmt.Errorf("failed to update tax fee: %w", err)
	}
	return nil
}

func (interactor *TaxFeeInteractor) checkExists(ctx context.Context, fee *domain.TaxFee) (err error) {
	exists, err := interactor.ds.TaxFeeRepo().Exists(ctx, domain.TaxFeeExistsParams{
		MerchantID: fee.MerchantID,
		StoreID:    fee.StoreID,
		Name:       fee.Name,
		ExcludeID:  fee.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to check tax fee exists: %w", err)
	}
	if exists {
		return domain.ParamsError(domain.ErrTaxFeeNameExists)
	}
	return nil
}
