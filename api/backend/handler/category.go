package handler

import (
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
	// r.PUT("/:id", h.Update())
	// r.DELETE("/:id", h.Delete())
	// r.GET("/:id", h.GetByID())
	// r.GET("", h.List())
}

func (h *CategoryHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		商品分类
//	@Security	BearerAuth
//	@Summary	创建一级商品分类
//	@Param		data	body	types.CreateRootCategoryReq	true	"请求信息"
//	@Success	200
//	@Router		/product/category [post]
func (h *CategoryHandler) CreateRoot() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("CategoryHandler.CreateRoot")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CreateRootCategoryReq
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
			if domain.IsConflict(err) {
				c.Error(errorx.New(http.StatusConflict, errcode.CategoryNameExists, err))
				return
			}
			err = fmt.Errorf("failed to create category: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// // Update
// //
// //	@Tags		商品分类
// //	@Security	BearerAuth
// //	@Summary	更新商品分类
// //	@Accept		json
// //	@Produce	json
// //	@Param		id		path		string					true	"分类ID"
// //	@Param		data	body		types.UpdateCategoryReq	true	"请求信息"
// //	@Success	200		{object}	types.CategoryResp		"成功"
// //	@Router		/category/{id} [put]
// func (h *CategoryHandler) Update() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx := c.Request.Context()
// 		logger := logging.FromContext(ctx).Named("CategoryHandler.Update")
// 		ctx = logging.NewContext(ctx, logger)
// 		c.Request = c.Request.Clone(ctx)

// 		idStr := c.Param("id")
// 		id, err := uuid.Parse(idStr)
// 		if err != nil {
// 			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("invalid category id: %w", err)))
// 			return
// 		}

// 		var req types.UpdateCategoryReq
// 		if err := c.ShouldBind(&req); err != nil {
// 			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
// 			return
// 		}

// 		category, err := h.CategoryInteractor.Update(ctx, domain.CategoryUpdateParams{
// 			ID:           id,
// 			Name:         req.Name,
// 			TaxRateID:    req.TaxRateID,
// 			DepartmentID: req.DepartmentID,
// 		})
// 		if err != nil {
// 			if domain.IsParamsError(err) {
// 				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
// 				return
// 			}
// 			if domain.IsConflict(err) {
// 				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
// 				return
// 			}
// 			if domain.IsNotFound(err) {
// 				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
// 				return
// 			}
// 			err = fmt.Errorf("failed to update category: %w", err)
// 			c.Error(err)
// 			return
// 		}

// 		response.Ok(c, convertCategoryToResp(category))
// 	}
// }

// // Delete
// //
// //	@Tags		商品分类
// //	@Security	BearerAuth
// //	@Summary	删除商品分类
// //	@Accept		json
// //	@Produce	json
// //	@Param		id	path		string	true	"分类ID"
// //	@Success	200	"No Content"
// //	@Router		/category/{id} [delete]
// func (h *CategoryHandler) Delete() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx := c.Request.Context()
// 		logger := logging.FromContext(ctx).Named("CategoryHandler.Delete")
// 		ctx = logging.NewContext(ctx, logger)
// 		c.Request = c.Request.Clone(ctx)

// 		idStr := c.Param("id")
// 		id, err := uuid.Parse(idStr)
// 		if err != nil {
// 			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("invalid category id: %w", err)))
// 			return
// 		}

// 		err = h.CategoryInteractor.Delete(ctx, id)
// 		if err != nil {
// 			if domain.IsParamsError(err) {
// 				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
// 				return
// 			}
// 			if domain.IsNotFound(err) {
// 				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
// 				return
// 			}
// 			err = fmt.Errorf("failed to delete category: %w", err)
// 			c.Error(err)
// 			return
// 		}

// 		response.Ok(c, nil)
// 	}
// }

// // GetByID
// //
// //	@Tags		商品分类
// //	@Security	BearerAuth
// //	@Summary	获取商品分类详情
// //	@Accept		json
// //	@Produce	json
// //	@Param		id		path		string				true	"分类ID"
// //	@Success	200		{object}	types.CategoryResp	"成功"
// //	@Router		/category/{id} [get]
// func (h *CategoryHandler) GetByID() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx := c.Request.Context()
// 		logger := logging.FromContext(ctx).Named("CategoryHandler.GetByID")
// 		ctx = logging.NewContext(ctx, logger)
// 		c.Request = c.Request.Clone(ctx)

// 		idStr := c.Param("id")
// 		id, err := uuid.Parse(idStr)
// 		if err != nil {
// 			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("invalid category id: %w", err)))
// 			return
// 		}

// 		category, err := h.CategoryInteractor.GetByID(ctx, id)
// 		if err != nil {
// 			if domain.IsNotFound(err) {
// 				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
// 				return
// 			}
// 			err = fmt.Errorf("failed to get category: %w", err)
// 			c.Error(err)
// 			return
// 		}

// 		response.Ok(c, convertCategoryToResp(category))
// 	}
// }

// // List
// //
// //	@Tags		商品分类
// //	@Security	BearerAuth
// //	@Summary	获取商品分类列表
// //	@Accept		json
// //	@Produce	json
// //	@Param		parent_id	query		string				false	"父分类ID，为空表示查询一级分类"
// //	@Param		store_id	query		string				true	"门店ID"
// //	@Success	200			{array}		types.CategoryResp	"成功"
// //	@Router		/category [get]
// func (h *CategoryHandler) List() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx := c.Request.Context()
// 		logger := logging.FromContext(ctx).Named("CategoryHandler.List")
// 		ctx = logging.NewContext(ctx, logger)
// 		c.Request = c.Request.Clone(ctx)

// 		var req types.ListCategoryReq
// 		if err := c.ShouldBindQuery(&req); err != nil {
// 			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
// 			return
// 		}

// 		storeIDStr := c.Query("store_id")
// 		if storeIDStr == "" {
// 			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, errors.New("store_id is required")))
// 			return
// 		}

// 		storeID, err := uuid.Parse(storeIDStr)
// 		if err != nil {
// 			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("invalid store_id: %w", err)))
// 			return
// 		}

// 		categories, err := h.CategoryInteractor.ListByStoreID(ctx, storeID, req.ParentID)
// 		if err != nil {
// 			err = fmt.Errorf("failed to list categories: %w", err)
// 			c.Error(err)
// 			return
// 		}

// 		res := make([]types.CategoryResp, 0, len(categories))
// 		for _, category := range categories {
// 			res = append(res, convertCategoryToResp(category))
// 		}

// 		response.Ok(c, res)
// 	}
// }

// func convertCategoryToResp(category *domain.Category) types.CategoryResp {
// 	res := types.CategoryResp{
// 		ID:           category.ID,
// 		Name:         category.Name,
// 		StoreID:      category.StoreID,
// 		ParentID:     category.ParentID,
// 		TaxRateID:    category.TaxRateID,
// 		DepartmentID: category.DepartmentID,
// 		ProductCount: category.ProductCount,
// 		CreatedAt:    category.CreatedAt.Format(time.RFC3339),
// 		UpdatedAt:    category.UpdatedAt.Format(time.RFC3339),
// 	}

// 	if category.Parent != nil {
// 		parentResp := convertCategoryToResp(category.Parent)
// 		res.Parent = &parentResp
// 	}

// 	if len(category.Children) > 0 {
// 		res.Children = make([]types.CategoryResp, 0, len(category.Children))
// 		for _, child := range category.Children {
// 			res.Children = append(res.Children, convertCategoryToResp(child))
// 		}
// 	}

// 	return res
// }
