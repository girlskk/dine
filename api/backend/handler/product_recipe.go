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

// ProductRecipeHandler 处理商品做法请求
type ProductRecipeHandler struct {
	RecipeInteractor domain.ProductRecipeInteractor
}

func NewProductRecipeHandler(interactor domain.ProductRecipeInteractor) *ProductRecipeHandler {
	return &ProductRecipeHandler{
		RecipeInteractor: interactor,
	}
}

// Routes 注册路由
func (h *ProductRecipeHandler) Routes(r gin.IRouter) {
	r = r.Group("/product/recipe")
	r.POST("/create", h.Create())
	r.POST("/update", h.Update())
	r.POST("/delete", h.Delete())
	r.POST("/list", h.List())
}

// Create 创建商品做法
//
//	@Tags		商品管理-做法
//	@Summary	创建商品做法
//	@Security	BearerAuth
//	@Param		data	body	types.RecipeCreateReq	true	"请求参数"
//	@Success	200
//	@Router		/product/recipe/create [post]
func (h *ProductRecipeHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeHandler.Create")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductRecipeHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RecipeCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		recipe := &domain.ProductRecipe{
			Name: req.Name,
		}

		if err := h.RecipeInteractor.Create(ctx, recipe); err != nil {
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

// Update 更新商品做法
//
//	@Tags		商品管理-做法
//	@Summary	更新商品做法
//	@Security	BearerAuth
//	@Param		data	body	types.RecipeUpdateReq	true	"请求参数"
//	@Success	200
//	@Router		/product/recipe/update [post]
func (h *ProductRecipeHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeHandler.Update")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductRecipeHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RecipeUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		recipe := &domain.ProductRecipe{
			ID:   req.ID,
			Name: req.Name,
		}

		if err := h.RecipeInteractor.Update(ctx, recipe); err != nil {
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

// Delete 删除商品做法
//
//	@Tags		商品管理-做法
//	@Summary	删除商品做法
//	@Security	BearerAuth
//	@Param		data	body	types.RecipeDeleteReq	true	"请求参数"
//	@Success	200
//	@Router		/product/recipe/delete [post]
func (h *ProductRecipeHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeHandler.Delete")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductRecipeHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RecipeDeleteReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.RecipeInteractor.Delete(ctx, req.ID); err != nil {
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

// List 获取做法列表
//
//	@Tags		商品管理-做法
//	@Summary	商品做法列表
//	@Param		data	body		types.RecipeListReq	true	"请求参数"
//	@Success	200		{object}	domain.RecipeSearchRes
//	@Router		/product/recipe/list [post]
func (h *ProductRecipeHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductRecipeHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductRecipeHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RecipeListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(req.Page, req.Size)
		params := domain.RecipeSearchParams{}

		res, err := h.RecipeInteractor.PagedListBySearch(ctx, page, params)
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
