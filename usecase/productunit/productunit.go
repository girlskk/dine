package productunit

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProductUnitInteractor = (*ProductUnitInteractor)(nil)

type ProductUnitInteractor struct {
	DS domain.DataStore
}

func NewProductUnitInteractor(ds domain.DataStore) *ProductUnitInteractor {
	return &ProductUnitInteractor{
		DS: ds,
	}
}

func (i *ProductUnitInteractor) Create(ctx context.Context, unit *domain.ProductUnit) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductUnitInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证名称在当前品牌商下是否唯一
		exists, err := ds.ProductUnitRepo().Exists(ctx, domain.ProductUnitExistsParams{
			MerchantID: unit.MerchantID,
			Name:       unit.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrProductUnitNameExists
		}

		// 3. 创建商品单位
		err = ds.ProductUnitRepo().Create(ctx, unit)
		if err != nil {
			return err
		}

		return nil
	})
}

func (i *ProductUnitInteractor) PagedListBySearch(ctx context.Context,
	page *upagination.Pagination, params domain.ProductUnitSearchParams,
) (res *domain.ProductUnitSearchRes, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductUnitInteractor.ListBySearch")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.ProductUnitRepo().PagedListBySearch(ctx, page, params)
}

func (i *ProductUnitInteractor) Update(ctx context.Context, unit *domain.ProductUnit) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductUnitInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 验证单位存在
		existingUnit, err := ds.ProductUnitRepo().FindByID(ctx, unit.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductUnitNotExists)
			}
			return err
		}

		// 验证更新后的名称在当前品牌商下是否唯一（排除自身）
		if unit.Name != existingUnit.Name {
			exists, err := ds.ProductUnitRepo().Exists(ctx, domain.ProductUnitExistsParams{
				MerchantID: existingUnit.MerchantID,
				Name:       unit.Name,
				ExcludeID:  unit.ID,
			})
			if err != nil {
				return err
			}
			if exists {
				return domain.ErrProductUnitNameExists
			}
		}

		existingUnit.Name = unit.Name
		existingUnit.Type = unit.Type

		err = ds.ProductUnitRepo().Update(ctx, existingUnit)
		if err != nil {
			return err
		}

		return nil
	})
}

func (i *ProductUnitInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProductUnitInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// 1. 查找单位，验证单位存在
		unit, err := ds.ProductUnitRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrProductUnitNotExists)
			}
			return err
		}

		// 2. 如果单位下有关联商品，不能删除
		if unit.ProductCount > 0 {
			return domain.ErrProductUnitDeleteHasProducts
		}

		// 3. 删除单位
		err = ds.ProductUnitRepo().Delete(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}
