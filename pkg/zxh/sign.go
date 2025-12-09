package zxh

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

func sign(params map[string]any, secret string) string {
	// 排序参数
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 生成待签名字符串
	var dataToSign []string
	for _, k := range keys {
		if k == "sign" {
			continue
		}
		value := params[k]
		var valueStr string
		// 检查值的类型并进行适当处理
		switch v := value.(type) {
		case string:
			valueStr = v
		case float64:
			valueStr = strconv.FormatFloat(v, 'f', -1, 64)
		case int:
			valueStr = strconv.Itoa(v)
		case int8:
			valueStr = strconv.Itoa(int(v))
		case int16:
			valueStr = strconv.Itoa(int(v))
		case int32:
			valueStr = strconv.Itoa(int(v))
		case int64:
			valueStr = strconv.FormatInt(v, 10)
		case bool:
			valueStr = strconv.FormatBool(v)
		case map[string]any, []any:
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return ""
			}
			valueStr = string(jsonBytes)
		default:
			return ""
		}
		dataToSign = append(dataToSign, k+"="+valueStr)
	}
	joinedData := strings.Join(dataToSign, "&")

	// 生成签名
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(joinedData))

	return hex.EncodeToString(h.Sum(nil))
}
