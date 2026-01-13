package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	entrefundorder "gitlab.jiguang.dev/pos-dine/dine/ent/refundorder"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RefundOrderRepository = (*RefundOrderRepository)(nil)

type RefundOrderRepository struct {
	Client *ent.Client
}

func NewRefundOrderRepository(client *ent.Client) *RefundOrderRepository {
	return &RefundOrderRepository{Client: client}
}

func (repo *RefundOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.RefundOrder, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RefundOrderRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	eo, err := repo.Client.RefundOrder.Query().
		Where(entrefundorder.ID(id)).
		WithRefundProducts().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get refund order: %w", err)
	}

	return convertRefundOrderToDomain(eo), nil
}

func (repo *RefundOrderRepository) Create(ctx context.Context, ro *domain.RefundOrder) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RefundOrderRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()

	builder := repo.Client.RefundOrder.Create().
		SetMerchantID(ro.MerchantID).
		SetStoreID(ro.StoreID).
		SetBusinessDate(ro.BusinessDate).
		SetRefundNo(ro.RefundNo).
		SetOriginOrderID(ro.OriginOrderID).
		SetOriginOrderNo(ro.OriginOrderNo).
		SetRefundType(ro.RefundType).
		SetStore(ro.Store).
		SetPos(ro.Pos).
		SetCashier(ro.Cashier).
		SetRefundAmount(ro.RefundAmount)

	if ro.ID != uuid.Nil {
		builder = builder.SetID(ro.ID)
	}
	if ro.ShiftNo != "" {
		builder = builder.SetShiftNo(ro.ShiftNo)
	}
	if !ro.OriginPaidAt.IsZero() {
		builder = builder.SetOriginPaidAt(ro.OriginPaidAt)
	}
	if !ro.OriginAmountPaid.IsZero() {
		builder = builder.SetOriginAmountPaid(ro.OriginAmountPaid)
	}
	if ro.RefundStatus != "" {
		builder = builder.SetRefundStatus(ro.RefundStatus)
	}
	if ro.RefundReasonCode != "" {
		builder = builder.SetRefundReasonCode(string(ro.RefundReasonCode))
	}
	if ro.RefundReason != "" {
		builder = builder.SetRefundReason(ro.RefundReason)
	}
	if ro.RefundedBy != uuid.Nil {
		builder = builder.SetRefundedBy(ro.RefundedBy)
	}
	if ro.RefundedByName != "" {
		builder = builder.SetRefundedByName(ro.RefundedByName)
	}
	if ro.ApprovedBy != uuid.Nil {
		builder = builder.SetApprovedBy(ro.ApprovedBy)
	}
	if ro.ApprovedByName != "" {
		builder = builder.SetApprovedByName(ro.ApprovedByName)
	}
	if !ro.ApprovedAt.IsZero() {
		builder = builder.SetApprovedAt(ro.ApprovedAt)
	}
	if !ro.RefundedAt.IsZero() {
		builder = builder.SetRefundedAt(ro.RefundedAt)
	}
	if ro.Channel != "" {
		builder = builder.SetChannel(ro.Channel)
	}
	if len(ro.RefundPayments) > 0 {
		builder = builder.SetRefundPayments(ro.RefundPayments)
	}
	if ro.Remark != "" {
		builder = builder.SetRemark(ro.Remark)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		if ent.IsValidationError(err) {
			return domain.ParamsError(fmt.Errorf("invalid refund order params: %w", err))
		}
		if ent.IsConstraintError(err) {
			return domain.ConflictError(err)
		}
		return fmt.Errorf("failed to create refund order: %w", err)
	}

	// 创建退款商品明细
	if len(ro.RefundProducts) > 0 {
		for i := range ro.RefundProducts {
			rp := &ro.RefundProducts[i]
			rp.RefundOrderID = created.ID
			if err := repo.createRefundProduct(ctx, rp); err != nil {
				return fmt.Errorf("failed to create refund product: %w", err)
			}
		}
	}

	ro.ID = created.ID
	return nil
}

