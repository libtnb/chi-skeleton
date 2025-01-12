package service

import (
	"net/http"

	"github.com/go-rat/chi-skeleton/internal/biz"
	"github.com/go-rat/chi-skeleton/internal/http/request"
)

type UserService struct {
	user biz.UserRepo
}

func NewUserService(user biz.UserRepo) *UserService {
	return &UserService{
		user: user,
	}
}

func (s *UserService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.Paginate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
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
}

func (s *UserService) Get(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.UserID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	user, err := s.user.Get(req.ID)
	if err != nil {
		ErrorSystem(w)
		return
	}

	Success(w, user)
}

func (s *UserService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.UserAdd](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	user := new(biz.User)
	user.Name = req.Name
	if err = s.user.Save(user); err != nil {
		ErrorSystem(w)
		return
	}

	Success(w, user)
}

func (s *UserService) Update(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.UserUpdate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
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
}

func (s *UserService) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.UserID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.user.Delete(req.ID); err != nil {
		ErrorSystem(w)
		return
	}

	Success(w, nil)
}
