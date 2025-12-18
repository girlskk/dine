package category

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (i *CategoryInteractor) CreateRoot(ctx context.Context, category *domain.Category) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "CategoryInteractor.CreateRoot")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return i.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// Check if the category name already exists
		exists, err := ds.CategoryRepo().Exists(ctx, domain.CategoryExistsParams{
			MerchantID: category.MerchantID,
			Name:       category.Name,
			IsRoot:     true,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ConflictError(domain.ErrCategoryNameExists)
		}

		// check tax rate id
		if category.TaxRateID != uuid.Nil {
			// @TODO
			// exists, err := ds.TaxRateRepo().Exists(ctx, domain.TaxRateExistsParams{
			// 	MerchantID: category.MerchantID,
			// 	ID:         category.TaxRateID,
			// })
			// if err != nil {
			// 	return err
			// }
			// if !exists {
			// 	return domain.NotFoundError(domain.ErrTaxRateNotExists)
			// }
		}

		if category.StallID != uuid.Nil {
			// @TODO
			// exists, err := ds.StallRepo().Exists(ctx, domain.StallExistsParams{
			// 	MerchantID: category.MerchantID,
			// 	ID:         category.StallID,
			// })
			// if err != nil {
			// 	return err
			// }
			// if !exists {
			// 	return domain.NotFoundError(domain.ErrStallNotExists)
			// }
		}

		// create root category
		err = ds.CategoryRepo().Create(ctx, category)
		if err != nil {
			return err
		}

		// create children categories
		if len(category.Childrens) > 0 {
			err = ds.CategoryRepo().CreateBulk(ctx, category.Childrens)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
