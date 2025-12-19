package productspec

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductSpecInteractor = (*ProductSpecInteractor)(nil)

type ProductSpecInteractor struct {
	DS domain.DataStore
}

func NewProductSpecInteractor(ds domain.DataStore) *ProductSpecInteractor {
	return &ProductSpecInteractor{
		DS: ds,
	}
}

func (i *ProductSpecInteractor) Create(ctx context.Context, spec *domain.ProductSpec) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductSpecInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证名称在当前门店下是否唯一
		exists, err := ds.ProductSpecRepo().Exists(ctx, domain.ProductSpecExistsParams{
			MerchantID: spec.MerchantID,
			Name:       spec.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrProductSpecNameExists
		}

		// 2. 创建商品规格
		err = ds.ProductSpecRepo().Create(ctx, spec)
		if err != nil {
			return err
		}

		return nil
	})
}

func (i *ProductSpecInteractor) Update(ctx context.Context, spec *domain.ProductSpec) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductSpecInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 验证规格存在
		existingSpec, err := ds.ProductSpecRepo().FindByID(ctx, spec.ID)
		if err != nil {
			return err
		}

		// 2. 验证更新后的名称在当前门店下是否唯一（排除自身）
		if spec.Name != existingSpec.Name {
			exists, err := ds.ProductSpecRepo().Exists(ctx, domain.ProductSpecExistsParams{
				MerchantID: existingSpec.MerchantID,
				Name:       spec.Name,
				ExcludeID:  spec.ID,
			})
			if err != nil {
				return err
			}
			if exists {
				return domain.ErrProductSpecNameExists
			}
		}

		// 3. 将请求数据合并到 existingSpec
		existingSpec.Name = spec.Name

		// 4. 执行更新操作
		err = ds.ProductSpecRepo().Update(ctx, existingSpec)
		if err != nil {
			return err
		}

		// 5. 更新返回的 spec 对象
		spec.UpdatedAt = existingSpec.UpdatedAt

		return nil
	})
}

func (i *ProductSpecInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductSpecInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 查找规格，验证规格存在
		spec, err := ds.ProductSpecRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductSpecNotExists)
			}
			return err
		}

		// 2. 如果规格下有关联商品，不能删除
		if spec.ProductCount > 0 {
			return domain.ErrProductSpecDeleteHasProducts
		}

		// 3. 删除规格
		err = ds.ProductSpecRepo().Delete(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (i *ProductSpecInteractor) PagedListBySearch(
	ctx context.Context,
	page *upagination.Pagination,
	params domain.ProductSpecSearchParams,
) (res *domain.ProductSpecSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductSpecInteractor.PagedListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.ProductSpecRepo().PagedListBySearch(ctx, page, params)
}
