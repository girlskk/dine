package handler

import (
	"context"
	"errors"
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
	"go.uber.org/fx"
)

// TaxFeeHandler 费用管理-税费管理
type TaxFeeHandler struct {
	TaxFeeInteractor domain.TaxFeeInteractor
	taxSeq           domain.IncrSequence
}

type TaxFeeHandlerParams struct {
	fx.In
	TaxFeeInteractor domain.TaxFeeInteractor
	TaxSeq           domain.IncrSequence `name:"backend_tax_seq"`
}

func NewTaxFeeHandler(p TaxFeeHandlerParams) *TaxFeeHandler {
	return &TaxFeeHandler{TaxFeeInteractor: p.TaxFeeInteractor, taxSeq: p.TaxSeq}
}

func (h *TaxFeeHandler) Routes(r gin.IRouter) {
	r = r.Group("/tax_fee")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("/:id", h.Get())
	r.GET("", h.List())
	r.PUT("/:id/Enable", h.Enable())
	r.PUT("/:id/disable", h.Disable())
}

// Create 创建税费
//
//	@Tags		费用管理-税费管理
//	@Security	BearerAuth
//	@Summary	创建税费
//	@Accept		json
//	@Produce	json
//	@Param		data	body	types.TaxFeeCreateReq	true	"请求信息"
//	@Success	200		"No Content"
//	@Router		/tax_fee [post]
func (h *TaxFeeHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("TaxFeeHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.TaxFeeCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		taxCode, err := h.generateTaxCode(ctx)
		if err != nil {
			err = fmt.Errorf("failed to generate tax code: %w", err)
			c.Error(err)
			return
		}
		user := domain.FromBackendUserContext(ctx)
		fee := &domain.TaxFee{
			Name:        req.Name,
			TaxFeeType:  domain.TaxFeeTypeMerchant,
			TaxCode:     taxCode,
			TaxRateType: domain.TaxRateTypeUnified,
			TaxRate:     req.TaxRate,
			DefaultTax:  false,
			MerchantID:  user.MerchantID,
		}

		if err := h.TaxFeeInteractor.Create(ctx, fee, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新税费
//
//	@Tags		费用管理-税费管理
//	@Security	BearerAuth
//	@Summary	更新税费
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string					true	"税费ID"
//	@Param		data	body	types.TaxFeeUpdateReq	true	"请求信息"
//	@Success	200		"No Content"
//	@Router		/tax_fee/{id} [put]
func (h *TaxFeeHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("TaxFeeHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.TaxFeeUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		fee := &domain.TaxFee{
			ID:          id,
			Name:        req.Name,
			TaxRateType: domain.TaxRateTypeUnified,
			TaxRate:     req.TaxRate,
			DefaultTax:  false,
		}

		if err := h.TaxFeeInteractor.Update(ctx, fee, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Delete 删除税费
//
//	@Tags		费用管理-税费管理
//	@Security	BearerAuth
//	@Summary	删除税费
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"税费ID"
//	@Success	200	"No Content"
//	@Success	204	"No Content"
//	@Router		/tax_fee/{id} [delete]
func (h *TaxFeeHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("TaxFeeHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		if err := h.TaxFeeInteractor.Delete(ctx, id, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Get 获取税费详情
//
//	@Tags		费用管理-税费管理
//	@Security	BearerAuth
//	@Summary	获取税费详情
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"税费ID"
//	@Success	200	{object}	response.Response{data=domain.TaxFee}
//	@Router		/tax_fee/{id} [get]
func (h *TaxFeeHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("TaxFeeHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		fee, err := h.TaxFeeInteractor.GetTaxFee(ctx, id, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, fee)
	}
}

// List 获取税费列表
//
//	@Tags		费用管理-税费管理
//	@Security	BearerAuth
//	@Summary	获取税费列表
//	@Accept		json
//	@Produce	json
//	@Param		data	query		types.TaxFeeListReq	true	"出品部门列表查询参数"
//	@Success	200		{object}	response.Response{data=types.TaxFeeListResp}
//	@Router		/tax_fee [get]
func (h *TaxFeeHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("TaxFeeHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.TaxFeeListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		pager := req.RequestPagination.ToPagination()
		user := domain.FromBackendUserContext(ctx)
		filter := &domain.TaxFeeListFilter{
			MerchantID: user.MerchantID,
			Name:       req.Name,
			TaxFeeType: domain.TaxFeeTypeMerchant,
		}

		fees, total, err := h.TaxFeeInteractor.GetTaxFees(ctx, pager, filter)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to list tax fees: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.TaxFeeListResp{TaxFees: fees, Total: total})
	}
}

// Enable 启用税费
//
//	@Tags			费用管理-税费管理
//	@Security		BearerAuth
//	@Summary		启用税费
//	@Description	将税费标记为默认
//	@Produce		json
//	@Param			id	path	string	true	"税费ID"
//	@Success		200	"No Content"
//	@Router			/tax_fee/{id}/enable [put]
func (h *TaxFeeHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("TaxFeeHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		fee := &domain.TaxFee{ID: id, DefaultTax: true}
		if err := h.TaxFeeInteractor.TaxFeeSimpleUpdate(ctx, domain.TaxFeeSimpleUpdateFieldDefault, fee, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 取消默认税费
//
//	@Tags			费用管理-税费管理
//	@Security		BearerAuth
//	@Summary		禁用税费
//	@Description	取消默认税费标记
//	@Produce		json
//	@Param			id	path	string	true	"税费ID"
//	@Success		200	"No Content"
//	@Router			/tax_fee/{id}/disable [put]
func (h *TaxFeeHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("TaxFeeHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		fee := &domain.TaxFee{ID: id, DefaultTax: false}
		if err := h.TaxFeeInteractor.TaxFeeSimpleUpdate(ctx, domain.TaxFeeSimpleUpdateFieldDefault, fee, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

func (h *TaxFeeHandler) generateTaxCode(ctx context.Context) (string, error) {
	if h.taxSeq == nil {
		return "", fmt.Errorf("tax fee sequence not initialized")
	}
	seq, err := h.taxSeq.Next(ctx)
	if err != nil {
		return "", err
	}
	return seq, nil
}

func (h *TaxFeeHandler) checkErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrTaxFeeNotExists):
		return errorx.New(http.StatusBadRequest, errcode.TaxFeeNotExists, err)
	case errors.Is(err, domain.ErrTaxFeeNameExists):
		return errorx.New(http.StatusConflict, errcode.TaxFeeNameExists, err)
	case domain.IsNotFound(err):
		return errorx.New(http.StatusNotFound, errcode.NotFound, err)
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return fmt.Errorf("tax fee handler error: %w", err)
	}
}
