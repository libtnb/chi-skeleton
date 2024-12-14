package app

import (
	"log/slog"

	"github.com/bddjr/hlfhr"
	"github.com/go-chi/chi/v5"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/go-rat/sessions"
	"github.com/knadh/koanf/v2"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type App struct {
	conf     *koanf.Koanf
	router   *chi.Mux
	http     *hlfhr.Server
	db       *gorm.DB
	migrator *gormigrate.Gormigrate
	cron     *cron.Cron
	session  *sessions.Manager
	log      *slog.Logger
}

func NewApp(conf *koanf.Koanf, router *chi.Mux, http *hlfhr.Server, db *gorm.DB, migrator *gormigrate.Gormigrate, cron *cron.Cron, session *sessions.Manager, log *slog.Logger) *App {
	return &App{
		conf:     conf,
		router:   router,
		http:     http,
		db:       db,
		migrator: migrator,
		cron:     cron,
		session:  session,
		log:      log,
	}
}

func (r *App) Run() error {
	// migrate database
	if err := r.migrator.Migrate(); err != nil {
		return err
	}

	// start cron
	r.cron.Start()

	// start http server
	return r.http.ListenAndServe()
}
