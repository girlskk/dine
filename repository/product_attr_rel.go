package repository

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductAttrRelRepository = (*ProductAttrRelRepository)(nil)

type ProductAttrRelRepository struct {
	Client *ent.Client
}

func NewProductAttrRelRepository(client *ent.Client) *ProductAttrRelRepository {
	return &ProductAttrRelRepository{
		Client: client,
	}
}

// CreateBulk 批量创建商品口味做法关联
func (repo *ProductAttrRelRepository) CreateBulk(ctx context.Context, relations domain.ProductAttrRelations) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProductAttrRelRepository.CreateBulk")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(relations) == 0 {
		return nil
	}

	builders := make([]*ent.ProductAttrRelationCreate, 0, len(relations))
	for _, relation := range relations {
		builder := repo.Client.ProductAttrRelation.Create().
			SetID(relation.ID).
			SetProductID(relation.ProductID).
			SetAttrID(relation.AttrID).
			SetAttrItemID(relation.AttrItemID).
			SetIsDefault(relation.IsDefault)

		builders = append(builders, builder)
	}

	_, err = repo.Client.ProductAttrRelation.CreateBulk(builders...).Save(ctx)
	return err
}

func convertProductAttrRelationToDomain(er *ent.ProductAttrRelation) *domain.ProductAttrRelation {
	if er == nil {
		return nil
	}

	return &domain.ProductAttrRelation{
		ID:         er.ID,
		ProductID:  er.ProductID,
		AttrID:     er.AttrID,
		AttrItemID: er.AttrItemID,
		IsDefault:  er.IsDefault,
		CreatedAt:  er.CreatedAt,
		UpdatedAt:  er.UpdatedAt,
	}
}
