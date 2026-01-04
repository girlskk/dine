package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.OrderInteractor = (*OrderInteractor)(nil)

type OrderInteractor struct {
	DS domain.DataStore
}

func NewOrderInteractor(ds domain.DataStore) *OrderInteractor {
	return &OrderInteractor{
		DS: ds,
	}
}

func (interactor *OrderInteractor) Create(ctx context.Context, order *domain.Order) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "OrderInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 创建订单和商品需要在同一事务内
	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		return ds.OrderRepo().Create(ctx, order)
	})
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
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

func (interactor *OrderInteractor) Update(ctx context.Context, order *domain.Order) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "OrderInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 更新订单和商品需要在同一事务内
	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		return ds.OrderRepo().Update(ctx, order)
	})
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
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
