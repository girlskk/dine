package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/device"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

// DeviceRepository implements Device CRUD and pagination.
type DeviceRepository struct {
	Client *ent.Client
}

var _ domain.DeviceRepository = (*DeviceRepository)(nil)

func NewDeviceRepository(client *ent.Client) *DeviceRepository {
	return &DeviceRepository{Client: client}
}

func (repo *DeviceRepository) FindByID(ctx context.Context, id uuid.UUID) (domainDevice *domain.Device, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "DeviceRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	es, err := repo.Client.Device.Query().
		Where(device.ID(id)).
		WithStore().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrDeviceNotExists)
			return
		}
		return
	}
	domainDevice = convertDeviceToDomain(es)
	return
}

func (repo *DeviceRepository) Create(ctx context.Context, domainDevice *domain.Device) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "DeviceRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainDevice == nil {
		return fmt.Errorf("device is nil")
	}

	builder := repo.Client.Device.Create().SetID(domainDevice.ID).
		SetName(domainDevice.Name).
		SetDeviceType(domainDevice.DeviceType).
		SetMerchantID(domainDevice.MerchantID).
		SetStoreID(domainDevice.StoreID).
		SetStatus(domainDevice.Status).
		SetLocation(domainDevice.Location).
		SetEnabled(domainDevice.Enabled).
		SetSortOrder(domainDevice.SortOrder).
		SetDeviceCode(domainDevice.DeviceCode).
		SetDeviceBrand(domainDevice.DeviceBrand).
		SetDeviceModel(domainDevice.DeviceModel)
	switch domainDevice.DeviceType {
	case domain.DeviceTypePrinter:
		builder = builder.SetIP(domainDevice.IP)
		builder = builder.SetPaperSize(domainDevice.PaperSize)
		builder = builder.SetStallID(domainDevice.StallID)
		builder = builder.SetOrderChannels(domainDevice.OrderChannels)
		builder = builder.SetDiningWays(domainDevice.DiningWays)
		builder = builder.SetDeviceStallPrintType(domainDevice.DeviceStallPrintType)
		builder = builder.SetDeviceStallReceiptType(domainDevice.DeviceStallReceiptType)
	case domain.DeviceTypeCashier:
		builder = builder.SetOpenCashDrawer(domainDevice.OpenCashDrawer)
	}
	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create device: %w", err)
		return
	}
	domainDevice.ID = created.ID
	domainDevice.CreatedAt = created.CreatedAt
	domainDevice.UpdatedAt = created.UpdatedAt
	return
}

func (repo *DeviceRepository) Update(ctx context.Context, domainDevice *domain.Device) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "DeviceRepository.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainDevice == nil {
		return fmt.Errorf("device is nil")
	}

	builder := repo.Client.Device.UpdateOneID(domainDevice.ID).
		SetName(domainDevice.Name).
		SetDeviceType(domainDevice.DeviceType).
		SetStatus(domainDevice.Status).
		SetLocation(domainDevice.Location).
		SetEnabled(domainDevice.Enabled).
		SetSortOrder(domainDevice.SortOrder).
		SetDeviceCode(domainDevice.DeviceCode).
		SetDeviceBrand(domainDevice.DeviceBrand).
		SetDeviceModel(domainDevice.DeviceModel)
	switch domainDevice.DeviceType {
	case domain.DeviceTypePrinter:
		builder = builder.SetIP(domainDevice.IP)
		builder = builder.SetPaperSize(domainDevice.PaperSize)
		builder = builder.SetStallID(domainDevice.StallID)
		builder = builder.SetOrderChannels(domainDevice.OrderChannels)
		builder = builder.SetDiningWays(domainDevice.DiningWays)
		builder = builder.SetDeviceStallPrintType(domainDevice.DeviceStallPrintType)
		builder = builder.SetDeviceStallReceiptType(domainDevice.DeviceStallReceiptType)
	case domain.DeviceTypeCashier:
		builder = builder.SetOpenCashDrawer(domainDevice.OpenCashDrawer)
	}
	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrDeviceNotExists)
			return
		}
		err = fmt.Errorf("failed to update device: %w", err)
		return
	}
	domainDevice.UpdatedAt = updated.UpdatedAt
	return
}

