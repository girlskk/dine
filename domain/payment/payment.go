package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/domain/order"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

type DomainService struct {
	MutexManager        domain.MutexManager
	DailySequence       domain.DailySequence
	OrderDomainService  *order.DomainService
	DataStore           domain.DataStore
	PaymentEventTrigger domain.PaymentEventTrigger
}

func NewDomainService(
	mutexManager domain.MutexManager,
	dailySequence domain.DailySequence,
	dataStore domain.DataStore,
	paymentEventTrigger domain.PaymentEventTrigger,
) *DomainService {
	return &DomainService{
		MutexManager:        mutexManager,
		DailySequence:       dailySequence,
		DataStore:           dataStore,
		PaymentEventTrigger: paymentEventTrigger,
	}
}

type PaymentParams struct {
	AuthCode   string
	NotifyURL  string
	Amount     decimal.Decimal
	GoodsDesc  string
	IPAddr     string
	PayBizType domain.PayBizType
	BizID      int
	Creator    any
	StoreID    int
	extra      map[string]any
	once       sync.Once
}

func (p *PaymentParams) Set(key string, value any) {
	p.once.Do(func() {
		p.extra = make(map[string]any)
	})
	p.extra[key] = value
}

func (p *PaymentParams) Get(key string) any {
	if v, ok := p.extra[key]; ok {
		return v
	}

	return nil
}

type PaymentResult struct {
	Req     json.RawMessage
	Resp    json.RawMessage
	Channel domain.PayChannel
}

type PaymentProvider interface {
	Provider() domain.PayProvider
	MchID() string
	Payment(ctx context.Context, seqNo string, params *PaymentParams) (*PaymentResult, error)
}

type PaymentCallbackParams struct {
	CallbackID int
	SeqNo      string
}

type PaymentCallbackResult struct {
	State      domain.PayState
	FailReason string
	MemberInfo domain.PaymentMemberInfo
}

type PaymentCallbackProvider interface {
	Callback(ctx context.Context, callback *domain.PaymentCallback, payment *domain.Payment) (*PaymentCallbackResult, error)
}

func (s *DomainService) genSeqNo(ctx context.Context) (seqNo string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentDomainService.genSeqNo")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	no, err := s.DailySequence.Next(ctx, domain.DailySequencePrefixPayNo)
	if err != nil {
		err = fmt.Errorf("failed to generate seq no: %w", err)
		return
	}
	currentTime := time.Now().Format("20060102150405")
	seqNo = fmt.Sprintf("%s%06d", currentTime, no)
	return
}

func (s *DomainService) ProcessPayment(ctx context.Context, provider PaymentProvider, params *PaymentParams) (seqNo string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentDomainService.ProcessPayment")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	seqNo, err = s.genSeqNo(ctx)
	if err != nil {
		return
	}
	logger := logging.FromContext(ctx).Named("PaymentDomainService.ProcessPayment")
	logger = logger.With("seqNo", seqNo)
	ctx = logging.NewContext(ctx, logger)

	mu := s.MutexManager.NewMutex(domain.NewMutexPaymentKey(seqNo))
	if err = mu.Lock(ctx); err != nil {
		err = fmt.Errorf("failed to lock payment: %w", err)
		return
	}
	defer func() {
		if _, err := mu.Unlock(ctx); err != nil {
			logger.Errorf("failed to unlock payment: %s", err)
		}
	}()

	res, err := provider.Payment(ctx, seqNo, params)
	if err != nil {
		return "", err
	}

	user := domain.ExtractOperatorInfo(params.Creator)

	payment := &domain.Payment{
		SeqNo:       seqNo,
		Provider:    provider.Provider(),
		Channel:     res.Channel,
		State:       domain.PayStateProcessing,
		MchID:       provider.MchID(),
		Amount:      params.Amount,
		GoodsDesc:   params.GoodsDesc,
		IPAddr:      params.IPAddr,
		Req:         res.Req,
		Resp:        res.Resp,
		PayBizType:  params.PayBizType,
		BizID:       params.BizID,
		CreatorType: user.Type,
		CreatorID:   user.ID,
		CreatorName: user.Name,
		StoreID:     params.StoreID,
	}

	if _, err = s.DataStore.PaymentRepo().Create(ctx, payment); err != nil {
		err = fmt.Errorf("failed to create payment: %w", err)
		return
	}

	return
}

