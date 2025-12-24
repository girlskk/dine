package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/productattr"
	"gitlab.jiguang.dev/pos-dine/dine/ent/productattritem"
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

// ============================================
// ProductAttr 相关操作
// ============================================

func (repo *ProductAttrRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.ProductAttr, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	ea, err := repo.Client.ProductAttr.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProductAttrNotExists)
		}
		return nil, err
	}

	res = convertProductAttrToDomain(ea)

	return res, nil
}

func (repo *ProductAttrRepository) GetDetail(ctx context.Context, id uuid.UUID) (res *domain.ProductAttr, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.GetDetail")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	ea, err := repo.Client.ProductAttr.Query().
		Where(productattr.ID(id)).
		WithItems(func(q *ent.ProductAttrItemQuery) {
			q.Order(ent.Asc(productattritem.FieldCreatedAt))
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProductAttrNotExists)
		}
		return nil, err
	}

	res = convertProductAttrToDomain(ea)

	return res, nil
}

func (repo *ProductAttrRepository) Create(ctx context.Context, attr *domain.ProductAttr) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.ProductAttr.Create().
		SetID(attr.ID).
		SetName(attr.Name).
		SetChannels(attr.Channels).
		SetMerchantID(attr.MerchantID).
		SetProductCount(attr.ProductCount)

	if attr.StoreID != uuid.Nil {
		builder = builder.SetStoreID(attr.StoreID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return err
	}

	attr.ID = created.ID
	attr.CreatedAt = created.CreatedAt

	return nil
}

func (repo *ProductAttrRepository) Update(ctx context.Context, attr *domain.ProductAttr) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.ProductAttr.UpdateOneID(attr.ID).
		SetName(attr.Name).
		SetChannels(attr.Channels).
		SetProductCount(attr.ProductCount)

	updated, err := builder.Save(ctx)
	if err != nil {
		return err
	}

	attr.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *ProductAttrRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return repo.Client.ProductAttr.DeleteOneID(id).Exec(ctx)
}

func (repo *ProductAttrRepository) Exists(ctx context.Context, params domain.ProductAttrExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductAttr.Query()
	if params.MerchantID != uuid.Nil {
		query.Where(productattr.MerchantID(params.MerchantID))
	}
	if params.Name != "" {
		query.Where(productattr.Name(params.Name))
	}
	// 排除指定的ID（用于更新时检查名称唯一性）
	if params.ExcludeID != uuid.Nil {
		query.Where(productattr.IDNEQ(params.ExcludeID))
	}
	return query.Exist(ctx)
}

func (repo *ProductAttrRepository) ListBySearch(
	ctx context.Context,
	params domain.ProductAttrSearchParams,
) (res domain.ProductAttrs, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.ListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductAttr.Query()

	if params.MerchantID != uuid.Nil {
		query.Where(productattr.MerchantID(params.MerchantID))
	}

	// 预加载所有子项（用于避免 N+1 查询）
	query.WithItems(func(q *ent.ProductAttrItemQuery) {
		q.Order(ent.Asc(productattritem.FieldCreatedAt))
	})

	// 按创建时间倒序排列，加载所有匹配的记录
	entAttrs, err := query.Order(ent.Desc(productattr.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.ProductAttrs, 0, len(entAttrs))
	for _, a := range entAttrs {
		items = append(items, convertProductAttrToDomain(a))
	}

	return items, nil
}

// ============================================
// ProductAttrItem 相关操作
// ============================================

func (repo *ProductAttrRepository) FindItemByID(ctx context.Context, id uuid.UUID) (res *domain.ProductAttrItem, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.FindItemByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	ei, err := repo.Client.ProductAttrItem.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProductAttrItemNotExists)
		}
		return nil, err
	}

	res = convertProductAttrItemToDomain(ei)
	return res, nil
}

