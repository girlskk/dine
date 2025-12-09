package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/domain/order"
	"gitlab.jiguang.dev/pos-dine/dine/domain/payment"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/huifu"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/zxh"
)

var _ domain.OrderInteractor = (*OrderInteractor)(nil)

type OrderInteractor struct {
	DataStore            domain.DataStore
	DailySequence        domain.DailySequence
	OrderEventTrigger    domain.OrderEventTrigger
	MutexManager         domain.MutexManager
	OrderDomainService   *order.DomainService
	bsPay                *huifu.BsPay
	PaymentDomainService *payment.DomainService
	zxhManager           *zxh.Manager
}

func NewOrderInteractor(
	dataStore domain.DataStore,
	seq domain.DailySequence,
	eventTrigger domain.OrderEventTrigger,
	mutexManager domain.MutexManager,
	bsPay *huifu.BsPay,
	paymentDomainService *payment.DomainService,
	zxhManager *zxh.Manager,
	orderDomainService *order.DomainService,
) *OrderInteractor {
	return &OrderInteractor{
		DataStore:            dataStore,
		DailySequence:        seq,
		OrderEventTrigger:    eventTrigger,
		MutexManager:         mutexManager,
		bsPay:                bsPay,
		PaymentDomainService: paymentDomainService,
		zxhManager:           zxhManager,
		OrderDomainService:   orderDomainService,
	}
}

// 创建订单
func (interactor *OrderInteractor) CreateOrder(ctx context.Context, params *domain.CreateOrderParams) (order *domain.Order, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.CreateOrder")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	store, table, creator := params.Store, params.Table, params.Creator

	// 餐厅类型订单必须有台桌
	if store.Type == domain.StoreTypeRestaurant && table == nil {
		return nil, domain.ParamsErrorf("餐厅类型订单必须选择台桌")
	}

	if table != nil && table.StoreID != store.ID {
		return nil, domain.ParamsErrorf("台桌 %d 不存在", table.ID)
	}

	var tableID int
	var tableName string
	if table != nil {
		tableID = table.ID
		tableName = table.Name
	}

	// 取出所有商品ID
	productIDs := lo.Map(params.Items, func(item *domain.CreateOrderItem, _ int) int {
		return item.ProductID
	})

	// 获取所有商品信息
	products, err := interactor.DataStore.ProductRepo().GetDetailsByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	productsMap := lo.KeyBy(products, func(item *domain.Product) int { return item.ID })

	// 构造订单商品结构 domain.OrderItem
	orderItems := make(domain.OrderItems, 0, len(params.Items))
	for _, item := range params.Items {
		product, ok := productsMap[item.ProductID]
		if !ok {
			return nil, domain.ParamsErrorf("商品 %d 不存在", item.ProductID)
		}

		if product.StoreID != store.ID {
			return nil, domain.ParamsErrorf("商品 %d 不存在", item.ProductID)
		}

		if !product.CheckSale() {
			return nil, domain.ParamsErrorf("商品 %d 不可售", item.ProductID)
		}

		orderItem, err := interactor.buildOrderItem(item, product)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, orderItem)
	}

	// 生成订单号
	i, err := interactor.DailySequence.Next(ctx, domain.DailySequencePrefixOrderNo)
	if err != nil {
		return nil, fmt.Errorf("failed to generate order number: %w", err)
	}
	// YYMMDDSSSNNNNN (SSS 是门店ID，NNNNN 是自增的订单号)
	orderNo := fmt.Sprintf("%s%03d%05d", time.Now().Format("060102"), store.ID, i)

	totalAmount := orderItems.TotalAmount()
	// 构造订单
	source := lo.Ternary(params.Source == "", domain.OrderSourceOffline, params.Source)
	order = &domain.Order{
		No:              orderNo,
		Type:            domain.OrderTypeDineIn,
		Source:          source,
		Status:          domain.OrderStatusUnpaid,
		TotalPrice:      totalAmount,
		PointsAvailable: orderItems.PointsAvailable(store.ID),
		RealPrice:       totalAmount,
		StoreID:         store.ID,
		StoreName:       store.Name,
		TableID:         tableID,
		TableName:       tableName,
		PeopleNumber:    params.PeopleNumber,
		CreatorID:       creator.GetOperatorID(),
		CreatorName:     creator.GetOperatorName(),
		CreatorType:     creator.GetOperatorType(),
		Items:           orderItems,
	}

	// 保存订单
	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		var err error
		order, err = ds.OrderRepo().Create(ctx, order)
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// 触发订单创建事件
		err = interactor.OrderEventTrigger.FireCreateOrder(ctx, &domain.OrderEventBaseParams{
			DataStore: ds,
			Order:     order,
			Operator:  creator,
		})
		if err != nil {
			return fmt.Errorf("failed to trigger create order event: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return
}

// CreateOrderFromCart 从购物车创建订单
func (interactor *OrderInteractor) CreateOrderFromCart(ctx context.Context, params *domain.CreateOrderParams) (order *domain.Order, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.CreateOrderFromCart")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	// 锁购物车
	unlock, err := interactor.lockCart(ctx, params.Table.ID)
	if err != nil {
		return
	}
	defer unlock()

	// 获取购物车内容
	cartItems, err := interactor.DataStore.OrderCartRepo().ListByTable(ctx, params.Table.ID, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	if len(cartItems) == 0 {
		return nil, domain.ParamsErrorf("购物车为空")
	}

	// 构造订单商品结构
	orderItems := make([]*domain.CreateOrderItem, 0, len(cartItems))
	for _, item := range cartItems {
		orderItems = append(orderItems, &domain.CreateOrderItem{
			ProductID:       item.ProductID,
			ProductSpecID:   item.ProductSpecID,
			ProductAttrID:   item.AttrID,
			ProductRecipeID: item.RecipeID,
			Quantity:        item.Quantity,
			Price:           item.Price,
			Remark:          "",
		})
	}

	params.Items = orderItems
	// 创建订单
	order, err = interactor.CreateOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	err = interactor.DataStore.OrderCartRepo().ClearByTable(ctx, params.Table.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}
	return order, nil
}

// 锁购物车
func (interactor *OrderInteractor) lockCart(ctx context.Context, tableID int) (unlock func(), err error) {
	mu := interactor.MutexManager.NewMutex(domain.NewMutexCartKey(tableID))
	if err = mu.Lock(ctx); err != nil {
		if domain.IsAlreadyTakenError(err) {
			err = domain.ParamsErrorf("购物车 %d 正在被其他操作修改", tableID)
		}
		err = fmt.Errorf("failed to lock cart: %w", err)
		return
	}
	unlock = func() {
		if _, err := mu.Unlock(ctx); err != nil {
			logger := logging.FromContext(ctx).Named("OrderInteractor.lockCart")
			logger.Errorf("failed to unlock cart: %s", err)
		}
	}
	return
}

// 获取订单
func (interactor *OrderInteractor) GetOrder(ctx context.Context, no string) (order *domain.Order, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.GetOrder")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	order, err = interactor.DataStore.OrderRepo().FindByNo(ctx, no)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsErrorf("订单 %s 不存在", no)
			return
		}
		err = fmt.Errorf("failed to find order: %w", err)
		return
	}

	return
}

