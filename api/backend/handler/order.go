package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
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
}

func (h *OrderHandler) NoAuths() []string {
	return []string{}
}

// SalesReport 销售报表
//
//	@Tags		数据分析
//	@Security	BearerAuth
//	@Summary	销售汇总表
//	@Accept		json
//	@Produce	json
//	@Param		merchant_id			query		string					true	"品牌商ID"
//	@Param		store_ids			query		string					false	"门店ID列表（逗号分隔）"
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

		merchantID, err := uuid.Parse(req.MerchantID)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams,
				fmt.Errorf("invalid merchant_id: %w", err)))
			return
		}

		var storeIDs []uuid.UUID
		if req.StoreIDs != "" {
			for _, s := range strings.Split(req.StoreIDs, ",") {
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				id, err := uuid.Parse(s)
				if err != nil {
					c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams,
						fmt.Errorf("invalid store_id: %w", err)))
					return
				}
				storeIDs = append(storeIDs, id)
			}
		}

		params := domain.OrderSalesReportParams{
			MerchantID:        merchantID,
			StoreIDs:          storeIDs,
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
