package service

import (
	"net/http"

	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/order/biz"
	"github.com/libtnb/chi-skeleton/internal/pkg/transport"
	"github.com/libtnb/chi-skeleton/internal/server"
)

func OrderRoutes(i do.Injector) (server.Endpoints, error) {
	order := do.MustInvoke[*OrderService](i)

	return server.Endpoints{
		{Method: http.MethodGet, Path: "/orders", Handler: order.List,
			Summary: "List orders", Tags: []string{"order"},
			Request: transport.Paginate{}, Response: transport.Envelope[transport.Page[*biz.Order]]{}},
		{Method: http.MethodPost, Path: "/orders", Handler: order.Create,
			Summary: "Place an order", Tags: []string{"order"},
			Request: OrderCreate{}, Response: transport.Envelope[biz.Order]{}},
		{Method: http.MethodGet, Path: "/orders/{id}", Handler: order.Get,
			Summary: "Get an order", Tags: []string{"order"},
			Request: OrderID{}, Response: transport.Envelope[biz.Order]{}},
		{Method: http.MethodDelete, Path: "/orders/{id}", Handler: order.Delete,
			Summary: "Delete an order", Tags: []string{"order"},
			Request: OrderID{}},
	}, nil
}
