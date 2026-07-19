package service

import (
	"net/http"

	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/pkg/transport"
	"github.com/libtnb/chi-skeleton/internal/user/biz"
)

func UserRoutes(i do.Injector) (transport.Endpoints, error) {
	user := do.MustInvoke[*UserService](i)

	return transport.Endpoints{
		{Method: http.MethodGet, Path: "/users", Handler: user.List,
			Summary: "List users", Tags: []string{"user"},
			Request: transport.Paginate{}, Response: transport.Envelope[transport.Page[*biz.User]]{}},
		{Method: http.MethodPost, Path: "/users", Handler: user.Create,
			Summary: "Create a user", Tags: []string{"user"},
			Request: UserAdd{}, Response: transport.Envelope[biz.User]{}},
		{Method: http.MethodGet, Path: "/users/{id}", Handler: user.Get,
			Summary: "Get a user", Tags: []string{"user"},
			Request: UserID{}, Response: transport.Envelope[biz.User]{}},
		{Method: http.MethodPut, Path: "/users/{id}", Handler: user.Update,
			Summary: "Update a user", Tags: []string{"user"},
			Request: UserUpdate{}, Response: transport.Envelope[biz.User]{}},
		{Method: http.MethodDelete, Path: "/users/{id}", Handler: user.Delete,
			Summary: "Delete a user", Tags: []string{"user"},
			Request: UserID{}},
	}, nil
}
