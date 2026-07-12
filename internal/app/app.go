package app

import (
	"context"
	_ "expvar" // registers /debug/vars on the default mux
	"fmt"
	"net/http"
	_ "net/http/pprof" // registers /debug/pprof on the default mux
	"time"

	"github.com/go-rio/migrate"
	"github.com/libtnb/cron"
	"github.com/libtnb/graceful"
	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/conf"
	"github.com/libtnb/chi-skeleton/internal/pkg/event"
	"github.com/libtnb/chi-skeleton/internal/pkg/registry"
)

type App struct {
	conf     *conf.Config
	server   *http.Server
	migrator *migrate.Migrator
	cron     *cron.Cron
}

func NewApp(i do.Injector) (*App, error) {
	// activate every subscriber so its handlers are on the bus before serving
	if _, err := registry.Collect[event.Subscription](i, registry.SubscriberPrefix); err != nil {
		return nil, err
	}

	return &App{
		conf:     do.MustInvoke[*conf.Config](i),
		server:   do.MustInvoke[*http.Server](i),
		migrator: do.MustInvoke[*migrate.Migrator](i),
		cron:     do.MustInvoke[*cron.Cron](i),
	}, nil
}

// Run migrates the database, then hands the lifecycle to graceful:
// SIGINT/SIGTERM drains everything, SIGHUP hot-upgrades the binary.
func (r *App) Run() error {
	if err := r.migrator.Up(context.Background()); err != nil {
		return err
	}
	fmt.Println("[DB] database migrated")

	g := graceful.New(
		graceful.WithUpgrade(),
		graceful.WithShutdownTimeout(30*time.Second),
	)
	// pprof/expvar live on http.DefaultServeMux, served on a private port
	if addr := r.conf.HTTP.DebugAddress; addr != "" {
		g.Listen("debug", addr, &http.Server{})
	}
	g.Add("cron", r.cron.Start, r.cron.Stop)
	g.Listen("http", r.conf.HTTP.Address, r.server)

	fmt.Println("[HTTP] listening and serving on", r.conf.HTTP.Address)
	return g.Run()
}
