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

// ProductAttrHandler 处理商品属性请求
type ProductAttrHandler struct {
	AttrInteractor domain.ProductAttrInteractor
}

func NewProductAttrHandler(interactor domain.ProductAttrInteractor) *ProductAttrHandler {
	return &ProductAttrHandler{
		AttrInteractor: interactor,
	}
}

// Routes 注册路由
func (h *ProductAttrHandler) Routes(r gin.IRouter) {
	r = r.Group("/product/attr")
	r.POST("/create", h.Create())
	r.POST("/update", h.Update())
	r.POST("/delete", h.Delete())
	r.POST("/list", h.List())
}

// Create 创建商品属性
//
//	@Tags		商品管理-属性
//	@Summary	创建商品属性
//	@Security	BearerAuth
//	@Param		data	body	types.AttrCreateReq	true	"请求参数"
//	@Success	200
//	@Router		/product/attr/create [post]
func (h *ProductAttrHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrHandler.Create")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductAttrHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AttrCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		attr := &domain.ProductAttr{
			Name: req.Name,
		}
		if err := h.AttrInteractor.Create(ctx, attr); err != nil {
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

// Update 更新商品属性
//
//	@Tags		商品管理-属性
//	@Summary	更新商品属性
//	@Security	BearerAuth
//	@Param		data	body	types.AttrUpdateReq	true	"请求参数"
//	@Success	200
//	@Router		/product/attr/update [post]
func (h *ProductAttrHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrHandler.Update")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductAttrHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AttrUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		attr := &domain.ProductAttr{
			ID:   req.ID,
			Name: req.Name,
		}

		if err := h.AttrInteractor.Update(ctx, attr); err != nil {
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

// Delete 删除商品属性
//
//	@Tags		商品管理-属性
//	@Summary	删除商品属性
//	@Security	BearerAuth
//	@Param		data	body	types.AttrDeleteReq	true	"请求参数"
//	@Success	200
//	@Router		/product/attr/delete [post]
func (h *ProductAttrHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrHandler.Delete")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductAttrHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AttrDeleteReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.AttrInteractor.Delete(ctx, req.ID); err != nil {
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

// List 获取属性列表
//
//	@Tags		商品管理-属性
//	@Summary	商品属性列表
//	@Param		data	body		types.AttrListReq	true	"请求参数"
//	@Success	200		{object}	domain.AttrSearchRes
//	@Router		/product/attr/list [post]
func (h *ProductAttrHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductAttrHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductAttrHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AttrListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(req.Page, req.Size)
		params := domain.AttrSearchParams{}

		res, err := h.AttrInteractor.PagedListBySearch(ctx, page, params)
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
