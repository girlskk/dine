package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// ProductHandler 处理商品请求
type ProductHandler struct {
	ProductInteractor domain.ProductInteractor
}

func NewProductHandler(interactor domain.ProductInteractor) *ProductHandler {
	return &ProductHandler{
		ProductInteractor: interactor,
	}
}

func (h *ProductHandler) Routes(r gin.IRouter) {
	r = r.Group("/product")
	r.POST("/list", h.List())
	r.POST("/detail", h.GetDetail())
	r.POST("/setmeal/list", h.ListSetmealDetails())
	r.POST("/clear-stock", h.ClearStock())
	r.POST("/restore-stock", h.RestoreStock())
}

// List 获取商品列表
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	商品列表
//	@Param		data	body		types.ProductListReq	true	"请求参数"
//	@Success	200		{object}	domain.ProductSearchRes	"成功"
//	@Router		/product/list [post]
func (h *ProductHandler) List() gin.HandlerFunc {
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

		page := upagination.New(req.Page, req.Size)
		user := domain.FromFrontendUserContext(ctx)
		params := domain.ProductSearchParams{
			CategoryID: req.CategoryID,
			SaleStatus: req.SaleStatus,
			StoreID:    user.Store.ID,
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

// GetDetail 获取商品详情
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	商品详情
//	@Param		data	body		types.ProductIDReq	true	"请求参数"
//	@Success	200		{object}	domain.Product		"成功"
//	@Router		/product/detail [post]
func (h *ProductHandler) GetDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.GetDetail")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.GetDetail")
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

		user := domain.FromFrontendUserContext(ctx)
		if user.Store.ID != product.StoreID {
			c.Error(uerr.BadRequest(domain.ErrProductNotExists.Error()))
			return
		}

		response.Ok(c, product)
	}
}

// ListSetmealDetails 套餐商品详情列表
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	套餐商品详情列表
//	@Param		data	body		types.ProductIDReq		true	"请求参数"
//	@Success	200		{object}	domain.SetMealDetails	"成功"
//	@Router		/product/setmeal/list [post]
func (h *ProductHandler) ListSetmealDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.ListSetmealDetails")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.ListSetmealDetails")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductIDReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		details, err := h.ProductInteractor.ListSetmealDetails(ctx, req.ID)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, details)
	}
}

// ClearStock 商品估清
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	商品估清
//	@Param		data	body	types.ProductClearStockReq	true	"请求参数"
//	@Success	200
//	@Router		/product/clear-stock [post]
func (h *ProductHandler) ClearStock() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.ClearStock")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.ClearStock")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductClearStockReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.ProductInteractor.ClearStock(ctx, req.ProductID, req.SpecIDs); err != nil {
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

// RestoreStock 取消估清
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	取消估清
//	@Param		data	body	types.ProductRestoreStockReq	true	"请求参数"
//	@Success	200
//	@Router		/product/restore-stock [post]
func (h *ProductHandler) RestoreStock() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.RestoreStock")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.RestoreStock")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductRestoreStockReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.ProductInteractor.RestoreStock(ctx, req.ProductID, req.SpecIDs); err != nil {
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
