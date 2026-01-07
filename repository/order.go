package repository

import (
	"context"
	"fmt"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	entorder "gitlab.jiguang.dev/pos-dine/dine/ent/order"
	entorderproduct "gitlab.jiguang.dev/pos-dine/dine/ent/orderproduct"
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

	if o.PlacedBy != uuid.Nil {
		builder = builder.SetPlacedBy(o.PlacedBy)
	}
	if o.OrderStatus != "" {
		builder = builder.SetOrderStatus(o.OrderStatus)
	}
	if o.PaymentStatus != "" {
		builder = builder.SetPaymentStatus(o.PaymentStatus)
	}

	if o.TableID != uuid.Nil {
		builder = builder.SetTableID(o.TableID)
	}
	if o.TableName != "" {
		builder = builder.SetTableName(o.TableName)
	}
	if o.GuestCount != 0 {
		builder = builder.SetGuestCount(o.GuestCount)
	}

	if o.Store.ID != uuid.Nil {
		builder = builder.SetStore(o.Store)
	}
	if o.Pos.ID != uuid.Nil {
		builder = builder.SetPos(o.Pos)
	}
	if o.Cashier.CashierID != uuid.Nil {
		builder = builder.SetCashier(o.Cashier)
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

	// Amount
	builder = builder.SetAmount(o.Amount)

	// Remark
	if o.Remark != "" {
		builder = builder.SetRemark(o.Remark)
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

	if !o.PlacedAt.IsZero() {
		builder = builder.SetPlacedAt(o.PlacedAt)
	}
	if !o.PaidAt.IsZero() {
		builder = builder.SetPaidAt(o.PaidAt)
	}
	if !o.CompletedAt.IsZero() {
		builder = builder.SetCompletedAt(o.CompletedAt)
	}

	if o.PlacedBy != uuid.Nil {
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

	if o.TableID != uuid.Nil {
		builder = builder.SetTableID(o.TableID)
	}
	if o.TableName != "" {
		builder = builder.SetTableName(o.TableName)
	}
	if o.GuestCount != 0 {
		builder = builder.SetGuestCount(o.GuestCount)
	}

	if o.Store.ID != uuid.Nil {
		builder = builder.SetStore(o.Store)
	}
	if o.Pos.ID != uuid.Nil {
		builder = builder.SetPos(o.Pos)
	}
	if o.Cashier.CashierID != uuid.Nil {
		builder = builder.SetCashier(o.Cashier)
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

	// Amount
	builder = builder.SetAmount(o.Amount)

	// Remark
	if o.Remark != "" {
		builder = builder.SetRemark(o.Remark)
	}

	_, err = builder.Save(ctx)
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

	// 更新订单商品（全量替换：先删除旧商品，再插入新商品）
	if len(o.OrderProducts) > 0 {
		// 删除旧商品
		_, err = repo.Client.OrderProduct.Delete().
			Where(entorderproduct.OrderID(o.ID)).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete old order products: %w", err)
		}

		// 插入新商品
		for i := range o.OrderProducts {
			op := &o.OrderProducts[i]
			op.OrderID = o.ID
			if err := repo.createOrderProduct(ctx, op); err != nil {
				return fmt.Errorf("failed to create order product: %w", err)
			}
		}
	}

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
		res = append(res, convertOrderToDomain(eo))
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
		SetQty(op.Qty).
		SetIsGift(op.IsGift)

	if op.ID != uuid.Nil {
		builder = builder.SetID(op.ID)
	}
	if op.CategoryID != uuid.Nil {
		builder = builder.SetCategoryID(op.CategoryID)
	}
	if op.UnitID != uuid.Nil {
		builder = builder.SetUnitID(op.UnitID)
	}
	if op.MainImage != "" {
		builder = builder.SetMainImage(op.MainImage)
	}
	if op.Description != "" {
		builder = builder.SetDescription(op.Description)
	}

	// 金额字段
	if !op.Price.IsZero() {
		builder = builder.SetPrice(op.Price)
	}
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
	if op.RefundedBy != uuid.Nil {
		builder = builder.SetRefundedBy(op.RefundedBy)
	}
	if !op.RefundedAt.IsZero() {
		builder = builder.SetRefundedAt(op.RefundedAt)
	}

	if op.Note != "" {
		builder = builder.SetNote(op.Note)
	}

	// 套餐信息
	if len(op.Groups) > 0 {
		builder = builder.SetGroups(op.Groups)
	}
	if len(op.SpecRelations) > 0 {
		builder = builder.SetSpecRelations(op.SpecRelations)
	}
	if len(op.AttrRelations) > 0 {
		builder = builder.SetAttrRelations(op.AttrRelations)
	}

	_, err := builder.Save(ctx)
	if err != nil {
		if ent.IsValidationError(err) {
			return domain.ParamsError(fmt.Errorf("invalid order product params: %w", err))
		}
		if ent.IsConstraintError(err) {
			return domain.ConflictError(err)
		}
		return fmt.Errorf("failed to create order product: %w", err)
	}
	return nil
}

func convertOrderToDomain(eo *ent.Order) *domain.Order {
	if eo == nil {
		return nil
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

		Store:   eo.Store,
		Channel: eo.Channel,
		Pos:     eo.Pos,
		Cashier: eo.Cashier,

		TaxRates: eo.TaxRates,
		Fees:     eo.Fees,
		Payments: eo.Payments,
		Amount:   eo.Amount,

		Remark: eo.Remark,

		OrderProducts: orderProducts,
	}
}

func convertOrderProductToDomain(eop *ent.OrderProduct) domain.OrderProduct {
	var refundedAt time.Time
	if eop.RefundedAt != nil {
		refundedAt = *eop.RefundedAt
	}

	var subtotal, discountAmount, amountBeforeTax, taxRate, tax, amountAfterTax, total decimal.Decimal
	var promotionDiscount, voidAmount, price decimal.Decimal

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
	if eop.Price != nil {
		price = *eop.Price
	}

	return domain.OrderProduct{
		ID:        eop.ID,
		CreatedAt: eop.CreatedAt,
		UpdatedAt: eop.UpdatedAt,

		OrderID:     eop.OrderID,
		OrderItemID: eop.OrderItemID,
		Index:       eop.Index,

		ProductID:   eop.ProductID,
		ProductName: eop.ProductName,
		ProductType: eop.ProductType,
		CategoryID:  eop.CategoryID,
		UnitID:      eop.UnitID,
		MainImage:   eop.MainImage,
		Description: eop.Description,

		Qty:             eop.Qty,
		IsGift:          eop.IsGift,
		Price:           price,
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

		Groups:        eop.Groups,
		SpecRelations: eop.SpecRelations,
		AttrRelations: eop.AttrRelations,
	}
}
