package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/go-rat/chi-skeleton/internal/service"
)

type Http struct {
	user *service.UserService
}

func NewHttp(user *service.UserService) *Http {
	return &Http{
		user: user,
	}
}

func (r *Http) Register(router *chi.Mux) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World 👋!"))
	})

	router.Get("/users", r.user.List)
	router.Post("/users", r.user.Create)
	router.Get("/users/:id", r.user.Get)
	router.Put("/users/:id", r.user.Update)
	router.Delete("/users/:id", r.user.Delete)
}
