package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
	"gorm.io/gorm"
)

type App struct {
	conf   *koanf.Koanf
	router *chi.Mux
	http   *http.Server
	db     *gorm.DB
}

func NewApp(conf *koanf.Koanf, router *chi.Mux, http *http.Server, db *gorm.DB) *App {
	return &App{
		conf:   conf,
		router: router,
		http:   http,
		db:     db,
	}
}

func (r *App) Run() error {
	return r.http.ListenAndServe()
}
