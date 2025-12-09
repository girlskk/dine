package repository

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	"entgo.io/ent/dialect/sql"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/order"
	"gitlab.jiguang.dev/pos-dine/dine/ent/orderitem"
	"gitlab.jiguang.dev/pos-dine/dine/ent/payment"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"

	"github.com/samber/lo"
)

var _ domain.OrderRepository = (*OrderRepository)(nil)

type OrderRepository struct {
	Client *ent.Client
}

func NewOrderRepository(client *ent.Client) *OrderRepository {
	return &OrderRepository{
		Client: client,
	}
}

func (r *OrderRepository) Create(ctx context.Context, dorder *domain.Order) (newOrder *domain.Order, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if _, err = r.Client.Tx(ctx); err == nil {
		panic("transaction not started")
	}

	// 创建订单
	od, err := r.Client.Order.Create().
		SetNo(dorder.No).
		SetType(dorder.Type).
		SetSource(dorder.Source).
		SetStatus(dorder.Status).
		SetTotalPrice(dorder.TotalPrice).
		SetDiscount(dorder.Discount).
		SetRealPrice(dorder.RealPrice).
		SetPointsAvailable(dorder.PointsAvailable).
		SetMemberID(dorder.MemberID).
		SetMemberName(dorder.MemberName).
		SetMemberPhone(dorder.MemberPhone).
		SetStoreID(dorder.StoreID).
		SetStoreName(dorder.StoreName).
		SetNillableTableID(lo.Ternary(dorder.TableID > 0, &dorder.TableID, nil)).
		SetTableName(dorder.TableName).
		SetPeopleNumber(dorder.PeopleNumber).
		SetCreatorID(dorder.CreatorID).
		SetCreatorName(dorder.CreatorName).
		SetCreatorType(dorder.CreatorType).
		Save(ctx)

	if err != nil {
		err = fmt.Errorf("failed to create order: %w", err)
		return
	}

	span.SetTag("order_id", od.ID)
	span.LogKV("event", "created order")

	// 创建订单商品
	if len(dorder.Items) > 0 {
		var items []*ent.OrderItem
		items, err = r.Client.OrderItem.MapCreateBulk(dorder.Items, func(c *ent.OrderItemCreate, idx int) {
			c.SetOrder(od).
				SetProductID(dorder.Items[idx].ProductID).
				SetName(dorder.Items[idx].Name).
				SetType(int(dorder.Items[idx].Type)).
				SetAllowPointPay(dorder.Items[idx].AllowPointPay).
				SetQuantity(dorder.Items[idx].Quantity).
				SetPrice(dorder.Items[idx].Price).
				SetAmount(dorder.Items[idx].Amount).
				SetProductSnapshot(dorder.Items[idx].ProductSnapshot).
				SetRemark(dorder.Items[idx].Remark)
		}).Save(ctx)
		if err != nil {
			err = fmt.Errorf("failed to create order items: %w", err)
			return
		}

		span.LogKV("event", "created order items")
		od.Edges.Items = items

		// 创建订单商品套餐详情
		for i, ditem := range dorder.Items {
			if ditem.Type == domain.ProductTypeSetMeal && len(ditem.SetMealDetails) > 0 {
				var details []*ent.OrderItemSetMealDetail
				details, err = r.Client.OrderItemSetMealDetail.MapCreateBulk(ditem.SetMealDetails, func(c *ent.OrderItemSetMealDetailCreate, idx int) {
					c.SetOrderItem(items[i]).
						SetName(ditem.SetMealDetails[idx].Name).
						SetType(int(ditem.SetMealDetails[idx].Type)).
						SetSetMealPrice(ditem.SetMealDetails[idx].SetMealPrice).
						SetSetMealID(ditem.SetMealDetails[idx].SetMealID).
						SetProductID(ditem.SetMealDetails[idx].ProductID).
						SetQuantity(ditem.SetMealDetails[idx].Quantity).
						SetProductSnapshot(ditem.SetMealDetails[idx].ProductSnapshot)
				}).Save(ctx)
				if err != nil {
					err = fmt.Errorf("failed to create order item set meal details: %w", err)
					return
				}

				span.LogKV("event", "created order item set meal details")
				items[i].Edges.SetMealDetails = details
			}
		}
	}

	newOrder = convertOrder(od)

	return
}

