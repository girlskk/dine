package repository

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/dinetable"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.TableRepository = (*TableRepository)(nil)

type TableRepository struct {
	Client *ent.Client
}

func NewTableRepository(client *ent.Client) *TableRepository {
	return &TableRepository{
		Client: client,
	}
}

func (repo *TableRepository) Create(ctx context.Context, table *domain.Table) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	created, err := repo.Client.DineTable.Create().
		SetAreaID(table.AreaID).
		SetName(table.Name).
		SetStoreID(table.StoreID).
		SetSeatCount(table.SeatCount).
		Save(ctx)

	if err != nil {
		return err
	}

	table.ID = created.ID
	table.CreatedAt = created.CreatedAt
	return nil
}

func (repo *TableRepository) Exists(ctx context.Context, params domain.TableExistsParams) (exists bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := repo.Client.DineTable.Query()
	if params.StoreID != 0 {
		query.Where(dinetable.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(dinetable.Name(params.Name))
	}
	return query.Exist(ctx)
}

func (repo *TableRepository) FindByID(ctx context.Context, id int) (res *domain.Table, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	table, err := repo.Client.DineTable.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrTableNotExists)
		}
		return nil, err
	}
	return convertTable(table), nil
}

func convertTable(table *ent.DineTable) *domain.Table {
	if table == nil {
		return nil
	}
	domainTable := &domain.Table{
		ID:        table.ID,
		Name:      table.Name,
		SeatCount: table.SeatCount,
		Status:    domain.TableStatus(table.Status),
		StoreID:   table.StoreID,
		AreaID:    table.AreaID,
		CreatedAt: table.CreatedAt,
		UpdatedAt: table.UpdatedAt,
		Order:     convertOrder(table.Edges.Order),
	}
	if table.Edges.Tablearea != nil {
		domainTable.Area = &domain.TableArea{
			ID:         table.Edges.Tablearea.ID,
			Name:       table.Edges.Tablearea.Name,
			StoreID:    table.Edges.Tablearea.StoreID,
			TableCount: table.Edges.Tablearea.TableCount,
		}
	}

	return domainTable
}

func (repo *TableRepository) Update(ctx context.Context, table *domain.Table) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 使用UpdateOneID进行精确更新
	updated, err := repo.Client.DineTable.UpdateOneID(table.ID).
		SetName(table.Name).
		SetSeatCount(table.SeatCount).
		SetAreaID(table.AreaID).
		Save(ctx)

	if err != nil {
		return err
	}
	// 更新领域对象时间戳
	table.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *TableRepository) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.DineTable.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TableRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.TableSearchParams,
) (res *domain.TableSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 构建基础查询
	query := repo.Client.DineTable.Query()

	// 应用搜索条件
	if params.StoreID > 0 {
		query.Where(dinetable.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(dinetable.NameContains(params.Name))
	}
	if params.AreaID > 0 {
		query.Where(dinetable.AreaID(params.AreaID))
	}
	if params.Status > 0 {
		query.Where(dinetable.Status(int(params.Status)))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	// 应用排序分页
	query.Order(ent.Desc(dinetable.FieldID)).
		WithTablearea().
		WithOrder().
		Offset(page.Offset()).
		Limit(page.Size)

	// 执行查询
	entTables, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为领域对象
	tables := make(domain.Tables, 0, len(entTables))
	for _, t := range entTables {
		tables = append(tables, convertTable(t))
	}

	return &domain.TableSearchRes{
		Pagination: page,
		Items:      tables,
	}, nil
}

func (repo *TableRepository) UpdateStatus(ctx context.Context, id int, status domain.TableStatus) (ok bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableRepository.UpdateStatus")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 使用Update进行精确更新
	affected, err := repo.Client.DineTable.Update().
		Where(dinetable.ID(id)).
		SetStatus(int(status)).
		Save(ctx)

	if err != nil {
		return
	}
	ok = affected > 0

	return
}

// 从状态到状态
func (repo *TableRepository) UpdateOrderIDAndStatusFrom(ctx context.Context, id int, orderID int, from, to domain.TableStatus) (ok bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableRepository.UpdateOrderIDAndStatusFrom")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 使用Update进行精确更新
	upd := repo.Client.DineTable.Update().
		Where(
			dinetable.ID(id),
			dinetable.Status(int(from)),
		).
		SetStatus(int(to))

	if orderID == 0 {
		upd.ClearOrderID()
	} else {
		upd.SetOrderID(orderID)
	}

	affected, err := upd.Save(ctx)
	if err != nil {
		return
	}
	ok = affected > 0

	return
}

func (repo *TableRepository) FindWithOrder(ctx context.Context, id int) (dtable *domain.Table, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TableRepository.FindWithOrder")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	table, err := repo.Client.DineTable.Query().
		Where(dinetable.ID(id)).
		WithOrder().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to get table by id: %w", err)
		return
	}

	dtable = convertTable(table)

	return
}
