package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/libtnb/validator"
	"github.com/libtnb/validator/contrib/openapi"
	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/conf"
	"github.com/libtnb/chi-skeleton/internal/pkg/registry"
	"github.com/libtnb/chi-skeleton/internal/pkg/transport"
)

// Package wires the HTTP server: the router, the http.Server around it, the
// middleware, the probes and the websocket endpoint.
var Package = do.Package(
	do.Lazy(NewMiddlewares),
	do.Lazy(NewHealthService),
	do.Lazy(NewRouter),
	do.Lazy(NewHttp),
	do.LazyNamed(registry.RoutePrefix+"health", HealthRoutes),
	do.LazyNamed(registry.RoutePrefix+"ws", WsRoutes),
)

func NewRouter(i do.Injector) (*chi.Mux, error) {
	middlewares := do.MustInvoke[*Middlewares](i)

	// handlers reach this instance through transport.Bind / validator.Default
	validator.SetDefault(do.MustInvoke[*validator.Validator](i))

	r := chi.NewRouter()
	r.Use(middlewares.Globals(r)...)

	if err := HTTP(i, r); err != nil {
		return nil, err
	}

	config := do.MustInvoke[*conf.Config](i)
	if config.HTTP.Docs {
		spec, err := SpecJSON(i, config.App.Name)
		if err != nil {
			return nil, err
		}
		docs := openapi.DocsHTML(config.App.Name, "/openapi.json")
		r.Get("/openapi.json", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(spec)
		})
		r.Get("/docs", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(docs)
		})
	}

	// framework-level errors leave as JSON in the same shape as the API
	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		transport.Error(w, http.StatusNotFound, "%s", http.StatusText(http.StatusNotFound))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		transport.Error(w, http.StatusMethodNotAllowed, "%s", http.StatusText(http.StatusMethodNotAllowed))
	})

	return r, nil
}

func NewHttp(i do.Injector) (*http.Server, error) {
	config := do.MustInvoke[*conf.Config](i)

	return &http.Server{
		Addr:           config.HTTP.Address,
		Handler:        http.AllowQuerySemicolons(do.MustInvoke[*chi.Mux](i)),
		MaxHeaderBytes: config.HTTP.HeaderLimit,
		ReadTimeout:    config.HTTP.ReadTimeout,
		WriteTimeout:   config.HTTP.WriteTimeout,
		IdleTimeout:    config.HTTP.IdleTimeout,
	}, nil
}
