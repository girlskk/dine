package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/uhttp"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type OrderHandler struct {
	OrderInteractor domain.OrderInteractor
	TableInteractor domain.TableInteractor
}

func NewOrderHandler(interactor domain.OrderInteractor, tableInteractor domain.TableInteractor) *OrderHandler {
	return &OrderHandler{
		OrderInteractor: interactor,
		TableInteractor: tableInteractor,
	}
}

func (h *OrderHandler) Routes(r gin.IRouter) {
	r = r.Group("/order")
	r.POST("/create", h.CreateOrder())
	r.POST("/detail", h.OrderDetail())
	r.POST("/list", h.OrderList())
	r.POST("/modify_price", h.OrderModifyPrice())
	r.POST("/append_items", h.OrderAppendItems())
	r.POST("/turn_table", h.OrderTurnTable())
	r.POST("/remove_items", h.OrderRemoveItems())
	r.POST("/cancel", h.OrderCancel())
	r.POST("/discount", h.OrderDiscount())
	r.POST("/cash_paid", h.OrderCashPaid())
	r.POST("/scan_paid", h.OrderScanPaid())
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

		store := domain.FromStoreContext(ctx)
		user := domain.FromFrontendUserContext(ctx)

		items := lo.Map(req.Items, func(item types.CreateOrderReqItem, _ int) *domain.CreateOrderItem {
			return &domain.CreateOrderItem{
				ProductID:       item.ProductID,
				Quantity:        item.Quantity,
				Price:           item.Price,
				Remark:          item.Remark,
				ProductAttrID:   item.AttrID,
				ProductSpecID:   item.SpecID,
				ProductRecipeID: item.RecipeID,
			}
		})

		var table *domain.Table
		var err error
		if req.TableID > 0 {
			table, err = h.TableInteractor.Get(ctx, req.TableID)
			if err != nil {
				if msg, ok := domain.GetParamsErrorMessage(err); ok {
					c.Error(uerr.BadRequest(msg))
				} else {
					c.Error(fmt.Errorf("failed to get table: %w", err))
				}
				return
			}

			if table.StoreID != store.ID {
				c.Error(uerr.BadRequest(domain.ErrTableNotExists.Error()))
				return
			}
		}

		params := &domain.CreateOrderParams{
			Store:        store,
			Creator:      user,
			Table:        table,
			Items:        items,
			PeopleNumber: lo.Ternary(req.PeopleNumber > 0, req.PeopleNumber, 1),
		}

		od, err := h.OrderInteractor.CreateOrder(ctx, params)
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

// OrderDetail	获取订单详情
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

		store := domain.FromStoreContext(ctx)
		if store.ID != order.StoreID {
			c.Error(uerr.BadRequest("订单不存在"))
			return
		}

		response.Ok(c, order)
	}
}

// OrderList	订单列表
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	订单列表
//	@Accept		json
//	@Produce	json
//	@Param		data	body		types.OrderListReq	true	"请求参数"
//	@Success	200		{object}	types.OrderListResp	"成功"
//	@Router		/order/list [post]
func (h *OrderHandler) OrderList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.OrderList")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.OrderList")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.OrderListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(req.Page, req.Size)
		store := domain.FromStoreContext(ctx)

		orders, total, err := h.OrderInteractor.GetOrders(ctx, page, &domain.OrderListFilter{
			StoreID: store.ID,
			Status:  req.Status,
		}, false)
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

// OrderAppendItems	添加订单商品
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

		user := domain.FromFrontendUserContext(ctx)

		items := lo.Map(req.Items, func(item types.CreateOrderReqItem, _ int) *domain.CreateOrderItem {
			return &domain.CreateOrderItem{
				ProductID:       item.ProductID,
				Quantity:        item.Quantity,
				Price:           item.Price,
				Remark:          item.Remark,
				ProductAttrID:   item.AttrID,
				ProductSpecID:   item.SpecID,
				ProductRecipeID: item.RecipeID,
			}
		})

		err := h.OrderInteractor.AppendItems(ctx, &domain.AppendItemParams{
			OrderNo:  req.No,
			Items:    items,
			Operator: user,
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

// OrderModifyPrice	修改订单商品价格
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	修改订单商品价格
//	@Accept		json
//	@Produce	json
//	@Param		data	body	types.OrderModifyPriceReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/order/modify_price [post]
func (h *OrderHandler) OrderModifyPrice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.OrderModifyPrice")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.OrderModifyPrice")

		var req types.OrderModifyPriceReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		logger = logger.With("order_no", req.No)
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromFrontendUserContext(ctx)

		err := h.OrderInteractor.ModifyItemPrice(ctx, &domain.ModifyItemPriceParams{
			OrderNo:  req.No,
			ItemID:   req.ItemID,
			Price:    req.Price,
			Operator: user,
		})
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to modify price: %w", err))
			}
			return
		}

		response.Ok(c, nil)
	}
}

