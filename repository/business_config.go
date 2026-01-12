package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/businessconfig"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.BusinessConfigRepository = (*BusinessConfigRepository)(nil)

type BusinessConfigRepository struct {
	Client *ent.Client
}

func NewBusinessConfigRepository(client *ent.Client) *BusinessConfigRepository {
	return &BusinessConfigRepository{
		Client: client,
	}
}

func (repo *BusinessConfigRepository) ListBySearch(
	ctx context.Context,
	params domain.BusinessConfigSearchParams,
) (res *domain.BusinessConfigSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "BusinessConfigRepository.ListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.Client.BusinessConfig.Query()

	if params.MerchantID != uuid.Nil {
		query.Where(businessconfig.MerchantID(params.MerchantID))
	}
	if params.StoreID != uuid.Nil {
		query.Where(businessconfig.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(businessconfig.NameContains(params.Name))
	}
	// 按创建时间倒序排列
	list, err := query.Order(ent.Asc(businessconfig.FieldSort)).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make(domain.BusinessConfigs, 0, len(list))
	for _, m := range list {
		items = append(items, convertBusinessConfigToDomain(m))
	}
	return &domain.BusinessConfigSearchRes{
		Items: items,
	}, nil
}

// ============================================
// 转换函数
// ============================================

func convertBusinessConfigToDomain(pm *ent.BusinessConfig) *domain.BusinessConfig {
	if pm == nil {
		return nil
	}
	m := &domain.BusinessConfig{
		ID:             pm.ID,
		SourceConfigID: pm.SourceConfigID,
		MerchantID:     pm.MerchantID,
		StoreID:        pm.StoreID,
		Group:          pm.Group,
		Name:           pm.Name,
		ConfigType:     pm.ConfigType,
		Key:            pm.Key,
		Value:          pm.Value,
		Sort:           pm.Sort,
		Tip:            pm.Tip,
		IsDefault:      pm.IsDefault,
		Status:         pm.Status,
		CreatedAt:      pm.CreatedAt,
		UpdatedAt:      pm.UpdatedAt,
	}
	return m
}
