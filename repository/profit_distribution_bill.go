package repository

import (
	"context"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProfitDistributionBillRepository = (*ProfitDistributionBillRepository)(nil)

type ProfitDistributionBillRepository struct {
	Client *ent.Client
}

func NewProfitDistributionBillRepository(client *ent.Client) *ProfitDistributionBillRepository {
	return &ProfitDistributionBillRepository{
		Client: client,
	}
}

func (repo *ProfitDistributionBillRepository) CreateBulk(ctx context.Context, bills []*domain.ProfitDistributionBill) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProfitDistributionBillRepository.CreateBulk")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(bills) == 0 {
		return nil
	}

	util.PrettyJson(bills)

	builders := make([]*ent.ProfitDistributionBillCreate, 0, len(bills))
	for _, bill := range bills {
		builder := repo.Client.ProfitDistributionBill.Create().
			SetID(bill.ID).
			SetNo(bill.No).
			SetMerchantID(bill.MerchantID).
			SetStoreID(bill.StoreID).
			SetReceivableAmount(bill.ReceivableAmount).
			SetPaymentAmount(bill.PaymentAmount).
			SetStatus(bill.Status).
			SetBillDate(bill.BillDate).
			SetStartDate(bill.StartDate).
			SetEndDate(bill.EndDate).
			SetRuleSnapshot(bill.RuleSnapshot)

		builders = append(builders, builder)
	}

	_, err = repo.Client.ProfitDistributionBill.CreateBulk(builders...).Save(ctx)
	return err
}
