package repository

import (
	"context"
	"strconv"

	"github.com/samber/lo"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/product"
	"gitlab.jiguang.dev/pos-dine/dine/ent/productspec"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductRepository = (*ProductRepository)(nil)

type ProductRepository struct {
	Client *ent.Client
}

func NewProductRepository(client *ent.Client) *ProductRepository {
	return &ProductRepository{
		Client: client,
	}
}

func (repo *ProductRepository) FindByID(ctx context.Context, id int) (res *domain.Product, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	product, err := repo.Client.Product.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProductNotExists)
		}
		return nil, err
	}
	return repo.convertToDomain(product), nil
}

func (repo *ProductRepository) Exists(ctx context.Context, params domain.ProductExistsParams) (exists bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Product.Query()
	if params.StoreID != 0 {
		query.Where(product.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(product.Name(params.Name))
	}
	if params.UnitID > 0 {
		query.Where(product.UnitID(params.UnitID))
	}
	if params.CategoryID > 0 {
		query.Where(product.CategoryID(params.CategoryID))
	}

	return query.Exist(ctx)
}

func (repo *ProductRepository) Create(ctx context.Context,
	p *domain.Product, attrIDs, recipeIDs []int,
) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 创建商品主体
	create := repo.Client.Product.Create().
		SetName(p.Name).
		SetType(int(p.Type)).
		SetPrice(p.Price).
		SetStatus(int(p.Status)).
		SetSaleStatus(int(p.SaleStatus)).
		SetStoreID(p.StoreID).
		SetImages(p.Images).
		SetAllowPointPay(p.AllowPointPay).
		SetCategoryID(p.CategoryID).
		AddAttrIDs(attrIDs...).
		AddRecipeIDs(recipeIDs...)

	if p.UnitID > 0 {
		create.SetUnitID(p.UnitID)
	}

	created, err := create.Save(ctx)
	if err != nil {
		return err
	}

	// 创建商品规格
	if len(p.Specs) > 0 {
		_, err = repo.Client.ProductSpec.MapCreateBulk(p.Specs, func(c *ent.ProductSpecCreate, idx int) {
			c.SetSpecID(p.Specs[idx].SpecID).
				SetPrice(p.Specs[idx].Price).
				SetSaleStatus(int(p.Specs[idx].SaleStatus)).
				SetProductID(created.ID)
		}).Save(ctx)
		if err != nil {
			return err
		}
	}

	p.ID = created.ID
	p.CreatedAt = created.CreatedAt

	return nil
}

func (repo *ProductRepository) Update(ctx context.Context,
	p *domain.Product, attrIDs, recipeIDs []int,
) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 更新商品主体
	update := repo.Client.Product.UpdateOneID(p.ID).
		SetName(p.Name).
		//SetType(int(p.Type)).
		SetPrice(p.Price).
		SetImages(p.Images).
		SetCategoryID(p.CategoryID)

	if p.UnitID > 0 {
		update.SetUnitID(p.UnitID)
	}
	if len(attrIDs) > 0 {
		update.ClearAttrs().AddAttrIDs(attrIDs...)
	}
	if len(recipeIDs) > 0 {
		update.ClearRecipes().AddRecipeIDs(recipeIDs...)
	}

	if _, err := update.Save(ctx); err != nil {
		return err
	}

	// 批量更新商品规格
	if len(p.Specs) > 0 {
		_, err = repo.Client.ProductSpec.Delete().
			Where(productspec.ProductID(p.ID)).
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = repo.Client.ProductSpec.MapCreateBulk(p.Specs, func(c *ent.ProductSpecCreate, idx int) {
			c.SetSpecID(p.Specs[idx].SpecID).
				SetPrice(p.Specs[idx].Price).
				SetSaleStatus(int(p.Specs[idx].SaleStatus)).
				SetProductID(p.ID)
		}).Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *ProductRepository) ListByIDs(ctx context.Context, ids []int) (res domain.Products, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRepository.ListByIDs")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	products, err := repo.Client.Product.Query().
		Where(product.IDIn(ids...)).
		All(ctx)

	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, nil
	}
	for _, u := range products {
		res = append(res, repo.convertToDomain(u))
	}
	return res, nil
}

func (repo *ProductRepository) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	affected, err := repo.Client.Product.Delete().
		Where(
			product.ID(id),
			product.StatusEQ(int(domain.ProductStatusUnApprove)),
		).
		Exec(ctx)

	if err != nil {
		return err
	}

	if affected == 0 {
		return domain.NotFoundError(domain.ErrProductNotExists)
	}
	return nil
}

