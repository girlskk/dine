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

type StoreAccountHandler struct {
	StoreAccountInteractor domain.StoreAccountInteractor
}

func NewStoreAccountHandler(interactor domain.StoreAccountInteractor,
) *StoreAccountHandler {
	return &StoreAccountHandler{
		StoreAccountInteractor: interactor,
	}
}

// Routes 注册路由
func (h *StoreAccountHandler) Routes(r gin.IRouter) {
	r = r.Group("/store-account")
	r.POST("/detail", h.Detail())
	r.POST("/list-transactions", h.ListTransactions())
}

// Detail 获取门店账户详情
//
//	@Tags		账户中心
//	@Summary	获取门店账户详情
//	@Security	BearerAuth
//	@Success	200	{object}	domain.StoreAccount
//	@Router		/store-account/detail [post]
func (h *StoreAccountHandler) Detail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreAccountHandler.Detail")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreAccountHandler.Detail")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromBackendUserContext(ctx)
		detail, err := h.StoreAccountInteractor.GetDetail(ctx, user.StoreID)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, detail)
	}
}

// ListTransactions 门店账户流水列表
//
//	@Tags		账户中心
//	@Summary	门店账户流水列表
//	@Security	BearerAuth
//	@Param		data	body		types.StoreAccountTransactionListReq	true	"请求参数"
//	@Success	200		{object}	domain.StoreAccountTransactionSearchRes
//	@Router		/store-account/list-transactions [post]
func (h *StoreAccountHandler) ListTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreAccountHandler.ListTransactions")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreAccountHandler.ListTransactions")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreAccountTransactionListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		page := upagination.New(req.Page, req.Size)
		params := domain.StoreAccountTransactionSearchParams{
			StartAt: req.StartAt.ToPtrStartOfDay(),
			EndAt:   req.EndAt.ToPtrEndOfDay(),
			StoreID: user.StoreID,
			Type:    req.Type,
		}

		res, err := h.StoreAccountInteractor.PagedListTransactions(ctx, page, params)
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