func (r *OrderRepository) GetOrders(ctx context.Context, pager *upagination.Pagination, filter *domain.OrderListFilter, orderBys ...domain.OrderListOrder) (dorders []*domain.Order, total int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.GetOrders")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := r.filterBuildQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count: %w", err)
		return
	}

	orders, err := query.Order(r.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get orders: %w", err)
	}

	dorders = lo.Map(orders, func(order *ent.Order, _ int) *domain.Order {
		return convertOrder(order)
	})

	return
}

func (r *OrderRepository) GetOrdersWithItems(ctx context.Context, pager *upagination.Pagination, filter *domain.OrderListFilter, orderBys ...domain.OrderListOrder) (dorders []*domain.Order, total int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.GetOrdersWithItems")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := r.filterBuildQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count: %w", err)
		return
	}

	orders, err := query.Order(r.orderBy(orderBys...)...).
		WithItems().
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get orders: %w", err)
	}

	dorders = lo.Map(orders, func(order *ent.Order, _ int) *domain.Order {
		return convertOrder(order)
	})

	return
}

func (r *OrderRepository) Find(ctx context.Context, id int) (dorder *domain.Order, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.Find")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	od, err := r.Client.Order.Query().
		Where(order.ID(id)).
		WithItems(func(q *ent.OrderItemQuery) {
			q.WithSetMealDetails()
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get order by no: %w", err)
	}

	dorder = convertOrder(od)

	return
}

func (r *OrderRepository) FindByNo(ctx context.Context, no string) (dorder *domain.Order, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.FindByNo")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	od, err := r.Client.Order.Query().
		Where(order.No(no)).
		WithItems(func(q *ent.OrderItemQuery) {
			q.WithSetMealDetails()
		}).
		WithLogs().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get order by no: %w", err)
	}

	dorder = convertOrder(od)

	return
}

func (r *OrderRepository) FindByItemID(ctx context.Context, itemID int) (dorder *domain.Order, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.FindByItemID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	od, err := r.Client.Order.Query().
		Where(order.HasItemsWith(orderitem.ID(itemID))).
		WithItems(func(q *ent.OrderItemQuery) {
			q.WithSetMealDetails()
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get order by item id: %w", err)
	}

	dorder = convertOrder(od)

	return
}

func (r *OrderRepository) CreateLog(ctx context.Context, dlog *domain.OrderLog) (newLog *domain.OrderLog, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.CreateLog")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	log, err := r.Client.OrderLog.Create().
		SetOrderID(dlog.OrderID).
		SetEvent(dlog.Event).
		SetOperatorType(dlog.OperatorType).
		SetOperatorID(dlog.OperatorID).
		SetOperatorName(dlog.OperatorName).
		Save(ctx)

	if err != nil {
		err = fmt.Errorf("failed to create order log: %w", err)
		return
	}

	newLog = convertOrderLog(log)

	return
}

func (r *OrderRepository) Update(ctx context.Context, dorder *domain.Order) (updatedOrder *domain.Order, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	od, err := r.Client.Order.UpdateOneID(dorder.ID).
		SetNillableTableID(lo.Ternary(dorder.TableID > 0, &dorder.TableID, nil)).
		SetTableName(dorder.TableName).
		SetStatus(dorder.Status).
		SetTotalPrice(dorder.TotalPrice).
		SetDiscount(dorder.Discount).
		SetRealPrice(dorder.RealPrice).
		SetPointsAvailable(dorder.PointsAvailable).
		SetMemberID(dorder.MemberID).
		SetMemberName(dorder.MemberName).
		SetMemberPhone(dorder.MemberPhone).
		SetPeopleNumber(dorder.PeopleNumber).
		SetPaid(dorder.Paid).
		SetRefunded(dorder.Refunded).
		SetPaidChannels(dorder.PaidChannels).
		SetCashPaid(dorder.CashPaid).
		SetWechatPaid(dorder.WechatPaid).
		SetWechatRefunded(dorder.WechatRefunded).
		SetAlipayPaid(dorder.AlipayPaid).
		SetAlipayRefunded(dorder.AlipayRefunded).
		SetPointsPaid(dorder.PointsPaid).
		SetPointsRefunded(dorder.PointsRefunded).
		SetPointsWalletPaid(dorder.PointsWalletPaid).
		SetPointsWalletRefunded(dorder.PointsWalletRefunded).
		SetNillableLastPaidAt(dorder.LastPaidAt).
		SetNillableFinishedAt(dorder.FinishedAt).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to update order: %w", err)
		return
	}

	updatedOrder = convertOrder(od)

	return
}

