package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
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
	r.GET("/:id", h.Get())
	r.GET("", h.List())
}

func (h *OrderHandler) NoAuths() []string {
	return []string{}
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

// List
//
//	@Tags		订单
//	@Security	BearerAuth
//	@Summary	获取订单列表
//	@Accept		json
//	@Produce	json
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

		user := domain.FromBackendUserContext(ctx)

		params := domain.OrderListParams{
			MerchantID:    user.MerchantID,
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
