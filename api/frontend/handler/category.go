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
	r.POST("/list", h.List())
}

// List 获取分类列表
//
//	@Tags		商品管理
//	@Security	BearerAuth
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
		user := domain.FromFrontendUserContext(ctx)
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
