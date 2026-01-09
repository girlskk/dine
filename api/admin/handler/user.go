package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/i18n"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"go.uber.org/fx"
)

type UserHandler struct {
	UserInteractor domain.AdminUserInteractor
	UserSeq        domain.IncrSequence
}

type UserHandlerParams struct {
	fx.In
	UserInteractor domain.AdminUserInteractor
	UserSeq        domain.IncrSequence `name:"admin_user_seq"`
}

func NewUserHandler(p UserHandlerParams) *UserHandler {
	return &UserHandler{UserInteractor: p.UserInteractor, UserSeq: p.UserSeq}
}

func (h *UserHandler) Routes(r gin.IRouter) {
	r = r.Group("/user")
	r.POST("/login", h.Login())
	r.POST("/logout", h.Logout())
	r.POST("/info", h.Info())
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("/:id", h.Get())
	r.GET("", h.List())
	r.PUT("/:id/enable", h.Enable())
	r.PUT("/:id/disable", h.Disable())
	r.PUT("/:id/reset_password", h.ResetPassword())
}

func (h *UserHandler) NoAuths() []string {
	return []string{
		"/user/login",
	}
}

// Login
//
//	@Tags		用户管理
//	@Security	BearerAuth
//	@Summary	用户登录
//	@Accept		json
//	@Produce	json
//	@Param		data	body		types.LoginReq	true	"请求信息"
//	@Success	200		{object}	types.LoginResp	"成功"
//	@Router		/user/login [post]
func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("UserHandler.Login")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.LoginReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		token, expAt, err := h.UserInteractor.Login(ctx, req.Username, req.Password)
		if err != nil {
			if domain.IsNotFound(err) {
				// 自定义错误，手动翻译
				translated := i18n.Translate(ctx, errcode.UserNotFound.String(), map[string]any{
					"Username": req.Username,
				})
				c.Error(errorx.New(http.StatusBadRequest, errcode.UserNotFound, err).WithMessage(translated))
				return
			}

			if errors.Is(err, domain.ErrMismatchedHashAndPassword) {
				// 默认错误，使用errcode
				c.Error(errorx.New(http.StatusBadRequest, errcode.UserNotFound, err))
				return
			}

			err = fmt.Errorf("failed to authenticate user: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, &types.LoginResp{
			Token:  token,
			Expire: expAt.Unix(),
		})
	}
}

// Logout
//
//	@Tags		用户管理
//	@Security	BearerAuth
//	@Summary	用户登出
//	@Accept		json
//	@Produce	json
//	@Success	200	"No Content"
//	@Router		/user/logout [post]
func (h *UserHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("UserHandler.Logout")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		if err := h.UserInteractor.Logout(ctx); err != nil {
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Info
//
//	@Tags		用户管理
//	@Security	BearerAuth
//	@Summary	获取当前用户信息
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	domain.AdminUser	"成功"
//	@Router		/user/info [post]
func (h *UserHandler) Info() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("UserHandler.Info")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromAdminUserContext(ctx)
		response.Ok(c, user)
	}
}

