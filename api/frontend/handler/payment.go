package handler

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/zxh"
)

type PaymentHandler struct {
	PayInteractor domain.PaymentInteractor
}

func NewPaymentHandler(interactor domain.PaymentInteractor) *PaymentHandler {
	return &PaymentHandler{
		PayInteractor: interactor,
	}
}

func (h *PaymentHandler) Routes(r gin.IRouter) {
	r = r.Group("/payment")
	r.POST("/polling", h.PayPolling())
	r.POST("/callback/pay/huifu", h.PayHuifuCallback())
	r.POST("/callback/pay/zxh", h.PayZxhCallback())
}

func (h *PaymentHandler) NoAuths() []string {
	return []string{
		"/payment/callback",
	}
}

// PayPolling 支付轮询
//
//	@Tags		支付管理
//	@Security	BearerAuth
//	@Summary	支付轮询
//	@Accept		json
//	@Produce	json
//	@Param		data	body		types.PayPollingReq	true	"请求参数"
//	@Success	200		{object}		types.PayPollingResp	"成功"
//	@Router		/payment/polling [post]
func (h *PaymentHandler) PayPolling() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentHandler.PayPolling")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("PaymentHandler.PayPolling")

		var req types.PayPollingReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		logger = logger.With("seq_no", req.SeqNo)
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromFrontendUserContext(ctx)

		state, reason, err := h.PayInteractor.PayPolling(ctx, req.SeqNo, user)
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to pay polling: %w", err))
			}
			return
		}

		response.Ok(c, &types.PayPollingResp{
			State:      state,
			FailReason: reason,
		})
	}
}

func (h *PaymentHandler) PayHuifuCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentHandler.PayHuifuCallback")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("PaymentHandler.PayHuifuCallback")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.PayHuifuCallback
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		if err := h.PayInteractor.PayHuifuCallback(ctx, req.Sign, req.RespData); err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
			} else {
				c.Error(fmt.Errorf("failed to pay huifu callback: %w", err))
			}
			return
		}

		response.Ok(c, nil)
	}
}

func (h *PaymentHandler) PayZxhCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentHandler.PayZxhCallback")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("PaymentHandler.PayZxhCallback")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req zxh.ZhixinhuaPointPayCallBack
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		// 立即返回200，异步处理业务逻辑
		response.Ok(c, nil)

		go func() {
			// 创建新的context，避免原请求的context被cancel
			asyncCtx := context.Background()
			asyncLogger := logging.FromContext(asyncCtx).Named("PaymentHandler.PayZxhCallback.async")
			// 复制必要的参数到新的context
			asyncCtx = logging.NewContext(asyncCtx, asyncLogger.With(
				"seqNo", req.OutOrderID,
				"request_id", c.GetString("request_id"), // 保留原请求ID用于追踪
			))
			asyncLogger.Infow("开始异步处理支付回调",
				"seqNo", req.OutOrderID,
				"status", req.Status)
			if err := h.PayInteractor.PayZxhCallback(asyncCtx, req); err != nil {
				asyncLogger.Errorf("异步处理支付回调失败: %v", err)
			} else {
				asyncLogger.Infow("异步处理支付回调成功", "seqNo", req.OutOrderID)
			}
		}()
	}
}
