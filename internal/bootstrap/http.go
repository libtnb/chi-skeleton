package bootstrap

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/libtnb/validator"
	"github.com/libtnb/validator/contrib/openapi"
	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/config"
	"github.com/libtnb/chi-skeleton/internal/middleware"
	"github.com/libtnb/chi-skeleton/internal/route"
	"github.com/libtnb/chi-skeleton/internal/service"
)

func NewRouter(i do.Injector) (*chi.Mux, error) {
	middlewares := do.MustInvoke[*middleware.Middlewares](i)

	// handlers reach this instance through service.Bind / validator.Default
	validator.SetDefault(do.MustInvoke[*validator.Validator](i))

	r := chi.NewRouter()
	r.Use(middlewares.Globals(r)...)

	if err := route.HTTP(i, r); err != nil {
		return nil, err
	}

	conf := do.MustInvoke[*config.Config](i)
	if conf.HTTP.Docs {
		spec, err := route.SpecJSON(i, conf.App.Name)
		if err != nil {
			return nil, err
		}
		docs := openapi.DocsHTML(conf.App.Name, "/openapi.json")
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
		service.Error(w, http.StatusNotFound, "%s", http.StatusText(http.StatusNotFound))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		service.Error(w, http.StatusMethodNotAllowed, "%s", http.StatusText(http.StatusMethodNotAllowed))
	})

	return r, nil
}

func NewHttp(i do.Injector) (*http.Server, error) {
	conf := do.MustInvoke[*config.Config](i)

	return &http.Server{
		Addr:           conf.HTTP.Address,
		Handler:        http.AllowQuerySemicolons(do.MustInvoke[*chi.Mux](i)),
		MaxHeaderBytes: conf.HTTP.HeaderLimit,
		ReadTimeout:    conf.HTTP.ReadTimeout,
		WriteTimeout:   conf.HTTP.WriteTimeout,
		IdleTimeout:    conf.HTTP.IdleTimeout,
	}, nil
}
