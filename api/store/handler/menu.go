package handler

import (
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
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type MenuHandler struct {
	MenuInteractor domain.MenuInteractor
}

func NewMenuHandler(menuInteractor domain.MenuInteractor) *MenuHandler {
	return &MenuHandler{
		MenuInteractor: menuInteractor,
	}
}

func (h *MenuHandler) Routes(r gin.IRouter) {
	r = r.Group("menu")
	r.GET("/:id", h.GetDetail())
	r.GET("", h.List())
}

func (h *MenuHandler) NoAuths() []string {
	return []string{}
}

// GetDetail
//
//	@Tags		菜单管理
//	@Security	BearerAuth
//	@Summary	获取菜单详情
//	@Param		id	path		string		true	"菜单ID"
//	@Success	200	{object}	domain.Menu	"成功"
//	@Router		/menu/{id} [get]
func (h *MenuHandler) GetDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MenuHandler.GetDetail")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取菜单ID
		idStr := c.Param("id")
		menuID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		menu, err := h.MenuInteractor.GetDetail(ctx, menuID, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to get menu detail: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, menu)
	}
}

// List
//
//	@Tags		菜单管理
//	@Security	BearerAuth
//	@Summary	查询菜单列表
//	@Param		name	query		string					false	"菜单名称（模糊匹配）"
//	@Param		page	query		int						false	"页码"
//	@Param		size	query		int						false	"每页数量"
//	@Success	200		{object}	domain.MenuSearchRes	"成功"
//	@Router		/menu [get]
func (h *MenuHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MenuHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.MenuListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		page := upagination.New(req.Page, req.Size)
		user := domain.FromStoreUserContext(ctx)

		params := domain.MenuSearchParams{
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
			Name:       req.Name,
		}

		res, err := h.MenuInteractor.PagedListBySearch(ctx, page, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to list menus: %w", err)
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}
