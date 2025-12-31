package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/profitdistributionrule"
	"gitlab.jiguang.dev/pos-dine/dine/ent/store"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProfitDistributionRuleRepository = (*ProfitDistributionRuleRepository)(nil)

type ProfitDistributionRuleRepository struct {
	Client *ent.Client
}

func NewProfitDistributionRuleRepository(client *ent.Client) *ProfitDistributionRuleRepository {
	return &ProfitDistributionRuleRepository{
		Client: client,
	}
}

func (repo *ProfitDistributionRuleRepository) FindByID(ctx context.Context, id uuid.UUID) (res *domain.ProfitDistributionRule, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProfitDistributionRuleRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	epdr, err := repo.Client.ProfitDistributionRule.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProfitDistributionRuleNotExists)
		}
		return nil, err
	}

	res = convertProfitDistributionRuleToDomain(epdr)
	return res, nil
}

func (repo *ProfitDistributionRuleRepository) GetDetail(ctx context.Context, id uuid.UUID) (res *domain.ProfitDistributionRule, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProfitDistributionRuleRepository.GetDetail")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	epdr, err := repo.Client.ProfitDistributionRule.Query().
		Where(profitdistributionrule.IDEQ(id)).
		WithStores().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrProfitDistributionRuleNotExists)
		}
		return nil, err
	}

	res = convertProfitDistributionRuleToDomain(epdr)
	return res, nil
}

func (repo *ProfitDistributionRuleRepository) Create(ctx context.Context, rule *domain.ProfitDistributionRule) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProfitDistributionRuleRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 创建分账方案
	ruleBuilder := repo.Client.ProfitDistributionRule.Create().
		SetID(rule.ID).
		SetMerchantID(rule.MerchantID).
		SetName(rule.Name).
		SetSplitRatio(rule.SplitRatio).
		SetBillingCycle(rule.BillingCycle).
		SetEffectiveDate(rule.EffectiveDate).
		SetExpiryDate(rule.ExpiryDate).
		SetBillGenerationDay(rule.BillGenerationDay).
		SetStatus(rule.Status).
		SetStoreCount(rule.StoreCount)

	// 设置关联门店（Many2Many）
	if len(rule.Stores) > 0 {
		storeIDs := make([]uuid.UUID, 0, len(rule.Stores))
		for _, s := range rule.Stores {
			storeIDs = append(storeIDs, s.ID)
		}
		ruleBuilder.AddStoreIDs(storeIDs...)
	}

	epdr, err := ruleBuilder.Save(ctx)
	if err != nil {
		return err
	}

	rule.CreatedAt = epdr.CreatedAt
	rule.UpdatedAt = epdr.UpdatedAt
	return nil
}

func (repo *ProfitDistributionRuleRepository) Update(ctx context.Context, rule *domain.ProfitDistributionRule) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProfitDistributionRuleRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 更新分账方案基本信息
	builder := repo.Client.ProfitDistributionRule.UpdateOneID(rule.ID).
		SetName(rule.Name).
		SetSplitRatio(rule.SplitRatio).
		SetBillingCycle(rule.BillingCycle).
		SetEffectiveDate(rule.EffectiveDate).
		SetExpiryDate(rule.ExpiryDate).
		SetBillGenerationDay(rule.BillGenerationDay).
		SetStatus(rule.Status).
		SetStoreCount(rule.StoreCount)

	// 更新关联门店（Many2Many）
	if len(rule.Stores) > 0 {
		storeIDs := make([]uuid.UUID, 0, len(rule.Stores))
		for _, s := range rule.Stores {
			storeIDs = append(storeIDs, s.ID)
		}
		builder = builder.ClearStores().AddStoreIDs(storeIDs...)
	} else {
		builder = builder.ClearStores()
	}

	_, err = builder.Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ProfitDistributionRuleRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProfitDistributionRuleRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.ProfitDistributionRule.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ProfitDistributionRuleRepository) Exists(ctx context.Context, params domain.ProfitDistributionRuleExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProfitDistributionRuleRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProfitDistributionRule.Query().
		Where(profitdistributionrule.MerchantID(params.MerchantID)).
		Where(profitdistributionrule.Name(params.Name))

	if params.ExcludeID != uuid.Nil {
		query.Where(profitdistributionrule.IDNEQ(params.ExcludeID))
	}

	return query.Exist(ctx)
}

func (repo *ProfitDistributionRuleRepository) CheckStoreBound(ctx context.Context, storeIDs []uuid.UUID, excludeRuleID uuid.UUID) (has bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProfitDistributionRuleRepository.CheckStoreBound")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(storeIDs) == 0 {
		return false, nil
	}

	query := repo.Client.ProfitDistributionRule.Query().
		Where(profitdistributionrule.HasStoresWith(store.IDIn(storeIDs...))).
		WithStores()

	if excludeRuleID != uuid.Nil {
		query.Where(profitdistributionrule.IDNEQ(excludeRuleID))
	}

	rules, err := query.All(ctx)
	if err != nil {
		return false, err
	}

	return len(rules) > 0, nil
}

func (repo *ProfitDistributionRuleRepository) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.ProfitDistributionRuleSearchParams,
) (res *domain.ProfitDistributionRuleSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "ProfitDistributionRuleRepository.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.ProfitDistributionRule.Query()

	if params.MerchantID != uuid.Nil {
		query.Where(profitdistributionrule.MerchantID(params.MerchantID))
	}

	if params.Name != "" {
		query.Where(profitdistributionrule.NameContains(params.Name))
	}

	if params.Status != "" {
		query.Where(profitdistributionrule.StatusEQ(params.Status))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}

	// 分页处理
	query = query.
		Offset(page.Offset()).
		Limit(page.Size)

	// 按创建时间倒序排列
	entRules, err := query.Order(ent.Desc(profitdistributionrule.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.ProfitDistributionRules, 0, len(entRules))
	for _, r := range entRules {
		items = append(items, convertProfitDistributionRuleToDomain(r))
	}

	page.SetTotal(total)

	return &domain.ProfitDistributionRuleSearchRes{
		Pagination: page,
		Items:      items,
	}, nil
}

// ============================================
// 转换函数
// ============================================

func convertProfitDistributionRuleToDomain(epdr *ent.ProfitDistributionRule) *domain.ProfitDistributionRule {
	if epdr == nil {
		return nil
	}

	r := &domain.ProfitDistributionRule{
		ID:                epdr.ID,
		MerchantID:        epdr.MerchantID,
		Name:              epdr.Name,
		SplitRatio:        epdr.SplitRatio,
		BillingCycle:      epdr.BillingCycle,
		EffectiveDate:     epdr.EffectiveDate,
		ExpiryDate:        epdr.ExpiryDate,
		BillGenerationDay: epdr.BillGenerationDay,
		Status:            epdr.Status,
		StoreCount:        epdr.StoreCount,
		CreatedAt:         epdr.CreatedAt,
		UpdatedAt:         epdr.UpdatedAt,
	}

	// 转换门店列表
	if len(epdr.Edges.Stores) > 0 {
		r.Stores = lo.Map(epdr.Edges.Stores, func(store *ent.Store, _ int) *domain.StoreSimple {
			return &domain.StoreSimple{
				ID:        store.ID,
				StoreName: store.StoreName,
			}
		})
	}

	return r
}
