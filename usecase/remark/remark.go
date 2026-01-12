package remark

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/i18n"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RemarkInteractor = (*RemarkInteractor)(nil)

type RemarkInteractor struct {
	DS domain.DataStore
}

func (interactor *RemarkInteractor) Create(ctx context.Context, remark *domain.CreateRemarkParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()

	domainRemark := &domain.Remark{
		Name:        remark.Name,
		RemarkType:  remark.RemarkType,
		Enabled:     remark.Enabled,
		SortOrder:   remark.SortOrder,
		RemarkScene: remark.RemarkScene,
		MerchantID:  remark.MerchantID,
		StoreID:     remark.StoreID,
	}

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.RemarkRepo().Exists(ctx, domain.RemarkExistsParams{
			RemarkType:  domainRemark.RemarkType,
			RemarkScene: domainRemark.RemarkScene,
			MerchantID:  domainRemark.MerchantID,
			StoreID:     domainRemark.StoreID,
			Name:        domainRemark.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrRemarkNameExists
		}

		domainRemark.ID = uuid.New()
		err = ds.RemarkRepo().Create(ctx, domainRemark)
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

func (interactor *RemarkInteractor) Update(ctx context.Context, remark *domain.UpdateRemarkParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		oldRemark, err := ds.RemarkRepo().FindByID(ctx, remark.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrRemarkNotExists
			}
			return err
		}
		domainRemark := &domain.Remark{
			ID:          remark.ID,
			Name:        remark.Name,
			Enabled:     remark.Enabled,
			SortOrder:   remark.SortOrder,
			RemarkType:  oldRemark.RemarkType,
			RemarkScene: oldRemark.RemarkScene,
			MerchantID:  oldRemark.MerchantID,
			StoreID:     oldRemark.StoreID,
		}

		exists, err := ds.RemarkRepo().Exists(ctx, domain.RemarkExistsParams{
			RemarkType:  domainRemark.RemarkType,
			RemarkScene: domainRemark.RemarkScene,
			MerchantID:  domainRemark.MerchantID,
			StoreID:     domainRemark.StoreID,
			Name:        domainRemark.Name,
			ExcludeID:   domainRemark.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrRemarkNameExists
		}

		err = ds.RemarkRepo().Update(ctx, domainRemark)
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

func (interactor *RemarkInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	remark, err := interactor.DS.RemarkRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ErrRemarkNotExists
		}
		return err
	}
	if remark.RemarkType == domain.RemarkTypeSystem {
		err = domain.ErrRemarkDeleteSystem
		return
	}
	err = interactor.DS.RemarkRepo().Delete(ctx, id)
	if err != nil {
		return err
	}
	return
}

func (interactor *RemarkInteractor) GetRemark(ctx context.Context, id uuid.UUID) (remark *domain.Remark, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.GetRemark")
	defer func() { util.SpanErrFinish(span, err) }()
	remark, err = interactor.DS.RemarkRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ErrRemarkNotExists
		}
		return
	}
	if msgID, ok := domain.RemarkSceneI18NMap[string(remark.RemarkScene)]; ok {
		name := i18n.Translate(ctx, msgID, nil)
		remark.RemarkSceneName = name
	}
	return
}

func (interactor *RemarkInteractor) GetRemarks(ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.RemarkListFilter,
	orderBys ...domain.RemarkOrderBy,
) (remarks domain.Remarks, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.GetRemarks")
	defer func() { util.SpanErrFinish(span, err) }()
	return interactor.DS.RemarkRepo().GetRemarks(ctx, pager, filter, orderBys...)
}

func (interactor *RemarkInteractor) RemarkSimpleUpdate(ctx context.Context,
	updateField domain.RemarkSimpleUpdateField,
	remark *domain.Remark,
) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkInteractor.RemarkSimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		oldRemark, err := ds.RemarkRepo().FindByID(ctx, remark.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrRemarkNotExists
			}
			return err
		}

		switch updateField {
		case domain.RemarkSimpleUpdateFieldEnabled:
			if oldRemark.Enabled == remark.Enabled {
				return nil
			}
			oldRemark.Enabled = remark.Enabled
		default:
			return fmt.Errorf("unsupported update field: %v", updateField)
		}

		err = ds.RemarkRepo().Update(ctx, oldRemark)
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

func NewRemarkInteractor(ds domain.DataStore) *RemarkInteractor {
	return &RemarkInteractor{DS: ds}
}
