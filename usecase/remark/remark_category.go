package remark

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RemarkCategoryInteractor = (*RemarkCategoryInteractor)(nil)

type RemarkCategoryInteractor struct {
	DataStore domain.DataStore
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
	err = interactor.DataStore.RemarkCategoryRepo().Create(ctx, remarkCategory)
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
	err = interactor.DataStore.RemarkCategoryRepo().Update(ctx, remarkCategory)
	if err != nil {
		err = fmt.Errorf("failed to update remark category: %w", err)
		return
	}
	return
}

func (interactor *RemarkCategoryInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RemarkCategoryInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	_, err = interactor.DataStore.RemarkCategoryRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrRemarkCategoryNotExists)
		}
		err = fmt.Errorf("failed to fetch remark category: %w", err)
		return
	}
	err = interactor.DataStore.RemarkCategoryRepo().Delete(ctx, id)
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
	remarkCategories, err = interactor.DataStore.RemarkCategoryRepo().GetRemarkCategories(ctx, filter)
	if err != nil {
		err = fmt.Errorf("failed to list remark categories: %w", err)
		return
	}
	categoriesIds := lo.Map(remarkCategories, func(item *domain.RemarkCategory, _ int) uuid.UUID {
		return item.ID
	})
	countRemark, err := interactor.DataStore.RemarkRepo().CountRemarkByCategories(ctx, domain.CountRemarkParams{
		CategoryIDs: categoriesIds,
		MerchantID:  filter.MerchantID,
		StoreID:     filter.StoreID,
		RemarkType:  filter.CountScene,
	})
	if err != nil {
		err = fmt.Errorf("failed to count remarks by categories: %w", err)
		return
	}
	for _, category := range remarkCategories {
		if count, ok := countRemark[category.ID]; ok {
			category.RemarkCount = count
		} else {
			category.RemarkCount = 0
		}
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
	exists, err = interactor.DataStore.RemarkCategoryRepo().Exists(ctx, params)
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
		DataStore: dataStore,
	}
}
