package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/types"
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
	r.POST("/list", h.List())
	r.POST("/approve", h.Approve())
	r.POST("/reject", h.Reject())
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
		page := upagination.New(req.Page, req.Size)
		params := domain.StoreWithdrawSearchParams{
			StartAt: req.StartAt.ToPtrStartOfDay(),
			EndAt:   req.EndAt.ToPtrEndOfDay(),
			StoreID: req.StoreID,
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

// Approve 审批提现单
//
//	@Tags		账户中心-提现
//	@Summary	审批提现单
//	@Security	BearerAuth
//	@Param		data	body	types.StoreWithdrawIDReq	true	"审批请求"
//	@Success	200
//	@Router		/store-withdraw/approve [post]
func (h *StoreWithdrawHandler) Approve() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawHandler.Approve")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreWithdrawHandler.Approve")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreWithdrawIDReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		// 审批提现单
		err := h.StoreWithdrawInteractor.Approve(ctx, req.ID)
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

// Reject 驳回提现单
//
//	@Tags		账户中心-提现
//	@Summary	驳回提现单
//	@Security	BearerAuth
//	@Param		data	body	types.StoreWithdrawIDReq	true	"驳回请求"
//	@Success	200
//	@Router		/store-withdraw/reject [post]
func (h *StoreWithdrawHandler) Reject() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreWithdrawHandler.Reject")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreWithdrawHandler.Reject")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreWithdrawIDReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		// 驳回提现单
		err := h.StoreWithdrawInteractor.Reject(ctx, req.ID)
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
