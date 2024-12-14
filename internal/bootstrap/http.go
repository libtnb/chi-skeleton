package bootstrap

import (
	"net/http"

	"github.com/bddjr/hlfhr"
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

func NewHttp(conf *koanf.Koanf, r *chi.Mux) (*hlfhr.Server, error) {
	srv := hlfhr.New(&http.Server{
		Addr:           conf.MustString("http.address"),
		Handler:        http.AllowQuerySemicolons(r),
		MaxHeaderBytes: 2048 << 20,
	})
	srv.HttpOnHttpsPortErrorHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hlfhr.RedirectToHttps(w, r, http.StatusTemporaryRedirect)
	})

	return srv, nil
}
