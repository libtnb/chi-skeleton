package route

import (
	"net/http"

	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/biz"
	"github.com/libtnb/chi-skeleton/internal/request"
	"github.com/libtnb/chi-skeleton/internal/service"
)

// HealthRoutes contributes the hello and probe endpoints; none of them are
// documented, so they carry no Request/Response samples.
func HealthRoutes(i do.Injector) (Endpoints, error) {
	health := do.MustInvoke[*service.HealthService](i)

	return Endpoints{
		{Method: http.MethodGet, Path: "/", Handler: func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Hello, World 👋!"))
		}},
		{Method: http.MethodGet, Path: "/healthz", Handler: health.Healthz},
		{Method: http.MethodGet, Path: "/readyz", Handler: health.Readyz},
	}, nil
}

// UserRoutes contributes the user endpoints.
func UserRoutes(i do.Injector) (Endpoints, error) {
	user := do.MustInvoke[*service.UserService](i)

	return Endpoints{
		{Method: http.MethodGet, Path: "/users", Handler: user.List,
			Summary: "List users", Tags: []string{"user"},
			Request: request.Paginate{}, Response: service.Envelope[service.Page[*biz.User]]{}},
		{Method: http.MethodPost, Path: "/users", Handler: user.Create,
			Summary: "Create a user", Tags: []string{"user"},
			Request: request.UserAdd{}, Response: service.Envelope[biz.User]{}},
		{Method: http.MethodGet, Path: "/users/{id}", Handler: user.Get,
			Summary: "Get a user", Tags: []string{"user"},
			Request: request.UserID{}, Response: service.Envelope[biz.User]{}},
		{Method: http.MethodPut, Path: "/users/{id}", Handler: user.Update,
			Summary: "Update a user", Tags: []string{"user"},
			Request: request.UserUpdate{}, Response: service.Envelope[biz.User]{}},
		{Method: http.MethodDelete, Path: "/users/{id}", Handler: user.Delete,
			Summary: "Delete a user", Tags: []string{"user"},
			Request: request.UserID{}},
	}, nil
}
