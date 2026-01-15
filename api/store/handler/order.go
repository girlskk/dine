package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/store/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

// OrderHandler 订单数据分析
type OrderHandler struct {
	OrderInteractor domain.OrderInteractor
}

func NewOrderHandler(orderInteractor domain.OrderInteractor) *OrderHandler {
	return &OrderHandler{
		OrderInteractor: orderInteractor,
	}
}

func (h *OrderHandler) Routes(r gin.IRouter) {
	r = r.Group("/order")
	r.GET("/sales-report", h.SalesReport())
	r.GET("/product-sales-summary", h.ProductSalesSummary())
	r.GET("", h.List())
	r.GET("/:id", h.Get())
	r.GET("/product-sales-detail", h.ProductSalesDetail())
}

func (h *OrderHandler) NoAuths() []string {
	return []string{}
}

// SalesReport 销售报表
//
//	@Tags		数据分析
//	@Summary	销售汇总表
//	@Accept		json
//	@Produce	json
//	@Param		merchant_id			query		string					true	"品牌商ID"
//	@Param		store_id			query		string					true	"门店ID"
//	@Param		business_date_start	query		string					true	"营业日开始"
//	@Param		business_date_end	query		string					true	"营业日结束"
//	@Param		page				query		int						false	"页码"
//	@Param		size				query		int						false	"每页数量"
//	@Success	200					{object}	types.SalesReportResp	"成功"
//	@Router		/order/sales-report [get]
func (h *OrderHandler) SalesReport() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.SalesReport")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.SalesReportReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		user := domain.FromStoreUserContext(ctx)
		params := domain.OrderSalesReportParams{
			MerchantID:        user.MerchantID,
			StoreIDs:          []uuid.UUID{user.StoreID},
			BusinessDateStart: req.BusinessDateStart,
			BusinessDateEnd:   req.BusinessDateEnd,
			Page:              req.Page,
			Size:              req.Size,
		}

		items, total, err := h.OrderInteractor.SalesReport(ctx, params)
		if err != nil {
			c.Error(fmt.Errorf("failed to get sales report: %w", err))
			return
		}

		p := req.ToPagination()
		p.SetTotal(total)

		response.Ok(c, &types.SalesReportResp{
			Items:      items,
			Pagination: p,
		})
	}
}

// ProductSalesSummary 商品销售汇总
//
//	@Tags		数据分析
//	@Summary	商品销售汇总表
//	@Accept		json
//	@Produce	json
//	@Param		business_date_start	query		string							true	"营业日开始"
//	@Param		business_date_end	query		string							true	"营业日结束"
//	@Param		order_channel		query		string							false	"订单来源"
//	@Param		category_id			query		string							false	"商品分类ID"
//	@Param		product_name		query		string							false	"商品名称（模糊搜索）"
//	@Param		product_type		query		string							false	"商品类型：normal/set_meal"
//	@Param		page				query		int								false	"页码"
//	@Param		size				query		int								false	"每页数量"
//	@Success	200					{object}	types.ProductSalesSummaryResp	"成功"
//	@Router		/order/product-sales-summary [get]
func (h *OrderHandler) ProductSalesSummary() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.ProductSalesSummary")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductSalesSummaryReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		var categoryID uuid.UUID
		if req.CategoryID != "" {
			categoryID, _ = uuid.Parse(req.CategoryID)
		}
		params := domain.ProductSalesSummaryParams{
			MerchantID:        user.MerchantID,
			StoreIDs:          []uuid.UUID{user.StoreID},
			BusinessDateStart: req.BusinessDateStart,
			BusinessDateEnd:   req.BusinessDateEnd,
			OrderChannel:      domain.Channel(req.OrderChannel),
			CategoryID:        categoryID,
			ProductName:       req.ProductName,
			ProductType:       domain.ProductType(req.ProductType),
			Page:              req.Page,
			Size:              req.Size,
		}

		items, total, err := h.OrderInteractor.ProductSalesSummary(ctx, params)
		if err != nil {
			c.Error(fmt.Errorf("failed to get product sales summary: %w", err))
			return
		}

		p := req.ToPagination()
		p.SetTotal(total)

		response.Ok(c, &types.ProductSalesSummaryResp{
			Items:      items,
			Pagination: p,
		})
	}
}