func (repo *ProductRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.ProductSearchParams,
) (res *domain.ProductSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Product.Query()

	if params.StoreID > 0 {
		query.Where(product.StoreID(params.StoreID))
	}
	if params.Name != "" {
		// 添加类型安全的条件判断
		if id, err := strconv.Atoi(params.Name); err == nil {
			// 当名称参数为数字时，同时匹配ID和名称
			query.Where(
				product.Or(
					product.NameContains(params.Name),
					product.ID(id),
				),
			)
		} else {
			// 非数字时仅进行名称模糊查询
			query.Where(product.NameContains(params.Name))
		}
	}

	if params.CategoryID > 0 {
		query.Where(product.CategoryID(params.CategoryID))
	}

	if params.Status > 0 {
		query.Where(product.Status(int(params.Status)))
	}

	if len(params.SaleStatus) > 0 {
		saleStatusInts := lo.Map(params.SaleStatus, func(s domain.ProductSaleStatus, _ int) int {
			return int(s)
		})
		query.Where(product.SaleStatusIn(saleStatusInts...))
	}

	if params.Type > 0 {
		query.Where(product.Type(int(params.Type)))
	}

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = total

	products, err := query.Order(ent.Desc(product.FieldID)).
		WithProductSpecs(func(q *ent.ProductSpecQuery) {
			q.WithSpec() // 加载规格基础信息
		}).
		WithUnit().
		WithSetMealDetails(func(q *ent.SetMealDetailQuery) {
			q.WithProduct()
		}).
		Offset(page.Offset()).
		Limit(page.Size).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.Products, 0, len(products))
	for _, c := range products {
		items = append(items, repo.convertToDomain(c))
	}

	return &domain.ProductSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

// 增强详情查询方法
func (repo *ProductRepository) GetDetail(ctx context.Context, id int) (res *domain.Product, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRepository.GetDetail")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	p, err := repo.Client.Product.Query().
		Where(product.ID(id)).
		WithCategory(). // 加载分类
		WithUnit().     // 加载单位
		WithAttrs().    // 加载属性
		WithRecipes().  // 加载做法
		WithProductSpecs(func(q *ent.ProductSpecQuery) {
			q.WithSpec() // 加载规格基础信息
		}).
		WithSetMealDetails(func(q *ent.SetMealDetailQuery) {
			q.WithProduct(func(q *ent.ProductQuery) {
				q.WithUnit().
					WithProductSpecs(func(q *ent.ProductSpecQuery) {
						q.WithSpec()
					})
			}).
				WithProductSpec(func(q *ent.ProductSpecQuery) {
					q.WithSpec()
				})
		}).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProductNotExists)
		}
		return nil, err
	}
	return repo.convertToDomain(p), nil
}

func (repo *ProductRepository) convertToDomain(p *ent.Product) *domain.Product {
	product := &domain.Product{
		ID:            p.ID,
		Name:          p.Name,
		Type:          domain.ProductType(p.Type),
		Price:         p.Price,
		Status:        domain.ProductStatus(p.Status),
		SaleStatus:    domain.ProductSaleStatus(p.SaleStatus),
		StoreID:       p.StoreID,
		Images:        p.Images,
		AllowPointPay: p.AllowPointPay,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
		UnitID:        p.UnitID,
		CategoryID:    p.CategoryID,
	}

	// 分类信息
	if p.Edges.Category != nil {
		product.Category = &domain.Category{
			ID:   p.Edges.Category.ID,
			Name: p.Edges.Category.Name,
		}
	}

	// 单位信息
	if p.Edges.Unit != nil {
		product.Unit = &domain.ProductUnit{
			ID:        p.Edges.Unit.ID,
			Name:      p.Edges.Unit.Name,
			CreatedAt: p.Edges.Unit.CreatedAt,
			UpdatedAt: p.Edges.Unit.UpdatedAt,
		}
	}

	// 商品属性
	for _, attr := range p.Edges.Attrs {
		product.Attrs = append(product.Attrs, &domain.ProductAttr{
			ID:        attr.ID,
			Name:      attr.Name,
			CreatedAt: attr.CreatedAt,
			UpdatedAt: attr.UpdatedAt,
		})
	}

	// 商品做法
	for _, recipe := range p.Edges.Recipes {
		product.Recipes = append(product.Recipes, &domain.ProductRecipe{
			ID:        recipe.ID,
			Name:      recipe.Name,
			CreatedAt: recipe.CreatedAt,
			UpdatedAt: recipe.UpdatedAt,
		})
	}

	// 转换规格信息
	for _, spec := range p.Edges.ProductSpecs {
		productSpec := &domain.ProductSpecRel{
			ID:         spec.ID,
			SpecID:     spec.SpecID,
			ProductID:  p.ID,
			Price:      spec.Price,
			SaleStatus: domain.ProductSaleStatus(spec.SaleStatus),
		}
		if spec.Edges.Spec != nil {
			productSpec.SpecName = spec.Edges.Spec.Name
		}
		product.Specs = append(product.Specs, productSpec)
	}

	var setmealDetails domain.SetMealDetails
	if len(p.Edges.SetMealDetails) > 0 {
		pd := p.Edges.SetMealDetails
		for _, d := range pd {
			item := &domain.SetMealDetail{
				ID:           d.ID,
				SetMealID:    d.SetMealID,
				ProductID:    d.ProductID,
				Quantity:     d.Quantity,
				SetMealPrice: d.Price,
				CreatedAt:    d.CreatedAt,
				UpdatedAt:    d.UpdatedAt,
			}
			// 具体规格
			if d.Edges.ProductSpec != nil {
				item.Spec = &domain.ProductSpecRel{
					ID:     d.ProductSpecID,
					SpecID: d.Edges.ProductSpec.SpecID,
					Price:  d.Edges.ProductSpec.Price,
				}
				if d.Edges.ProductSpec.Edges.Spec != nil {
					item.Spec.SpecName = d.Edges.ProductSpec.Edges.Spec.Name
				}
			}
			// 详情商品信息
			if d.Edges.Product != nil {
				item.Name = d.Edges.Product.Name
				item.Price = d.Edges.Product.Price
				item.ProductType = domain.ProductType(d.Edges.Product.Type)
				item.Images = d.Edges.Product.Images
				unit := d.Edges.Product.Edges.Unit
				if unit != nil {
					item.UnitID = unit.ID
					item.Unit = &domain.ProductUnit{
						ID:        unit.ID,
						Name:      unit.Name,
						CreatedAt: unit.CreatedAt,
						UpdatedAt: unit.UpdatedAt,
					}
				}
				specs := d.Edges.Product.Edges.ProductSpecs
				for _, spec := range specs {
					item.Specs = append(item.Specs, &domain.ProductSpecRel{
						ID:         spec.ID,
						SpecID:     spec.SpecID,
						SpecName:   lo.TernaryF(spec.Edges.Spec != nil, func() string { return spec.Edges.Spec.Name }, func() string { return "" }),
						Price:      spec.Price,
						SaleStatus: domain.ProductSaleStatus(spec.SaleStatus),
					})
				}
			}
			setmealDetails = append(setmealDetails, item)
		}
		product.SetMealDetails = setmealDetails
	}

	return product
}

