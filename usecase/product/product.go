package product

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductInteractor = (*ProductInteractor)(nil)

type ProductInteractor struct {
	ds domain.DataStore
}

func NewProductInteractor(ds domain.DataStore) *ProductInteractor {
	return &ProductInteractor{
		ds: ds,
	}
}

func (i *ProductInteractor) Create(ctx context.Context, params domain.ProductUpsetParams) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	user := domain.FromBackendUserContext(ctx)
	product := params.Product
	product.StoreID = user.Store.ID
	params.RecipeIDs = lo.Uniq(params.RecipeIDs)
	params.AttrIDs = lo.Uniq(params.AttrIDs)
	// 非套餐商品必须传递单位
	if product.Type != domain.ProductTypeSetMeal && product.UnitID == 0 {
		return domain.ParamsError(domain.ErrProductNeedUnit)
	}
	// 多规格商品必须传递规格参数
	if product.Type == domain.ProductTypeMulti {
		if len(product.Specs) == 0 {
			return domain.ParamsError(domain.ErrProductTypeMulti)
		}
		product.Specs = lo.UniqBy(product.Specs, func(item *domain.ProductSpecRel) int {
			return item.SpecID
		})
	}
	// 套餐商品必须传递套餐详情
	if product.Type == domain.ProductTypeSetMeal {
		if len(params.SetMealDetails) == 0 {
			return domain.ParamsError(domain.ErrSetMealProductNotExists)
		}
		params.SetMealDetails = lo.UniqBy(params.SetMealDetails, func(item *domain.SetMealDetail) int {
			return item.ProductID
		})
		// 计算总价
		product.Price = lo.Reduce(params.SetMealDetails, func(acc decimal.Decimal, item *domain.SetMealDetail, _ int) decimal.Decimal {
			return acc.Add(item.SetMealPrice.Mul(item.Quantity))
		}, decimal.Zero)
	}

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 名称唯一性校验
		exists, err := ds.ProductRepo().Exists(ctx, domain.ProductExistsParams{
			StoreID: product.StoreID,
			Name:    product.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrProductNameExists)
		}
		// 2. 校验基础数据
		if err := i.validateBasicData(ctx, ds, product, params.AttrIDs, params.RecipeIDs); err != nil {
			return err
		}

		// 3. 校验商品规格
		if err := i.validateSpecData(ctx, ds, product); err != nil {
			return err
		}

		// 4. 校验套餐商品
		if err := i.validateSetMealData(ctx, ds, product, params.SetMealDetails); err != nil {
			return err
		}

		// 5. 创建商品
		product.Status = lo.Ternary(user.Store.NeedAudit, domain.ProductStatusUnApprove, domain.ProductStatusApproved)
		product.SaleStatus = domain.ProductSaleStatusOn

		err = ds.ProductRepo().Create(ctx, product, params.AttrIDs, params.RecipeIDs)
		if err != nil {
			return err
		}

		// 6. 创建套餐详情
		if product.Type == domain.ProductTypeSetMeal {
			for _, item := range params.SetMealDetails {
				item.SetMealID = product.ID
			}
			if err := ds.SetMealDetailRepo().BatchCreate(ctx, params.SetMealDetails); err != nil {
				return err
			}
		}
		return nil
	})
}

// 校验基础数据（分类、单位）
func (i *ProductInteractor) validateBasicData(ctx context.Context, ds domain.DataStore,
	product *domain.Product,
	attrIDs, recipeIDs []int,
) error {
	// 校验分类
	if product.CategoryID > 0 {
		category, err := ds.ProductCategoryRepo().FindByID(ctx, product.CategoryID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrCategoryNotExists)
			}
			return err
		}
		if category.StoreID != product.StoreID {
			return domain.ParamsError(domain.ErrCategoryNotExists)
		}
	}

	// 校验单位
	if product.UnitID > 0 {
		unit, err := ds.ProductUnitRepo().FindByID(ctx, product.UnitID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrUnitNotExists)
			}
			return err
		}
		if unit.StoreID != product.StoreID {
			return domain.ParamsError(domain.ErrUnitNotExists)
		}
	}

	// 校验属性
	if len(attrIDs) > 0 {
		existAttrs, err := ds.ProductAttrRepo().ListByIDs(ctx, attrIDs)
		if err != nil {
			return err
		}
		if len(existAttrs) != len(attrIDs) {
			return domain.ParamsError(domain.ErrAttrNotExists)
		}
		storeNotMatch := lo.SomeBy(existAttrs, func(attr *domain.ProductAttr) bool {
			return attr.StoreID != product.StoreID
		})
		if storeNotMatch {
			return domain.ParamsError(domain.ErrAttrNotExists)
		}
	}

	// 校验做法
	if len(recipeIDs) > 0 {
		existRecipes, err := ds.ProductRecipeRepo().ListByIDs(ctx, recipeIDs)
		if err != nil {
			return err
		}
		if len(existRecipes) != len(recipeIDs) {
			return domain.ParamsError(domain.ErrRecipeNotExists)
		}
		storeNotMatch := lo.SomeBy(existRecipes, func(recipe *domain.ProductRecipe) bool {
			return recipe.StoreID != product.StoreID
		})
		if storeNotMatch {
			return domain.ParamsError(domain.ErrRecipeNotExists)
		}
	}
	return nil
}

