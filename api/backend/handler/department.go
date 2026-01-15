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

type DepartmentHandler struct {
	Interactor domain.DepartmentInteractor
	DeptSeq    domain.IncrSequence
}

type DepartmentHandlerParams struct {
	fx.In
	Interactor domain.DepartmentInteractor
	DeptSeq    domain.IncrSequence `name:"backend_department_seq"`
}

func NewDepartmentHandler(p DepartmentHandlerParams) *DepartmentHandler {
	return &DepartmentHandler{Interactor: p.Interactor, DeptSeq: p.DeptSeq}
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
//	@Description	在品牌后台创建部门
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.DepartmentCreateReq	true	"创建部门请求"
//	@Success		200		"No Content"
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

		deptCode, err := h.generateDepartmentCode(ctx)
		if err != nil {
			err = fmt.Errorf("failed to generate department code: %w", err)
			c.Error(err)
			return
		}
		user := domain.FromBackendUserContext(ctx)
		params := &domain.CreateDepartmentParams{
			Name:           req.Name,
			Code:           deptCode,
			DepartmentType: domain.DepartmentBackend,
			Enabled:        req.Enabled,
			MerchantID:     user.MerchantID,
		}

		if err := h.Interactor.CreateDepartment(ctx, params, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新部门
//
//	@Tags			部门管理
//	@Summary		更新部门
//	@Description	修改指定品牌部门
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string						true	"部门ID"
//	@Param			data	body	types.DepartmentUpdateReq	true	"更新部门请求"
//	@Success		200		"No Content"
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

		user := domain.FromBackendUserContext(ctx)
		params := &domain.UpdateDepartmentParams{ID: id, Name: req.Name, Enabled: req.Enabled}
		if err := h.Interactor.UpdateDepartment(ctx, params, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Delete 删除部门
//
//	@Tags			部门管理
//	@Summary		删除部门
//	@Description	删除指定品牌部门
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"部门ID"
//	@Success		200	"No Content"
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

		user := domain.FromBackendUserContext(ctx)
		if err := h.Interactor.DeleteDepartment(ctx, id, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Get 获取部门详情
//
//	@Tags			部门管理
//	@Summary		获取部门
//	@Description	查询指定品牌部门详情
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"部门ID"
//	@Success		200	{object}	domain.Department
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

		user := domain.FromBackendUserContext(ctx)
		dept, err := h.Interactor.GetDepartment(ctx, id, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, dept)
	}
}

// List 部门列表
//
//	@Tags			部门管理
//	@Summary		部门列表
//	@Description	分页查询品牌部门
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	query		types.DepartmentListReq	true	"部门列表请求"
//	@Success		200		{object}	types.DepartmentListResp
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

		user := domain.FromBackendUserContext(ctx)
		pager := req.RequestPagination.ToPagination()
		filter := &domain.DepartmentListFilter{
			Name:           req.Name,
			Code:           req.Code,
			DepartmentType: domain.DepartmentBackend,
			Enabled:        req.Enabled,
			MerchantID:     user.MerchantID,
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
//	@Description	启用指定品牌部门
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"部门ID"
//	@Success		200	"No Content"
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

		user := domain.FromBackendUserContext(ctx)
		params := domain.DepartmentSimpleUpdateParams{ID: id, Enabled: true}
		err = h.Interactor.SimpleUpdate(
			ctx,
			domain.DepartmentSimpleUpdateFieldEnabled,
			params,
			user,
		)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 禁用部门
//
//	@Tags			部门管理
//	@Summary		禁用部门
//	@Description	禁用指定品牌部门
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"部门ID"
//	@Success		200	"No Content"
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

		user := domain.FromBackendUserContext(ctx)
		params := domain.DepartmentSimpleUpdateParams{ID: id, Enabled: false}
		err = h.Interactor.SimpleUpdate(
			ctx,
			domain.DepartmentSimpleUpdateFieldEnabled,
			params,
			user,
		)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

func (h *DepartmentHandler) generateDepartmentCode(ctx context.Context) (string, error) {
	if h.DeptSeq == nil {
		return "", fmt.Errorf("department sequence not initialized")
	}
	seq, err := h.DeptSeq.Next(ctx)
	if err != nil {
		return "", err
	}
	return seq, nil
}

func (h *DepartmentHandler) checkErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrDepartmentNotExists):
		return errorx.New(http.StatusBadRequest, errcode.DepartmentNotExists, err)
	case errors.Is(err, domain.ErrDepartmentNameExists):
		return errorx.New(http.StatusConflict, errcode.DepartmentNameExists, err)
	case errors.Is(err, domain.ErrDepartmentCodeExists):
		return errorx.New(http.StatusConflict, errcode.DepartmentCodeExists, err)
	case errors.Is(err, domain.ErrDepartmentHasUsersCannotDisable):
		return errorx.New(http.StatusForbidden, errcode.DepartmentHasUserCannotDisable, err)
	case errors.Is(err, domain.ErrDepartmentHasUsersCannotDelete):
		return errorx.New(http.StatusForbidden, errcode.DepartmentHasUserCannotDelete, err)
	case domain.IsNotFound(err):
		return errorx.New(http.StatusNotFound, errcode.NotFound, err)
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return fmt.Errorf("department handler error: %w", err)
	}
}
