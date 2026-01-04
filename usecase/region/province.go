package region

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ProvinceInteractor = (*ProvinceInteractor)(nil)

type ProvinceInteractor struct {
	ds domain.DataStore
}

func NewProvinceInteractor(ds domain.DataStore) *ProvinceInteractor {
	return &ProvinceInteractor{ds: ds}
}

func (interactor *ProvinceInteractor) GetProvinces(ctx context.Context, countryID uuid.UUID) (provinceList []*domain.Province, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProvinceInteractor.GetAllProvinces")
	defer func() { util.SpanErrFinish(span, err) }()

	provinceList, err = interactor.ds.ProvinceRepo().GetAll(ctx, countryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provinces: %w", err)
	}
	return
}

func (interactor *ProvinceInteractor) GetProvince(ctx context.Context, id uuid.UUID) (province *domain.Province, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProvinceInteractor.GetProvince")
	defer func() { util.SpanErrFinish(span, err) }()

	province, err = interactor.ds.ProvinceRepo().FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get province: %w", err)
	}
	return
}

func (interactor *ProvinceInteractor) CreateProvince(ctx context.Context, province *domain.Province) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProvinceInteractor.CreateProvince")
	defer func() { util.SpanErrFinish(span, err) }()

	province.ID = uuid.New()
	err = interactor.ds.ProvinceRepo().Create(ctx, province)
	if err != nil {
		return fmt.Errorf("failed to create province: %w", err)
	}
	return nil
}

func (interactor *ProvinceInteractor) UpdateProvince(ctx context.Context, province *domain.Province) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProvinceInteractor.UpdateProvince")
	defer func() { util.SpanErrFinish(span, err) }()

	if err = interactor.ds.ProvinceRepo().Update(ctx, province); err != nil {
		return fmt.Errorf("failed to update province: %w", err)
	}
	return nil
}

func (interactor *ProvinceInteractor) DeleteProvince(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProvinceInteractor.DeleteProvince")
	defer func() { util.SpanErrFinish(span, err) }()

	if err = interactor.ds.ProvinceRepo().Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete province: %w", err)
	}
	return nil
}

func (interactor *ProvinceInteractor) GetProvincesByFilter(ctx context.Context, filter *domain.ProvinceListFilter) (provinceList []*domain.Province, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "ProvinceInteractor.GetProvincesByFilter")
	defer func() { util.SpanErrFinish(span, err) }()

	provinceList, err = interactor.ds.ProvinceRepo().GetByFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get provinces by filter: %w", err)
	}
	return
}
