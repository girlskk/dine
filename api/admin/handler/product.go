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
	r.POST("/list", h.List())
	r.POST("/detail", h.Detail())
	r.POST("/approve", h.Approve())
	r.POST("/unapprove", h.UnApprove())

}

// List 获取商品列表
//
//	@Tags		商品管理
//	@Summary	商品列表
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

		page := upagination.New(req.Page, req.Size)
		params := domain.ProductSearchParams{
			Name:    req.Name,
			StoreID: req.StoreID,
			Status:  req.Status,
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

		response.Ok(c, product)
	}
}

// Approve 审批接口处理
//
//	@Tags		商品管理
//	@Summary	商品审批
//	@Security	BearerAuth
//	@Param		data	body		types.ProductApproveReq	true	"请求参数"
//	@Success	200		{object}	domain.Product
//	@Router		/product/approve [post]
func (h *ProductHandler) Approve() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.Approve")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.Approve")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductApproveReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		err := h.ProductInteractor.Approve(ctx, req.IDs, req.AllowPointPay)
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

// UnApprove 审批接口处理
//
//	@Tags		商品管理
//	@Summary	商品反审批
//	@Security	BearerAuth
//	@Param		data	body		types.ProductUnApproveReq	true	"请求参数"
//	@Success	200		{object}	domain.Product
//	@Router		/product/unapprove [post]
func (h *ProductHandler) UnApprove() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductHandler.UnApprove")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductHandler.UnApprove")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductUnApproveReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		err := h.ProductInteractor.UnApprove(ctx, req.IDs)
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
