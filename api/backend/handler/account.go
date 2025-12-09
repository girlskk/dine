package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type AccountHandler struct {
	FrontendUserInteractor domain.FrontendUserInteractor
}

func NewAccountHandler(interactor domain.FrontendUserInteractor) *AccountHandler {
	return &AccountHandler{
		FrontendUserInteractor: interactor,
	}
}

// Routes 注册路由
func (h *AccountHandler) Routes(r gin.IRouter) {
	r = r.Group("/account")
	r.POST("/create", h.Create())
	r.POST("/list", h.List())
	r.POST("/update", h.Update())
}

// Create 创建账号
//
//	@Tags		账号管理
//	@Summary	创建账号
//	@Security	BearerAuth
//	@Param		data	body	types.AccountCreateReq	true	"请求参数"
//	@Success	200
//	@Router		/account/create [post]
func (h *AccountHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "AccountHandler.Create")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("AccountHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AccountCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		currentUser := domain.FromBackendUserContext(ctx)
		user := &domain.FrontendUser{
			Nickname: req.NickName,
			Username: req.Username,
			StoreID:  currentUser.StoreID,
		}
		err := user.SetPassword(req.Password)
		if err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		if err := h.FrontendUserInteractor.Create(ctx, user); err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, nil)
	}
}

// List 账号列表
//
//	@Tags		账号管理
//	@Summary	账号列表
//	@Security	BearerAuth
//	@Param		data	body		types.AccountListReq	true	"请求参数"
//	@Success	200		{object}	types.AccountListResp	"成功"
//	@Router		/account/list [post]
func (h *AccountHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "AccountHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("AccountHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AccountListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		store := domain.FromStoreContext(ctx)

		users, total, err := h.FrontendUserInteractor.List(ctx, req.ToPagination(), &domain.FrontendUserListFilter{
			StoreID: store.ID,
		})
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, &types.AccountListResp{
			Items: users,
			Total: total,
		})
	}
}

// Update 账号更新
//
//	@Tags		账号管理
//	@Summary	账号更新
//	@Security	BearerAuth
//	@Param		data	body	types.AccountUpdateReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/account/update [post]
func (h *AccountHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "AccountHandler.Update")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("AccountHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.AccountUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		store := domain.FromStoreContext(ctx)

		user, err := h.FrontendUserInteractor.Find(ctx, req.ID)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		if user.StoreID != store.ID {
			c.Error(uerr.BadRequest(domain.ErrUserNotFound.Error()))
			return
		}

		if req.Password != "" {
			err := user.SetPassword(req.Password)
			if err != nil {
				c.Error(uerr.BadRequest(err.Error()))
				return
			}
		}

		if req.NickName != "" {
			user.Nickname = req.NickName
		}

		if req.Username != "" {
			user.Username = req.Username
		}

		if err := h.FrontendUserInteractor.Update(ctx, user); err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, nil)
	}
}
