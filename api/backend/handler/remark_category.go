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
)

type RemarkCategoryHandler struct {
	RemarkCategoryInteractor domain.RemarkCategoryInteractor
}

func NewRemarkCategoryHandler(remarkCategoryInteractor domain.RemarkCategoryInteractor) *RemarkCategoryHandler {
	return &RemarkCategoryHandler{
		RemarkCategoryInteractor: remarkCategoryInteractor,
	}
}

func (h *RemarkCategoryHandler) Routes(r gin.IRouter) {
	r = r.Group("/remark_category")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("", h.GetRemarkCategories())
}

// Create 创建备注分类
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		创建备注分类
//	@Description	创建备注分类
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.RemarkCategoryCreateReq	true	"请求信息"
//	@Success		200		"No Content"
//	@Failure		400		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/remark_category [post]
func (h *RemarkCategoryHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkCategoryHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RemarkCategoryCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		remarkCategory := &domain.RemarkCategory{
			Name:        req.Name,
			RemarkScene: req.RemarkScene,
			Description: req.Description,
			SortOrder:   req.SortOrder,
			MerchantID:  user.MerchantID,
		}

		if err := h.RemarkCategoryInteractor.Create(ctx, remarkCategory); err != nil {
			if errors.Is(err, domain.ErrRemarkCategoryNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.RemarkCategoryNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create remark category: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新备注分类
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		更新备注分类
//	@Description	更新备注分类
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string							true	"备注分类ID"
//	@Param			data	body	types.RemarkCategoryUpdateReq	true	"请求信息"
//	@Success		200		"No Content"
//	@Failure		400		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/remark_category/{id} [put]
func (h *RemarkCategoryHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkCategoryHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.RemarkCategoryUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		remarkCategory := &domain.RemarkCategory{
			ID:          id,
			Name:        req.Name,
			RemarkScene: req.RemarkScene,
			Description: req.Description,
			SortOrder:   req.SortOrder,
			MerchantID:  user.MerchantID,
		}

		if err := h.RemarkCategoryInteractor.Update(ctx, remarkCategory); err != nil {
			if errors.Is(err, domain.ErrRemarkCategoryNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.RemarkCategoryNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update remark category: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete 删除备注分类
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		删除备注分类
//	@Description	删除备注分类
//	@Param			id	path	string	true	"备注分类ID"
//	@Success		200	"No Content"
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/remark_category/{id} [delete]
func (h *RemarkCategoryHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkCategoryHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		if err := h.RemarkCategoryInteractor.Delete(ctx, id); err != nil {
			if errors.Is(err, domain.ErrRemarkCategoryNotExists) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to delete remark category: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// GetRemarkCategories 获取备注分类列表
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		获取备注分类列表
//	@Description	获取备注分类列表
//	@Param			data	query		types.RemarkCategoryListReq	true	"查询参数"
//	@Success		200		{object}	response.Response{data=types.RemarkCategoryListResp}
//	@Failure		400		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/remark_category [get]
func (h *RemarkCategoryHandler) GetRemarkCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkCategoryHandler.GetRemarkCategories")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RemarkCategoryListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		filter := &domain.RemarkCategoryListFilter{
			MerchantID: user.MerchantID,
		}

		remarkCategories, err := h.RemarkCategoryInteractor.GetRemarkCategories(ctx, filter)
		if err != nil {
			err = fmt.Errorf("failed to get remark categories: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.RemarkCategoryListResp{RemarkCategories: remarkCategories})
	}
}
