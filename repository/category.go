package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/category"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.CategoryRepository = (*CategoryRepository)(nil)

type CategoryRepository struct {
	Client *ent.Client
}

func NewCategoryRepository(client *ent.Client) *CategoryRepository {
	return &CategoryRepository{
		Client: client,
	}
}

func (repo *CategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.Category, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CategoryRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	ec, err := repo.Client.Category.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrCategoryNotExists)
		}
		return nil, err
	}

	res = convertCategoryToDomain(ec)

	return res, nil
}

func (repo *CategoryRepository) Create(ctx context.Context, cat *domain.Category) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CategoryRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.Category.Create().
		SetID(cat.ID).
		SetName(cat.Name).
		SetMerchantID(cat.MerchantID).
		SetInheritTaxRate(cat.InheritTaxRate).
		SetInheritStall(cat.InheritStall).
		SetSortOrder(cat.SortOrder)

	if cat.StoreID != uuid.Nil {
		builder = builder.SetStoreID(cat.StoreID)
	}
	if cat.ParentID != uuid.Nil {
		builder = builder.SetParentID(cat.ParentID)
	}
	if cat.TaxRateID != uuid.Nil {
		builder = builder.SetTaxRateID(cat.TaxRateID)
	}
	if cat.StallID != uuid.Nil {
		builder = builder.SetStallID(cat.StallID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return err
	}

	cat.ID = created.ID
	cat.CreatedAt = created.CreatedAt

	return nil
}

func (repo *CategoryRepository) CreateBulk(ctx context.Context, categories []*domain.Category) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CategoryRepository.CreateBulk")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(categories) == 0 {
		return nil
	}

	builders := make([]*ent.CategoryCreate, 0, len(categories))
	for _, cat := range categories {
		builder := repo.Client.Category.Create().
			SetID(cat.ID).
			SetName(cat.Name).
			SetMerchantID(cat.MerchantID).
			SetInheritTaxRate(cat.InheritTaxRate).
			SetInheritStall(cat.InheritStall).
			SetSortOrder(cat.SortOrder)

		if cat.StoreID != uuid.Nil {
			builder = builder.SetStoreID(cat.StoreID)
		}
		if cat.ParentID != uuid.Nil {
			builder = builder.SetParentID(cat.ParentID)
		}
		if cat.TaxRateID != uuid.Nil {
			builder = builder.SetTaxRateID(cat.TaxRateID)
		}
		if cat.StallID != uuid.Nil {
			builder = builder.SetStallID(cat.StallID)
		}

		builders = append(builders, builder)
	}

	_, err = repo.Client.Category.CreateBulk(builders...).Save(ctx)
	return err
}

func (repo *CategoryRepository) Update(ctx context.Context, cat *domain.Category) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CategoryRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.Category.UpdateOneID(cat.ID).
		SetName(cat.Name).
		SetInheritTaxRate(cat.InheritTaxRate).
		SetInheritStall(cat.InheritStall).
		SetSortOrder(cat.SortOrder).
		SetProductCount(cat.ProductCount)

	if cat.ParentID == uuid.Nil {
		builder = builder.ClearParentID()
	} else {
		builder = builder.SetParentID(cat.ParentID)
	}

	if cat.TaxRateID == uuid.Nil {
		builder = builder.ClearTaxRateID()
	} else {
		builder = builder.SetTaxRateID(cat.TaxRateID)
	}

	if cat.StallID == uuid.Nil {
		builder = builder.ClearStallID()
	} else {
		builder = builder.SetStallID(cat.StallID)
	}

	updated, err := builder.Save(ctx)

	if err != nil {
		return err
	}

	cat.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CategoryRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.Category.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *CategoryRepository) Exists(ctx context.Context, params domain.CategoryExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CategoryRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Category.Query()
	if params.MerchantID != uuid.Nil {
		query.Where(category.MerchantID(params.MerchantID))
	}
	if params.IsRoot {
		query.Where(category.ParentIDIsNil())
	} else if params.ParentID != uuid.Nil {
		query.Where(category.ParentID(params.ParentID))
	}
	if params.Name != "" {
		query.Where(category.Name(params.Name))
	}
	return query.Exist(ctx)
}

func convertCategoryToDomain(ec *ent.Category) *domain.Category {
	if ec == nil {
		return nil
	}

	cat := &domain.Category{
		ID:             ec.ID,
		Name:           ec.Name,
		MerchantID:     ec.MerchantID,
		StoreID:        ec.StoreID,
		ParentID:       ec.ParentID,
		InheritTaxRate: ec.InheritTaxRate,
		TaxRateID:      ec.TaxRateID,
		InheritStall:   ec.InheritStall,
		StallID:        ec.StallID,
		ProductCount:   ec.ProductCount,
		SortOrder:      ec.SortOrder,
		CreatedAt:      ec.CreatedAt,
		UpdatedAt:      ec.UpdatedAt,
	}

	return cat
}
