package i18n

import (
	"context"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	DefaultLocale = "en-US"
	LocaleEnUS    = "en-US"
	LocaleZhCN    = "zh-CN"
)

type contextKey string

const (
	localizerKey contextKey = "i18n_localizer"
)

var (
	bundle *i18n.Bundle
)

// SetBundle 设置 i18n bundle（由 bootstrap/i18n 调用）
func SetBundle(b *i18n.Bundle) {
	bundle = b
}

// getBundle 获取 i18n bundle
func getBundle() *i18n.Bundle {
	return bundle
}

// getLocalizer 从 context 获取或创建 Localizer（带缓存）
func getLocalizer(ctx context.Context) *i18n.Localizer {
	// 尝试从 context 中获取缓存的 Localizer
	if localizer, ok := ctx.Value(localizerKey).(*i18n.Localizer); ok {
		return localizer
	}
	return nil
}

// WithLocalizer 将 Localizer 添加到 context 中（用于中间件优化）
func WithLocalizer(ctx context.Context, locale string) context.Context {
	if locale == "" {
		locale = DefaultLocale
	}

	b := getBundle()
	if b == nil {
		return ctx
	}

	localizer := i18n.NewLocalizer(b, locale)
	return context.WithValue(ctx, localizerKey, localizer)
}

// Translate 翻译消息，如果翻译失败则返回 messageID
// ctx: 上下文，用于获取语言环境
// messageID: 消息 ID
// templateData: 模板数据（可选，用于参数化翻译）
func Translate(ctx context.Context, messageID string, templateData map[string]any) string {
	localizer := getLocalizer(ctx)
	if localizer == nil {
		return messageID
	}

	config := &i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	}

	// 如果 templateData 为 nil，显式设置为 nil
	if templateData == nil {
		config.TemplateData = nil
	}

	msg, err := localizer.Localize(config)
	if err != nil || msg == "" {
		// 翻译失败，返回原始 messageID
		return messageID
	}

	return msg
}
