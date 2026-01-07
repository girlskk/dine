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

type PaymentAccountHandler struct {
	PaymentAccountInteractor domain.PaymentAccountInteractor
}

func NewPaymentAccountHandler(paymentAccountInteractor domain.PaymentAccountInteractor) *PaymentAccountHandler {
	return &PaymentAccountHandler{
		PaymentAccountInteractor: paymentAccountInteractor,
	}
}

func (h *PaymentAccountHandler) Routes(r gin.IRouter) {
	r = r.Group("payment/account")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("", h.List())
}

func (h *PaymentAccountHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		收款账户
//	@Security	BearerAuth
//	@Summary	创建收款账户
//	@Param		data	body	types.PaymentAccountCreateReq	true	"请求信息"
//	@Success	200
//	@Router		/payment/account [post]
func (h *PaymentAccountHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("PaymentAccountHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.PaymentAccountCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		account := &domain.PaymentAccount{
			ID:             uuid.New(),
			MerchantID:     user.MerchantID,
			Channel:        req.Channel,
			MerchantNumber: req.MerchantNumber,
			MerchantName:   req.MerchantName,
		}

		err := h.PaymentAccountInteractor.Create(ctx, account)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create payment account: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update
//
//	@Tags		收款账户
//	@Security	BearerAuth
//	@Summary	更新收款账户
//	@Param		id		path	string							true	"收款账户ID"
//	@Param		data	body	types.PaymentAccountUpdateReq	true	"请求信息"
//	@Success	200
//	@Router		/payment/account/{id} [put]
func (h *PaymentAccountHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("PaymentAccountHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("invalid payment account id: %w", err)))
			return
		}

		var req types.PaymentAccountUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		account := &domain.PaymentAccount{
			ID:             id,
			Channel:        req.Channel,
			MerchantNumber: req.MerchantNumber,
			MerchantName:   req.MerchantName,
		}

		err = h.PaymentAccountInteractor.Update(ctx, account, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update payment account: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete
//
//	@Tags		收款账户
//	@Security	BearerAuth
//	@Summary	删除收款账户
//	@Param		id	path	string	true	"收款账户ID"
//	@Success	200
//	@Router		/payment/account/{id} [delete]
func (h *PaymentAccountHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("PaymentAccountHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("invalid payment account id: %w", err)))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		err = h.PaymentAccountInteractor.Delete(ctx, id, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to delete payment account: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// List
//
//	@Tags		收款账户
//	@Security	BearerAuth
//	@Summary	查询收款账户列表
//	@Param		data	query		types.PaymentAccountListReq		true	"请求信息"
//	@Success	200		{object}	domain.PaymentAccountSearchRes	"成功"
//	@Router		/payment/account [get]
func (h *PaymentAccountHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("PaymentAccountHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.PaymentAccountListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		page := upagination.New(req.Page, req.Size)
		user := domain.FromBackendUserContext(ctx)

		params := domain.PaymentAccountSearchParams{
			MerchantID:     user.MerchantID,
			Channel:        req.Channel,
			MerchantName:   req.MerchantName,
			CreatedAtStart: nil,
			CreatedAtEnd:   nil,
		}

		var startAt, endAt time.Time
		var err error
		if req.CreatedAtStart != "" {
			startAt, err = time.Parse(time.DateOnly, req.CreatedAtStart)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			params.CreatedAtStart = &startAt
		}
		if req.CreatedAtEnd != "" {
			endAt, err = time.Parse(time.DateOnly, req.CreatedAtEnd)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			params.CreatedAtEnd = &endAt
		}

		res, err := h.PaymentAccountInteractor.PagedListBySearch(ctx, page, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to list payment accounts: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, res)
	}
}
