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

type DepartmentHandler struct {
	Interactor domain.DepartmentInteractor
}

func NewDepartmentHandler(interactor domain.DepartmentInteractor) *DepartmentHandler {
	return &DepartmentHandler{Interactor: interactor}
}

func (h *DepartmentHandler) Routes(r gin.IRouter) {
	r = r.Group("/common/department")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("/:id", h.Get())
	r.GET("", h.List())
	r.PUT("/:id/enable", h.Enable())
	r.PUT("/:id/disable", h.Disable())
}

// Create 创建部门
//
//	@Tags			部门管理
//	@Summary		创建部门
//	@Description	新建一个部门
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.DepartmentCreateReq	true	"创建部门请求"
//	@Success		200		"No Content"
//	@Failure		400		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/common/department [post]
func (h *DepartmentHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DepartmentHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.DepartmentCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		dept := &domain.CreateDepartmentParams{
			MerchantID:     req.MerchantID,
			StoreID:        req.StoreID,
			Name:           req.Name,
			Code:           req.Code,
			DepartmentType: req.DepartmentType,
			Enable:         req.Enable,
		}

		if err := h.Interactor.CreateDepartment(ctx, dept); err != nil {
			if errors.Is(err, domain.ErrDepartmentNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if errors.Is(err, domain.ErrDepartmentCodeExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create department: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新部门
//
//	@Tags			部门管理
//	@Summary		更新部门
//	@Description	修改指定部门
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string						true	"部门ID"
//	@Param			data	body	types.DepartmentUpdateReq	true	"更新部门请求"
//	@Success		200		"No Content"
//	@Failure		400		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/common/department/{id} [put]
func (h *DepartmentHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DepartmentHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.DepartmentUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		dept := &domain.UpdateDepartmentParams{
			ID:             id,
			Name:           req.Name,
			Code:           req.Code,
			DepartmentType: req.DepartmentType,
			Enable:         req.Enable,
		}

		if err := h.Interactor.UpdateDepartment(ctx, dept); err != nil {
			if errors.Is(err, domain.ErrDepartmentNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if errors.Is(err, domain.ErrDepartmentCodeExists) {
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
			err = fmt.Errorf("failed to update department: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete 删除部门
//
//	@Tags			部门管理
//	@Summary		删除部门
//	@Description	删除指定部门
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"部门ID"
//	@Success		200	"No Content"
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/common/department/{id} [delete]
func (h *DepartmentHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DepartmentHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		if err := h.Interactor.DeleteDepartment(ctx, id); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to delete department: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Get 获取部门详情
//
//	@Tags			部门管理
//	@Summary		获取部门
//	@Description	查询指定部门详情
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"部门ID"
//	@Success		200	{object}	domain.Department
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/common/department/{id} [get]
func (h *DepartmentHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DepartmentHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		dept, err := h.Interactor.GetDepartment(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to get department: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, dept)
	}
}

// List 部门列表
//
//	@Tags			部门管理
//	@Summary		部门列表
//	@Description	查询部门列表
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	query		types.DepartmentListReq	true	"部门列表请求"
//	@Success		200		{object}	types.DepartmentListResp
//	@Failure		400		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/common/department [get]
func (h *DepartmentHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DepartmentHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.DepartmentListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		pager := req.RequestPagination.ToPagination()
		filter := &domain.DepartmentListFilter{
			MerchantID:     req.MerchantID,
			StoreID:        req.StoreID,
			Name:           req.Name,
			Code:           req.Code,
			DepartmentType: req.DepartmentType,
			Enable:         req.Enable,
		}

		depts, total, err := h.Interactor.GetDepartments(ctx, pager, filter, domain.NewDepartmentListOrderByCreatedAt(true))
		if err != nil {
			err = fmt.Errorf("failed to get departments: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.DepartmentListResp{Departments: depts, Total: total})
	}
}

// Enable 启用部门
//
//	@Tags			部门管理
//	@Summary		启用部门
//	@Description	将部门设置为启用
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"部门ID"
//	@Success		200	"No Content"
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/common/department/{id}/enable [put]
func (h *DepartmentHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DepartmentHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		dept := &domain.UpdateDepartmentParams{ID: id, Enable: true}
		if err := h.Interactor.UpdateDepartment(ctx, dept); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to enable department: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 禁用部门
//
//	@Tags			部门管理
//	@Summary		禁用部门
//	@Description	将部门设置为禁用
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"部门ID"
//	@Success		200	"No Content"
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/common/department/{id}/disable [put]
func (h *DepartmentHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DepartmentHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		dept := &domain.UpdateDepartmentParams{ID: id, Enable: false}
		if err := h.Interactor.UpdateDepartment(ctx, dept); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to disable department: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}
