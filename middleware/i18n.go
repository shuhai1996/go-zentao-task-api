package middleware

import (
	"go-zentao-task/core"
	"go-zentao-task/pkg/i18n"
	"golang.org/x/text/language"
)

func init() {
	i18n.Setup("zh-CN", "docs/i18n/locale_zh-CN.ini") //nolint
	i18n.Setup("en-US", "docs/i18n/locale_en-US.ini") //nolint
}

func I18n(c *core.Context) {
	lang := c.GetHeader("Request-Language")
	if lang == "" || !i18n.Exist(lang) {
		tag, _, err := language.ParseAcceptLanguage(c.GetHeader("Accept-Language"))
		if err == nil {
			for _, v := range tag {
				if i18n.Exist(v.String()) {
					lang = v.String()
					break
				}
			}
		}
	}

	c.Locale = &i18n.Locale{Lang: lang}
	c.Next()
}