// 校验商品规格
func (i *ProductInteractor) validateSpecData(ctx context.Context, ds domain.DataStore,
	product *domain.Product,
) error {
	if product.Type != domain.ProductTypeMulti {
		return nil
	}
	specIDs := lo.Map(product.Specs, func(spec *domain.ProductSpecRel, _ int) int {
		return spec.SpecID
	})

	existSpecs, err := ds.ProductSpecRepo().ListByIDs(ctx, specIDs)
	if err != nil {
		return err
	}
	if len(existSpecs) != len(specIDs) {
		return domain.ParamsError(domain.ErrSpecNotExists)
	}

	storeNotMatch := lo.SomeBy(existSpecs, func(spec *domain.ProductSpec) bool {
		return spec.StoreID != product.StoreID
	})
	if storeNotMatch {
		return domain.ParamsError(domain.ErrSpecNotExists)
	}

	return nil
}

// 校验套餐数据
func (i *ProductInteractor) validateSetMealData(ctx context.Context, ds domain.DataStore,
	product *domain.Product,
	setMealDetails domain.SetMealDetails,
) error {
	if product.Type != domain.ProductTypeSetMeal {
		return nil
	}
	productIDs := lo.Map(setMealDetails, func(detail *domain.SetMealDetail, _ int) int {
		return detail.ProductID
	})
	existProducts, err := ds.ProductRepo().ListByIDs(ctx, productIDs)
	if err != nil {
		return err
	}
	if len(existProducts) != len(productIDs) {
		return domain.ParamsError(domain.ErrProductNotExists)
	}

	// 校验商品基础状态
	invalidProduct := lo.SomeBy(existProducts, func(p *domain.Product) bool {
		return p.StoreID != product.StoreID || p.Type == domain.ProductTypeSetMeal
	})
	if invalidProduct {
		return domain.ParamsError(domain.ErrSetMealProductInvalid)
	}

	// 校验商品规格信息
	existProductMap := lo.SliceToMap(existProducts, func(item *domain.Product) (int, *domain.Product) {
		return item.ID, item
	})
	var productSpecIDs []int
	productSpecIDMap := make(map[int]*domain.Product)
	for _, item := range setMealDetails {
		if item.Spec.ID == 0 {
			continue
		}
		productSpecIDs = append(productSpecIDs, item.Spec.ID)
		existProduct, ok := existProductMap[item.ProductID]
		if !ok {
			return domain.ParamsError(domain.ErrSetMealProductInvalid)
		}
		productSpecIDMap[item.Spec.ID] = existProduct
	}
	existProductSpecs, err := ds.ProductSpecRelRepo().ListByIDs(ctx, productSpecIDs)
	if err != nil {
		return err
	}
	if len(existProductSpecs) != len(productSpecIDs) {
		return domain.ParamsError(domain.ErrSetMealProductInvalid)
	}
	for _, item := range existProductSpecs {
		// 检查跟product是不是对应的
		_, ok := productSpecIDMap[item.ID]
		if !ok {
			return domain.ParamsError(domain.ErrSetMealProductInvalid)
		}
	}
	return nil
}

