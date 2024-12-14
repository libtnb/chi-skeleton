package middleware

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-rat/sessions"
	sessionmiddleware "github.com/go-rat/sessions/middleware"
	"github.com/golang-cz/httplog"
)

// GlobalMiddleware is a collection of global middleware that will be applied to every request.
func GlobalMiddleware(r *chi.Mux, log *slog.Logger, session *sessions.Manager) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		//middleware.SupressNotFound(r),// bug https://github.com/go-chi/chi/pull/940
		sessionmiddleware.StartSession(session),
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.Compress(5),
		middleware.RequestID,
		middleware.RealIP,
		httplog.RequestLogger(log, &httplog.Options{
			Level:             slog.LevelInfo,
			LogRequestHeaders: []string{"User-Agent"},
		}),
		middleware.Recoverer,
	}
}
