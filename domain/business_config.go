package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// 错误定义
// ------------------------------------------------------------

var (
	ErrBusinessConfigNotExists = errors.New("经营配置不存在")
)

// ------------------------------------------------------------
// 枚举定义
// ------------------------------------------------------------

// BusinessConfigGroup 计入规则
type BusinessConfigGroup string

const (
	BusinessConfigGroupPrint BusinessConfigGroup = "print" // print
)

func (BusinessConfigGroup) Values() []string {
	return []string{
		string(BusinessConfigGroupPrint),
	}
}

type BusinessConfigConfigType string

const (
	BusinessConfigConfigTypeString   BusinessConfigConfigType = "string"   // string
	BusinessConfigConfigTypeInt      BusinessConfigConfigType = "int"      // int
	BusinessConfigConfigTypeUint     BusinessConfigConfigType = "uint"     // uint
	BusinessConfigConfigTypeDatetime BusinessConfigConfigType = "datetime" // datetime
	BusinessConfigConfigTypeDate     BusinessConfigConfigType = "date"     // date
)

func (BusinessConfigConfigType) Values() []string {
	return []string{
		string(BusinessConfigConfigTypeString),
		string(BusinessConfigConfigTypeInt),
		string(BusinessConfigConfigTypeUint),
		string(BusinessConfigConfigTypeDatetime),
		string(BusinessConfigConfigTypeDate),
	}
}

type BusinessConfig struct {
	ID             uuid.UUID                `json:"id"`
	SourceConfigID uuid.UUID                `json:"source_config_id"` // 来源配置ID
	MerchantID     uuid.UUID                `json:"merchant_id"`      // 品牌商ID
	StoreID        uuid.UUID                `json:"store_id"`         // 门店ID
	Group          BusinessConfigGroup      `json:"group"`            // 配置分组
	Name           string                   `json:"name"`             // 参数名称
	ConfigType     BusinessConfigConfigType `json:"config_type"`      // 键值类型
	Key            string                   `json:"key"`              // 参数键名
	Value          string                   `json:"value"`            // 参数键值
	Sort           int32                    `json:"sort"`             // 排序
	Tip            string                   `json:"tip"`              // 变量描述
	IsDefault      bool                     `json:"is_default"`       // 是否为系统默认
	Status         bool                     `json:"status"`           // 状态
	CreatedAt      time.Time                `json:"created_at"`       // 创建时间
	UpdatedAt      time.Time                `json:"updated_at"`       // 更新时间
}

type BusinessConfigs []*BusinessConfig
type BusinessConfigSearchParams struct {
	MerchantID uuid.UUID
	StoreID    uuid.UUID
	Name       string // 结算方式名称（模糊匹配）
}

// BusinessConfigSearchRes 查询结果
type BusinessConfigSearchRes struct {
	*upagination.Pagination
	Items BusinessConfigs `json:"items"`
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/business_config_repository.go -package=mock . BusinessConfigRepository
type BusinessConfigRepository interface {
	ListBySearch(ctx context.Context, params BusinessConfigSearchParams) (*BusinessConfigSearchRes, error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/business_config_interactor.go -package=mock . BusinessConfigInteractor
type BusinessConfigInteractor interface {
	ListBySearch(ctx context.Context, params BusinessConfigSearchParams) (*BusinessConfigSearchRes, error)
}