func (i *ProductInteractor) Update(ctx context.Context, params domain.ProductUpsetParams) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	user := domain.FromBackendUserContext(ctx)
	product := params.Product
	product.StoreID = user.Store.ID
	params.RecipeIDs = lo.Uniq(params.RecipeIDs)
	params.AttrIDs = lo.Uniq(params.AttrIDs)

	// 多规格商品必须传递规格参数
	if product.Type == domain.ProductTypeMulti {
		if len(product.Specs) == 0 {
			return domain.ParamsError(domain.ErrProductTypeMulti)
		}
		product.Price = decimal.Zero
		product.Specs = lo.UniqBy(product.Specs, func(item *domain.ProductSpecRel) int {
			return item.SpecID
		})
	}
	// 套餐商品必须传递套餐详情
	if product.Type == domain.ProductTypeSetMeal {
		if len(params.SetMealDetails) == 0 {
			return domain.ParamsError(domain.ErrSetMealProductNotExists)
		}
		params.SetMealDetails = lo.UniqBy(params.SetMealDetails, func(item *domain.SetMealDetail) int {
			return item.ProductID
		})
		// 计算总价
		product.Price = lo.Reduce(params.SetMealDetails, func(acc decimal.Decimal, item *domain.SetMealDetail, _ int) decimal.Decimal {
			return acc.Add(item.SetMealPrice.Mul(item.Quantity))
		}, decimal.Zero)
	}

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证商品是否存在
		existProduct, err := ds.ProductRepo().FindByID(ctx, product.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductNotExists)
			}
			return nil
		}
		if existProduct.StoreID != product.StoreID {
			return domain.ParamsError(domain.ErrProductNotExists)
		}
		if existProduct.Status != domain.ProductStatusUnApprove {
			return domain.ParamsError(domain.ErrProductStatus)
		}

		// 非套餐商品必须传递单位
		if existProduct.Type != domain.ProductTypeSetMeal && product.UnitID == 0 {
			return domain.ParamsError(domain.ErrProductNeedUnit)
		}

		// 名称唯一性校验
		if existProduct.Name != product.Name {
			exists, err := ds.ProductRepo().Exists(ctx, domain.ProductExistsParams{
				StoreID: product.StoreID,
				Name:    product.Name,
			})
			if err != nil {
				return err
			}
			if exists {
				return domain.ParamsError(domain.ErrProductNameExists)
			}
		}

		// 2. 校验基础数据
		if err := i.validateBasicData(ctx, ds, product, params.AttrIDs, params.RecipeIDs); err != nil {
			return err
		}

		// 3. 校验商品规格
		if err := i.validateSpecData(ctx, ds, product); err != nil {
			return err
		}

		// 4. 校验套餐商品
		if err := i.validateSetMealData(ctx, ds, product, params.SetMealDetails); err != nil {
			return err
		}

		// 5. 更新商品
		product.Status = lo.Ternary(user.Store.NeedAudit, domain.ProductStatusUnApprove, domain.ProductStatusApproved)
		err = ds.ProductRepo().Update(ctx, product, params.AttrIDs, params.RecipeIDs)
		if err != nil {
			return err
		}

		// 6. 更新套餐详情
		if product.Type == domain.ProductTypeSetMeal {
			if err := ds.SetMealDetailRepo().DeleteBySetMealID(ctx, product.ID); err != nil {
				return err
			}

			for _, item := range params.SetMealDetails {
				item.SetMealID = product.ID
			}

			if err := ds.SetMealDetailRepo().BatchCreate(ctx, params.SetMealDetails); err != nil {
				return err
			}
		}
		return nil
	})
}

func (i *ProductInteractor) Delete(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromBackendUserContext(ctx)
	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 获取商品信息
		p, err := ds.ProductRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductNotExists)
			}
			return err
		}

		// 2. 权限校验
		if p.StoreID != user.Store.ID {
			return domain.ParamsError(domain.ErrProductNotExists)
		}

		// 3. 状态校验
		if p.Status != domain.ProductStatusUnApprove {
			return domain.ParamsError(domain.ErrProductStatus)
		}

		// 4. 执行删除
		err = ds.ProductRepo().Delete(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductNotExists)
			}
			return err
		}
		return nil
	})
}

func (i *ProductInteractor) PagedListBySearch(ctx context.Context,
	page *upagination.Pagination, params domain.ProductSearchParams,
) (res *domain.ProductSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.ds.ProductRepo().PagedListBySearch(ctx, page, params)
}

// 实现用例层详情方法
func (i *ProductInteractor) GetDetail(ctx context.Context, id int) (res *domain.Product, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductInteractor.GetDetail")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	product, err := i.ds.ProductRepo().GetDetail(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ParamsError(domain.ErrProductNotExists)
		}
		return nil, err
	}
	return product, nil
}

// ListSetmealDetails 加载套餐详情
func (i *ProductInteractor) ListSetmealDetails(ctx context.Context, setMealID int) (details domain.SetMealDetails, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductInteractor.ListSetmealDetails")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return i.ds.SetMealDetailRepo().ListBySetMealID(ctx, setMealID)
}

func (i *ProductInteractor) Approve(ctx context.Context, ids []int, allowPointPay *bool) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductInteractor.Approve")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	products, err := i.ds.ProductRepo().ListByIDs(ctx, ids)
	if err != nil {
		return err
	}
	if len(products) != len(ids) {
		return domain.ParamsError(domain.ErrProductNotExists)
	}
	for _, product := range products {
		if product.Status != domain.ProductStatusUnApprove {
			return domain.ParamsError(domain.ErrProductStatus)
		}
	}
	return i.ds.ProductRepo().BatchUpdateAttr(ctx, ids, domain.ProductUpdateAttrs{
		Status:        domain.ProductStatusApproved,
		AllowPointPay: allowPointPay,
	})
}