// 构造订单商品
func (interactor *OrderInteractor) buildOrderItem(item *domain.CreateOrderItem, product *domain.Product) (*domain.OrderItem, error) {
	specMap := lo.KeyBy(product.Specs, func(spec *domain.ProductSpecRel) int { return spec.ID })
	attrMap := lo.KeyBy(product.Attrs, func(attr *domain.ProductAttr) int { return attr.ID })
	recipeMap := lo.KeyBy(product.Recipes, func(recipe *domain.ProductRecipe) int { return recipe.ID })

	var curSpec *domain.ProductSpecRel
	var curAttr *domain.ProductAttr
	var curRecipe *domain.ProductRecipe

	var ok bool
	if product.Type != domain.ProductTypeSetMeal {
		if item.ProductSpecID != 0 {
			curSpec, ok = specMap[item.ProductSpecID]
			if !ok {
				return nil, domain.ParamsErrorf("规格 %d 不存在", item.ProductSpecID)
			}

			if curSpec.SaleStatus == domain.ProductSaleStatusOff {
				return nil, domain.ParamsErrorf("规格 %d 售罄", item.ProductSpecID)
			}
		}

		if item.ProductAttrID != 0 {
			if curAttr, ok = attrMap[item.ProductAttrID]; !ok {
				return nil, domain.ParamsErrorf("属性 %d 不存在", item.ProductAttrID)
			}
		}

		if item.ProductRecipeID != 0 {
			if curRecipe, ok = recipeMap[item.ProductRecipeID]; !ok {
				return nil, domain.ParamsErrorf("做法 %d 不存在", item.ProductRecipeID)
			}
		}
	}

	snapshot := domain.OrderProductInfoSnapshot{
		Price:  product.Price,
		Images: product.Images,
		UnitID: product.UnitID,
	}

	if product.Unit != nil {
		snapshot.UnitName = product.Unit.Name
	}

	if curSpec != nil {
		snapshot.SpecID = curSpec.ID
		snapshot.SpecName = curSpec.SpecName
		snapshot.SpecPrice = curSpec.Price
	}

	if curAttr != nil {
		snapshot.AttrID = curAttr.ID
		snapshot.AttrName = curAttr.Name
	}

	if curRecipe != nil {
		snapshot.RecipeID = curRecipe.ID
		snapshot.RecipeName = curRecipe.Name
	}

	orderItem := &domain.OrderItem{
		ProductID:       item.ProductID,
		Name:            product.Name,
		Type:            product.Type,
		AllowPointPay:   product.AllowPointPay,
		Quantity:        item.Quantity,
		Price:           item.Price,
		Amount:          item.Price.Mul(item.Quantity),
		ProductSnapshot: snapshot,
		Remark:          item.Remark,
	}

	for _, setMealDetail := range product.SetMealDetails {
		snapshot := domain.OrderProductInfoSnapshot{
			Price:    setMealDetail.Price,
			Images:   setMealDetail.Images,
			UnitID:   setMealDetail.UnitID,
			UnitName: setMealDetail.Unit.Name,
		}

		if setMealDetail.Spec != nil {
			snapshot.SpecID = setMealDetail.Spec.ID
			snapshot.SpecName = setMealDetail.Spec.SpecName
			snapshot.SpecPrice = setMealDetail.Spec.Price
		}

		orderItem.SetMealDetails = append(orderItem.SetMealDetails, &domain.OrderItemSetMealDetail{
			Name:            setMealDetail.Name,
			Type:            setMealDetail.ProductType,
			SetMealPrice:    setMealDetail.SetMealPrice,
			SetMealID:       setMealDetail.SetMealID,
			ProductID:       setMealDetail.ProductID,
			Quantity:        setMealDetail.Quantity,
			ProductSnapshot: snapshot,
		})
	}

	return orderItem, nil
}

