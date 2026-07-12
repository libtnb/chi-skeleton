// Package service adapts HTTP to the order usecase and owns the module's
// request DTOs, route contribution and event subscribers.
package service

import (
	"net/http"

	"github.com/libtnb/chi-skeleton/internal/order/biz"
	"github.com/libtnb/chi-skeleton/internal/pkg/transport"
)

// OrderService adapts HTTP to the order usecase: bind, validate, delegate, respond.
type OrderService struct {
	order *biz.OrderUsecase
}

func NewOrderService(order *biz.OrderUsecase) *OrderService {
	return &OrderService{
		order: order,
	}
}

func (r *OrderService) List(w http.ResponseWriter, req *http.Request) {
	paginate, err := transport.Bind[transport.Paginate](req)
	if err != nil {
		transport.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	orders, total, err := r.order.List(req.Context(), paginate.Page, paginate.Limit)
	if err != nil {
		transport.ErrorFrom(w, req, err)
		return
	}

	transport.Success(w, transport.Page[*biz.Order]{
		Total: total,
		Items: orders,
	})
}

func (r *OrderService) Get(w http.ResponseWriter, req *http.Request) {
	orderID, err := transport.Bind[OrderID](req)
	if err != nil {
		transport.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	order, err := r.order.Get(req.Context(), orderID.ID)
	if err != nil {
		transport.ErrorFrom(w, req, err)
		return
	}

	transport.Success(w, order)
}

func (r *OrderService) Create(w http.ResponseWriter, req *http.Request) {
	create, err := transport.Bind[OrderCreate](req)
	if err != nil {
		transport.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	order, err := r.order.Place(req.Context(), create.UserID, create.Amount)
	if err != nil {
		transport.ErrorFrom(w, req, err)
		return
	}

	transport.Success(w, order)
}

func (r *OrderService) Delete(w http.ResponseWriter, req *http.Request) {
	orderID, err := transport.Bind[OrderID](req)
	if err != nil {
		transport.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = r.order.Delete(req.Context(), orderID.ID); err != nil {
		transport.ErrorFrom(w, req, err)
		return
	}

	transport.Success(w, nil)
}
