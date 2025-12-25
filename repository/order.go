package repository

import (
	"context"
	"encoding/json"
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

	res, err = convertOrderToDomain(eo)
	if err != nil {
		return nil, fmt.Errorf("failed to convert order: %w", err)
	}
	return res, nil
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
		SetDiningMode(o.DiningMode)

	if o.ID != uuid.Nil {
		builder = builder.SetID(o.ID)
	}
	if o.ShiftNo != "" {
		builder = builder.SetShiftNo(o.ShiftNo)
	}
	if o.OrderType != "" {
		builder = builder.SetOrderType(o.OrderType)
	}
	if o.OriginOrderID != "" {
		builder = builder.SetOriginOrderID(o.OriginOrderID)
	}
	if o.Refund != nil {
		b, mErr := json.Marshal(o.Refund)
		if mErr != nil {
			return fmt.Errorf("failed to marshal refund: %w", mErr)
		}
		builder = builder.SetRefund(b)
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
		builder = builder.SetOrderStatus(o.OrderStatus)
	}
	if o.PaymentStatus != "" {
		builder = builder.SetPaymentStatus(o.PaymentStatus)
	}
	if o.FulfillmentStatus != "" {
		builder = builder.SetFulfillmentStatus(o.FulfillmentStatus)
	}
	if o.TableStatus != "" {
		builder = builder.SetTableStatus(o.TableStatus)
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

	if o.Store != nil {
		b, mErr := json.Marshal(o.Store)
		if mErr != nil {
			return fmt.Errorf("failed to marshal store: %w", mErr)
		}
		builder = builder.SetStore(b)
	}
	if o.Channel != nil {
		b, mErr := json.Marshal(o.Channel)
		if mErr != nil {
			return fmt.Errorf("failed to marshal channel: %w", mErr)
		}
		builder = builder.SetChannel(b)
	}
	if o.Pos != nil {
		b, mErr := json.Marshal(o.Pos)
		if mErr != nil {
			return fmt.Errorf("failed to marshal pos: %w", mErr)
		}
		builder = builder.SetPos(b)
	}
	if o.Cashier != nil {
		b, mErr := json.Marshal(o.Cashier)
		if mErr != nil {
			return fmt.Errorf("failed to marshal cashier: %w", mErr)
		}
		builder = builder.SetCashier(b)
	}

	if o.Member != nil {
		b, mErr := json.Marshal(o.Member)
		if mErr != nil {
			return fmt.Errorf("failed to marshal member: %w", mErr)
		}
		builder = builder.SetMember(b)
	}
	if o.Takeaway != nil {
		b, mErr := json.Marshal(o.Takeaway)
		if mErr != nil {
			return fmt.Errorf("failed to marshal takeaway: %w", mErr)
		}
		builder = builder.SetTakeaway(b)
	}

	if o.Cart != nil {
		b, mErr := json.Marshal(*o.Cart)
		if mErr != nil {
			return fmt.Errorf("failed to marshal cart: %w", mErr)
		}
		builder = builder.SetCart(b)
	}
	if o.Products != nil {
		b, mErr := json.Marshal(*o.Products)
		if mErr != nil {
			return fmt.Errorf("failed to marshal products: %w", mErr)
		}
		builder = builder.SetProducts(b)
	}
	if o.Promotions != nil {
		b, mErr := json.Marshal(*o.Promotions)
		if mErr != nil {
			return fmt.Errorf("failed to marshal promotions: %w", mErr)
		}
		builder = builder.SetPromotions(b)
	}
	if o.Coupons != nil {
		b, mErr := json.Marshal(*o.Coupons)
		if mErr != nil {
			return fmt.Errorf("failed to marshal coupons: %w", mErr)
		}
		builder = builder.SetCoupons(b)
	}
	if o.TaxRates != nil {
		b, mErr := json.Marshal(*o.TaxRates)
		if mErr != nil {
			return fmt.Errorf("failed to marshal tax_rates: %w", mErr)
		}
		builder = builder.SetTaxRates(b)
	}
	if o.Fees != nil {
		b, mErr := json.Marshal(*o.Fees)
		if mErr != nil {
			return fmt.Errorf("failed to marshal fees: %w", mErr)
		}
		builder = builder.SetFees(b)
	}
	if o.Payments != nil {
		b, mErr := json.Marshal(*o.Payments)
		if mErr != nil {
			return fmt.Errorf("failed to marshal payments: %w", mErr)
		}
		builder = builder.SetPayments(b)
	}
	if o.RefundsProducts != nil {
		b, mErr := json.Marshal(*o.RefundsProducts)
		if mErr != nil {
			return fmt.Errorf("failed to marshal refunds_products: %w", mErr)
		}
		builder = builder.SetRefundsProducts(b)
	}
	if o.Amount != nil {
		b, mErr := json.Marshal(o.Amount)
		if mErr != nil {
			return fmt.Errorf("failed to marshal amount: %w", mErr)
		}
		builder = builder.SetAmount(b)
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
		builder = builder.SetOrderType(domain.OrderType(o.OrderType))
	}
	if o.OriginOrderID != "" {
		builder = builder.SetOriginOrderID(o.OriginOrderID)
	}
	if o.Refund != nil {
		b, mErr := json.Marshal(o.Refund)
		if mErr != nil {
			return fmt.Errorf("failed to marshal refund: %w", mErr)
		}
		builder = builder.SetRefund(b)
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
		builder = builder.SetDiningMode(o.DiningMode)
	}
	if o.OrderStatus != "" {
		builder = builder.SetOrderStatus(o.OrderStatus)
	}
	if o.PaymentStatus != "" {
		builder = builder.SetPaymentStatus(o.PaymentStatus)
	}
	if o.FulfillmentStatus != "" {
		builder = builder.SetFulfillmentStatus(o.FulfillmentStatus)
	}
	if o.TableStatus != "" {
		builder = builder.SetTableStatus(o.TableStatus)
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

	if o.Store != nil {
		b, mErr := json.Marshal(o.Store)
		if mErr != nil {
			return fmt.Errorf("failed to marshal store: %w", mErr)
		}
		builder = builder.SetStore(b)
	}
	if o.Channel != nil {
		b, mErr := json.Marshal(o.Channel)
		if mErr != nil {
			return fmt.Errorf("failed to marshal channel: %w", mErr)
		}
		builder = builder.SetChannel(b)
	}
	if o.Pos != nil {
		b, mErr := json.Marshal(o.Pos)
		if mErr != nil {
			return fmt.Errorf("failed to marshal pos: %w", mErr)
		}
		builder = builder.SetPos(b)
	}
	if o.Cashier != nil {
		b, mErr := json.Marshal(o.Cashier)
		if mErr != nil {
			return fmt.Errorf("failed to marshal cashier: %w", mErr)
		}
		builder = builder.SetCashier(b)
	}
	if o.Member != nil {
		b, mErr := json.Marshal(o.Member)
		if mErr != nil {
			return fmt.Errorf("failed to marshal member: %w", mErr)
		}
		builder = builder.SetMember(b)
	}
	if o.Takeaway != nil {
		b, mErr := json.Marshal(o.Takeaway)
		if mErr != nil {
			return fmt.Errorf("failed to marshal takeaway: %w", mErr)
		}
		builder = builder.SetTakeaway(b)
	}
	if o.Cart != nil {
		b, mErr := json.Marshal(*o.Cart)
		if mErr != nil {
			return fmt.Errorf("failed to marshal cart: %w", mErr)
		}
		builder = builder.SetCart(b)
	}
	if o.Products != nil {
		b, mErr := json.Marshal(*o.Products)
		if mErr != nil {
			return fmt.Errorf("failed to marshal products: %w", mErr)
		}
		builder = builder.SetProducts(b)
	}
	if o.Promotions != nil {
		b, mErr := json.Marshal(*o.Promotions)
		if mErr != nil {
			return fmt.Errorf("failed to marshal promotions: %w", mErr)
		}
		builder = builder.SetPromotions(b)
	}
	if o.Coupons != nil {
		b, mErr := json.Marshal(*o.Coupons)
		if mErr != nil {
			return fmt.Errorf("failed to marshal coupons: %w", mErr)
		}
		builder = builder.SetCoupons(b)
	}
	if o.TaxRates != nil {
		b, mErr := json.Marshal(*o.TaxRates)
		if mErr != nil {
			return fmt.Errorf("failed to marshal tax_rates: %w", mErr)
		}
		builder = builder.SetTaxRates(b)
	}
	if o.Fees != nil {
		b, mErr := json.Marshal(*o.Fees)
		if mErr != nil {
			return fmt.Errorf("failed to marshal fees: %w", mErr)
		}
		builder = builder.SetFees(b)
	}
	if o.Payments != nil {
		b, mErr := json.Marshal(*o.Payments)
		if mErr != nil {
			return fmt.Errorf("failed to marshal payments: %w", mErr)
		}
		builder = builder.SetPayments(b)
	}
	if o.RefundsProducts != nil {
		b, mErr := json.Marshal(*o.RefundsProducts)
		if mErr != nil {
			return fmt.Errorf("failed to marshal refunds_products: %w", mErr)
		}
		builder = builder.SetRefundsProducts(b)
	}
	if o.Amount != nil {
		b, mErr := json.Marshal(o.Amount)
		if mErr != nil {
			return fmt.Errorf("failed to marshal amount: %w", mErr)
		}
		builder = builder.SetAmount(b)
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

	query := repo.Client.Order.Query()
	if params.MerchantID != uuid.Nil {
		query.Where(entorder.MerchantID(params.MerchantID))
	}
	if params.StoreID != uuid.Nil {
		query.Where(entorder.StoreID(params.StoreID))
	}
	if params.BusinessDate != "" {
		query.Where(entorder.BusinessDate(params.BusinessDate))
	}
	if params.OrderNo != "" {
		query.Where(entorder.OrderNo(params.OrderNo))
	}
	if params.OrderType != "" {
		query.Where(entorder.OrderTypeEQ(params.OrderType))
	}
	if params.OrderStatus != "" {
		query.Where(entorder.OrderStatusEQ(params.OrderStatus))
	}
	if params.PaymentStatus != "" {
		query.Where(entorder.PaymentStatusEQ(params.PaymentStatus))
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
		o, cErr := convertOrderToDomain(eo)
		if cErr != nil {
			return nil, 0, fmt.Errorf("failed to convert order: %w", cErr)
		}
		res = append(res, o)
	}
	return res, total, nil
}

func convertOrderToDomain(eo *ent.Order) (*domain.Order, error) {
	if eo == nil {
		return nil, nil
	}

	var refund *domain.OrderRefund
	if len(eo.Refund) > 0 {
		var v domain.OrderRefund
		if uErr := json.Unmarshal(eo.Refund, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal refund: %w", uErr)
		}
		refund = &v
	}

	var store *domain.OrderStore
	if len(eo.Store) > 0 {
		var v domain.OrderStore
		if uErr := json.Unmarshal(eo.Store, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal store: %w", uErr)
		}
		store = &v
	}
	var channel *domain.OrderChannel
	if len(eo.Channel) > 0 {
		var v domain.OrderChannel
		if uErr := json.Unmarshal(eo.Channel, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal channel: %w", uErr)
		}
		channel = &v
	}
	var pos *domain.OrderPOS
	if len(eo.Pos) > 0 {
		var v domain.OrderPOS
		if uErr := json.Unmarshal(eo.Pos, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal pos: %w", uErr)
		}
		pos = &v
	}
	var cashier *domain.OrderCashier
	if len(eo.Cashier) > 0 {
		var v domain.OrderCashier
		if uErr := json.Unmarshal(eo.Cashier, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal cashier: %w", uErr)
		}
		cashier = &v
	}

	var member *domain.OrderMember
	if len(eo.Member) > 0 {
		var v domain.OrderMember
		if uErr := json.Unmarshal(eo.Member, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal member: %w", uErr)
		}
		member = &v
	}
	var takeaway *domain.OrderTakeaway
	if len(eo.Takeaway) > 0 {
		var v domain.OrderTakeaway
		if uErr := json.Unmarshal(eo.Takeaway, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal takeaway: %w", uErr)
		}
		takeaway = &v
	}

	var cart *[]domain.OrderProduct
	if len(eo.Cart) > 0 {
		var v []domain.OrderProduct
		if uErr := json.Unmarshal(eo.Cart, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal cart: %w", uErr)
		}
		cart = &v
	}
	var products *[]domain.OrderProduct
	if len(eo.Products) > 0 {
		var v []domain.OrderProduct
		if uErr := json.Unmarshal(eo.Products, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal products: %w", uErr)
		}
		products = &v
	}
	var promotions *[]domain.OrderPromotion
	if len(eo.Promotions) > 0 {
		var v []domain.OrderPromotion
		if uErr := json.Unmarshal(eo.Promotions, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal promotions: %w", uErr)
		}
		promotions = &v
	}
	var coupons *[]domain.OrderCoupon
	if len(eo.Coupons) > 0 {
		var v []domain.OrderCoupon
		if uErr := json.Unmarshal(eo.Coupons, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal coupons: %w", uErr)
		}
		coupons = &v
	}
	var taxRates *[]domain.OrderTaxRate
	if len(eo.TaxRates) > 0 {
		var v []domain.OrderTaxRate
		if uErr := json.Unmarshal(eo.TaxRates, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal tax_rates: %w", uErr)
		}
		taxRates = &v
	}
	var fees *[]domain.OrderFee
	if len(eo.Fees) > 0 {
		var v []domain.OrderFee
		if uErr := json.Unmarshal(eo.Fees, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal fees: %w", uErr)
		}
		fees = &v
	}
	var payments *[]domain.OrderPayment
	if len(eo.Payments) > 0 {
		var v []domain.OrderPayment
		if uErr := json.Unmarshal(eo.Payments, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal payments: %w", uErr)
		}
		payments = &v
	}
	var refundsProducts *[]domain.OrderProduct
	if len(eo.RefundsProducts) > 0 {
		var v []domain.OrderProduct
		if uErr := json.Unmarshal(eo.RefundsProducts, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal refunds_products: %w", uErr)
		}
		refundsProducts = &v
	}
	var amount *domain.OrderAmount
	if len(eo.Amount) > 0 {
		var v domain.OrderAmount
		if uErr := json.Unmarshal(eo.Amount, &v); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal amount: %w", uErr)
		}
		amount = &v
	}

	return &domain.Order{
		ID:        eo.ID,
		CreatedAt: eo.CreatedAt,
		UpdatedAt: eo.UpdatedAt,

		MerchantID: eo.MerchantID,
		StoreID:    eo.StoreID,

		BusinessDate: eo.BusinessDate,
		ShiftNo:      eo.ShiftNo,
		OrderNo:      eo.OrderNo,

		OrderType:     eo.OrderType,
		OriginOrderID: eo.OriginOrderID,
		Refund:        refund,

		OpenedAt:    eo.OpenedAt,
		PlacedAt:    eo.PlacedAt,
		PaidAt:      eo.PaidAt,
		CompletedAt: eo.CompletedAt,

		OpenedBy: eo.OpenedBy,
		PlacedBy: eo.PlacedBy,
		PaidBy:   eo.PaidBy,

		DiningMode:        eo.DiningMode,
		OrderStatus:       eo.OrderStatus,
		PaymentStatus:     eo.PaymentStatus,
		FulfillmentStatus: eo.FulfillmentStatus,
		TableStatus:       eo.TableStatus,

		TableID:       eo.TableID,
		TableName:     eo.TableName,
		TableCapacity: eo.TableCapacity,
		GuestCount:    eo.GuestCount,

		MergedToOrderID: eo.MergedToOrderID,
		MergedAt:        eo.MergedAt,

		Store:   store,
		Channel: channel,
		Pos:     pos,
		Cashier: cashier,

		Member:   member,
		Takeaway: takeaway,

		Cart:            cart,
		Products:        products,
		Promotions:      promotions,
		Coupons:         coupons,
		TaxRates:        taxRates,
		Fees:            fees,
		Payments:        payments,
		RefundsProducts: refundsProducts,
		Amount:          amount,
	}, nil
}
