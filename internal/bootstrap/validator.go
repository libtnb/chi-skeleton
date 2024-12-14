package bootstrap

import (
	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
)

func NewValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}

func NewTranslator(validate *validator.Validate) (*ut.Translator, error) {
	translator := zh_Hans_CN.New()
	uni := ut.New(translator, translator)
	trans, _ := uni.GetTranslator("zh_Hans_CN")

	if err := zh.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, err
	}

	return &trans, nil
}
