package service

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/go-rat/chi-skeleton/internal/biz"
	"github.com/go-rat/chi-skeleton/internal/http/request"
)

type UserService struct {
	validator *validator.Validate
	user      biz.UserRepo
}

func NewUserService(validator *validator.Validate, user biz.UserRepo) *UserService {
	return &UserService{
		validator: validator,
		user:      user,
	}
}

func (s *UserService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.Paginate](r, s.validator)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	users, total, err := s.user.List(req.Page, req.Limit)
	if err != nil {
		ErrorSystem(w)
		return
	}

	Success(w, map[string]any{
		"total": total,
		"items": users,
	})
	return
}

func (s *UserService) Get(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.UserID](r, s.validator)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := s.user.Get(req.ID)
	if err != nil {
		ErrorSystem(w)
		return
	}

	Success(w, user)
	return
}

func (s *UserService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.AddUser](r, s.validator)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user := new(biz.User)
	user.Name = req.Name
	if err = s.user.Save(user); err != nil {
		ErrorSystem(w)
		return
	}

	Success(w, user)
	return
}

func (s *UserService) Update(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.UpdateUser](r, s.validator)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user := new(biz.User)
	user.ID = req.ID
	user.Name = req.Name
	if err = s.user.Save(user); err != nil {
		ErrorSystem(w)
		return
	}

	Success(w, user)
	return
}

func (s *UserService) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.UserID](r, s.validator)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.user.Delete(req.ID); err != nil {
		ErrorSystem(w)
		return
	}

	Success(w, nil)
	return
}
