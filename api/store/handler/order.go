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
