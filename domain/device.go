package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrDeviceNotExists  = errors.New("设备不存在")
	ErrDeviceNameExists = errors.New("设备名称已存在")
	ErrDeviceCodeExists = errors.New("设备编号已存在")
)

// DeviceRepository 设备仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/device_repository.go -package=mock . DeviceRepository
type DeviceRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (device *Device, err error)
	Create(ctx context.Context, device *Device) (err error)
	Update(ctx context.Context, device *Device) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetDevices(ctx context.Context, pager *upagination.Pagination, filter *DeviceListFilter, orderBys ...DeviceOrderBy) (devices []*Device, total int, err error)
	Exists(ctx context.Context, params DeviceExistsParams) (exists bool, err error)
}

// DeviceInteractor 设备用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/device_interactor.go -package=mock . DeviceInteractor
type DeviceInteractor interface {
	Create(ctx context.Context, device *Device) (err error)
	Update(ctx context.Context, device *Device) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetDevice(ctx context.Context, id uuid.UUID) (*Device, error)
	GetDevices(ctx context.Context, pager *upagination.Pagination, filter *DeviceListFilter, orderBys ...DeviceOrderBy) (devices []*Device, total int, err error)
	DeviceSimpleUpdate(ctx context.Context, updateField DeviceSimpleUpdateType, device *Device) (err error)
}

// DeviceType 设备类型
// 收银机用于收银/展示，打印机用于小票或标签打印。
type DeviceType string

const (
	DeviceTypeCashier DeviceType = "cashier" // 收银机
	DeviceTypePrinter DeviceType = "printer" // 打印机
)

func (DeviceType) Values() []string {
	return []string{string(DeviceTypeCashier), string(DeviceTypePrinter)}
}

// DeviceStatus 在线/离线状态
// online 表示设备已连接，offline 表示未连接。
type DeviceStatus string

const (
	DeviceStatusOnline  DeviceStatus = "online"  // 在线
	DeviceStatusOffline DeviceStatus = "offline" // 离线
)

func (DeviceStatus) Values() []string {
	return []string{string(DeviceStatusOnline), string(DeviceStatusOffline)}
}

// DeviceLocation 设备所在区域
// front_hall 前厅，back_kitchen 后厨。
type DeviceLocation string

const (
	DeviceLocationFrontHall   DeviceLocation = "front_hall"   // 前厅
	DeviceLocationBackKitchen DeviceLocation = "back_kitchen" // 后厨
)

func (DeviceLocation) Values() []string {
	return []string{string(DeviceLocationFrontHall), string(DeviceLocationBackKitchen)}
}

type PaperSize string // 打印纸张尺寸

const (
	PaperSize58mm PaperSize = "58mm"
	PaperSize80mm PaperSize = "80mm"
)

func (PaperSize) Values() []string {
	return []string{string(PaperSize58mm), string(PaperSize80mm)}
}

// DeviceStallPrintType 出品部门总分单
type DeviceStallPrintType string

const (
	DeviceStallPrintTypeAll      DeviceStallPrintType = "all"      // 全部打印 总单+分单
	DeviceStallPrintTypeCombined DeviceStallPrintType = "combined" // 总单
	DeviceStallPrintTypeSeparate DeviceStallPrintType = "separate" // 分单
)

func (DeviceStallPrintType) Values() []string {
	return []string{
		string(DeviceStallPrintTypeAll),
		string(DeviceStallPrintTypeCombined),
		string(DeviceStallPrintTypeSeparate),
	}
}

// DeviceStallReceiptType 打印出品部门全部票据
type DeviceStallReceiptType string

const (
	DeviceStallReceiptTypeAll     DeviceStallReceiptType = "all"     // 全部打印
	DeviceStallReceiptTypeExclude DeviceStallReceiptType = "exclude" // 剔除部门商户
)

func (DeviceStallReceiptType) Values() []string {
	return []string{
		string(DeviceStallReceiptTypeAll),
		string(DeviceStallReceiptTypeExclude),
	}
}

// Device 设备实体
type Device struct {
	ID                     uuid.UUID              `json:"id"`                        // 设备 ID
	MerchantID             uuid.UUID              `json:"merchant_id"`               // 所属商户 ID
	StoreID                uuid.UUID              `json:"store_id"`                  // 所属门店 ID
	StoreName              string                 `json:"store_name"`                // 门店名称
	Name                   string                 `json:"name"`                      // 设备名称
	DeviceType             DeviceType             `json:"device_type"`               // 设备类型
	DeviceCode             string                 `json:"device_code"`               // 设备编号/序列号
	DeviceBrand            string                 `json:"device_brand"`              // 设备品牌
	DeviceModel            string                 `json:"device_model"`              // 设备型号
	Location               DeviceLocation         `json:"location"`                  // 设备位置
	Enabled                bool                   `json:"enabled"`                   // 启用/停用状态
	IP                     string                 `json:"ip"`                        // 设备 IP 地址
	Status                 DeviceStatus           `json:"status"`                    // 设备状态
	PaperSize              PaperSize              `json:"paper_size"`                // 打印纸张尺寸
	StallID                uuid.UUID              `json:"stall_id"`                  // 出品部门 ID
	StallName              string                 `json:"stall_name"`                // 出品部门名称
	OrderChannels          []OrderChannel         `json:"order_channels"`            // 订单来源
	DiningWays             []DiningWay            `json:"dining_ways"`               // 订单类型/就餐方式
	DeviceStallPrintType   DeviceStallPrintType   `json:"device_stall_print_type"`   // 出品部门总分单
	DeviceStallReceiptType DeviceStallReceiptType `json:"device_stall_receipt_type"` // 打印出品部门全部票据
	OpenCashDrawer         bool                   `json:"open_cash_drawer"`          // 收银机 开钱箱
	SortOrder              int                    `json:"sort_order"`                // 排序
	CreatedAt              time.Time              `json:"created_at"`                // 创建时间
	UpdatedAt              time.Time              `json:"updated_at"`                // 更新时间
}

// DeviceExistsParams 存在性检查参数
type DeviceExistsParams struct {
	MerchantID uuid.UUID
	StoreID    uuid.UUID
	Name       string
	DeviceCode string
	ExcludeID  uuid.UUID
}

// DeviceListFilter 查询过滤参数
type DeviceListFilter struct {
	MerchantID uuid.UUID    // 商户 ID
	StoreID    uuid.UUID    // 门店 ID
	DeviceType DeviceType   // 设备类型
	Status     DeviceStatus // 设备状态
	Name       string       // 设备名称模糊查询
}

type DeviceOrderByType int

const (
	_ DeviceOrderByType = iota
	DeviceOrderByID
	DeviceOrderByCreatedAt
	DeviceOrderBySortOrder
)

type DeviceOrderBy struct {
	OrderBy DeviceOrderByType
	Desc    bool
}

func NewDeviceOrderByID(desc bool) DeviceOrderBy {
	return DeviceOrderBy{OrderBy: DeviceOrderByID, Desc: desc}
}

func NewDeviceOrderByCreatedAt(desc bool) DeviceOrderBy {
	return DeviceOrderBy{OrderBy: DeviceOrderByCreatedAt, Desc: desc}
}

func NewDeviceOrderBySortOrder(desc bool) DeviceOrderBy {
	return DeviceOrderBy{OrderBy: DeviceOrderBySortOrder, Desc: desc}
}

// DeviceSimpleUpdateType 简易字段更新类型
type DeviceSimpleUpdateType string

const (
	DeviceSimpleUpdateTypeEnabled DeviceSimpleUpdateType = "enabled"
)