func (repo *ProductAttrRepository) CreateItems(ctx context.Context, items []*domain.ProductAttrItem) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.CreateItems")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(items) == 0 {
		return nil
	}

	builders := make([]*ent.ProductAttrItemCreate, 0, len(items))
	for _, item := range items {
		builder := repo.Client.ProductAttrItem.Create().
			SetID(item.ID).
			SetAttrID(item.AttrID).
			SetName(item.Name).
			SetBasePrice(item.BasePrice).
			SetProductCount(item.ProductCount).
			SetImage(item.Image)

		builders = append(builders, builder)
	}

	_, err = repo.Client.ProductAttrItem.CreateBulk(builders...).Save(ctx)
	return err
}

// SaveItems 批量保存口味做法项（新增或更新，如果ID存在则覆盖）
func (repo *ProductAttrRepository) SaveItems(ctx context.Context, items []*domain.ProductAttrItem) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.SaveItems")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(items) == 0 {
		return nil
	}

	// 构建所有的 Create builders
	builders := make([]*ent.ProductAttrItemCreate, 0, len(items))
	for _, item := range items {
		builder := repo.Client.ProductAttrItem.Create().
			SetID(item.ID).
			SetAttrID(item.AttrID).
			SetName(item.Name).
			SetBasePrice(item.BasePrice).
			SetProductCount(item.ProductCount).
			SetImage(item.Image)
		builders = append(builders, builder)
	}

	// 使用 Upsert：如果 ID 冲突则更新，否则创建
	err = repo.Client.ProductAttrItem.CreateBulk(builders...).
		OnConflictColumns(productattritem.FieldID).
		UpdateNewValues().
		Exec(ctx)

	return err
}

func (repo *ProductAttrRepository) DeleteItem(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.DeleteItem")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return repo.Client.ProductAttrItem.DeleteOneID(id).Exec(ctx)
}

func (repo *ProductAttrRepository) DeleteItems(ctx context.Context, ids []uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.DeleteItems")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	if len(ids) == 0 {
		return nil
	}
	_, err = repo.Client.ProductAttrItem.Delete().Where(productattritem.IDIn(ids...)).Exec(ctx)
	return err
}

func (repo *ProductAttrRepository) ListItemsByIDs(ctx context.Context, ids []uuid.UUID) (res domain.ProductAttrItems, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRepository.ListItemsByIDs")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductAttrItem.Query().Where(productattritem.IDIn(ids...))

	entItems, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	res = make(domain.ProductAttrItems, 0, len(entItems))
	for _, i := range entItems {
		res = append(res, convertProductAttrItemToDomain(i))
	}
	return res, nil
}

// ============================================
// 转换函数
// ============================================

func convertProductAttrToDomain(ea *ent.ProductAttr) *domain.ProductAttr {
	if ea == nil {
		return nil
	}

	attr := &domain.ProductAttr{
		ID:           ea.ID,
		Name:         ea.Name,
		Channels:     ea.Channels,
		MerchantID:   ea.MerchantID,
		StoreID:      ea.StoreID,
		ProductCount: ea.ProductCount,
		CreatedAt:    ea.CreatedAt,
		UpdatedAt:    ea.UpdatedAt,
	}

	// 加载关联的口味做法项
	if items, err := ea.Edges.ItemsOrErr(); err == nil && len(items) > 0 {
		attr.Items = make([]*domain.ProductAttrItem, 0, len(items))
		for _, item := range items {
			attr.Items = append(attr.Items, convertProductAttrItemToDomain(item))
		}
	} else {
		attr.Items = make([]*domain.ProductAttrItem, 0)
	}
	return attr
}

func convertProductAttrItemToDomain(ei *ent.ProductAttrItem) *domain.ProductAttrItem {
	if ei == nil {
		return nil
	}

	item := &domain.ProductAttrItem{
		ID:           ei.ID,
		AttrID:       ei.AttrID,
		Name:         ei.Name,
		Image:        ei.Image,
		BasePrice:    ei.BasePrice,
		ProductCount: ei.ProductCount,
		CreatedAt:    ei.CreatedAt,
		UpdatedAt:    ei.UpdatedAt,
	}

	return item
}