func (repo *RefundOrderRepository) createRefundProduct(ctx context.Context, rp *domain.RefundOrderProduct) error {
	builder := repo.Client.RefundOrderProduct.Create().
		SetRefundOrderID(rp.RefundOrderID).
		SetOriginOrderProductID(rp.OriginOrderProductID).
		SetProductID(rp.ProductID).
		SetProductName(rp.ProductName).
		SetOriginQty(rp.OriginQty).
		SetRefundQty(rp.RefundQty)

	if rp.ID != uuid.Nil {
		builder = builder.SetID(rp.ID)
	}
	if rp.OriginOrderItemID != "" {
		builder = builder.SetOriginOrderItemID(rp.OriginOrderItemID)
	}
	if rp.ProductType != "" {
		builder = builder.SetProductType(rp.ProductType)
	}
	if rp.Category.ID != uuid.Nil {
		builder = builder.SetCategory(rp.Category)
	}
	if rp.MainImage != "" {
		builder = builder.SetMainImage(rp.MainImage)
	}
	if rp.Description != "" {
		builder = builder.SetDescription(rp.Description)
	}
	if !rp.OriginPrice.IsZero() {
		builder = builder.SetOriginPrice(rp.OriginPrice)
	}
	if !rp.OriginSubtotal.IsZero() {
		builder = builder.SetOriginSubtotal(rp.OriginSubtotal)
	}
	if !rp.OriginDiscount.IsZero() {
		builder = builder.SetOriginDiscount(rp.OriginDiscount)
	}
	if !rp.OriginTax.IsZero() {
		builder = builder.SetOriginTax(rp.OriginTax)
	}
	if !rp.OriginTotal.IsZero() {
		builder = builder.SetOriginTotal(rp.OriginTotal)
	}
	if !rp.RefundSubtotal.IsZero() {
		builder = builder.SetRefundSubtotal(rp.RefundSubtotal)
	}
	if !rp.RefundDiscount.IsZero() {
		builder = builder.SetRefundDiscount(rp.RefundDiscount)
	}
	if !rp.RefundTax.IsZero() {
		builder = builder.SetRefundTax(rp.RefundTax)
	}
	if !rp.RefundTotal.IsZero() {
		builder = builder.SetRefundTotal(rp.RefundTotal)
	}
	if len(rp.Groups) > 0 {
		builder = builder.SetGroups(rp.Groups)
	}
	if len(rp.SpecRelations) > 0 {
		builder = builder.SetSpecRelations(rp.SpecRelations)
	}
	if len(rp.AttrRelations) > 0 {
		builder = builder.SetAttrRelations(rp.AttrRelations)
	}
	if rp.RefundReason != "" {
		builder = builder.SetRefundReason(rp.RefundReason)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create refund product: %w", err)
	}
	rp.ID = created.ID
	return nil
}

func (repo *RefundOrderRepository) Update(ctx context.Context, ro *domain.RefundOrder) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RefundOrderRepository.Update")
	defer func() { util.SpanErrFinish(span, err) }()

	builder := repo.Client.RefundOrder.UpdateOneID(ro.ID)

	if ro.RefundStatus != "" {
		builder = builder.SetRefundStatus(ro.RefundStatus)
	}
	if ro.RefundReasonCode != "" {
		builder = builder.SetRefundReasonCode(string(ro.RefundReasonCode))
	}
	if ro.RefundReason != "" {
		builder = builder.SetRefundReason(ro.RefundReason)
	}
	if ro.ApprovedBy != uuid.Nil {
		builder = builder.SetApprovedBy(ro.ApprovedBy)
	}
	if ro.ApprovedByName != "" {
		builder = builder.SetApprovedByName(ro.ApprovedByName)
	}
	if !ro.ApprovedAt.IsZero() {
		builder = builder.SetApprovedAt(ro.ApprovedAt)
	}
	if !ro.RefundedAt.IsZero() {
		builder = builder.SetRefundedAt(ro.RefundedAt)
	}
	if ro.RefundedBy != uuid.Nil {
		builder = builder.SetRefundedBy(ro.RefundedBy)
	}
	if ro.RefundedByName != "" {
		builder = builder.SetRefundedByName(ro.RefundedByName)
	}
	if len(ro.RefundPayments) > 0 {
		builder = builder.SetRefundPayments(ro.RefundPayments)
	}
	if ro.Remark != "" {
		builder = builder.SetRemark(ro.Remark)
	}

	_, err = builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.NotFoundError(err)
		}
		if ent.IsConstraintError(err) {
			return domain.ConflictError(err)
		}
		return fmt.Errorf("failed to update refund order: %w", err)
	}
	return nil
}

func (repo *RefundOrderRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RefundOrderRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.RefundOrder.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.NotFoundError(err)
		}
		return fmt.Errorf("failed to delete refund order: %w", err)
	}
	return nil
}

func (repo *RefundOrderRepository) List(ctx context.Context, params domain.RefundOrderListParams) (res []*domain.RefundOrder, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RefundOrderRepository.List")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.RefundOrder.Query()

	if params.MerchantID != uuid.Nil {
		query = query.Where(entrefundorder.MerchantID(params.MerchantID))
	}
	if params.StoreID != uuid.Nil {
		query = query.Where(entrefundorder.StoreID(params.StoreID))
	}
	if params.OriginOrderID != uuid.Nil {
		query = query.Where(entrefundorder.OriginOrderID(params.OriginOrderID))
	}

	if params.RefundNo != "" {
		query = query.Where(entrefundorder.RefundNo(params.RefundNo))
	}
	if params.RefundType != "" {
		query = query.Where(entrefundorder.RefundTypeEQ(params.RefundType))
	}
	if params.RefundStatus != "" {
		query = query.Where(entrefundorder.RefundStatusEQ(params.RefundStatus))
	}

	total, err = query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count refund orders: %w", err)
	}

	if params.Size > 0 {
		query = query.Limit(params.Size)
	}
	if params.Page > 0 && params.Size > 0 {
		query = query.Offset((params.Page - 1) * params.Size)
	}

	eos, err := query.
		Order(entrefundorder.ByCreatedAt()).
		WithRefundProducts().
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list refund orders: %w", err)
	}

	res = make([]*domain.RefundOrder, len(eos))
	for i, eo := range eos {
		res[i] = convertRefundOrderToDomain(eo)
	}
	return res, total, nil
}

