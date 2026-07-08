package bootstrap

import (
	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/config"
	"github.com/libtnb/chi-skeleton/internal/middleware"
)

// Package wires the infrastructure.
var Package = do.Package(
	do.Lazy(func(i do.Injector) (*config.Config, error) { return config.Load() }),
	do.Lazy(NewLogger),
	do.Lazy(NewSlog),
	do.Lazy(NewSession),
	do.Lazy(middleware.NewMiddlewares),
	do.Lazy(NewValidator),
	do.Lazy(NewRouter),
	do.Lazy(NewHttp),
	do.Lazy(NewMigrate),
	do.Lazy(NewCron),
	do.Lazy(NewCrypter),
	do.Lazy(NewCli),
)
