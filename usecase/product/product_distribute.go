package product

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

const (
	maxWorkers = 10 // 最大并发数
)

type jobResult struct {
	storeID uuid.UUID
	err     error
}

func (i *ProductInteractor) Distribute(ctx context.Context, params domain.ProductDistributeParams, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductInteractor.Distribute")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	logger := logging.FromContext(ctx).Named("ProductInteractor.Distribute")
	// 1. 验证商品存在且属于当前品牌商
	product, err := i.DS.ProductRepo().GetDetail(ctx, params.ProductID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrProductNotExists)
		}
		return err
	}
	if product.MerchantID != params.MerchantID {
		return domain.ParamsError(domain.ErrProductDistributeStoreInvalid)
	}

	// 2. 验证是否可以操作该商品
	if err := verifyProductOwnership(user, product); err != nil {
		return err
	}

	// 3. 验证门店存在且属于当前品牌商
	storeIDs := lo.Uniq(params.StoreIDs)
	stores, err := i.DS.StoreRepo().ListByIDs(ctx, storeIDs)
	if err != nil {
		return err
	}
	if len(stores) != len(storeIDs) {
		return domain.ParamsError(domain.ErrStoreNotExists)
	}

	for _, store := range stores {
		if store.MerchantID != params.MerchantID {
			return domain.ParamsError(domain.ErrStoreNotExists)
		}
	}

	jobs := make(chan uuid.UUID, len(stores))
	results := make(chan jobResult, len(stores))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// workers
	worker := func(workerID int) {
		for storeID := range jobs {
			select {
			case <-ctx.Done():
				return
			default:
				logger.Infof("开始处理门店下任务：worker_id=%d, store_id=%s", workerID, storeID.String())
				err := i.distributeToStore(ctx, product, storeID)
				results <- jobResult{
					storeID: storeID,
					err:     err,
				}
			}

		}
	}
	// 启动 worker
	workerCount := maxWorkers
	if len(storeIDs) < maxWorkers {
		workerCount = len(storeIDs)
	}
	for i := 0; i < workerCount; i++ {
		go worker(i)
	}

	// 投递任务
	for _, storeID := range storeIDs {
		jobs <- storeID
	}
	close(jobs)

	// 收集结果
	var failed []string

	for range len(storeIDs) {
		r := <-results
		if r.err != nil {
			failed = append(
				failed,
				fmt.Sprintf("store_id=%s err=%v", r.storeID.String(), r.err),
			)
		}
	}

	if len(failed) > 0 {
		return domain.ParamsErrorf("部分门店下发失败: %s", strings.Join(failed, "; "))
	}
	return nil
}

