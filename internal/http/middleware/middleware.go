package middleware

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/libtnb/sessions"
	sessionmiddleware "github.com/libtnb/sessions/middleware"
	"github.com/golang-cz/httplog"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewMiddlewares)

type Middlewares struct {
	log     *slog.Logger
	session *sessions.Manager
}

func NewMiddlewares(log *slog.Logger, session *sessions.Manager) *Middlewares {
	return &Middlewares{
		log:     log,
		session: session,
	}
}

// Globals is a collection of global middleware that will be applied to every request.
func (r *Middlewares) Globals(mux *chi.Mux) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		middleware.Recoverer,
		//middleware.SupressNotFound(r),// bug https://github.com/go-chi/chi/pull/940
		middleware.StripSlashes,
		middleware.RequestID,
		middleware.RealIP,
		httplog.RequestLogger(r.log, &httplog.Options{
			Level:             slog.LevelInfo,
			LogRequestHeaders: []string{"User-Agent"},
		}),
		middleware.Compress(5),
		sessionmiddleware.StartSession(r.session),
	}
}
