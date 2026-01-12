package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/store"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type DeviceRepositoryTestSuite struct {
	RepositoryTestSuite
	repo *DeviceRepository
	ctx  context.Context
}

func TestDeviceRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceRepositoryTestSuite))
}

func (s *DeviceRepositoryTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &DeviceRepository{Client: s.client}
	s.ctx = context.Background()
}

func (s *DeviceRepositoryTestSuite) newPrinter(tag string) *domain.Device {
	return &domain.Device{
		ID:                     uuid.New(),
		MerchantID:             uuid.New(),
		StoreID:                uuid.New(),
		Name:                   "打印机-" + tag,
		DeviceType:             domain.DeviceTypePrinter,
		DeviceCode:             "CODE-" + tag,
		DeviceBrand:            "品牌-" + tag,
		DeviceModel:            "型号-" + tag,
		Location:               domain.DeviceLocationBackKitchen,
		Enabled:                true,
		IP:                     "192.168.1." + tag[len(tag)-1:],
		Status:                 domain.DeviceStatusOnline,
		PaperSize:              domain.PaperSize80mm,
		ConnectType:            domain.DeviceConnectTypeInside,
		StallID:                uuid.New(),
		OrderChannels:          []domain.OrderChannel{domain.OrderChannelPOS},
		DiningWays:             []domain.DiningWay{domain.DiningWayDineIn},
		DeviceStallPrintType:   domain.DeviceStallPrintTypeAll,
		DeviceStallReceiptType: domain.DeviceStallReceiptTypeAll,
		SortOrder:              5,
	}
}

func (s *DeviceRepositoryTestSuite) newCashier(tag string) *domain.Device {
	return &domain.Device{
		ID:             uuid.New(),
		MerchantID:     uuid.New(),
		StoreID:        uuid.New(),
		Name:           "收银机-" + tag,
		DeviceType:     domain.DeviceTypeCashier,
		DeviceCode:     "CODE-CASH-" + tag,
		DeviceBrand:    "品牌-" + tag,
		DeviceModel:    "型号-" + tag,
		Location:       domain.DeviceLocationFrontHall,
		Enabled:        true,
		Status:         domain.DeviceStatusOffline,
		OpenCashDrawer: true,
		SortOrder:      8,
	}
}

