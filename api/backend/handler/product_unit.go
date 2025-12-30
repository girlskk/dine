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

type ProductUnitHandler struct {
	ProductUnitInteractor domain.ProductUnitInteractor
}

func NewProductUnitHandler(productUnitInteractor domain.ProductUnitInteractor) *ProductUnitHandler {
	return &ProductUnitHandler{
		ProductUnitInteractor: productUnitInteractor,
	}
}

func (h *ProductUnitHandler) Routes(r gin.IRouter) {
	r = r.Group("product/unit")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("", h.List())
}

func (h *ProductUnitHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		商品单位
//	@Security	BearerAuth
//	@Summary	创建商品单位
//	@Param		data	body	types.ProductUnitCreateReq	true	"请求信息"
//	@Success	200
//	@Router		/product/unit [post]
func (h *ProductUnitHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductUnitHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductUnitCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		unit := &domain.ProductUnit{
			ID:         uuid.New(),
			Name:       req.Name,
			Type:       domain.ProductUnitType(req.Type),
			MerchantID: user.MerchantID,
		}

		err := h.ProductUnitInteractor.Create(ctx, unit)

		if err != nil {
			if errors.Is(err, domain.ErrProductUnitNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.ProductUnitNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create product unit: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update
//
//	@Tags		商品单位
//	@Security	BearerAuth
//	@Summary	更新商品单位
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string						true	"单位ID"
//	@Param		data	body	types.ProductUnitUpdateReq	true	"请求信息"
//	@Success	200
//	@Router		/product/unit/{id} [put]
func (h *ProductUnitHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductUnitHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取单位ID
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.ProductUnitUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		// 将请求数据映射到 domain.ProductUnit
		unit := &domain.ProductUnit{
			ID:   id,
			Name: req.Name,
			Type: domain.ProductUnitType(req.Type),
		}

		user := domain.FromBackendUserContext(ctx)
		err = h.ProductUnitInteractor.Update(ctx, unit, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductUnitNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.ProductUnitNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update product unit: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete
//
//	@Tags		商品单位
//	@Security	BearerAuth
//	@Summary	删除商品单位
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"单位ID"
//	@Success	200	"No Content"
//	@Router		/product/unit/{id} [delete]
func (h *ProductUnitHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductUnitHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		err = h.ProductUnitInteractor.Delete(ctx, id, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductUnitDeleteHasProducts) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.ProductUnitDeleteHasProducts, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}

			err = fmt.Errorf("failed to delete product unit: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// List
//
//	@Tags		商品单位
//	@Security	BearerAuth
//	@Summary	获取商品单位列表
//	@Param		data	query		types.ProductUnitListReq	true	"请求信息"
//	@Success	200		{object}	domain.ProductUnitSearchRes	"成功"
//	@Router		/product/unit [get]
func (h *ProductUnitHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductUnitHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductUnitListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		page := upagination.New(req.Page, req.Size)
		user := domain.FromBackendUserContext(ctx)

		params := domain.ProductUnitSearchParams{
			MerchantID:   user.MerchantID,
			Name:         req.Name,
			Type:         domain.ProductUnitType(req.Type),
			OnlyMerchant: true,
		}

		res, err := h.ProductUnitInteractor.PagedListBySearch(ctx, page, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to list product units: %w", err)
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}
