package repository

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/payment"
	"gitlab.jiguang.dev/pos-dine/dine/ent/paymentcallback"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.PaymentRepository = (*PaymentRepository)(nil)

type PaymentRepository struct {
	Client *ent.Client
}

func NewPaymentRepository(client *ent.Client) *PaymentRepository {
	return &PaymentRepository{
		Client: client,
	}
}

func (r *PaymentRepository) Create(ctx context.Context, dpayment *domain.Payment) (newPayment *domain.Payment, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	pmt, err := r.Client.Payment.Create().
		SetSeqNo(dpayment.SeqNo).
		SetProvider(dpayment.Provider).
		SetChannel(dpayment.Channel).
		SetState(dpayment.State).
		SetAmount(dpayment.Amount).
		SetGoodsDesc(dpayment.GoodsDesc).
		SetIPAddr(dpayment.IPAddr).
		SetPayBizType(dpayment.PayBizType).
		SetMchID(dpayment.MchID).
		SetReq(dpayment.Req).
		SetResp(dpayment.Resp).
		SetCallback(dpayment.Callback).
		SetRefunded(dpayment.Refunded).
		SetBizID(dpayment.BizID).
		SetCreatorType(dpayment.CreatorType).
		SetCreatorID(dpayment.CreatorID).
		SetCreatorName(dpayment.CreatorName).
		SetStoreID(dpayment.StoreID).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create payment: %w", err)
		return
	}

	newPayment = convertPayment(pmt)

	return
}

func (r *PaymentRepository) Update(ctx context.Context, dpayment *domain.Payment) (newPayment *domain.Payment, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	pmt, err := r.Client.Payment.UpdateOneID(dpayment.ID).
		SetState(dpayment.State).
		SetRefunded(dpayment.Refunded).
		SetFailReason(dpayment.FailReason).
		SetCallback(dpayment.Callback).
		SetNillableFinishedAt(dpayment.FinishedAt).Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to update payment: %w", err)
		return
	}
	newPayment = convertPayment(pmt)

	return
}

func (r *PaymentRepository) FindBySeqNo(ctx context.Context, seqNo string) (dpayment *domain.Payment, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentRepository.FindBySeqNo")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	pmt, err := r.Client.Payment.Query().
		Where(payment.SeqNo(seqNo)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		} else {
			err = fmt.Errorf("failed to find payment by seq no: %w", err)
		}
		return
	}

	dpayment = convertPayment(pmt)

	return
}

func (r *PaymentRepository) CreateCallback(ctx context.Context, callback *domain.PaymentCallback) (newCallback *domain.PaymentCallback, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentRepository.CreateCallback")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	cb, err := r.Client.PaymentCallback.Create().
		SetSeqNo(callback.SeqNo).
		SetType(callback.Type).
		SetRaw(callback.Raw).
		SetProvider(callback.Provider).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create payment callback: %w", err)
		return
	}

	newCallback = convertPaymentCallback(cb)

	return
}

func (r *PaymentRepository) GetCallback(ctx context.Context, id int) (callback *domain.PaymentCallback, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentRepository.GetCallback")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	cb, err := r.Client.PaymentCallback.Query().
		Where(paymentcallback.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to get payment callback by id: %w", err)
		return
	}

	callback = convertPaymentCallback(cb)

	return
}

func (r *PaymentRepository) RemoveCallback(ctx context.Context, callbackID int) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PaymentRepository.RemoveCallback")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = r.Client.PaymentCallback.DeleteOneID(callbackID).Exec(ctx)
	if err != nil {
		err = fmt.Errorf("failed to remove payment callback: %w", err)
		return
	}

	return
}

func convertPayment(pmt *ent.Payment) *domain.Payment {
	return &domain.Payment{
		ID:          pmt.ID,
		SeqNo:       pmt.SeqNo,
		Provider:    pmt.Provider,
		Channel:     pmt.Channel,
		State:       pmt.State,
		Amount:      pmt.Amount,
		GoodsDesc:   pmt.GoodsDesc,
		MchID:       pmt.MchID,
		IPAddr:      pmt.IPAddr,
		Req:         pmt.Req,
		Resp:        pmt.Resp,
		Callback:    pmt.Callback,
		FinishedAt:  pmt.FinishedAt,
		Refunded:    pmt.Refunded,
		FailReason:  pmt.FailReason,
		PayBizType:  pmt.PayBizType,
		BizID:       pmt.BizID,
		CreatorType: pmt.CreatorType,
		CreatorID:   pmt.CreatorID,
		CreatorName: pmt.CreatorName,
		StoreID:     pmt.StoreID,
		CreatedAt:   pmt.CreatedAt,
		UpdatedAt:   pmt.UpdatedAt,
	}
}

func convertPaymentCallback(cb *ent.PaymentCallback) *domain.PaymentCallback {
	return &domain.PaymentCallback{
		ID:       cb.ID,
		Type:     cb.Type,
		Raw:      cb.Raw,
		Provider: cb.Provider,
	}
}
