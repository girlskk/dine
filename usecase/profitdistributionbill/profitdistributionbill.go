package profitdistributionbill

import (
	"context"
	"fmt"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProfitDistributionBillInteractor = (*ProfitDistributionBillInteractor)(nil)

type ProfitDistributionBillInteractor struct {
	DS domain.DataStore
}

func NewProfitDistributionBillInteractor(ds domain.DataStore) *ProfitDistributionBillInteractor {
	return &ProfitDistributionBillInteractor{
		DS: ds,
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

	return nil, nil
}

func (i *ProfitDistributionBillInteractor) GenerateProfitDistributionBills(ctx context.Context) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProfitDistributionBillInteractor.GenerateProfitDistributionBills")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	fmt.Println("开始生成分账账单")

	return nil
}
