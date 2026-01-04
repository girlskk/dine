package types

import (
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// DeviceCreateReq 创建设备请求
type DeviceCreateReq struct {
	StoreID       uuid.UUID             `json:"store_id" binding:"required"`                               // 门店 ID
	Name          string                `json:"name" binding:"required,max=50"`                            // 设备名称
	DeviceType    domain.DeviceType     `json:"device_type" binding:"required,oneof=cashier printer"`      // 设备类型
	DeviceCode    string                `json:"device_code" binding:"required,max=100"`                    // 设备编号
	DeviceBrand   string                `json:"device_brand" binding:"omitempty,max=100"`                  // 设备品牌
	DeviceModel   string                `json:"device_model" binding:"omitempty,max=100"`                  // 设备型号
	Location      domain.DeviceLocation `json:"location" binding:"required,oneof=front_hall back_kitchen"` // 设备位置
	Enabled       bool                  `json:"enabled"`                                                   // 是否启用
	SortOrder     int                   `json:"sort_order" binding:"omitempty,gte=0"`                      // 排序值，越小越靠前
	DevicePrint   DevicePrint           `json:"device_print" binding:"omitempty"`                          // 打印机设备配置
	DeviceCashier DeviceCashier         `json:"device_cashier" binding:"omitempty"`                        // 收银机设备配置
}

// DevicePrint 打印机设备
type DevicePrint struct {
	IP                     string                        `json:"ip" binding:"required,max=50"`                                                                                    // 打印机 IP 地址
	PaperSize              domain.PaperSize              `json:"paper_size" binding:"required,oneof=58mm 80mm"`                                                                   // 纸张大小
	StallID                uuid.UUID                     `json:"stall_id" binding:"required"`                                                                                     // 出品部门 ID
	OrderChannels          []domain.OrderChannel         `json:"order_channels" binding:"required,dive,oneof=pos self_order mini_program mobile_order scan_order third_delivery"` // 订单来源
	DiningWays             []domain.DiningWay            `json:"dining_ways" binding:"required,dive,oneof=dine_in take_out delivery"`                                             // 用餐方式
	DeviceStallPrintType   domain.DeviceStallPrintType   `json:"device_stall_print_type" binding:"required,oneof=all combined separate"`                                          // 打印出品部门总分单
	DeviceStallReceiptType domain.DeviceStallReceiptType `json:"device_stall_receipt_type" binding:"required,oneof=all exclude"`                                                  // 打印出品部门全部票据
}

// DeviceCashier 收银机设备
type DeviceCashier struct {
	OpenCashDrawer bool `json:"open_cash_drawer"` // 开启钱箱
}

// DeviceUpdateReq 更新设备请求
type DeviceUpdateReq struct {
	StoreID       uuid.UUID             `json:"store_id" binding:"required"`                               // 门店 ID
	Name          string                `json:"name" binding:"required,max=50"`                            // 设备名称
	DeviceType    domain.DeviceType     `json:"device_type" binding:"required,oneof=cashier printer"`      // 设备类型
	DeviceCode    string                `json:"device_code" binding:"omitempty,max=100"`                   // 设备编号
	DeviceBrand   string                `json:"device_brand" binding:"omitempty,max=100"`                  // 设备品牌
	DeviceModel   string                `json:"device_model" binding:"omitempty,max=100"`                  // 设备型号
	Location      domain.DeviceLocation `json:"location" binding:"required,oneof=front_hall back_kitchen"` // 设备位置
	Enabled       bool                  `json:"enabled"`                                                   // 是否启用
	SortOrder     int                   `json:"sort_order" binding:"omitempty,gte=0"`                      // 排序值，越小越靠前
	DevicePrint   DevicePrint           `json:"device_print" binding:"omitempty"`                          // 打印机设备配置
	DeviceCashier DeviceCashier         `json:"device_cashier" binding:"omitempty"`                        // 收银机设备配置
}

// DeviceListReq 设备列表查询
type DeviceListReq struct {
	upagination.RequestPagination
	StoreID    uuid.UUID           `form:"store_id"`    // 门店 ID
	DeviceType domain.DeviceType   `form:"device_type"` // 设备类型
	Status     domain.DeviceStatus `form:"status"`      // 设备状态
	Name       string              `form:"name"`        // 设备名称模糊查询
}

type DeviceListResp struct {
	Devices []*domain.Device `json:"devices"` // 设备列表
	Total   int              `json:"total"`   // 总数
}
