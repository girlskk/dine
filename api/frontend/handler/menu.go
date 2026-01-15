package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
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
	r.GET("", h.ListAll())
}

func (h *MenuHandler) NoAuths() []string {
	return []string{}
}

// ListAll
//
//	@Tags		菜单管理
//	@Security	BearerAuth
//	@Summary	查询所有菜单
//	@Param		store_id	query		string					false	"门店ID"
//	@Success	200			{object}	domain.MenuSearchRes	"成功"
//	@Router		/menu [get]
func (h *MenuHandler) ListAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MenuHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		storeIDStr := c.Query("store_id")
		if storeIDStr == "" {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("store_id is required")))
			return
		}
		storeID, err := uuid.Parse(storeIDStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("invalid store_id: %w", err)))
			return
		}

		user := domain.FromFrontendContext(ctx)
		params := domain.MenuListAllParams{
			MerchantID: user.MerchantID,
			StoreID:    storeID,
		}

		res, err := h.MenuInteractor.ListAllStoreMenus(ctx, params)
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
