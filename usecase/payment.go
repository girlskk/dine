package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/huifu"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/zxh"
)

var _ domain.PaymentInteractor = (*PaymentInteractor)(nil)

type PaymentInteractor struct {
	DataStore   domain.DataStore
	bsPay       *huifu.BsPay
	asynqClient *asynq.Client
}

func NewPaymentInteractor(dataStore domain.DataStore, bsPay *huifu.BsPay, asynqClient *asynq.Client) *PaymentInteractor {
	return &PaymentInteractor{
		DataStore:   dataStore,
		bsPay:       bsPay,
		asynqClient: asynqClient,
	}
}

func (interactor *PaymentInteractor) PayHuifuCallback(ctx context.Context, sign, respData string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentInteractor.PayHuifuCallback")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	logger := logging.FromContext(ctx).Named("PaymentInteractor.PayHuifuCallback")

	if _, err = huifu.RsaSignVerify(sign, respData, interactor.bsPay.Msc); err != nil {
		logger.Warnw("failed to verify sign", "error", err)
		err = domain.ParamsErrorf("sign verify failed")
		return
	}

	var respDataMap map[string]any
	if err = json.Unmarshal([]byte(respData), &respDataMap); err != nil {
		err = fmt.Errorf("failed to unmarshal response data: %w", err)
		return
	}

	seqNo, ok := respDataMap["req_seq_id"].(string)
	if !ok {
		err = errors.New("not found req_seq_id")
		return
	}

	provider := domain.PayProviderHuifu
	var callback *domain.PaymentCallback
	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) (err error) {
		callback, err = ds.PaymentRepo().CreateCallback(ctx, &domain.PaymentCallback{
			SeqNo:    seqNo,
			Type:     domain.PaymentCallbackTypePay,
			Raw:      json.RawMessage(respData),
			Provider: provider,
		})
		if err != nil {
			err = fmt.Errorf("failed to create payment callback: %w", err)
			return
		}

		payload := domain.PaymentCallbackTaskPayload{
			CallbackID: callback.ID,
			SeqNo:      seqNo,
			Provider:   provider,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			err = fmt.Errorf("failed to marshal payload: %w", err)
			return
		}

		t := asynq.NewTask(domain.TaskTypePaymentCallback, payloadBytes)
		if _, err = interactor.asynqClient.Enqueue(t, asynq.Timeout(time.Minute)); err != nil {
			err = fmt.Errorf("failed to enqueue payment callback task: %w", err)
			return
		}

		return
	})

	return
}

func (interactor *PaymentInteractor) PayZxhCallback(ctx context.Context, params zxh.ZhixinhuaPointPayCallBack) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentInteractor.PayZxhCallback")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	logger := logging.FromContext(ctx).Named("PaymentInteractor.PayZxhCallback")

	seqNo := params.OutOrderID

	// 添加重试机制
	var payment *domain.Payment
	maxRetries := 4
	retryDelay := 2 * time.Second
	for i := 0; i < maxRetries; i++ {
		payment, err = interactor.DataStore.PaymentRepo().FindBySeqNo(ctx, seqNo)
		if err == nil {
			// 找到记录，跳出循环
			break
		}

		if !domain.IsNotFound(err) {
			// 其他错误，直接返回
			return fmt.Errorf("failed to get payment: %w", err)
		}

		// 记录未找到，记录日志并等待
		logger.Infow("payment record not found, retrying...",
			"seqNo", seqNo,
			"attempt", i+1,
			"maxRetries", maxRetries)

		// 最后一次尝试不等待
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	// 重试后仍未找到
	if err != nil && domain.IsNotFound(err) {
		logger.Warnw("not found payment", "seqNo", seqNo)
		err = nil
		return
	}

	if payment.FinishedAt != nil {
		logger.Warnw("payment already finished", "seqNo", seqNo)
		err = nil
		return
	}

	store, err := interactor.DataStore.StoreRepo().Find(ctx, payment.StoreID)
	if err != nil {
		err = fmt.Errorf("failed to get store: %w", err)
		return
	}

	if store.ZxhID != params.MchId {
		err = errors.New("mch_id not match")
		return
	}

	if !params.Verify(store.ZxhSecret) {
		err = errors.New("verify failed")
		return
	}

	raw, err := json.Marshal(params)
	if err != nil {
		err = fmt.Errorf("failed to marshal params: %w", err)
		return
	}

	provider := lo.Ternary(payment.Channel == domain.PayChannelPointWallet, domain.PayProviderZhiXinHuaWallet, domain.PayProviderZhiXinHua)

	var callback *domain.PaymentCallback
	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) (err error) {
		callback, err = ds.PaymentRepo().CreateCallback(ctx, &domain.PaymentCallback{
			SeqNo:    seqNo,
			Type:     domain.PaymentCallbackTypePay,
			Raw:      raw,
			Provider: provider,
		})
		if err != nil {
			err = fmt.Errorf("failed to create payment callback: %w", err)
			return
		}

		payload := domain.PaymentCallbackTaskPayload{
			CallbackID: callback.ID,
			SeqNo:      seqNo,
			Provider:   provider,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			err = fmt.Errorf("failed to marshal payload: %w", err)
			return
		}

		t := asynq.NewTask(domain.TaskTypePaymentCallback, payloadBytes)
		if _, err = interactor.asynqClient.Enqueue(t, asynq.Timeout(time.Minute)); err != nil {
			err = fmt.Errorf("failed to enqueue payment callback task: %w", err)
			return
		}

		return
	})

	return
}

func (interactor *PaymentInteractor) PayPolling(ctx context.Context, seqNo string, operator *domain.FrontendUser) (state domain.PayState, reason string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentInteractor.PayPolling")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	payment, err := interactor.DataStore.PaymentRepo().FindBySeqNo(ctx, seqNo)
	if err != nil {
		err = fmt.Errorf("failed to get payment: %w", err)
		return
	}

	if operator.StoreID != payment.StoreID {
		err = domain.ParamsErrorf("支付流水号不存在")
		return
	}

	state = payment.State

	if state == domain.PayStateWaiting || state == domain.PayStateFailure {
		reason = payment.FailReason
	}

	return
}
