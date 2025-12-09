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

// ProductHandler 处理商品请求
type ProductHandler struct {
	ProductInteractor domain.ProductInteractor
}

func NewProductHandler(interactor domain.ProductInteractor) *ProductHandler {
	return &ProductHandler{
		ProductInteractor: interactor,
	}
}

// Routes 注册路由
func (h *ProductHandler) Routes(r gin.IRouter) {
	r = r.Group("/product")
	r.POST("/create", h.Create())
	r.POST("/update", h.Update())
	r.POST("/delete", h.Delete())
	r.POST("/list", h.List())
	r.POST("/detail", h.Detail())
}

// Create 创建商品
//
//	@Tags		商品管理
//	@Summary	创建商品
//	@Security	BearerAuth
//	@Param		data	body	types.ProductCreateReq	true	"请求参数"
//	@Success	200
//	@Router		/product/create [post]
func (h *ProductHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.Create")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		params := types.ToProductUpsetReq(req, 0)

		if err := h.ProductInteractor.Create(ctx, params); err != nil {
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

// Update 更新商品
//
//	@Tags		商品管理
//	@Summary	更新商品
//	@Security	BearerAuth
//	@Param		data	body	types.ProductUpdateReq	true	"请求参数"
//	@Success	200
//	@Router		/product/update [post]
func (h *ProductHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.Update")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		params := types.ToProductUpdateReq(req)

		if err := h.ProductInteractor.Update(ctx, params); err != nil {
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

// Delete 删除商品
//
//	@Tags		商品管理
//	@Summary	删除商品
//	@Security	BearerAuth
//	@Param		data	body	types.ProductDeleteReq	true	"请求参数"
//	@Success	200
//	@Router		/product/delete [post]
func (h *ProductHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.Delete")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductDeleteReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.ProductInteractor.Delete(ctx, req.ID); err != nil {
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

// List 获取商品列表
//
//	@Tags		商品管理
//	@Summary	商品列表
//	@Security	BearerAuth
//	@Param		data	body		types.ProductListReq	true	"请求参数"
//	@Success	200		{object}	domain.ProductSearchRes
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

		user := domain.FromBackendUserContext(ctx)
		page := upagination.New(req.Page, req.Size)
		params := domain.ProductSearchParams{
			CategoryID: req.CategoryID,
			Name:       req.Name,
			Status:     req.Status,
			SaleStatus: req.SaleStatus,
			StoreID:    user.Store.ID,
			Type:       req.Type,
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

// Detail 详情接口处理
//
//	@Tags		商品管理
//	@Summary	商品详情
//	@Security	BearerAuth
//	@Param		data	body		types.ProductDetailReq	true	"请求参数"
//	@Success	200		{object}	domain.Product
//	@Router		/product/detail [post]
func (h *ProductHandler) Detail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.Detail")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.Detail")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductDetailReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		product, err := h.ProductInteractor.GetDetail(ctx, req.ID)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		user := domain.FromBackendUserContext(ctx)
		if user.Store.ID != product.StoreID {
			c.Error(uerr.BadRequest(domain.ErrProductNotExists.Error()))
			return
		}

		response.Ok(c, product)
	}
}
