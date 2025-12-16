package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/i18n"
)

type Locale struct{}

func NewLocale() *Locale {
	return &Locale{}
}

func (m *Locale) Name() string {
	return "Locale"
}

func (m *Locale) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := extractLocale(c)
		// 优化：预先创建 Localizer 并缓存到 context
		ctx := i18n.WithLocalizer(c.Request.Context(), locale)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// extractLocale 从请求中提取语言环境
// 优先级：1. Query 参数 ?locale=en-US
//  2. Header Accept-Language: en-US
//  3. Header X-Locale: en-US
//  4. 默认 zh-CN
func extractLocale(c *gin.Context) string {
	// 1. 从 Query 参数获取
	if locale := c.Query("locale"); locale != "" {
		return normalizeLocale(locale)
	}

	// 2. 从自定义 Header 获取
	if locale := c.GetHeader("X-Locale"); locale != "" {
		return normalizeLocale(locale)
	}

	// 3. 从 Accept-Language Header 获取
	if acceptLang := c.GetHeader("Accept-Language"); acceptLang != "" {
		locale := parseAcceptLanguage(acceptLang)
		if locale != "" {
			return locale
		}
	}

	// 4. 默认返回中文
	return i18n.DefaultLocale
}

// normalizeLocale 规范化 locale 格式
func normalizeLocale(locale string) string {
	locale = strings.TrimSpace(locale)
	locale = strings.ToLower(locale)

	// 支持简写：en -> en-US, zh -> zh-CN
	switch locale {
	case "en":
		return i18n.LocaleEnUS
	case "zh":
		return i18n.LocaleZhCN
	case "zh-cn", "zh_cn":
		return i18n.LocaleZhCN
	case "en-us", "en_us":
		return i18n.LocaleEnUS
	default:
		return i18n.DefaultLocale
	}
}

// parseAcceptLanguage 解析 Accept-Language Header
// 例如: "en-US,en;q=0.9,zh-CN;q=0.8" -> "en-US"
func parseAcceptLanguage(acceptLang string) string {
	// 简化实现：取第一个语言标签
	parts := strings.Split(acceptLang, ",")
	if len(parts) > 0 {
		lang := strings.TrimSpace(parts[0])
		// 移除质量值，例如: "en-US;q=0.9" -> "en-US"
		if idx := strings.Index(lang, ";"); idx > 0 {
			lang = lang[:idx]
		}
		return normalizeLocale(lang)
	}
	return ""
}