// distributeToStore 下发商品到单个门店
func (i *ProductInteractor) distributeToStore(
	ctx context.Context,
	brandProduct *domain.Product,
	storeID uuid.UUID,
) error {
	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 检查门店是否存在同名商品
		_, err := ds.ProductRepo().FindByNameInStore(ctx, storeID, brandProduct.Name)
		if err == nil {
			// 门店已存在同名商品，直接结束
			return nil
		}
		if !domain.IsNotFound(err) {
			return err
		}

		// 2. 门店不存在同名商品，创建新商品

		// 2.1 处理分类：如果门店没有同名分类则创建，如果有则使用门店的
		storeCategoryID, err := i.findOrCreateCategory(ctx, ds, brandProduct.Category, storeID)
		if err != nil {
			return err
		}

		// 2.2 处理单位：如果门店没有同名单位则创建，如果有则使用门店的
		storeUnitID, err := i.findOrCreateUnit(ctx, ds, brandProduct.Unit, storeID)
		if err != nil {
			return err
		}

		// 2.3 处理标签：如果门店没有同名标签则创建，如果有则使用门店的
		storeTags := make(domain.ProductTags, 0, len(brandProduct.Tags))
		if len(brandProduct.Tags) > 0 {
			storeTags, err = i.findOrCreateTags(ctx, ds, brandProduct.Tags, storeID)
			if err != nil {
				return err
			}
		}

		// 2.4 创建门店商品
		storeProduct := brandProduct
		storeProduct.ID = uuid.New()
		storeProduct.StoreID = storeID
		storeProduct.CategoryID = storeCategoryID
		storeProduct.UnitID = storeUnitID
		storeProduct.Tags = storeTags

		if err = ds.ProductRepo().Create(ctx, storeProduct); err != nil {
			return err
		}

		// 2.5 处理规格：如果门店没有同名规格则创建，如果有则使用门店的
		if len(brandProduct.SpecRelations) > 0 {
			storeSpecRelations, err := i.findOrCreateSpecs(ctx, ds, brandProduct.SpecRelations, storeProduct.ID, brandProduct.MerchantID, storeID)
			if err != nil {
				return err
			}
			if len(storeSpecRelations) > 0 {
				if err = ds.ProductSpecRelRepo().CreateBulk(ctx, storeSpecRelations); err != nil {
					return err
				}
			}
		}

		// 2.6 处理口味做法：如果门店没有同名口味做法则创建，如果有则使用门店的
		if len(brandProduct.AttrRelations) > 0 {
			storeAttrRelations, err := i.findOrCreateAttrs(ctx, ds, brandProduct.AttrRelations, storeProduct.ID, brandProduct.MerchantID, storeID)
			if err != nil {
				return err
			}
			if len(storeAttrRelations) > 0 {
				if err = ds.ProductAttrRelRepo().CreateBulk(ctx, storeAttrRelations); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// findOrCreateCategory 查找或创建分类
func (i *ProductInteractor) findOrCreateCategory(
	ctx context.Context,
	ds domain.DataStore,
	brandCategory *domain.Category,
	storeID uuid.UUID,
) (uuid.UUID, error) {
	// 如果是一级分类
	if brandCategory.ParentID == uuid.Nil {
		// 查找门店中是否存在同名一级分类
		storeCategory, err := ds.CategoryRepo().FindByNameInStore(ctx, brandCategory.Name, storeID, uuid.Nil)
		if err != nil && !domain.IsNotFound(err) {
			return uuid.Nil, err
		}
		if storeCategory != nil {
			return storeCategory.ID, nil
		}
		// 门店不存在同名分类，创建新分类
		storeCategory = &domain.Category{
			ID:         uuid.New(),
			Name:       brandCategory.Name,
			MerchantID: brandCategory.MerchantID,
			StoreID:    storeID,
		}
		if err := ds.CategoryRepo().Create(ctx, storeCategory); err != nil {
			return uuid.Nil, err
		}
		return storeCategory.ID, nil
	}

	// 如果是二级分类
	// 检查一级分类是否存在
	storeRootCategory, err := ds.CategoryRepo().FindByNameInStore(ctx, brandCategory.Parent.Name, storeID, uuid.Nil)
	if err != nil && !domain.IsNotFound(err) {
		return uuid.Nil, err
	}
	// 不存在，同时创建一级分类和二级分类
	if storeRootCategory == nil {
		storeRootCategory = &domain.Category{
			ID:         uuid.New(),
			Name:       brandCategory.Parent.Name,
			MerchantID: brandCategory.MerchantID,
			StoreID:    storeID,
		}
		if err := ds.CategoryRepo().Create(ctx, storeRootCategory); err != nil {
			return uuid.Nil, err
		}
		storeCategory := &domain.Category{
			ID:         uuid.New(),
			Name:       brandCategory.Name,
			MerchantID: brandCategory.MerchantID,
			StoreID:    storeID,
			ParentID:   storeRootCategory.ID,
		}
		if err := ds.CategoryRepo().Create(ctx, storeCategory); err != nil {
			return uuid.Nil, err
		}
		return storeCategory.ID, nil
	}

	// 存在，检查二级分类是否存在
	// 存在一级，再检查二级分类是否存在
	storeCategory, err := ds.CategoryRepo().FindByNameInStore(ctx, brandCategory.Name, storeID, storeRootCategory.ID)
	if err != nil && !domain.IsNotFound(err) {
		return uuid.Nil, err
	}
	if storeCategory != nil {
		return storeCategory.ID, nil
	}

	// 不存在，创建二级分类
	storeCategory = &domain.Category{
		ID:         uuid.New(),
		Name:       brandCategory.Name,
		MerchantID: brandCategory.MerchantID,
		StoreID:    storeID,
		ParentID:   storeRootCategory.ID,
	}
	if err := ds.CategoryRepo().Create(ctx, storeCategory); err != nil {
		return uuid.Nil, err
	}
	return storeCategory.ID, nil
}

// findOrCreateUnit 查找或创建单位
func (i *ProductInteractor) findOrCreateUnit(
	ctx context.Context,
	ds domain.DataStore,
	brandUnit *domain.ProductUnit,
	storeID uuid.UUID,
) (uuid.UUID, error) {
	storeUnit, err := ds.ProductUnitRepo().FindByNameInStore(ctx, storeID, brandUnit.Name)
	if err != nil && !domain.IsNotFound(err) {
		return uuid.Nil, err
	}
	if storeUnit != nil {
		return storeUnit.ID, nil
	}
	// 门店不存在同名单位，创建新单位
	storeUnit = &domain.ProductUnit{
		ID:         uuid.New(),
		Name:       brandUnit.Name,
		Type:       brandUnit.Type,
		MerchantID: brandUnit.MerchantID,
		StoreID:    storeID,
	}
	if err := ds.ProductUnitRepo().Create(ctx, storeUnit); err != nil {
		return uuid.Nil, err
	}
	return storeUnit.ID, nil
}

// findOrCreateSpecs 查找或创建规格
func (i *ProductInteractor) findOrCreateSpecs(
	ctx context.Context,
	ds domain.DataStore,
	brandSpecRelations domain.ProductSpecRelations,
	productID, merchantID, storeID uuid.UUID,
) (domain.ProductSpecRelations, error) {
	// 1. 提取规格名称
	specNames := lo.Map(brandSpecRelations, func(specRel *domain.ProductSpecRelation, _ int) string {
		return specRel.SpecName
	})
	// 2. 查找门店中已存在的规格
	storeSpecs, err := ds.ProductSpecRepo().FindByNamesInStore(ctx, storeID, specNames)
	if err != nil {
		return nil, err
	}
	// 3. 构建门店规格名称到ID的映射
	storeSpecNameIDMap := lo.SliceToMap(storeSpecs, func(spec *domain.ProductSpec) (string, uuid.UUID) {
		return spec.Name, spec.ID
	})

	// 4. 找出需要创建的规格
	notExistsSpecs := make(domain.ProductSpecs, 0)
	for _, specRel := range brandSpecRelations {
		if _, exists := storeSpecNameIDMap[specRel.SpecName]; !exists {
			spec := &domain.ProductSpec{
				ID:         uuid.New(),
				Name:       specRel.SpecName,
				MerchantID: merchantID,
				StoreID:    storeID,
			}
			notExistsSpecs = append(notExistsSpecs, spec)
			storeSpecNameIDMap[specRel.SpecName] = spec.ID
		}
	}

	// 5. 批量创建不存在的规格
	if len(notExistsSpecs) > 0 {
		if err := ds.ProductSpecRepo().CreateBulk(ctx, notExistsSpecs); err != nil {
			return nil, err
		}
	}

	// 6. 构建门店规格关联
	for _, specRel := range brandSpecRelations {
		specID, ok := storeSpecNameIDMap[specRel.SpecName]
		if !ok {
			return nil, domain.ParamsErrorf("规格名称 %s 不存在", specRel.SpecName)
		}
		specRel.SpecID = specID
		specRel.ProductID = productID
		specRel.ID = uuid.New()
		// 打包费置空
		specRel.PackingFeeID = uuid.Nil
	}

	return brandSpecRelations, nil
}

// findOrCreateAttrs 查找或创建口味做法
func (i *ProductInteractor) findOrCreateAttrs(
	ctx context.Context,
	ds domain.DataStore,
	brandAttrRelations domain.ProductAttrRelations,
	productID, merchantID, storeID uuid.UUID,
) (domain.ProductAttrRelations, error) {
	// 1. 按照 AttrID 进行分组
	attrIDMap := make(map[uuid.UUID]domain.ProductAttrRelations)
	attrNames := make([]string, 0)
	for _, attrRel := range brandAttrRelations {
		attrIDMap[attrRel.AttrID] = append(attrIDMap[attrRel.AttrID], attrRel)
		attrNames = append(attrNames, attrRel.Attr.Name)
	}
	attrNames = lo.Uniq(attrNames)
	// 2. 获取门店中已存在的属性组
	storeAttrs, err := ds.ProductAttrRepo().FindByNamesInStore(ctx, storeID, attrNames)
	if err != nil {
		return nil, err
	}
	storeAttrNameIDMap := lo.SliceToMap(storeAttrs, func(attr *domain.ProductAttr) (string, uuid.UUID) {
		return attr.Name, attr.ID
	})

	// 3. 找出需要创建的属性组
	notExistsAttrs := make(domain.ProductAttrs, 0)
	notExistsAttrItems := make(domain.ProductAttrItems, 0)
	existsAttrIDMap := make(map[uuid.UUID]domain.ProductAttrRelations)
	for key, attrRels := range attrIDMap {
		attrName := attrRels[0].Attr.Name
		storeAttrID, ok := storeAttrNameIDMap[attrName]
		// 不存在，创建属性组，并且创建属性组里面的属性值
		if !ok {
			attr := &domain.ProductAttr{
				ID:         uuid.New(),
				Name:       attrName,
				MerchantID: merchantID,
				StoreID:    storeID,
				Channels:   attrRels[0].Attr.Channels,
			}
			notExistsAttrs = append(notExistsAttrs, attr)

			for _, attrRel := range attrRels {
				attrItem := &domain.ProductAttrItem{
					ID:        uuid.New(),
					AttrID:    attr.ID,
					Name:      attrRel.AttrItem.Name,
					Image:     attrRel.AttrItem.Image,
					BasePrice: attrRel.AttrItem.BasePrice,
				}
				notExistsAttrItems = append(notExistsAttrItems, attrItem)
				// 更新原始的值
				attrRel.AttrID = attr.ID
				attrRel.ID = uuid.New()
				attrRel.AttrItemID = attrItem.ID
				attrRel.ProductID = productID
			}
		} else {
			// 更新原始的值
			for _, attrRel := range attrRels {
				attrRel.AttrID = storeAttrID
				attrRel.ID = uuid.New()
				attrRel.AttrItemID = uuid.Nil
				attrRel.ProductID = productID
			}
			existsAttrIDMap[key] = attrRels
		}
	}

	// 4. 对于已存在的属性组，检查里面的属性值是否存在，如果不存在，创建属性值
	for _, attrRels := range existsAttrIDMap {
		attrItemNames := lo.Map(attrRels, func(attrRel *domain.ProductAttrRelation, _ int) string {
			return attrRel.AttrItem.Name
		})
		attrItemNames = lo.Uniq(attrItemNames)
		attrID := attrRels[0].AttrID

		storeAttrItems, err := ds.ProductAttrRepo().FindItemsByNamesInAttr(ctx, attrID, attrItemNames)
		if err != nil {
			return nil, err
		}

		storeAttrItemNameIDMap := lo.SliceToMap(storeAttrItems, func(attrItem *domain.ProductAttrItem) (string, uuid.UUID) {
			return attrItem.Name, attrItem.ID
		})

		for _, attrRel := range attrRels {
			attrItemName := attrRel.AttrItem.Name
			storeAttrItemID, ok := storeAttrItemNameIDMap[attrItemName]
			if !ok {
				attrItem := &domain.ProductAttrItem{
					ID:        uuid.New(),
					AttrID:    attrID,
					Name:      attrItemName,
					Image:     attrRel.AttrItem.Image,
					BasePrice: attrRel.AttrItem.BasePrice,
				}
				notExistsAttrItems = append(notExistsAttrItems, attrItem)
				attrRel.AttrItemID = attrItem.ID
				attrRel.ID = uuid.New()
				attrRel.ProductID = productID
			} else {
				attrRel.AttrItemID = storeAttrItemID
				attrRel.ID = uuid.New()
				attrRel.ProductID = productID
			}
		}
	}

	if len(notExistsAttrs) > 0 {
		if err := ds.ProductAttrRepo().CreateBulk(ctx, notExistsAttrs); err != nil {
			return nil, err
		}
	}

	if len(notExistsAttrItems) > 0 {
		if err := ds.ProductAttrRepo().CreateItems(ctx, notExistsAttrItems); err != nil {
			return nil, err
		}
	}
	return brandAttrRelations, nil
}

// findOrCreateTags 查找或创建标签
func (i *ProductInteractor) findOrCreateTags(
	ctx context.Context,
	ds domain.DataStore,
	brandTags domain.ProductTags,
	storeID uuid.UUID,
) (domain.ProductTags, error) {
	tagNames := lo.Map(brandTags, func(tag *domain.ProductTag, _ int) string {
		return tag.Name
	})
	storeTags, err := ds.ProductTagRepo().FindByNamesInStore(ctx, storeID, tagNames)
	if err != nil {
		return nil, err
	}
	if len(storeTags) == len(brandTags) {
		return storeTags, nil
	}

	notExistsTags := make(domain.ProductTags, 0, len(brandTags)-len(storeTags))
	existsTags := make(domain.ProductTags, 0, len(storeTags))
	for _, tag := range brandTags {
		if !lo.ContainsBy(storeTags, func(t *domain.ProductTag) bool {
			return t.Name == tag.Name
		}) {
			notExistsTags = append(notExistsTags, &domain.ProductTag{
				ID:         uuid.New(),
				Name:       tag.Name,
				MerchantID: tag.MerchantID,
				StoreID:    storeID,
			})
		} else {
			existsTags = append(existsTags, tag)
		}
	}

	// 创建门店标签
	if err := ds.ProductTagRepo().CreateBulk(ctx, notExistsTags); err != nil {
		return nil, err
	}
	return append(existsTags, notExistsTags...), nil
}
