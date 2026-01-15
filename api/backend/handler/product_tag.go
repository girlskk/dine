package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type ProductTagHandler struct {
	ProductTagInteractor domain.ProductTagInteractor
}

func NewProductTagHandler(productTagInteractor domain.ProductTagInteractor) *ProductTagHandler {
	return &ProductTagHandler{
		ProductTagInteractor: productTagInteractor,
	}
}

func (h *ProductTagHandler) Routes(r gin.IRouter) {
	r = r.Group("product/tag")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("", h.List())
}

func (h *ProductTagHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		商品标签
//	@Security	BearerAuth
//	@Summary	创建商品标签
//	@Param		data	body	types.ProductTagCreateReq	true	"请求信息"
//	@Success	200
//	@Router		/product/tag [post]
func (h *ProductTagHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductTagHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductTagCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		tag := &domain.ProductTag{
			ID:         uuid.New(),
			Name:       req.Name,
			MerchantID: user.MerchantID,
		}

		err := h.ProductTagInteractor.Create(ctx, tag)

		if err != nil {
			if errors.Is(err, domain.ErrProductTagNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.ProductTagNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create product tag: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update
//
//	@Tags		商品标签
//	@Security	BearerAuth
//	@Summary	更新商品标签
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string						true	"标签ID"
//	@Param		data	body	types.ProductTagUpdateReq	true	"请求信息"
//	@Success	200
//	@Router		/product/tag/{id} [put]
func (h *ProductTagHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductTagHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取标签ID
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.ProductTagUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		// 将请求数据映射到 domain.ProductTag
		tag := &domain.ProductTag{
			ID:   id,
			Name: req.Name,
		}

		user := domain.FromBackendUserContext(ctx)

		err = h.ProductTagInteractor.Update(ctx, tag, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductTagNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.ProductTagNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update product tag: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete
//
//	@Tags		商品标签
//	@Security	BearerAuth
//	@Summary	删除商品标签
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"标签ID"
//	@Success	200	"No Content"
//	@Router		/product/tag/{id} [delete]
func (h *ProductTagHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductTagHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		err = h.ProductTagInteractor.Delete(ctx, id, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductTagDeleteHasProducts) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.ProductTagDeleteHasProducts, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}

			err = fmt.Errorf("failed to delete product tag: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// List
//
//	@Tags		商品标签
//	@Security	BearerAuth
//	@Summary	获取商品标签列表
//	@Param		data	query		types.ProductTagListReq		true	"请求信息"
//	@Success	200		{object}	domain.ProductTagSearchRes	"成功"
//	@Router		/product/tag [get]
func (h *ProductTagHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductTagHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductTagListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		page := upagination.New(req.Page, req.Size)
		user := domain.FromBackendUserContext(ctx)

		params := domain.ProductTagSearchParams{
			MerchantID:   user.MerchantID,
			Name:         req.Name,
			OnlyMerchant: true,
		}

		res, err := h.ProductTagInteractor.PagedListBySearch(ctx, page, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to list product tags: %w", err)
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}
