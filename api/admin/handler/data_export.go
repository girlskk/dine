package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type DataExportHandler struct {
	DataExportInteractor domain.DataExportInteractor
}

func NewDataExportHandler(interactor domain.DataExportInteractor) *DataExportHandler {
	return &DataExportHandler{
		DataExportInteractor: interactor,
	}
}

func (h *DataExportHandler) Routes(r gin.IRouter) {
	r = r.Group("/data-export")
	r.POST("/list", h.List())
	r.POST("/retry", h.Retry())
}

// List	数据导出列表
//
//	@Tags		数据导出
//	@Security	BearerAuth
//	@Summary	数据导出列表
//	@Accept		json
//	@Produce	json
//	@Param		data	body		types.DataExportListReq		true	"请求参数"
//	@Success	200		{object}	types.DataExportListResp	"成功"
//	@Router		/data-export/list [post]
func (h *DataExportHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "DataExportHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("DataExportHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.DataExportListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		dataExports, total, err := h.DataExportInteractor.List(ctx, req.ToPagination(), &domain.DataExportFilter{
			StoreID:      0,
			Type:         req.Type,
			Status:       req.Status,
			CreatedAtGte: req.CreatedAtStart.ToPtrStartOfDay(),
			CreatedAtLte: req.CreatedAtEnd.ToPtrEndOfDay(),
		})
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to get data exports: %w", err))
			}
			return
		}

		response.Ok(c, &types.DataExportListResp{
			DataExports: dataExports,
			Total:       total,
		})
	}
}

// Retry 重试
//
//	@Tags		数据导出
//	@Security	BearerAuth
//	@Summary	重试
//	@Accept		json
//	@Produce	json
//	@Param		data	body	types.DataExportRetryReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/data-export/retry [post]
func (h *DataExportHandler) Retry() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "DataExportHandler.Retry")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("DataExportHandler.Retry")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.DataExportRetryReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.DataExportInteractor.Retry(ctx, 0, req.ID); err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to retry data export: %w", err))
			}
			return
		}

		response.Ok(c, nil)
	}
}
