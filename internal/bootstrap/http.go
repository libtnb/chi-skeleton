package bootstrap

import (
	"net/http"

	"github.com/bddjr/hlfhr"
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"

	"github.com/libtnb/chi-skeleton/internal/http/middleware"
	"github.com/libtnb/chi-skeleton/internal/route"
)

func NewRouter(middlewares *middleware.Middlewares, http *route.Http, ws *route.Ws) (*chi.Mux, error) {
	r := chi.NewRouter()

	// add middleware
	r.Use(middlewares.Globals(r)...)
	// add http route
	http.Register(r)
	// add ws route
	ws.Register(r)

	return r, nil
}

func NewHttp(conf *koanf.Koanf, r *chi.Mux) (*hlfhr.Server, error) {
	srv := hlfhr.New(&http.Server{
		Addr:           conf.MustString("http.address"),
		Handler:        http.AllowQuerySemicolons(r),
		MaxHeaderBytes: 2048 << 20,
	})
	srv.Listen80RedirectTo443 = true

	return srv, nil
}
