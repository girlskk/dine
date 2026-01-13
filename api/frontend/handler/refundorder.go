package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type RefundOrderHandler struct {
	RefundOrderInteractor domain.RefundOrderInteractor
	Seq                   domain.DailySequence
}

func NewRefundOrderHandler(refundOrderInteractor domain.RefundOrderInteractor, seq domain.DailySequence) *RefundOrderHandler {
	return &RefundOrderHandler{
		RefundOrderInteractor: refundOrderInteractor,
		Seq:                   seq,
	}
}

func (h *RefundOrderHandler) Routes(r gin.IRouter) {
	r = r.Group("/refund-order")
	r.POST("", h.Create())
	r.GET("/:id", h.Get())
	r.PUT("/:id", h.Update())
	r.POST("/:id/cancel", h.Cancel())
	r.GET("", h.List())
}

func (h *RefundOrderHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		退款订单
//	@Security	BearerAuth
//	@Summary	创建退款订单
//	@Accept		json
//	@Produce	json
//	@Param		data	body		types.CreateRefundOrderReq	true	"请求信息"
//	@Success	200		{object}	domain.RefundOrder			"成功"
//	@Router		/refund-order [post]
func (h *RefundOrderHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RefundOrderHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CreateRefundOrderReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromFrontendUserContext(ctx)

		// 构建退款商品明细
		refundProducts := make([]domain.RefundOrderProduct, len(req.RefundProducts))
		for i, rp := range req.RefundProducts {
			refundProducts[i] = domain.RefundOrderProduct{
				OriginOrderProductID: rp.OriginOrderProductID,
				OriginOrderItemID:    rp.OriginOrderItemID,
				ProductID:            rp.ProductID,
				ProductName:          rp.ProductName,
				ProductType:          rp.ProductType,
				Category:             rp.Category,
				MainImage:            rp.MainImage,
				Description:          rp.Description,
				OriginQty:            rp.OriginQty,
				OriginPrice:          rp.OriginPrice,
				OriginSubtotal:       rp.OriginSubtotal,
				OriginDiscount:       rp.OriginDiscount,
				OriginTax:            rp.OriginTax,
				OriginTotal:          rp.OriginTotal,
				RefundQty:            rp.RefundQty,
				RefundSubtotal:       rp.RefundSubtotal,
				RefundDiscount:       rp.RefundDiscount,
				RefundTax:            rp.RefundTax,
				RefundTotal:          rp.RefundTotal,
				Groups:               rp.Groups,
				SpecRelations:        rp.SpecRelations,
				AttrRelations:        rp.AttrRelations,
				RefundReason:         rp.RefundReason,
			}
		}

		ro := &domain.RefundOrder{
			ID:               uuid.New(),
			MerchantID:       user.MerchantID,
			StoreID:          user.StoreID,
			BusinessDate:     req.BusinessDate,
			ShiftNo:          req.ShiftNo,
			RefundNo:         req.RefundNo,
			OriginOrderID:    req.OriginOrderID,
			RefundType:       req.RefundType,
			RefundStatus:     domain.RefundStatusPending,
			RefundReasonCode: req.RefundReasonCode,
			RefundReason:     req.RefundReason,
			Store:            req.Store,
			Channel:          domain.ChannelPOS,
			Pos:              req.Pos,
			Cashier:          req.Cashier,
			RefundAmount:     req.RefundAmount,
			RefundPayments:   req.RefundPayments,
			RefundProducts:   refundProducts,
			Remark:           req.Remark,
		}

		// 自动生成退款单号
		if ro.RefundNo == "" {
			refundNo, err := h.generateRefundNo(ctx, ro)
			if err != nil {
				c.Error(fmt.Errorf("failed to generate refund_no: %w", err))
				return
			}
			ro.RefundNo = refundNo
		}

		err := h.RefundOrderInteractor.Create(ctx, ro)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsConflict(err) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			c.Error(fmt.Errorf("failed to create refund order: %w", err))
			return
		}

		response.Ok(c, ro)
	}
}

