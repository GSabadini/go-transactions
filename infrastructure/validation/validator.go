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

func TranslateErr(err error) map[string]string {
	var errs = make(map[string]string, 0)

	for _, e := range err.(validator.ValidationErrors) {
		errs[formatErrName(e.Namespace())] = e.Translate(translate)
	}

	return errs
}

func formatErrName(n string) string {
	return strings.Join(strings.Split(strings.ToLower(n), ".")[1:], ".")
}
