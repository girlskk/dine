package objectstorage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ali/oss"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.ObjectStorage = (*Storage)(nil)

type Storage struct {
	client *oss.Client
}

func NewStorage(client *oss.Client) *Storage {
	return &Storage{
		client: client,
	}
}

// PutObject 上传对象
// scene: 业务场景
// filename: 文件名
// reader: 读取器
// optFns: 选项函数
func (s *Storage) PutObject(
	ctx context.Context,
	scene domain.ObjectStorageScene,
	filename string,
	reader io.Reader,
	optFns ...func(*domain.ObjectStorageOption),
) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Storage.PutObject")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	opOpt := new(domain.ObjectStorageOption)

	for _, fn := range optFns {
		fn(opOpt)
	}

	key := domain.GenerateObjectKey(scene, filename)

	// 创建上传对象的请求
	req := s.client.NewDefaultPutObjectRequest(key)

	if opOpt.ForDownload {
		req.ContentDisposition = lo.ToPtr(oss.ContentDispositionAttachmentFilename(filename))
	}

	if _, err = s.client.NewUploader().UploadFrom(ctx, req, reader); err != nil {
		err = fmt.Errorf("failed to put object: %w", err)
		return
	}

	if url, err = s.client.FullURL(key); err != nil {
		err = fmt.Errorf("failed to get object url: %w", err)
		return
	}

	return
}

// ExportExcel 导出 Excel 文件
// scene: 业务场景
// readableName: 下载时的文件名(不包含后缀),可读性强
// headers: 表头
// contents: 数据
func (s *Storage) ExportExcel(
	ctx context.Context,
	scene domain.ObjectStorageScene,
	readableName string,
	headers []string,
	contents [][]string,
) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Storage.ExportExcel")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 创建 excel 文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logger := logging.FromContext(ctx).Named("Storage.ExportExcel")
			logger.Errorw("failed to close file", "error", err)
		}
	}()

	sheetName := "Sheet1"
	// 设置表头
	for i, header := range headers {
		var cell string
		cell, err = excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			err = fmt.Errorf("failed to get cell name: %w", err)
			return
		}
		f.SetCellValue(sheetName, cell, header)
	}
	// 写入数据
	for rowIndex, row := range contents {
		// 写入每一列的数据
		for colIndex, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 写入内存
	var buf bytes.Buffer
	if err = f.Write(&buf); err != nil {
		err = fmt.Errorf("failed to write excel: %w", err)
		return
	}

	readableName, _ = util.GetFileNameAndExt(readableName)
	filename := fmt.Sprintf("%s.xlsx", readableName)
	// 上传文件
	url, err = s.PutObject(ctx, scene, filename, &buf, domain.ObjectStorageWithForDownload())
	if err != nil {
		err = fmt.Errorf("failed to put object: %w", err)
		return
	}
	return
}

// ExportExcelWithBlankMerge 导出 Excel 文件，空列会自动向上合并
// scene: 业务场景
// readableName: 下载时的文件名(不包含后缀)，可读性强
// headers: 表头
// contents: 数据
func (s *Storage) ExportExcelWithBlankMerge(
	ctx context.Context,
	scene domain.ObjectStorageScene,
	readableName string,
	headers []string,
	contents [][]string,
) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Storage.ExportExcelWithBlankMerge")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 创建 Excel 文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logger := logging.FromContext(ctx).Named("Storage.ExportExcelWithBlankMerge")
			logger.Errorw("failed to close file", "error", err)
		}
	}()

	sheetName := "Sheet1"

	// 设置表头
	for i, header := range headers {
		var cell string
		cell, err = excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			err = fmt.Errorf("failed to get cell name: %w", err)
			return
		}
		f.SetCellValue(sheetName, cell, header)
	}

	// 记录每列最后一次非空的行号（用于合并单元格）
	lastNonEmptyRow := make([]int, len(headers))

	// 写入数据
	for rowIndex, row := range contents {
		excelRow := rowIndex + 2 // 从第二行开始写数据（第一行是表头）
		for colIndex, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, excelRow)

			if value != "" {
				// 记录当前列最新的非空行索引
				lastNonEmptyRow[colIndex] = excelRow
				f.SetCellValue(sheetName, cell, value)
			} else if lastNonEmptyRow[colIndex] > 1 { // 该列为空，并且上方有非空数据
				prevCell, _ := excelize.CoordinatesToCellName(colIndex+1, lastNonEmptyRow[colIndex])
				if err = f.MergeCell(sheetName, prevCell, cell); err != nil {
					err = fmt.Errorf("failed to merge cell: %w", err)
					return
				}
			}
		}
	}

	// 写入内存
	var buf bytes.Buffer
	if err = f.Write(&buf); err != nil {
		err = fmt.Errorf("failed to write excel: %w", err)
		return
	}

	readableName, _ = util.GetFileNameAndExt(readableName)
	filename := fmt.Sprintf("%s.xlsx", readableName)
	// 上传文件
	url, err = s.PutObject(ctx, scene, filename, &buf, domain.ObjectStorageWithForDownload())
	if err != nil {
		err = fmt.Errorf("failed to put object: %w", err)
		return
	}
	return
}
