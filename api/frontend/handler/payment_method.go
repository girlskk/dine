package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type PaymentMethodHandler struct {
	PaymentMethodInteractor domain.PaymentMethodInteractor
}

func NewPaymentMethodHandler(PaymentMethodInteractor domain.PaymentMethodInteractor) *PaymentMethodHandler {
	return &PaymentMethodHandler{
		PaymentMethodInteractor: PaymentMethodInteractor,
	}
}

func (h *PaymentMethodHandler) Routes(r gin.IRouter) {
	r = r.Group("payment/method")
	r.GET("", h.List())
}

func (h *PaymentMethodHandler) NoAuths() []string {
	return []string{}
}

// List
//
//	@Tags		结算方式管理
//	@Security	BearerAuth
//	@Summary	查询结算方式列表
//	@Param		data	query		types.PaymentMethodListReq		true	"结算方式列表查询参数"
//	@Success	200		{object}	domain.PaymentMethodSearchRes	"成功"
//	@Router		/payment/method [get]
func (h *PaymentMethodHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("PaymentMethodHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.PaymentMethodListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		params := domain.PaymentMethodSearchParams{}
		if req.StoreID != "" {
			storeID, err := uuid.Parse(req.StoreID)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			params.StoreID = storeID
		}

		user := domain.FromFrontendContext(ctx)
		if user != nil {
			params.MerchantID = user.GetMerchantID()
		}
		page := upagination.New(req.Page, req.Size)
		res, err := h.PaymentMethodInteractor.PagedListBySearch(ctx, page, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to list paymentMethods: %w", err)
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}
