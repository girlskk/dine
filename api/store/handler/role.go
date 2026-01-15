package handler

import (
	"context"
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
	"go.uber.org/fx"
)

type RoleHandler struct {
	Interactor         domain.RoleInteractor
	RoleMenuInteractor domain.RoleMenuInteractor
	RoleSequence       domain.IncrSequence
}

type RoleHandlerParams struct {
	fx.In
	Interactor         domain.RoleInteractor
	RoleMenuInteractor domain.RoleMenuInteractor
	RoleSequence       domain.IncrSequence `name:"store_role_seq"`
}

func NewRoleHandler(p RoleHandlerParams) *RoleHandler {
	return &RoleHandler{Interactor: p.Interactor, RoleMenuInteractor: p.RoleMenuInteractor, RoleSequence: p.RoleSequence}
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
//	@Description	在门店后台创建角色
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
		user := domain.FromStoreUserContext(ctx)

		roleCode, err := h.generateRoleCode(ctx)
		if err != nil {
			err = fmt.Errorf("failed to generate role code: %w", err)
			c.Error(err)
			return
		}

		params := &domain.CreateRoleParams{
			Name:       req.Name,
			Code:       roleCode,
			RoleType:   domain.RoleTypeStore,
			DataScope:  domain.RoleDataScopeAll,
			Enabled:    req.Enabled,
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
		}

		if err := h.Interactor.CreateRole(ctx, params, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新角色
//
//	@Tags			角色管理
//	@Summary		更新角色
//	@Description	修改指定门店后台角色
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

		user := domain.FromStoreUserContext(ctx)
		params := &domain.UpdateRoleParams{
			ID:        id,
			Name:      req.Name,
			RoleType:  domain.RoleTypeStore,
			DataScope: domain.RoleDataScopeAll,
			Enabled:   req.Enabled,
		}

		if err := h.Interactor.UpdateRole(ctx, params, user); err != nil {
			c.Error(h.checkErr(err))
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

		user := domain.FromStoreUserContext(ctx)
		if err := h.Interactor.DeleteRole(ctx, id, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Get 获取角色详情
//
//	@Tags			角色管理
//	@Summary		获取角色详情
//	@Description	查询指定角色
//	@Security		BearerAuth
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

		user := domain.FromStoreUserContext(ctx)
		role, err := h.Interactor.GetRole(ctx, id, user)
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
//	@Description	分页查询门店角色
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

		user := domain.FromStoreUserContext(ctx)

		pager := req.RequestPagination.ToPagination()
		filter := &domain.RoleListFilter{
			Name:       req.Name,
			RoleType:   domain.RoleTypeStore,
			Enabled:    req.Enabled,
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
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
//	@Description	启用指定门店角色
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

		user := domain.FromStoreUserContext(ctx)
		err = h.Interactor.SimpleUpdate(ctx, domain.RoleSimpleUpdateFieldEnabled, domain.RoleSimpleUpdateParams{
			ID:      id,
			Enabled: true,
		}, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 禁用角色
//
//	@Tags			角色管理
//	@Summary		禁用角色
//	@Description	禁用指定门店角色
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

		user := domain.FromStoreUserContext(ctx)
		err = h.Interactor.SimpleUpdate(ctx, domain.RoleSimpleUpdateFieldEnabled, domain.RoleSimpleUpdateParams{
			ID:      id,
			Enabled: false,
		}, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// SetMenus 设置角色菜单
//
//	@Tags			角色管理
//	@Summary		设置角色菜单
//	@Description	为指定角色设置菜单路径（交集保留，新增/删除按 paths 调整）
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string				true	"角色ID"
//	@Param			data	body	types.SetMenusReq	true	"设置菜单请求"
//	@Success		200		"No Content"
//	@Router			/common/role/{id}/menus [post]
func (h *RoleHandler) SetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RoleHandler.SetMenus")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.SetMenusReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		if err := h.RoleMenuInteractor.SetRoleMenu(ctx, id, req.Paths); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// RoleMenuList 角色菜单列表
//
//	@Tags			角色管理
//	@Summary		角色菜单列表
//	@Description	分页或非分页获取指定角色的菜单路径列表
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"角色ID"
//	@Success		200	{object}	types.RoleMenusResp
//	@Router			/common/role/{id}/menus [get]
func (h *RoleHandler) RoleMenuList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RoleHandler.RoleMenuList")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		paths, err := h.RoleMenuInteractor.RoleMenuList(ctx, id)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, types.RoleMenusResp{Paths: paths})
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

func (h *RoleHandler) checkErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserRoleNotExists):
		return errorx.New(http.StatusBadRequest, errcode.UserRoleNotExists, err)
	case errors.Is(err, domain.ErrRoleNameExists), errors.Is(err, domain.ErrRoleCodeExists):
		return errorx.New(http.StatusConflict, errcode.Conflict, err)
	case errors.Is(err, domain.ErrRoleAssignedCannotDisable):
		return errorx.New(http.StatusForbidden, errcode.RoleAssignedCannotDisable, err)
	case errors.Is(err, domain.ErrRoleAssignedCannotDelete):
		return errorx.New(http.StatusForbidden, errcode.RoleAssignedCannotDelete, err)
	case domain.IsNotFound(err):
		return errorx.New(http.StatusNotFound, errcode.NotFound, err)
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return fmt.Errorf("role handler error: %w", err)
	}
}
