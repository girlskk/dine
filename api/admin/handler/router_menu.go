package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type RouterMenuHandler struct {
	Interactor domain.RouterMenuInteractor
}

func NewRouterMenuHandler(interactor domain.RouterMenuInteractor) *RouterMenuHandler {
	return &RouterMenuHandler{Interactor: interactor}
}

func (h *RouterMenuHandler) Routes(r gin.IRouter) {
	r = r.Group("/common/router_menu")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("/:id", h.Get())
	r.GET("", h.List())
	r.PUT("/:id/enable", h.Enable())
	r.PUT("/:id/disable", h.Disable())
}

// Create 创建菜单
//
//	@Tags			菜单管理
//	@Summary		创建菜单
//	@Description	新建一个菜单
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.RouterMenuCreateReq	true	"创建菜单请求"
//	@Success		200		"No Content"
//	@Router			/common/router_menu [post]
func (h *RouterMenuHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RouterMenuHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RouterMenuCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		params := &domain.CreateRouterMenuParams{
			UserType:  domain.UserTypeAdmin,
			ParentID:  req.ParentID,
			Name:      req.Name,
			Path:      req.Path,
			Component: req.Component,
			Icon:      req.Icon,
			Sort:      req.Sort,
			Enabled:   req.Enabled,
		}

		if err := h.Interactor.CreateRouterMenu(ctx, params); err != nil {
			if errors.Is(err, domain.ErrRouterMenuNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if errors.Is(err, domain.ErrRouterMenuForbidenAddChild) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create router menu: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新菜单
//
//	@Tags			菜单管理
//	@Summary		更新菜单
//	@Description	修改指定菜单
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string						true	"菜单ID"
//	@Param			data	body	types.RouterMenuUpdateReq	true	"更新菜单请求"
//	@Success		200		"No Content"
//	@Router			/common/router_menu/{id} [put]
func (h *RouterMenuHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RouterMenuHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.RouterMenuUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		params := &domain.UpdateRouterMenuParams{
			ID:        id,
			ParentID:  req.ParentID,
			Name:      req.Name,
			Path:      req.Path,
			Component: req.Component,
			Icon:      req.Icon,
			Sort:      req.Sort,
			Enabled:   req.Enabled,
		}

		if err := h.Interactor.UpdateRouterMenu(ctx, params); err != nil {
			if errors.Is(err, domain.ErrRouterMenuNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if errors.Is(err, domain.ErrRouterMenuForbidenAddChild) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
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
			err = fmt.Errorf("failed to update router menu: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete 删除菜单
//
//	@Tags			菜单管理
//	@Summary		删除菜单
//	@Description	删除指定菜单
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	string	true	"菜单ID"
//	@Success		200	"No Content"
//	@Router			/common/router_menu/{id} [delete]
func (h *RouterMenuHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RouterMenuHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		if err := h.Interactor.DeleteRouterMenu(ctx, id); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to delete router menu: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Get 获取菜单详情
//
//	@Tags			菜单管理
//	@Summary		获取菜单详情
//	@Description	查询指定菜单
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"菜单ID"
//	@Success		200	{object}	domain.RouterMenu
//	@Router			/common/router_menu/{id} [get]
func (h *RouterMenuHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RouterMenuHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		menu, err := h.Interactor.GetRouterMenu(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to get router menu: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, menu)
	}
}

// List 菜单列表
//
//	@Tags			菜单管理
//	@Summary		菜单列表
//	@Description	分页查询菜单
//	@Security		BearerAuth
//	@Produce		json
//	@Param			data	query		types.RouterMenuListReq	true	"菜单列表请求"
//	@Success		200		{object}	types.RouterMenuListResp
//	@Router			/common/router_menu [get]
func (h *RouterMenuHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RouterMenuHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RouterMenuListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		enable := true
		filter := &domain.RouterMenuListFilter{
			UserType: domain.UserTypeAdmin,
			Enabled:  &enable,
		}

		menus, total, err := h.Interactor.GetRouterMenus(ctx, filter, domain.NewRouterMenuListOrderBySort(false))
		if err != nil {
			err = fmt.Errorf("failed to get router menus: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.RouterMenuListResp{Menus: menus, Total: total})
	}
}

// Enable 启用菜单
//
//	@Tags			菜单管理
//	@Summary		启用菜单
//	@Description	启用指定菜单
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	string	true	"菜单ID"
//	@Success		200	"No Content"
//	@Router			/common/router_menu/{id}/enable [put]
func (h *RouterMenuHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RouterMenuHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		params := &domain.UpdateRouterMenuParams{ID: id, Enabled: true}
		if err := h.Interactor.UpdateRouterMenu(ctx, params); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to enable router menu: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 禁用菜单
//
//	@Tags			菜单管理
//	@Summary		禁用菜单
//	@Description	禁用指定菜单
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	string	true	"菜单ID"
//	@Success		200	"No Content"
//	@Router			/common/router_menu/{id}/disable [put]
func (h *RouterMenuHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RouterMenuHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		params := &domain.UpdateRouterMenuParams{ID: id, Enabled: false}
		if err := h.Interactor.UpdateRouterMenu(ctx, params); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to disable router menu: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}
