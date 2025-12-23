package repository

import (
	"context"
	"fmt"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	entorder "gitlab.jiguang.dev/pos-dine/dine/ent/order"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.OrderRepository = (*OrderRepository)(nil)

type OrderRepository struct {
	Client *ent.Client
}

func NewOrderRepository(client *ent.Client) *OrderRepository {
	return &OrderRepository{Client: client}
}

func (repo *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.Order, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "OrderRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eo, err := repo.Client.Order.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return convertOrderToDomain(eo), nil
}

func (repo *OrderRepository) Create(ctx context.Context, o *domain.Order) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "OrderRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.Order.Create().
		SetMerchantID(o.MerchantID).
		SetStoreID(o.StoreID).
		SetBusinessDate(o.BusinessDate).
		SetOrderNo(o.OrderNo).
		SetDiningMode(entorder.DiningMode(o.DiningMode))

	if o.ID != uuid.Nil {
		builder = builder.SetID(o.ID)
	}
	if o.ShiftNo != "" {
		builder = builder.SetShiftNo(o.ShiftNo)
	}
	if o.OrderType != "" {
		builder = builder.SetOrderType(entorder.OrderType(o.OrderType))
	}
	if o.OriginOrderID != "" {
		builder = builder.SetOriginOrderID(o.OriginOrderID)
	}
	if len(o.Refund) > 0 {
		builder = builder.SetRefund(o.Refund)
	}

	if o.OpenedAt != nil {
		builder = builder.SetOpenedAt(*o.OpenedAt)
	}
	if o.PlacedAt != nil {
		builder = builder.SetPlacedAt(*o.PlacedAt)
	}
	if o.PaidAt != nil {
		builder = builder.SetPaidAt(*o.PaidAt)
	}
	if o.CompletedAt != nil {
		builder = builder.SetCompletedAt(*o.CompletedAt)
	}

	if o.OpenedBy != "" {
		builder = builder.SetOpenedBy(o.OpenedBy)
	}
	if o.PlacedBy != "" {
		builder = builder.SetPlacedBy(o.PlacedBy)
	}
	if o.PaidBy != "" {
		builder = builder.SetPaidBy(o.PaidBy)
	}

	if o.OrderStatus != "" {
		builder = builder.SetOrderStatus(entorder.OrderStatus(o.OrderStatus))
	}
	if o.PaymentStatus != "" {
		builder = builder.SetPaymentStatus(entorder.PaymentStatus(o.PaymentStatus))
	}
	if o.FulfillmentStatus != "" {
		builder = builder.SetFulfillmentStatus(entorder.FulfillmentStatus(o.FulfillmentStatus))
	}
	if o.TableStatus != "" {
		builder = builder.SetTableStatus(entorder.TableStatus(o.TableStatus))
	}

	if o.TableID != "" {
		builder = builder.SetTableID(o.TableID)
	}
	if o.TableName != "" {
		builder = builder.SetTableName(o.TableName)
	}
	if o.TableCapacity != 0 {
		builder = builder.SetTableCapacity(o.TableCapacity)
	}
	if o.GuestCount != 0 {
		builder = builder.SetGuestCount(o.GuestCount)
	}

	if o.MergedToOrderID != "" {
		builder = builder.SetMergedToOrderID(o.MergedToOrderID)
	}
	if o.MergedAt != nil {
		builder = builder.SetMergedAt(*o.MergedAt)
	}

	if len(o.Store) > 0 {
		builder = builder.SetStore(o.Store)
	}
	if len(o.Channel) > 0 {
		builder = builder.SetChannel(o.Channel)
	}
	if len(o.Pos) > 0 {
		builder = builder.SetPos(o.Pos)
	}
	if len(o.Cashier) > 0 {
		builder = builder.SetCashier(o.Cashier)
	}

	if len(o.Member) > 0 {
		builder = builder.SetMember(o.Member)
	}
	if len(o.Takeaway) > 0 {
		builder = builder.SetTakeaway(o.Takeaway)
	}

	if len(o.Cart) > 0 {
		builder = builder.SetCart(o.Cart)
	}
	if len(o.Products) > 0 {
		builder = builder.SetProducts(o.Products)
	}
	if len(o.Promotions) > 0 {
		builder = builder.SetPromotions(o.Promotions)
	}
	if len(o.Coupons) > 0 {
		builder = builder.SetCoupons(o.Coupons)
	}
	if len(o.TaxRates) > 0 {
		builder = builder.SetTaxRates(o.TaxRates)
	}
	if len(o.Fees) > 0 {
		builder = builder.SetFees(o.Fees)
	}
	if len(o.Payments) > 0 {
		builder = builder.SetPayments(o.Payments)
	}
	if len(o.RefundsProducts) > 0 {
		builder = builder.SetRefundsProducts(o.RefundsProducts)
	}
	if len(o.Amount) > 0 {
		builder = builder.SetAmount(o.Amount)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		if ent.IsValidationError(err) {
			return domain.ParamsError(fmt.Errorf("invalid order params: %w", err))
		}
		if ent.IsConstraintError(err) {
			return domain.ConflictError(err)
		}
		return fmt.Errorf("failed to create order: %w", err)
	}

	o.ID = created.ID
	o.CreatedAt = created.CreatedAt
	o.UpdatedAt = created.UpdatedAt

	return nil
}

