package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type ReconciliationHandler struct {
	ReconciliationRecordInteractor domain.ReconciliationRecordInteractor
	DataExportInteractor           domain.DataExportInteractor
}

func NewReconciliationHandler(interactor domain.ReconciliationRecordInteractor, dataExportInteractor domain.DataExportInteractor) *ReconciliationHandler {
	return &ReconciliationHandler{
		ReconciliationRecordInteractor: interactor,
		DataExportInteractor:           dataExportInteractor,
	}
}

// Routes 注册路由
func (h *ReconciliationHandler) Routes(r gin.IRouter) {
	r = r.Group("/reconciliation")
	r.POST("/list", h.List())
	r.POST("/list/export", h.ExportList())
	r.POST("/list-details", h.ListDetails())
	r.POST("/list-details/export", h.ExportListDetails())
	r.POST("/summary", h.Summary())
}

// List 财务对账单列表
//
//	@Tags		财务对账
//	@Summary	财务对账单列表
//	@Security	BearerAuth
//	@Param		data	body		types.ReconciliationListReq	true	"请求参数"
//	@Success	200		{object}	domain.ReconciliationSearchRes
//	@Router		/reconciliation/list [post]
func (h *ReconciliationHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ReconciliationHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ReconciliationListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		user := domain.FromBackendUserContext(ctx)
		page := upagination.New(req.Page, req.Size)
		params := domain.ReconciliationSearchParams{
			StartAt: req.StartAt.ToPtrStartOfDay(),
			EndAt:   req.EndAt.ToPtrEndOfDay(),
			StoreID: user.StoreID,
			Channel: req.Channel,
		}

		res, err := h.ReconciliationRecordInteractor.PagedListBySearch(ctx, page, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}

// ListDetails 财务对账单明细
//
//	@Tags	财务对账
//	@Summary		财务对账单明细
//	@Security	BearerAuth
//	@Param		data	body	types.ReconciliationDetailReq	true	"请求参数"
//	@Success		200		{object}	domain.ReconciliationDetails
//	@Router		/reconciliation/list-details [post]
func (h *ReconciliationHandler) ListDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationHandler.ListDetails")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ReconciliationHandler.ListDetails")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ReconciliationDetailReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		user := domain.FromBackendUserContext(ctx)
		res, err := h.ReconciliationRecordInteractor.ListDetails(ctx, req.ID, user.StoreID)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}
		response.Ok(c, res)
	}
}

// Summary 财务对账汇总
//
//	@Tags	财务对账
//	@Summary		财务对账汇总
//	@Security	BearerAuth
//	@Param		data	body	types.ReconciliationSummaryReq	true	"请求参数"
//	@Success		200		{object}	domain.ReconciliationSummaryRes
//	@Router		/reconciliation/summary [post]
func (h *ReconciliationHandler) Summary() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationHandler.Summary")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ReconciliationHandler.Summary")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ReconciliationSummaryReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		user := domain.FromBackendUserContext(ctx)
		params := domain.ReconciliationSearchParams{
			StartAt: req.StartAt.ToPtrStartOfDay(),
			EndAt:   req.EndAt.ToPtrEndOfDay(),
			StoreID: user.StoreID,
			Channel: req.Channel,
		}

		res, err := h.ReconciliationRecordInteractor.Summary(ctx, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}
		response.Ok(c, res)
	}
}

