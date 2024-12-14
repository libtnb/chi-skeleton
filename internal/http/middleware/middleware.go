package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// GlobalMiddleware is a collection of global middleware that will be applied to every request.
func GlobalMiddleware(r *chi.Mux) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		middleware.SupressNotFound(r),
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Compress(5),
	}
}
