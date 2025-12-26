package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/productunit"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductUnitRepository = (*ProductUnitRepository)(nil)

type ProductUnitRepository struct {
	Client *ent.Client
}

func NewProductUnitRepository(client *ent.Client) *ProductUnitRepository {
	return &ProductUnitRepository{
		Client: client,
	}
}

func (repo *ProductUnitRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.ProductUnit, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductUnitRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eu, err := repo.Client.ProductUnit.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProductUnitNotExists)
		}
		return nil, err
	}

	res = convertProductUnitToDomain(eu)

	return res, nil
}

func (repo *ProductUnitRepository) Create(ctx context.Context, unit *domain.ProductUnit) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductUnitRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.ProductUnit.Create().
		SetID(unit.ID).
		SetName(unit.Name).
		SetType(unit.Type).
		SetMerchantID(unit.MerchantID).
		SetProductCount(unit.ProductCount)

	if unit.StoreID != uuid.Nil {
		builder.SetStoreID(unit.StoreID)
	}
	created, err := builder.Save(ctx)
	if err != nil {
		return err
	}

	unit.ID = created.ID
	unit.CreatedAt = created.CreatedAt

	return nil
}

func (repo *ProductUnitRepository) Update(ctx context.Context, unit *domain.ProductUnit) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductUnitRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.ProductUnit.UpdateOneID(unit.ID).
		SetName(unit.Name).
		SetType(unit.Type).
		SetProductCount(unit.ProductCount)

	updated, err := builder.Save(ctx)

	if err != nil {
		return err
	}

	unit.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *ProductUnitRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductUnitRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.ProductUnit.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ProductUnitRepository) Exists(ctx context.Context, params domain.ProductUnitExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductUnitRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductUnit.Query()
	if params.MerchantID != uuid.Nil {
		query.Where(productunit.MerchantID(params.MerchantID))
	}
	if params.Name != "" {
		query.Where(productunit.Name(params.Name))
	}
	// 排除指定的ID（用于更新时检查名称唯一性）
	if params.ExcludeID != uuid.Nil {
		query.Where(productunit.IDNEQ(params.ExcludeID))
	}
	return query.Exist(ctx)
}

func (repo *ProductUnitRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.ProductUnitSearchParams,
) (res *domain.ProductUnitSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductUnitRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductUnit.Query()

	if params.MerchantID != uuid.Nil {
		query.Where(productunit.MerchantID(params.MerchantID))
	}

	if params.Name != "" {
		query.Where(productunit.NameContains(params.Name))
	}
	if params.Type != "" {
		query.Where(productunit.TypeEQ(params.Type))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	// 应用排序分页
	query.Order(ent.Desc(productunit.FieldCreatedAt)).
		Offset(page.Offset()).
		Limit(page.Size)

	// 执行查询
	entUnits, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.ProductUnits, 0, len(entUnits))
	for _, u := range entUnits {
		items = append(items, convertProductUnitToDomain(u))
	}

	return &domain.ProductUnitSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

func convertProductUnitToDomain(eu *ent.ProductUnit) *domain.ProductUnit {
	if eu == nil {
		return nil
	}

	unit := &domain.ProductUnit{
		ID:           eu.ID,
		Name:         eu.Name,
		Type:         eu.Type,
		MerchantID:   eu.MerchantID,
		StoreID:      eu.StoreID,
		ProductCount: eu.ProductCount,
		CreatedAt:    eu.CreatedAt,
		UpdatedAt:    eu.UpdatedAt,
	}

	return unit
}
