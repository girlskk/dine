package repository

import (
	"context"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
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
		SetStoreID(cat.StoreID).
		SetParentID(cat.ParentID).
		SetInheritTaxRate(cat.InheritTaxRate).
		SetInheritStall(cat.InheritStall)

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
			SetStoreID(cat.StoreID).
			SetParentID(cat.ParentID).
			SetInheritTaxRate(cat.InheritTaxRate).
			SetInheritStall(cat.InheritStall)

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
		SetParentID(cat.ParentID).
		SetProductCount(cat.ProductCount)

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

	query := repo.Client.Category.Query().
		Where(category.MerchantID(params.MerchantID)).
		Where(category.StoreID(params.StoreID))

	if params.IsRoot {
		query.Where(category.ParentID(uuid.Nil))
	} else {
		query.Where(category.ParentID(params.ParentID))
	}
	if params.Name != "" {
		query.Where(category.Name(params.Name))
	}
	// 排除指定的ID（用于更新时检查名称唯一性）
	if params.ExcludeID != uuid.Nil {
		query.Where(category.IDNEQ(params.ExcludeID))
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

func (repo *CategoryRepository) ListBySearch(
	ctx context.Context,
	params domain.CategorySearchParams,
) (res domain.Categories, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CategoryRepository.ListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Category.Query()

	// 默认只查询一级分类
	query.Where(category.ParentID(uuid.Nil))

	if params.MerchantID != uuid.Nil {
		query.Where(category.MerchantID(params.MerchantID))
	}

	if params.OnlyMerchant {
		query.Where(category.StoreID(uuid.Nil))
	} else if params.StoreID != uuid.Nil {
		query.Where(category.StoreID(params.StoreID))
	}

	// 预加载子分类
	query.WithChildren(func(q *ent.CategoryQuery) {
		q.Order(
			category.BySortOrder(),            // 先按 SortOrder 升序（值越小越靠前）
			ent.Desc(category.FieldCreatedAt), // 如果 SortOrder 相同，按创建时间倒序
		)
	})

	// 预加载税率和出品部门
	query.WithTaxRate().WithStall()

	entCats, err := query.Order(
		category.BySortOrder(),
		ent.Desc(category.FieldCreatedAt),
	).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.Categories, 0, len(entCats))
	for _, c := range entCats {
		items = append(items, convertCategoryToDomainWithChildren(c))
	}

	return items, nil
}

func (repo *CategoryRepository) FindByNameInStore(ctx context.Context, name string, storeID, parentID uuid.UUID) (res *domain.Category, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CategoryRepository.FindByNameInStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := repo.Client.Category.Query().
		Where(category.Name(name)).
		Where(category.StoreID(storeID))
	if parentID != uuid.Nil {
		query.Where(category.ParentID(parentID))
	}
	ec, err := query.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrCategoryNotExists)
		}
		return nil, err
	}
	return convertCategoryToDomain(ec), nil
}

// 在 repository/category.go 实现：
func (repo *CategoryRepository) ListByParentID(ctx context.Context, merchantID, storeID, parentID uuid.UUID) (res domain.Categories, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CategoryRepository.ListByParentID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Category.Query().
		Where(category.MerchantID(merchantID)).
		Where(category.ParentID(parentID)).
		Where(category.StoreID(storeID))

	entCats, err := query.Order(category.BySortOrder()).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.Categories, 0, len(entCats))
	for _, c := range entCats {
		items = append(items, convertCategoryToDomain(c))
	}

	return items, nil
}

func (repo *CategoryRepository) UpdateSortOrders(ctx context.Context, updates map[uuid.UUID]int) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "CategoryRepository.UpdateSortOrders")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(updates) == 0 {
		return nil
	}

	// 收集所有需要更新的 ID
	ids := make([]uuid.UUID, 0, len(updates))
	for id := range updates {
		ids = append(ids, id)
	}

	// 使用 Ent 的 Modify 功能构建 CASE WHEN 批量更新
	_, err = repo.Client.Category.Update().
		Where(category.IDIn(ids...)).
		Modify(func(u *sql.UpdateBuilder) {
			var args []interface{}
			var caseExpr strings.Builder

			caseExpr.WriteString("CASE `id`")
			// 为每个 ID 添加 WHEN 分支
			for _, id := range ids {
				caseExpr.WriteString(" WHEN ? THEN ?")
				args = append(args, id, updates[id])
			}
			caseExpr.WriteString(" ELSE `sort_order` END")
			u.Set(category.FieldSortOrder, sql.Expr(caseExpr.String(), args...))
		}).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to batch update sort_order: %w", err)
	}

	return nil
}

// ============================================
// 转换函数
// ============================================

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

	if ec.Edges.Parent != nil {
		cat.Parent = convertCategoryToDomain(ec.Edges.Parent)
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
	// 转换税率和出品部门
	if ec.Edges.TaxRate != nil {
		cat.TaxRate = convertTaxFeeToDomain(ec.Edges.TaxRate)
	}
	if ec.Edges.Stall != nil {
		cat.Stall = convertStallToDomain(ec.Edges.Stall)
	}
	return cat
}
