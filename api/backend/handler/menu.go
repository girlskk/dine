package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
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
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("/:id", h.GetDetail())
	r.GET("", h.List())
}

func (h *MenuHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		菜单管理
//	@Security	BearerAuth
//	@Summary	创建菜单
//	@Param		data	body	types.MenuCreateReq	true	"请求信息"
//	@Success	200
//	@Router		/menu [post]
func (h *MenuHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MenuHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.MenuCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		req.StoreIDs = lo.Uniq(req.StoreIDs)

		// 构建 domain.Menu
		menu := &domain.Menu{
			ID:               uuid.New(),
			MerchantID:       user.MerchantID,
			Name:             req.Name,
			DistributionRule: req.DistributionRule,
			Stores: lo.Map(req.StoreIDs, func(storeID uuid.UUID, _ int) *domain.StoreSimple {
				return &domain.StoreSimple{
					ID: storeID,
				}
			}),
			StoreCount: len(req.StoreIDs),
			ItemCount:  len(req.Items),
		}

		// 转换菜单项
		menu.Items = make(domain.MenuItems, 0, len(req.Items))
		productIDMap := make(map[uuid.UUID]struct{})
		for _, itemReq := range req.Items {
			if _, ok := productIDMap[itemReq.ProductID]; ok {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, errors.New("菜单商品重复")))
				return
			}
			productIDMap[itemReq.ProductID] = struct{}{}
			menu.Items = append(menu.Items, &domain.MenuItem{
				ID:          uuid.New(),
				ProductID:   itemReq.ProductID,
				SaleRule:    itemReq.SaleRule,
				BasePrice:   itemReq.BasePrice,
				MemberPrice: itemReq.MemberPrice,
			})
		}

		err := h.MenuInteractor.Create(ctx, menu)
		if err != nil {
			if errors.Is(err, domain.ErrMenuNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create menu: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, menu)
	}
}

// Update
//
//	@Tags		菜单管理
//	@Security	BearerAuth
//	@Summary	更新菜单
//	@Param		id		path	string				true	"菜单ID"
//	@Param		data	body	types.MenuUpdateReq	true	"请求信息"
//	@Success	200
//	@Router		/menu/{id} [put]
func (h *MenuHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MenuHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取菜单ID
		idStr := c.Param("id")
		menuID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.MenuUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		req.StoreIDs = lo.Uniq(req.StoreIDs)

		// 构建 domain.Menu
		menu := &domain.Menu{
			ID:               menuID,
			MerchantID:       user.MerchantID,
			Name:             req.Name,
			DistributionRule: req.DistributionRule,
			Stores: lo.Map(req.StoreIDs, func(storeID uuid.UUID, _ int) *domain.StoreSimple {
				return &domain.StoreSimple{
					ID: storeID,
				}
			}),
			StoreCount: len(req.StoreIDs),
			ItemCount:  len(req.Items),
		}

		// 转换菜单项
		menu.Items = make(domain.MenuItems, 0, len(req.Items))
		productIDMap := make(map[uuid.UUID]struct{})
		for _, itemReq := range req.Items {
			if _, ok := productIDMap[itemReq.ProductID]; ok {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, errors.New("菜单商品重复")))
				return
			}
			productIDMap[itemReq.ProductID] = struct{}{}
			menu.Items = append(menu.Items, &domain.MenuItem{
				ID:          uuid.New(),
				ProductID:   itemReq.ProductID,
				SaleRule:    itemReq.SaleRule,
				BasePrice:   itemReq.BasePrice,
				MemberPrice: itemReq.MemberPrice,
			})
		}

		err = h.MenuInteractor.Update(ctx, menu)
		if err != nil {
			if errors.Is(err, domain.ErrMenuNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update menu: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, menu)
	}
}

// Delete
//
//	@Tags		菜单管理
//	@Security	BearerAuth
//	@Summary	删除菜单
//	@Param		id	path	string	true	"菜单ID"
//	@Success	200
//	@Router		/menu/{id} [delete]
func (h *MenuHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MenuHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取菜单ID
		idStr := c.Param("id")
		menuID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		err = h.MenuInteractor.Delete(ctx, menuID)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to delete menu: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
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

		menu, err := h.MenuInteractor.GetDetail(ctx, menuID)
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
		user := domain.FromBackendUserContext(ctx)

		params := domain.MenuSearchParams{
			MerchantID: user.MerchantID,
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
