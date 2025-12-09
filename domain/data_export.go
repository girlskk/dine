package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

const (
	TaskTypeDataExport = "data_export"
)

type DataExportTaskPayload struct {
	ID int `json:"id"`
}

type DataExportType string

const (
	DataExportTypeOrderListExport                   DataExportType = "order_list"                    // 订单列表导出
	DataExportTypeReconciliationRecordListExport    DataExportType = "reconciliation_record_list"    // 财务对账单列表导出
	DataExportTypeReconciliationRecordDetailsExport DataExportType = "reconciliation_record_details" // 财务对账单详情导出
	DataExportTypePointSettlementListExport         DataExportType = "point_settlement_list"         // 积分结算账单列表导出
	DataExportTypePointSettlementDetailsExport      DataExportType = "point_settlement_details"      // 积分结算账单详情导出
)

func (t DataExportType) Values() []string {
	return []string{
		string(DataExportTypeOrderListExport),
		string(DataExportTypeReconciliationRecordListExport),
		string(DataExportTypeReconciliationRecordDetailsExport),
		string(DataExportTypePointSettlementListExport),
		string(DataExportTypePointSettlementDetailsExport),
	}
}

type DataExportStatus string

const (
	DataExportStatusPending DataExportStatus = "pending" // 待导出
	DataExportStatusSuccess DataExportStatus = "success" // 导出成功
	DataExportStatusFailed  DataExportStatus = "failed"  // 导出失败
)

func (s DataExportStatus) Values() []string {
	return []string{
		string(DataExportStatusPending),
		string(DataExportStatusSuccess),
		string(DataExportStatusFailed),
	}
}

var DataExportStatusError = ParamsErrorf("数据导出记录的状态不支持")

// 数据导出
type DataExport struct {
	ID           int              `json:"id"`
	StoreID      int              `json:"store_id"`      // 门店ID
	Type         DataExportType   `json:"type"`          // 导出类型 order_list: 订单列表
	Status       DataExportStatus `json:"status"`        // 导出状态
	Params       json.RawMessage  `json:"param"`         // 导出参数
	FailedReason string           `json:"failed_reason"` // 导出失败原因
	OperatorType OperatorType     `json:"operator_type"` // 操作人类型
	OperatorID   int              `json:"operator_id"`   // 操作人ID
	OperatorName string           `json:"operator_name"` // 操作人姓名
	FileName     string           `json:"file_name"`     // 导出文件名（包含后缀）
	URL          string           `json:"url"`           // 下载地址
	CreatedAt    time.Time        `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time        `json:"updated_at"`    // 更新时间
}

type DataExportFilter struct {
	StoreID      int
	Type         DataExportType
	Status       DataExportStatus
	CreatedAtGte *time.Time
	CreatedAtLte *time.Time
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/data_export_repository.go -package=mock . DataExportRepository
type DataExportRepository interface {
	Create(ctx context.Context, dataExport *DataExport) (*DataExport, error)
	CreateBulk(ctx context.Context, dataExports []*DataExport) ([]*DataExport, error)
	Update(ctx context.Context, dataExport *DataExport) (*DataExport, error)
	Find(ctx context.Context, id int) (*DataExport, error)
	List(ctx context.Context, pager *upagination.Pagination, filter *DataExportFilter) ([]*DataExport, int, error)
}

type SubmitDataExportParams struct {
	StoreID  int
	Type     DataExportType
	Params   json.RawMessage
	FileName string
	Operator any
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/data_export_interactor.go -package=mock . DataExportInteractor
type DataExportInteractor interface {
	Submit(ctx context.Context, params ...*SubmitDataExportParams) ([]*DataExport, error)
	Run(ctx context.Context, id int) error
	Retry(ctx context.Context, storeID int, id int) error
	List(ctx context.Context, pager *upagination.Pagination, filter *DataExportFilter) ([]*DataExport, int, error)
}

type DataExportBizError struct {
	msg string
}

func (e *DataExportBizError) Error() string {
	return e.msg
}

func DataExportBizErrorf(format string, args ...any) *DataExportBizError {
	return &DataExportBizError{msg: fmt.Sprintf(format, args...)}
}

func IsDataExportBizError(err error) bool {
	if err == nil {
		return false
	}
	var e *DataExportBizError
	return errors.As(err, &e)
}

func BuildDataExportSubmitParams[Param any, Operator any](storeID int, type_ DataExportType, params []Param, fileName string, operator Operator) ([]*SubmitDataExportParams, error) {
	if len(params) == 0 {
		return nil, nil
	}

	name, ext := util.GetFileNameAndExt(fileName)

	submitParams := make([]*SubmitDataExportParams, 0, len(params))
	for i, param := range params {
		paramBytes, err := json.Marshal(param)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %w", err)
		}

		fileName := fmt.Sprintf("%s（%d）%s", name, i+1, ext)

		submitParams = append(submitParams, &SubmitDataExportParams{
			StoreID:  storeID,
			Type:     type_,
			Params:   paramBytes,
			FileName: fileName,
			Operator: operator,
		})
	}

	if len(submitParams) == 1 {
		submitParams[0].FileName = fileName
	}

	return submitParams, nil
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/data_exporter.go -package=mock . DataExporter

type DataExporter interface {
	NewParams() any
	Export(ctx context.Context, filename string, params any) (url string, err error)
}
