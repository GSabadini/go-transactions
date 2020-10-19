package validation

import (
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// use a single instance , it caches struct info
var (
	uni       *ut.UniversalTranslator
	validate  *validator.Validate
	translate ut.Translator
)

// NewValidator create new validator.Validate
func NewValidator() *validator.Validate {
	en := en.New()
	uni = ut.New(en, en)
	translate, _ = uni.GetTranslator("en")

	validate = validator.New()
	if err := en_translations.RegisterDefaultTranslations(validate, translate); err != nil {
		log.Fatal(err)
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}

func ErrMessages(err error) []string {
	var msgs []string
	for _, e := range err.(validator.ValidationErrors) {
		msgs = append(msgs, e.Translate(translate))
	}
	return msgs
}
