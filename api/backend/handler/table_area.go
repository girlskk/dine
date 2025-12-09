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

// TableAreaHandler 处理台桌区域请求
type TableAreaHandler struct {
	AreaInteractor domain.TableAreaInteractor
}

func NewTableAreaHandler(interactor domain.TableAreaInteractor) *TableAreaHandler {
	return &TableAreaHandler{
		AreaInteractor: interactor,
	}
}

// Routes 注册路由
func (h *TableAreaHandler) Routes(r gin.IRouter) {
	r = r.Group("/table-area")
	r.POST("/create", h.Create())
	r.POST("/update", h.Update())
	r.POST("/delete", h.Delete())
	r.POST("/list", h.List())
}

// Create 创建区域
//
//	@Tags		台桌管理-区域
//	@Summary	创建台桌区域
//	@Security	BearerAuth
//	@Param		data	body	types.AreaCreateReq	true	"请求参数"
//	@Success	200
//	@Router		/table-area/create [post]
func (h *TableAreaHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaHandler.Create")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("TableAreaHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AreaCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		area := &domain.TableArea{
			Name: req.Name,
		}

		if err := h.AreaInteractor.Create(ctx, area); err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, area)
	}
}

// Update 更新区域
//
//	@Tags		台桌管理-区域
//	@Summary	更新台桌区域
//	@Security	BearerAuth
//	@Param		data	body	types.AreaUpdateReq	true	"请求参数"
//	@Success	200
//	@Router		/table-area/update [post]
func (h *TableAreaHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaHandler.Update")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("TableAreaHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AreaUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		area := &domain.TableArea{
			ID:   req.ID,
			Name: req.Name,
		}

		if err := h.AreaInteractor.Update(ctx, area); err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, area)
	}
}

// Delete 删除区域
//
//	@Tags		台桌管理-区域
//	@Summary	删除台桌区域
//	@Security	BearerAuth
//	@Param		data	body	types.AreaDeleteReq	true	"请求参数"
//	@Success	200
//	@Router		/table-area/delete [post]
func (h *TableAreaHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaHandler.Delete")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("TableAreaHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AreaDeleteReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.AreaInteractor.Delete(ctx, req.ID); err != nil {
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

// List 获取区域列表
//
//	@Tags		台桌管理-区域
//	@Summary	台桌区域列表
//	@Param		data	body		types.AreaListReq	true	"请求参数"
//	@Success	200		{object}	domain.AreaSearchRes
//	@Router		/table-area/list [post]
func (h *TableAreaHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("TableAreaHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AreaListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(req.Page, req.Size)
		user := domain.FromBackendUserContext(ctx)
		params := domain.AreaSearchParams{
			StoreID: user.Store.ID,
		}

		res, err := h.AreaInteractor.PagedListBySearch(ctx, page, params)
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
