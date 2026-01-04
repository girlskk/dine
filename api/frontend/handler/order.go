package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type OrderHandler struct {
	OrderInteractor domain.OrderInteractor
	Seq             domain.DailySequence
}

func NewOrderHandler(orderInteractor domain.OrderInteractor, seq domain.DailySequence) *OrderHandler {
	return &OrderHandler{
		OrderInteractor: orderInteractor,
		Seq:             seq,
	}
}

func (h *OrderHandler) Routes(r gin.IRouter) {
	r = r.Group("/order")
	r.POST("", h.Create())
	r.GET("/:id", h.Get())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("", h.List())
}

func (h *OrderHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		订单
//	@Security	BearerAuth
//	@Summary	创建订单
//	@Accept		json
//	@Produce	json
//	@Param		data	body		types.CreateOrderReq	true	"请求信息"
//	@Success	200		{object}	domain.Order			"成功"
//	@Router		/order [post]
func (h *OrderHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CreateOrderReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		o := &domain.Order{
			ID:            uuid.New(),
			MerchantID:    req.MerchantID,
			StoreID:       req.StoreID,
			BusinessDate:  req.BusinessDate,
			ShiftNo:       req.ShiftNo,
			OrderNo:       req.OrderNo,
			DiningMode:    domain.DiningMode(req.DiningMode),
			TableID:       req.TableID,
			TableName:     req.TableName,
			GuestCount:    req.GuestCount,
			PlacedAt:      req.PlacedAt,
			PlacedBy:      req.PlacedBy,
			Refund:        req.Refund,
			Store:         req.Store,
			Pos:           req.Pos,
			Cashier:       req.Cashier,
			OrderProducts: req.OrderProducts,
			TaxRates:      req.TaxRates,
			Fees:          req.Fees,
			Payments:      req.Payments,
			Amount:        req.Amount,
		}

		if req.OrderType != "" {
			o.OrderType = domain.OrderType(req.OrderType)
		}
		if req.OrderStatus != "" {
			o.OrderStatus = domain.OrderStatus(req.OrderStatus)
		}
		if req.PaymentStatus != "" {
			o.PaymentStatus = domain.PaymentStatus(req.PaymentStatus)
		}
		if req.Channel != "" {
			o.Channel = domain.Channel(req.Channel)
		}

		// 自动生成订单号
		if o.OrderNo == "" {
			orderNo, err := h.generateOrderNo(ctx, o)
			if err != nil {
				c.Error(fmt.Errorf("failed to generate order_no: %w", err))
				return
			}
			o.OrderNo = orderNo
		}

		err := h.OrderInteractor.Create(ctx, o)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			if domain.IsConflict(err) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			c.Error(fmt.Errorf("failed to create order: %w", err))
			return
		}

		response.Ok(c, o)
	}
}

// Get
//
//	@Tags		订单
//	@Security	BearerAuth
//	@Summary	获取订单详情
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string			true	"订单ID"
//	@Success	200	{object}	domain.Order	"成功"
//	@Router		/order/{id} [get]
func (h *OrderHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		o, err := h.OrderInteractor.Get(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			c.Error(fmt.Errorf("failed to get order: %w", err))
			return
		}

		response.Ok(c, o)
	}
}

