package request

import "net/http"

type UserID struct {
	ID uint `uri:"id" validate:"required,number"`
}

func (r *UserID) PrepareForValidation(_ *http.Request) error {
	return nil
}

type AddUser struct {
	Name string `json:"name" form:"name" validate:"required,min=3,max=255" comment:"用户名"`
}

func (r *AddUser) PrepareForValidation(_ *http.Request) error {
	return nil
}

type UpdateUser struct {
	ID   uint   `uri:"id" validate:"required,number"`
	Name string `json:"name" form:"name" validate:"required,min=3,max=255" comment:"用户名"`
}

func (r *UpdateUser) PrepareForValidation(_ *http.Request) error {
	return nil
}
