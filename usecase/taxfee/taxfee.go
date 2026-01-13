package taxfee

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

// Ensure implementation satisfies the interface.
var _ domain.TaxFeeInteractor = (*TaxFeeInteractor)(nil)

type TaxFeeInteractor struct {
	DS domain.DataStore
}

func NewTaxFeeInteractor(ds domain.DataStore) *TaxFeeInteractor {
	return &TaxFeeInteractor{DS: ds}
}

func (interactor *TaxFeeInteractor) Create(ctx context.Context, fee *domain.TaxFee, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()

	if err = verifyTaxFeeOwnership(user, fee); err != nil {
		return err
	}

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.TaxFeeRepo().Exists(ctx, domain.TaxFeeExistsParams{
			MerchantID: fee.MerchantID,
			StoreID:    fee.StoreID,
			Name:       fee.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrTaxFeeNameExists
		}
		fee.ID = uuid.New()
		err = interactor.DS.TaxFeeRepo().Create(ctx, fee)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (interactor *TaxFeeInteractor) Update(ctx context.Context, fee *domain.TaxFee, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		oldFee, err := ds.TaxFeeRepo().FindByID(ctx, fee.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrTaxFeeNotExists
			}
			return err
		}

		if err = verifyTaxFeeOwnership(user, oldFee); err != nil {
			return err
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

		exists, err := ds.TaxFeeRepo().Exists(ctx, domain.TaxFeeExistsParams{
			MerchantID: fee.MerchantID,
			StoreID:    fee.StoreID,
			Name:       fee.Name,
			ExcludeID:  fee.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrTaxFeeNameExists
		}
		err = ds.TaxFeeRepo().Update(ctx, updatedFee)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (interactor *TaxFeeInteractor) Delete(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	fee, err := interactor.DS.TaxFeeRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrTaxFeeNotExists)
		}
		return err
	}
	if err = verifyTaxFeeOwnership(user, fee); err != nil {
		return err
	}

	return interactor.DS.TaxFeeRepo().Delete(ctx, id)
}

func (interactor *TaxFeeInteractor) GetTaxFee(ctx context.Context, id uuid.UUID, user domain.User) (fee *domain.TaxFee, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeInteractor.GetTaxFee")
	defer func() { util.SpanErrFinish(span, err) }()

	fee, err = interactor.DS.TaxFeeRepo().FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err = verifyTaxFeeOwnership(user, fee); err != nil {
		return nil, err
	}
	return fee, nil
}

func (interactor *TaxFeeInteractor) GetTaxFees(ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.TaxFeeListFilter,
	orderBys ...domain.TaxFeeOrderBy,
) (fees []*domain.TaxFee, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeInteractor.GetTaxFees")
	defer func() { util.SpanErrFinish(span, err) }()

	return interactor.DS.TaxFeeRepo().GetTaxFees(ctx, pager, filter, orderBys...)
}

func (interactor *TaxFeeInteractor) TaxFeeSimpleUpdate(ctx context.Context,
	updateField domain.TaxFeeSimpleUpdateField,
	fee *domain.TaxFee,
	user domain.User,
) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "TaxFeeInteractor.TaxFeeSimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		oldFee, err := ds.TaxFeeRepo().FindByID(ctx, fee.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrTaxFeeNotExists
			}
			return err
		}

		if err = verifyTaxFeeOwnership(user, oldFee); err != nil {
			return err
		}

		switch updateField {
		case domain.TaxFeeSimpleUpdateFieldDefault:
			if oldFee.DefaultTax == fee.DefaultTax {
				return nil
			}
			oldFee.DefaultTax = fee.DefaultTax
		default:
			return fmt.Errorf("unsupported update field")
		}

		err = ds.TaxFeeRepo().Update(ctx, oldFee)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func verifyTaxFeeOwnership(user domain.User, fee *domain.TaxFee) error {
	switch user.GetUserType() {
	case domain.UserTypeAdmin:
	case domain.UserTypeBackend:
		if !domain.VerifyOwnerMerchant(user, fee.MerchantID) {
			return domain.ErrTaxFeeNotExists
		}
	case domain.UserTypeStore:
		if !domain.VerifyOwnerShip(user, fee.MerchantID, fee.StoreID) {
			return domain.ErrTaxFeeNotExists
		}
	}

	return nil
}
