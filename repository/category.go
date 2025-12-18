package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/category"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
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
		SetInheritStall(cat.InheritStall)

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
			SetInheritStall(cat.InheritStall)

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

func (repo *CategoryRepository) CountChildrenByParentID(ctx context.Context, parentID uuid.UUID) (count int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CategoryRepository.CountChildrenByParentID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	count, err = repo.Client.Category.Query().
		Where(category.ParentID(parentID)).
		Count(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *CategoryRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.CategorySearchParams,
) (res *domain.CategorySearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CategoryRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Category.Query()

	// 默认只查询一级分类
	query.Where(category.ParentIDIsNil())

	if params.MerchantID != uuid.Nil {
		query.Where(category.MerchantID(params.MerchantID))
	}
	if params.ID != uuid.Nil {
		query.Where(category.ID(params.ID))
	}
	if params.Name != "" {
		query.Where(category.NameContains(params.Name))
	}

	// 预加载子分类
	query.WithChildren(func(q *ent.CategoryQuery) {
		q.Order(
			category.BySortOrder(),            // 先按 SortOrder 升序（值越小越靠前）
			ent.Desc(category.FieldCreatedAt), // 如果 SortOrder 相同，按创建时间倒序
		)
	})

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	entCats, err := query.Order(
		category.BySortOrder(),
		ent.Desc(category.FieldCreatedAt),
	).
		Offset(page.Offset()).
		Limit(page.Size).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.Categories, 0, len(entCats))
	for _, c := range entCats {
		items = append(items, convertCategoryToDomainWithChildren(c))
	}

	return &domain.CategorySearchRes{
		Pagination: page,
		Items:      items,
	}, nil
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

func convertCategoryToDomainWithChildren(ec *ent.Category) *domain.Category {
	if ec == nil {
		return nil
	}

	cat := convertCategoryToDomain(ec)

	// 转换子分类
	if children, err := ec.Edges.ChildrenOrErr(); err == nil && len(children) > 0 {
		cat.Childrens = make([]*domain.Category, 0, len(children))
		for _, child := range children {
			cat.Childrens = append(cat.Childrens, convertCategoryToDomain(child))
		}
	}

	return cat
}
