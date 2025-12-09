package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/customer/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type ProductHandler struct {
	ProductInteractor  domain.ProductInteractor
	CategoryInteractor domain.CategoryInteractor
}

func NewProductHandler(productInteractor domain.ProductInteractor,
	categoryInteractor domain.CategoryInteractor,
) *ProductHandler {
	return &ProductHandler{
		ProductInteractor:  productInteractor,
		CategoryInteractor: categoryInteractor,
	}
}

func (h *ProductHandler) Routes(r gin.IRouter) {
	r = r.Group("/product")
	r.POST("/category/list", h.ListCategory())
	r.POST("/list", h.ListProduct())
	r.POST("/detail", h.GetProductDetail())
}

func (h *ProductHandler) NoAuths() []string {
	return []string{
		"/product/category/list",
		"/product/list",
		"/product/detail",
	}
}

// ListCategory 获取分类列表
//
//	@Tags		商品管理
//	@Summary	商品分类列表
//	@Param		data	body		types.CategoryListReq		true	"请求参数"
//	@Success	200		{object}	domain.CategorySearchRes	"成功"
//	@Router		/product/category/list [post]
func (h *ProductHandler) ListCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.ListCategory")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.ListCategory")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CategoryListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(1, upagination.MaxSize)

		params := domain.CategorySearchParams{
			StoreID: req.StoreID,
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

// ListProduct 获取商品列表
//
//	@Tags		商品管理
//	@Summary	商品列表
//	@Param		data	body		types.ProductListReq	true	"请求参数"
//	@Success	200		{object}	domain.ProductSearchRes	"成功"
//	@Router		/product/list [post]
func (h *ProductHandler) ListProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(1, upagination.MaxSize)
		params := domain.ProductSearchParams{
			CategoryID: req.CategoryID,
			StoreID:    req.StoreID,
			Status:     domain.ProductStatusApproved,
		}

		res, err := h.ProductInteractor.PagedListBySearch(ctx, page, params)
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

// GetProductDetail 获取商品详情
//
//	@Tags		商品管理
//	@Summary	商品详情
//	@Param		data	body		types.ProductIDReq	true	"请求参数"
//	@Success	200		{object}	domain.Product		"成功"
//	@Router		/product/detail [post]
func (h *ProductHandler) GetProductDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.GetProductDetail")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.GetProductDetail")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductIDReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		product, err := h.ProductInteractor.GetDetail(c.Request.Context(), req.ID)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, product)
	}
}
