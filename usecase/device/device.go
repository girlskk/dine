package device

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.DeviceInteractor = (*DeviceInteractor)(nil)

// DeviceInteractor implements device use cases.
type DeviceInteractor struct {
	DS domain.DataStore
}

func NewDeviceInteractor(ds domain.DataStore) *DeviceInteractor {
	return &DeviceInteractor{DS: ds}
}

func (interactor *DeviceInteractor) DeviceSimpleUpdate(ctx context.Context,
	updateField domain.DeviceSimpleUpdateType,
	device *domain.Device,
	user domain.User,
) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DeviceInteractor.DeviceSimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()
	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		oldDevice, err := ds.DeviceRepo().FindByID(ctx, device.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrDeviceNotExists
			}
			return err
		}

		if err = verifyDeviceOwnership(user, oldDevice); err != nil {
			return err
		}

		switch updateField {
		case domain.DeviceSimpleUpdateTypeEnabled:
			if oldDevice.Enabled == device.Enabled {
				return nil
			}
			oldDevice.Enabled = device.Enabled
		default:
			return fmt.Errorf("unsupported simple update field: %s", updateField)
		}

		err = ds.DeviceRepo().Update(ctx, oldDevice)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (interactor *DeviceInteractor) Create(ctx context.Context, domainDevice *domain.Device, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DeviceInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if err = verifyDeviceOwnership(user, domainDevice); err != nil {
		return err
	}

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.DeviceRepo().Exists(ctx, domain.DeviceExistsParams{
			MerchantID: domainDevice.MerchantID,
			StoreID:    domainDevice.StoreID,
			Name:       domainDevice.Name,
			ExcludeID:  domainDevice.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrDeviceNameExists
		}

		if domainDevice.DeviceCode != "" {
			exists, err = ds.DeviceRepo().Exists(ctx, domain.DeviceExistsParams{
				MerchantID: domainDevice.MerchantID,
				StoreID:    domainDevice.StoreID,
				DeviceCode: domainDevice.DeviceCode,
				ExcludeID:  domainDevice.ID,
			})
			if err != nil {
				return err
			}
			if exists {
				return domain.ErrDeviceCodeExists
			}
		}
		domainDevice.ID = uuid.New()
		err = ds.DeviceRepo().Create(ctx, domainDevice)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (interactor *DeviceInteractor) Update(ctx context.Context, domainDevice *domain.Device, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DeviceInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainDevice == nil {
		return fmt.Errorf("device is nil")
	}
	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		old, err := ds.DeviceRepo().FindByID(ctx, domainDevice.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrDeviceNotExists
			}
			return err
		}
		if err = verifyDeviceOwnership(user, old); err != nil {
			return err
		}

		exists, err := ds.DeviceRepo().Exists(ctx, domain.DeviceExistsParams{
			MerchantID: domainDevice.MerchantID,
			StoreID:    domainDevice.StoreID,
			Name:       domainDevice.Name,
			ExcludeID:  domainDevice.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrDeviceNameExists
		}

		if domainDevice.DeviceCode != "" {
			exists, err = ds.DeviceRepo().Exists(ctx, domain.DeviceExistsParams{
				MerchantID: domainDevice.MerchantID,
				StoreID:    domainDevice.StoreID,
				DeviceCode: domainDevice.DeviceCode,
				ExcludeID:  domainDevice.ID,
			})
			if err != nil {
				return err
			}
			if exists {
				return domain.ErrDeviceCodeExists
			}
		}
		err = ds.DeviceRepo().Update(ctx, domainDevice)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (interactor *DeviceInteractor) Delete(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DeviceInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()
	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		device, err := ds.DeviceRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrDeviceNotExists
			}
			return err
		}
		if err = verifyDeviceOwnership(user, device); err != nil {
			return err
		}
		err = ds.DeviceRepo().Delete(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (interactor *DeviceInteractor) GetDevice(ctx context.Context, id uuid.UUID, user domain.User) (domainDevice *domain.Device, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DeviceInteractor.GetDevice")
	defer func() { util.SpanErrFinish(span, err) }()
	domainDevice, err = interactor.DS.DeviceRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ErrDeviceNotExists
		}
		return nil, err
	}
	if err = verifyDeviceOwnership(user, domainDevice); err != nil {
		return nil, err
	}
	return domainDevice, nil
}

func (interactor *DeviceInteractor) GetDevices(ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.DeviceListFilter,
	orderBys ...domain.DeviceOrderBy,
) (domainDevices []*domain.Device, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DeviceInteractor.GetDevices")
	defer func() { util.SpanErrFinish(span, err) }()
	domainDevices, total, err = interactor.DS.DeviceRepo().GetDevices(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get devices: %w", err)
		return
	}
	return domainDevices, total, nil
}

func verifyDeviceOwnership(user domain.User, device *domain.Device) error {
	switch user.GetUserType() {
	case domain.UserTypeAdmin:
	case domain.UserTypeBackend:
		if !domain.VerifyOwnerMerchant(user, device.MerchantID) {
			return domain.ErrDeviceNotExists
		}
	case domain.UserTypeStore:
		if !domain.VerifyOwnerShip(user, device.MerchantID, device.StoreID) {
			return domain.ErrDeviceNotExists
		}
	}
	return nil
}
