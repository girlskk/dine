package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StoreInteractor = (*StoreInteractor)(nil)

type StoreInteractor struct {
	ds domain.DataStore
}

func NewStoreInteractor(dataStore domain.DataStore) *StoreInteractor {
	return &StoreInteractor{
		ds: dataStore,
	}
}

func (i *StoreInteractor) Create(ctx context.Context, store *domain.Store, user *domain.BackendUser) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 校验门店名称唯一性
		exists, err := ds.StoreRepo().Exists(ctx, domain.StoreExistsParams{
			Name: store.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrStoreNameExists)
		}
		// 检查用户名是否存在
		exists, err = ds.BackendUserRepo().Exists(ctx, user.Username)
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrUserExists)
		}

		// 创建门店
		if err := ds.StoreRepo().Create(ctx, store); err != nil {
			return err
		}
		user.StoreID = store.ID
		// 创建门店管理员
		if err := ds.BackendUserRepo().Create(ctx, user); err != nil {
			return err
		}
		// 创建门店账户
		return ds.StoreAccountRepo().Create(ctx, &domain.StoreAccount{
			StoreID: store.ID,
		})
	})
}

func (i *StoreInteractor) Update(ctx context.Context, store *domain.Store, user *domain.BackendUser) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		existStore, err := ds.StoreRepo().Find(ctx, store.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
		}
		// 名称唯一性校验（排除自身）
		if existStore.Name != store.Name {
			exists, err := ds.StoreRepo().Exists(ctx, domain.StoreExistsParams{Name: store.Name})
			if err != nil {
				return err
			}
			if exists {
				return domain.ParamsError(domain.ErrStoreNameExists)
			}
		}
		err = ds.StoreRepo().Update(ctx, store)
		if err != nil {
			return err
		}

		// 更新用户密码
		if user.HashedPassword != "" {
			if err := ds.BackendUserRepo().Update(ctx, user); err != nil {
				return err
			}
		}
		return nil
	})
}

// 门店后台更新
func (i *StoreInteractor) UpdateByStore(ctx context.Context, store *domain.Store, user *domain.BackendUser) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreInteractor.UpdateByStore")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.ds.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		existStore, err := ds.StoreRepo().Find(ctx, store.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(err)
			}
		}
		existStore.Type = store.Type
		existStore.Info = store.Info
		existStore.Finance = store.Finance

		existUser, err := ds.BackendUserRepo().FindByStoreID(ctx, store.ID)
		if err != nil {
			return err
		}

		err = ds.StoreRepo().Update(ctx, existStore)
		if err != nil {
			return err
		}
		if existUser.HashedPassword != user.HashedPassword {
			if err := ds.BackendUserRepo().Update(ctx, user); err != nil {
				return err
			}
		}
		return nil
	})
}

func (i *StoreInteractor) PagedListBySearch(ctx context.Context, page *upagination.Pagination,
	params domain.StoreSearchParams,
) (res *domain.StoreSearchRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	return i.ds.StoreRepo().PagedListBySearch(ctx, page, params)
}

func (i *StoreInteractor) GetDetail(ctx context.Context, id int) (res *domain.Store, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreInteractor.GetDetail")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	store, err := i.ds.StoreRepo().GetDetail(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ParamsError(domain.ErrStoreNotExists)
		}
		return nil, err
	}
	return store, nil
}
