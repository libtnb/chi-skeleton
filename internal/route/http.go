package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/go-rat/chi-skeleton/internal/service"
)

func Http(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World 👋!"))
	})

	user := service.NewUserService()
	r.Get("/users", user.List)
	r.Post("/users", user.Create)
	r.Get("/users/:id", user.Get)
	r.Put("/users/:id", user.Update)
	r.Delete("/users/:id", user.Delete)
}
