package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/i18n"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type UserHandler struct {
	UserInteractor domain.AdminUserInteractor
}

func NewUserHandler(userInteractor domain.AdminUserInteractor) *UserHandler {
	return &UserHandler{
		UserInteractor: userInteractor,
	}
}

func (h *UserHandler) Routes(r gin.IRouter) {
	r = r.Group("/user")
	r.POST("/login", h.Login())
	r.POST("/logout", h.Logout())
	r.POST("/info", h.Info())
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
		logger := logging.FromContext(ctx).Named("UserHandler.Login")
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
