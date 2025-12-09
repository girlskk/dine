package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type OrderHandler struct {
	OrderInteractor      domain.OrderInteractor
	DataExportInteractor domain.DataExportInteractor
}

func NewOrderHandler(orderInteractor domain.OrderInteractor, dataExportInteractor domain.DataExportInteractor) *OrderHandler {
	return &OrderHandler{
		OrderInteractor:      orderInteractor,
		DataExportInteractor: dataExportInteractor,
	}
}

func (h *OrderHandler) Routes(r gin.IRouter) {
	r = r.Group("/order")
	r.POST("/list", h.OrderList())
	r.POST("/detail", h.OrderDetail())
	r.POST("/list/export", h.OrderListExport())
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
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.OrderListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		orders, total, err := h.OrderInteractor.GetOrders(ctx, req.ToPagination(), &domain.OrderListFilter{
			StoreID:           req.StoreID,
			Status:            req.Status,
			HasItemName:       req.HasItemName,
			MemberNameOrPhone: req.MemberNameOrPhone,
			CreatedAtGte:      req.CreatedAtGte.ToPtrStartOfDay(),
			CreatedAtLte:      req.CreatedAtLte.ToPtrEndOfDay(),
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

		response.Ok(c, order)
	}
}

// OrderListExport	订单导出
//
//	@Tags		订单管理
//	@Security	BearerAuth
//	@Summary	订单导出
//	@Accept		json
//	@Produce	json
//	@Param		data	body	types.OrderListExportReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/order/list/export [post]
func (h *OrderHandler) OrderListExport() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "OrderHandler.OrderListExport")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OrderHandler.OrderListExport")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.OrderListExportReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if !req.CreatedAtGte.IsValid() || !req.CreatedAtLte.IsValid() {
			c.Error(uerr.BadRequest("时间范围不能为空"))
			return
		}

		filter := &domain.OrderListFilter{
			StoreID:           req.StoreID,
			Status:            req.Status,
			HasItemName:       req.HasItemName,
			MemberNameOrPhone: req.MemberNameOrPhone,
			CreatedAtGte:      req.CreatedAtGte.ToPtrStartOfDay(),
			CreatedAtLte:      req.CreatedAtLte.ToPtrEndOfDay(),
		}

		orderRange, err := h.OrderInteractor.GetOrderRange(ctx, filter)
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to get order range: %w", err))
			}
			return
		}

		if orderRange.Count < 1 {
			c.Error(uerr.BadRequest("没有订单可导出"))
			return
		}

		filter.IDGte = orderRange.MinID
		filter.IDLte = orderRange.MaxID

		totalPages := upagination.TotalPages(orderRange.Count, domain.OrderListExportSingleMaxSize)
		params := make([]*domain.OrderListExportParams, 0, totalPages)
		for i := range totalPages {
			page := i + 1
			params = append(params, &domain.OrderListExportParams{
				Filter: *filter,
				Pager:  *upagination.New(page, domain.OrderListExportSingleMaxSize),
			})
		}

		user := domain.FromAdminUserContext(ctx)

		fileName := fmt.Sprintf(
			"%s-%s_%d_订单列表.xlsx",
			req.CreatedAtGte.ToTime().Format(time.DateOnly),
			req.CreatedAtLte.ToTime().Format(time.DateOnly),
			time.Now().Unix(),
		)

		submitParams, err := domain.BuildDataExportSubmitParams(0, domain.DataExportTypeOrderListExport, params, fileName, user)
		if err != nil {
			c.Error(fmt.Errorf("failed to build data export submit params: %w", err))
			return
		}

		if _, err = h.DataExportInteractor.Submit(ctx, submitParams...); err != nil {
			c.Error(fmt.Errorf("failed to submit data export: %w", err))
			return
		}

		response.Ok(c, nil)
	}
}
