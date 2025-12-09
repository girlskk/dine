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
	r.POST("/create", h.Create())
	r.POST("/update", h.Update())
	r.POST("/list", h.List())
	r.POST("/detail", h.Detail())
}

// Create 创建门店接口处理
//
//	@Tags		门店管理
//	@Summary	创建门店
//	@Security	BearerAuth
//	@Param		data	body	types.StoreCreateReq	true	"请求参数"
//	@Success	200
//	@Router		/store/create [post]
func (h *StoreHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreHandler.Create")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		store := types.ToStoreCreateReq(req)
		user := &domain.BackendUser{
			Username: req.UserName,
			Nickname: req.UserName,
		}
		err := user.SetPassword(req.Password)
		if err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		err = h.StoreInteractor.Create(ctx, store, user)
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
		store := types.ToStoreUpdateReq(req)
		user := &domain.BackendUser{StoreID: req.ID}
		if req.Password != "" {
			err := user.SetPassword(req.Password)
			if err != nil {
				c.Error(uerr.BadRequest(err.Error()))
				return
			}
		}

		err := h.StoreInteractor.Update(ctx, store, user)
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

// List 获取门店列表接口处理
//
//	@Tags		门店管理
//	@Summary	获取门店列表
//	@Security	BearerAuth
//	@Param		data	body		types.StoreListReq		true	"请求参数"
//	@Success	200		{object}	domain.StoreSearchRes	"成功"
//	@Router		/store/list [post]
func (h *StoreHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)
		var req types.StoreListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		page := upagination.New(req.Page, req.Size)
		params := domain.StoreSearchParams{
			Name: req.Name,
			City: req.City,
		}

		res, err := h.StoreInteractor.PagedListBySearch(ctx, page, params)
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

// Detail 获取门店详情接口处理
//
//	@Tags		门店管理
//	@Summary	获取门店详情
//	@Security	BearerAuth
//	@Param		data	body		types.StoreIDReq	true	"请求参数"
//	@Success	200		{object}	domain.Store		"成功"
//	@Router		/store/detail [post]
func (h *StoreHandler) Detail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "StoreHandler.Detail")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("StoreHandler.Detail")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)
		var req types.StoreIDReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		store, err := h.StoreInteractor.GetDetail(ctx, req.ID)
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
