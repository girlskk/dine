package repository

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/attr"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductAttrRepository = (*ProductAttrRepository)(nil)

type ProductAttrRepository struct {
	Client *ent.Client
}

func NewProductAttrRepository(client *ent.Client) *ProductAttrRepository {
	return &ProductAttrRepository{
		Client: client,
	}
}

func (repo *ProductAttrRepository) FindByID(ctx context.Context, id int) (res *domain.ProductAttr, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	entAttr, err := repo.Client.Attr.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrAttrNotExists)
		}
		return nil, err
	}
	return repo.convertToDomain(entAttr), nil
}

func (repo *ProductAttrRepository) Create(ctx context.Context, attr *domain.ProductAttr) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	created, err := repo.Client.Attr.Create().
		SetName(attr.Name).
		SetStoreID(attr.StoreID).
		Save(ctx)

	if err != nil {
		return err
	}

	attr.ID = created.ID
	attr.CreatedAt = created.CreatedAt
	return nil
}

func (repo *ProductAttrRepository) Exists(ctx context.Context, params domain.AttrExistsParams) (exists bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := repo.Client.Attr.Query()
	if params.StoreID != 0 {
		query.Where(attr.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(attr.Name(params.Name))
	}
	return query.Exist(ctx)
}

func (repo *ProductAttrRepository) Update(ctx context.Context, attr *domain.ProductAttr) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 使用UpdateOneID进行精确更新
	updated, err := repo.Client.Attr.UpdateOneID(attr.ID).
		SetName(attr.Name).
		Save(ctx)

	if err != nil {
		return err
	}
	// 更新领域对象时间戳
	attr.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *ProductAttrRepository) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.Attr.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductAttrRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.AttrSearchParams,
) (res *domain.AttrSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 构建基础查询
	query := repo.Client.Attr.Query()

	// 应用搜索条件
	if params.StoreID != 0 {
		query.Where(attr.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(attr.NameContains(params.Name))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	// 应用排序分页
	query.Order(ent.Desc(attr.FieldID)).
		Offset(page.Offset()).
		Limit(page.Size)

	// 执行查询
	entAttrs, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为领域对象
	attrs := make(domain.ProductAttrs, 0, len(entAttrs))
	for _, u := range entAttrs {
		attrs = append(attrs, repo.convertToDomain(u))
	}

	return &domain.AttrSearchRes{
		Pagination: page,
		Items:      attrs,
	}, nil
}

func (repo *ProductAttrRepository) ListByIDs(ctx context.Context, ids []int) (res domain.ProductAttrs, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrRepository.ListByIDs")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	attrs, err := repo.Client.Attr.Query().
		Where(attr.IDIn(ids...)).
		All(ctx)

	if err != nil {
		return nil, err
	}
	if len(attrs) == 0 {
		return nil, nil
	}
	for _, u := range attrs {
		res = append(res, repo.convertToDomain(u))
	}
	return res, nil
}

func (repo *ProductAttrRepository) IsUsedByProduct(ctx context.Context, id int) (used bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrRepository.IsUsedByProduct")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	used, err = repo.Client.Attr.Query().
		Where(attr.ID(id)).
		QueryProducts().
		Exist(ctx)
	if err != nil {
		return false, err
	}
	return used, nil
}

func (repo *ProductAttrRepository) convertToDomain(attr *ent.Attr) *domain.ProductAttr {
	if attr == nil {
		return nil
	}
	return &domain.ProductAttr{
		ID:        attr.ID,
		Name:      attr.Name,
		StoreID:   attr.StoreID,
		CreatedAt: attr.CreatedAt,
		UpdatedAt: attr.UpdatedAt,
	}
}