func (r *OrderRepository) UpdateItem(ctx context.Context, ditem *domain.OrderItem) (updatedItem *domain.OrderItem, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.UpdateItem")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	item, err := r.Client.OrderItem.UpdateOneID(ditem.ID).
		SetQuantity(ditem.Quantity).
		SetPrice(ditem.Price).
		SetAmount(ditem.Amount).
		SetRemark(ditem.Remark).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to update order item: %w", err)
		return
	}

	updatedItem = convertOrderItem(item)

	return
}

func (r *OrderRepository) RemoveItems(ctx context.Context, orderID int, itemIDs ...int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.RemoveItems")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = r.Client.OrderItem.Delete().
		Where(orderitem.OrderID(orderID)).
		Where(orderitem.IDIn(itemIDs...)).
		Exec(ctx)
	if err != nil {
		err = fmt.Errorf("failed to remove order items: %w", err)
		return
	}

	return
}

func (r *OrderRepository) AppendItems(ctx context.Context, orderID int, ditems []*domain.OrderItem) (newItems []*domain.OrderItem, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.AppendItems")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	items, err := r.Client.OrderItem.MapCreateBulk(ditems, func(c *ent.OrderItemCreate, idx int) {
		c.SetOrderID(orderID).
			SetProductID(ditems[idx].ProductID).
			SetName(ditems[idx].Name).
			SetType(int(ditems[idx].Type)).
			SetAllowPointPay(ditems[idx].AllowPointPay).
			SetQuantity(ditems[idx].Quantity).
			SetPrice(ditems[idx].Price).
			SetAmount(ditems[idx].Amount).
			SetProductSnapshot(ditems[idx].ProductSnapshot).
			SetRemark(ditems[idx].Remark)
	}).Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to append order items: %w", err)
		return
	}
	span.LogKV("event", "created order items")

	for i, ditem := range ditems {
		if ditem.Type == domain.ProductTypeSetMeal && len(ditem.SetMealDetails) > 0 {
			var details []*ent.OrderItemSetMealDetail
			details, err = r.Client.OrderItemSetMealDetail.MapCreateBulk(ditem.SetMealDetails, func(c *ent.OrderItemSetMealDetailCreate, idx int) {
				c.SetOrderItem(items[i]).
					SetName(ditem.SetMealDetails[idx].Name).
					SetType(int(ditem.SetMealDetails[idx].Type)).
					SetSetMealPrice(ditem.SetMealDetails[idx].SetMealPrice).
					SetSetMealID(ditem.SetMealDetails[idx].SetMealID).
					SetProductID(ditem.SetMealDetails[idx].ProductID).
					SetQuantity(ditem.SetMealDetails[idx].Quantity).
					SetProductSnapshot(ditem.SetMealDetails[idx].ProductSnapshot)
			}).Save(ctx)
			if err != nil {
				err = fmt.Errorf("failed to create order item set meal details: %w", err)
				return
			}

			span.LogKV("event", "created order item set meal details")
			items[i].Edges.SetMealDetails = details
		}
	}

	newItems = convertOrderItems(items)

	return
}

func (r *OrderRepository) CreateFinanceLog(ctx context.Context, dlog *domain.OrderFinanceLog) (newLog *domain.OrderFinanceLog, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.CreateFinanceLog")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	log, err := r.Client.OrderFinanceLog.Create().
		SetOrderID(dlog.OrderID).
		SetAmount(dlog.Amount).
		SetType(dlog.Type).
		SetChannel(dlog.Channel).
		SetSeqNo(dlog.SeqNo).
		SetCreatorType(dlog.CreatorType).
		SetCreatorID(dlog.CreatorID).
		SetCreatorName(dlog.CreatorName).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create order finance log: %w", err)
		return
	}

	newLog = convertOrderFinanceLog(log)

	return
}

