package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type PointSettlementHandler struct {
	PointSettlementInteractor domain.PointSettlementInteractor
	DataExportInteractor      domain.DataExportInteractor
}

func NewPointSettlementHandler(interactor domain.PointSettlementInteractor, dataExportInteractor domain.DataExportInteractor) *PointSettlementHandler {
	return &PointSettlementHandler{
		PointSettlementInteractor: interactor,
		DataExportInteractor:      dataExportInteractor,
	}
}

// Routes 注册路由
func (h *PointSettlementHandler) Routes(r gin.IRouter) {
	r = r.Group("/point-settlement")
	r.POST("/list", h.List())
	r.POST("/list/export", h.ExportList())
	r.POST("/list-details", h.ListDetails())
	r.POST("/list-details/export", h.ExportListDetails())
	r.POST("/approve", h.Approve())
	r.POST("/unapprove", h.UnApprove())
}

// List 积分结算账单列表
//
//	@Tags		积分结算
//	@Summary	积分结算账单列表
//	@Security	BearerAuth
//	@Param		data	body		types.PointSettlementListReq	true	"请求参数"
//	@Success	200		{object}	domain.PointSettlementSearchRes
//	@Router		/point-settlement/list [post]
func (h *PointSettlementHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("PointSettlementHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.PointSettlementListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(req.Page, req.Size)
		params := domain.PointSettlementSearchParams{
			StartAt: req.StartAt.ToPtrStartOfDay(),
			EndAt:   req.EndAt.ToPtrEndOfDay(),
			StoreID: req.StoreID,
		}

		res, err := h.PointSettlementInteractor.PagedListBySearch(ctx, page, params)
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

// ListDetails 积分结算单明细
//
//	@Tags	积分结算
//	@Summary		积分结算单明细
//	@Security	BearerAuth
//	@Param		data	body	types.PointSettlementIDReq	true	"请求参数"
//	@Success		200		{object}	domain.PointSettlementDetails
//	@Router		/point-settlement/list-details [post]
func (h *PointSettlementHandler) ListDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementHandler.ListDetails")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("PointSettlementHandler.ListDetails")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.PointSettlementIDReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}
		res, err := h.PointSettlementInteractor.ListDetails(ctx, req.ID, 0)
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

// Approve 积分结算账单审批
//
//	@Tags		积分结算
//	@Summary	积分结算账单审批
//	@Security	BearerAuth
//	@Param		data	body	types.PointSettlementIDReq	true	"请求参数"
//	@Success	200
//	@Router		/point-settlement/approve [post]
func (h *PointSettlementHandler) Approve() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementHandler.Approve")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("PointSettlementHandler.Approve")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.PointSettlementIDReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.PointSettlementInteractor.Approve(ctx, req.ID); err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, nil)
	}
}

// UnApprove 积分结算账单反审批
//
//	@Tags		积分结算
//	@Summary	积分结算账单反审批
//	@Security	BearerAuth
//	@Param		data	body	types.PointSettlementIDReq	true	"请求参数"
//	@Success	200
//	@Router		/point-settlement/unapprove [post]
func (h *PointSettlementHandler) UnApprove() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementHandler.UnApprove")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("PointSettlementHandler.UnApprove")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.PointSettlementIDReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.PointSettlementInteractor.UnApprove(ctx, req.ID); err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, nil)
	}
}

// ExportList	导出积分结算单列表
//
//	@Tags		积分结算
//	@Security	BearerAuth
//	@Summary	导出积分结算单列表
//	@Param		data	body	types.PointSettlementListExportReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/point-settlement/list/export [post]
func (h *PointSettlementHandler) ExportList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementHandler.ExportList")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("PointSettlementHandler.ExportList")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.PointSettlementListExportReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if !req.StartAt.IsValid() || !req.EndAt.IsValid() {
			c.Error(uerr.BadRequest("时间范围不能为空"))
			return
		}

		filter := domain.PointSettlementSearchParams{
			StartAt: req.StartAt.ToPtrStartOfDay(),
			EndAt:   req.EndAt.ToPtrEndOfDay(),
			StoreID: req.StoreID,
		}

		pointSettlementRange, err := h.PointSettlementInteractor.GetPointSettlementRange(ctx, filter)
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to get point settlement range: %w", err))
			}
			return
		}

		if pointSettlementRange.Count < 1 {
			c.Error(uerr.BadRequest("没有积分结算单可导出"))
			return
		}

		filter.IDGte = pointSettlementRange.MinID
		filter.IDLte = pointSettlementRange.MaxID

		totalPages := upagination.TotalPages(pointSettlementRange.Count, domain.PointSettlementListExportSingleMaxSize)
		params := make([]*domain.PointSettlementListExportParams, 0, totalPages)
		for i := range totalPages {
			page := i + 1
			params = append(params, &domain.PointSettlementListExportParams{
				Filter: filter,
				Pager:  *upagination.New(page, domain.PointSettlementListExportSingleMaxSize),
			})
		}

		user := domain.FromAdminUserContext(ctx)

		fileName := fmt.Sprintf(
			"%s-%s_%d_积分结算单列表.xlsx",
			req.StartAt.ToTime().Format(time.DateOnly),
			req.EndAt.ToTime().Format(time.DateOnly),
			time.Now().Unix(),
		)

		submitParams, err := domain.BuildDataExportSubmitParams(0, domain.DataExportTypePointSettlementListExport, params, fileName, user)
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

// ExportListDetails	导出积分结算单明细
//
//	@Tags		积分结算
//	@Security	BearerAuth
//	@Summary	导出积分结算单明细
//	@Param		data	body	types.PointSettlementListExportReq	true	"请求参数"
//	@Success	200		"No Content"
//	@Router		/point-settlement/list-details/export [post]
func (h *PointSettlementHandler) ExportListDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementHandler.ExportListDetails")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("PointSettlementHandler.ExportListDetails")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.PointSettlementListExportReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if !req.StartAt.IsValid() || !req.EndAt.IsValid() {
			c.Error(uerr.BadRequest("时间范围不能为空"))
			return
		}

		filter := domain.PointSettlementSearchParams{
			StartAt: req.StartAt.ToPtrStartOfDay(),
			EndAt:   req.EndAt.ToPtrEndOfDay(),
			StoreID: req.StoreID,
		}

		pointSettlementRange, err := h.PointSettlementInteractor.GetPointSettlementRange(ctx, filter)
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to get point settlement range: %w", err))
			}
			return
		}

		if pointSettlementRange.Count < 1 {
			c.Error(uerr.BadRequest("没有积分结算单可导出"))
			return
		}

		filter.IDGte = pointSettlementRange.MinID
		filter.IDLte = pointSettlementRange.MaxID

		totalPages := upagination.TotalPages(pointSettlementRange.Count, domain.PointSettlementDetailExportSingleMaxSize)
		params := make([]*domain.PointSettlementListExportParams, 0, totalPages)
		for i := range totalPages {
			page := i + 1
			params = append(params, &domain.PointSettlementListExportParams{
				Filter: filter,
				Pager:  *upagination.New(page, domain.PointSettlementDetailExportSingleMaxSize),
			})
		}

		user := domain.FromAdminUserContext(ctx)

		fileName := fmt.Sprintf(
			"%s-%s_%d_积分结算单明细.xlsx",
			req.StartAt.ToTime().Format(time.DateOnly),
			req.EndAt.ToTime().Format(time.DateOnly),
			time.Now().Unix(),
		)

		submitParams, err := domain.BuildDataExportSubmitParams(0, domain.DataExportTypePointSettlementDetailsExport, params, fileName, user)
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
