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
	ds domain.DataStore
}

func (interactor *StallInteractor) StallSimpleUpdate(ctx context.Context, updateField domain.StallSimpleUpdateType, stall *domain.Stall) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StallInteractor.StallSimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	if stall == nil {
		return fmt.Errorf("stall is nil")
	}
	oldStall, err := interactor.ds.StallRepo().FindByID(ctx, stall.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrStallNotExists)
		}
		return fmt.Errorf("failed to fetch stall: %w", err)
	}

	switch updateField {
	case domain.StallSimpleUpdateTypeEnabled:
		if oldStall.Enabled == stall.Enabled {
			return
		}
		oldStall.Enabled = stall.Enabled
	default:
		return domain.ParamsError(errors.New("unsupported update field"))
	}

	err = interactor.ds.StallRepo().Update(ctx, oldStall)
	if err != nil {
		return fmt.Errorf("failed to update stall: %w", err)
	}
	return
}

func NewStallInteractor(ds domain.DataStore) *StallInteractor {
	return &StallInteractor{ds: ds}
}

func (interactor *StallInteractor) Create(ctx context.Context, domainStall *domain.Stall) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StallInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainStall == nil {
		return fmt.Errorf("stall is nil")
	}

	if err = interactor.checkExists(ctx, domainStall); err != nil {
		return err
	}
	domainStall.ID = uuid.New()
	err = interactor.ds.StallRepo().Create(ctx, domainStall)
	if err != nil {
		return fmt.Errorf("failed to create stall: %w", err)
	}
	return
}

func (interactor *StallInteractor) Update(ctx context.Context, domainStall *domain.Stall) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StallInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainStall == nil {
		return fmt.Errorf("stall is nil")
	}
	if err = interactor.checkExists(ctx, domainStall); err != nil {
		return err
	}
	err = interactor.ds.StallRepo().Update(ctx, domainStall)
	if err != nil {
		return fmt.Errorf("failed to update stall: %w", err)
	}
	return
}

func (interactor *StallInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StallInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()
	err = interactor.ds.StallRepo().Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete stall: %w", err)
	}
	return
}

func (interactor *StallInteractor) GetStall(ctx context.Context, id uuid.UUID) (domainStall *domain.Stall, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StallInteractor.GetStall")
	defer func() { util.SpanErrFinish(span, err) }()
	domainStall, err = interactor.ds.StallRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrStallNotExists)
			return
		}
		err = fmt.Errorf("failed to fetch stall: %w", err)
		return
	}
	return
}

func (interactor *StallInteractor) GetStalls(ctx context.Context, pager *upagination.Pagination, filter *domain.StallListFilter, orderBys ...domain.StallOrderBy) (domainStalls []*domain.Stall, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StallInteractor.GetStalls")
	defer func() { util.SpanErrFinish(span, err) }()
	if filter == nil {
		err = domain.ParamsError(errors.New("filter is required"))
	}
	domainStalls, total, err = interactor.ds.StallRepo().GetStalls(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get stalls: %w", err)
		return
	}
	return
}

func (interactor *StallInteractor) checkExists(ctx context.Context, domainStall *domain.Stall) (err error) {
	exists, err := interactor.ds.StallRepo().Exists(ctx, domain.StallExistsParams{
		MerchantID: domainStall.MerchantID,
		StoreID:    domainStall.StoreID,
		Name:       domainStall.Name,
		ExcludeID:  domainStall.ID,
	})
	if err != nil {
		return err
	}
	if exists {
		return domain.ParamsError(domain.ErrStallNameExists)
	}
	return nil
}
