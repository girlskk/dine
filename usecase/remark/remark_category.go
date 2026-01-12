package remark

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/i18n"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RemarkCategoryInteractor = (*RemarkCategoryInteractor)(nil)

type RemarkCategoryInteractor struct {
	DS domain.DataStore
}

func (interactor *RemarkCategoryInteractor) Create(ctx context.Context, remarkCategory *domain.RemarkCategory) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkCategoryInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()

	if remarkCategory == nil {
		return domain.ParamsError(errors.New("remark category is required"))
	}

	if err = interactor.checkExists(ctx, remarkCategory); err != nil {
		return
	}

	remarkCategory.ID = uuid.New()
	err = interactor.DS.RemarkCategoryRepo().Create(ctx, remarkCategory)
	if err != nil {
		err = fmt.Errorf("failed to create remark category: %w", err)
		return
	}
	return
}

func (interactor *RemarkCategoryInteractor) Update(ctx context.Context, remarkCategory *domain.RemarkCategory) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkCategoryInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()

	if remarkCategory == nil {
		return domain.ParamsError(errors.New("remark category is required"))
	}

	if err = interactor.checkExists(ctx, remarkCategory); err != nil {
		return
	}
	err = interactor.DS.RemarkCategoryRepo().Update(ctx, remarkCategory)
	if err != nil {
		err = fmt.Errorf("failed to update remark category: %w", err)
		return
	}
	return
}

func (interactor *RemarkCategoryInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkCategoryInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	_, err = interactor.DS.RemarkCategoryRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrRemarkCategoryNotExists)
		}
		err = fmt.Errorf("failed to fetch remark category: %w", err)
		return
	}
	err = interactor.DS.RemarkCategoryRepo().Delete(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to delete remark category: %w", err)
		return
	}
	return
}

func (interactor *RemarkCategoryInteractor) GetRemarkCategories(ctx context.Context,
	filter *domain.RemarkCategoryListFilter,
) (remarkCategories domain.RemarkCategories, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkCategoryInteractor.GetRemarkCategories")
	defer func() { util.SpanErrFinish(span, err) }()
	if filter == nil {
		filter = &domain.RemarkCategoryListFilter{}
	}
	remarkCategories, err = interactor.DS.RemarkCategoryRepo().GetRemarkCategories(ctx, filter)
	if err != nil {
		err = fmt.Errorf("failed to list remark categories: %w", err)
		return
	}
	return
}

func (interactor *RemarkCategoryInteractor) GetRemarkGroup(ctx context.Context, params domain.RemarkGroupListFilter) (remarkGroups []domain.RemarkGroup, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkCategoryInteractor.GetRemarkGroup")
	defer func() { util.SpanErrFinish(span, err) }()

	countRemark, err := interactor.DS.RemarkRepo().CountRemarkByScene(ctx, domain.CountRemarkParams{
		RemarkScenes: domain.RemarkSceneList,
		MerchantID:   params.MerchantID,
		StoreID:      params.StoreID,
		RemarkType:   params.CountScene,
	})
	if err != nil {
		err = fmt.Errorf("failed to count remarks by scenes: %w", err)
		return
	}
	// Build remark groups from the RemarkScene enum entries with i18n support
	for _, e := range domain.RemarkSceneEntries {
		name := i18n.Translate(ctx, e.MsgID, nil)
		remarkGroup := domain.RemarkGroup{
			Name:        name,
			RemarkScene: domain.RemarkScene(e.Code),
		}
		if count, ok := countRemark[remarkGroup.RemarkScene]; ok {
			remarkGroup.RemarkCount = count
		} else {
			remarkGroup.RemarkCount = 0
		}
		remarkGroups = append(remarkGroups, remarkGroup)
	}

	return
}

func (interactor *RemarkCategoryInteractor) checkExists(ctx context.Context, remarkCategory *domain.RemarkCategory) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkCategoryInteractor.checkExists")
	defer func() { util.SpanErrFinish(span, err) }()

	params := domain.RemarkCategoryExistsParams{
		MerchantID: remarkCategory.MerchantID,
		Name:       remarkCategory.Name,
		ExcludeID:  remarkCategory.ID,
	}
	var exists bool
	exists, err = interactor.DS.RemarkCategoryRepo().Exists(ctx, params)
	if err != nil {
		err = fmt.Errorf("failed to check remark category exists: %w", err)
		return
	}
	if exists {
		err = domain.ConflictError(domain.ErrRemarkCategoryNameExists)
		return
	}
	return
}

func NewRemarkCategoryInteractor(dataStore domain.DataStore) *RemarkCategoryInteractor {
	return &RemarkCategoryInteractor{
		DS: dataStore,
	}
}