// GetOrders 获取订单列表
func (interactor *OrderInteractor) GetOrders(
	ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.OrderListFilter,
	withItems bool,
) (orders []*domain.Order, total int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.GetOrders")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if withItems {
		orders, total, err = interactor.DataStore.OrderRepo().GetOrdersWithItems(ctx, pager, filter)
	} else {
		orders, total, err = interactor.DataStore.OrderRepo().GetOrders(ctx, pager, filter)
	}

	if err != nil {
		err = fmt.Errorf("failed to get orders: %w", err)
		return
	}

	return
}

// 锁订单
func (interactor *OrderInteractor) lockOrder(ctx context.Context, no string) (unlock func(), err error) {
	mu := interactor.MutexManager.NewMutex(domain.NewMutexOrderKey(no))
	if err = mu.Lock(ctx); err != nil {
		if domain.IsAlreadyTakenError(err) {
			err = domain.ParamsErrorf("订单 %s 正在被其他操作修改", no)
		}
		err = fmt.Errorf("failed to lock order: %w", err)
		return
	}
	unlock = func() {
		if _, err := mu.Unlock(ctx); err != nil {
			logger := logging.FromContext(ctx).Named("OrderInteractor.AppendItem")
			logger.Errorf("failed to unlock order: %s", err)
		}
	}
	return
}

