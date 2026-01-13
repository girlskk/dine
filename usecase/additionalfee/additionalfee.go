package additionalfee

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

// Ensure implementation satisfies the interface.
var _ domain.AdditionalFeeInteractor = (*AdditionalFeeInteractor)(nil)

type AdditionalFeeInteractor struct {
	DS domain.DataStore
}

func NewAdditionalFeeInteractor(ds domain.DataStore) *AdditionalFeeInteractor {
	return &AdditionalFeeInteractor{DS: ds}
}

func (interactor *AdditionalFeeInteractor) Create(ctx context.Context, fee *domain.AdditionalFee, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()

	if err = verifyAdditionalFeeOwnership(user, fee); err != nil {
		return err
	}

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.AdditionalFeeRepo().Exists(ctx, domain.AdditionalFeeExistsParams{
			FeeType:    fee.FeeType,
			MerchantID: fee.MerchantID,
			StoreID:    fee.StoreID,
			Name:       fee.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrRemarkNameExists
		}
		fee.ID = uuid.New()
		err = ds.AdditionalFeeRepo().Create(ctx, fee)
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

func (interactor *AdditionalFeeInteractor) Update(ctx context.Context, fee *domain.AdditionalFee, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		oldFee, err := ds.AdditionalFeeRepo().FindByID(ctx, fee.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrAdditionalFeeNotExists
			}
			return err
		}

		// ownership check based on old record
		if err = verifyAdditionalFeeOwnership(user, oldFee); err != nil {
			return err
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
		exists, err := ds.AdditionalFeeRepo().Exists(ctx, domain.AdditionalFeeExistsParams{
			FeeType:    fee.FeeType,
			MerchantID: fee.MerchantID,
			StoreID:    fee.StoreID,
			Name:       fee.Name,
			ExcludeID:  fee.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrRemarkNameExists
		}
		err = ds.AdditionalFeeRepo().Update(ctx, updatedFee)
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

func (interactor *AdditionalFeeInteractor) Delete(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	fee, err := interactor.DS.AdditionalFeeRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ErrAdditionalFeeNotExists
		}
		return err
	}
	if err = verifyAdditionalFeeOwnership(user, fee); err != nil {
		return err
	}

	return interactor.DS.AdditionalFeeRepo().Delete(ctx, id)
}

func (interactor *AdditionalFeeInteractor) GetAdditionalFee(ctx context.Context, id uuid.UUID, user domain.User) (fee *domain.AdditionalFee, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeInteractor.GetAdditionalFee")
	defer func() { util.SpanErrFinish(span, err) }()

	fee, err = interactor.DS.AdditionalFeeRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ErrAdditionalFeeNotExists
		}
		return nil, err
	}
	if err = verifyAdditionalFeeOwnership(user, fee); err != nil {
		return nil, err
	}
	return fee, nil
}

func (interactor *AdditionalFeeInteractor) GetAdditionalFees(ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.AdditionalFeeListFilter,
	orderBys ...domain.AdditionalFeeOrderBy,
) (fees []*domain.AdditionalFee, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeInteractor.GetAdditionalFees")
	defer func() { util.SpanErrFinish(span, err) }()

	return interactor.DS.AdditionalFeeRepo().GetAdditionalFees(ctx, pager, filter, orderBys...)
}

func (interactor *AdditionalFeeInteractor) AdditionalFeeSimpleUpdate(ctx context.Context,
	updateField domain.AdditionalFeeSimpleUpdateType,
	fee *domain.AdditionalFee,
	user domain.User,
) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdditionalFeeInteractor.AdditionalFeeSimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		oldFee, err := ds.AdditionalFeeRepo().FindByID(ctx, fee.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrAdditionalFeeNotExists
			}
			return err
		}

		if err = verifyAdditionalFeeOwnership(user, oldFee); err != nil {
			return err
		}

		switch updateField {
		case domain.AdditionalFeeSimpleUpdateTypeEnabled:
			if oldFee.Enabled == fee.Enabled {
				return nil
			}
			oldFee.Enabled = fee.Enabled
		default:
			return fmt.Errorf("unsupported update field")
		}

		err = ds.AdditionalFeeRepo().Update(ctx, oldFee)
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

func verifyAdditionalFeeOwnership(user domain.User, fee *domain.AdditionalFee) error {
	switch user.GetUserType() {
	case domain.UserTypeAdmin:
	case domain.UserTypeBackend:
		if !domain.VerifyOwnerMerchant(user, fee.MerchantID) {
			return domain.ErrAdditionalFeeNotExists
		}
	case domain.UserTypeStore:
		if !domain.VerifyOwnerShip(user, fee.MerchantID, fee.StoreID) {
			return domain.ErrAdditionalFeeNotExists
		}
	}
	return nil
}
