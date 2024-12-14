package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/v2"
	"gorm.io/gorm"
)

type App struct {
	conf       *koanf.Koanf
	router     *chi.Mux
	http       *http.Server
	db         *gorm.DB
	validator  *validator.Validate
	translator *ut.Translator
}

func NewApp(conf *koanf.Koanf, router *chi.Mux, http *http.Server, db *gorm.DB, validator *validator.Validate, translator *ut.Translator) *App {
	return &App{
		conf:       conf,
		router:     router,
		http:       http,
		db:         db,
		validator:  validator,
		translator: translator,
	}
}

func (r *App) Run() error {
	return r.http.ListenAndServe()
}
