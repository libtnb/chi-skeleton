package transport

import (
	"errors"
	"net/http"

	"github.com/libtnb/chix/v2"
	"github.com/libtnb/validator"
)

// Bind binds and validates the request against the validator installed via
// validator.SetDefault.
func Bind[T any](r *http.Request) (*T, error) {
	v := validator.Default()

	req := new(T)

	binder := chix.NewBind(r)
	defer binder.Release()
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		if r.ContentLength > 0 {
			if err := binder.Body(req); err != nil {
				return nil, err
			}
		}
	}
	if err := binder.Query(req); err != nil {
		return nil, err
	}
	if err := binder.URI(req); err != nil {
		return nil, err
	}

	if hook, ok := any(req).(WithPrepare); ok {
		if err := hook.Prepare(r); err != nil {
			return nil, err
		}
	}

	vd := v.Struct(req)
	if hook, ok := any(req).(WithRules); ok {
		for field, expr := range hook.Rules(r) {
			if err := vd.AddRules(field, expr); err != nil {
				return nil, err
			}
		}
	}
	if hook, ok := any(req).(WithFilters); ok {
		for field, filters := range hook.Filters(r) {
			if err := vd.AddFilters(field, filters); err != nil {
				return nil, err
			}
		}
	}
	if hook, ok := any(req).(WithMessages); ok {
		if messages := hook.Messages(r); messages != nil {
			if err := vd.AddMessages(messages); err != nil {
				return nil, err
			}
		}
	}

	vd.Validate(r.Context())
	if vd.Fails() {
		return nil, errors.New(vd.Errors().One())
	}

	// write filtered values (trim, lower, ...) back into the request struct
	if err := vd.SafeBind(req); err != nil {
		return nil, err
	}

	return req, nil
}
