package product

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *ProductInteractor) Create(ctx context.Context, product *domain.Product) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if err = validateProductParams(product); err != nil {
		return err
	}

	err = i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		if err = validateProductBusinessRules(ctx, ds, product); err != nil {
			return err
		}
		// 创建商品
		err := ds.ProductRepo().Create(ctx, product)
		if err != nil {
			return err
		}

		// 创建规格关联
		if len(product.SpecRelations) > 0 {
			err = ds.ProductSpecRelRepo().CreateBulk(ctx, product.SpecRelations)
			if err != nil {
				return err
			}
		}

		// 创建口味做法关联
		if len(product.AttrRelations) > 0 {
			err = ds.ProductAttrRelRepo().CreateBulk(ctx, product.AttrRelations)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// validateProductParams 校验商品参数格式
func validateProductParams(product *domain.Product) (err error) {
	// 1. 验证规格关联
	if len(product.SpecRelations) == 0 {
		return domain.ParamsError(domain.ErrProductSpecRelationNoDefault)
	}

	defaultSpecCount := 0
	specIDMap := make(map[uuid.UUID]bool)

	for _, specRel := range product.SpecRelations {
		// 验证规格ID不重复
		if specIDMap[specRel.SpecID] {
			return domain.ParamsError(domain.ErrProductSpecInvalid)
		}
		specIDMap[specRel.SpecID] = true

		// 统计默认项数量
		if specRel.IsDefault {
			defaultSpecCount++
		}
	}

	if defaultSpecCount != 1 {
		return domain.ParamsError(domain.ErrProductSpecRelationNoDefault)
	}

	// 2. 验证口味做法关联（如果存在）
	if len(product.AttrRelations) > 0 {
		// 验证口味做法项ID不重复
		attrItemIDMap := make(map[uuid.UUID]bool)
		// 每个属性分组下面的属性值默认值数量
		attrGroupDefaultCountMap := make(map[uuid.UUID]int)

		for _, attrRel := range product.AttrRelations {
			if attrItemIDMap[attrRel.AttrItemID] {
				return domain.ParamsError(domain.ErrProductAttrInvalid)
			}
			attrItemIDMap[attrRel.AttrItemID] = true

			if attrRel.IsDefault {
				attrGroupDefaultCountMap[attrRel.AttrID]++
			}
		}

		// 当点单限制为"必选一项"时，必须设置其中一个为默认项
		for _, count := range attrGroupDefaultCountMap {
			if count != 1 {
				return domain.ParamsError(domain.ErrProductAttrRelationNoDefault)
			}
		}
	}

	// 3. 验证生效日期
	if product.EffectiveDateType == domain.EffectiveDateTypeCustom {
		if product.EffectiveStartTime == nil || product.EffectiveEndTime == nil {
			return domain.ParamsError(domain.ErrProductEffectiveDateInvalid)
		}
		if product.EffectiveStartTime.After(*product.EffectiveEndTime) ||
			product.EffectiveStartTime.Equal(*product.EffectiveEndTime) {
			return domain.ParamsError(domain.ErrProductEffectiveDateInvalid)
		}
	}

	// 4. 验证税率和出品部门
	if !product.InheritTaxRate && product.TaxRateID == uuid.Nil {
		return domain.ParamsError(domain.ErrProductTaxRateNotExists)
	}

	if !product.InheritStall && product.StallID == uuid.Nil {
		return domain.ParamsError(domain.ErrProductStallNotExists)
	}

	return nil
}

// validateProductBusinessRules 校验商品业务规则
func validateProductBusinessRules(ctx context.Context, ds domain.DataStore, product *domain.Product) error {
	// 1. 验证商品名称在当前商户下是否唯一
	exists, err := ds.ProductRepo().Exists(ctx, domain.ProductExistsParams{
		MerchantID: product.MerchantID,
		Name:       product.Name,
	})
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrProductNameExists
	}

	// 2. 验证分类是否存在且有效
	category, err := ds.CategoryRepo().FindByID(ctx, product.CategoryID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrProductCategoryNotExists)
		}
		return err
	}
	// 验证分类是否属于当前商户
	if category.MerchantID != product.MerchantID {
		return domain.ParamsError(domain.ErrProductCategoryInvalid)
	}

	// 3. 验证单位是否存在且有效
	unit, err := ds.ProductUnitRepo().FindByID(ctx, product.UnitID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrProductUnitInvalid)
		}
		return err
	}
	// 验证单位是否属于当前商户
	if unit.MerchantID != product.MerchantID {
		return domain.ParamsError(domain.ErrProductUnitInvalid)
	}

	// 4. 验证规格关联
	if len(product.SpecRelations) > 0 {
		specIDs := lo.Map(product.SpecRelations, func(specRel *domain.ProductSpecRelation, _ int) uuid.UUID {
			return specRel.SpecID
		})
		specs, err := ds.ProductSpecRepo().ListByIDs(ctx, specIDs)
		if err != nil {
			return err
		}
		if len(specs) != len(specIDs) {
			return domain.ParamsError(domain.ErrProductSpecInvalid)
		}
		for _, spec := range specs {
			if spec.MerchantID != product.MerchantID {
				return domain.ParamsError(domain.ErrProductSpecInvalid)
			}
		}
	}

	// 5. 验证口味做法关联（如果存在）
	if len(product.AttrRelations) > 0 {
		attrItemIDMap := make(map[uuid.UUID]uuid.UUID)
		attrItemIDs := make([]uuid.UUID, 0, len(product.AttrRelations))
		for _, attrRel := range product.AttrRelations {
			attrItemIDMap[attrRel.AttrItemID] = attrRel.AttrID
			attrItemIDs = append(attrItemIDs, attrRel.AttrItemID)
		}
		attrItems, err := ds.ProductAttrRepo().ListItemsByIDs(ctx, attrItemIDs)
		if err != nil {
			return err
		}
		if len(attrItems) != len(attrItemIDs) {
			return domain.ParamsError(domain.ErrProductAttrInvalid)
		}
		for _, attrItem := range attrItems {
			attrID, ok := attrItemIDMap[attrItem.ID]
			if !ok || attrItem.AttrID != attrID {
				return domain.ParamsError(domain.ErrProductAttrInvalid)
			}
		}
	}

	// 6. 验证标签有效性（如果存在）
	if len(product.Tags) > 0 {
		tagIDs := lo.Map(product.Tags, func(tag *domain.ProductTag, _ int) uuid.UUID {
			return tag.ID
		})
		tags, err := ds.ProductTagRepo().ListByIDs(ctx, tagIDs)
		if err != nil {
			return err
		}
		if len(tags) != len(tagIDs) {
			return domain.ParamsError(domain.ErrProductTagInvalid)
		}
		for _, tag := range tags {
			if tag.MerchantID != product.MerchantID {
				return domain.ParamsError(domain.ErrProductTagInvalid)
			}
		}
	}

	// 7. 验证税率和出品部门（如果指定）
	if !product.InheritTaxRate && product.TaxRateID != uuid.Nil {
		// @TODO: 需要实现税率配置的验证逻辑
		// 这里暂时跳过，后续需要添加税率配置的 Repository
	}

	if !product.InheritStall && product.StallID != uuid.Nil {
		// @TODO: 需要实现出品部门配置的验证逻辑
		// 这里暂时跳过，后续需要添加出品部门配置的 Repository
	}

	return nil
}
