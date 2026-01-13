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
	r.GET("", h.List())
	r.PUT("/:id/enable", h.Enable())
	r.PUT("/:id/disable", h.Disable())
}

// Create 创建备注
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		创建备注
//	@Description	创建备注
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.RemarkCreateReq	true	"请求信息"
//	@Success		200		"No Content"
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

		user := domain.FromStoreUserContext(ctx)
		remark := &domain.CreateRemarkParams{
			Name:        req.Name,
			RemarkType:  domain.RemarkTypeStore,
			Enabled:     req.Enabled,
			SortOrder:   req.SortOrder,
			RemarkScene: req.RemarkScene,
			MerchantID:  user.MerchantID,
			StoreID:     user.StoreID,
		}

		if err := h.RemarkInteractor.Create(ctx, remark, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新备注
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		更新备注
//	@Description	更新备注
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"备注ID"
//	@Param			data	body	types.RemarkUpdateReq	true	"请求信息"
//	@Success		200		"No Content"
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
		user := domain.FromStoreUserContext(ctx)
		remark := &domain.UpdateRemarkParams{
			ID:        id,
			Name:      req.Name,
			Enabled:   req.Enabled,
			SortOrder: req.SortOrder,
		}

		if err := h.RemarkInteractor.Update(ctx, remark, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Delete 删除备注
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		删除备注
//	@Description	删除备注
//	@Param			id	path	string	true	"备注ID"
//	@Success		200	"No Content"
//	@Success		204	"No Content"
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
		user := domain.FromStoreUserContext(ctx)
		if err := h.RemarkInteractor.Delete(ctx, id, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Get 获取备注详情
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		获取备注详情
//	@Description	根据备注ID获取详情
//	@Param			id	path		string	true	"备注ID"
//	@Success		200	{object}	response.Response{data=domain.Remark}
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
		user := domain.FromStoreUserContext(ctx)
		remark, err := h.RemarkInteractor.GetRemark(ctx, id, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, remark)
	}
}

// List 获取备注列表
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		获取备注列表
//	@Description	分页查询备注列表
//	@Param			data	query		types.RemarkListReq	true	"备注列表查询参数"
//	@Success		200		{object}	response.Response{data=types.RemarkListResp}
//	@Router			/remark [get]
func (h *RemarkHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RemarkListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		pager := req.RequestPagination.ToPagination()
		filter := &domain.RemarkListFilter{
			MerchantID:  user.MerchantID,
			StoreID:     user.StoreID,
			Enabled:     req.Enabled,
			RemarkType:  domain.RemarkTypeStore,
			RemarkScene: req.RemarkScene,
		}

		remarks, total, err := h.RemarkInteractor.GetRemarks(ctx, pager, filter)
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

// Enable 启用备注
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		启用备注
//	@Description	将备注置为启用
//	@Produce		json
//	@Param			id	path	string	true	"备注ID"
//	@Success		200	"No Content"
//	@Router			/remark/{id}/enable [put]
func (h *RemarkHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		user := domain.FromStoreUserContext(ctx)
		remark := &domain.Remark{ID: id, Enabled: true}
		if err := h.RemarkInteractor.RemarkSimpleUpdate(ctx, domain.RemarkSimpleUpdateFieldEnabled, remark, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 禁用备注
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		禁用备注
//	@Description	将备注置为禁用
//	@Produce		json
//	@Param			id	path	string	true	"备注ID"
//	@Success		200	"No Content"
//	@Router			/remark/{id}/disable [put]
func (h *RemarkHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		user := domain.FromStoreUserContext(ctx)
		remark := &domain.Remark{ID: id, Enabled: false}
		if err := h.RemarkInteractor.RemarkSimpleUpdate(ctx, domain.RemarkSimpleUpdateFieldEnabled, remark, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

func (h *RemarkHandler) checkErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrRemarkNotExists):
		return errorx.New(http.StatusBadRequest, errcode.RemarkNotExists, err)
	case errors.Is(err, domain.ErrRemarkNameExists):
		return errorx.New(http.StatusConflict, errcode.RemarkNameExists, err)
	case errors.Is(err, domain.ErrRemarkDeleteSystem):
		return errorx.New(http.StatusForbidden, errcode.RemarkDeleteSystem, err)
	case domain.IsNotFound(err):
		return errorx.New(http.StatusNotFound, errcode.NotFound, err)
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return fmt.Errorf("remark handler error: %w", err)
	}
}
