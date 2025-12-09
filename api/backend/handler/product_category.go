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

// CategoryHandler 处理商品分类请求
type CategoryHandler struct {
	CategoryInteractor domain.CategoryInteractor
}

func NewCategoryHandler(interactor domain.CategoryInteractor) *CategoryHandler {
	return &CategoryHandler{
		CategoryInteractor: interactor,
	}
}

func (h *CategoryHandler) Routes(r gin.IRouter) {
	r = r.Group("/product/category")
	r.POST("/create", h.Create())
	r.POST("/update", h.Update())
	r.POST("/delete", h.Delete())
	r.POST("/list", h.List())
}

// Create 创建商品分类
//
//	@Tags		商品管理-分类
//	@Summary	创建商品分类
//	@Security	BearerAuth
//	@Param		data	body	types.CategoryCreateReq	true	"请求参数"
//	@Success	200
//	@Router		/product/category/create [post]
func (h *CategoryHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryHandler.Create")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("CategoryHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CategoryCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		category := &domain.Category{
			Name: req.Name,
		}

		if err := h.CategoryInteractor.Create(ctx, category); err != nil {
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

// Update 更新商品分类
//
//	@Tags		商品管理-分类
//	@Summary	更新商品分类
//	@Security	BearerAuth
//	@Param		data	body	types.CategoryUpdateReq	true	"请求参数"
//	@Success	200
//	@Router		/product/category/update [post]
func (h *CategoryHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryHandler.Update")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("CategoryHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CategoryUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		category := &domain.Category{
			ID:   req.ID,
			Name: req.Name,
		}

		if err := h.CategoryInteractor.Update(ctx, category); err != nil {
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

// Delete 删除商品分类
//
//	@Tags		商品管理-分类
//	@Summary	删除商品分类
//	@Security	BearerAuth
//	@Param		data	body	types.CategoryDeleteReq	true	"请求参数"
//	@Success	200
//	@Router		/product/category/delete [post]
func (h *CategoryHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryHandler.Delete")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("CategoryHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CategoryDeleteReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.CategoryInteractor.Delete(ctx, req.ID); err != nil {
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

// List 获取分类列表
//
//	@Tags		商品管理-分类
//	@Summary	商品分类列表
//	@Param		data	body		types.CategoryListReq		true	"请求参数"
//	@Success	200		{object}	domain.CategorySearchRes	"成功"
//	@Router		/product/category/list [post]
func (h *CategoryHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "CategoryHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("CategoryHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CategoryListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(req.Page, req.Size)
		user := domain.FromBackendUserContext(ctx)
		params := domain.CategorySearchParams{
			StoreID: user.Store.ID,
		}
		res, err := h.CategoryInteractor.PagedListBySearch(ctx, page, params)
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
