package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/product"
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

func (repo *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.Product, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	ep, err := repo.Client.Product.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProductNotExists)
		}
		return nil, err
	}

	res = convertProductToDomain(ep)
	return res, nil
}

func (repo *ProductRepository) Create(ctx context.Context, p *domain.Product) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.Product.Create().
		SetID(p.ID).
		SetName(p.Name).
		SetType(p.Type).
		SetMerchantID(p.MerchantID).
		SetCategoryID(p.CategoryID).
		SetUnitID(p.UnitID).
		SetMnemonic(p.Mnemonic).
		SetShelfLife(p.ShelfLife).
		SetSupportTypes(p.SupportTypes).
		SetSaleStatus(p.SaleStatus).
		SetSaleChannels(p.SaleChannels).
		SetMinSaleQuantity(p.MinSaleQuantity).
		SetAddSaleQuantity(p.AddSaleQuantity).
		SetInheritTaxRate(p.InheritTaxRate).
		SetInheritStall(p.InheritStall).
		SetMainImage(p.MainImage).
		SetDescription(p.Description)

	// 套餐属性（仅套餐商品使用）
	if p.EstimatedCostPrice != nil {
		builder = builder.SetEstimatedCostPrice(*p.EstimatedCostPrice)
	}
	if p.DeliveryCostPrice != nil {
		builder = builder.SetDeliveryCostPrice(*p.DeliveryCostPrice)
	}

	// 可选字段
	if p.StoreID != uuid.Nil {
		builder = builder.SetStoreID(p.StoreID)
	}
	if p.MenuID != uuid.Nil {
		builder = builder.SetMenuID(p.MenuID)
	}
	if p.EffectiveDateType != "" {
		builder = builder.SetEffectiveDateType(p.EffectiveDateType)
	}
	if p.EffectiveStartTime != nil {
		builder = builder.SetEffectiveStartTime(*p.EffectiveStartTime)
	}
	if p.EffectiveEndTime != nil {
		builder = builder.SetEffectiveEndTime(*p.EffectiveEndTime)
	}

	if p.TaxRateID != uuid.Nil {
		builder = builder.SetTaxRateID(p.TaxRateID)
	}
	if p.StallID != uuid.Nil {
		builder = builder.SetStallID(p.StallID)
	}
	if len(p.DetailImages) > 0 {
		builder = builder.SetDetailImages(p.DetailImages)
	}

	// 设置标签（Many2Many）
	if len(p.Tags) > 0 {
		tagIDs := make([]uuid.UUID, 0, len(p.Tags))
		for _, tag := range p.Tags {
			tagIDs = append(tagIDs, tag.ID)
		}
		builder = builder.AddTagIDs(tagIDs...)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return err
	}

	p.ID = created.ID
	p.CreatedAt = created.CreatedAt
	p.UpdatedAt = created.UpdatedAt

	return nil
}

func (repo *ProductRepository) Update(ctx context.Context, p *domain.Product) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.Product.UpdateOneID(p.ID).
		SetName(p.Name).
		SetCategoryID(p.CategoryID).
		SetUnitID(p.UnitID).
		SetMnemonic(p.Mnemonic).
		SetShelfLife(p.ShelfLife).
		SetSupportTypes(p.SupportTypes).
		SetSaleStatus(p.SaleStatus).
		SetSaleChannels(p.SaleChannels).
		SetMinSaleQuantity(p.MinSaleQuantity).
		SetAddSaleQuantity(p.AddSaleQuantity).
		SetInheritTaxRate(p.InheritTaxRate).
		SetInheritStall(p.InheritStall).
		SetMainImage(p.MainImage).
		SetDescription(p.Description)

	// 套餐属性（仅套餐商品使用）
	if p.EstimatedCostPrice != nil {
		builder = builder.SetEstimatedCostPrice(*p.EstimatedCostPrice)
	} else {
		builder = builder.ClearEstimatedCostPrice()
	}
	if p.DeliveryCostPrice != nil {
		builder = builder.SetDeliveryCostPrice(*p.DeliveryCostPrice)
	} else {
		builder = builder.ClearDeliveryCostPrice()
	}

	// 可选字段
	if p.MenuID != uuid.Nil {
		builder = builder.SetMenuID(p.MenuID)
	} else {
		builder = builder.ClearMenuID()
	}
	if p.EffectiveDateType != "" {
		builder = builder.SetEffectiveDateType(p.EffectiveDateType)
	} else {
		builder = builder.ClearEffectiveDateType()
	}
	if p.EffectiveStartTime != nil {
		builder = builder.SetEffectiveStartTime(*p.EffectiveStartTime)
	} else {
		builder = builder.ClearEffectiveStartTime()
	}
	if p.EffectiveEndTime != nil {
		builder = builder.SetEffectiveEndTime(*p.EffectiveEndTime)
	} else {
		builder = builder.ClearEffectiveEndTime()
	}
	if p.TaxRateID != uuid.Nil {
		builder = builder.SetTaxRateID(p.TaxRateID)
	} else {
		builder = builder.ClearTaxRateID()
	}
	if p.StallID != uuid.Nil {
		builder = builder.SetStallID(p.StallID)
	} else {
		builder = builder.ClearStallID()
	}
	if len(p.DetailImages) > 0 {
		builder = builder.SetDetailImages(p.DetailImages)
	} else {
		builder = builder.ClearDetailImages()
	}

	_, err = builder.Save(ctx)
	return err
}

