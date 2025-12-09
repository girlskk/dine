package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/customer/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type CustomerHandler struct {
	CustomerInteractor domain.CustomerInteractor
}

func NewCustomerHandler(customerInteractor domain.CustomerInteractor) *CustomerHandler {
	return &CustomerHandler{
		CustomerInteractor: customerInteractor,
	}
}

func (h *CustomerHandler) Routes(r gin.IRouter) {
	r = r.Group("/customer")
	r.POST("/wx_login", h.WXLogin())
	r.POST("/logout", h.Logout())
	r.POST("/info", h.Info())
}

func (h *CustomerHandler) NoAuths() []string {
	return []string{
		"/customer/wx_login",
	}
}

// WXLogin
//
//	@Tags		客户管理
//	@Security	BearerAuth
//	@Summary	微信小程序登录
//	@Accept		json
//	@Produce	json
//	@Param		data	body		types.WXLoginReq	true	"请求信息"
//	@Success	200		{object}	types.WXLoginResp	"成功"
//	@Router		/customer/wx_login [post]
func (h *CustomerHandler) WXLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "CustomerHandler.WXLogin")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("CustomerHandler.WXLogin")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.WXLoginReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		token, expAt, err := h.CustomerInteractor.WXLogin(ctx, req.Code)
		if err != nil {
			c.Error(err)
			return
		}

		response.Ok(c, &types.WXLoginResp{
			Token:  token,
			Expire: expAt.Unix(),
		})
	}
}

// Logout
//
//	@Tags		客户管理
//	@Security	BearerAuth
//	@Summary	客户登出
//	@Accept		json
//	@Produce	json
//	@Success	200	"No Content"
//	@Router		/customer/logout [post]
func (h *CustomerHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "CustomerHandler.Logout")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("CustomerHandler.Logout")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		if err := h.CustomerInteractor.Logout(ctx); err != nil {
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Info
//
//	@Tags		客户管理
//	@Security	BearerAuth
//	@Summary	获取当前客户信息
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	domain.Customer	"成功"
//	@Router		/customer/info [post]
func (h *CustomerHandler) Info() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "CustomerHandler.Info")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("CustomerHandler.Info")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		customer := domain.FromCustomerContext(ctx)
		response.Ok(c, customer)
	}
}
