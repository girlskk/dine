package businessconfig

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.BusinessConfigInteractor = (*BusinessConfigInteractor)(nil)

type BusinessConfigInteractor struct {
	DS domain.DataStore
}

func NewBusinessConfigInteractor(ds domain.DataStore) *BusinessConfigInteractor {
	return &BusinessConfigInteractor{
		DS: ds,
	}
}

func (i *BusinessConfigInteractor) ListBySearch(
	ctx context.Context,
	params domain.BusinessConfigSearchParams,
) (res *domain.BusinessConfigSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BusinessConfigInteractor.ListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.BusinessConfigRepo().ListBySearch(ctx, params)
}

func (i *BusinessConfigInteractor) UpsertConfig(ctx context.Context, configs []*domain.BusinessConfig, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BusinessConfigInteractor.UpsertConfig")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		return ds.BusinessConfigRepo().UpsertConfig(ctx, configs)
	})
}
func (i *BusinessConfigInteractor) Distribute(ctx context.Context, ids, storeIDs []uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BusinessConfigInteractor.UpsertConfig")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	configList, err := i.DS.BusinessConfigRepo().ListBySearch(ctx, domain.BusinessConfigSearchParams{
		Ids: ids,
	})
	if err != nil {
		return err
	}
	if len(configList.Items) == 0 {
		return nil
	}
	configs := make([]*domain.BusinessConfig, 0, len(configList.Items))
	for _, config := range configList.Items {
		for _, storeID := range storeIDs {
			config.MerchantID = user.GetMerchantID()
			config.StoreID = storeID
			config.IsDefault = false
			if config.SourceConfigID == uuid.Nil {
				config.SourceConfigID = config.ID
			}
			configs = append(configs, config)
		}
	}
	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		return ds.BusinessConfigRepo().UpsertConfig(ctx, configs)
	})
}
