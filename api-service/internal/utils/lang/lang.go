package lang

import (
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/logger"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/utils/json"
	"golang.org/x/text/language"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const localesDir = "internal/utils/lang/locales"

var bundle *i18n.Bundle

var supportedLanguages = []string{"en", "fr"}

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, lang := range supportedLanguages {
		translation, err := loadTranslation(lang)
		if err != nil {
			logger.Error(err)
			continue
		}
		err = bundle.AddMessages(language.Make(lang), translation.Messages...)
		if err != nil {
			logger.Error(err)
			continue
		}
	}
}

func loadTranslation(lang string) (*i18n.MessageFile, error) {
	translationFile := filepath.Join(localesDir, strings.ToLower(lang)+".json")
	return bundle.LoadMessageFile(translationFile)
}

func Translate(lang string, messageID string) string {
	return i18n.NewLocalizer(bundle, lang).
		MustLocalize(&i18n.LocalizeConfig{
			MessageID: messageID,
		})
}
