package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.PointSettlementInteractor = (*PointSettlementInteractor)(nil)

type PointSettlementInteractor struct {
	ds domain.DataStore
}

func NewPointSettlementInteractor(dataStore domain.DataStore) *PointSettlementInteractor {
	return &PointSettlementInteractor{
		ds: dataStore,
	}
}

func (p *PointSettlementInteractor) PagedListBySearch(ctx context.Context, page *upagination.Pagination,
	params domain.PointSettlementSearchParams,
) (res *domain.PointSettlementSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return p.ds.PointSettlementRepo().PagedListBySearch(ctx, page, params)
}

func (p *PointSettlementInteractor) Approve(ctx context.Context, id int) error {
	user := domain.FromAdminUserContext(ctx)
	return p.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 获取账单并加锁
		settlement, err := ds.PointSettlementRepo().FindByIDForUpdate(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrPointSettlementNotExists)
			}
			return err
		}
		if settlement.Status != domain.PointSettlementStatusPending {
			return domain.ParamsError(domain.ErrPointSettlementStatusInvalid)
		}
		// 2. 获取门店账户并加锁
		account, err := ds.StoreAccountRepo().FindByStoreForUpdate(ctx, settlement.StoreID)
		if err != nil {
			return err
		}
		// 3. 更新账单状态
		if err = ds.PointSettlementRepo().UpdateStatus(ctx, id, domain.PointSettlementStatusApproved, &user.ID); err != nil {
			return err
		}
		// 4. 更新账户： 加账户余额；加总收益
		if err = ds.StoreAccountRepo().AdjustAmount(ctx, settlement.StoreID, domain.StoreAccountAdjustments{
			BalanceDelta: settlement.Amount,
			TotalDelta:   settlement.Amount,
		}); err != nil {
			return err
		}
		// 5. 记录账户金额变更日志
		transaction := &domain.StoreAccountTransaction{
			StoreID: account.StoreID,
			No:      settlement.No,
			Amount:  settlement.Amount,
			After:   account.Balance.Add(settlement.Amount),
			Type:    domain.TransactionTypeSaleIncome,
		}
		if err = ds.StoreAccountRepo().RecordTransaction(ctx, transaction); err != nil {
			return err
		}
		return nil
	})
}

func (p *PointSettlementInteractor) UnApprove(ctx context.Context, id int) error {
	return p.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 获取账单并加锁
		settlement, err := ds.PointSettlementRepo().FindByIDForUpdate(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrPointSettlementNotExists)
			}
			return err
		}
		if settlement.Status != domain.PointSettlementStatusApproved {
			return domain.ParamsError(domain.ErrPointSettlementStatusInvalid)
		}
		// 2. 获取门店账户并加锁
		account, err := ds.StoreAccountRepo().FindByStoreForUpdate(ctx, settlement.StoreID)
		if err != nil {
			return err
		}
		// 3. 更新账单状态
		if err = ds.PointSettlementRepo().UpdateStatus(ctx, id, domain.PointSettlementStatusPending, nil); err != nil {
			return err
		}
		// 4. 更新账户： 减账户余额；减总收益
		if err = ds.StoreAccountRepo().AdjustAmount(ctx, settlement.StoreID, domain.StoreAccountAdjustments{
			BalanceDelta: settlement.Amount.Neg(),
			TotalDelta:   settlement.Amount.Neg(),
		}); err != nil {
			return err
		}
		// 5. 记录账户金额变更日志
		transaction := &domain.StoreAccountTransaction{
			StoreID: account.StoreID,
			No:      settlement.No,
			Amount:  settlement.Amount.Neg(),
			After:   account.Balance.Sub(settlement.Amount),
			Type:    domain.TransactionTypeSaleRevert,
		}
		if err = ds.StoreAccountRepo().RecordTransaction(ctx, transaction); err != nil {
			return err
		}
		return nil
	})
}

func (p *PointSettlementInteractor) ListDetails(ctx context.Context, id, storeID int) (res domain.PointSettlementDetails, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementInteractor.ListDetails")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	record, err := p.ds.PointSettlementRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ParamsError(domain.ErrPointSettlementNotExists)
		}
		return nil, err
	}
	if storeID > 0 && record.StoreID != storeID {
		return nil, domain.ParamsError(domain.ErrPointSettlementNotExists)
	}

	dayStart := util.DayStart(record.Date)
	dayEnd := util.DayEnd(record.Date)
	filter := &domain.OrderListFilter{
		Status:        domain.OrderStatusPaid,
		FinishedAtGte: &dayStart,
		FinishedAtLte: &dayEnd,
		PointsPaidGt0: true,
	}
	orders, total, err := p.ds.OrderRepo().GetOrders(ctx, &upagination.Pagination{
		Page: 1,
		Size: upagination.MaxSize,
	}, filter)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return res, nil
	}
	var orderIDs []int
	for _, order := range orders {
		item := &domain.PointSettlementDetail{
			OrderID:         order.ID,
			OrderNo:         order.No,
			OrderAmount:     order.TotalPrice,
			OrderFinishedAt: *order.FinishedAt,
			Amount:          order.PointsPaid,
			MemberName:      order.MemberName,
		}
		res = append(res, item)
		orderIDs = append(orderIDs, order.ID)
	}
	if len(orderIDs) == 0 {
		return res, nil
	}
	orderItemNamesMap, err := p.ds.OrderRepo().ListItemNamesByOrders(ctx, orderIDs)
	if err != nil {
		return nil, err
	}
	for _, item := range res {
		item.ProductInfo = orderItemNamesMap[item.OrderID]
	}
	return res, nil
}

func (p *PointSettlementInteractor) GetPointSettlementRange(ctx context.Context, params domain.PointSettlementSearchParams,
) (res domain.PointSettlementRange, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PointSettlementInteractor.GetPointSettlementRange")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	res, err = p.ds.PointSettlementRepo().GetPointSettlementRange(ctx, params)
	return
}
