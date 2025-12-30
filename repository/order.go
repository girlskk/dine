package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

	eo, err := repo.Client.Order.Query().
		Where(entorder.ID(id)).
		WithOrderProducts().
		Only(ctx)
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
		SetDiningMode(o.DiningMode).
		SetChannel(o.Channel)

	if o.ID != uuid.Nil {
		builder = builder.SetID(o.ID)
	}
	if o.ShiftNo != "" {
		builder = builder.SetShiftNo(o.ShiftNo)
	}
	if o.OrderType != "" {
		builder = builder.SetOrderType(o.OrderType)
	}

	// Refund
	if o.Refund.OriginOrderID != "" || o.Refund.OriginOrderNo != "" || o.Refund.Reason != "" {
		b, mErr := json.Marshal(o.Refund)
		if mErr != nil {
			return fmt.Errorf("failed to marshal refund: %w", mErr)
		}
		builder = builder.SetRefund(b)
	}

	// 时间字段
	if !o.PlacedAt.IsZero() {
		builder = builder.SetPlacedAt(o.PlacedAt)
	}
	if !o.PaidAt.IsZero() {
		builder = builder.SetPaidAt(o.PaidAt)
	}
	if !o.CompletedAt.IsZero() {
		builder = builder.SetCompletedAt(o.CompletedAt)
	}

	if o.PlacedBy != "" {
		builder = builder.SetPlacedBy(o.PlacedBy)
	}
	if o.OrderStatus != "" {
		builder = builder.SetOrderStatus(o.OrderStatus)
	}
	if o.PaymentStatus != "" {
		builder = builder.SetPaymentStatus(o.PaymentStatus)
	}

	if o.TableID != "" {
		builder = builder.SetTableID(o.TableID)
	}
	if o.TableName != "" {
		builder = builder.SetTableName(o.TableName)
	}
	if o.GuestCount != 0 {
		builder = builder.SetGuestCount(o.GuestCount)
	}

	// JSON 字段
	if o.Store.ID != uuid.Nil {
		b, mErr := json.Marshal(o.Store)
		if mErr != nil {
			return fmt.Errorf("failed to marshal store: %w", mErr)
		}
		builder = builder.SetStore(b)
	}
	if o.Pos.PosID != "" {
		b, mErr := json.Marshal(o.Pos)
		if mErr != nil {
			return fmt.Errorf("failed to marshal pos: %w", mErr)
		}
		builder = builder.SetPos(b)
	}
	if o.Cashier.CashierID != "" {
		b, mErr := json.Marshal(o.Cashier)
		if mErr != nil {
			return fmt.Errorf("failed to marshal cashier: %w", mErr)
		}
		builder = builder.SetCashier(b)
	}

	if len(o.TaxRates) > 0 {
		b, mErr := json.Marshal(o.TaxRates)
		if mErr != nil {
			return fmt.Errorf("failed to marshal tax_rates: %w", mErr)
		}
		builder = builder.SetTaxRates(b)
	}
	if len(o.Fees) > 0 {
		b, mErr := json.Marshal(o.Fees)
		if mErr != nil {
			return fmt.Errorf("failed to marshal fees: %w", mErr)
		}
		builder = builder.SetFees(b)
	}
	if len(o.Payments) > 0 {
		b, mErr := json.Marshal(o.Payments)
		if mErr != nil {
			return fmt.Errorf("failed to marshal payments: %w", mErr)
		}
		builder = builder.SetPayments(b)
	}

	// Amount
	b, mErr := json.Marshal(o.Amount)
	if mErr != nil {
		return fmt.Errorf("failed to marshal amount: %w", mErr)
	}
	builder = builder.SetAmount(b)

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

	// 创建订单商品明细
	if len(o.OrderProducts) > 0 {
		for i := range o.OrderProducts {
			op := &o.OrderProducts[i]
			op.OrderID = created.ID
			if err := repo.createOrderProduct(ctx, op); err != nil {
				return fmt.Errorf("failed to create order product: %w", err)
			}
		}
	}

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
		builder = builder.SetOrderType(o.OrderType)
	}

	// Refund
	if o.Refund.OriginOrderID != "" || o.Refund.OriginOrderNo != "" || o.Refund.Reason != "" {
		b, mErr := json.Marshal(o.Refund)
		if mErr != nil {
			return fmt.Errorf("failed to marshal refund: %w", mErr)
		}
		builder = builder.SetRefund(b)
	}

	if !o.PlacedAt.IsZero() {
		builder = builder.SetPlacedAt(o.PlacedAt)
	}
	if !o.PaidAt.IsZero() {
		builder = builder.SetPaidAt(o.PaidAt)
	}
	if !o.CompletedAt.IsZero() {
		builder = builder.SetCompletedAt(o.CompletedAt)
	}

	if o.PlacedBy != "" {
		builder = builder.SetPlacedBy(o.PlacedBy)
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
	if o.Channel != "" {
		builder = builder.SetChannel(o.Channel)
	}

	if o.TableID != "" {
		builder = builder.SetTableID(o.TableID)
	}
	if o.TableName != "" {
		builder = builder.SetTableName(o.TableName)
	}
	if o.GuestCount != 0 {
		builder = builder.SetGuestCount(o.GuestCount)
	}

	// JSON 字段
	if o.Store.ID != uuid.Nil {
		b, mErr := json.Marshal(o.Store)
		if mErr != nil {
			return fmt.Errorf("failed to marshal store: %w", mErr)
		}
		builder = builder.SetStore(b)
	}
	if o.Pos.PosID != "" {
		b, mErr := json.Marshal(o.Pos)
		if mErr != nil {
			return fmt.Errorf("failed to marshal pos: %w", mErr)
		}
		builder = builder.SetPos(b)
	}
	if o.Cashier.CashierID != "" {
		b, mErr := json.Marshal(o.Cashier)
		if mErr != nil {
			return fmt.Errorf("failed to marshal cashier: %w", mErr)
		}
		builder = builder.SetCashier(b)
	}

	if len(o.TaxRates) > 0 {
		b, mErr := json.Marshal(o.TaxRates)
		if mErr != nil {
			return fmt.Errorf("failed to marshal tax_rates: %w", mErr)
		}
		builder = builder.SetTaxRates(b)
	}
	if len(o.Fees) > 0 {
		b, mErr := json.Marshal(o.Fees)
		if mErr != nil {
			return fmt.Errorf("failed to marshal fees: %w", mErr)
		}
		builder = builder.SetFees(b)
	}
	if len(o.Payments) > 0 {
		b, mErr := json.Marshal(o.Payments)
		if mErr != nil {
			return fmt.Errorf("failed to marshal payments: %w", mErr)
		}
		builder = builder.SetPayments(b)
	}

	// Amount
	b, mErr := json.Marshal(o.Amount)
	if mErr != nil {
		return fmt.Errorf("failed to marshal amount: %w", mErr)
	}
	builder = builder.SetAmount(b)

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
		WithOrderProducts().
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

func (repo *OrderRepository) createOrderProduct(ctx context.Context, op *domain.OrderProduct) error {
	builder := repo.Client.OrderProduct.Create().
		SetOrderID(op.OrderID).
		SetOrderItemID(op.OrderItemID).
		SetIndex(op.Index).
		SetProductID(op.ProductID).
		SetProductName(op.ProductName).
		SetProductType(op.ProductType).
		SetQty(op.Qty)

	if op.ID != uuid.Nil {
		builder = builder.SetID(op.ID)
	}
	if op.CategoryID != uuid.Nil {
		builder = builder.SetCategoryID(op.CategoryID)
	}
	if op.MenuID != uuid.Nil {
		builder = builder.SetMenuID(op.MenuID)
	}
	if op.UnitID != uuid.Nil {
		builder = builder.SetUnitID(op.UnitID)
	}
	if len(op.SupportTypes) > 0 {
		builder = builder.SetSupportTypes(op.SupportTypes)
	}
	if op.SaleStatus != "" {
		builder = builder.SetSaleStatus(op.SaleStatus)
	}
	if len(op.SaleChannels) > 0 {
		builder = builder.SetSaleChannels(op.SaleChannels)
	}
	if op.MainImage != "" {
		builder = builder.SetMainImage(op.MainImage)
	}
	if op.Description != "" {
		builder = builder.SetDescription(op.Description)
	}

	// 金额字段
	if !op.Subtotal.IsZero() {
		builder = builder.SetSubtotal(op.Subtotal)
	}
	if !op.DiscountAmount.IsZero() {
		builder = builder.SetDiscountAmount(op.DiscountAmount)
	}
	if !op.AmountBeforeTax.IsZero() {
		builder = builder.SetAmountBeforeTax(op.AmountBeforeTax)
	}
	if !op.TaxRate.IsZero() {
		builder = builder.SetTaxRate(op.TaxRate)
	}
	if !op.Tax.IsZero() {
		builder = builder.SetTax(op.Tax)
	}
	if !op.AmountAfterTax.IsZero() {
		builder = builder.SetAmountAfterTax(op.AmountAfterTax)
	}
	if !op.Total.IsZero() {
		builder = builder.SetTotal(op.Total)
	}
	if !op.PromotionDiscount.IsZero() {
		builder = builder.SetPromotionDiscount(op.PromotionDiscount)
	}

	// 退菜信息
	if op.VoidQty != 0 {
		builder = builder.SetVoidQty(op.VoidQty)
	}
	if !op.VoidAmount.IsZero() {
		builder = builder.SetVoidAmount(op.VoidAmount)
	}
	if op.RefundReason != "" {
		builder = builder.SetRefundReason(op.RefundReason)
	}
	if op.RefundedBy != "" {
		builder = builder.SetRefundedBy(op.RefundedBy)
	}
	if !op.RefundedAt.IsZero() {
		builder = builder.SetRefundedAt(op.RefundedAt)
	}

	if op.Note != "" {
		builder = builder.SetNote(op.Note)
	}

	// 套餐信息
	if !op.EstimatedCostPrice.IsZero() {
		builder = builder.SetEstimatedCostPrice(op.EstimatedCostPrice)
	}
	if !op.DeliveryCostPrice.IsZero() {
		builder = builder.SetDeliveryCostPrice(op.DeliveryCostPrice)
	}
	if len(op.SetMealGroups) > 0 {
		builder = builder.SetSetMealGroups(op.SetMealGroups)
	}
	if len(op.SpecRelations) > 0 {
		builder = builder.SetSpecRelations(op.SpecRelations)
	}
	if len(op.AttrRelations) > 0 {
		builder = builder.SetAttrRelations(op.AttrRelations)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		if ent.IsValidationError(err) {
			return domain.ParamsError(fmt.Errorf("invalid order product params: %w", err))
		}
		if ent.IsConstraintError(err) {
			return domain.ConflictError(err)
		}
		return fmt.Errorf("failed to create order product: %w", err)
	}

	op.ID = created.ID
	op.CreatedAt = created.CreatedAt
	op.UpdatedAt = created.UpdatedAt
	return nil
}

func convertOrderToDomain(eo *ent.Order) (*domain.Order, error) {
	if eo == nil {
		return nil, nil
	}

	var refund domain.OrderRefund
	if len(eo.Refund) > 0 {
		if uErr := json.Unmarshal(eo.Refund, &refund); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal refund: %w", uErr)
		}
	}

	var store domain.OrderStore
	if len(eo.Store) > 0 {
		if uErr := json.Unmarshal(eo.Store, &store); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal store: %w", uErr)
		}
	}

	var pos domain.OrderPOS
	if len(eo.Pos) > 0 {
		if uErr := json.Unmarshal(eo.Pos, &pos); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal pos: %w", uErr)
		}
	}

	var cashier domain.OrderCashier
	if len(eo.Cashier) > 0 {
		if uErr := json.Unmarshal(eo.Cashier, &cashier); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal cashier: %w", uErr)
		}
	}

	var taxRates []domain.OrderTaxRate
	if len(eo.TaxRates) > 0 {
		if uErr := json.Unmarshal(eo.TaxRates, &taxRates); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal tax_rates: %w", uErr)
		}
	}

	var fees []domain.OrderFee
	if len(eo.Fees) > 0 {
		if uErr := json.Unmarshal(eo.Fees, &fees); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal fees: %w", uErr)
		}
	}

	var payments []domain.OrderPayment
	if len(eo.Payments) > 0 {
		if uErr := json.Unmarshal(eo.Payments, &payments); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal payments: %w", uErr)
		}
	}

	var amount domain.OrderAmount
	if len(eo.Amount) > 0 {
		if uErr := json.Unmarshal(eo.Amount, &amount); uErr != nil {
			return nil, fmt.Errorf("failed to unmarshal amount: %w", uErr)
		}
	}

	// 转换时间字段
	var placedAt, paidAt, completedAt time.Time
	if eo.PlacedAt != nil {
		placedAt = *eo.PlacedAt
	}
	if eo.PaidAt != nil {
		paidAt = *eo.PaidAt
	}
	if eo.CompletedAt != nil {
		completedAt = *eo.CompletedAt
	}

	// 转换订单商品
	var orderProducts []domain.OrderProduct
	if eo.Edges.OrderProducts != nil {
		orderProducts = make([]domain.OrderProduct, 0, len(eo.Edges.OrderProducts))
		for _, eop := range eo.Edges.OrderProducts {
			op := convertOrderProductToDomain(eop)
			orderProducts = append(orderProducts, op)
		}
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

		OrderType: eo.OrderType,
		Refund:    refund,

		PlacedAt:    placedAt,
		PaidAt:      paidAt,
		CompletedAt: completedAt,

		PlacedBy: eo.PlacedBy,

		DiningMode:    eo.DiningMode,
		OrderStatus:   eo.OrderStatus,
		PaymentStatus: eo.PaymentStatus,

		TableID:    eo.TableID,
		TableName:  eo.TableName,
		GuestCount: eo.GuestCount,

		Store:   store,
		Channel: eo.Channel,
		Pos:     pos,
		Cashier: cashier,

		TaxRates: taxRates,
		Fees:     fees,
		Payments: payments,
		Amount:   amount,

		OrderProducts: orderProducts,
	}, nil
}