// 修改订单商品价格
func (interactor *OrderInteractor) ModifyItemPrice(ctx context.Context, params *domain.ModifyItemPriceParams) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.ModifyItemPrice")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 锁订单
	unlock, err := interactor.lockOrder(ctx, params.OrderNo)
	if err != nil {
		return
	}
	defer unlock()

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) (err error) {
		order, err := ds.OrderRepo().FindByNo(ctx, params.OrderNo)
		if err != nil {
			if domain.IsNotFound(err) {
				err = domain.ParamsErrorf("订单不存在")
				return
			}
			err = fmt.Errorf("failed to find order: %w", err)
			return
		}

		if order.StoreID != params.Operator.StoreID {
			err = domain.ParamsErrorf("订单 %s 不存在", params.OrderNo)
			return
		}

		if order.Status != domain.OrderStatusUnpaid {
			err = domain.ParamsErrorf("订单 %s 当前状态无法修改商品价格", params.OrderNo)
			return
		}

		has, err := ds.OrderRepo().HasIncompletePayment(ctx, order.ID)
		if err != nil {
			err = fmt.Errorf("failed to check incomplete payment: %w", err)
			return
		}
		if has {
			err = domain.ParamsErrorf("订单 %s 正在支付中，请稍后重试", params.OrderNo)
			return
		}

		item, ok := lo.Find(order.Items, func(item *domain.OrderItem) bool {
			return item.ID == params.ItemID
		})

		if !ok {
			err = domain.ParamsErrorf("订单 %s 中不存在商品 %d", params.OrderNo, params.ItemID)
			return
		}

		item.Price = params.Price
		item.Amount = item.Price.Mul(item.Quantity)

		// 更新商品
		updatedItem, err := ds.OrderRepo().UpdateItem(ctx, item)
		if err != nil {
			err = fmt.Errorf("failed to update order item: %w", err)
			return
		}

		order.Items = lo.ReplaceAll(order.Items, item, updatedItem)
		totalAmount := order.Items.TotalAmount()
		order.TotalPrice = totalAmount
		order.RealPrice = totalAmount
		order.Discount = decimal.Zero
		order.PointsAvailable = order.Items.PointsAvailable(order.StoreID)

		// 更新订单
		updatedOrder, err := ds.OrderRepo().Update(ctx, order)
		if err != nil {
			err = fmt.Errorf("failed to update order: %w", err)
			return
		}
		updatedOrder.Items = order.Items

		// 触发改价事件
		err = interactor.OrderEventTrigger.FireModifyPrice(ctx, &domain.OrderEventBaseParams{
			DataStore: ds,
			Order:     updatedOrder,
			Operator:  params.Operator,
		})
		if err != nil {
			err = fmt.Errorf("failed to trigger modify price event: %w", err)
			return
		}

		return nil
	})

	return
}

// 添加商品
func (interactor *OrderInteractor) AppendItems(ctx context.Context, params *domain.AppendItemParams) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.AppendItems")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 锁订单
	unlock, err := interactor.lockOrder(ctx, params.OrderNo)
	if err != nil {
		return
	}
	defer unlock()

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) (err error) {
		order, err := ds.OrderRepo().FindByNo(ctx, params.OrderNo)
		if err != nil {
			if domain.IsNotFound(err) {
				err = domain.ParamsErrorf("订单不存在")
				return
			}
			err = fmt.Errorf("failed to find order: %w", err)
			return
		}

		if params.Operator.GetOperatorType() == domain.OperatorTypeCustomer {
			if order.TableID != params.TableID {
				err = domain.ParamsErrorf("订单 %s 不存在", params.OrderNo)
				return
			}
		} else {
			if order.StoreID != params.Operator.GetOperatorStoreID() {
				err = domain.ParamsErrorf("订单 %s 不存在", params.OrderNo)
				return
			}
		}

		if order.Status != domain.OrderStatusUnpaid {
			err = domain.ParamsErrorf("订单 %s 当前状态无法添加商品", params.OrderNo)
			return
		}

		has, err := ds.OrderRepo().HasIncompletePayment(ctx, order.ID)
		if err != nil {
			err = fmt.Errorf("failed to check incomplete payment: %w", err)
			return
		}
		if has {
			err = domain.ParamsErrorf("订单 %s 正在支付中，请稍后重试", params.OrderNo)
			return
		}

		// 取出所有商品ID
		productIDs := lo.Map(params.Items, func(item *domain.CreateOrderItem, _ int) int {
			return item.ProductID
		})

		// 获取所有商品信息
		products, err := interactor.DataStore.ProductRepo().GetDetailsByIDs(ctx, productIDs)
		if err != nil {
			return fmt.Errorf("failed to get products: %w", err)
		}

		productsMap := lo.KeyBy(products, func(item *domain.Product) int { return item.ID })

		// 构造订单商品结构 domain.OrderItem
		orderItems := make(domain.OrderItems, 0, len(params.Items))
		for _, item := range params.Items {
			product, ok := productsMap[item.ProductID]
			if !ok {
				err = domain.ParamsErrorf("商品 %d 不存在", item.ProductID)
				return
			}

			if product.StoreID != order.StoreID {
				err = domain.ParamsErrorf("商品 %d 不存在", item.ProductID)
				return
			}

			if !product.CheckSale() {
				err = domain.ParamsErrorf("商品 %d 不可售", item.ProductID)
				return
			}

			var orderItem *domain.OrderItem
			orderItem, err = interactor.buildOrderItem(item, product)
			if err != nil {
				err = fmt.Errorf("failed to build order item: %w", err)
				return
			}

			orderItems = append(orderItems, orderItem)
		}

		// 插入新商品
		newItems, err := ds.OrderRepo().AppendItems(ctx, order.ID, orderItems)
		if err != nil {
			err = fmt.Errorf("failed to append order items: %w", err)
			return
		}

		order.Items = append(order.Items, newItems...)
		totalAmount := order.Items.TotalAmount()
		order.TotalPrice = totalAmount
		order.RealPrice = totalAmount
		order.Discount = decimal.Zero
		order.PointsAvailable = order.Items.PointsAvailable(order.StoreID)

		// 更新订单
		updatedOrder, err := ds.OrderRepo().Update(ctx, order)
		if err != nil {
			err = fmt.Errorf("failed to update order: %w", err)
			return
		}
		updatedOrder.Items = order.Items

		// 触发增加商品事件
		err = interactor.OrderEventTrigger.FireAppendItem(ctx, &domain.OrderEventBaseParams{
			DataStore:     ds,
			Order:         updatedOrder,
			OperatedItems: newItems,
			Operator:      params.Operator,
		})
		if err != nil {
			err = fmt.Errorf("failed to trigger append item event: %w", err)
			return
		}

		return nil
	})

	return
}