func (repo *ProductRepository) UpdateAttr(ctx context.Context, id int, attr domain.ProductUpdateAttrs) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRepository.UpdateAttr")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	update := repo.Client.Product.UpdateOneID(id)

	if attr.Status > 0 {
		update.SetStatus(int(attr.Status))
	}
	if attr.SaleStatus > 0 {
		update.SetSaleStatus(int(attr.SaleStatus))
	}
	if attr.AllowPointPay != nil {
		update.SetAllowPointPay(*attr.AllowPointPay)
	}

	if _, err := update.Save(ctx); err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepository) BatchUpdateAttr(ctx context.Context, ids []int, attr domain.ProductUpdateAttrs) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRepository.BatchUpdateAttr")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(ids) == 0 {
		return domain.ErrProductNotExists
	}

	update := repo.Client.Product.Update().Where(product.IDIn(ids...))
	if attr.Status > 0 {
		update.SetStatus(int(attr.Status))
	}
	if attr.SaleStatus > 0 {
		update.SetSaleStatus(int(attr.SaleStatus))
	}
	if attr.AllowPointPay != nil {
		update.SetAllowPointPay(*attr.AllowPointPay)
	}

	if _, err := update.Save(ctx); err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepository) GetDetailWithSpec(ctx context.Context, id int) (res *domain.Product, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRepository.GetDetailWithSpec")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	p, err := repo.Client.Product.Query().
		Where(product.ID(id)).
		WithProductSpecs().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProductNotExists)
		}
		return nil, err
	}
	return repo.convertToDomain(p), nil
}

func (repo *ProductRepository) GetDetailsByIDs(ctx context.Context, ids []int) (res domain.Products, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRepository.GetDetailsByIDs")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	products, err := repo.Client.Product.Query().
		Where(product.IDIn(ids...)).
		WithCategory(). // 加载分类
		WithUnit().     // 加载单位
		WithAttrs().    // 加载属性
		WithRecipes().  // 加载做法
		WithProductSpecs(func(q *ent.ProductSpecQuery) {
			q.WithSpec() // 加载规格基础信息
		}).
		WithSetMealDetails(func(q *ent.SetMealDetailQuery) {
			q.WithProduct(func(q *ent.ProductQuery) {
				q.WithUnit().
					WithProductSpecs()
			}).WithProductSpec(func(q *ent.ProductSpecQuery) {
				q.WithSpec()
			})
		}).
		All(ctx)

	if err != nil {
		return nil, err
	}

	for _, p := range products {
		res = append(res, repo.convertToDomain(p))
	}
	return res, nil
}
