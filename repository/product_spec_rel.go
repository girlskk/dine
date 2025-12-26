package repository

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductSpecRelRepository = (*ProductSpecRelRepository)(nil)

type ProductSpecRelRepository struct {
	Client *ent.Client
}

func NewProductSpecRelRepository(client *ent.Client) *ProductSpecRelRepository {
	return &ProductSpecRelRepository{
		Client: client,
	}
}

// CreateBulk 批量创建商品规格关联
func (repo *ProductSpecRelRepository) CreateBulk(ctx context.Context, relations domain.ProductSpecRelations) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductSpecRelRepository.CreateBulk")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(relations) == 0 {
		return nil
	}

	builders := make([]*ent.ProductSpecRelationCreate, 0, len(relations))
	for _, relation := range relations {
		builder := repo.Client.ProductSpecRelation.Create().
			SetID(relation.ID).
			SetProductID(relation.ProductID).
			SetSpecID(relation.SpecID).
			SetBasePrice(relation.BasePrice).
			SetPackingFeeID(relation.PackingFeeID).
			SetBarcode(relation.Barcode).
			SetIsDefault(relation.IsDefault)

		if relation.MemberPrice != nil {
			builder = builder.SetMemberPrice(*relation.MemberPrice)
		}
		if relation.EstimatedCostPrice != nil {
			builder = builder.SetEstimatedCostPrice(*relation.EstimatedCostPrice)
		}
		if relation.OtherPrice1 != nil {
			builder = builder.SetOtherPrice1(*relation.OtherPrice1)
		}
		if relation.OtherPrice2 != nil {
			builder = builder.SetOtherPrice2(*relation.OtherPrice2)
		}
		if relation.OtherPrice3 != nil {
			builder = builder.SetOtherPrice3(*relation.OtherPrice3)
		}

		builders = append(builders, builder)
	}

	_, err = repo.Client.ProductSpecRelation.CreateBulk(builders...).Save(ctx)
	return err
}

func convertProductSpecRelationToDomain(er *ent.ProductSpecRelation) *domain.ProductSpecRelation {
	if er == nil {
		return nil
	}

	relation := &domain.ProductSpecRelation{
		ID:           er.ID,
		ProductID:    er.ProductID,
		SpecID:       er.SpecID,
		BasePrice:    er.BasePrice,
		PackingFeeID: er.PackingFeeID,
		Barcode:      er.Barcode,
		IsDefault:    er.IsDefault,
		CreatedAt:    er.CreatedAt,
		UpdatedAt:    er.UpdatedAt,
	}

	if er.MemberPrice != nil {
		relation.MemberPrice = er.MemberPrice
	}
	if er.EstimatedCostPrice != nil {
		relation.EstimatedCostPrice = er.EstimatedCostPrice
	}
	if er.OtherPrice1 != nil {
		relation.OtherPrice1 = er.OtherPrice1
	}
	if er.OtherPrice2 != nil {
		relation.OtherPrice2 = er.OtherPrice2
	}
	if er.OtherPrice3 != nil {
		relation.OtherPrice3 = er.OtherPrice3
	}

	return relation
}