// AppendItemsFromCart 从购物车追加商品到订单
func (interactor *OrderInteractor) AppendItemsFromCart(ctx context.Context, params *domain.AppendItemParams) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.AppendItemsFromCart")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 锁购物车
	unlock, err := interactor.lockCart(ctx, params.TableID)
	if err != nil {
		return err
	}
	defer unlock()

	// 获取购物车内容
	cartItems, err := interactor.DataStore.OrderCartRepo().ListByTable(ctx, params.TableID, true)
	if err != nil {
		return fmt.Errorf("failed to get cart items: %w", err)
	}
	if len(cartItems) == 0 {
		return domain.ParamsErrorf("购物车为空")
	}

	// 4. 构造订单商品结构
	orderItems := make([]*domain.CreateOrderItem, 0, len(cartItems))
	for _, item := range cartItems {
		orderItems = append(orderItems, &domain.CreateOrderItem{
			ProductID:       item.ProductID,
			ProductSpecID:   item.ProductSpecID,
			ProductAttrID:   item.AttrID,
			ProductRecipeID: item.RecipeID,
			Quantity:        item.Quantity,
			Price:           item.Price,
			Remark:          "",
		})
	}

	// 追加商品到订单
	params.Items = orderItems
	err = interactor.AppendItems(ctx, params)
	if err != nil {
		return err
	}

	// 清空购物车
	err = interactor.DataStore.OrderCartRepo().ClearByTable(ctx, params.TableID)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

// 取消订单
func (interactor *OrderInteractor) CancelOrder(ctx context.Context, no string, operator any) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.CancelOrder")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 锁订单
	unlock, err := interactor.lockOrder(ctx, no)
	if err != nil {
		return
	}
	defer unlock()

	// 获取订单
	order, err := interactor.DataStore.OrderRepo().FindByNo(ctx, no)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsErrorf("订单 %s 不存在", no)
			return
		}
		err = fmt.Errorf("failed to find order: %w", err)
		return
	}

	if operator != nil {
		var store *domain.Store
		switch v := operator.(type) {
		default:
			return fmt.Errorf("invalid operator type")
		case *domain.FrontendUser:
			store = v.Store
		case *domain.BackendUser:
			store = v.Store
		}

		if order.StoreID != store.ID {
			err = domain.ParamsErrorf("订单 %s 不存在", no)
		}
	}

	// 检查订单状态
	if order.Status != domain.OrderStatusUnpaid {
		err = domain.ParamsErrorf("订单 %s 当前状态无法取消", no)
		return
	}

	has, err := interactor.DataStore.OrderRepo().HasIncompletePayment(ctx, order.ID)
	if err != nil {
		err = fmt.Errorf("failed to check incomplete payment: %w", err)
		return
	}
	if has {
		err = domain.ParamsErrorf("订单 %s 正在支付中，请稍后重试", no)
		return
	}

	order.Status = domain.OrderStatusCanceled

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) (err error) {
		// 更新订单状态
		updatedOrder, err := ds.OrderRepo().Update(ctx, order)
		if err != nil {
			err = fmt.Errorf("failed to update order: %w", err)
			return
		}
		updatedOrder.Items = order.Items

		// 触发订单取消事件
		err = interactor.OrderEventTrigger.FireCancel(ctx, &domain.OrderEventBaseParams{
			DataStore: ds,
			Order:     updatedOrder,
			Operator:  operator,
		})
		if err != nil {
			err = fmt.Errorf("failed to trigger cancel order event: %w", err)
			return
		}

		return nil
	})

	return
}

