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

// TableHandler 处理台桌请求
type TableHandler struct {
	TableInteractor domain.TableInteractor
}

func NewTableHandler(interactor domain.TableInteractor) *TableHandler {
	return &TableHandler{
		TableInteractor: interactor,
	}
}

// Routes 注册路由
func (h *TableHandler) Routes(r gin.IRouter) {
	r = r.Group("/table")
	r.POST("/create", h.Create())
	r.POST("/update", h.Update())
	r.POST("/delete", h.Delete())
	r.POST("/list", h.List())
}

// Create 创建台桌
//
//	@Tags		台桌管理
//	@Summary	创建台桌
//	@Security	BearerAuth
//	@Param		data	body	types.TableCreateReq	true	"请求参数"
//	@Success	200
//	@Router		/table/create [post]
func (h *TableHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "TableHandler.Create")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("TableHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.TableCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		table := &domain.Table{
			Name:      req.Name,
			SeatCount: req.SeatCount,
			AreaID:    req.AreaID,
		}

		if err := h.TableInteractor.Create(ctx, table); err != nil {
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

// Update 更新台桌
//
//	@Tags		台桌管理
//	@Summary	更新台桌
//	@Security	BearerAuth
//	@Param		data	body	types.TableUpdateReq	true	"请求参数"
//	@Success	200
//	@Router		/table/update [post]
func (h *TableHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "TableHandler.Update")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("TableHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.TableUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		table := &domain.Table{
			ID:        req.ID,
			Name:      req.Name,
			SeatCount: req.SeatCount,
			AreaID:    req.AreaID,
		}

		if err := h.TableInteractor.Update(ctx, table); err != nil {
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

// Delete 删除台桌
//
//	@Tags		台桌管理
//	@Summary	删除台桌
//	@Security	BearerAuth
//	@Param		data	body	types.TableDeleteReq	true	"请求参数"
//	@Success	200
//	@Router		/table/delete [post]
func (h *TableHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "TableHandler.Delete")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("TableHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.TableDeleteReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.TableInteractor.Delete(ctx, req.ID); err != nil {
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

// List 获取台桌列表
//
//	@Tags		台桌管理
//	@Summary	台桌列表
//	@Param		data	body		types.TableListReq	true	"请求参数"
//	@Success	200		{object}	domain.TableSearchRes
//	@Router		/table/list [post]
func (h *TableHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "TableHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("TableHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.TableListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(req.Page, req.Size)
		user := domain.FromBackendUserContext(ctx)
		params := domain.TableSearchParams{
			AreaID:  req.AreaID,
			StoreID: user.Store.ID,
		}

		res, err := h.TableInteractor.PagedListBySearch(ctx, page, params)
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