func convertOrderProductToDomain(eop *ent.OrderProduct) domain.OrderProduct {
	var refundedAt time.Time
	if eop.RefundedAt != nil {
		refundedAt = *eop.RefundedAt
	}

	var subtotal, discountAmount, amountBeforeTax, taxRate, tax, amountAfterTax, total decimal.Decimal
	var promotionDiscount, voidAmount, estimatedCostPrice, deliveryCostPrice decimal.Decimal

	if eop.Subtotal != nil {
		subtotal = *eop.Subtotal
	}
	if eop.DiscountAmount != nil {
		discountAmount = *eop.DiscountAmount
	}
	if eop.AmountBeforeTax != nil {
		amountBeforeTax = *eop.AmountBeforeTax
	}
	if eop.TaxRate != nil {
		taxRate = *eop.TaxRate
	}
	if eop.Tax != nil {
		tax = *eop.Tax
	}
	if eop.AmountAfterTax != nil {
		amountAfterTax = *eop.AmountAfterTax
	}
	if eop.Total != nil {
		total = *eop.Total
	}
	if eop.PromotionDiscount != nil {
		promotionDiscount = *eop.PromotionDiscount
	}
	if eop.VoidAmount != nil {
		voidAmount = *eop.VoidAmount
	}
	if eop.EstimatedCostPrice != nil {
		estimatedCostPrice = *eop.EstimatedCostPrice
	}
	if eop.DeliveryCostPrice != nil {
		deliveryCostPrice = *eop.DeliveryCostPrice
	}

	return domain.OrderProduct{
		ID:        eop.ID,
		CreatedAt: eop.CreatedAt,
		UpdatedAt: eop.UpdatedAt,

		OrderID:     eop.OrderID,
		OrderItemID: eop.OrderItemID,
		Index:       eop.Index,

		ProductID:    eop.ProductID,
		ProductName:  eop.ProductName,
		ProductType:  eop.ProductType,
		CategoryID:   eop.CategoryID,
		MenuID:       eop.MenuID,
		UnitID:       eop.UnitID,
		SupportTypes: eop.SupportTypes,
		SaleStatus:   eop.SaleStatus,
		SaleChannels: eop.SaleChannels,
		MainImage:    eop.MainImage,
		Description:  eop.Description,

		Qty:             eop.Qty,
		Subtotal:        subtotal,
		DiscountAmount:  discountAmount,
		AmountBeforeTax: amountBeforeTax,
		TaxRate:         taxRate,
		Tax:             tax,
		AmountAfterTax:  amountAfterTax,
		Total:           total,

		PromotionDiscount: promotionDiscount,

		VoidQty:      eop.VoidQty,
		VoidAmount:   voidAmount,
		RefundReason: eop.RefundReason,
		RefundedBy:   eop.RefundedBy,
		RefundedAt:   refundedAt,

		Note: eop.Note,

		EstimatedCostPrice: estimatedCostPrice,
		DeliveryCostPrice:  deliveryCostPrice,
		SetMealGroups:      eop.SetMealGroups,
		SpecRelations:      eop.SpecRelations,
		AttrRelations:      eop.AttrRelations,
	}
}
