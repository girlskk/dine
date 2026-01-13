package repository

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
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

	query := repo.Client.BusinessConfig.Query().Where(businessconfig.Or(
		businessconfig.MerchantID(params.MerchantID),
		businessconfig.And(businessconfig.MerchantIDIsNil(), businessconfig.IsDefaultEQ(true)),
	))
	if params.StoreID != uuid.Nil {
		query.Where(businessconfig.StoreID(params.StoreID))
	}
	if params.Name != "" {
		query.Where(businessconfig.NameContains(params.Name))
	}
	if params.Group != "" {
		query.Where(businessconfig.GroupEQ(params.Group))
	}
	if len(params.Ids) > 0 {
		query.Where(businessconfig.IDIn(params.Ids...))
	}
	// 按创建时间倒序排列
	list, err := query.Order(ent.Asc(businessconfig.FieldSort)).All(ctx)
	if err != nil {
		return nil, err
	}
	items := lo.MapToSlice(mergeConfigs(list), func(key string, m *ent.BusinessConfig) *domain.BusinessConfig {
		return convertBusinessConfigToDomain(m)
	})
	sort.Slice(items, func(i, j int) bool {
		return items[i].Sort < items[j].Sort
	})
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

// 合并配置：门店配置覆盖默认配置
func mergeConfigs(configs []*ent.BusinessConfig) map[string]*ent.BusinessConfig {
	result := make(map[string]*ent.BusinessConfig)
	defaultConfigs := make(map[string]*ent.BusinessConfig)
	storeConfigs := make(map[string]*ent.BusinessConfig)
	// 分离默认配置和门店配置
	for _, config := range configs {
		key := fmt.Sprintf("%s.%s", config.Group, config.Key)
		if config.MerchantID != uuid.Nil || config.StoreID != uuid.Nil {
			storeConfigs[key] = config
		} else if config.IsDefault {
			defaultConfigs[key] = config
		}
	}
	for key, defaultConfig := range defaultConfigs {
		result[key] = defaultConfig
	}
	for key, storeConfig := range storeConfigs {
		result[key] = storeConfig
	}

	return result
}

func (repo *BusinessConfigRepository) UpsertConfig(ctx context.Context, configs []*domain.BusinessConfig) error {
	now := time.Now()
	builders := make([]*ent.BusinessConfigCreate, 0, len(configs))

	for _, c := range configs {
		builder := repo.Client.BusinessConfig.
			Create().
			SetID(c.ID).
			SetMerchantID(c.MerchantID).
			SetGroup(c.Group).
			SetName(c.Name).
			SetConfigType(c.ConfigType).
			SetKey(c.Key).
			SetValue(c.Value).
			SetSort(c.Sort).
			SetTip(c.Tip).
			SetIsDefault(c.IsDefault).
			SetStatus(c.Status).
			SetCreatedAt(now).
			SetUpdatedAt(now).
			SetDeletedAt(0)
		if c.SourceConfigID != uuid.Nil {
			builder = builder.SetSourceConfigID(c.SourceConfigID)
		}
		if c.StoreID != uuid.Nil {
			builder = builder.SetStoreID(c.StoreID)
		}
		if c.MerchantID != uuid.Nil {
			builder = builder.SetMerchantID(c.MerchantID)
		}
		builders = append(builders, builder)
	}
	err := repo.Client.BusinessConfig.
		CreateBulk(builders...).
		OnConflict().
		UpdateValue().
		UpdateUpdatedAt().
		Exec(ctx)
	return err
}
