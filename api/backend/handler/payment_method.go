package handler

import (
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
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("/:id", h.GetDetail())
	r.GET("", h.List())
	r.GET("/stat", h.Stat())
}

func (h *PaymentMethodHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		结算方式管理
//	@Security	BearerAuth
//	@Summary	创建结算方式
//	@Param		data	body	types.PaymentMethodCreateReq	true	"请求信息"
//	@Success	200
//	@Router		/payment/method [post]
func (h *PaymentMethodHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("PaymentMethodHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.PaymentMethodCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		user := domain.FromBackendUserContext(ctx)
		// 构建 domain.PaymentMethod
		paymentMethod := &domain.PaymentMethod{
			ID:               uuid.New(),
			MerchantID:       user.MerchantID,
			Name:             req.Name,
			AccountingRule:   req.AccountingRule,
			PaymentType:      req.PaymentType,
			FeeRate:          req.FeeRate,
			InvoiceRule:      req.InvoiceRule,
			CashDrawerStatus: req.CashDrawerStatus,
			DisplayChannels:  req.DisplayChannels,
			Source:           req.Source,
			Status:           req.Status,
		}
		err := h.PaymentMethodInteractor.Create(ctx, paymentMethod)
		if err != nil {
			if errors.Is(err, domain.ErrPaymentMethodNotExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create paymentMethod: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update
//
//	@Tags		结算方式管理
//	@Security	BearerAuth
//	@Summary	更新结算方式
//	@Param		id		path	string							true	"结算方式ID"
//	@Param		data	body	types.PaymentMethodUpdateReq	true	"请求信息"
//	@Success	200
//	@Router		/payment/method/{id} [put]
func (h *PaymentMethodHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("PaymentMethodHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)
		idStr := c.Param("id")
		paymentMethodID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.PaymentMethodUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		// 构建 domain.PaymentMethod
		paymentMethod := &domain.PaymentMethod{
			ID:               paymentMethodID,
			MerchantID:       user.MerchantID,
			Name:             req.Name,
			AccountingRule:   req.AccountingRule,
			PaymentType:      req.PaymentType,
			FeeRate:          req.FeeRate,
			InvoiceRule:      req.InvoiceRule,
			CashDrawerStatus: req.CashDrawerStatus,
			DisplayChannels:  req.DisplayChannels,
			Status:           req.Status,
		}

		err = h.PaymentMethodInteractor.Update(ctx, paymentMethod, user)
		if err != nil {
			if errors.Is(err, domain.ErrPaymentMethodNotExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update paymentMethod: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, paymentMethod)
	}
}

// Delete
//
//	@Tags		结算方式管理
//	@Security	BearerAuth
//	@Summary	删除结算方式
//	@Param		id	path	string	true	"结算方式ID"
//	@Success	200
//	@Router		/payment/method/{id} [delete]
func (h *PaymentMethodHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("PaymentMethodHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)
		idStr := c.Param("id")
		paymentMethodID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		err = h.PaymentMethodInteractor.Delete(ctx, paymentMethodID, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to delete paymentMethod: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// GetDetail
//
//	@Tags		结算方式管理
//	@Security	BearerAuth
//	@Summary	获取结算方式详情
//	@Param		id	path		string					true	"结算方式ID"
//	@Success	200	{object}	domain.PaymentMethod	"成功"
//	@Router		/payment/method/{id} [get]
func (h *PaymentMethodHandler) GetDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("PaymentMethodHandler.GetDetail")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)
		idStr := c.Param("id")
		paymentMethodID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		user := domain.FromBackendUserContext(ctx)
		paymentMethod, err := h.PaymentMethodInteractor.GetDetail(ctx, paymentMethodID, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to get paymentMethod detail: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, paymentMethod)
	}
}

// List
//
//	@Tags		结算方式管理
//	@Security	BearerAuth
//	@Summary	查询结算方式列表
//	@Param		name	query		string							false	"结算方式名称（模糊匹配）"
//	@Param		source	query		string							false	"来源:brand-品牌,store-门店,system-系统"
//	@Param		page	query		int								false	"页码"
//	@Param		size	query		int								false	"每页数量"
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

		page := upagination.New(req.Page, req.Size)
		user := domain.FromBackendUserContext(ctx)
		params := domain.PaymentMethodSearchParams{
			MerchantID: user.MerchantID,
			Name:       req.Name,
			Source:     req.Source,
		}
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

// Stat
//
//	@Tags		结算方式管理
//	@Security	BearerAuth
//	@Summary	统计各个结算分类对应的结算方式数量
//	@Param		name	query		string						false	"结算方式名称（模糊匹配）"
//	@Param		source	query		string						false	"来源:brand-品牌,store-门店,system-系统"
//	@Success	200		{object}	domain.PaymentMethodStatRes	"成功"
//	@Router		/payment/method/stat [get]
func (h *PaymentMethodHandler) Stat() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("PaymentMethodHandler.Stat")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.PaymentMethodStatReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		ctx = domain.NewBackendUserContext(ctx, &domain.BackendUser{
			ID:         uuid.New(),
			MerchantID: uuid.New(),
		})
		user := domain.FromBackendUserContext(ctx)
		params := domain.PaymentMethodStatParams{
			MerchantID: user.MerchantID,
			Name:       req.Name,
			Source:     req.Source,
		}
		res, err := h.PaymentMethodInteractor.Stat(ctx, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to query stat count : %w", err)
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}