func (repo *OrderRepository) Update(ctx context.Context, o *domain.Order) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "OrderRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.Order.UpdateOneID(o.ID)

	if o.BusinessDate != "" {
		builder = builder.SetBusinessDate(o.BusinessDate)
	}
	if o.ShiftNo != "" {
		builder = builder.SetShiftNo(o.ShiftNo)
	}
	if o.OrderNo != "" {
		builder = builder.SetOrderNo(o.OrderNo)
	}

	if o.OrderType != "" {
		builder = builder.SetOrderType(entorder.OrderType(o.OrderType))
	}
	if o.OriginOrderID != "" {
		builder = builder.SetOriginOrderID(o.OriginOrderID)
	}
	if len(o.Refund) > 0 {
		builder = builder.SetRefund(o.Refund)
	}

	if o.OpenedAt != nil {
		builder = builder.SetOpenedAt(*o.OpenedAt)
	}
	if o.PlacedAt != nil {
		builder = builder.SetPlacedAt(*o.PlacedAt)
	}
	if o.PaidAt != nil {
		builder = builder.SetPaidAt(*o.PaidAt)
	}
	if o.CompletedAt != nil {
		builder = builder.SetCompletedAt(*o.CompletedAt)
	}

	if o.OpenedBy != "" {
		builder = builder.SetOpenedBy(o.OpenedBy)
	}
	if o.PlacedBy != "" {
		builder = builder.SetPlacedBy(o.PlacedBy)
	}
	if o.PaidBy != "" {
		builder = builder.SetPaidBy(o.PaidBy)
	}

	if o.DiningMode != "" {
		builder = builder.SetDiningMode(entorder.DiningMode(o.DiningMode))
	}
	if o.OrderStatus != "" {
		builder = builder.SetOrderStatus(entorder.OrderStatus(o.OrderStatus))
	}
	if o.PaymentStatus != "" {
		builder = builder.SetPaymentStatus(entorder.PaymentStatus(o.PaymentStatus))
	}
	if o.FulfillmentStatus != "" {
		builder = builder.SetFulfillmentStatus(entorder.FulfillmentStatus(o.FulfillmentStatus))
	}
	if o.TableStatus != "" {
		builder = builder.SetTableStatus(entorder.TableStatus(o.TableStatus))
	}

	if o.TableID != "" {
		builder = builder.SetTableID(o.TableID)
	}
	if o.TableName != "" {
		builder = builder.SetTableName(o.TableName)
	}
	if o.TableCapacity != 0 {
		builder = builder.SetTableCapacity(o.TableCapacity)
	}
	if o.GuestCount != 0 {
		builder = builder.SetGuestCount(o.GuestCount)
	}

	if o.MergedToOrderID != "" {
		builder = builder.SetMergedToOrderID(o.MergedToOrderID)
	}
	if o.MergedAt != nil {
		builder = builder.SetMergedAt(*o.MergedAt)
	}

	if len(o.Store) > 0 {
		builder = builder.SetStore(o.Store)
	}
	if len(o.Channel) > 0 {
		builder = builder.SetChannel(o.Channel)
	}
	if len(o.Pos) > 0 {
		builder = builder.SetPos(o.Pos)
	}
	if len(o.Cashier) > 0 {
		builder = builder.SetCashier(o.Cashier)
	}
	if len(o.Member) > 0 {
		builder = builder.SetMember(o.Member)
	}
	if len(o.Takeaway) > 0 {
		builder = builder.SetTakeaway(o.Takeaway)
	}
	if len(o.Cart) > 0 {
		builder = builder.SetCart(o.Cart)
	}
	if len(o.Products) > 0 {
		builder = builder.SetProducts(o.Products)
	}
	if len(o.Promotions) > 0 {
		builder = builder.SetPromotions(o.Promotions)
	}
	if len(o.Coupons) > 0 {
		builder = builder.SetCoupons(o.Coupons)
	}
	if len(o.TaxRates) > 0 {
		builder = builder.SetTaxRates(o.TaxRates)
	}
	if len(o.Fees) > 0 {
		builder = builder.SetFees(o.Fees)
	}
	if len(o.Payments) > 0 {
		builder = builder.SetPayments(o.Payments)
	}
	if len(o.RefundsProducts) > 0 {
		builder = builder.SetRefundsProducts(o.RefundsProducts)
	}
	if len(o.Amount) > 0 {
		builder = builder.SetAmount(o.Amount)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.NotFoundError(err)
		}
		if ent.IsValidationError(err) {
			return domain.ParamsError(fmt.Errorf("invalid order params: %w", err))
		}
		if ent.IsConstraintError(err) {
			return domain.ConflictError(err)
		}
		return fmt.Errorf("failed to update order: %w", err)
	}

	o.UpdatedAt = updated.UpdatedAt
	return nil
}