// 优惠订单
func (interactor *OrderInteractor) DiscountOrder(ctx context.Context, params *domain.DiscountOrderParams) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.DiscountOrder")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 锁订单
	unlock, err := interactor.lockOrder(ctx, params.OrderNo)
	if err != nil {
		return
	}
	defer unlock()

	// 获取订单
	order, err := interactor.DataStore.OrderRepo().FindByNo(ctx, params.OrderNo)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsErrorf("订单 %s 不存在", params.OrderNo)
			return
		}
		err = fmt.Errorf("failed to find order: %w", err)
		return
	}

	if order.Status != domain.OrderStatusUnpaid {
		err = domain.ParamsErrorf("订单 %s 当前状态无法折扣", params.OrderNo)
		return
	}

	has, err := interactor.DataStore.OrderRepo().HasIncompletePayment(ctx, order.ID)
	if err != nil {
		err = fmt.Errorf("failed to check incomplete payment: %w", err)
		return
	}
	if has {
		err = domain.ParamsErrorf("订单 %s 正在支付中，请稍后重试", params.OrderNo)
		return
	}

	order.Discount = params.Discount
	order.RealPrice = order.TotalPrice.Sub(order.Discount)
	if order.RealPrice.LessThan(decimal.Zero) {
		err = domain.ParamsErrorf("订单 %s 折扣金额大于总价", params.OrderNo)
		return
	}

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) (err error) {
		// 更新订单
		updatedOrder, err := ds.OrderRepo().Update(ctx, order)
		if err != nil {
			err = fmt.Errorf("failed to update order: %w", err)
			return
		}
		updatedOrder.Items = order.Items

		// 触发折扣事件
		err = interactor.OrderEventTrigger.FireDiscount(ctx, &domain.OrderEventBaseParams{
			DataStore: ds,
			Order:     updatedOrder,
			Operator:  params.Operator,
		})
		if err != nil {
			err = fmt.Errorf("failed to trigger discount order event: %w", err)
			return
		}

		return nil
	})

	return
}

// 删除商品
func (interactor *OrderInteractor) RemoveItems(ctx context.Context, params *domain.RemoveItemParams) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.RemoveItems")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 锁订单
	unlock, err := interactor.lockOrder(ctx, params.OrderNo)
	if err != nil {
		return
	}
	defer unlock()

	// 获取订单
	order, err := interactor.DataStore.OrderRepo().FindByNo(ctx, params.OrderNo)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsErrorf("订单 %s 不存在", params.OrderNo)
		}
		return
	}

	if params.Operator.StoreID != order.StoreID {
		err = domain.ParamsErrorf("订单 %s 不存在", params.OrderNo)
		return
	}

	if order.Status != domain.OrderStatusUnpaid {
		err = domain.ParamsErrorf("订单 %s 当前状态无法删除商品", params.OrderNo)
		return
	}

	has, err := interactor.DataStore.OrderRepo().HasIncompletePayment(ctx, order.ID)
	if err != nil {
		err = fmt.Errorf("failed to check incomplete payment: %w", err)
		return
	}
	if has {
		err = domain.ParamsErrorf("订单 %s 正在支付中，请稍后重试", params.OrderNo)
		return
	}

	item, ok := lo.Find(order.Items, func(item *domain.OrderItem) bool {
		return item.ID == params.ItemID
	})

	if !ok {
		err = domain.ParamsErrorf("订单 %s 中不存在商品 %d", params.OrderNo, params.ItemID)
		return
	}

	// 删除的商品项
	operatedItems := make(domain.OrderItems, 0)
	removeItem := *item
	removeItem.Quantity = params.Quantity
	removeItem.Amount = removeItem.Price.Mul(removeItem.Quantity)
	operatedItems = append(operatedItems, &removeItem)

	item.Quantity = item.Quantity.Sub(params.Quantity)
	if item.Quantity.LessThan(decimal.Zero) {
		err = domain.ParamsErrorf("订单 %s 商品 %d 数量不足", params.OrderNo, params.ItemID)
		return
	}

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) (err error) {
		if item.Quantity.Equal(decimal.Zero) {
			// 删除商品
			if err = ds.OrderRepo().RemoveItems(ctx, order.ID, params.ItemID); err != nil {
				err = fmt.Errorf("failed to remove order item: %w", err)
				return
			}
			order.Items = lo.Filter(order.Items, func(item *domain.OrderItem, _ int) bool {
				return item.ID != params.ItemID
			})
		} else {
			// 更新商品
			item.Amount = item.Price.Mul(item.Quantity)
			var updatedItem *domain.OrderItem
			updatedItem, err = ds.OrderRepo().UpdateItem(ctx, item)
			if err != nil {
				err = fmt.Errorf("failed to update order item: %w", err)
				return
			}
			order.Items = lo.ReplaceAll(order.Items, item, updatedItem)
		}

		// 更新订单
		totalAmount := order.Items.TotalAmount()
		order.TotalPrice = totalAmount
		order.RealPrice = totalAmount
		order.Discount = decimal.Zero
		order.PointsAvailable = order.Items.PointsAvailable(order.StoreID)

		// 更新订单
		updatedOrder, err := ds.OrderRepo().Update(ctx, order)
		if err != nil {
			err = fmt.Errorf("failed to update order: %w", err)
			return
		}
		updatedOrder.Items = order.Items

		// 触发删除商品事件
		err = interactor.OrderEventTrigger.FireRemoveItem(ctx, &domain.OrderEventBaseParams{
			DataStore:     ds,
			Order:         updatedOrder,
			OperatedItems: operatedItems,
			Operator:      params.Operator,
		})
		if err != nil {
			err = fmt.Errorf("failed to trigger remove item event: %w", err)
			return
		}

		return nil
	})

	return
}