func (r *OrderRepository) HasIncompletePayment(ctx context.Context, orderID int) (has bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.HasIncompletePayment")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	has, err = r.Client.Payment.Query().
		Where(
			payment.PayBizTypeEQ(domain.PayBizTypeOrder),
			payment.BizID(orderID),
			payment.FinishedAtIsNil(),
		).Exist(ctx)
	if err != nil {
		err = fmt.Errorf("failed to check incomplete payment: %w", err)
		return
	}

	return
}

func (r *OrderRepository) ListItemNamesByOrders(ctx context.Context, orderIDs []int) (res map[int][]string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.ListItemNamesByOrders")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 查询 order_id 和 name
	items, err := r.Client.OrderItem.
		Query().
		Where(orderitem.OrderIDIn(orderIDs...)).
		Select(orderitem.FieldOrderID, orderitem.FieldName).
		All(ctx)
	if err != nil {
		return nil, err
	}

	res = make(map[int][]string)
	for _, item := range items {
		res[item.OrderID] = append(res[item.OrderID], item.Name)
	}

	return res, nil
}

func (r *OrderRepository) GetOrderRange(ctx context.Context, filter *domain.OrderListFilter) (rg domain.OrderRange, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderRepository.GetOrderRange")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := r.filterBuildQuery(filter)

	var v []struct {
		Min, Max, Count int
	}

	err = query.Aggregate(
		ent.Min(order.FieldID),
		ent.Max(order.FieldID),
		ent.Count(),
	).Scan(ctx, &v)

	if err != nil {
		err = fmt.Errorf("failed to get max order id: %w", err)
		return
	}

	if len(v) > 0 {
		rg = domain.OrderRange{
			MinID: v[0].Min,
			MaxID: v[0].Max,
			Count: v[0].Count,
		}
	}

	return
}

func (r *OrderRepository) filterBuildQuery(filter *domain.OrderListFilter) *ent.OrderQuery {
	query := r.Client.Order.Query()

	if filter.StoreID > 0 {
		query = query.Where(order.StoreID(filter.StoreID))
	}
	if filter.Status != "" {
		query = query.Where(order.StatusEQ(filter.Status))
	}

	if filter.FinishedAtGte != nil {
		query = query.Where(order.FinishedAtGTE(*filter.FinishedAtGte))
	}
	if filter.FinishedAtLte != nil {
		query = query.Where(order.FinishedAtLTE(*filter.FinishedAtLte))
	}

	if filter.HasItemName != "" {
		query = query.Where(order.HasItemsWith(orderitem.NameContains(filter.HasItemName)))
	}
	if filter.MemberNameOrPhone != "" {
		query = query.Where(
			order.Or(
				order.MemberNameContains(filter.MemberNameOrPhone),
				order.MemberPhoneContains(filter.MemberNameOrPhone),
			),
		)
	}
	if filter.CreatedAtGte != nil {
		query = query.Where(order.CreatedAtGTE(*filter.CreatedAtGte))
	}
	if filter.CreatedAtLte != nil {
		query = query.Where(order.CreatedAtLTE(*filter.CreatedAtLte))
	}

	if filter.PointsPaidGt0 {
		query = query.Where(order.PointsPaidGT(decimal.Zero))
	}

	if filter.IDGte > 0 {
		query = query.Where(order.IDGTE(filter.IDGte))
	}
	if filter.IDLte > 0 {
		query = query.Where(order.IDLTE(filter.IDLte))
	}

	if filter.CreatorID > 0 {
		query = query.Where(order.CreatorID(filter.CreatorID))
	}

	if filter.CreatorType != "" {
		query = query.Where(order.CreatorTypeEQ(filter.CreatorType))
	}

	return query
}

func (r *OrderRepository) orderBy(orderBys ...domain.OrderListOrder) []order.OrderOption {
	var opts []order.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.OrderListOrderByID:
			opts = append(opts, order.ByID(rule))
		case domain.OrderListOrderByCreatedAt:
			opts = append(opts, order.ByCreatedAt(rule))
		}
	}

	if len(opts) == 0 {
		opts = append(opts, order.ByCreatedAt(sql.OrderDesc()))
	}

	return opts
}
