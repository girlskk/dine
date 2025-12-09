package repository

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"github.com/shopspring/decimal"

	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/ordercart"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.OrderCartRepository = (*OrderCartRepository)(nil)

type OrderCartRepository struct {
	Client *ent.Client
}

func NewOrderCartRepository(client *ent.Client) *OrderCartRepository {
	return &OrderCartRepository{
		Client: client,
	}
}

func (r *OrderCartRepository) ListByTable(ctx context.Context, tableID int, withExtra bool) (res domain.OrderCarts, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderCartRepository.ListByTable")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	// 查询购物车项，并关联商品、规格、属性和做法信息
	query := r.Client.OrderCart.Query().
		Where(
			ordercart.TableID(tableID),
		)

	if withExtra {
		query = query.WithProduct(func(q *ent.ProductQuery) {
			q.WithSetMealDetails(func(q *ent.SetMealDetailQuery) {
				q.WithProduct()
			}).
				WithCategory()
		}).
			WithProductSpec(func(q *ent.ProductSpecQuery) {
				q.WithSpec()
			}).
			WithAttr().
			WithRecipe()
	}

	entItems, err := query.
		Order(ent.Desc(ordercart.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		return nil, err
	}
	// 使用DTO转换函数将Ent实体转换为领域对象
	return convertOrderCarts(entItems), nil
}

// 根据唯一键查找购物车项
func (r *OrderCartRepository) FindByUniqueKey(ctx context.Context,
	key domain.OrderCartItemUniqueKey,
) (item *domain.OrderCart, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderCartRepository.FindByUniqueKey")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := r.Client.OrderCart.Query().
		Where(
			ordercart.TableID(key.TableID),
			ordercart.ProductID(key.ProductID),
		)

	if key.ProductSpecID > 0 {
		query = query.Where(ordercart.ProductSpecID(key.ProductSpecID))
	} else {
		query = query.Where(ordercart.ProductSpecIDIsNil())
	}
	if key.AttrID > 0 {
		query = query.Where(ordercart.AttrID(key.AttrID))
	} else {
		query = query.Where(ordercart.AttrIDIsNil())
	}
	if key.RecipeID > 0 {
		query = query.Where(ordercart.RecipeID(key.RecipeID))
	} else {
		query = query.Where(ordercart.RecipeIDIsNil())
	}

	res, err := query.First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrOrderCartNotFound)
		}
		return nil, err
	}
	return convertOrderCart(res), nil
}

func (r *OrderCartRepository) Create(ctx context.Context, item *domain.OrderCart) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderCartRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	sqlBuilder := r.Client.OrderCart.Create().
		SetTableID(item.TableID).
		SetProductID(item.ProductID).
		SetQuantity(item.Quantity)

	if item.ProductSpecID > 0 {
		sqlBuilder.SetProductSpecID(item.ProductSpecID)
	}
	if item.AttrID > 0 {
		sqlBuilder.SetAttrID(item.AttrID)
	}
	if item.RecipeID > 0 {
		sqlBuilder.SetRecipeID(item.RecipeID)
	}

	_, err = sqlBuilder.Save(ctx)

	return err
}

func (r *OrderCartRepository) FindByID(ctx context.Context, id int) (item *domain.OrderCart, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderCartRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	res, err := r.Client.OrderCart.Query().
		Where(ordercart.ID(id)).
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrOrderCartNotFound)
		}
		return nil, err
	}
	return convertOrderCart(res), nil
}

// DecrementQuantity 原子减少数量，如果数量为1，则删除记录
func (r *OrderCartRepository) DecrementQuantity(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderCartRepository.DecrementQuantity")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 尝试删除数量等于1的记录
	result, err := r.Client.OrderCart.Delete().
		Where(
			ordercart.ID(id),
			ordercart.QuantityEQ(decimal.NewFromInt(1)),
		).
		Exec(ctx)

	if err != nil {
		return err
	}

	// 如果删除成功（有记录被删除），则返回
	if result > 0 {
		return nil
	}

	// 如果没有记录被删除，尝试更新数量大于1的记录
	_, err = r.Client.OrderCart.UpdateOneID(id).
		Where(ordercart.QuantityGT(decimal.NewFromInt(1))). // 只更新数量大于1的记录
		Modify(func(u *sql.UpdateBuilder) {
			u.Set(ordercart.FieldQuantity, sql.ExprFunc(func(b *sql.Builder) {
				b.Ident(ordercart.FieldQuantity).WriteOp(sql.OpSub).Arg(decimal.NewFromInt(1))
			}))
		}).
		Save(ctx)

	return err
}

// IncrementQuantity 原子增加数量
func (r *OrderCartRepository) IncrementQuantity(ctx context.Context, id int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderCartRepository.IncrementQuantity")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 使用SQL修改器来原子增加quantity
	_, err = r.Client.OrderCart.UpdateOneID(id).
		Modify(func(u *sql.UpdateBuilder) {
			u.Set(ordercart.FieldQuantity, sql.ExprFunc(func(b *sql.Builder) {
				b.Ident(ordercart.FieldQuantity).WriteOp(sql.OpAdd).Arg(decimal.NewFromInt(1))
			}))
		}).
		Save(ctx)

	return err
}

func (r *OrderCartRepository) ClearByTable(ctx context.Context, tableID int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderCartRepository.ClearByTable")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = r.Client.OrderCart.Delete().
		Where(ordercart.TableID(tableID)).
		Exec(ctx)

	return err
}