// Update
//
//	@Tags		订单
//	@Security	BearerAuth
//	@Summary	更新订单
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string					true	"订单ID"
//	@Param		data	body		types.UpdateOrderReq	true	"请求信息"
//	@Success	200		{object}	domain.Order			"成功"
//	@Router		/order/{id} [put]
func (h *OrderHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.UpdateOrderReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		o := &domain.Order{
			ID:            id,
			BusinessDate:  req.BusinessDate,
			ShiftNo:       req.ShiftNo,
			OrderNo:       req.OrderNo,
			Refund:        req.Refund,
			TableID:       req.TableID,
			TableName:     req.TableName,
			GuestCount:    req.GuestCount,
			PlacedAt:      req.PlacedAt,
			PaidAt:        req.PaidAt,
			PlacedBy:      req.PlacedBy,
			Store:         req.Store,
			Pos:           req.Pos,
			Cashier:       req.Cashier,
			OrderProducts: req.OrderProducts,
			TaxRates:      req.TaxRates,
			Fees:          req.Fees,
			Payments:      req.Payments,
			Amount:        req.Amount,
		}

		if req.OrderType != "" {
			o.OrderType = domain.OrderType(req.OrderType)
		}
		if req.DiningMode != "" {
			o.DiningMode = domain.DiningMode(req.DiningMode)
		}
		if req.OrderStatus != "" {
			o.OrderStatus = domain.OrderStatus(req.OrderStatus)
		}
		if req.PaymentStatus != "" {
			o.PaymentStatus = domain.PaymentStatus(req.PaymentStatus)
		}
		if req.Channel != "" {
			o.Channel = domain.Channel(req.Channel)
		}

		err = h.OrderInteractor.Update(ctx, o)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			c.Error(fmt.Errorf("failed to update order: %w", err))
			return
		}

		response.Ok(c, o)
	}
}

// Delete
//
//	@Tags		订单
//	@Security	BearerAuth
//	@Summary	删除订单
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"订单ID"
//	@Success	200	"No Content"
//	@Router		/order/{id} [delete]
func (h *OrderHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		if err := h.OrderInteractor.Delete(ctx, id); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			c.Error(fmt.Errorf("failed to delete order: %w", err))
			return
		}

		response.Ok(c, nil)
	}
}

// List
//
//	@Tags		订单
//	@Security	BearerAuth
//	@Summary	获取订单列表
//	@Accept		json
//	@Produce	json
//	@Param		merchant_id		query		string				true	"品牌商ID"
//	@Param		store_id		query		string				true	"门店ID"
//	@Param		business_date	query		string				false	"营业日"
//	@Param		order_no		query		string				false	"订单号"
//	@Param		order_type		query		string				false	"订单类型"	Enums(SALE,REFUND,PARTIAL_REFUND)
//	@Param		order_status	query		string				false	"订单状态"	Enums(PLACED,COMPLETED,CANCELLED)
//	@Param		payment_status	query		string				false	"支付状态"	Enums(UNPAID,PAYING,PAID,REFUNDED)
//	@Param		page			query		int					false	"页码"
//	@Param		size			query		int					false	"每页数量"
//	@Success	200				{object}	types.ListOrderResp	"成功"
//	@Router		/order [get]
func (h *OrderHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ListOrderReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		merchantID, err := uuid.Parse(req.MerchantID)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams,
				fmt.Errorf("invalid merchant_id: %w", err)))
			return
		}

		storeID, err := uuid.Parse(req.StoreID)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams,
				fmt.Errorf("invalid store_id: %w", err)))
			return
		}

		params := domain.OrderListParams{
			MerchantID:    merchantID,
			StoreID:       storeID,
			BusinessDate:  req.BusinessDate,
			OrderNo:       req.OrderNo,
			OrderType:     domain.OrderType(req.OrderType),
			OrderStatus:   domain.OrderStatus(req.OrderStatus),
			PaymentStatus: domain.PaymentStatus(req.PaymentStatus),
			Page:          req.Page,
			Size:          req.Size,
		}

		items, total, err := h.OrderInteractor.List(ctx, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			c.Error(fmt.Errorf("failed to list orders: %w", err))
			return
		}

		p := req.ToPagination()
		p.SetTotal(total)

		response.Ok(c, &types.ListOrderResp{
			Items:      items,
			Pagination: p,
		})
	}
}

func (h *OrderHandler) generateOrderNo(ctx context.Context, o *domain.Order) (string, error) {
	storePart := ""
	if o.Store.StoreCode != "" {
		storePart = o.Store.StoreCode
	}

	datePart := strings.ReplaceAll(o.BusinessDate, "-", "")
	prefix := fmt.Sprintf("%s:%s", domain.DailySequencePrefixOrderNo, o.StoreID.String())
	seq, err := h.Seq.Next(ctx, prefix)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s%06d", storePart, datePart, seq), nil
}