func (repo *OrderRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "OrderRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.Order.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.NotFoundError(err)
		}
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}

func (repo *OrderRepository) List(ctx context.Context, params domain.OrderListParams) (res []*domain.Order, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "OrderRepository.List")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if params.OrderType != "" {
		if vErr := entorder.OrderTypeValidator(entorder.OrderType(params.OrderType)); vErr != nil {
			return nil, 0, domain.ParamsError(fmt.Errorf("invalid order_type: %w", vErr))
		}
	}
	if params.OrderStatus != "" {
		if vErr := entorder.OrderStatusValidator(entorder.OrderStatus(params.OrderStatus)); vErr != nil {
			return nil, 0, domain.ParamsError(fmt.Errorf("invalid order_status: %w", vErr))
		}
	}
	if params.PaymentStatus != "" {
		if vErr := entorder.PaymentStatusValidator(entorder.PaymentStatus(params.PaymentStatus)); vErr != nil {
			return nil, 0, domain.ParamsError(fmt.Errorf("invalid payment_status: %w", vErr))
		}
	}

	query := repo.Client.Order.Query()
	if params.MerchantID != "" {
		query.Where(entorder.MerchantID(params.MerchantID))
	}
	if params.StoreID != "" {
		query.Where(entorder.StoreID(params.StoreID))
	}
	if params.BusinessDate != "" {
		query.Where(entorder.BusinessDate(params.BusinessDate))
	}
	if params.OrderNo != "" {
		query.Where(entorder.OrderNo(params.OrderNo))
	}
	if params.OrderType != "" {
		query.Where(entorder.OrderTypeEQ(entorder.OrderType(params.OrderType)))
	}
	if params.OrderStatus != "" {
		query.Where(entorder.OrderStatusEQ(entorder.OrderStatus(params.OrderStatus)))
	}
	if params.PaymentStatus != "" {
		query.Where(entorder.PaymentStatusEQ(entorder.PaymentStatus(params.PaymentStatus)))
	}

	pageInfo := upagination.New(params.Page, params.Size)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	items, err := query.
		Order(entorder.ByCreatedAt(entsql.OrderDesc())).
		Limit(pageInfo.Size).
		Offset(pageInfo.Offset()).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}

	res = make([]*domain.Order, 0, len(items))
	for _, eo := range items {
		res = append(res, convertOrderToDomain(eo))
	}
	return res, total, nil
}

func convertOrderToDomain(eo *ent.Order) *domain.Order {
	if eo == nil {
		return nil
	}

	return &domain.Order{
		ID:        eo.ID,
		CreatedAt: eo.CreatedAt,
		UpdatedAt: eo.UpdatedAt,
		DeletedAt: eo.DeletedAt,

		MerchantID: eo.MerchantID,
		StoreID:    eo.StoreID,

		BusinessDate: eo.BusinessDate,
		ShiftNo:      eo.ShiftNo,
		OrderNo:      eo.OrderNo,

		OrderType:     eo.OrderType.String(),
		OriginOrderID: eo.OriginOrderID,
		Refund:        eo.Refund,

		OpenedAt:    eo.OpenedAt,
		PlacedAt:    eo.PlacedAt,
		PaidAt:      eo.PaidAt,
		CompletedAt: eo.CompletedAt,

		OpenedBy: eo.OpenedBy,
		PlacedBy: eo.PlacedBy,
		PaidBy:   eo.PaidBy,

		DiningMode:        eo.DiningMode.String(),
		OrderStatus:       eo.OrderStatus.String(),
		PaymentStatus:     eo.PaymentStatus.String(),
		FulfillmentStatus: eo.FulfillmentStatus.String(),
		TableStatus:       eo.TableStatus.String(),

		TableID:       eo.TableID,
		TableName:     eo.TableName,
		TableCapacity: eo.TableCapacity,
		GuestCount:    eo.GuestCount,

		MergedToOrderID: eo.MergedToOrderID,
		MergedAt:        eo.MergedAt,

		Store:   eo.Store,
		Channel: eo.Channel,
		Pos:     eo.Pos,
		Cashier: eo.Cashier,

		Member:   eo.Member,
		Takeaway: eo.Takeaway,

		Cart:            eo.Cart,
		Products:        eo.Products,
		Promotions:      eo.Promotions,
		Coupons:         eo.Coupons,
		TaxRates:        eo.TaxRates,
		Fees:            eo.Fees,
		Payments:        eo.Payments,
		RefundsProducts: eo.RefundsProducts,
		Amount:          eo.Amount,
	}
}
