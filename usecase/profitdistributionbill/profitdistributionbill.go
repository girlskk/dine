package profitdistributionbill

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProfitDistributionBillInteractor = (*ProfitDistributionBillInteractor)(nil)

type ProfitDistributionBillInteractor struct {
	DS  domain.DataStore
	Seq domain.DailySequence
}

func NewProfitDistributionBillInteractor(ds domain.DataStore, seq domain.DailySequence) *ProfitDistributionBillInteractor {
	return &ProfitDistributionBillInteractor{
		DS:  ds,
		Seq: seq,
	}
}

func (i *ProfitDistributionBillInteractor) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.ProfitDistributionBillSearchParams,
) (res *domain.ProfitDistributionBillSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProfitDistributionBillInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return i.DS.ProfitDistributionBillRepo().PagedListBySearch(ctx, page, params)
}

func (i *ProfitDistributionBillInteractor) Pay(ctx context.Context, id uuid.UUID, paymentAmount decimal.Decimal, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProfitDistributionBillInteractor.Pay")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证账单存在
		bill, err := ds.ProfitDistributionBillRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProfitDistributionBillNotExists)
			}
			return err
		}

		// 2. 验证账单是否属于当前品牌商
		if bill.MerchantID != user.GetMerchantID() {
			return domain.ParamsError(domain.ErrProfitDistributionBillNotExists)
		}

		// 3. 验证账单状态：只有未打款状态的账单才能打款
		if bill.Status != domain.ProfitDistributionBillStatusUnpaid {
			return domain.ParamsError(domain.ErrProfitDistributionBillStatusInvalid)
		}

		// 4. 更新账单状态为已打款，并设置打款金额
		bill.Status = domain.ProfitDistributionBillStatusPaid
		bill.PaymentAmount = paymentAmount

		// 5. 保存更新
		return ds.ProfitDistributionBillRepo().Update(ctx, bill)
	})
}
