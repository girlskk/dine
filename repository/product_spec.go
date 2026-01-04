package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/productspec"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductSpecRepository = (*ProductSpecRepository)(nil)

type ProductSpecRepository struct {
	Client *ent.Client
}

func NewProductSpecRepository(client *ent.Client) *ProductSpecRepository {
	return &ProductSpecRepository{
		Client: client,
	}
}

func (repo *ProductSpecRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.ProductSpec, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductSpecRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	es, err := repo.Client.ProductSpec.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProductSpecNotExists)
		}
		return nil, err
	}

	res = convertProductSpecToDomain(es)

	return res, nil
}

func (repo *ProductSpecRepository) Create(ctx context.Context, spec *domain.ProductSpec) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductSpecRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.ProductSpec.Create().
		SetID(spec.ID).
		SetName(spec.Name).
		SetMerchantID(spec.MerchantID).
		SetStoreID(spec.StoreID).
		SetProductCount(spec.ProductCount)

	created, err := builder.Save(ctx)
	if err != nil {
		return err
	}

	spec.ID = created.ID
	spec.CreatedAt = created.CreatedAt

	return nil
}

func (repo *ProductSpecRepository) Update(ctx context.Context, spec *domain.ProductSpec) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductSpecRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.ProductSpec.UpdateOneID(spec.ID).
		SetName(spec.Name).
		SetProductCount(spec.ProductCount)

	updated, err := builder.Save(ctx)

	if err != nil {
		return err
	}

	spec.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *ProductSpecRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductSpecRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.ProductSpec.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ProductSpecRepository) Exists(ctx context.Context, params domain.ProductSpecExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductSpecRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductSpec.Query().
		Where(productspec.MerchantID(params.MerchantID)).
		Where(productspec.StoreID(params.StoreID))

	if params.Name != "" {
		query.Where(productspec.Name(params.Name))
	}
	// 排除指定的ID（用于更新时检查名称唯一性）
	if params.ExcludeID != uuid.Nil {
		query.Where(productspec.IDNEQ(params.ExcludeID))
	}
	return query.Exist(ctx)
}

func (repo *ProductSpecRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.ProductSpecSearchParams,
) (res *domain.ProductSpecSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductSpecRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductSpec.Query()

	if params.MerchantID != uuid.Nil {
		query.Where(productspec.MerchantID(params.MerchantID))
	}
	if params.Name != "" {
		query.Where(productspec.NameContains(params.Name))
	}

	if params.OnlyMerchant {
		query.Where(productspec.StoreID(uuid.Nil))
	} else if params.StoreID != uuid.Nil {
		query.Where(productspec.StoreID(params.StoreID))
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
	entSpecs, err := query.Order(ent.Desc(productspec.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.ProductSpecs, 0, len(entSpecs))
	for _, s := range entSpecs {
		items = append(items, convertProductSpecToDomain(s))
	}

	page.SetTotal(total)

	return &domain.ProductSpecSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

func (repo *ProductSpecRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) (res domain.ProductSpecs, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductSpecRepository.ListByIDs")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductSpec.Query().Where(productspec.IDIn(ids...))

	entSpecs, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	res = make(domain.ProductSpecs, 0, len(entSpecs))
	for _, s := range entSpecs {
		res = append(res, convertProductSpecToDomain(s))
	}

	return res, nil
}

func (repo *ProductSpecRepository) FindByNamesInStore(ctx context.Context, storeID uuid.UUID, names []string) (res domain.ProductSpecs, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductSpecRepository.FindByNamesInStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProductSpec.Query().Where(productspec.StoreID(storeID)).Where(productspec.NameIn(names...))
	entSpecs, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	res = make(domain.ProductSpecs, 0, len(entSpecs))
	for _, s := range entSpecs {
		res = append(res, convertProductSpecToDomain(s))
	}
	return res, nil
}

func (repo *ProductSpecRepository) CreateBulk(ctx context.Context, specs domain.ProductSpecs) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductSpecRepository.CreateBulk")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	if len(specs) == 0 {
		return nil
	}

	builders := make([]*ent.ProductSpecCreate, 0, len(specs))
	for _, spec := range specs {
		builder := repo.Client.ProductSpec.Create().
			SetID(spec.ID).
			SetName(spec.Name).
			SetMerchantID(spec.MerchantID).
			SetStoreID(spec.StoreID).
			SetProductCount(spec.ProductCount)

		builders = append(builders, builder)
	}
	_, err = repo.Client.ProductSpec.CreateBulk(builders...).Save(ctx)
	return err
}

func convertProductSpecToDomain(es *ent.ProductSpec) *domain.ProductSpec {
	if es == nil {
		return nil
	}

	spec := &domain.ProductSpec{
		ID:           es.ID,
		Name:         es.Name,
		MerchantID:   es.MerchantID,
		StoreID:      es.StoreID,
		ProductCount: es.ProductCount,
		CreatedAt:    es.CreatedAt,
		UpdatedAt:    es.UpdatedAt,
	}

	return spec
}
