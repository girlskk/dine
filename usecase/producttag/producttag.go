package producttag

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductTagInteractor = (*ProductTagInteractor)(nil)

type ProductTagInteractor struct {
	DS domain.DataStore
}

func NewProductTagInteractor(ds domain.DataStore) *ProductTagInteractor {
	return &ProductTagInteractor{
		DS: ds,
	}
}

func (i *ProductTagInteractor) Create(ctx context.Context, tag *domain.ProductTag) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductTagInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证名称在当前品牌商/门店下是否唯一
		exists, err := ds.ProductTagRepo().Exists(ctx, domain.ProductTagExistsParams{
			MerchantID: tag.MerchantID,
			StoreID:    tag.StoreID,
			Name:       tag.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrProductTagNameExists
		}

		// 2. 创建商品标签
		err = ds.ProductTagRepo().Create(ctx, tag)
		if err != nil {
			return err
		}

		return nil
	})
}

func (i *ProductTagInteractor) Update(ctx context.Context, tag *domain.ProductTag, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductTagInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证标签存在
		existingTag, err := ds.ProductTagRepo().FindByID(ctx, tag.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductTagNotExists)
			}
			return err
		}

		if err := verifyProductTagOwnership(user, existingTag); err != nil {
			return err
		}

		// 2. 验证更新后的名称在当前品牌商/门店下是否唯一（排除自身）
		if tag.Name != existingTag.Name {
			exists, err := ds.ProductTagRepo().Exists(ctx, domain.ProductTagExistsParams{
				MerchantID: existingTag.MerchantID,
				StoreID:    existingTag.StoreID,
				Name:       tag.Name,
				ExcludeID:  tag.ID,
			})
			if err != nil {
				return err
			}
			if exists {
				return domain.ErrProductTagNameExists
			}
		}

		// 3. 将请求数据合并到 existingTag
		existingTag.Name = tag.Name

		// 4. 执行更新操作
		err = ds.ProductTagRepo().Update(ctx, existingTag)
		if err != nil {
			return err
		}

		// 5. 更新返回的 tag 对象
		tag.UpdatedAt = existingTag.UpdatedAt

		return nil
	})
}

func (i *ProductTagInteractor) Delete(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductTagInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 查找标签，验证标签存在
		tag, err := ds.ProductTagRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductTagNotExists)
			}
			return err
		}

		if err := verifyProductTagOwnership(user, tag); err != nil {
			return err
		}

		// 2. 如果标签下有关联商品，不能删除
		if tag.ProductCount > 0 {
			return domain.ErrProductTagDeleteHasProducts
		}

		// 3. 删除标签
		err = ds.ProductTagRepo().Delete(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (i *ProductTagInteractor) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.ProductTagSearchParams,
) (res *domain.ProductTagSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductTagInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.ProductTagRepo().PagedListBySearch(ctx, page, params)
}

func verifyProductTagOwnership(user domain.User, tag *domain.ProductTag) error {
	if user.GetMerchantID() != tag.MerchantID || user.GetStoreID() != tag.StoreID {
		return domain.ParamsError(domain.ErrProductTagNotExists)
	}
	return nil
}