func (h *RefundOrderHandler) generateRefundNo(ctx context.Context, ro *domain.RefundOrder) (string, error) {
	storePart := ""
	if ro.Store.StoreCode != "" {
		storePart = ro.Store.StoreCode
	}

	datePart := strings.ReplaceAll(ro.BusinessDate, "-", "")
	prefix := fmt.Sprintf("%s:%s", "RF", ro.StoreID.String())
	seq, err := h.Seq.Next(ctx, prefix)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s%06d", storePart, datePart, seq), nil
}

// Get
//
//	@Tags		退款订单
//	@Security	BearerAuth
//	@Summary	获取退款订单详情
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string				true	"退款订单ID"
//	@Success	200	{object}	domain.RefundOrder	"成功"
//	@Router		/refund-order/{id} [get]
func (h *RefundOrderHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RefundOrderHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		ro, err := h.RefundOrderInteractor.Get(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			c.Error(fmt.Errorf("failed to get refund order: %w", err))
			return
		}

		response.Ok(c, ro)
	}
}

// Update
//
//	@Tags		退款订单
//	@Security	BearerAuth
//	@Summary	更新退款订单
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string						true	"退款订单ID"
//	@Param		data	body		types.UpdateRefundOrderReq	true	"请求信息"
//	@Success	200		{object}	domain.RefundOrder			"成功"
//	@Router		/refund-order/{id} [put]
func (h *RefundOrderHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RefundOrderHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.UpdateRefundOrderReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		ro := &domain.RefundOrder{
			ID:               id,
			RefundStatus:     req.RefundStatus,
			RefundReasonCode: req.RefundReasonCode,
			RefundReason:     req.RefundReason,
			ApprovedBy:       req.ApprovedBy,
			ApprovedByName:   req.ApprovedByName,
			ApprovedAt:       req.ApprovedAt,
			RefundedAt:       req.RefundedAt,
			RefundPayments:   req.RefundPayments,
			Remark:           req.Remark,
		}

		err = h.RefundOrderInteractor.Update(ctx, ro)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsConflict(err) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			c.Error(fmt.Errorf("failed to update refund order: %w", err))
			return
		}

		response.Ok(c, ro)
	}
}

// Cancel
//
//	@Tags		退款订单
//	@Security	BearerAuth
//	@Summary	取消退款订单
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string				true	"退款订单ID"
//	@Success	200	{object}	response.Response	"成功"
//	@Router		/refund-order/{id}/cancel [post]
func (h *RefundOrderHandler) Cancel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RefundOrderHandler.Cancel")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		err = h.RefundOrderInteractor.Cancel(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			c.Error(fmt.Errorf("failed to cancel refund order: %w", err))
			return
		}

		response.Ok(c, nil)
	}
}

// List
//
//	@Tags		退款订单
//	@Security	BearerAuth
//	@Summary	退款订单列表
//	@Accept		json
//	@Produce	json
//	@Param		origin_order_id	query		string						false	"原订单ID"
//	@Param		business_date	query		string						false	"营业日"
//	@Param		refund_no		query		string						false	"退款单号"
//	@Param		refund_type		query		string						false	"退款类型"
//	@Param		refund_status	query		string						false	"退款状态"
//	@Param		page			query		int							false	"页码"
//	@Param		size			query		int							false	"每页数量"
//	@Success	200				{object}	types.RefundOrderListResp	"成功"
//	@Router		/refund-order [get]
func (h *RefundOrderHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RefundOrderHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.RefundOrderListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromFrontendUserContext(ctx)
		params := domain.RefundOrderListParams{
			MerchantID:    user.MerchantID,
			StoreID:       user.StoreID,
			OriginOrderID: req.OriginOrderID,
			RefundNo:      req.RefundNo,
			Page:          req.Page,
			Size:          req.Size,
		}
		if req.RefundType != "" {
			params.RefundType = domain.RefundType(req.RefundType)
		}
		if req.RefundStatus != "" {
			params.RefundStatus = domain.RefundStatus(req.RefundStatus)
		}

		items, total, err := h.RefundOrderInteractor.List(ctx, params)
		if err != nil {
			c.Error(fmt.Errorf("failed to list refund orders: %w", err))
			return
		}

		response.Ok(c, types.RefundOrderListResp{
			Items: items,
			Total: total,
		})
	}
}
