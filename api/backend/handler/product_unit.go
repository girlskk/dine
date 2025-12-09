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

// ProductUnitHandler 处理商品单位请求
type ProductUnitHandler struct {
	UnitInteractor domain.ProductUnitInteractor
}

func NewProductUnitHandler(interactor domain.ProductUnitInteractor) *ProductUnitHandler {
	return &ProductUnitHandler{
		UnitInteractor: interactor,
	}
}

// Routes 注册路由
func (h *ProductUnitHandler) Routes(r gin.IRouter) {
	r = r.Group("/product/unit")
	r.POST("/create", h.Create())
	r.POST("/update", h.Update())
	r.POST("/delete", h.Delete())
	r.POST("/list", h.List())
}

// Create 创建商品单位
//
//	@Tags		商品管理-单位
//	@Summary	创建商品单位
//	@Security	BearerAuth
//	@Param		data	body	types.UnitCreateReq	true	"请求参数"
//	@Success	200
//	@Router		/product/unit/create [post]
func (h *ProductUnitHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitHandler.Create")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductUnitHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.UnitCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		unit := &domain.ProductUnit{
			Name: req.Name,
		}

		if err := h.UnitInteractor.Create(ctx, unit); err != nil {
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

// Update 更新商品单位
//
//	@Tags		商品管理-单位
//	@Summary	更新商品单位
//	@Security	BearerAuth
//	@Param		data	body	types.UnitUpdateReq	true	"请求参数"
//	@Success	200
//	@Router		/product/unit/update [post]
func (h *ProductUnitHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitHandler.Update")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductUnitHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.UnitUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		unit := &domain.ProductUnit{
			ID:   req.ID,
			Name: req.Name,
		}

		if err := h.UnitInteractor.Update(ctx, unit); err != nil {
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

// Delete 删除商品单位
//
//	@Tags		商品管理-单位
//	@Summary	删除商品单位
//	@Security	BearerAuth
//	@Param		data	body	types.UnitDeleteReq	true	"请求参数"
//	@Success	200
//	@Router		/product/unit/delete [post]
func (h *ProductUnitHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitHandler.Delete")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductUnitHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.UnitDeleteReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.UnitInteractor.Delete(ctx, req.ID); err != nil {
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

// List 获取单位列表
//
//	@Tags		商品管理-单位
//	@Summary	商品单位列表
//	@Param		data	body		types.UnitListReq	true	"请求参数"
//	@Success	200		{object}	domain.UnitSearchRes
//	@Router		/product/unit/list [post]
func (h *ProductUnitHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUnitHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductUnitHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.UnitListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(req.Page, req.Size)
		params := domain.UnitSearchParams{}

		res, err := h.UnitInteractor.PagedListBySearch(ctx, page, params)
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