// ProductSalesDetail 商品销售明细
//
//	@Tags		数据分析
//	@Summary	商品销售明细表
//	@Accept		json
//	@Produce	json
//	@Param		business_date_start	query		string							true	"营业日开始"
//	@Param		business_date_end	query		string							true	"营业日结束"
//	@Param		order_channel		query		string							false	"订单来源"
//	@Param		category_id			query		string							false	"商品分类ID"
//	@Param		product_name		query		string							false	"商品名称（模糊搜索）"
//	@Param		product_type		query		string							false	"商品类型：normal/set_meal"
//	@Param		order_no			query		string							false	"订单号"
//	@Param		page				query		int								false	"页码"
//	@Param		size				query		int								false	"每页数量"
//	@Success	200					{object}	types.ProductSalesDetailResp	"成功"
//	@Router		/order/product-sales-detail [get]
func (h *OrderHandler) ProductSalesDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.ProductSalesDetail")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductSalesDetailReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		var categoryID uuid.UUID
		if req.CategoryID != "" {
			categoryID, _ = uuid.Parse(req.CategoryID)
		}
		params := domain.ProductSalesDetailParams{
			MerchantID:        user.MerchantID,
			StoreIDs:          []uuid.UUID{user.StoreID},
			BusinessDateStart: req.BusinessDateStart,
			BusinessDateEnd:   req.BusinessDateEnd,
			OrderChannel:      domain.Channel(req.OrderChannel),
			CategoryID:        categoryID,
			ProductName:       req.ProductName,
			ProductType:       domain.ProductType(req.ProductType),
			OrderNo:           req.OrderNo,
			Page:              req.Page,
			Size:              req.Size,
		}

		items, total, err := h.OrderInteractor.ProductSalesDetail(ctx, params)
		if err != nil {
			c.Error(fmt.Errorf("failed to get product sales detail: %w", err))
			return
		}

		p := req.ToPagination()
		p.SetTotal(total)

		response.Ok(c, &types.ProductSalesDetailResp{
			Items:      items,
			Pagination: p,
		})
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

// List
//
//	@Tags		订单
//	@Security	BearerAuth
//	@Summary	获取订单列表
//	@Accept		json
//	@Produce	json
//	@Param		business_date_start	query		string				false	"营业日开始"
//	@Param		business_date_end	query		string				false	"营业日结束"
//	@Param		order_no			query		string				false	"订单号"
//	@Param		order_type			query		string				false	"订单类型"	Enums(SALE,REFUND,PARTIAL_REFUND)
//	@Param		order_status		query		string				false	"订单状态"	Enums(PLACED,COMPLETED,CANCELLED)
//	@Param		payment_status		query		string				false	"支付状态"	Enums(UNPAID,PAYING,PAID,REFUNDED)
//	@Param		page				query		int					false	"页码"
//	@Param		size				query		int					false	"每页数量"
//	@Success	200					{object}	types.ListOrderResp	"成功"
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

		user := domain.FromStoreUserContext(ctx)
		params := domain.OrderListParams{
			MerchantID:        user.MerchantID,
			StoreID:           user.StoreID,
			BusinessDateStart: req.BusinessDateStart,
			BusinessDateEnd:   req.BusinessDateEnd,
			OrderNo:           req.OrderNo,
			OrderType:         domain.OrderType(req.OrderType),
			OrderStatus:       domain.OrderStatus(req.OrderStatus),
			PaymentStatus:     domain.PaymentStatus(req.PaymentStatus),
			Page:              req.Page,
			Size:              req.Size,
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
