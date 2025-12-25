package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/remarkcategory"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RemarkCategoryRepository = (*RemarkCategoryRepository)(nil)

type RemarkCategoryRepository struct {
	Client *ent.Client
}

func NewRemarkCategoryRepository(client *ent.Client) *RemarkCategoryRepository {
	return &RemarkCategoryRepository{Client: client}
}

func (repo *RemarkCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (domainRemarkCategory *domain.RemarkCategory, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RemarkCategoryRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	ec, err := repo.Client.RemarkCategory.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrRemarkCategoryNotExists)
		}
		return nil, err
	}
	domainRemarkCategory = convertRemarkCategoryToDomain(ec)
	return
}

func (repo *RemarkCategoryRepository) Create(ctx context.Context, remarkCategory *domain.RemarkCategory) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RemarkCategoryRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	if remarkCategory == nil {
		return fmt.Errorf("remark category is nil")
	}

	_, err = repo.Client.RemarkCategory.Create().
		SetID(remarkCategory.ID).
		SetName(remarkCategory.Name).
		SetRemarkScene(remarkCategory.RemarkScene).
		SetSortOrder(remarkCategory.SortOrder).
		SetDescription(remarkCategory.Description).
		SetMerchantID(remarkCategory.MerchantID).
		Save(ctx)
	if err != nil {
		if err != nil {
			err = fmt.Errorf("failed to create remark category: %w", err)
			return err
		}
		return err
	}

	return
}

func (repo *RemarkCategoryRepository) Update(ctx context.Context, c *domain.RemarkCategory) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RemarkCategoryRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	if c == nil {
		return fmt.Errorf("remark category is nil")
	}

	_, err = repo.Client.RemarkCategory.UpdateOneID(c.ID).
		SetName(c.Name).
		SetRemarkScene(c.RemarkScene).
		SetSortOrder(c.SortOrder).
		SetDescription(c.Description).
		SetMerchantID(c.MerchantID).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to update remark category: %w", err)
		return err
	}

	return
}

func (repo *RemarkCategoryRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RemarkCategoryRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	err = repo.Client.RemarkCategory.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to delete remark category: %w", err)
		return
	}
	return
}

func (repo *RemarkCategoryRepository) GetRemarkCategories(ctx context.Context, filter *domain.RemarkCategoryListFilter) (domainRemarkCategories domain.RemarkCategories, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RemarkCategoryRepository.List")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if filter == nil {
		filter = &domain.RemarkCategoryListFilter{}
	}
	query := repo.Client.RemarkCategory.Query()
	if filter.MerchantID != uuid.Nil {
		query = query.Where(remarkcategory.Or(remarkcategory.MerchantID(filter.MerchantID), remarkcategory.MerchantIDIsNil(), remarkcategory.MerchantID(uuid.Nil)))
	}

	remarkCategories, err := query.Order(remarkcategory.BySortOrder()).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query remark category: %w", err)
		return nil, err
	}
	domainRemarkCategories = lo.Map(remarkCategories, func(item *ent.RemarkCategory, _ int) *domain.RemarkCategory {
		return convertRemarkCategoryToDomain(item)
	})
	return
}

func (repo *RemarkCategoryRepository) Exists(ctx context.Context, params domain.RemarkCategoryExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RemarkCategoryRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.RemarkCategory.Query()
	if params.MerchantID != uuid.Nil {
		query = query.Where(remarkcategory.MerchantID(params.MerchantID))
	}
	if params.Name != "" {
		query = query.Where(remarkcategory.Name(params.Name))
	}
	if params.ExcludeID != uuid.Nil {
		query = query.Where(remarkcategory.IDNEQ(params.ExcludeID))
	}
	exists, err = query.Exist(ctx)
	if err != nil {
		err = fmt.Errorf("failed to check remark category exists: %w", err)
		return
	}
	return
}

func convertRemarkCategoryToDomain(ec *ent.RemarkCategory) *domain.RemarkCategory {
	if ec == nil {
		return nil
	}
	return &domain.RemarkCategory{
		ID:          ec.ID,
		Name:        ec.Name,
		RemarkScene: ec.RemarkScene,
		MerchantID:  ec.MerchantID,
		Description: ec.Description,
		SortOrder:   ec.SortOrder,
		CreatedAt:   ec.CreatedAt,
		UpdatedAt:   ec.UpdatedAt,
	}
}
