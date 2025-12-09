package zxh

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
)

const (
	CodeSuccess = 0
)

type Client struct {
	mac *Credentials
	rc  *resty.Client
}

type response[T any] struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    T      `json:"data"`
}

type request[T any] struct {
	client *Client
	Path   string
	Data   map[string]any
	Res    T
}

func newRequest[T any](c *Client, path string, data map[string]any, res T) *request[T] {
	return &request[T]{
		client: c,
		Path:   path,
		Data:   data,
		Res:    res,
	}
}

// generateRandomString32 generates a 32-character random string
func generateRandomString32() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func (r *request[T]) do(ctx context.Context) error {
	logger := logging.FromContext(ctx).Named("ZhixinhuaClient.do")

	data := r.Data
	data["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
	data["nonceStr"] = generateRandomString32()

	r.client.mac.auth(data)

	resp := response[T]{
		Data: r.Res,
	}

	logger.Infof("ponintpay request: %v", data)

	res, err := r.client.rc.R().
		SetContext(ctx).
		SetBody(data).
		SetResult(&resp).
		Post(r.Path)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	logger.Infof("ponintpay response: %v", res.String())

	if res.IsError() {
		return fmt.Errorf("request status: %s, error: %s", res.Status(), res.String())
	}

	if resp.ErrCode != CodeSuccess {
		return &APIError{
			ErrCode: resp.ErrCode,
			ErrMsg:  resp.ErrMsg,
			ErrDlt:  resp.ErrDlt,
		}
	}

	return nil
}
