package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

// StoreHandler 处理门店请求
type StoreHandler struct {
	StoreInteractor domain.StoreInteractor
}

func NewStoreHandler(interactor domain.StoreInteractor) *StoreHandler {
	return &StoreHandler{
		StoreInteractor: interactor,
	}
}

// Routes 注册路由
func (h *StoreHandler) Routes(r gin.IRouter) {
	r = r.Group("/store")
	r.POST("/update", h.Update())
	r.POST("/detail", h.Detail())
}

// Update 编辑门店接口处理
//
//	@Tags		门店管理
//	@Summary	更新门店
//	@Security	BearerAuth
//	@Param		data	body	types.StoreUpdateReq	true	"请求参数"
//	@Success	200
//	@Router		/store/update [post]
func (h *StoreHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreHandler.Update")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		store := types.ToStoreUpsetReq(req)
		user := domain.FromBackendUserContext(ctx)
		if req.Password != "" {
			err := user.SetPassword(req.Password)
			if err != nil {
				c.Error(uerr.BadRequest(err.Error()))
				return
			}
		}
		store.ID = user.StoreID

		err := h.StoreInteractor.UpdateByStore(ctx, store, user)
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

// Detail 获取门店详情接口处理
//
//	@Tags		门店管理
//	@Summary	获取门店详情
//	@Security	BearerAuth
//	@Success	200	{object}	domain.Store	"成功"
//	@Router		/store/detail [post]
func (h *StoreHandler) Detail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreHandler.Detail")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreHandler.Detail")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromBackendUserContext(ctx)
		store, err := h.StoreInteractor.GetDetail(ctx, user.StoreID)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}
		response.Ok(c, store)
	}
}
