package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// DeviceListReq 设备列表查询
type DeviceListReq struct {
	DeviceType domain.DeviceType   `form:"device_type"` // 设备类型
	Status     domain.DeviceStatus `form:"status"`      // 设备状态
}

// DeviceListResp 设备列表
type DeviceListResp struct {
	Devices []*domain.Device `json:"devices"` // 设备列表
	Total   int              `json:"total"`   // 总数
}

type SyncStoreDeviceStatusReq struct {
	DeviceCodes []string `json:"device_codes"` // 设备编号列表
}
