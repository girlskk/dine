package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ali/oss"
)

// OssTokenReq 用于获取OSS上传凭证的请求参数
type OssTokenReq struct {
	Scene       domain.ObjectStorageScene `json:"scene" binding:"required,oneof=store product"` // 业务场景，store：门店相关 product: 商品相关
	Filename    string                    `json:"filename" binding:"required"`                  // 文件名
	ForDownload bool                      `json:"for_download"`                                 // 是否用于下载
}

// OssTokenResp 用于获取OSS上传凭证的响应参数
type OssTokenResp struct {
	oss.PolicyToken
	Key                string `json:"key"`
	ContentDisposition string `json:"content_disposition"`
}