// ExportList	导出对账单列表
//
//	@Tags		财务对账
//	@Security	BearerAuth
//	@Summary	导出对账单列表
//	@Param		data	body	types.ReconciliationListExportReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/reconciliation/list/export [post]
func (h *ReconciliationHandler) ExportList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationHandler.ExportList")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ReconciliationHandler.ExportList")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ReconciliationListExportReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if !req.StartAt.IsValid() || !req.EndAt.IsValid() {
			c.Error(uerr.BadRequest("时间范围不能为空"))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		filter := domain.ReconciliationSearchParams{
			StartAt: req.StartAt.ToPtrStartOfDay(),
			EndAt:   req.EndAt.ToPtrEndOfDay(),
			StoreID: user.StoreID,
			Channel: req.Channel,
		}

		reconciliationRange, err := h.ReconciliationRecordInteractor.GetReconciliationRange(ctx, filter)
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to get reconciliation range: %w", err))
			}
			return
		}

		if reconciliationRange.Count < 1 {
			c.Error(uerr.BadRequest("没有财务对账单可导出"))
			return
		}

		filter.IDGte = reconciliationRange.MinID
		filter.IDLte = reconciliationRange.MaxID

		totalPages := upagination.TotalPages(reconciliationRange.Count, domain.ReconciliationRecordListExportSingleMaxSize)
		params := make([]*domain.ReconciliationRecordListExportParams, 0, totalPages)
		for i := range totalPages {
			page := i + 1
			params = append(params, &domain.ReconciliationRecordListExportParams{
				Filter: filter,
				Pager:  *upagination.New(page, domain.ReconciliationRecordListExportSingleMaxSize),
			})
		}

		fileName := fmt.Sprintf(
			"%s-%s_%d_财务对账单列表.xlsx",
			req.StartAt.ToTime().Format(time.DateOnly),
			req.EndAt.ToTime().Format(time.DateOnly),
			time.Now().Unix(),
		)

		submitParams, err := domain.BuildDataExportSubmitParams(user.StoreID, domain.DataExportTypeReconciliationRecordListExport, params, fileName, user)
		if err != nil {
			c.Error(fmt.Errorf("failed to build data export submit params: %w", err))
			return
		}

		if _, err = h.DataExportInteractor.Submit(ctx, submitParams...); err != nil {
			c.Error(fmt.Errorf("failed to submit data export: %w", err))
			return
		}

		response.Ok(c, nil)
	}
}

// ExportListDetails	导出对账单明细
//
//	@Tags		财务对账
//	@Security	BearerAuth
//	@Summary	导出对账单明细
//	@Param		data	body	types.ReconciliationListExportReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/reconciliation/list-details/export [post]
func (h *ReconciliationHandler) ExportListDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "ReconciliationHandler.ExportListDetails")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("ReconciliationHandler.ExportListDetails")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ReconciliationListExportReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if !req.StartAt.IsValid() || !req.EndAt.IsValid() {
			c.Error(uerr.BadRequest("时间范围不能为空"))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		filter := domain.ReconciliationSearchParams{
			StartAt: req.StartAt.ToPtrStartOfDay(),
			EndAt:   req.EndAt.ToPtrEndOfDay(),
			StoreID: user.StoreID,
			Channel: req.Channel,
		}

		reconciliationRange, err := h.ReconciliationRecordInteractor.GetReconciliationRange(ctx, filter)
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to get reconciliation range: %w", err))
			}
			return
		}

		if reconciliationRange.Count < 1 {
			c.Error(uerr.BadRequest("没有财务对账单可导出"))
			return
		}

		filter.IDGte = reconciliationRange.MinID
		filter.IDLte = reconciliationRange.MaxID

		totalPages := upagination.TotalPages(reconciliationRange.Count, domain.ReconciliationRecordDetailExportSingleMaxSize)
		params := make([]*domain.ReconciliationRecordListExportParams, 0, totalPages)
		for i := range totalPages {
			page := i + 1
			params = append(params, &domain.ReconciliationRecordListExportParams{
				Filter: filter,
				Pager:  *upagination.New(page, domain.ReconciliationRecordDetailExportSingleMaxSize),
			})
		}

		fileName := fmt.Sprintf(
			"%s-%s_%d_财务对账单明细.xlsx",
			req.StartAt.ToTime().Format(time.DateOnly),
			req.EndAt.ToTime().Format(time.DateOnly),
			time.Now().Unix(),
		)

		submitParams, err := domain.BuildDataExportSubmitParams(user.StoreID, domain.DataExportTypeReconciliationRecordDetailsExport, params, fileName, user)
		if err != nil {
			c.Error(fmt.Errorf("failed to build data export submit params: %w", err))
			return
		}

		if _, err = h.DataExportInteractor.Submit(ctx, submitParams...); err != nil {
			c.Error(fmt.Errorf("failed to submit data export: %w", err))
			return
		}

		response.Ok(c, nil)
	}
}
