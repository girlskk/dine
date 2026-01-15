package handler

import (
	"errors"
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

// AdditionalFeeHandler 费用管理-附加费管理
// No interface assertion needed; handler delegates to interactor.
type AdditionalFeeHandler struct {
	AdditionalFeeInteractor domain.AdditionalFeeInteractor
}

func NewAdditionalFeeHandler(interactor domain.AdditionalFeeInteractor) *AdditionalFeeHandler {
	return &AdditionalFeeHandler{AdditionalFeeInteractor: interactor}
}

func (h *AdditionalFeeHandler) Routes(r gin.IRouter) {
	r = r.Group("/additional_fee")
	r.GET("/:id", h.Get())
	r.GET("", h.List())
}

// Get 获取附加费详情
//
//	@Tags		费用管理-附加费管理
//	@Security	BearerAuth
//	@Summary	获取附加费详情
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"附加费ID"
//	@Success	200	{object}	response.Response{data=domain.AdditionalFee}
//	@Router		/additional_fee/{id} [get]
func (h *AdditionalFeeHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("AdditionalFeeHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromFrontendUserContext(ctx)
		fee, err := h.AdditionalFeeInteractor.GetAdditionalFee(ctx, id, user)
		if err != nil {
			if errors.Is(err, domain.ErrAdditionalFeeNotExists) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.AdditionalFeeNotExists, err))
				return
			}
			err = fmt.Errorf("failed to get additional fee: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, fee)
	}
}

// List 获取附加费列表
//
//	@Tags		费用管理-附加费管理
//	@Security	BearerAuth
//	@Summary	获取附加费列表
//	@Accept		json
//	@Produce	json
//	@Param		data	query		types.AdditionalFeeListReq	true	"附加费列表查询参数"
//	@Success	200		{object}	response.Response{data=types.AdditionalFeeListResp}
//	@Router		/additional_fee [get]
func (h *AdditionalFeeHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("AdditionalFeeHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AdditionalFeeListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromFrontendUserContext(ctx)
		filter := &domain.AdditionalFeeListFilter{
			MerchantID: user.MerchantID,
			FeeType:    domain.AdditionalFeeTypeStore,
			Enabled:    req.Enabled,
		}

		pager := upagination.New(1, upagination.MaxSize)
		fees, total, err := h.AdditionalFeeInteractor.GetAdditionalFees(ctx, pager, filter)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to list additional fees: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.AdditionalFeeListResp{AdditionalFees: fees, Total: total})
	}
}
