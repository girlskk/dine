package handler

import (
	"errors"
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
)

type StallHandler struct {
	StallInteractor domain.StallInteractor
}

func NewStallHandler(stallInteractor domain.StallInteractor) *StallHandler {
	return &StallHandler{StallInteractor: stallInteractor}
}

func (h *StallHandler) Routes(r gin.IRouter) {
	r = r.Group("/restaurant/stall")
	r.GET("/:id", h.Get())
	r.GET("", h.List())
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

		user := domain.FromFrontendUserContext(ctx)
		stall, err := h.StallInteractor.GetStall(ctx, id, user)
		if err != nil {
			if errors.Is(err, domain.ErrStallNotExists) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.StallNotExists, err))
				return
			}

			err = fmt.Errorf("failed to get stall: %w", err)
			c.Error(err)
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

		user := domain.FromFrontendUserContext(ctx)

		filter := &domain.StallListFilter{
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
			StallType:  domain.StallTypeStore,
			PrintType:  req.PrintType,
			Enabled:    req.Enabled,
		}

		pager := upagination.New(1, upagination.MaxSize)
		stalls, total, err := h.StallInteractor.GetStalls(ctx, pager, filter)
		if err != nil {
			err = fmt.Errorf("failed to get stalls: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.StallListResp{Stalls: stalls, Total: total})
	}
}