func (s *DomainService) ProcessPaymentCallback(ctx context.Context, provider PaymentCallbackProvider, params *PaymentCallbackParams) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentDomainService.ProcessPaymentCallback")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	logger := logging.FromContext(ctx).Named("PaymentDomainService.ProcessPaymentCallback")
	logger = logger.With("seqNo", params.SeqNo)
	ctx = logging.NewContext(ctx, logger)

	mu := s.MutexManager.NewMutex(domain.NewMutexPaymentKey(params.SeqNo))
	if err = mu.Lock(ctx); err != nil {
		err = fmt.Errorf("failed to lock payment: %w", err)
		return
	}
	defer func() {
		if _, err := mu.Unlock(ctx); err != nil {
			logger.Errorf("failed to unlock payment: %s", err)
		}
	}()

	err = s.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) (err error) {
		callback, err := s.DataStore.PaymentRepo().GetCallback(ctx, params.CallbackID)
		if err != nil {
			if domain.IsNotFound(err) {
				logger.Warnf("callback %d not found", params.CallbackID)
				err = nil
				return
			}
			err = fmt.Errorf("failed to get callback: %w", err)
			return
		}

		if err = ds.PaymentRepo().RemoveCallback(ctx, callback.ID); err != nil {
			err = fmt.Errorf("failed to remove callback: %w", err)
			return
		}

		span.LogKV("event", "remove callback")

		payment, err := s.DataStore.PaymentRepo().FindBySeqNo(ctx, params.SeqNo)
		if err != nil {
			if domain.IsNotFound(err) { // 发起阶段就失败了，没有入库
				logger.Warnf("payment %s not found", params.SeqNo)
				err = nil
				return
			}
			err = fmt.Errorf("failed to get payment: %w", err)
			return
		}
		if payment.FinishedAt != nil {
			logger.Infof("payment %s already finished", params.SeqNo)
			return
		}

		res, err := provider.Callback(ctx, callback, payment)
		if err != nil {
			err = fmt.Errorf("failed to callback: %w", err)
			return
		}

		span.LogKV("event", "provider callback")

		payment.State = res.State
		payment.FailReason = res.FailReason
		payment.Callback = callback.Raw
		if payment.State.IsFinished() {
			payment.FinishedAt = lo.ToPtr(time.Now())
		}

		payment, err = ds.PaymentRepo().Update(ctx, payment)
		if err != nil {
			err = fmt.Errorf("failed to update payment: %w", err)
			return
		}

		span.LogKV("event", "update payment")

		if payment.State != domain.PayStateSuccess { // 未成功不走后续流程
			return nil
		}

		var store *domain.Store
		store, err = ds.StoreRepo().Find(ctx, payment.StoreID)
		if err != nil {
			err = fmt.Errorf("failed to get store: %w", err)
			return
		}

		operator := &domain.FrontendUser{
			ID:       payment.CreatorID,
			Nickname: payment.CreatorName,
			StoreID:  store.ID,
			Store:    store,
		}

		// 触发支付成功事件
		err = s.PaymentEventTrigger.FireSuccess(ctx, &domain.PaymentEventSuccessParams{
			PaymentEventBaseParams: domain.PaymentEventBaseParams{
				DataStore: ds,
				Operator:  operator,
				Payment:   payment,
			},
			MemberInfo: res.MemberInfo,
		})

		span.LogKV("event", "fire success")

		return
	})

	return
}