func (i *ProductInteractor) UnApprove(ctx context.Context, ids []int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductInteractor.UnApprove")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	products, err := i.ds.ProductRepo().ListByIDs(ctx, ids)
	if err != nil {
		return err
	}
	if len(products) != len(ids) {
		return domain.ParamsError(domain.ErrProductNotExists)
	}

	for _, product := range products {
		if product.Status != domain.ProductStatusApproved {
			return domain.ParamsError(domain.ErrProductStatus)
		}
	}

	return i.ds.ProductRepo().BatchUpdateAttr(ctx, ids, domain.ProductUpdateAttrs{
		Status: domain.ProductStatusUnApprove,
	})
}

func (i *ProductInteractor) ClearStock(ctx context.Context, productID int, specIDs []int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductInteractor.ClearStock")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromFrontendUserContext(ctx)

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		product, err := i.ds.ProductRepo().GetDetailWithSpec(ctx, productID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductNotExists)
			}
			return err
		}
		if product.StoreID != user.Store.ID {
			return domain.ParamsError(domain.ErrProductNotExists)
		}
		// 单规格商品和套餐商品，直接修改商品状态
		switch product.Type {
		case domain.ProductTypeSingle, domain.ProductTypeSetMeal:
			return i.ds.ProductRepo().UpdateAttr(ctx, productID, domain.ProductUpdateAttrs{
				SaleStatus: domain.ProductSaleStatusOff,
			})
			// 多规格商品，修改对应规格的状态
		case domain.ProductTypeMulti:
			// 获取商品的规格ID列表
			var productSpecIDs, saleOffProductSpecIDs []int
			for _, item := range product.Specs {
				productSpecIDs = append(productSpecIDs, item.ID)
				if item.SaleStatus == domain.ProductSaleStatusOff {
					saleOffProductSpecIDs = append(saleOffProductSpecIDs, item.ID)
				}
			}
			diff, _ := lo.Difference(specIDs, productSpecIDs)
			if len(diff) > 0 {
				return domain.ParamsError(domain.ErrSpecNotExists)
			}
			// 更新商品规格的售卖状态
			err = i.ds.ProductSpecRelRepo().UpdateSaleStatusByIDs(ctx, specIDs, domain.ProductSaleStatusOff)
			if err != nil {
				return err
			}
			// 更新商品的售卖状态：部分规格被估清；全部估清
			partSaleOff := len(lo.Union(specIDs, saleOffProductSpecIDs)) != len(product.Specs)
			saleStatus := domain.ProductSaleStatusOff
			if partSaleOff {
				saleStatus = domain.ProductSaleStatusPartOff
			}
			return i.ds.ProductRepo().UpdateAttr(ctx, productID, domain.ProductUpdateAttrs{
				SaleStatus: saleStatus,
			})
		default:
			return nil
		}
	})
}

func (i *ProductInteractor) RestoreStock(ctx context.Context, productID int, specIDs []int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductInteractor.RestoreStock")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user := domain.FromFrontendUserContext(ctx)

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		product, err := i.ds.ProductRepo().GetDetailWithSpec(ctx, productID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductNotExists)
			}
			return err
		}
		if product.StoreID != user.Store.ID {
			return domain.ParamsError(domain.ErrProductNotExists)
		}
		// 单规格商品和套餐商品，直接修改商品状态
		switch product.Type {
		case domain.ProductTypeSingle, domain.ProductTypeSetMeal:
			return i.ds.ProductRepo().UpdateAttr(ctx, productID, domain.ProductUpdateAttrs{
				SaleStatus: domain.ProductSaleStatusOn,
			})
			// 多规格商品，修改对应规格的状态
		case domain.ProductTypeMulti:
			// 获取商品的规格ID列表
			var productSpecIDs, saleOnProductSpecIDs []int
			for _, item := range product.Specs {
				productSpecIDs = append(productSpecIDs, item.ID)
				if item.SaleStatus == domain.ProductSaleStatusOn {
					saleOnProductSpecIDs = append(saleOnProductSpecIDs, item.ID)
				}
			}
			diff, _ := lo.Difference(specIDs, productSpecIDs)
			if len(diff) > 0 {
				return domain.ParamsError(domain.ErrSpecNotExists)
			}
			// 更新商品规格的售卖状态
			err = i.ds.ProductSpecRelRepo().UpdateSaleStatusByIDs(ctx, specIDs, domain.ProductSaleStatusOn)
			if err != nil {
				return err
			}
			// 更新商品的售卖状态：部分规格被估清；在售
			partSaleOn := len(lo.Union(specIDs, saleOnProductSpecIDs)) != len(product.Specs)
			saleStatus := domain.ProductSaleStatusOn
			if partSaleOn {
				saleStatus = domain.ProductSaleStatusPartOff
			}
			return i.ds.ProductRepo().UpdateAttr(ctx, productID, domain.ProductUpdateAttrs{
				SaleStatus: saleStatus,
			})
		default:
			return nil
		}
	})
}
