package oss

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/url"
	"path"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	openapicred "github.com/aliyun/credentials-go/credentials"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

const (
	product                  = "oss"
	contentDispositionInline = "inline"
)

// PolicyToken 结构体用于存储生成的表单数据
type PolicyToken struct {
	Policy           string `json:"policy"`
	SecurityToken    string `json:"security_token"`
	SignatureVersion string `json:"x_oss_signature_version"`
	Credential       string `json:"x_oss_credential"`
	Date             string `json:"x_oss_date"`
	Signature        string `json:"signature"`
	Host             string `json:"host"`
	Dir              string `json:"dir"`
}

type Client struct {
	Config
	host string
	*oss.Client
}

func New(config Config) *Client {
	provider := credentials.NewStaticCredentialsProvider(config.AccessKeyID, config.AccessKeySecret)
	cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider).WithRegion(config.Region)

	return &Client{
		Config: config,
		Client: oss.NewClient(cfg),
		host:   lo.Ternary(config.Domain != "", config.Domain, config.host()),
	}
}

func (c *Client) GeneratePolicyToken() (*PolicyToken, error) {
	config := new(openapicred.Config).
		SetType("ram_role_arn").
		SetAccessKeyId(c.AccessKeyID).
		SetAccessKeySecret(c.AccessKeySecret).
		SetRoleArn(c.RoleARN).
		SetRoleSessionName(c.RoleSessionName).
		SetPolicy("").
		SetRoleSessionExpiration(c.Expiration)

	// 根据配置创建凭证提供器
	provider, err := openapicred.NewCredential(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	// 从凭证提供器获取凭证
	cred, err := provider.GetCredential()
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	// 构建policy
	utcTime := time.Now().UTC()
	date := utcTime.Format("20060102")
	expiration := utcTime.Add(time.Duration(c.Expiration) * time.Second)
	policyMap := map[string]any{
		"expiration": expiration.Format("2006-01-02T15:04:05.000Z"),
		"conditions": []any{
			map[string]string{"bucket": c.Bucket},
			map[string]string{"x-oss-signature-version": "OSS4-HMAC-SHA256"},
			map[string]string{"x-oss-credential": fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", *cred.AccessKeyId, date, c.Region, product)},
			map[string]string{"x-oss-date": utcTime.Format("20060102T150405Z")},
			map[string]string{"x-oss-security-token": *cred.SecurityToken},
		},
	}

	// 将policy转换为 JSON 格式
	policy, err := json.Marshal(policyMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal policy: %w", err)
	}

	// 构造待签名字符串（StringToSign）
	stringToSign := base64.StdEncoding.EncodeToString([]byte(policy))

	hmacHash := func() hash.Hash { return sha256.New() }
	// 构建signing key
	signingKey := "aliyun_v4" + *cred.AccessKeySecret
	h1 := hmac.New(hmacHash, []byte(signingKey))
	io.WriteString(h1, date)
	h1Key := h1.Sum(nil)

	h2 := hmac.New(hmacHash, h1Key)
	io.WriteString(h2, c.Config.Region)
	h2Key := h2.Sum(nil)

	h3 := hmac.New(hmacHash, h2Key)
	io.WriteString(h3, product)
	h3Key := h3.Sum(nil)

	h4 := hmac.New(hmacHash, h3Key)
	io.WriteString(h4, "aliyun_v4_request")
	h4Key := h4.Sum(nil)

	// 生成签名
	h := hmac.New(hmacHash, h4Key)
	io.WriteString(h, stringToSign)
	signature := hex.EncodeToString(h.Sum(nil))

	return &PolicyToken{
		Policy:           stringToSign,
		SecurityToken:    *cred.SecurityToken,
		SignatureVersion: "OSS4-HMAC-SHA256",
		Credential:       fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", *cred.AccessKeyId, date, c.Region, product),
		Date:             utcTime.UTC().Format("20060102T150405Z"),
		Signature:        signature,
		Host:             c.host,
		Dir:              c.Dir,
	}, nil
}

func (c *Client) NewDefaultPutObjectRequest(key string) *oss.PutObjectRequest {
	return &oss.PutObjectRequest{
		Bucket:       oss.Ptr(c.Bucket),              // 存储空间名称
		Key:          oss.Ptr(path.Join(c.Dir, key)), // 对象名称
		StorageClass: oss.StorageClassStandard,       // 指定对象的存储类型为标准存储
		Acl:          oss.ObjectACLDefault,           // 指定对象的访问权限为私有访问
	}
}

// FullURL 获取文件的完整 URL
func (c *Client) FullURL(key string) (string, error) {
	return url.JoinPath(c.host, c.Dir, key)
}

// ContentDispositionAttachmentFilename 生成指定下载文件名的Content-Disposition头
func ContentDispositionAttachmentFilename(filename string) string {
	name, ext := util.GetFileNameAndExt(filename)
	return fmt.Sprintf(`attachment; filename="%s%s"`, url.QueryEscape(name), ext)
}

// ContentDispositionInline 生成预览的Content-Disposition头
func ContentDispositionInline() string {
	return contentDispositionInline
}
