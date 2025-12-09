package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/product"
	"gitlab.jiguang.dev/pos-dine/dine/ent/setmealdetail"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.SetMealDetailRepository = (*SetMealDetailRepository)(nil)

type SetMealDetailRepository struct {
	Client *ent.Client
}

func NewSetMealDetailRepository(client *ent.Client) *SetMealDetailRepository {
	return &SetMealDetailRepository{
		Client: client,
	}
}

func (repo *SetMealDetailRepository) BatchCreate(ctx context.Context, details domain.SetMealDetails) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SetMealDetailRepository.BatchCreate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.SetMealDetail.MapCreateBulk(details, func(c *ent.SetMealDetailCreate, idx int) {
		c.SetProductID(details[idx].ProductID).
			SetSetMealID(details[idx].SetMealID).
			SetQuantity(details[idx].Quantity).
			SetPrice(details[idx].SetMealPrice)
		if details[idx].Spec.ID > 0 {
			c.SetProductSpecID(details[idx].Spec.ID)
		}
	}).Save(ctx)
	return err
}

func (repo *SetMealDetailRepository) DeleteBySetMealID(ctx context.Context, setMealID int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SetMealDetailRepository.DeleteBySetMealID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.SetMealDetail.Delete().
		Where(
			setmealdetail.HasSetMealWith(
				product.ID(setMealID), // 通过 edge 关系定位
			),
		).
		Exec(ctx)
	return err
}

func (repo *SetMealDetailRepository) ListBySetMealID(ctx context.Context, id int) (details domain.SetMealDetails, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SetMealDetailRepository.ListBySetMealID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	res, err := repo.Client.SetMealDetail.Query().
		Where(
			setmealdetail.SetMealID(id),
		).
		WithProduct().
		All(ctx)

	if err != nil {
		return nil, err
	}

	for _, item := range res {
		details = append(details, repo.convertToDomain(item))
	}
	return details, nil
}

func (repo *SetMealDetailRepository) convertToDomain(d *ent.SetMealDetail) *domain.SetMealDetail {
	if d == nil {
		return nil
	}

	item := &domain.SetMealDetail{
		ID:           d.ID,
		ProductID:    d.ProductID,
		Quantity:     d.Quantity,
		SetMealID:    d.SetMealID,
		SetMealPrice: d.Price,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}

	if d.Edges.Product != nil {
		item.UnitID = d.Edges.Product.UnitID
		item.CategoryID = d.Edges.Product.CategoryID
		item.Name = d.Edges.Product.Name
		item.Price = d.Edges.Product.Price
	}

	return item
}
