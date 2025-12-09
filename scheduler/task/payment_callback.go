package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/domain/payment"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
	"gitlab.jiguang.dev/pos-dine/dine/scheduler"
)

var _ scheduler.Handler = (*PaymentCallbackTask)(nil)

type PaymentCallbackTask struct {
	DataStore            domain.DataStore
	PaymentDomainService *payment.DomainService
	Alert                alert.Alert
}

func NewPaymentCallbackTask(dataStore domain.DataStore, paymentDomainService *payment.DomainService, alert alert.Alert) *PaymentCallbackTask {
	return &PaymentCallbackTask{
		DataStore:            dataStore,
		PaymentDomainService: paymentDomainService,
		Alert:                alert,
	}
}

func (task *PaymentCallbackTask) Type() string {
	return domain.TaskTypePaymentCallback
}

func (task *PaymentCallbackTask) ProcessTask(ctx context.Context, t *asynq.Task) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentCallbackTask.ProcessTask")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	logger := logging.FromContext(ctx).Named("PaymentCallbackTask.ProcessTask")
	ctx = logging.NewContext(ctx, logger)

	var payload *domain.PaymentCallbackTaskPayload
	if err = json.Unmarshal(t.Payload(), &payload); err != nil {
		err = fmt.Errorf("json.Unmarshal payload: %w", err)
		return
	}

	defer func() {
		if err != nil {
			logger.Errorf("支付回调失败[%d]: %v", payload.CallbackID, err)
			task.Alert.Notify(ctx, fmt.Sprintf("支付回调失败[%d]: %v", payload.CallbackID, err))
		}
	}()

	logger.Infof("处理支付回调[%d]", payload.CallbackID)
	var provider payment.PaymentCallbackProvider

	switch payload.Provider {
	case domain.PayProviderZhiXinHua:
		provider = payment.NewZxhPaymentCallbackProvider()
	case domain.PayProviderZhiXinHuaWallet:
		provider = payment.NewZxhWalletPaymentCallbackProvider()
	case domain.PayProviderHuifu:
		provider = payment.NewHuifuPaymentCallbackProvider()
	default:
		err = fmt.Errorf("unknown provider: %s", payload.Provider)
		return
	}

	if err = task.PaymentDomainService.ProcessPaymentCallback(ctx, provider, &payment.PaymentCallbackParams{
		CallbackID: payload.CallbackID,
		SeqNo:      payload.SeqNo,
	}); err != nil {
		err = fmt.Errorf("PaymentDomainService.ProcessPaymentCallback: %w", err)
		return
	}

	return
}
