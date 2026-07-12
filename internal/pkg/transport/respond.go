package transport

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-rio/rio"
	"github.com/libtnb/chix/v2"
	"github.com/samber/oops"
)

type SuccessResponse struct {
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Envelope mirrors SuccessResponse with a typed payload; route declarations
// use it to document response bodies.
type Envelope[T any] struct {
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// Page is the typed payload of list responses.
type Page[T any] struct {
	Total int64 `json:"total"`
	Items []T   `json:"items"`
}

type ErrorResponse struct {
	Msg  string `json:"msg"`
	Code string `json:"code,omitempty"`
}

func Success(w http.ResponseWriter, data any) {
	render := chix.NewRender(w)
	defer render.Release()
	render.JSON(&SuccessResponse{
		Msg:  "success",
		Data: data,
	})
}

func Error(w http.ResponseWriter, code int, format string, args ...any) {
	render := chix.NewRender(w)
	defer render.Release()
	render.Header(chix.HeaderContentType, chix.MIMEApplicationJSONCharsetUTF8) // must before Status()
	render.Status(code)
	render.JSON(&ErrorResponse{
		Msg: fmt.Sprintf(format, args...),
	})
}

// ErrorSystem writes a generic 500 without leaking details.
func ErrorSystem(w http.ResponseWriter) {
	render := chix.NewRender(w)
	defer render.Release()
	render.Header(chix.HeaderContentType, chix.MIMEApplicationJSONCharsetUTF8) // must before Status()
	render.Status(http.StatusInternalServerError)
	render.JSON(&ErrorResponse{
		Msg: http.StatusText(http.StatusInternalServerError),
	})
}

// ErrorFrom maps an error to an HTTP response. A not-found becomes 404; an oops
// error whose Code is known becomes that status and returns the error's public
// message and code; anything else is logged with its full structured context
// (stack trace, domain, attributes) and answered as a 500 without leaking it.
func ErrorFrom(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, rio.ErrNotFound) {
		Error(w, http.StatusNotFound, "not found")
		return
	}

	if oopsErr, ok := oops.AsError[oops.OopsError](err); ok {
		code, _ := oopsErr.Code().(string)
		if status := statusFromCode(code); status != 0 {
			render := chix.NewRender(w)
			defer render.Release()
			render.Header(chix.HeaderContentType, chix.MIMEApplicationJSONCharsetUTF8)
			render.Status(status)
			render.JSON(&ErrorResponse{
				Msg:  oops.GetPublic(err, http.StatusText(status)),
				Code: code,
			})
			return
		}
	}

	slog.ErrorContext(r.Context(), "request failed",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Any("error", err),
	)
	ErrorSystem(w)
}

// statusFromCode maps an application error code to an HTTP status, or 0 when
// the code is unknown — an unknown code is an unexpected error, not a
// client-facing one. Add a case here when a module introduces a new code.
func statusFromCode(code string) int {
	switch code {
	case "user.name_taken":
		return http.StatusConflict
	case "order.user_not_found":
		return http.StatusUnprocessableEntity
	default:
		return 0
	}
}
