package repository

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/unit"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
	"time"
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

func (repo *ProductUnitRepository) FindByID(ctx context.Context, id int) (res *domain.ProductUnit, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	unit, err := repo.Client.Unit.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrUnitNotExists)
		}
		return nil, err
	}
	return repo.convertToDomain(unit), nil
}

func (repo *ProductUnitRepository) Create(ctx context.Context, unit *domain.ProductUnit) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 转换领域模型到Ent模型
	created, err := repo.Client.Unit.Create().
		SetName(unit.Name).
		SetStoreID(unit.StoreID).
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return err
	}
	// 回填生成字段
	unit.ID = created.ID
	unit.CreatedAt = created.CreatedAt
	return nil
}

func (repo *ProductUnitRepository) Update(ctx context.Context, unit *domain.ProductUnit) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 使用UpdateOneID进行精确更新
	updated, err := repo.Client.Unit.UpdateOneID(unit.ID).
		SetName(unit.Name).
		Save(ctx)

	if err != nil {
		return err
	}
	// 更新领域对象时间戳
	unit.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *ProductUnitRepository) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.Unit.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return
}

func (repo *ProductUnitRepository) Exists(ctx context.Context, params domain.UnitExistsParams) (exists bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := repo.Client.Unit.Query()
	if params.StoreID != 0 {
		query.Where(unit.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(unit.Name(params.Name))
	}
	return query.Exist(ctx)
}

func (repo *ProductUnitRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.UnitSearchParams,
) (res *domain.UnitSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 构建基础查询
	query := repo.Client.Unit.Query()

	// 应用搜索条件
	if params.StoreID != 0 {
		query.Where(unit.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(unit.NameContains(params.Name))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	// 应用排序分页
	query.Order(ent.Desc(unit.FieldID)).
		Offset(page.Offset()).
		Limit(page.Size)

	// 执行查询
	entUnits, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为领域对象
	units := make(domain.ProductUnits, 0, len(entUnits))
	for _, u := range entUnits {
		units = append(units, repo.convertToDomain(u))
	}

	return &domain.UnitSearchRes{
		Pagination: page,
		Items:      units,
	}, nil
}

func (repo *ProductUnitRepository) convertToDomain(unit *ent.Unit) *domain.ProductUnit {
	if unit == nil {
		return nil
	}
	return &domain.ProductUnit{
		ID:        unit.ID,
		Name:      unit.Name,
		StoreID:   unit.StoreID,
		CreatedAt: unit.CreatedAt,
		UpdatedAt: unit.UpdatedAt,
	}
}
