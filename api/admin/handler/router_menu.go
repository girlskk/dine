package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	r.GET("", h.List())
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
		Enabled := true
		filter := &domain.RouterMenuListFilter{
			UserType: domain.UserTypeAdmin,
			Enabled:  &Enabled,
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
