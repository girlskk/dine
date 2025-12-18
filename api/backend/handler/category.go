package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type CategoryHandler struct {
	CategoryInteractor domain.CategoryInteractor
}

func NewCategoryHandler(categoryInteractor domain.CategoryInteractor) *CategoryHandler {
	return &CategoryHandler{
		CategoryInteractor: categoryInteractor,
	}
}

func (h *CategoryHandler) Routes(r gin.IRouter) {
	r = r.Group("product/category")
	r.POST("", h.CreateRoot())
	r.POST("/:id", h.CreateChild())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("", h.List())
}

func (h *CategoryHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		商品分类
//	@Security	BearerAuth
//	@Summary	创建一级商品分类
//	@Param		data	body	types.CategoryCreateRootReq	true	"请求信息"
//	@Success	200
//	@Router		/product/category [post]
func (h *CategoryHandler) CreateRoot() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("CategoryHandler.CreateRoot")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CategoryCreateRootReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		category := &domain.Category{
			ID:         uuid.New(),
			Name:       req.Name,
			MerchantID: user.MerchantID,
		}

		if req.TaxRateID != nil {
			category.TaxRateID = *req.TaxRateID
		}

		if req.StallID != nil {
			category.StallID = *req.StallID
		}

		if len(req.ChildrenNames) > 0 {
			req.ChildrenNames = lo.Uniq(req.ChildrenNames)
			category.Childrens = make([]*domain.Category, 0, len(req.ChildrenNames))
			for _, name := range req.ChildrenNames {
				category.Childrens = append(category.Childrens, &domain.Category{
					ID:         uuid.New(),
					Name:       name,
					MerchantID: user.MerchantID,
					ParentID:   category.ID,
					// 默认继承父分类的税率和出品部门
					InheritTaxRate: true,
					InheritStall:   true,
				})
			}
		}

		err := h.CategoryInteractor.CreateRoot(ctx, category)

		if err != nil {
			if errors.Is(err, domain.ErrCategoryNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.CategoryNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create root category: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// CreateChild
//
//	@Tags		商品分类
//	@Security	BearerAuth
//	@Summary	创建二级商品分类
//	@Param		id		path	string							true	"父分类ID"
//	@Param		data	body	types.CategoryCreateChildReq	true	"请求信息"
//	@Success	200
//	@Router		/product/category/{id} [post]
func (h *CategoryHandler) CreateChild() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("CategoryHandler.CreateChild")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取父分类ID
		parentIDStr := c.Param("id")
		parentID, err := uuid.Parse(parentIDStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.CategoryCreateChildReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		if !req.InheritTaxRate && req.TaxRateID == nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, nil))
			return
		}

		if !req.InheritStall && req.StallID == nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, nil))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		category := &domain.Category{
			ID:             uuid.New(),
			Name:           req.Name,
			ParentID:       parentID,
			MerchantID:     user.MerchantID,
			InheritTaxRate: req.InheritTaxRate,
			InheritStall:   req.InheritStall,
		}

		// 如果设置了税率ID，则不继承
		if req.TaxRateID != nil {
			category.TaxRateID = *req.TaxRateID
			category.InheritTaxRate = false
		}

		// 如果设置了出品部门ID，则不继承
		if req.StallID != nil {
			category.StallID = *req.StallID
			category.InheritStall = false
		}

		err = h.CategoryInteractor.CreateChild(ctx, category)

		if err != nil {
			if errors.Is(err, domain.ErrCategoryNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.CategoryNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create child category: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update
//
//	@Tags		商品分类
//	@Security	BearerAuth
//	@Summary	更新商品分类
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string					true	"分类ID"
//	@Param		data	body	types.UpdateCategoryReq	true	"请求信息"
//	@Success	200
//	@Router		/product/category/{id} [put]
func (h *CategoryHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("CategoryHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取分类ID
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("invalid category id: %w", err)))
			return
		}

		var req types.UpdateCategoryReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		// 将请求数据映射到 domain.Category
		category := &domain.Category{
			ID:             id,
			Name:           req.Name,
			InheritTaxRate: req.InheritTaxRate,
			InheritStall:   req.InheritStall,
		}

		// 如果设置了税率ID，则不继承
		if req.TaxRateID != nil {
			category.TaxRateID = *req.TaxRateID
			category.InheritTaxRate = false
		}

		// 如果设置了出品部门ID，则不继承
		if req.StallID != nil {
			category.StallID = *req.StallID
			category.InheritStall = false
		}

		err = h.CategoryInteractor.Update(ctx, category)
		if err != nil {
			if errors.Is(err, domain.ErrCategoryNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.CategoryNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update category: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete
//
//	@Tags		商品分类
//	@Security	BearerAuth
//	@Summary	删除商品分类
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"分类ID"
//	@Success	200	"No Content"
//	@Router		/product/category/{id} [delete]
func (h *CategoryHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("CategoryHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		err = h.CategoryInteractor.Delete(ctx, id)
		if err != nil {
			if errors.Is(err, domain.ErrCategoryDeleteHasChildren) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.CategoryDeleteHasChildren, err))
				return
			}
			if errors.Is(err, domain.ErrCategoryDeleteHasProducts) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.CategoryDeleteHasProducts, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}

			err = fmt.Errorf("failed to delete category: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// List
//
//	@Tags		商品分类
//	@Security	BearerAuth
//	@Summary	获取商品分类列表
//	@Success	200	{array}	domain.Categories	"成功"
//	@Router		/product/category [get]
func (h *CategoryHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("CategoryHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromBackendUserContext(ctx)

		params := domain.CategorySearchParams{
			MerchantID: user.MerchantID,
		}

		res, err := h.CategoryInteractor.ListBySearch(ctx, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to paged list categories: %w", err)
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}
