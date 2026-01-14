package handler

import (
	"errors"
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
)

type StallHandler struct {
	StallInteractor domain.StallInteractor
}

func NewStallHandler(stallInteractor domain.StallInteractor) *StallHandler {
	return &StallHandler{StallInteractor: stallInteractor}
}

func (h *StallHandler) Routes(r gin.IRouter) {
	r = r.Group("/restaurant/stall")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("/:id", h.Get())
	r.GET("", h.List())
	r.PUT("/:id/enable", h.Enable())
	r.PUT("/:id/disable", h.Disable())
}

// Create 创建出品部门
//
//	@Tags			后厨管理
//	@Security		BearerAuth
//	@Summary		创建出品部门
//	@Description	创建出品部门
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.StallCreateReq	true	"请求信息"
//	@Success		200		"No Content"
//	@Router			/restaurant/stall [post]
func (h *StallHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StallHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StallCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		stall := &domain.Stall{
			Name:       req.Name,
			StallType:  domain.StallTypeStore,
			PrintType:  req.PrintType,
			Enabled:    req.Enabled,
			SortOrder:  req.SortOrder,
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
		}

		if err := h.StallInteractor.Create(ctx, stall, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新出品部门
//
//	@Tags			后厨管理
//	@Security		BearerAuth
//	@Summary		更新出品部门
//	@Description	更新出品部门
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"出品部门ID"
//	@Param			data	body	types.StallUpdateReq	true	"请求信息"
//	@Success		200		"No Content"
//	@Router			/restaurant/stall/{id} [put]
func (h *StallHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StallHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.StallUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		stall := &domain.Stall{
			ID:        id,
			Name:      req.Name,
			PrintType: req.PrintType,
			Enabled:   req.Enabled,
			SortOrder: req.SortOrder,
		}

		if err := h.StallInteractor.Update(ctx, stall, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Delete 删除出品部门
//
//	@Tags			后厨管理
//	@Security		BearerAuth
//	@Summary		删除出品部门
//	@Description	删除出品部门
//	@Param			id	path	string	true	"出品部门ID"
//	@Success		200	"No Content"
//	@Success		204	"No Content"
//	@Router			/restaurant/stall/{id} [delete]
func (h *StallHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StallHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		if err := h.StallInteractor.Delete(ctx, id, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Get 获取出品部门详情
//
//	@Tags			后厨管理
//	@Security		BearerAuth
//	@Summary		获取出品部门详情
//	@Description	根据出品部门ID获取详情
//	@Param			id	path		string	true	"出品部门ID"
//	@Success		200	{object}	response.Response{data=domain.Stall}
//	@Router			/restaurant/stall/{id} [get]
func (h *StallHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StallHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		stall, err := h.StallInteractor.GetStall(ctx, id, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, stall)
	}
}

// List 获取出品部门列表
//
//	@Tags			后厨管理
//	@Security		BearerAuth
//	@Summary		获取出品部门列表
//	@Description	分页查询出品部门列表
//	@Param			data	query		types.StallListReq	true	"出品部门列表查询参数"
//	@Success		200		{object}	response.Response{data=types.StallListResp}
//	@Router			/restaurant/stall [get]
func (h *StallHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StallHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StallListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		pager := req.RequestPagination.ToPagination()
		filter := &domain.StallListFilter{
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
			StallType:  domain.StallTypeStore,
			PrintType:  req.PrintType,
			Enabled:    req.Enabled,
			Name:       req.Name,
		}

		stalls, total, err := h.StallInteractor.GetStalls(ctx, pager, filter)
		if err != nil {
			err = fmt.Errorf("failed to get stalls: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.StallListResp{Stalls: stalls, Total: total})
	}
}

// Enable 启用出品部门
//
//	@Tags			后厨管理
//	@Security		BearerAuth
//	@Summary		启用出品部门
//	@Description	将出品部门置为启用
//	@Produce		json
//	@Param			id	path	string	true	"出品部门ID"
//	@Success		200	"No Content"
//	@Router			/restaurant/stall/{id}/enable [put]
func (h *StallHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StallHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		stall := &domain.Stall{ID: id, Enabled: true}

		err = h.StallInteractor.StallSimpleUpdate(
			ctx,
			domain.StallSimpleUpdateFieldEnabled,
			stall,
			user,
		)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 禁用出品部门
//
//	@Tags			后厨管理
//	@Security		BearerAuth
//	@Summary		禁用出品部门
//	@Description	将出品部门置为禁用
//	@Produce		json
//	@Param			id	path	string	true	"出品部门ID"
//	@Success		200	"No Content"
//	@Router			/restaurant/stall/{id}/disable [put]
func (h *StallHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StallHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		stall := &domain.Stall{ID: id, Enabled: false}
		err = h.StallInteractor.StallSimpleUpdate(
			ctx,
			domain.StallSimpleUpdateFieldEnabled,
			stall,
			user,
		)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

func (h *StallHandler) checkErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrStallNotExists):
		return errorx.New(http.StatusBadRequest, errcode.StallNotExists, err)
	case errors.Is(err, domain.ErrStallNameExists):
		return errorx.New(http.StatusConflict, errcode.StallNameExists, err)
	case domain.IsNotFound(err):
		return errorx.New(http.StatusNotFound, errcode.NotFound, err)
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return fmt.Errorf("stall handler error: %w", err)
	}
}
