package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-rat/chi-skeleton/internal/app"
)

// GlobalMiddleware is a collection of global middleware that will be applied to every request.
func GlobalMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		middleware.SupressNotFound(app.Http),
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Compress(5),
	}
}
