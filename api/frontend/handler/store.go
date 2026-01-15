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

type StoreHandler struct {
	StoreInteractor domain.StoreInteractor
}

func NewStoreHandler(storeInteractor domain.StoreInteractor) *StoreHandler {
	return &StoreHandler{StoreInteractor: storeInteractor}
}

func (h *StoreHandler) Routes(r gin.IRouter) {
	r = r.Group("store")
	r.GET("/:id", h.GetStore())
	r.GET("/list", h.List())
}

// GetStore 获取门店
//
//	@Summary		获取门店
//	@Description	根据门店ID获取门店信息
//	@Tags			门店管理
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"门店ID"
//	@Success		200	{object}	response.Response{data=domain.Store}
//	@Router			/store/{id} [get]
func (h *StoreHandler) GetStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.GetStore")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		storeID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromFrontendUserContext(ctx)
		user.StoreID = storeID
		domainStore, err := h.StoreInteractor.GetStore(ctx, storeID, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, domainStore)
	}
}

// List 门店列表
//
//	@Summary		门店列表
//	@Description	分页查询门店列表
//	@Tags			门店管理
//	@Security		BearerAuth
//	@Produce		json
//	@Param			data	query		types.StoreListReq	true	"门店列表查询参数"
//	@Success		200		{object}	response.Response{data=types.StoreListResp}
//	@Router			/store/list [get]
func (h *StoreHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromFrontendUserContext(ctx)

		filter := &domain.StoreListFilter{
			MerchantID:       user.MerchantID,
			BusinessTypeCode: req.BusinessTypeCode,
			Status:           req.Status,
			BusinessModel:    req.BusinessModel,
		}

		pager := upagination.New(1, upagination.MaxSize)
		stores, total, err := h.StoreInteractor.GetStores(ctx, pager, filter)
		if err != nil {
			c.Error(fmt.Errorf("failed to list stores: %w", err))
			return
		}

		response.Ok(c, &types.StoreListResp{Stores: stores, Total: total})
	}
}

func (h *StoreHandler) checkErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrStoreNotExists):
		return errorx.New(http.StatusBadRequest, errcode.StoreNotExists, err)
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return fmt.Errorf("failed to process store: %w", err)
	}
}
