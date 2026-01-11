package i18n

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	i18npkg "gitlab.jiguang.dev/pos-dine/dine/pkg/i18n"
	"golang.org/x/text/language"
)

var (
	bundle *i18n.Bundle
)

// Config i18n 配置
type Config struct {
	LanguageDir string `default:"etc/language"` // 语言文件目录
}

// Init 初始化 i18n，从指定目录加载 TOML 翻译文件
func Init(cfg Config) error {
	bundle = i18n.NewBundle(language.English)

	// 注册 TOML 解析器
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 递归读取语言目录下的所有 TOML 文件（包含子目录）
	err := filepath.WalkDir(cfg.LanguageDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".toml") {
			_, err := bundle.LoadMessageFile(path)
			if err != nil {
				return fmt.Errorf("failed to load language file %s: %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk language directory %s: %w", cfg.LanguageDir, err)
	}

	// 将初始化好的 bundle 设置到 pkg/i18n
	i18npkg.SetBundle(bundle)

	return nil
}
