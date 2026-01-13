package refundorder

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RefundOrderInteractor = (*RefundOrderInteractor)(nil)

type RefundOrderInteractor struct {
	DS domain.DataStore
}

func NewRefundOrderInteractor(ds domain.DataStore) *RefundOrderInteractor {
	return &RefundOrderInteractor{DS: ds}
}

func (uc *RefundOrderInteractor) Create(ctx context.Context, refundOrder *domain.RefundOrder) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RefundOrderInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()

	// 校验原订单存在且已支付
	originOrder, err := uc.DS.OrderRepo().FindByID(ctx, refundOrder.OriginOrderID)
	if err != nil {
		return fmt.Errorf("failed to find origin order: %w", err)
	}
	if originOrder.PaymentStatus != domain.PaymentStatusPaid {
		return domain.ParamsError(fmt.Errorf("origin order is not paid"))
	}

	// 填充原订单信息
	refundOrder.OriginOrderNo = originOrder.OrderNo
	refundOrder.OriginPaidAt = originOrder.PaidAt
	refundOrder.OriginAmountPaid = originOrder.Amount.AmountPaid

	err = uc.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		return ds.RefundOrderRepo().Create(ctx, refundOrder)
	})
	if err != nil {
		return fmt.Errorf("failed to create refund order: %w", err)
	}

	return nil
}

func (uc *RefundOrderInteractor) Get(ctx context.Context, id uuid.UUID) (res *domain.RefundOrder, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RefundOrderInteractor.Get")
	defer func() { util.SpanErrFinish(span, err) }()

	res, err = uc.DS.RefundOrderRepo().FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find refund order: %w", err)
	}
	return res, nil
}

func (uc *RefundOrderInteractor) Update(ctx context.Context, refundOrder *domain.RefundOrder) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RefundOrderInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()

	err = uc.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		return ds.RefundOrderRepo().Update(ctx, refundOrder)
	})
	if err != nil {
		return fmt.Errorf("failed to update refund order: %w", err)
	}
	return nil
}

func (uc *RefundOrderInteractor) Cancel(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RefundOrderInteractor.Cancel")
	defer func() { util.SpanErrFinish(span, err) }()

	// 查询退款单
	refundOrder, err := uc.DS.RefundOrderRepo().FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find refund order: %w", err)
	}

	// 只有待处理状态可以取消
	if refundOrder.RefundStatus != domain.RefundStatusPending {
		return domain.ParamsError(fmt.Errorf("refund order status is not pending"))
	}

	refundOrder.RefundStatus = domain.RefundStatusCancelled
	err = uc.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		return ds.RefundOrderRepo().Update(ctx, refundOrder)
	})
	if err != nil {
		return fmt.Errorf("failed to cancel refund order: %w", err)
	}
	return nil
}

func (uc *RefundOrderInteractor) List(ctx context.Context, params domain.RefundOrderListParams) (res []*domain.RefundOrder, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RefundOrderInteractor.List")
	defer func() { util.SpanErrFinish(span, err) }()

	res, total, err = uc.DS.RefundOrderRepo().List(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list refund orders: %w", err)
	}
	return res, total, nil
}
