package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type ProfitDistributionBillHandler struct {
	ProfitDistributionBillInteractor domain.ProfitDistributionBillInteractor
}

func NewProfitDistributionBillHandler(profitDistributionBillInteractor domain.ProfitDistributionBillInteractor) *ProfitDistributionBillHandler {
	return &ProfitDistributionBillHandler{
		ProfitDistributionBillInteractor: profitDistributionBillInteractor,
	}
}

func (h *ProfitDistributionBillHandler) Routes(r gin.IRouter) {
	r = r.Group("profit/distribution/bill")
	r.GET("", h.List())
	r.POST("/:id/pay", h.Pay())
}

func (h *ProfitDistributionBillHandler) NoAuths() []string {
	return []string{}
}

// List
//
//	@Tags		分账账单
//	@Security	BearerAuth
//	@Summary	查询分账账单列表
//	@Param		data	query		types.ProfitDistributionBillListReq		true	"请求信息"
//	@Success	200		{object}	domain.ProfitDistributionBillSearchRes	"成功"
//	@Router		/profit/distribution/bill [get]
func (h *ProfitDistributionBillHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProfitDistributionBillHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProfitDistributionBillListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		page := upagination.New(req.Page, req.Size)
		user := domain.FromBackendUserContext(ctx)

		params := domain.ProfitDistributionBillSearchParams{
			MerchantID:    user.MerchantID,
			BillStartDate: nil,
			BillEndDate:   nil,
			Status:        domain.ProfitDistributionBillStatus(req.Status),
		}

		// 转换门店ID列表
		if len(req.StoreIDs) > 0 {
			storeIDs := make([]uuid.UUID, 0, len(req.StoreIDs))
			for _, storeIDStr := range req.StoreIDs {
				if storeIDStr == "" {
					continue
				}
				storeID, err := uuid.Parse(storeIDStr)
				if err != nil {
					c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
					return
				}
				storeIDs = append(storeIDs, storeID)
			}
			params.StoreIDs = storeIDs
		}

		// 处理日期字符串转换（如果需要）
		if req.BillStartDate != "" {
			startAt, err := time.Parse(time.DateOnly, req.BillStartDate)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			params.BillStartDate = &startAt
		}
		if req.BillEndDate != "" {
			endDate, err := time.Parse(time.DateOnly, req.BillEndDate)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			params.BillEndDate = &endDate
		}

		res, err := h.ProfitDistributionBillInteractor.PagedListBySearch(ctx, page, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to list profit distribution bills: %w", err)
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}

// Pay
//
//	@Tags		分账账单
//	@Security	BearerAuth
//	@Summary	打款分账账单
//	@Param		id		path	string								true	"分账账单ID"
//	@Param		data	body	types.ProfitDistributionBillPayReq	true	"打款信息"
//	@Success	200
//	@Router		/profit/distribution/bill/{id}/pay [post]
func (h *ProfitDistributionBillHandler) Pay() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProfitDistributionBillHandler.Pay")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取分账账单ID
		idStr := c.Param("id")
		billID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		// 从请求体获取打款信息
		var req types.ProfitDistributionBillPayReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		err = h.ProfitDistributionBillInteractor.Pay(ctx, billID, req.PaymentAmount, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to pay profit distribution bill: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}
