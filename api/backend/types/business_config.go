package types

import (
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// BusinessConfigListReq 经营设置列表列表请求
type BusinessConfigListReq struct {
	Name  string                     `json:"name" form:"name"`   // 设置名称（模糊匹配）
	Group domain.BusinessConfigGroup `json:"group" form:"group"` // 配置分组
}

// BusinessConfigUpsertReq 经营设置更新请求
type BusinessConfigUpsertReq struct {
	Configs []BusinessConfig `json:"configs" binding:"required,min=1,dive"` // 配置分组
}

type BusinessConfig struct {
	ID             string                          `json:"id" binding:"required"`                                                   // 记录ID
	SourceConfigID string                          `json:"source_config_id"`                                                        // 来源配置ID
	Group          domain.BusinessConfigGroup      `json:"group" binding:"required,oneof=print"`                                    // 配置分组
	Name           string                          `json:"name"`                                                                    // 参数名称
	ConfigType     domain.BusinessConfigConfigType `json:"config_type" binding:"required,oneof=string int uint bool datetime date"` // 键值类型
	Key            string                          `json:"key"`                                                                     // 参数键名
	Value          string                          `json:"value"`                                                                   // 参数键值
	Sort           int32                           `json:"sort"`                                                                    // 排序
	Tip            string                          `json:"tip"`                                                                     // 变量描述
}

// BusinessConfigDistributeReq 下发门店
type BusinessConfigDistributeReq struct {
	Ids          []uuid.UUID `json:"ids"  binding:"required,min=1"`      // 记录ID列表
	StoreIDs     []uuid.UUID `json:"store_ids" binding:"required,min=1"` // 门店ID列表（必选，多选）
	ModifyStatus bool        `json:"modify_status" binding:"required"`   // 下发后是否可以进行修改: true-可以, false-不可以
}
