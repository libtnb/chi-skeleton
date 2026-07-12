package app

import (
	"time"

	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/bootstrap"
	"github.com/libtnb/chi-skeleton/internal/conf"
	"github.com/libtnb/chi-skeleton/internal/order"
	"github.com/libtnb/chi-skeleton/internal/server"
	"github.com/libtnb/chi-skeleton/internal/user"
)

// NewInjector assembles every package of the application.
func NewInjector() do.Injector {
	return do.NewWithOpts(&do.InjectorOpts{
		// keeps /readyz bounded even when a dependency hangs
		HealthCheckGlobalTimeout: 5 * time.Second,
	},
		do.Lazy(func(i do.Injector) (*conf.Config, error) { return conf.Load() }),

		// boot-time infrastructure and the HTTP server
		bootstrap.Package,
		server.Package,

		// business modules
		user.Package,
		order.Package,

		// application lifecycle
		do.Lazy(newRootCommand),
		do.Lazy(NewApp),
		do.Lazy(NewCli),
	)
}
