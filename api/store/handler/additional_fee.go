package handler

import (
	"errors"
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
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("/:id", h.Get())
	r.GET("", h.List())
	r.PUT("/:id/enable", h.Enable())
	r.PUT("/:id/disable", h.Disable())
}

// Create 创建附加费
//
//	@Tags		费用管理-附加费管理
//	@Security	BearerAuth
//	@Summary	创建附加费
//	@Accept		json
//	@Produce	json
//	@Param		data	body	types.AdditionalFeeCreateReq	true	"请求信息"
//	@Success	200		"No Content"
//	@Router		/additional_fee [post]
func (h *AdditionalFeeHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("AdditionalFeeHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AdditionalFeeCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		fee := &domain.AdditionalFee{
			Name:                req.Name,
			FeeType:             domain.AdditionalFeeTypeStore,
			FeeCategory:         req.FeeCategory,
			ChargeMode:          req.ChargeMode,
			FeeValue:            req.FeeValue,
			IncludeInReceivable: req.IncludeInReceivable,
			Taxable:             req.Taxable,
			DiscountScope:       req.DiscountScope,
			OrderChannels:       req.OrderChannels,
			DiningWays:          req.DiningWays,
			Enabled:             req.Enabled,
			SortOrder:           req.SortOrder,
			MerchantID:          user.MerchantID,
			StoreID:             user.StoreID,
		}

		if err := h.AdditionalFeeInteractor.Create(ctx, fee, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新附加费
//
//	@Tags		费用管理-附加费管理
//	@Security	BearerAuth
//	@Summary	更新附加费
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string							true	"附加费ID"
//	@Param		data	body	types.AdditionalFeeUpdateReq	true	"请求信息"
//	@Success	200		"No Content"
//	@Router		/additional_fee/{id} [put]
func (h *AdditionalFeeHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("AdditionalFeeHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.AdditionalFeeUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		fee := &domain.AdditionalFee{
			ID:                  id,
			Name:                req.Name,
			FeeType:             domain.AdditionalFeeTypeStore,
			FeeCategory:         req.FeeCategory,
			ChargeMode:          req.ChargeMode,
			FeeValue:            req.FeeValue,
			IncludeInReceivable: req.IncludeInReceivable,
			Taxable:             req.Taxable,
			DiscountScope:       req.DiscountScope,
			OrderChannels:       req.OrderChannels,
			DiningWays:          req.DiningWays,
			Enabled:             req.Enabled,
			SortOrder:           req.SortOrder,
		}

		if err := h.AdditionalFeeInteractor.Update(ctx, fee, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Delete 删除附加费
//
//	@Tags		费用管理-附加费管理
//	@Security	BearerAuth
//	@Summary	删除附加费
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"附加费ID"
//	@Success	200	"No Content"
//	@Router		/additional_fee/{id} [delete]
func (h *AdditionalFeeHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("AdditionalFeeHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		if err := h.AdditionalFeeInteractor.Delete(ctx, id, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
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

		user := domain.FromStoreUserContext(ctx)
		fee, err := h.AdditionalFeeInteractor.GetAdditionalFee(ctx, id, user)
		if err != nil {
			c.Error(h.checkErr(err))
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

		pager := req.RequestPagination.ToPagination()
		user := domain.FromStoreUserContext(ctx)
		filter := &domain.AdditionalFeeListFilter{
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
			Name:       req.Name,
			FeeType:    domain.AdditionalFeeTypeStore,
			Enabled:    req.Enabled,
		}

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

// Enable 启用附加费
//
//	@Tags			费用管理-附加费管理
//	@Security		BearerAuth
//	@Summary		启用附加费
//	@Description	将附加费置为启用
//	@Produce		json
//	@Param			id	path	string	true	"附加费ID"
//	@Success		200	"No Content"
//	@Router			/additional_fee/{id}/enable [put]
func (h *AdditionalFeeHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("AdditionalFeeHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		fee := &domain.AdditionalFee{ID: id, Enabled: true}
		if err := h.AdditionalFeeInteractor.AdditionalFeeSimpleUpdate(ctx, domain.AdditionalFeeSimpleUpdateTypeEnabled, fee, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 禁用附加费
//
//	@Tags			费用管理-附加费管理
//	@Security		BearerAuth
//	@Summary		禁用附加费
//	@Description	将附加费置为禁用
//	@Produce		json
//	@Param			id	path	string	true	"附加费ID"
//	@Success		200	"No Content"
//	@Router			/additional_fee/{id}/disable [put]
func (h *AdditionalFeeHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("AdditionalFeeHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		fee := &domain.AdditionalFee{ID: id, Enabled: false}
		if err := h.AdditionalFeeInteractor.AdditionalFeeSimpleUpdate(ctx, domain.AdditionalFeeSimpleUpdateTypeEnabled, fee, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}
func (h *AdditionalFeeHandler) checkErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrAdditionalFeeNotExists):
		return errorx.New(http.StatusBadRequest, errcode.AdditinalFeeNotExists, err)
	case errors.Is(err, domain.ErrAdditionalFeeNameExists):
		return errorx.New(http.StatusConflict, errcode.AdditionalFeeNameExists, err)
	case domain.IsNotFound(err):
		return errorx.New(http.StatusNotFound, errcode.NotFound, err)
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return fmt.Errorf("additional fee handler error: %w", err)
	}
}
