package stall

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StallInteractor = (*StallInteractor)(nil)

type StallInteractor struct {
	DS domain.DataStore
}

func NewStallInteractor(ds domain.DataStore) *StallInteractor {
	return &StallInteractor{DS: ds}
}

func (interactor *StallInteractor) Create(ctx context.Context, domainStall *domain.Stall, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StallInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainStall == nil {
		return fmt.Errorf("stall is nil")
	}

	// verify ownership
	if err = verifyStallOwnership(user, domainStall); err != nil {
		return err
	}
	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.StallRepo().Exists(ctx, domain.StallExistsParams{
			MerchantID: domainStall.MerchantID,
			StoreID:    domainStall.StoreID,
			Name:       domainStall.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrStallNameExists
		}

		domainStall.ID = uuid.New()
		err = ds.StallRepo().Create(ctx, domainStall)
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

func (interactor *StallInteractor) Update(ctx context.Context, domainStall *domain.Stall, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StallInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainStall == nil {
		return fmt.Errorf("stall is nil")
	}
	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		old, err := ds.StallRepo().FindByID(ctx, domainStall.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrStallNotExists
			}
			return err
		}
		if err = verifyStallOwnership(user, old); err != nil {
			return err
		}
		if domainStall.StallType == domain.StallTypeSystem {
			return domain.ErrStallCannotUpdateSystem
		}
		domainStall.MerchantID = old.MerchantID
		domainStall.StoreID = old.StoreID
		domainStall.StallType = old.StallType

		exists, err := ds.StallRepo().Exists(ctx, domain.StallExistsParams{
			MerchantID: domainStall.MerchantID,
			StoreID:    domainStall.StoreID,
			Name:       domainStall.Name,
			ExcludeID:  domainStall.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrStallNameExists
		}
		err = ds.StallRepo().Update(ctx, domainStall)
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

func (interactor *StallInteractor) Delete(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StallInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()
	domainStall, err := interactor.DS.StallRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ErrStallNotExists
		}
		return err
	}
	if err = verifyStallOwnership(user, domainStall); err != nil {
		return err
	}
	if domainStall.StallType == domain.StallTypeSystem {
		return domain.ErrStallCannotDeleteSystem
	}
	err = interactor.DS.StallRepo().Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete stall: %w", err)
	}
	return
}

func (interactor *StallInteractor) GetStall(ctx context.Context, id uuid.UUID, user domain.User) (domainStall *domain.Stall, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StallInteractor.GetStall")
	defer func() { util.SpanErrFinish(span, err) }()
	domainStall, err = interactor.DS.StallRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ErrStallNotExists
			return
		}
		return
	}
	if err = verifyStallOwnership(user, domainStall); err != nil {
		return nil, err
	}
	return
}

func (interactor *StallInteractor) GetStalls(ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.StallListFilter,
	orderBys ...domain.StallOrderBy,
) (domainStalls []*domain.Stall, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StallInteractor.GetStalls")
	defer func() { util.SpanErrFinish(span, err) }()
	if filter == nil {
		err = domain.ParamsError(errors.New("filter is required"))
	}
	domainStalls, total, err = interactor.DS.StallRepo().GetStalls(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get stalls: %w", err)
		return
	}
	return
}

func (interactor *StallInteractor) StallSimpleUpdate(ctx context.Context,
	updateField domain.StallSimpleUpdateField,
	stall *domain.Stall,
	user domain.User,
) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StallInteractor.StallSimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		oldStall, err := ds.StallRepo().FindByID(ctx, stall.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrStallNotExists
			}
			return err
		}
		if err = verifyStallOwnership(user, oldStall); err != nil {
			return err
		}

		switch updateField {
		case domain.StallSimpleUpdateFieldEnabled:
			if oldStall.Enabled == stall.Enabled {
				return nil
			}
			oldStall.Enabled = stall.Enabled
		default:
			return domain.ParamsError(errors.New("unsupported update field"))
		}

		err = ds.StallRepo().Update(ctx, oldStall)
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

func verifyStallOwnership(user domain.User, stall *domain.Stall) error {
	switch user.GetUserType() {
	case domain.UserTypeAdmin:
	case domain.UserTypeBackend:
		if !domain.VerifyOwnerMerchant(user, stall.MerchantID) {
			return domain.ErrStallNotExists
		}
	case domain.UserTypeStore:
		if !domain.VerifyOwnerShip(user, stall.MerchantID, stall.StoreID) {
			return domain.ErrStallNotExists
		}
	}
	return nil
}
