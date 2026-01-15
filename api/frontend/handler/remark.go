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

type RemarkHandler struct {
	RemarkInteractor domain.RemarkInteractor
}

func NewRemarkHandler(remarkInteractor domain.RemarkInteractor) *RemarkHandler {
	return &RemarkHandler{
		RemarkInteractor: remarkInteractor,
	}
}

func (h *RemarkHandler) Routes(r gin.IRouter) {
	r = r.Group("/remark")
	r.GET("/:id", h.Get())
	r.GET("", h.List())
}

// Get 获取备注详情
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		获取备注详情
//	@Description	根据备注ID获取详情
//	@Param			id	path		string	true	"备注ID"
//	@Success		200	{object}	response.Response{data=domain.Remark}
//	@Router			/remark/{id} [get]
func (h *RemarkHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		user := domain.FromFrontendUserContext(ctx)
		remark, err := h.RemarkInteractor.GetRemark(ctx, id, user)
		if err != nil {
			if errors.Is(err, domain.ErrRemarkNotExists) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.RemarkNotExists, err))
				return
			}
			err = fmt.Errorf("failed to get remark: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, remark)
	}
}

// List 获取备注列表
//
//	@Tags			前厅管理
//	@Security		BearerAuth
//	@Summary		获取备注列表
//	@Description	分页查询备注列表
//	@Param			data	query		types.RemarkListReq	true	"备注列表查询参数"
//	@Success		200		{object}	response.Response{data=types.RemarkListResp}
//	@Router			/remark [get]
func (h *RemarkHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RemarkHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RemarkListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromFrontendUserContext(ctx)

		filter := &domain.RemarkListFilter{
			MerchantID:  user.MerchantID,
			StoreID:     user.StoreID,
			Enabled:     req.Enabled,
			RemarkType:  domain.RemarkTypeStore,
			RemarkScene: req.RemarkScene,
		}

		pager := upagination.New(1, upagination.MaxSize)
		remarks, total, err := h.RemarkInteractor.GetRemarks(ctx, pager, filter)
		if err != nil {
			err = fmt.Errorf("failed to get remarks: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.RemarkListResp{
			Remarks: remarks,
			Total:   total,
		})
	}
}
