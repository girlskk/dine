package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/customer/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type OrderHandler struct {
	OrderInteractor domain.OrderInteractor
	TableInteractor domain.TableInteractor
	StoreInteractor domain.StoreInteractor
}

func NewOrderHandler(
	interactor domain.OrderInteractor,
	tableInteractor domain.TableInteractor,
	storeInteractor domain.StoreInteractor,
) *OrderHandler {
	return &OrderHandler{
		OrderInteractor: interactor,
		TableInteractor: tableInteractor,
		StoreInteractor: storeInteractor,
	}
}

func (h *OrderHandler) Routes(r gin.IRouter) {
	r = r.Group("/order")
	r.POST("/create", h.CreateOrder())
	r.POST("/append_items", h.OrderAppendItems())
	r.POST("/list", h.OrderList())
	r.POST("/detail", h.OrderDetail())
}

// CreateOrder 创建订单
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	创建订单
//	@Accept		json
//	@Produce	json
//	@Param		data	body		types.CreateOrderReq	true	"请求参数"
//	@Success	200		{object}	types.CreateOrderResp	"成功"
//	@Router		/order/create [post]
func (h *OrderHandler) CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.CreateOrder")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.CreateOrder")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CreateOrderReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		user := domain.FromCustomerContext(ctx)
		table, err := h.TableInteractor.Get(ctx, req.TableID)
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to get table: %w", err))
			}
			return
		}

		store, err := h.StoreInteractor.GetDetail(ctx, table.StoreID)
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to get store: %w", err))
			}
		}

		params := &domain.CreateOrderParams{
			Store:        store,
			Creator:      user,
			Table:        table,
			PeopleNumber: req.PeopleNumber,
			Source:       domain.OrderSourceMiniProgram,
		}

		od, err := h.OrderInteractor.CreateOrderFromCart(ctx, params)
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to create order: %w", err))
			}
			return
		}

		response.Ok(c, &types.CreateOrderResp{No: od.No})
	}
}

// OrderAppendItems 添加订单商品
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	添加订单商品
//	@Accept		json
//	@Produce	json
//	@Param		data	body	types.OrderAppendItemsReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/order/append_items [post]
func (h *OrderHandler) OrderAppendItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.OrderAppendItems")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.OrderAppendItems")

		var req types.OrderAppendItemsReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		logger = logger.With("order_no", req.No)
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromCustomerContext(ctx)

		err := h.OrderInteractor.AppendItemsFromCart(ctx, &domain.AppendItemParams{
			OrderNo:  req.No,
			Operator: user,
			TableID:  req.TableID,
		})
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to append items: %w", err))
			}
			return
		}

		response.Ok(c, nil)
	}
}

// OrderList 订单列表
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	订单列表
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	types.OrderListResp	"成功"
//	@Router		/order/list [post]
func (h *OrderHandler) OrderList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.OrderList")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.OrderList")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		page := upagination.New(1, upagination.MaxSize)
		user := domain.FromCustomerContext(ctx)
		orders, total, err := h.OrderInteractor.GetOrders(ctx, page, &domain.OrderListFilter{
			CreatorID:   user.ID,
			CreatorType: domain.OperatorTypeCustomer,
		}, true)
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to get orders: %w", err))
			}
			return
		}

		response.Ok(c, &types.OrderListResp{
			Orders: orders,
			Total:  total,
		})
	}
}

// OrderDetail 获取订单详情
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	订单详情
//	@Accept		json
//	@Produce	json
//	@Param		data	body		types.OrderDetailReq	true	"请求参数"
//	@Success	200		{object}	domain.Order			"成功"
//	@Router		/order/detail [post]
func (h *OrderHandler) OrderDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.OrderDetail")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.OrderDetail")

		var req types.OrderDetailReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		logger = logger.With("order_no", req.No)
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		order, err := h.OrderInteractor.GetOrder(ctx, req.No)
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to get order: %w", err))
			}
			return
		}

		response.Ok(c, order)
	}
}
