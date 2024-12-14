package service

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-rat/chix"

	"github.com/go-rat/chi-skeleton/internal/http/request"
)

// SuccessResponse 通用成功响应
type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// ErrorResponse 通用错误响应
type ErrorResponse struct {
	Message string `json:"message"`
}

// Success 响应成功
func Success(w http.ResponseWriter, data any) {
	render := chix.NewRender(w)
	defer render.Release()
	render.JSON(&SuccessResponse{
		Message: "success",
		Data:    data,
	})
}

// Error 响应错误
func Error(w http.ResponseWriter, code int, message string) {
	render := chix.NewRender(w)
	defer render.Release()
	render.Status(code)
	render.JSON(&ErrorResponse{
		Message: message,
	})
}

// ErrorSystem 响应系统错误
func ErrorSystem(w http.ResponseWriter) {
	render := chix.NewRender(w)
	defer render.Release()
	render.Status(http.StatusInternalServerError)
	render.JSON(&ErrorResponse{
		Message: http.StatusText(http.StatusInternalServerError),
	})
}

// Bind 验证并绑定请求参数
func Bind[T any](r *http.Request, validate *validator.Validate) (*T, error) {
	req := new(T)

	// 绑定参数
	binder := chix.NewBind(r)
	defer binder.Release()
	if err := binder.URI(req); err != nil {
		return nil, err
	}
	if err := binder.Query(req); err != nil {
		return nil, err
	}
	if slices.Contains([]string{"POST", "PUT", "PATCH"}, strings.ToUpper(r.Method)) {
		if err := binder.Body(req); err != nil {
			return nil, err
		}
	}

	// 准备验证
	if reqWithPrepare, ok := any(req).(request.WithPrepare); ok {
		if err := reqWithPrepare.Prepare(r); err != nil {
			return nil, err
		}
	}
	if reqWithAuthorize, ok := any(req).(request.WithAuthorize); ok {
		if err := reqWithAuthorize.Authorize(r); err != nil {
			return nil, err
		}
	}
	if reqWithRules, ok := any(req).(request.WithRules); ok {
		if rules := reqWithRules.Rules(r); rules != nil {
			validate.RegisterStructValidationMapRules(rules, req)
		}
	}

	// 验证参数
	err := validate.Struct(req)
	if err == nil {
		return req, nil
	}

	// 翻译错误信息
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, e := range errs {
			if reqWithMessages, ok := any(req).(request.WithMessages); ok {
				if msg, found := reqWithMessages.Messages(r)[fmt.Sprintf("%s.%s", e.Field(), e.Tag())]; found {
					return nil, errors.New(msg)
				}
			}
		}
	}

	return nil, err
}

// Paginate 取分页条目
func Paginate[T any](r request.Paginate, allItems []T) (pagedItems []T, total uint) {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Limit == 0 {
		r.Limit = 10
	}
	total = uint(len(allItems))
	startIndex := (r.Page - 1) * r.Limit
	endIndex := r.Page * r.Limit

	if total == 0 {
		return []T{}, 0
	}
	if startIndex > total {
		return []T{}, total
	}
	if endIndex > total {
		endIndex = total
	}

	return allItems[startIndex:endIndex], total
}
