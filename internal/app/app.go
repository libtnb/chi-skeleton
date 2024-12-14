package app

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-rat/sessions"
	"github.com/knadh/koanf/v2"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type App struct {
	conf    *koanf.Koanf
	router  *chi.Mux
	http    *http.Server
	db      *gorm.DB
	cron    *cron.Cron
	session *sessions.Manager
	log     *slog.Logger
}

func NewApp(conf *koanf.Koanf, router *chi.Mux, http *http.Server, db *gorm.DB, cron *cron.Cron, session *sessions.Manager, log *slog.Logger) *App {
	return &App{
		conf:    conf,
		router:  router,
		http:    http,
		db:      db,
		cron:    cron,
		session: session,
		log:     log,
	}
}

func (r *App) Run() error {
	return r.http.ListenAndServe()
}