func (repo *ProductRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return repo.Client.Product.DeleteOneID(id).Exec(ctx)
}

func (repo *ProductRepository) Exists(ctx context.Context, params domain.ProductExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Product.Query().
		Where(
			product.MerchantID(params.MerchantID),
			product.Name(params.Name),
		)

	if params.ExcludeID != uuid.Nil {
		query = query.Where(product.IDNEQ(params.ExcludeID))
	}

	exists, err = query.Exist(ctx)
	return exists, err
}

func (repo *ProductRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) (res domain.Products, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductRepository.ListByIDs")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.Product.Query().Where(product.IDIn(ids...))
	entProducts, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	res = make(domain.Products, 0, len(entProducts))
	for _, p := range entProducts {
		res = append(res, convertProductToDomain(p))
	}
	return res, nil
}

// ============================================
// 转换函数
// ============================================

func convertProductToDomain(ep *ent.Product) *domain.Product {
	if ep == nil {
		return nil
	}

	p := &domain.Product{
		ID:                ep.ID,
		Name:              ep.Name,
		MerchantID:        ep.MerchantID,
		StoreID:           ep.StoreID,
		Type:              ep.Type,
		CategoryID:        ep.CategoryID,
		MenuID:            ep.MenuID,
		Mnemonic:          ep.Mnemonic,
		ShelfLife:         ep.ShelfLife,
		SupportTypes:      ep.SupportTypes,
		UnitID:            ep.UnitID,
		SaleStatus:        ep.SaleStatus,
		SaleChannels:      ep.SaleChannels,
		EffectiveDateType: ep.EffectiveDateType,
		MinSaleQuantity:   ep.MinSaleQuantity,
		AddSaleQuantity:   ep.AddSaleQuantity,
		InheritTaxRate:    ep.InheritTaxRate,
		TaxRateID:         ep.TaxRateID,
		InheritStall:      ep.InheritStall,
		StallID:           ep.StallID,
		MainImage:         ep.MainImage,
		Description:       ep.Description,
		CreatedAt:         ep.CreatedAt,
		UpdatedAt:         ep.UpdatedAt,
	}

	// 可选字段
	if ep.EstimatedCostPrice != nil {
		p.EstimatedCostPrice = ep.EstimatedCostPrice
	}
	if ep.DeliveryCostPrice != nil {
		p.DeliveryCostPrice = ep.DeliveryCostPrice
	}
	if ep.EffectiveStartTime != nil {
		p.EffectiveStartTime = ep.EffectiveStartTime
	}
	if ep.EffectiveEndTime != nil {
		p.EffectiveEndTime = ep.EffectiveEndTime
	}
	if len(ep.DetailImages) > 0 {
		p.DetailImages = ep.DetailImages
	}

	// 规格字段
	for _, spec := range ep.Edges.ProductSpecs {
		p.SpecRelations = append(p.SpecRelations, convertProductSpecRelationToDomain(spec))
	}
	for _, attr := range ep.Edges.ProductAttrs {
		p.AttrRelations = append(p.AttrRelations, convertProductAttrRelationToDomain(attr))
	}

	if len(ep.Edges.Tags) > 0 {
		p.Tags = make(domain.ProductTags, 0, len(ep.Edges.Tags))
		for _, tag := range ep.Edges.Tags {
			p.Tags = append(p.Tags, convertProductTagToDomain(tag))
		}
	}

	return p
}
