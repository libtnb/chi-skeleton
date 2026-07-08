package middleware

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-cz/httplog"
	"github.com/libtnb/sessions"
	sessionmiddleware "github.com/libtnb/sessions/middleware"
	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/config"
)

type Middlewares struct {
	conf    *config.Config
	log     *slog.Logger
	session *sessions.Manager
}

func NewMiddlewares(i do.Injector) (*Middlewares, error) {
	return &Middlewares{
		conf:    do.MustInvoke[*config.Config](i),
		log:     do.MustInvoke[*slog.Logger](i),
		session: do.MustInvoke[*sessions.Manager](i),
	}, nil
}

// Globals is a collection of global middleware that will be applied to every request.
func (r *Middlewares) Globals(mux *chi.Mux) []func(http.Handler) http.Handler {
	handlers := []func(http.Handler) http.Handler{
		chimiddleware.Recoverer,
		chimiddleware.RequestSize(int64(r.conf.HTTP.BodyLimit) << 10),
		chimiddleware.StripSlashes,
		// chi's SupressNotFound is deliberately absent: it matches the raw
		// URL path and method, so it breaks StripSlashes (trailing-slash
		// requests 404) and answers every 405 with a 404.
		chimiddleware.RequestID,
		// middleware.RealIP is deliberately absent: it blindly trusts
		// X-Forwarded-For and is spoofable (GHSA-3fxj-6jh8-hvhx). Behind a
		// trusted proxy, add a middleware that only honors headers set by it.
	}

	// CORS only when origins are explicitly allowed; empty = same-origin
	if len(r.conf.HTTP.CorsOrigins) > 0 {
		handlers = append(handlers, cors.Handler(cors.Options{
			AllowedOrigins: r.conf.HTTP.CorsOrigins,
		}))
	}

	return append(handlers,
		httplog.RequestLogger(r.log, &httplog.Options{
			Level:             slog.LevelInfo,
			LogRequestHeaders: []string{"User-Agent"},
			// probes are noise
			Skip: func(req *http.Request, respStatus int) bool {
				return isProbe(req.URL.Path)
			},
		}),
		requestIDAttr,
		chimiddleware.Compress(5),
		// probes carry no cookies; a session per hit would be created,
		// persisted and garbage-collected for nothing
		skipProbes(sessionmiddleware.StartSession(r.session)),
	)
}

func isProbe(path string) bool {
	return path == "/healthz" || path == "/readyz"
}

// skipProbes bypasses mw for probe requests.
func skipProbes(mw func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		wrapped := mw(next)
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if isProbe(req.URL.Path) {
				next.ServeHTTP(w, req)
				return
			}
			wrapped.ServeHTTP(w, req)
		})
	}
}

// requestIDAttr returns the request id in X-Request-Id and attaches it to
// the access log record.
func requestIDAttr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if id := chimiddleware.GetReqID(req.Context()); id != "" {
			w.Header().Set("X-Request-Id", id)
			httplog.SetAttrs(req.Context(), slog.String("request_id", id))
		}
		next.ServeHTTP(w, req)
	})
}
