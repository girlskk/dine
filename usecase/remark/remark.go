package remark

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RemarkInteractor = (*RemarkInteractor)(nil)

type RemarkInteractor struct {
	DataStore domain.DataStore
}

func (interactor *RemarkInteractor) Create(ctx context.Context, remark *domain.CreateRemarkParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()

	if remark == nil {
		return domain.ParamsError(errors.New("remark is required"))
	}

	domainRemark := &domain.Remark{
		Name:       remark.Name,
		Enabled:    remark.Enabled,
		SortOrder:  remark.SortOrder,
		CategoryID: remark.CategoryID,
		MerchantID: remark.MerchantID,
		StoreID:    remark.StoreID,
	}

	err = interactor.checkExists(ctx, domainRemark)
	if err != nil {
		return
	}
	domainRemark.ID = uuid.New()
	err = interactor.DataStore.RemarkRepo().Create(ctx, domainRemark)
	if err != nil {
		err = fmt.Errorf("failed to create remark: %w", err)
		return
	}
	return
}

func (interactor *RemarkInteractor) Update(ctx context.Context, remark *domain.UpdateRemarkParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()

	if remark == nil {
		return domain.ParamsError(errors.New("remark is required"))
	}

	oldRemark, err := interactor.DataStore.RemarkRepo().FindByID(ctx, remark.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrRemarkNotExists)
		}
		return fmt.Errorf("failed to fetch remark: %w", err)
	}
	domainRemark := &domain.Remark{
		ID:         remark.ID,
		Name:       remark.Name,
		Enabled:    remark.Enabled,
		SortOrder:  remark.SortOrder,
		RemarkType: oldRemark.RemarkType,
		CategoryID: oldRemark.CategoryID,
		MerchantID: oldRemark.MerchantID,
		StoreID:    oldRemark.StoreID,
	}
	err = interactor.checkExists(ctx, domainRemark)
	if err != nil {
		return
	}
	err = interactor.DataStore.RemarkRepo().Update(ctx, domainRemark)
	if err != nil {
		return fmt.Errorf("failed to update remark: %w", err)
	}
	return
}

func (interactor *RemarkInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	remark, err := interactor.DataStore.RemarkRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrRemarkNotExists)
		}
		err = fmt.Errorf("failed to fetch remark: %w", err)
		return
	}
	if remark.RemarkType == domain.RemarkTypeSystem {
		err = domain.ErrRemarkDeleteSystem
		return
	}
	err = interactor.DataStore.RemarkRepo().Delete(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to delete remark: %w", err)
		return
	}
	return
}

func (interactor *RemarkInteractor) GetRemark(ctx context.Context, id uuid.UUID) (remark *domain.Remark, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.GetRemark")
	defer func() { util.SpanErrFinish(span, err) }()

	remark, err = interactor.DataStore.RemarkRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrRemarkNotExists)
			return
		}
		err = fmt.Errorf("failed to get remark: %w", err)
		return
	}
	return
}

func (interactor *RemarkInteractor) GetRemarks(ctx context.Context, pager *upagination.Pagination, filter *domain.RemarkListFilter, orderBys ...domain.RemarkOrderBy) (remarks domain.Remarks, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.GetRemarks")
	defer func() { util.SpanErrFinish(span, err) }()
	if filter == nil {
		err = domain.ParamsError(errors.New("filter is required"))
	}
	remarks, total, err = interactor.DataStore.RemarkRepo().GetRemarks(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to list remarks: %w", err)
		return
	}
	return
}

func (interactor *RemarkInteractor) Exists(ctx context.Context, filter domain.RemarkExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.Exists")
	defer func() { util.SpanErrFinish(span, err) }()

	exists, err = interactor.DataStore.RemarkRepo().Exists(ctx, filter)
	if err != nil {
		err = fmt.Errorf("failed to check remark exists: %w", err)
	}
	return
}

func (interactor *RemarkInteractor) RemarkSimpleUpdate(ctx context.Context, updateField domain.RemarkSimpleUpdateType, remark *domain.Remark) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.RemarkSimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	if remark == nil {
		return domain.ParamsError(errors.New("remark is required"))
	}

	oldRemark, err := interactor.DataStore.RemarkRepo().FindByID(ctx, remark.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrRemarkNotExists)
		}
		err = fmt.Errorf("failed to fetch remark: %w", err)
		return
	}

	switch updateField {
	case domain.RemarkSimpleUpdateTypeEnabled:
		if oldRemark.Enabled == remark.Enabled {
			return nil
		}
		oldRemark.Enabled = remark.Enabled
	default:
		return domain.ParamsError(fmt.Errorf("unsupported update field: %v", updateField))
	}

	err = interactor.DataStore.RemarkRepo().Update(ctx, oldRemark)
	if err != nil {
		err = fmt.Errorf("failed to simple update remark: %w", err)
	}
	return
}

func (interactor *RemarkInteractor) checkExists(ctx context.Context, remark *domain.Remark) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.checkExists")
	defer func() { util.SpanErrFinish(span, err) }()

	filter := domain.RemarkExistsParams{
		RemarkType: remark.RemarkType,
		Name:       remark.Name,
		CategoryID: remark.CategoryID,
		MerchantID: remark.MerchantID,
		StoreID:    remark.StoreID,
		ExcludeID:  remark.ID,
	}
	exists, err := interactor.DataStore.RemarkRepo().Exists(ctx, filter)
	if err != nil {
		err = fmt.Errorf("failed to check remark exists: %w", err)
		return
	}
	if exists {
		err = domain.ConflictError(domain.ErrRemarkNameExists)
		return
	}
	return
}

func NewRemarkInteractor(ds domain.DataStore) *RemarkInteractor {
	return &RemarkInteractor{DataStore: ds}
}
