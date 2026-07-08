package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/libtnb/chix"

	"github.com/libtnb/chi-skeleton/internal/biz"
)

// SuccessResponse is the envelope for successful responses.
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

// ErrorResponse is the envelope for error responses.
type ErrorResponse struct {
	Msg string `json:"msg"`
}

// Success writes data in the success envelope.
func Success(w http.ResponseWriter, data any) {
	render := chix.NewRender(w)
	defer render.Release()
	render.JSON(&SuccessResponse{
		Msg:  "success",
		Data: data,
	})
}

// Error writes a formatted message with the given status code.
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

// ErrorFrom maps business errors to HTTP responses: not-found becomes 404,
// anything else is logged and answered as a 500 without leaking details.
func ErrorFrom(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, biz.ErrNotFound) {
		Error(w, http.StatusNotFound, "%v", err)
		return
	}

	slog.ErrorContext(r.Context(), "request failed",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Any("err", err),
	)
	ErrorSystem(w)
}