func (s *DeviceRepositoryTestSuite) TestDevice_Create() {

	s.T().Run("创建打印机成功", func(t *testing.T) {
		device := s.newPrinter("p1")
		err := s.repo.Create(s.ctx, device)
		require.NoError(t, err)

		saved := s.client.Device.GetX(s.ctx, device.ID)
		require.Equal(t, device.Name, saved.Name)
		require.Equal(t, device.IP, saved.IP)
		require.Equal(t, device.PaperSize, saved.PaperSize)
		require.Equal(t, device.OrderChannels, saved.OrderChannels)
	})

	s.T().Run("创建收银机成功", func(t *testing.T) {
		device := s.newCashier("c1")
		err := s.repo.Create(s.ctx, device)
		require.NoError(t, err)

		saved := s.client.Device.GetX(s.ctx, device.ID)
		require.Equal(t, device.Name, saved.Name)
		require.Equal(t, device.OpenCashDrawer, saved.OpenCashDrawer)
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Create(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *DeviceRepositoryTestSuite) TestDevice_FindByID() {
	device := s.newPrinter("find")
	require.NoError(s.T(), s.repo.Create(s.ctx, device))

	s.T().Run("查询成功", func(t *testing.T) {
		found, err := s.repo.FindByID(s.ctx, device.ID)
		require.NoError(t, err)
		require.Equal(t, device.ID, found.ID)
		require.Equal(t, device.StoreID, found.StoreID)
		require.Equal(t, store.StoreName, found.StoreName)
	})

	s.T().Run("不存在", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *DeviceRepositoryTestSuite) TestDevice_Update() {
	printer := s.newPrinter("update")
	require.NoError(s.T(), s.repo.Create(s.ctx, printer))
	cashier := s.newCashier("update")
	require.NoError(s.T(), s.repo.Create(s.ctx, cashier))

	s.T().Run("更新打印机成功", func(t *testing.T) {
		printer.Name = "更新-" + printer.Name
		printer.Enabled = false
		printer.PaperSize = domain.PaperSize58mm
		printer.SortOrder = 1

		err := s.repo.Update(s.ctx, printer)
		require.NoError(t, err)

		updated := s.client.Device.GetX(s.ctx, printer.ID)
		require.Equal(t, printer.Name, updated.Name)
		require.Equal(t, printer.PaperSize, updated.PaperSize)
		require.Equal(t, printer.SortOrder, updated.SortOrder)
	})

	s.T().Run("更新收银机成功", func(t *testing.T) {
		cashier.Name = "更新-" + cashier.Name
		cashier.OpenCashDrawer = false

		err := s.repo.Update(s.ctx, cashier)
		require.NoError(t, err)

		updated := s.client.Device.GetX(s.ctx, cashier.ID)
		require.Equal(t, cashier.Name, updated.Name)
		require.Equal(t, cashier.OpenCashDrawer, updated.OpenCashDrawer)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		missing := s.newPrinter("missing")
		err := s.repo.Update(s.ctx, missing)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})

	s.T().Run("空入参", func(t *testing.T) {
		err := s.repo.Update(s.ctx, nil)
		require.Error(t, err)
	})
}

func (s *DeviceRepositoryTestSuite) TestDevice_Delete() {
	device := s.newPrinter("delete")
	require.NoError(s.T(), s.repo.Create(s.ctx, device))

	s.T().Run("删除成功", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, device.ID)
		require.NoError(t, err)
		_, err = s.client.Device.Get(s.ctx, device.ID)
		require.Error(t, err)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *DeviceRepositoryTestSuite) TestDevice_GetDevices() {
	device1 := s.newPrinter("001")
	require.NoError(s.T(), s.repo.Create(s.ctx, device1))
	time.Sleep(10 * time.Millisecond)
	device2 := s.newPrinter("002")
	device2.Enabled = false
	device2.SortOrder = 2
	device2.Status = domain.DeviceStatusOffline
	require.NoError(s.T(), s.repo.Create(s.ctx, device2))
	time.Sleep(10 * time.Millisecond)
	device3 := s.newCashier("003")
	device3.SortOrder = 1
	require.NoError(s.T(), s.repo.Create(s.ctx, device3))

	pager := upagination.New(1, 10)

	s.T().Run("按门店筛选默认排序", func(t *testing.T) {
		list, total, err := s.repo.GetDevices(s.ctx, pager, &domain.DeviceListFilter{})
		require.NoError(t, err)
		require.Equal(t, 3, total)
		require.Len(t, list, 3)
		require.Equal(t, device3.ID, list[0].ID)
	})

	s.T().Run("按排序字段升序", func(t *testing.T) {
		order := domain.NewDeviceOrderBySortOrder(false)
		list, total, err := s.repo.GetDevices(s.ctx, pager, &domain.DeviceListFilter{}, order)
		require.NoError(t, err)
		require.Equal(t, 3, total)
		require.Equal(t, device3.ID, list[0].ID)
	})

	s.T().Run("按类型与状态筛选", func(t *testing.T) {
		list, total, err := s.repo.GetDevices(s.ctx, pager, &domain.DeviceListFilter{DeviceType: domain.DeviceTypePrinter, Status: domain.DeviceStatusOffline})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, device2.ID, list[0].ID)
	})

	s.T().Run("按名称模糊筛选", func(t *testing.T) {
		list, total, err := s.repo.GetDevices(s.ctx, pager, &domain.DeviceListFilter{Name: "002"})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, device2.ID, list[0].ID)
	})
}

func (s *DeviceRepositoryTestSuite) TestDevice_Exists() {
	device := s.newPrinter("exists")
	require.NoError(s.T(), s.repo.Create(s.ctx, device))

	s.T().Run("同名存在", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.DeviceExistsParams{Name: device.Name})
		require.NoError(t, err)
		require.True(t, exists)
	})

	s.T().Run("排除自身", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.DeviceExistsParams{Name: device.Name, ExcludeID: device.ID})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("不同门店不冲突", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.DeviceExistsParams{Name: device.Name, MerchantID: uuid.New(), StoreID: uuid.New()})
		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("按设备编号查询", func(t *testing.T) {
		exists, err := s.repo.Exists(s.ctx, domain.DeviceExistsParams{DeviceCode: device.DeviceCode})
		require.NoError(t, err)
		require.True(t, exists)
	})
}