// Create 创建管理员用户
//
//	@Tags			用户管理
//	@Summary		创建管理员用户
//	@Description	新建一个管理员用户
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.AdminUserCreateReq	true	"创建管理员用户请求"
//	@Success		200		"No Content"
//	@Router			/user [post]
func (h *UserHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("UserHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AdminUserCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		userCode, err := h.generateUserCode(ctx)
		if err != nil {
			err = fmt.Errorf("failed to generate admin user code: %w", err)
			c.Error(err)
			return
		}
		createUser := &domain.AdminUser{
			Username:     req.Username,
			Nickname:     req.Nickname,
			DepartmentID: req.DepartmentID,
			Code:         userCode,
			RealName:     req.RealName,
			Gender:       req.Gender,
			Email:        req.Email,
			PhoneNumber:  req.PhoneNumber,
			Enabled:      req.Enabled,
			RoleIDs:      req.RoleIDs,
		}
		if err := createUser.SetPassword(req.Password); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		if err := h.UserInteractor.Create(ctx, createUser); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新管理员用户
//
//	@Tags			用户管理
//	@Summary		更新管理员用户
//	@Description	修改指定管理员用户信息
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string						true	"管理员用户ID"
//	@Param			data	body	types.AdminUserUpdateReq	true	"更新管理员用户请求"
//	@Success		200		"No Content"
//	@Router			/user/{id} [put]
func (h *UserHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("UserHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.AdminUserUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := &domain.AdminUser{
			ID:           id,
			Username:     req.Username,
			Nickname:     req.Nickname,
			DepartmentID: req.DepartmentID,
			RealName:     req.RealName,
			Gender:       req.Gender,
			Email:        req.Email,
			PhoneNumber:  req.PhoneNumber,
			Enabled:      req.Enabled,
			RoleIDs:      req.RoleIDs,
		}
		if req.Password != "" {
			if err := user.SetPassword(req.Password); err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
		}

		if err := h.UserInteractor.Update(ctx, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Delete 删除管理员用户
//
//	@Tags			用户管理
//	@Summary		删除管理员用户
//	@Description	删除指定管理员用户
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"管理员用户ID"
//	@Success		200	"No Content"
//	@Router			/user/{id} [delete]
func (h *UserHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("UserHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		if err := h.UserInteractor.Delete(ctx, id); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Get 获取管理员用户详情
//
//	@Tags			用户管理
//	@Summary		获取管理员用户
//	@Description	查询指定管理员用户详情
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"管理员用户ID"
//	@Success		200	{object}	domain.AdminUser
//	@Router			/user/{id} [get]
func (h *UserHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("UserHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user, err := h.UserInteractor.GetUser(ctx, id)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, user)
	}
}

// List 管理员用户列表
//
//	@Tags			用户管理
//	@Summary		管理员用户列表
//	@Description	查询管理员用户列表
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	query		types.AdminUserListReq	true	"管理员用户列表请求"
//	@Success		200		{object}	types.AdminUserListResp
//	@Router			/user [get]
func (h *UserHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("UserHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AdminUserListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		pager := req.RequestPagination.ToPagination()
		filter := &domain.AdminUserListFilter{
			Code:        req.Code,
			RealName:    req.RealName,
			Gender:      req.Gender,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			Enabled:     req.Enabled,
		}

		// parse RoleID if provided (parsed but not applied to filter here)
		if req.RoleID != "" {
			rid, err := uuid.Parse(req.RoleID)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			filter.RoleID = rid
		}

		users, total, err := h.UserInteractor.GetUsers(ctx, pager, filter, domain.NewAdminUserOrderByCreatedAt(true))
		if err != nil {
			err = fmt.Errorf("failed to get admin users: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.AdminUserListResp{Users: users, Total: total})
	}
}

// Enable 启用管理员用户
//
//	@Tags			用户管理
//	@Summary		启用管理员用户
//	@Description	启用指定管理员用户
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"管理员用户ID"
//	@Success		200	"No Content"
//	@Router			/user/{id}/enable [put]
func (h *UserHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("UserHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		err = h.UserInteractor.SimpleUpdate(ctx, domain.AdminUserSimpleUpdateFieldEnable, domain.AdminUserSimpleUpdateParams{
			ID:      id,
			Enabled: true,
		})
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 禁用管理员用户
//
//	@Tags			用户管理
//	@Summary		禁用管理员用户
//	@Description	禁用指定管理员用户
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"管理员用户ID"
//	@Success		200	"No Content"
//	@Router			/user/{id}/disable [put]
func (h *UserHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("UserHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		err = h.UserInteractor.SimpleUpdate(ctx, domain.AdminUserSimpleUpdateFieldEnable, domain.AdminUserSimpleUpdateParams{
			ID:      id,
			Enabled: false,
		})
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// ResetPassword 重置密码
//
//	@Tags			用户管理
//	@Security		BearerAuth
//	@Summary		重置用户密码
//	@Description	重置用户密码
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"用户ID"
//	@Param			data	body	types.ResetPasswordReq	true	"重置密码请求"
//	@Success		200		"No Content"
//	@Router			/user/{id}/reset_password [put]
func (h *UserHandler) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("UserHandler.ResetPassword")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		var req types.ResetPasswordReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		err = h.UserInteractor.SimpleUpdate(ctx, domain.AdminUserSimpleUpdateFieldPassword, domain.AdminUserSimpleUpdateParams{
			ID:       id,
			Password: req.NewPassword,
		})
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

func (h *UserHandler) generateUserCode(ctx context.Context) (string, error) {
	seq, err := h.UserSeq.Next(ctx)
	if err != nil {
		return "", err
	}
	return seq, nil
}

func (h *UserHandler) checkErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserNotExists):
		return errorx.New(http.StatusNotFound, errcode.NotFound, err)
	case errors.Is(err, domain.ErrUsernameExist):
		return errorx.New(http.StatusConflict, errcode.Conflict, err)
	case errors.Is(err, domain.ErrSuperUserCannotDelete):
		return errorx.New(http.StatusBadRequest, errcode.SuperUserCannotDelete, err)
	case errors.Is(err, domain.ErrSuperUserCannotDisable):
		return errorx.New(http.StatusBadRequest, errcode.SuperUserCannotDisable, err)
	case errors.Is(err, domain.ErrSuperUserCannotUpdate):
		return errorx.New(http.StatusBadRequest, errcode.SuperUserCannotUpdate, err)
	case errors.Is(err, domain.ErrUserDisabled):
		return errorx.New(http.StatusBadRequest, errcode.UserDisabled, err)
	case errors.Is(err, domain.ErrDepartmentDisabled):
		return errorx.New(http.StatusBadRequest, errcode.DepartmentDisabled, err)
	case errors.Is(err, domain.ErrRoleDisabled):
		return errorx.New(http.StatusBadRequest, errcode.RoleDisabled, err)
	case errors.Is(err, domain.ErrUserRoleRequired):
		return errorx.New(http.StatusBadRequest, errcode.UserRoleRequired, err)
	case errors.Is(err, domain.ErrUserDepartmentRequired):
		return errorx.New(http.StatusBadRequest, errcode.UserDepartmentRequired, err)
	case errors.Is(err, domain.ErrUserRoleTypeMismatch):
		return errorx.New(http.StatusBadRequest, errcode.UserRoleTypeMismatch, err)
	case errors.Is(err, domain.ErrUserDepartmentTypeMismatch):
		return errorx.New(http.StatusBadRequest, errcode.UserDepartmentTypeMismatch, err)
	case errors.Is(err, domain.ErrPasswordCannotBeEmpty):
		return errorx.New(http.StatusBadRequest, errcode.PasswordCannotBeEmpty, err)
	case domain.IsNotFound(err):
		return errorx.New(http.StatusNotFound, errcode.NotFound, err)
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return err
	}
}
