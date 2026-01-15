package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"go.uber.org/fx"
)

type UserHandler struct {
	UserInteractor domain.StoreUserInteractor
	UserSeq        domain.IncrSequence
}

type UserHandlerParams struct {
	fx.In
	UserInteractor domain.StoreUserInteractor
	UserSeq        domain.IncrSequence `name:"store_user_seq"`
}

func NewUserHandler(p UserHandlerParams) *UserHandler {
	return &UserHandler{UserInteractor: p.UserInteractor, UserSeq: p.UserSeq}
}

func (h *UserHandler) Routes(r gin.IRouter) {
	r = r.Group("/user")
	r.GET("/:id", h.Get())
	r.GET("", h.List())
}

func (h *UserHandler) NoAuths() []string {
	return []string{}
}

// Get 门店后台用户详情
//
//	@Tags			用户管理
//	@Summary		获取门店用户
//	@Description	查询指定门店用户详情
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"用户ID"
//	@Success		200	{object}	domain.StoreUser
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
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.UserNotFound, err))
				return
			}
			err = fmt.Errorf("failed to get store user: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, user)
	}
}

// List 门店后台用户列表
//
//	@Tags			用户管理
//	@Summary		门店用户列表
//	@Description	查询门店用户列表
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	query		types.AccountListReq	true	"门店用户列表请求"
//	@Success		200		{object}	types.AccountListResp
//	@Router			/user [get]
func (h *UserHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("UserHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AccountListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		user := domain.FromFrontendUserContext(ctx)
		filter := &domain.StoreUserListFilter{
			Enabled:    req.Enabled,
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
		}

		pager := upagination.New(1, upagination.MaxSize)
		users, total, err := h.UserInteractor.GetUsers(ctx, pager, filter, domain.NewStoreUserOrderByCreatedAt(true))
		if err != nil {
			err = fmt.Errorf("failed to get store users: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.AccountListResp{Users: users, Total: total})
	}
}
