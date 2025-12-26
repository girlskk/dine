package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/producttag"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductTagRepository = (*ProductTagRepository)(nil)

type ProductTagRepository struct {
	Client *ent.Client
}

func NewProductTagRepository(client *ent.Client) *ProductTagRepository {
	return &ProductTagRepository{
		Client: client,
	}
}

func (repo *ProductTagRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.ProductTag, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductTagRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	et, err := repo.Client.ProductTag.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProductTagNotExists)
		}
		return nil, err
	}

	res = convertProductTagToDomain(et)

	return res, nil
}

func (repo *ProductTagRepository) Create(ctx context.Context, tag *domain.ProductTag) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductTagRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.ProductTag.Create().
		SetID(tag.ID).
		SetName(tag.Name).
		SetMerchantID(tag.MerchantID).
		SetProductCount(tag.ProductCount)

	if tag.StoreID != uuid.Nil {
		builder.SetStoreID(tag.StoreID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return err
	}

	tag.ID = created.ID
	tag.CreatedAt = created.CreatedAt

	return nil
}

func (repo *ProductTagRepository) Update(ctx context.Context, tag *domain.ProductTag) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductTagRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.ProductTag.UpdateOneID(tag.ID).
		SetName(tag.Name).
		SetProductCount(tag.ProductCount)

	updated, err := builder.Save(ctx)

	if err != nil {
		return err
	}

	tag.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *ProductTagRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductTagRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.ProductTag.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ProductTagRepository) Exists(ctx context.Context, params domain.ProductTagExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductTagRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductTag.Query()
	if params.MerchantID != uuid.Nil {
		query.Where(producttag.MerchantID(params.MerchantID))
	}
	if params.Name != "" {
		query.Where(producttag.Name(params.Name))
	}
	// 排除指定的ID（用于更新时检查名称唯一性）
	if params.ExcludeID != uuid.Nil {
		query.Where(producttag.IDNEQ(params.ExcludeID))
	}
	return query.Exist(ctx)
}

func (repo *ProductTagRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.ProductTagSearchParams,
) (res *domain.ProductTagSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductTagRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductTag.Query()

	if params.MerchantID != uuid.Nil {
		query.Where(producttag.MerchantID(params.MerchantID))
	}
	if params.Name != "" {
		query.Where(producttag.NameContains(params.Name))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}

	// 分页处理
	query = query.
		Offset(page.Offset()).
		Limit(page.Size)

	// 按创建时间倒序排列
	entTags, err := query.Order(ent.Desc(producttag.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.ProductTags, 0, len(entTags))
	for _, t := range entTags {
		items = append(items, convertProductTagToDomain(t))
	}

	page.SetTotal(total)

	return &domain.ProductTagSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

func (repo *ProductTagRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) (res domain.ProductTags, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductTagRepository.ListByIDs")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductTag.Query().Where(producttag.IDIn(ids...))

	entTags, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	res = make(domain.ProductTags, 0, len(entTags))
	for _, t := range entTags {
		res = append(res, convertProductTagToDomain(t))
	}
	return res, nil
}

func (repo *ProductTagRepository) FindByNamesInStore(ctx context.Context, storeID uuid.UUID, names []string) (res domain.ProductTags, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductTagRepository.FindByNamesInStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductTag.Query().Where(producttag.StoreID(storeID)).Where(producttag.NameIn(names...))
	entTags, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	res = make(domain.ProductTags, 0, len(entTags))
	for _, t := range entTags {
		res = append(res, convertProductTagToDomain(t))
	}
	return res, nil
}

func (repo *ProductTagRepository) CreateBulk(ctx context.Context, tags domain.ProductTags) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductTagRepository.CreateBulk")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	if len(tags) == 0 {
		return nil
	}

	builders := make([]*ent.ProductTagCreate, 0, len(tags))
	for _, tag := range tags {
		builder := repo.Client.ProductTag.Create().
			SetID(tag.ID).
			SetName(tag.Name).
			SetMerchantID(tag.MerchantID).
			SetProductCount(tag.ProductCount)

		if tag.StoreID != uuid.Nil {
			builder = builder.SetStoreID(tag.StoreID)
		}
		builders = append(builders, builder)
	}

	_, err = repo.Client.ProductTag.CreateBulk(builders...).Save(ctx)
	return err
}

func convertProductTagToDomain(et *ent.ProductTag) *domain.ProductTag {
	if et == nil {
		return nil
	}

	tag := &domain.ProductTag{
		ID:           et.ID,
		Name:         et.Name,
		MerchantID:   et.MerchantID,
		StoreID:      et.StoreID,
		ProductCount: et.ProductCount,
		CreatedAt:    et.CreatedAt,
		UpdatedAt:    et.UpdatedAt,
	}

	return tag
}
