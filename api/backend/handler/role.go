package handler

import (
	"context"
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
	"go.uber.org/fx"
)

type RoleHandler struct {
	Interactor   domain.RoleInteractor
	RoleSequence domain.IncrSequence
}

type RoleHandlerParams struct {
	fx.In
	Interactor   domain.RoleInteractor
	RoleSequence domain.IncrSequence `name:"backend_role_seq"`
}

func NewRoleHandler(p RoleHandlerParams) *RoleHandler {
	return &RoleHandler{Interactor: p.Interactor, RoleSequence: p.RoleSequence}
}

func (h *RoleHandler) Routes(r gin.IRouter) {
	r = r.Group("/common/role")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("/:id", h.Get())
	r.GET("", h.List())
	r.PUT("/:id/enable", h.Enable())
	r.PUT("/:id/disable", h.Disable())
}

// Create 创建角色
//
//	@Tags			角色管理
//	@Summary		创建角色
//	@Description	在品牌后台创建角色
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.RoleCreateReq	true	"创建角色请求"
//	@Success		200		"No Content"
//	@Router			/common/role [post]
func (h *RoleHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RoleHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RoleCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		user := domain.FromBackendUserContext(ctx)

		roleCode, err := h.generateRoleCode(ctx)
		if err != nil {
			err = fmt.Errorf("failed to generate role code: %w", err)
			c.Error(err)
			return
		}

		params := &domain.CreateRoleParams{
			Name:       req.Name,
			Code:       roleCode,
			RoleType:   domain.RoleTypeAdmin,
			DataScope:  domain.RoleDataScopeAll,
			Enable:     req.Enable,
			MerchantID: user.MerchantID,
		}

		if err := h.Interactor.CreateRole(ctx, params); err != nil {
			if errors.Is(err, domain.ErrRoleNameExists) || errors.Is(err, domain.ErrRoleCodeExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create role: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新角色
//
//	@Tags			角色管理
//	@Summary		更新角色
//	@Description	修改指定品牌后台角色
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string				true	"角色ID"
//	@Param			data	body	types.RoleUpdateReq	true	"更新角色请求"
//	@Success		200		"No Content"
//	@Router			/common/role/{id} [put]
func (h *RoleHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RoleHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.RoleUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		params := &domain.UpdateRoleParams{
			ID:        id,
			Name:      req.Name,
			RoleType:  domain.RoleTypeAdmin,
			DataScope: domain.RoleDataScopeAll,
			Enable:    req.Enable,
		}

		if err := h.Interactor.UpdateRole(ctx, params); err != nil {
			if errors.Is(err, domain.ErrRoleNameExists) || errors.Is(err, domain.ErrRoleCodeExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to update role: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete 删除角色
//
//	@Tags			角色管理
//	@Summary		删除角色
//	@Description	删除指定角色
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	string	true	"角色ID"
//	@Success		200	"No Content"
//	@Router			/common/role/{id} [delete]
func (h *RoleHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RoleHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		if err := h.Interactor.DeleteRole(ctx, id); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to delete role: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Get 获取角色详情
//
//	@Tags			角色管理
//	@Summary		获取角色详情
//	@Description	查询指定品牌角色
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"角色ID"
//	@Success		200	{object}	domain.Role
//	@Router			/common/role/{id} [get]
func (h *RoleHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RoleHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		role, err := h.Interactor.GetRole(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to get role: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, role)
	}
}

// List 角色列表
//
//	@Tags			角色管理
//	@Summary		角色列表
//	@Description	分页查询品牌角色
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	query		types.RoleListReq	true	"角色列表请求"
//	@Success		200		{object}	types.RoleListResp
//	@Router			/common/role [get]
func (h *RoleHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RoleHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RoleListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		pager := req.RequestPagination.ToPagination()
		filter := &domain.RoleListFilter{
			Name:       req.Name,
			RoleType:   domain.RoleTypeAdmin,
			Enable:     req.Enable,
			MerchantID: user.MerchantID,
		}

		roles, total, err := h.Interactor.GetRoles(ctx, pager, filter, domain.NewRoleListOrderByCreatedAt(true))
		if err != nil {
			err = fmt.Errorf("failed to get roles: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.RoleListResp{Roles: roles, Total: total})
	}
}

// Enable 启用角色
//
//	@Tags			角色管理
//	@Summary		启用角色
//	@Description	启用指定品牌角色
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"角色ID"
//	@Success		200	"No Content"
//	@Router			/common/role/{id}/enable [put]
func (h *RoleHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RoleHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		params := &domain.UpdateRoleParams{ID: id, Enable: true}
		if err := h.Interactor.UpdateRole(ctx, params); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to enable role: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 禁用角色
//
//	@Tags			角色管理
//	@Summary		禁用角色
//	@Description	禁用指定品牌角色
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"角色ID"
//	@Success		200	"No Content"
//	@Router			/common/role/{id}/disable [put]
func (h *RoleHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RoleHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		params := &domain.UpdateRoleParams{ID: id, Enable: false}
		if err := h.Interactor.UpdateRole(ctx, params); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to disable role: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

func (h *RoleHandler) generateRoleCode(ctx context.Context) (string, error) {
	if h.RoleSequence == nil {
		return "", fmt.Errorf("role sequence not initialized")
	}
	seq, err := h.RoleSequence.Next(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate role code: %w", err)
	}
	return seq, nil
}
