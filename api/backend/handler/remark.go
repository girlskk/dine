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

type RemarkHandler struct {
	RemarkInteractor domain.RemarkInteractor
}

func NewRemarkHandler(remarkInteractor domain.RemarkInteractor) *RemarkHandler {
	return &RemarkHandler{
		RemarkInteractor: remarkInteractor,
	}
}

func (h *RemarkHandler) Routes(r gin.IRouter) {
	r = r.Group("/remark")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("/:id", h.Get())
	r.GET("", h.GetRemarks())
	r.PATCH("/:id", h.RemarkSimpleUpdate())
}

// Create 创建备注
//
//	@Tags			备注
//	@Security		BearerAuth
//	@Summary		创建备注
//	@Description	创建备注
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.RemarkCreateReq	true	"请求信息"
//	@Success		200		"No Content"
//	@Failure		400		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/remark [post]
func (h *RemarkHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RemarkCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		remark := &domain.Remark{
			Name:       req.Name,
			RemarkType: req.RemarkType,
			Enabled:    req.Enabled,
			SortOrder:  req.SortOrder,
			CategoryID: req.CategoryID,
			MerchantID: user.MerchantID,
			StoreID:    req.StoreID,
		}

		if err := h.RemarkInteractor.Create(ctx, remark); err != nil {
			if errors.Is(err, domain.ErrRemarkNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.RemarkNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create remark: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新备注
//
//	@Tags			备注
//	@Security		BearerAuth
//	@Summary		更新备注
//	@Description	更新备注
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"备注ID"
//	@Param			data	body	types.RemarkUpdateReq	true	"请求信息"
//	@Success		200		"No Content"
//	@Failure		400		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/remark/{id} [put]
func (h *RemarkHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.RemarkUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		remark := &domain.Remark{
			ID:         id,
			Name:       req.Name,
			RemarkType: req.RemarkType,
			Enabled:    req.Enabled,
			SortOrder:  req.SortOrder,
			CategoryID: req.CategoryID,
			MerchantID: user.MerchantID,
			StoreID:    req.StoreID,
		}

		if err := h.RemarkInteractor.Update(ctx, remark); err != nil {
			if errors.Is(err, domain.ErrRemarkNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.RemarkNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update remark: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete 删除备注
//
//	@Tags			备注
//	@Security		BearerAuth
//	@Summary		删除备注
//	@Description	删除备注
//	@Param			id	path	string	true	"备注ID"
//	@Success		200	"No Content"
//	@Success		204	"No Content"
//	@Failure		400	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/remark/{id} [delete]
func (h *RemarkHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		if err := h.RemarkInteractor.Delete(ctx, id); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNoContent, errcode.NotFound, err))
				return
			}
			if errors.Is(err, domain.ErrRemarkDeleteSystem) {
				c.Error(errorx.New(http.StatusForbidden, errcode.RemarkDeleteSystem, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to delete remark: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Get 获取备注详情
//
//	@Tags			备注
//	@Security		BearerAuth
//	@Summary		获取备注详情
//	@Description	根据备注ID获取详情
//	@Param			id	path		string	true	"备注ID"
//	@Success		200	{object}	response.Response{data=domain.Remark}
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/remark/{id} [get]
func (h *RemarkHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		remark, err := h.RemarkInteractor.GetRemark(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to get remark: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, remark)
	}
}

// GetRemarks 获取备注列表
//
//	@Tags			备注
//	@Security		BearerAuth
//	@Summary		获取备注列表
//	@Description	分页查询备注列表
//	@Param			data	query		types.RemarkListReq	true	"备注列表查询参数"
//	@Success		200		{object}	response.Response{data=types.RemarkListResp}
//	@Failure		400		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/remark [get]
func (h *RemarkHandler) GetRemarks() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkHandler.GetRemarks")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RemarkListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		pager := req.RequestPagination.ToPagination()
		filter := &domain.RemarkListFilter{
			CategoryID: req.CategoryID,
			MerchantID: user.MerchantID,
			StoreID:    req.StoreID,
			Enabled:    req.Enabled,
			RemarkType: req.RemarkType,
		}

		remarks, total, err := h.RemarkInteractor.GetRemarks(ctx, pager, filter, domain.NewRemarkOrderByCreatedAt(true))
		if err != nil {
			err = fmt.Errorf("failed to get remarks: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.RemarkListResp{
			Remarks: remarks,
			Total:   total,
		})
	}
}

// RemarkSimpleUpdate 快速切换启用状态
//
//	@Tags			备注
//	@Security		BearerAuth
//	@Summary		快速切换启用状态
//	@Description	快速切换启用状态
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string						true	"备注ID"
//	@Param			data	body	types.RemarkSimpleUpdateReq	true	"请求信息"
//	@Success		200		"No Content"
//	@Failure		400		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/remark/{id} [patch]
func (h *RemarkHandler) RemarkSimpleUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkHandler.RemarkSimpleUpdate")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.RemarkSimpleUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		remark := &domain.Remark{ID: id, Enabled: req.Enabled}
		if err := h.RemarkInteractor.RemarkSimpleUpdate(ctx, req.SimpleUpdateType, remark); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to toggle enabled: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}
