package errorx

import (
	"fmt"
	"net/http"

	"github.com/DrReMain/cgin/pkg/encoding/json"
)

type ErrorX struct {
	ID     string `json:"id"`
	Code   int32  `json:"code"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

func (e *ErrorX) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// ErrorBadRequest 400
func ErrorBadRequest(id, format string, a ...interface{}) error {
	if id == "" {
		id = DefaultBadRequestID
	}
	return &ErrorX{
		ID:     id,
		Code:   http.StatusBadRequest,
		Status: http.StatusText(http.StatusBadRequest),
		Detail: fmt.Sprintf(format, a...),
	}
}

// ErrorUnauthorized 401
func ErrorUnauthorized(id, format string, a ...interface{}) error {
	if id == "" {
		id = DefaultUnauthorizedID
	}
	return &ErrorX{
		ID:     id,
		Code:   http.StatusUnauthorized,
		Status: http.StatusText(http.StatusUnauthorized),
		Detail: fmt.Sprintf(format, a...),
	}
}

// ErrorForbidden 403
func ErrorForbidden(id, format string, a ...interface{}) error {
	if id == "" {
		id = DefaultForbiddenID
	}
	return &ErrorX{
		ID:     id,
		Code:   http.StatusForbidden,
		Status: http.StatusText(http.StatusForbidden),
		Detail: fmt.Sprintf(format, a...),
	}
}

// ErrorNotFound 404
func ErrorNotFound(id, format string, args ...any) error {
	if id == "" {
		id = DefaultNotFoundID
	}
	return &ErrorX{
		ID:     id,
		Code:   http.StatusNotFound,
		Status: http.StatusText(http.StatusNotFound),
		Detail: fmt.Sprintf(format, args...),
	}
}

// ErrorMethodNotAllowed 405
func ErrorMethodNotAllowed(id, format string, args ...any) error {
	if id == "" {
		id = DefaultMethodNotAllowedID
	}
	return &ErrorX{
		ID:     id,
		Code:   http.StatusMethodNotAllowed,
		Status: http.StatusText(http.StatusMethodNotAllowed),
		Detail: fmt.Sprintf(format, args...),
	}
}

// ErrorTimeout 408
func ErrorTimeout(id, format string, a ...interface{}) error {
	if id == "" {
		id = DefaultRequestTimeoutID
	}
	return &ErrorX{
		ID:     id,
		Code:   http.StatusRequestTimeout,
		Status: http.StatusText(http.StatusRequestTimeout),
		Detail: fmt.Sprintf(format, a...),
	}
}

// ErrorConflict 409
func ErrorConflict(id, format string, a ...interface{}) error {
	if id == "" {
		id = DefaultConflictID
	}
	return &ErrorX{
		ID:     id,
		Code:   http.StatusConflict,
		Status: http.StatusText(http.StatusConflict),
		Detail: fmt.Sprintf(format, a...),
	}
}

// ErrorRequestEntityTooLarge 413
func ErrorRequestEntityTooLarge(id, format string, args ...any) error {
	if id == "" {
		id = DefaultRequestEntityTooLargeID
	}
	return &ErrorX{
		ID:     id,
		Code:   http.StatusRequestEntityTooLarge,
		Status: http.StatusText(http.StatusRequestEntityTooLarge),
		Detail: fmt.Sprintf(format, args...),
	}
}

// ErrorTooManyRequests 429
func ErrorTooManyRequests(id, format string, args ...any) error {
	if id == "" {
		id = DefaultTooManyRequestsID
	}
	return &ErrorX{
		ID:     id,
		Code:   http.StatusTooManyRequests,
		Status: http.StatusText(http.StatusTooManyRequests),
		Detail: fmt.Sprintf(format, args...),
	}
}

// ErrorInternalServerError 500
func ErrorInternalServerError(id, format string, args ...any) error {
	if id == "" {
		id = DefaultInternalServerErrorID
	}
	return &ErrorX{
		ID:     id,
		Code:   http.StatusInternalServerError,
		Status: http.StatusText(http.StatusInternalServerError),
		Detail: fmt.Sprintf(format, args...),
	}
}
