package order

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.OrderInteractor = (*OrderInteractor)(nil)

type OrderInteractor struct {
	DS  domain.DataStore
	Seq domain.DailySequence
}

func NewOrderInteractor(ds domain.DataStore, seq domain.DailySequence) *OrderInteractor {
	return &OrderInteractor{
		DS:  ds,
		Seq: seq,
	}
}

func (interactor *OrderInteractor) Create(ctx context.Context, order *domain.Order) (res *domain.Order, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "OrderInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if order == nil {
		return nil, domain.ParamsErrorf("order is nil")
	}
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	if order.MerchantID == uuid.Nil {
		return nil, domain.ParamsErrorf("merchant_id is required")
	}
	if order.StoreID == uuid.Nil {
		return nil, domain.ParamsErrorf("store_id is required")
	}
	if order.BusinessDate == "" {
		return nil, domain.ParamsErrorf("business_date is required")
	}
	if order.DiningMode == "" {
		return nil, domain.ParamsErrorf("dining_mode is required")
	}

	if order.OrderNo == "" {
		order.OrderNo, err = interactor.generateOrderNo(ctx, order)
		if err != nil {
			return nil, fmt.Errorf("failed to generate order_no: %w", err)
		}
	}

	err = interactor.DS.OrderRepo().Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	res, err = interactor.DS.OrderRepo().FindByID(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find created order: %w", err)
	}
	return res, nil
}

func (interactor *OrderInteractor) Get(ctx context.Context, id uuid.UUID) (res *domain.Order, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "OrderInteractor.Get")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	res, err = interactor.DS.OrderRepo().FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	return res, nil
}

func (interactor *OrderInteractor) Update(ctx context.Context, order *domain.Order) (res *domain.Order, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "OrderInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if order == nil {
		return nil, domain.ParamsErrorf("order is nil")
	}
	if order.ID == uuid.Nil {
		return nil, domain.ParamsErrorf("id is required")
	}

	err = interactor.DS.OrderRepo().Update(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	res, err = interactor.DS.OrderRepo().FindByID(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find updated order: %w", err)
	}
	return res, nil
}

func (interactor *OrderInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "OrderInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = interactor.DS.OrderRepo().Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}

func (interactor *OrderInteractor) List(ctx context.Context, params domain.OrderListParams) (res []*domain.Order, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "OrderInteractor.List")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	res, total, err = interactor.DS.OrderRepo().List(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	return res, total, nil
}

func (interactor *OrderInteractor) generateOrderNo(ctx context.Context, o *domain.Order) (orderNo string, err error) {
	storePart := ""
	if o.Store != nil {
		if o.Store.StoreCode != "" {
			storePart = o.Store.StoreCode
		}
	}

	datePart := strings.ReplaceAll(o.BusinessDate, "-", "")
	prefix := fmt.Sprintf("%s:%s", domain.DailySequencePrefixOrderNo, o.StoreID.String())
	seq, err := interactor.Seq.Next(ctx, prefix)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s%06d", storePart, datePart, seq), nil
}
