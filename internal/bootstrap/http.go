package bootstrap

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"

	"github.com/go-rat/chi-skeleton/internal/http/middleware"
	"github.com/go-rat/chi-skeleton/internal/route"
)

func NewRouter(http *route.Http) (*chi.Mux, error) {
	r := chi.NewRouter()

	// add middleware
	r.Use(middleware.GlobalMiddleware(r)...)
	// add http route
	http.Register(r)

	return r, nil
}

func NewHttp(conf *koanf.Koanf, r *chi.Mux) (*http.Server, error) {
	server := &http.Server{
		Addr:           conf.MustString("http.address"),
		Handler:        http.AllowQuerySemicolons(r),
		MaxHeaderBytes: conf.MustInt("http.headerLimit") << 10,
	}

	return server, nil
}
