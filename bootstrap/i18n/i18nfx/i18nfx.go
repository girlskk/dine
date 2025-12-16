package i18nfx

import (
	"gitlab.jiguang.dev/pos-dine/dine/bootstrap/i18n"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"i18n",
	fx.Provide(
		func() i18n.Config {
			return i18n.Config{
				LanguageDir: "etc/language",
			}
		},
	),
	fx.Invoke(i18n.Init),
)
