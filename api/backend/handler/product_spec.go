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

// ProductSpecHandler 处理商品规格请求
type ProductSpecHandler struct {
	SpecInteractor domain.ProductSpecInteractor
}

func NewProductSpecHandler(interactor domain.ProductSpecInteractor) *ProductSpecHandler {
	return &ProductSpecHandler{
		SpecInteractor: interactor,
	}
}

// Routes 注册路由
func (h *ProductSpecHandler) Routes(r gin.IRouter) {
	r = r.Group("/product/spec")
	r.POST("/create", h.Create())
	r.POST("/update", h.Update())
	r.POST("/delete", h.Delete())
	r.POST("/list", h.List())
}

// Create 创建商品规格
//
//	@Tags		商品管理-规格
//	@Summary	创建商品规格
//	@Security	BearerAuth
//	@Param		data	body	types.SpecCreateReq	true	"请求参数"
//	@Success	200
//	@Router		/product/spec/create [post]
func (h *ProductSpecHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecHandler.Create")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductSpecHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.SpecCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		spec := &domain.ProductSpec{
			Name: req.Name,
		}

		if err := h.SpecInteractor.Create(ctx, spec); err != nil {
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

// Update 更新商品规格
//
//	@Tags		商品管理-规格
//	@Summary	更新商品规格
//	@Security	BearerAuth
//	@Param		data	body	types.SpecUpdateReq	true	"请求参数"
//	@Success	200
//	@Router		/product/spec/update [post]
func (h *ProductSpecHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecHandler.Update")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductSpecHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.SpecUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		spec := &domain.ProductSpec{
			ID:   req.ID,
			Name: req.Name,
		}

		if err := h.SpecInteractor.Update(ctx, spec); err != nil {
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

// Delete 删除商品规格
//
//	@Tags		商品管理-规格
//	@Summary	删除商品规格
//	@Security	BearerAuth
//	@Param		data	body	types.SpecDeleteReq	true	"请求参数"
//	@Success	200
//	@Router		/product/spec/delete [post]
func (h *ProductSpecHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecHandler.Delete")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductSpecHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.SpecDeleteReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.SpecInteractor.Delete(ctx, req.ID); err != nil {
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

// List 获取规格列表
//
//	@Tags		商品管理-规格
//	@Summary	商品规格列表
//	@Param		data	body		types.SpecListReq	true	"请求参数"
//	@Success	200		{object}	domain.SpecSearchRes
//	@Router		/product/spec/list [post]
func (h *ProductSpecHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ProductSpecHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ProductSpecHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.SpecListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(req.Page, req.Size)
		params := domain.SpecSearchParams{}

		res, err := h.SpecInteractor.PagedListBySearch(ctx, page, params)
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
