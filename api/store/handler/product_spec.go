package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/store/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type ProductSpecHandler struct {
	ProductSpecInteractor domain.ProductSpecInteractor
}

func NewProductSpecHandler(productSpecInteractor domain.ProductSpecInteractor) *ProductSpecHandler {
	return &ProductSpecHandler{
		ProductSpecInteractor: productSpecInteractor,
	}
}

func (h *ProductSpecHandler) Routes(r gin.IRouter) {
	r = r.Group("product/spec")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("", h.List())
}

func (h *ProductSpecHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		商品规格
//	@Security	BearerAuth
//	@Summary	创建商品规格
//	@Param		data	body	types.ProductSpecCreateReq	true	"请求信息"
//	@Success	200
//	@Router		/product/spec [post]
func (h *ProductSpecHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductSpecHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductSpecCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)

		spec := &domain.ProductSpec{
			ID:         uuid.New(),
			Name:       req.Name,
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
		}

		err := h.ProductSpecInteractor.Create(ctx, spec)

		if err != nil {
			if errors.Is(err, domain.ErrProductSpecNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.ProductSpecNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create product spec: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update
//
//	@Tags		商品规格
//	@Security	BearerAuth
//	@Summary	更新商品规格
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string						true	"规格ID"
//	@Param		data	body	types.ProductSpecUpdateReq	true	"请求信息"
//	@Success	200
//	@Router		/product/spec/{id} [put]
func (h *ProductSpecHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductSpecHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取规格ID
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.ProductSpecUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		// 将请求数据映射到 domain.ProductSpec
		spec := &domain.ProductSpec{
			ID:   id,
			Name: req.Name,
		}

		user := domain.FromStoreUserContext(ctx)
		err = h.ProductSpecInteractor.Update(ctx, spec, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductSpecNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.ProductSpecNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update product spec: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete
//
//	@Tags		商品规格
//	@Security	BearerAuth
//	@Summary	删除商品规格
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"规格ID"
//	@Success	200	"No Content"
//	@Router		/product/spec/{id} [delete]
func (h *ProductSpecHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductSpecHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		err = h.ProductSpecInteractor.Delete(ctx, id, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductSpecDeleteHasProducts) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.ProductSpecDeleteHasProducts, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}

			err = fmt.Errorf("failed to delete product spec: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// List
//
//	@Tags		商品规格
//	@Security	BearerAuth
//	@Summary	获取商品规格列表
//	@Param		data	query		types.ProductSpecListReq	true	"请求信息"
//	@Success	200		{object}	domain.ProductSpecSearchRes	"成功"
//	@Router		/product/spec [get]
func (h *ProductSpecHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductSpecHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductSpecListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		page := upagination.New(req.Page, req.Size)
		user := domain.FromStoreUserContext(ctx)

		params := domain.ProductSpecSearchParams{
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
			Name:       req.Name,
		}

		res, err := h.ProductSpecInteractor.PagedListBySearch(ctx, page, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to list product specs: %w", err)
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}
