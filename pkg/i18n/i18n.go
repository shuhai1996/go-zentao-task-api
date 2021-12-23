package i18n

import (
	"fmt"
	"github.com/verystar/ini"
	"strings"
)

var (
	locales     = make(map[string]*ini.Ini)
	defaultLang = "zh-CN"
)

type Locale struct {
	Lang string
}

func Setup(lang, langFile string, langFiles ...string) error {
	if _, ok := locales[lang]; ok {
		return nil
	}

	cfg, err := ini.Load(langFile, langFiles...)
	if err != nil {
		return err
	}
	locales[lang] = cfg
	return nil
}

func SetDefaultLang(lang string) {
	defaultLang = lang
}

func Exist(lang string) bool {
	_, ok := locales[lang]
	return ok
}

func Tr(lang, key string, args ...interface{}) string {
	if key == "" {
		return ""
	}
	if _, ok := locales[lang]; !ok {
		lang = defaultLang
	}
	if locales[lang] == nil {
		return ""
	}

	var section string
	dotIndex := strings.Index(key, ".")
	if dotIndex != -1 {
		section = key[:dotIndex]
		key = key[dotIndex+1:]
	}

	cfg := locales[lang].Read(section, key)
	if cfg == "" && section != "" {
		cfg = locales[lang].Read("", section+"."+key)
	}

	if len(args) > 0 {
		return fmt.Sprintf(cfg, args...)
	}
	return cfg
}

func (l *Locale) Tr(key string, args ...interface{}) string {
	return Tr(l.Lang, key, args...)
}
