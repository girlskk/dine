package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/store/types"
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
		user := domain.FromStoreUserContext(ctx)

		params := domain.ProfitDistributionBillSearchParams{
			MerchantID:    user.MerchantID,
			StoreIDs:      []uuid.UUID{user.StoreID}, // 门店端只能查询自己的账单
			BillStartDate: nil,
			BillEndDate:   nil,
			Status:        domain.ProfitDistributionBillStatus(req.Status),
		}

		// 处理日期字符串转换
		if req.BillStartDate != "" {
			startAt, err := time.ParseInLocation(time.DateOnly, req.BillStartDate, time.Local)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			params.BillStartDate = &startAt
		}
		if req.BillEndDate != "" {
			endDate, err := time.ParseInLocation(time.DateOnly, req.BillEndDate, time.Local)
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
