package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type StoreWithdrawHandler struct {
	StoreWithdrawInteractor domain.StoreWithdrawInteractor
}

func NewStoreWithdrawHandler(interactor domain.StoreWithdrawInteractor) *StoreWithdrawHandler {
	return &StoreWithdrawHandler{
		StoreWithdrawInteractor: interactor,
	}
}

// Routes 注册路由
func (h *StoreWithdrawHandler) Routes(r gin.IRouter) {
	r = r.Group("/store-withdraw")
	r.POST("/apply", h.Apply())
	r.POST("/commit", h.Commit())
	r.POST("/cancel", h.Cancel())
	r.POST("/delete", h.Delete())
	r.POST("/update", h.Update())
	r.POST("/list", h.List())
}

// Apply 创建提现单
//
//	@Tags		账户中心-提现
//	@Summary	创建提现单
//	@Security	BearerAuth
//	@Param		data	body	types.StoreWithdrawApplyReq	true	"请求参数"
//	@Success	200
//	@Router		/store-withdraw/apply [post]
func (h *StoreWithdrawHandler) Apply() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawHandler.Apply")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreWithdrawHandler.Apply")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreWithdrawApplyReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		user := domain.FromBackendUserContext(ctx)

		withdraw := &domain.StoreWithdraw{
			StoreID:             user.StoreID,
			StoreName:           user.Store.Name,
			Amount:              req.Amount,
			PointWithdrawalRate: user.Store.PointWithdrawalRate,
			ActualAmount:        req.Amount.Sub(req.Amount.Mul(user.Store.PointWithdrawalRate)),
			AccountType:         req.AccountType,
			BankAccount:         req.BankAccount,
			BankCardName:        req.BankCardName,
			BankName:            req.BankName,
			BankBranch:          req.BankBranch,
			InvoiceAmount:       req.InvoiceAmount,
			Status:              domain.StoreWithdrawStatusUncommitted,
		}
		err := h.StoreWithdrawInteractor.Apply(ctx, withdraw)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}
		response.Ok(c, nil)
	}
}

// Commit 提交提现单
//
//	@Tags		账户中心-提现
//	@Summary	提交提现单
//	@Security	BearerAuth
//	@Param		data	body	types.StoreWithdrawIDReq	true	"请求参数"
//	@Success	200
//	@Router		/store-withdraw/commit [post]
func (h *StoreWithdrawHandler) Commit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawHandler.Commit")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreWithdrawHandler.Commit")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreWithdrawIDReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		user := domain.FromBackendUserContext(ctx)

		err := h.StoreWithdrawInteractor.Commit(ctx, req.ID, user.StoreID)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}
		response.Ok(c, nil)
	}
}

// List 提现单列表
//
//	@Tags		账户中心-提现
//	@Summary	提现单列表
//	@Security	BearerAuth
//	@Param		data	body		types.StoreWithdrawListReq	true	"请求参数"
//	@Success	200		{object}	domain.StoreWithdrawSearchRes
//	@Router		/store-withdraw/list [post]
func (h *StoreWithdrawHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreWithdrawHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreWithdrawListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		user := domain.FromBackendUserContext(ctx)
		page := upagination.New(req.Page, req.Size)
		params := domain.StoreWithdrawSearchParams{
			StartAt: req.StartAt.ToPtrStartOfDay(),
			EndAt:   req.EndAt.ToPtrEndOfDay(),
			StoreID: user.StoreID,
			Status:  req.Status,
		}

		res, err := h.StoreWithdrawInteractor.PagedListBySearch(ctx, page, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}

// Cancel 撤回提现单
//
//	@Tags		账户中心-提现
//	@Summary	撤回提现单
//	@Security	BearerAuth
//	@Param		data	body	types.StoreWithdrawIDReq	true	"请求参数"
//	@Success	200
//	@Router		/store-withdraw/cancel [post]
func (h *StoreWithdrawHandler) Cancel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawHandler.Cancel")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreWithdrawHandler.Cancel")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreWithdrawIDReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		user := domain.FromBackendUserContext(ctx)

		err := h.StoreWithdrawInteractor.Cancel(ctx, req.ID, user.StoreID)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}
		response.Ok(c, nil)
	}
}

// Delete 删除提现单
//
//	@Tags		账户中心-提现
//	@Summary	删除提现单
//	@Security	BearerAuth
//	@Param		data	body	types.StoreWithdrawIDReq	true	"请求参数"
//	@Success	200
//	@Router		/store-withdraw/delete [post]
func (h *StoreWithdrawHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawHandler.Delete")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreWithdrawHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreWithdrawIDReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		user := domain.FromBackendUserContext(ctx)

		err := h.StoreWithdrawInteractor.Delete(ctx, req.ID, user.StoreID)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}
		response.Ok(c, nil)
	}
}

// Update 编辑提现单
//
//	@Tags		账户中心-提现
//	@Summary	编辑提现单
//	@Security	BearerAuth
//	@Param		data	body	types.StoreWithdrawUpdateReq	true	"请求参数"
//	@Success	200
//	@Router		/store-withdraw/update [post]
func (h *StoreWithdrawHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawHandler.Update")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreWithdrawHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreWithdrawUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		user := domain.FromBackendUserContext(ctx)

		withdraw := &domain.StoreWithdraw{
			ID:                  req.ID,
			Amount:              req.Amount,
			StoreID:             user.StoreID,
			PointWithdrawalRate: user.Store.PointWithdrawalRate,
			ActualAmount:        req.Amount.Sub(req.Amount.Mul(user.Store.PointWithdrawalRate)),
			AccountType:         req.AccountType,
			BankAccount:         req.BankAccount,
			BankCardName:        req.BankCardName,
			BankName:            req.BankName,
			BankBranch:          req.BankBranch,
			InvoiceAmount:       req.InvoiceAmount,
		}

		err := h.StoreWithdrawInteractor.Update(ctx, withdraw)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}
		response.Ok(c, nil)
	}
}
