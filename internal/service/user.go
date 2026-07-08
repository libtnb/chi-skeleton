package service

import (
	"net/http"

	"github.com/libtnb/chi-skeleton/internal/biz"
	"github.com/libtnb/chi-skeleton/internal/request"
)

// UserService adapts HTTP to the user usecase: bind, validate, delegate, respond.
type UserService struct {
	user *biz.UserUsecase
}

func NewUserService(user *biz.UserUsecase) *UserService {
	return &UserService{
		user: user,
	}
}

func (r *UserService) List(w http.ResponseWriter, req *http.Request) {
	paginate, err := Bind[request.Paginate](req)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	users, total, err := r.user.List(req.Context(), paginate.Page, paginate.Limit)
	if err != nil {
		ErrorFrom(w, req, err)
		return
	}

	Success(w, Page[*biz.User]{
		Total: total,
		Items: users,
	})
}

func (r *UserService) Get(w http.ResponseWriter, req *http.Request) {
	userID, err := Bind[request.UserID](req)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	user, err := r.user.Get(req.Context(), userID.ID)
	if err != nil {
		ErrorFrom(w, req, err)
		return
	}

	Success(w, user)
}

func (r *UserService) Create(w http.ResponseWriter, req *http.Request) {
	userAdd, err := Bind[request.UserAdd](req)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	user, err := r.user.Create(req.Context(), userAdd.Name)
	if err != nil {
		ErrorFrom(w, req, err)
		return
	}

	Success(w, user)
}

func (r *UserService) Update(w http.ResponseWriter, req *http.Request) {
	userUpdate, err := Bind[request.UserUpdate](req)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	user, err := r.user.Update(req.Context(), userUpdate.ID, userUpdate.Name)
	if err != nil {
		ErrorFrom(w, req, err)
		return
	}

	Success(w, user)
}

func (r *UserService) Delete(w http.ResponseWriter, req *http.Request) {
	userID, err := Bind[request.UserID](req)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = r.user.Delete(req.Context(), userID.ID); err != nil {
		ErrorFrom(w, req, err)
		return
	}

	Success(w, nil)
}