// OrderTurnTable 转台
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	转台
//	@Accept		json
//	@Produce	json
//	@Param		data	body	types.OrderTurnTableReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/order/turn_table [post]
func (h *OrderHandler) OrderTurnTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.OrderTurnTable")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.OrderTurnTable")

		var req types.OrderTurnTableReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		logger = logger.With("order_no", req.No)
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromFrontendUserContext(ctx)

		err := h.OrderInteractor.TurnTable(ctx, &domain.TurnTableParams{
			OrderNo:  req.No,
			TableID:  req.TableID,
			Operator: user,
		})
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to turn table: %w", err))
			}
			return
		}

		response.Ok(c, nil)
	}
}

// OrderRemoveItems 退菜
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	退菜
//	@Accept		json
//	@Produce	json
//	@Param		data	body	types.OrderRemoveItemsReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/order/remove_items [post]
func (h *OrderHandler) OrderRemoveItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.OrderRemoveItems")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.OrderRemoveItems")

		var req types.OrderRemoveItemsReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		logger = logger.With("order_no", req.No)
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromFrontendUserContext(ctx)

		err := h.OrderInteractor.RemoveItems(ctx, &domain.RemoveItemParams{
			OrderNo:  req.No,
			ItemID:   req.ItemID,
			Quantity: req.Quantity,
			Operator: user,
		})
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to remove items: %w", err))
			}
			return
		}

		response.Ok(c, nil)
	}
}

// OrderCancel 撤单
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	撤单
//	@Accept		json
//	@Produce	json
//	@Param		data	body	types.OrderCancelReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/order/cancel [post]
func (h *OrderHandler) OrderCancel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.OrderCancel")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.OrderCancel")

		var req types.OrderCancelReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		logger = logger.With("order_no", req.No)
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromFrontendUserContext(ctx)

		err := h.OrderInteractor.CancelOrder(ctx, req.No, user)
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to cancel: %w", err))
			}
			return
		}

		response.Ok(c, nil)
	}
}

// OrderDiscount 订单折扣
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	订单折扣
//	@Accept		json
//	@Produce	json
//	@Param		data	body	types.OrderDiscountReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/order/discount [post]
func (h *OrderHandler) OrderDiscount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.OrderDiscount")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.OrderDiscount")

		var req types.OrderDiscountReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		logger = logger.With("order_no", req.No)
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromFrontendUserContext(ctx)

		err := h.OrderInteractor.DiscountOrder(ctx, &domain.DiscountOrderParams{
			OrderNo:  req.No,
			Discount: req.Discount,
			Operator: user,
		})
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to discount: %w", err))
			}
			return
		}

		response.Ok(c, nil)
	}
}

// OrderCashPaid 订单现金支付
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	订单现金支付
//	@Accept		json
//	@Produce	json
//	@Param		data	body	types.OrderCashPaidReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/order/cash_paid [post]
func (h *OrderHandler) OrderCashPaid() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.OrderCashPaid")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.OrderCashPaid")

		var req types.OrderCashPaidReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		logger = logger.With("order_no", req.No)
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromFrontendUserContext(ctx)

		err := h.OrderInteractor.CashPaid(ctx, &domain.OrderCashPaidParams{
			OrderNo:  req.No,
			Amount:   req.Amount,
			Operator: user,
		})
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to cash paid: %w", err))
			}
			return
		}

		response.Ok(c, nil)
	}
}

// OrderScanPaid 订单扫码支付
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	订单扫码支付
//	@Accept		json
//	@Produce	json
//	@Param		data	body		types.OrderScanPaidReq	true	"请求参数"
//	@Success	200		{object}	types.OrderScanPaidResp	"成功"
//	@Router		/order/scan_paid [post]
func (h *OrderHandler) OrderScanPaid() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.OrderScanPaid")
		defer span.Finish()

		logger := logging.FromContext(ctx).Named("OrderHandler.OrderScanPaid")

		var req types.OrderScanPaidReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		logger = logger.With("order_no", req.No)
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		pointpay := req.Channel == types.OrderScanPaidChannelPoint
		isPointCode := domain.IsPointCode(req.AuthCode)
		pointWalletPay := domain.IsPointWalletCode(req.AuthCode)

		if pointWalletPay {
			c.Error(uerr.BadRequest("暂不支持该付款方式"))
			return
		}
		if pointpay && !isPointCode {
			c.Error(uerr.BadRequest("请使用知心话付款码"))
			return
		} else if !pointpay && isPointCode {
			c.Error(uerr.BadRequest("请选择积分支付"))
			return
		}

		user := domain.FromFrontendUserContext(ctx)

		seqNo, err := h.OrderInteractor.ScanPaid(ctx, &domain.OrderScanPaidParams{
			OrderNo:        req.No,
			Amount:         req.Amount,
			AuthCode:       req.AuthCode,
			Operator:       user,
			IPAddr:         uhttp.GetClientIP(c.Request),
			HuifuNotifyURL: uhttp.RelativeEndpoint(c.Request, frontend.ApiPrefixV1+"/payment/callback/pay/huifu"),
			ZxhNotifyURL:   uhttp.RelativeEndpoint(c.Request, frontend.ApiPrefixV1+"/payment/callback/pay/zxh"),
		})
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to scan paid: %w", err))
			}
			return
		}

		response.Ok(c, &types.OrderScanPaidResp{SeqNo: seqNo})
	}
}
