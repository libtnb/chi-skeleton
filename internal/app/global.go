package app

import (
	"github.com/go-chi/chi/v5"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/v2"
	"gorm.io/gorm"
)

var (
	Conf       *koanf.Koanf
	Http       *chi.Mux
	Orm        *gorm.DB
	Validator  *validator.Validate
	Translator *ut.Translator
)
