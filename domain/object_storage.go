package domain

import (
	"context"
	"io"
	"path"
	"time"

	"github.com/google/uuid"
)

type ObjectStorageScene string

const (
	ObjectStorageSceneMerchant ObjectStorageScene = "merchant" // 商户
	ObjectStorageSceneStore    ObjectStorageScene = "store"    // 门店
)

type ObjectStorageOption struct {
	ForDownload bool // 是否用于下载
}

// ObjectStorageWithForDownload 用于下载
func ObjectStorageWithForDownload() func(*ObjectStorageOption) {
	return func(o *ObjectStorageOption) {
		o.ForDownload = true
	}
}

// ObjectStorage 对象存储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/object_storage.go -package=mock . ObjectStorage
type ObjectStorage interface {
	PutObject(ctx context.Context, scene ObjectStorageScene, filename string, reader io.Reader, optFns ...func(*ObjectStorageOption)) (url string, err error)
	ExportExcel(ctx context.Context, scene ObjectStorageScene, readableName string, headers []string, contents [][]string) (url string, err error)
	ExportExcelWithBlankMerge(ctx context.Context, scene ObjectStorageScene, readableName string, headers []string, contents [][]string) (url string, err error)
}

// GenerateObjectKey 生成对象存储的key
func GenerateObjectKey(scene ObjectStorageScene, filename string) string {
	filename = uuid.New().String() + path.Ext(filename)
	return path.Join(string(scene), time.Now().Format("20060102"), filename)
}