func (repo *RefundOrderRepository) FindByOriginOrderID(ctx context.Context, originOrderID uuid.UUID) (res []*domain.RefundOrder, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RefundOrderRepository.FindByOriginOrderID")
	defer func() { util.SpanErrFinish(span, err) }()

	eos, err := repo.Client.RefundOrder.Query().
		Where(entrefundorder.OriginOrderID(originOrderID)).
		WithRefundProducts().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find refund orders by origin order id: %w", err)
	}

	res = make([]*domain.RefundOrder, len(eos))
	for i, eo := range eos {
		res[i] = convertRefundOrderToDomain(eo)
	}
	return res, nil
}

func convertRefundOrderToDomain(eo *ent.RefundOrder) *domain.RefundOrder {
	ro := &domain.RefundOrder{
		ID:               eo.ID,
		CreatedAt:        eo.CreatedAt,
		UpdatedAt:        eo.UpdatedAt,
		MerchantID:       eo.MerchantID,
		StoreID:          eo.StoreID,
		BusinessDate:     eo.BusinessDate,
		ShiftNo:          eo.ShiftNo,
		RefundNo:         eo.RefundNo,
		OriginOrderID:    eo.OriginOrderID,
		OriginOrderNo:    eo.OriginOrderNo,
		RefundType:       eo.RefundType,
		RefundStatus:     eo.RefundStatus,
		RefundReasonCode: domain.RefundReasonCode(eo.RefundReasonCode),
		RefundReason:     eo.RefundReason,
		RefundedBy:       eo.RefundedBy,
		RefundedByName:   eo.RefundedByName,
		ApprovedBy:       eo.ApprovedBy,
		ApprovedByName:   eo.ApprovedByName,
		Store:            eo.Store,
		Channel:          eo.Channel,
		Pos:              eo.Pos,
		Cashier:          eo.Cashier,
		RefundAmount:     eo.RefundAmount,
		RefundPayments:   eo.RefundPayments,
		Remark:           eo.Remark,
	}

	if eo.OriginPaidAt != nil {
		ro.OriginPaidAt = *eo.OriginPaidAt
	}
	if eo.OriginAmountPaid != nil {
		ro.OriginAmountPaid = *eo.OriginAmountPaid
	}
	if eo.ApprovedAt != nil {
		ro.ApprovedAt = *eo.ApprovedAt
	}
	if eo.RefundedAt != nil {
		ro.RefundedAt = *eo.RefundedAt
	}

	if eo.Edges.RefundProducts != nil {
		ro.RefundProducts = make([]domain.RefundOrderProduct, len(eo.Edges.RefundProducts))
		for i, ep := range eo.Edges.RefundProducts {
			ro.RefundProducts[i] = convertRefundProductToDomain(ep)
		}
	}

	return ro
}

func convertRefundProductToDomain(ep *ent.RefundOrderProduct) domain.RefundOrderProduct {
	rp := domain.RefundOrderProduct{
		ID:                   ep.ID,
		CreatedAt:            ep.CreatedAt,
		UpdatedAt:            ep.UpdatedAt,
		RefundOrderID:        ep.RefundOrderID,
		OriginOrderProductID: ep.OriginOrderProductID,
		OriginOrderItemID:    ep.OriginOrderItemID,
		ProductID:            ep.ProductID,
		ProductName:          ep.ProductName,
		ProductType:          ep.ProductType,
		Category:             ep.Category,
		MainImage:            ep.MainImage,
		Description:          ep.Description,
		OriginQty:            ep.OriginQty,
		RefundQty:            ep.RefundQty,
		Groups:               ep.Groups,
		SpecRelations:        ep.SpecRelations,
		AttrRelations:        ep.AttrRelations,
		RefundReason:         ep.RefundReason,
	}

	if ep.OriginPrice != nil {
		rp.OriginPrice = *ep.OriginPrice
	}
	if ep.OriginSubtotal != nil {
		rp.OriginSubtotal = *ep.OriginSubtotal
	}
	if ep.OriginDiscount != nil {
		rp.OriginDiscount = *ep.OriginDiscount
	}
	if ep.OriginTax != nil {
		rp.OriginTax = *ep.OriginTax
	}
	if ep.OriginTotal != nil {
		rp.OriginTotal = *ep.OriginTotal
	}
	if ep.RefundSubtotal != nil {
		rp.RefundSubtotal = *ep.RefundSubtotal
	}
	if ep.RefundDiscount != nil {
		rp.RefundDiscount = *ep.RefundDiscount
	}
	if ep.RefundTax != nil {
		rp.RefundTax = *ep.RefundTax
	}
	if ep.RefundTotal != nil {
		rp.RefundTotal = *ep.RefundTotal
	}

	return rp
}