// 订单转台
func (interactor *OrderInteractor) TurnTable(ctx context.Context, params *domain.TurnTableParams) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.TurnTable")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 锁订单
	unlock, err := interactor.lockOrder(ctx, params.OrderNo)
	if err != nil {
		return
	}
	defer unlock()

	// 获取订单
	order, err := interactor.DataStore.OrderRepo().FindByNo(ctx, params.OrderNo)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsErrorf("订单 %s 不存在", params.OrderNo)
		}
		return
	}

	if params.Operator.StoreID != order.StoreID {
		err = domain.ParamsErrorf("订单 %s 不存在", params.OrderNo)
		return
	}

	if order.Status != domain.OrderStatusUnpaid {
		err = domain.ParamsErrorf("订单 %s 当前状态无法转台", params.OrderNo)
		return
	}

	if order.TableID == 0 {
		err = domain.ParamsErrorf("订单 %s 未选择台桌，无法转台", params.OrderNo)
		return
	}

	has, err := interactor.DataStore.OrderRepo().HasIncompletePayment(ctx, order.ID)
	if err != nil {
		err = fmt.Errorf("failed to check incomplete payment: %w", err)
		return
	}
	if has {
		err = domain.ParamsErrorf("订单 %s 正在支付中，请稍后重试", params.OrderNo)
		return
	}

	// 获取新桌子
	table, err := interactor.DataStore.TableRepo().FindByID(ctx, params.TableID)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsErrorf("桌子 %d 不存在", params.TableID)
			return
		}
		err = fmt.Errorf("failed to get table: %w", err)
		return
	}

	if table.StoreID != order.StoreID {
		err = domain.ParamsErrorf("桌子 %d 不存在", params.TableID)
		return
	}

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) (err error) {
		// 更新订单
		oldTableID := order.TableID
		order.TableID = params.TableID
		order.TableName = table.Name
		updatedOrder, err := ds.OrderRepo().Update(ctx, order)
		if err != nil {
			err = fmt.Errorf("failed to update order: %w", err)
			return
		}
		updatedOrder.Items = order.Items

		// 触发转台事件
		err = interactor.OrderEventTrigger.FireTurnTable(ctx, &domain.OrderEventTurnTableParams{
			OrderEventBaseParams: domain.OrderEventBaseParams{
				DataStore: ds,
				Order:     updatedOrder,
				Operator:  params.Operator,
			},
			OldTableID: oldTableID,
		})
		if err != nil {
			err = fmt.Errorf("failed to trigger turn table event: %w", err)
			return
		}

		return nil
	})

	return
}

