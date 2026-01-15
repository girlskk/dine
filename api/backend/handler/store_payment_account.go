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

type StorePaymentAccountHandler struct {
	StorePaymentAccountInteractor domain.StorePaymentAccountInteractor
}

func NewStorePaymentAccountHandler(storePaymentAccountInteractor domain.StorePaymentAccountInteractor) *StorePaymentAccountHandler {
	return &StorePaymentAccountHandler{
		StorePaymentAccountInteractor: storePaymentAccountInteractor,
	}
}

func (h *StorePaymentAccountHandler) Routes(r gin.IRouter) {
	r = r.Group("payment/store-account")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("", h.List())
}

func (h *StorePaymentAccountHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		门店收款账户
//	@Security	BearerAuth
//	@Summary	创建门店收款账户
//	@Param		data	body	types.StorePaymentAccountCreateReq	true	"请求信息"
//	@Success	200
//	@Router		/payment/store-account [post]
func (h *StorePaymentAccountHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StorePaymentAccountHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StorePaymentAccountCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		account := &domain.StorePaymentAccount{
			ID:               uuid.New(),
			MerchantID:       user.MerchantID,
			StoreID:          req.StoreID,
			PaymentAccountID: req.PaymentAccountID,
			MerchantNumber:   req.MerchantNumber,
		}

		err := h.StorePaymentAccountInteractor.Create(ctx, account, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create store payment account: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update
//
//	@Tags		门店收款账户
//	@Security	BearerAuth
//	@Summary	更新门店收款账户
//	@Param		id		path	string								true	"门店收款账户ID"
//	@Param		data	body	types.StorePaymentAccountUpdateReq	true	"请求信息"
//	@Success	200
//	@Router		/payment/store-account/{id} [put]
func (h *StorePaymentAccountHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StorePaymentAccountHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("invalid store payment account id: %w", err)))
			return
		}

		var req types.StorePaymentAccountUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		account := &domain.StorePaymentAccount{
			ID:             id,
			MerchantID:     user.MerchantID,
			MerchantNumber: req.MerchantNumber,
		}

		err = h.StorePaymentAccountInteractor.Update(ctx, account, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update store payment account: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete
//
//	@Tags		门店收款账户
//	@Security	BearerAuth
//	@Summary	删除门店收款账户
//	@Param		id	path	string	true	"门店收款账户ID"
//	@Success	200
//	@Router		/store/payment/account/{id} [delete]
func (h *StorePaymentAccountHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StorePaymentAccountHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("invalid store payment account id: %w", err)))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		err = h.StorePaymentAccountInteractor.Delete(ctx, id, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to delete store payment account: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// List
//
//	@Tags		门店收款账户
//	@Security	BearerAuth
//	@Summary	查询门店收款账户列表
//	@Param		data	query		types.StorePaymentAccountListReq	true	"请求信息"
//	@Success	200		{object}	domain.StorePaymentAccountSearchRes	"成功"
//	@Router		/payment/store-account [get]
func (h *StorePaymentAccountHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StorePaymentAccountHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StorePaymentAccountListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		page := upagination.New(req.Page, req.Size)
		user := domain.FromBackendUserContext(ctx)

		params := domain.StorePaymentAccountSearchParams{
			MerchantID:     user.MerchantID,
			MerchantName:   req.MerchantName,
			CreatedAtStart: nil,
			CreatedAtEnd:   nil,
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

		var startAt, endAt time.Time
		var err error
		if req.CreatedAtStart != "" {
			startAt, err = time.ParseInLocation(time.DateOnly, req.CreatedAtStart, time.Local)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			params.CreatedAtStart = &startAt
		}
		if req.CreatedAtEnd != "" {
			endAt, err = time.ParseInLocation(time.DateOnly, req.CreatedAtEnd, time.Local)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			params.CreatedAtEnd = &endAt
		}

		res, err := h.StorePaymentAccountInteractor.PagedListBySearch(ctx, page, params, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to list store payment accounts: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, res)
	}
}
