package i18n

import (
	"os"
	"path/filepath"
	"strings"
	"template/global"

	"github.com/jingyuexing/go-utils"
	"github.com/jingyuexing/i18n"
	"go.uber.org/zap"
)


const (
	LanguageEN      = "en"    // 英文
	LanguageJP      = "ja"    // 日文
	LanguageTW      = "zh-TW" // 台湾
	LanguageZH      = "zh-CN" // 中文
	LanguageKorean  = "ko"    // 韩文
	LanguageSpanish = "es"    // 西班牙文
	LanguageRussian = "ru"    // 俄文
	LanguageThai    = "th"    // 泰文
	LanguageHebrew  = "he"    // 希伯来文
	LanguageFrench  = "fr"    // 法语
	LanguageArabic  = "ar"    // 阿拉伯语
	LanguagePersian = "fa"    // 波斯语
)

var getLanguages = (func() []string {
	// 获取所有支持的语言文件
	lang := make([]string, 0)
	err := filepath.Walk(global.Config.System.Locale, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			global.Logger.Error("Error walking the path", zap.String("path", path))
			// fmt.Println("Error walking the path:", err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			// 提取语言代码（假设文件名格式为 lang.json）
			langCode := strings.TrimSuffix(info.Name(), ".json")
			lang = append(lang, langCode)
		}
		return nil
	})
	if err != nil {
		global.Logger.Error("Error enumerating languages", zap.String("error", err.Error()))
	}
	return lang
})()

var initLocal = "en"

// Translate the translate information map
var Translate = (func() map[string]any {
	message := map[string]any{}
	// 获取所有支持的语言
	for _, lang := range getLanguages {
		if _, ok := message[lang]; !ok {
			// 动态加载语言文件
			message[lang] = utils.LoadConfig[Locale](
				utils.Template(global.Config.System.Locale+"/{lang}.json", map[string]any{
					"lang": lang,
				}, "{}"))
		}
	}
	return message
})()

var loggerLanguageInit string = ""

// Translate the translate information map for logger
var TranslateLogger = (func() map[string]any {
	// initial language
	loggerLanguageInit = global.Config.Env.LoggerLanguage
	if global.Config.Env.LoggerLanguage == "" {
		loggerLanguageInit = initLocal
	}

	message := map[string]any{}
	// 获取所有支持的语言
	for _, lang := range getLanguages {
		if _, ok := message[lang]; !ok {
			// 动态加载语言文件
			message[lang] = utils.LoadConfig[Locale](
				utils.Template("locale/{lang}.json", map[string]any{
					"lang": lang,
				}, "{}")).Logger
		}
	}
	return message
})()

var I18N *i18n.I18n = i18n.CreateI18n(&i18n.Options{
	Local:          initLocal,
	FallbackLocale: "zh",
	Message:        Translate,
})

var LocaleLogger = i18n.CreateI18n(&i18n.Options{
	Local:          loggerLanguageInit,
	FallbackLocale: initLocal,
	Message:        TranslateLogger,
})