func (repo *DeviceRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "DeviceRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.Device.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to delete device: %w", err)
		return
	}
	return nil
}

func (repo *DeviceRepository) GetDevices(ctx context.Context, pager *upagination.Pagination, filter *domain.DeviceListFilter, orderBys ...domain.DeviceOrderBy) (domainDevices []*domain.Device, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "DeviceRepository.GetDevices")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.buildFilterQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count device: %w", err)
		return
	}

	devices, err := query.
		Order(repo.orderBy(orderBys...)...).
		WithStore().
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query device: %w", err)
		return
	}

	domainDevices = lo.Map(devices, func(item *ent.Device, _ int) *domain.Device {
		return convertDeviceToDomain(item)
	})
	return
}

func (repo *DeviceRepository) Exists(ctx context.Context, params domain.DeviceExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "DeviceRepository.Exists")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.Device.Query()
	if params.MerchantID != uuid.Nil {
		query = query.Where(device.MerchantIDEQ(params.MerchantID))
	}
	if params.StoreID != uuid.Nil {
		query = query.Where(device.StoreIDEQ(params.StoreID))
	}
	if params.Name != "" {
		query = query.Where(device.NameEQ(params.Name))
	}
	if params.DeviceCode != "" {
		query = query.Where(device.DeviceCodeEQ(params.DeviceCode))
	}
	if params.ExcludeID != uuid.Nil {
		query = query.Where(device.IDNEQ(params.ExcludeID))
	}

	exists, err = query.Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check device existence: %w", err)
	}
	return
}

func (repo *DeviceRepository) buildFilterQuery(filter *domain.DeviceListFilter) *ent.DeviceQuery {
	query := repo.Client.Device.Query()
	if filter == nil {
		return query
	}

	if filter.MerchantID != uuid.Nil {
		query = query.Where(device.MerchantID(filter.MerchantID))
	}
	if filter.StoreID != uuid.Nil {
		query = query.Where(device.StoreID(filter.StoreID))
	}
	if filter.DeviceType != "" {
		query = query.Where(device.DeviceTypeEQ(filter.DeviceType))
	}
	if filter.Status != "" {
		query = query.Where(device.StatusEQ(filter.Status))
	}
	if filter.Name != "" {
		query = query.Where(device.NameContains(filter.Name))
	}

	return query
}

func (repo *DeviceRepository) orderBy(orderBys ...domain.DeviceOrderBy) []device.OrderOption {
	var opts []device.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.DeviceOrderByID:
			opts = append(opts, device.ByID(rule))
		case domain.DeviceOrderByCreatedAt:
			opts = append(opts, device.ByCreatedAt(rule))
		case domain.DeviceOrderBySortOrder:
			opts = append(opts, device.BySortOrder(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, device.ByCreatedAt(sql.OrderDesc()))
	}
	return opts
}

func convertDeviceToDomain(es *ent.Device) (d *domain.Device) {
	d = &domain.Device{
		ID:                     es.ID,
		MerchantID:             es.MerchantID,
		StoreID:                es.StoreID,
		Name:                   es.Name,
		DeviceType:             es.DeviceType,
		DeviceCode:             es.DeviceCode,
		DeviceBrand:            es.DeviceBrand,
		DeviceModel:            es.DeviceModel,
		Location:               es.Location,
		Enabled:                es.Enabled,
		IP:                     es.IP,
		Status:                 es.Status,
		PaperSize:              es.PaperSize,
		StallID:                es.StallID,
		OrderChannels:          es.OrderChannels,
		DiningWays:             es.DiningWays,
		DeviceStallPrintType:   es.DeviceStallPrintType,
		DeviceStallReceiptType: es.DeviceStallReceiptType,
		OpenCashDrawer:         es.OpenCashDrawer,
		SortOrder:              es.SortOrder,
		CreatedAt:              es.CreatedAt,
		UpdatedAt:              es.UpdatedAt,
	}
	if es.Edges.Store != nil {
		d.StoreName = es.Edges.Store.StoreName
	}
	return d
}