// 订单现金支付
func (interactor *OrderInteractor) CashPaid(ctx context.Context, params *domain.OrderCashPaidParams) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.Paid")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 锁订单
	unlock, err := interactor.lockOrder(ctx, params.OrderNo)
	if err != nil {
		return
	}
	defer unlock()

	// 获取订单
	od, err := interactor.DataStore.OrderRepo().FindByNo(ctx, params.OrderNo)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsErrorf("订单 %s 不存在", params.OrderNo)
		}
		err = fmt.Errorf("failed to get order: %w", err)
		return
	}

	if params.Operator.StoreID != od.StoreID {
		err = domain.ParamsErrorf("订单 %s 不存在", params.OrderNo)
		return
	}

	has, err := interactor.DataStore.OrderRepo().HasIncompletePayment(ctx, od.ID)
	if err != nil {
		err = fmt.Errorf("failed to check incomplete payment: %w", err)
		return
	}
	if has {
		err = domain.ParamsErrorf("订单 %s 正在支付中，请稍后重试", params.OrderNo)
		return
	}

	if err = od.CanPaid(params.Amount, false); err != nil {
		return
	}

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) (err error) {
		_, err = interactor.OrderDomainService.Paid(ctx, ds, &order.PaidParams{
			Order:    od,
			Operator: params.Operator,
			Channel:  domain.OrderPaidChannelCash,
			Amount:   params.Amount,
		})
		if err != nil {
			return fmt.Errorf("failed to paid order: %w", err)
		}
		return
	})

	return
}

func (interactor *OrderInteractor) ScanPaid(ctx context.Context, params *domain.OrderScanPaidParams) (seqNo string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.ScanPaid")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 锁订单
	unlock, err := interactor.lockOrder(ctx, params.OrderNo)
	if err != nil {
		return
	}
	defer unlock()

	// 获取订单
	od, err := interactor.DataStore.OrderRepo().FindByNo(ctx, params.OrderNo)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsErrorf("订单 %s 不存在", params.OrderNo)
		}
		err = fmt.Errorf("failed to get order: %w", err)
		return
	}

	if params.Operator.StoreID != od.StoreID {
		err = domain.ParamsErrorf("订单 %s 不存在", params.OrderNo)
		return
	}

	has, err := interactor.DataStore.OrderRepo().HasIncompletePayment(ctx, od.ID)
	if err != nil {
		err = fmt.Errorf("failed to check incomplete payment: %w", err)
		return
	}
	if has {
		err = domain.ParamsErrorf("订单 %s 正在支付中，请稍后重试", params.OrderNo)
		return
	}

	pointPay := domain.IsPointCode(params.AuthCode)
	pointWalletPay := domain.IsPointWalletCode(params.AuthCode)

	if err = od.CanPaid(params.Amount, pointPay); err != nil {
		return
	}

	store := params.Operator.Store
	var provider payment.PaymentProvider
	var notifyURL string
	if pointPay {
		notifyURL = params.ZxhNotifyURL
		if store.ZxhID == "" || store.ZxhSecret == "" {
			err = domain.ParamsErrorf("门店未配置积分支付")
			return
		}
		provider = payment.NewZxhPaymentProvider(interactor.zxhManager, store.ZxhID, store.ZxhSecret)
	} else if pointWalletPay {
		notifyURL = params.ZxhNotifyURL
		if store.ZxhID == "" || store.ZxhSecret == "" {
			err = domain.ParamsErrorf("门店的知心话支付配置有误")
			return
		}
		provider = payment.NewZxhWalletPaymentProvider(interactor.zxhManager, store.ZxhID, store.ZxhSecret)
	} else {
		notifyURL = params.HuifuNotifyURL
		if store.HuifuID == "" {
			err = domain.ParamsErrorf("门店未配置微信/支付宝支付")
			return
		}
		provider = payment.NewHuifuPaymentProvider(interactor.bsPay, store.HuifuID)
	}

	paymentParams := &payment.PaymentParams{
		AuthCode:   params.AuthCode,
		NotifyURL:  notifyURL,
		Amount:     params.Amount,
		GoodsDesc:  od.GoodsDesc(),
		IPAddr:     params.IPAddr,
		PayBizType: domain.PayBizTypeOrder,
		BizID:      od.ID,
		Creator:    params.Operator,
		StoreID:    od.StoreID,
	}
	paymentParams.Set("merchant_name", store.Name)

	seqNo, err = interactor.PaymentDomainService.ProcessPayment(ctx, provider, paymentParams)

	if err != nil {
		err = fmt.Errorf("failed to process payment: %w", err)
		return
	}

	return
}

func (interactor *OrderInteractor) GetOrderRange(ctx context.Context, filter *domain.OrderListFilter) (rg domain.OrderRange, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderInteractor.GetOrderRange")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	rg, err = interactor.DataStore.OrderRepo().GetOrderRange(ctx, filter)
	if err != nil {
		err = fmt.Errorf("failed to get order range: %w", err)
		return
	}

	return
}
