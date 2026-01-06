package device

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.DeviceInteractor = (*DeviceInteractor)(nil)

// DeviceInteractor implements device use cases.
type DeviceInteractor struct {
	ds domain.DataStore
}

func NewDeviceInteractor(ds domain.DataStore) *DeviceInteractor {
	return &DeviceInteractor{ds: ds}
}

func (interactor *DeviceInteractor) DeviceSimpleUpdate(ctx context.Context,
	updateField domain.DeviceSimpleUpdateType,
	device *domain.Device,
) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DeviceInteractor.DeviceSimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	if device == nil {
		return fmt.Errorf("device is nil")
	}
	oldDevice, err := interactor.ds.DeviceRepo().FindByID(ctx, device.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrDeviceNotExists)
		}
		return fmt.Errorf("failed to fetch device: %w", err)
	}

	switch updateField {
	case domain.DeviceSimpleUpdateTypeEnabled:
		if oldDevice.Enabled == device.Enabled {
			return
		}
		oldDevice.Enabled = device.Enabled
	default:
		return domain.ParamsError(errors.New("unsupported update field"))
	}

	err = interactor.ds.DeviceRepo().Update(ctx, oldDevice)
	if err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}
	return
}

func (interactor *DeviceInteractor) Create(ctx context.Context, domainDevice *domain.Device) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DeviceInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainDevice == nil {
		return fmt.Errorf("device is nil")
	}

	if err = interactor.checkExists(ctx, domainDevice); err != nil {
		return err
	}
	domainDevice.ID = uuid.New()
	err = interactor.ds.DeviceRepo().Create(ctx, domainDevice)
	if err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}
	return
}

func (interactor *DeviceInteractor) Update(ctx context.Context, domainDevice *domain.Device) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DeviceInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainDevice == nil {
		return fmt.Errorf("device is nil")
	}
	if err = interactor.checkExists(ctx, domainDevice); err != nil {
		return err
	}
	err = interactor.ds.DeviceRepo().Update(ctx, domainDevice)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrDeviceNotExists)
		}
		return fmt.Errorf("failed to update device: %w", err)
	}
	return
}

func (interactor *DeviceInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DeviceInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()
	err = interactor.ds.DeviceRepo().Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}
	return
}

func (interactor *DeviceInteractor) GetDevice(ctx context.Context, id uuid.UUID) (domainDevice *domain.Device, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DeviceInteractor.GetDevice")
	defer func() { util.SpanErrFinish(span, err) }()
	domainDevice, err = interactor.ds.DeviceRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrDeviceNotExists)
			return
		}
		err = fmt.Errorf("failed to fetch device: %w", err)
		return
	}
	return
}

func (interactor *DeviceInteractor) GetDevices(ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.DeviceListFilter,
	orderBys ...domain.DeviceOrderBy,
) (domainDevices []*domain.Device, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DeviceInteractor.GetDevices")
	defer func() { util.SpanErrFinish(span, err) }()
	if filter == nil {
		err = domain.ParamsError(errors.New("filter is required"))
	}
	domainDevices, total, err = interactor.ds.DeviceRepo().GetDevices(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get devices: %w", err)
		return
	}
	return
}

func (interactor *DeviceInteractor) checkExists(ctx context.Context, domainDevice *domain.Device) (err error) {
	exists, err := interactor.ds.DeviceRepo().Exists(ctx, domain.DeviceExistsParams{
		MerchantID: domainDevice.MerchantID,
		StoreID:    domainDevice.StoreID,
		Name:       domainDevice.Name,
		ExcludeID:  domainDevice.ID,
	})
	if err != nil {
		return err
	}
	if exists {
		return domain.ParamsError(domain.ErrDeviceNameExists)
	}

	if domainDevice.DeviceCode != "" {
		exists, err = interactor.ds.DeviceRepo().Exists(ctx, domain.DeviceExistsParams{
			MerchantID: domainDevice.MerchantID,
			StoreID:    domainDevice.StoreID,
			DeviceCode: domainDevice.DeviceCode,
			ExcludeID:  domainDevice.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrDeviceCodeExists)
		}
	}
	return nil
}
