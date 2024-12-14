package bootstrap

import (
	"log/slog"
	"net/http"

	"github.com/bddjr/hlfhr"
	"github.com/go-chi/chi/v5"
	"github.com/go-rat/sessions"
	"github.com/knadh/koanf/v2"

	"github.com/go-rat/chi-skeleton/internal/http/middleware"
	"github.com/go-rat/chi-skeleton/internal/route"
)

func NewRouter(log *slog.Logger, session *sessions.Manager, http *route.Http, ws *route.Ws) (*chi.Mux, error) {
	r := chi.NewRouter()

	// add middleware
	r.Use(middleware.GlobalMiddleware(r, log, session)...)
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
	srv.HttpOnHttpsPortErrorHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hlfhr.RedirectToHttps(w, r, http.StatusTemporaryRedirect)
	})

	return srv, nil
}
