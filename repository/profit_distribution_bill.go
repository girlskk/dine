package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	entprofitdistributionbill "gitlab.jiguang.dev/pos-dine/dine/ent/profitdistributionbill"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
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

func (repo *ProfitDistributionBillRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.ProfitDistributionBillSearchParams,
) (res *domain.ProfitDistributionBillSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProfitDistributionBillRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProfitDistributionBill.Query()

	// 必填条件：品牌商ID
	if params.MerchantID != uuid.Nil {
		query.Where(entprofitdistributionbill.MerchantID(params.MerchantID))
	}

	// 可选条件：门店ID列表
	if len(params.StoreIDs) > 0 {
		query.Where(entprofitdistributionbill.StoreIDIn(params.StoreIDs...))
	}

	// 可选条件：账单开始日期
	if params.BillStartDate != nil {
		query.Where(entprofitdistributionbill.BillDateGTE(util.DayStart(*params.BillStartDate)))
	}

	// 可选条件：账单结束日期
	if params.BillEndDate != nil {
		query.Where(entprofitdistributionbill.BillDateLTE(util.DayEnd(*params.BillEndDate)))
	}

	// 可选条件：分账状态
	if params.Status != "" {
		query.Where(entprofitdistributionbill.StatusEQ(params.Status))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}

	// 分页处理
	query = query.
		WithMerchant().
		WithStore().
		Offset(page.Offset()).
		Limit(page.Size)

	// 按创建时间倒序排列
	entBills, err := query.Order(ent.Desc(entprofitdistributionbill.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.ProfitDistributionBills, 0, len(entBills))
	for _, b := range entBills {
		items = append(items, convertProfitDistributionBillToDomain(b))
	}

	page.SetTotal(total)

	return &domain.ProfitDistributionBillSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

func (repo *ProfitDistributionBillRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.ProfitDistributionBill, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProfitDistributionBillRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	epdb, err := repo.Client.ProfitDistributionBill.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProfitDistributionBillNotExists)
		}
		return nil, err
	}
	res = convertProfitDistributionBillToDomain(epdb)
	return res, nil
}

func (repo *ProfitDistributionBillRepository) Update(ctx context.Context, bill *domain.ProfitDistributionBill) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProfitDistributionBillRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.ProfitDistributionBill.UpdateOneID(bill.ID).
		SetStatus(bill.Status).
		SetPaymentAmount(bill.PaymentAmount)

	updated, err := builder.Save(ctx)
	if err != nil {
		return err
	}

	bill.UpdatedAt = updated.UpdatedAt
	return nil
}

// ============================================
// 转换函数
// ============================================

func convertProfitDistributionBillToDomain(epdb *ent.ProfitDistributionBill) *domain.ProfitDistributionBill {
	if epdb == nil {
		return nil
	}

	bill := &domain.ProfitDistributionBill{
		ID:               epdb.ID,
		No:               epdb.No,
		MerchantID:       epdb.MerchantID,
		StoreID:          epdb.StoreID,
		ReceivableAmount: epdb.ReceivableAmount,
		PaymentAmount:    epdb.PaymentAmount,
		Status:           epdb.Status,
		BillDate:         epdb.BillDate,
		StartDate:        epdb.StartDate,
		EndDate:          epdb.EndDate,
		RuleSnapshot:     epdb.RuleSnapshot,
		CreatedAt:        epdb.CreatedAt,
		UpdatedAt:        epdb.UpdatedAt,
	}

	if epdb.Edges.Merchant != nil {
		bill.Merchant = &domain.MerchantSimple{
			ID:           epdb.Edges.Merchant.ID,
			MerchantName: epdb.Edges.Merchant.MerchantName,
		}
	}
	if epdb.Edges.Store != nil {
		bill.Store = &domain.StoreSimple{
			ID:        epdb.Edges.Store.ID,
			StoreName: epdb.Edges.Store.StoreName,
		}
	}
	return bill
}
