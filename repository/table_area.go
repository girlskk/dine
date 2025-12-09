package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/tablearea"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.TableAreaRepository = (*TableAreaRepository)(nil)

type TableAreaRepository struct {
	Client *ent.Client
}

func NewTableAreaRepository(client *ent.Client) *TableAreaRepository {
	return &TableAreaRepository{
		Client: client,
	}
}

func (repo *TableAreaRepository) FindByID(ctx context.Context, id int) (res *domain.TableArea, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	area, err := repo.Client.TableArea.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrAreaNotExists)
		}
		return nil, err
	}
	return repo.convertToDomain(area), nil
}

func (repo *TableAreaRepository) convertToDomain(area *ent.TableArea) *domain.TableArea {
	if area == nil {
		return nil
	}
	return &domain.TableArea{
		ID:         area.ID,
		Name:       area.Name,
		TableCount: area.TableCount,
		StoreID:    area.StoreID,
		CreatedAt:  area.CreatedAt,
		UpdatedAt:  area.UpdatedAt,
	}
}

func (repo *TableAreaRepository) Exists(ctx context.Context, params domain.AreaExistsParams) (exists bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := repo.Client.TableArea.Query()
	if params.StoreID != 0 {
		query.Where(tablearea.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(tablearea.Name(params.Name))
	}
	return query.Exist(ctx)
}

func (repo *TableAreaRepository) Create(ctx context.Context, area *domain.TableArea) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	created, err := repo.Client.TableArea.Create().
		SetName(area.Name).
		SetStoreID(area.StoreID).
		Save(ctx)

	if err != nil {
		return err
	}

	area.ID = created.ID
	area.CreatedAt = created.CreatedAt
	return nil
}

func (repo *TableAreaRepository) Update(ctx context.Context, area *domain.TableArea) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	updated, err := repo.Client.TableArea.UpdateOneID(area.ID).
		SetName(area.Name).
		Save(ctx)

	if err != nil {
		return err
	}
	// 更新领域对象时间戳
	area.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *TableAreaRepository) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.TableArea.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TableAreaRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.AreaSearchParams,
) (res *domain.AreaSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 构建基础查询
	query := repo.Client.TableArea.Query()

	// 应用搜索条件
	if params.StoreID != 0 {
		query.Where(tablearea.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(tablearea.NameContains(params.Name))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	// 应用排序分页
	query.Order(ent.Desc(tablearea.FieldID)).
		Offset(page.Offset()).
		Limit(page.Size)

	// 执行查询
	entAreas, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为领域对象
	areas := make(domain.TableAreas, 0, len(entAreas))
	for _, u := range entAreas {
		areas = append(areas, repo.convertToDomain(u))
	}

	return &domain.AreaSearchRes{
		Pagination: page,
		Items:      areas,
	}, nil
}

func (repo *TableAreaRepository) IncreaseTableCount(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaRepository.IncreaseTableCount")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.TableArea.UpdateOneID(id).
		AddTableCount(1).
		Save(ctx)

	return err
}

func (repo *TableAreaRepository) DecreaseTableCount(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaRepository.DecreaseTableCount")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.TableArea.UpdateOneID(id).
		AddTableCount(-1).
		Save(ctx)

	return err
}
