package apierror

import (
	"errors"
	"fmt"
	"net/http"
)

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *ApiError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *ApiError) Unwrap() error {
	return e.Err
}

func New(code int, message string, err error) *ApiError {
	return &ApiError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func BadRequest(message string, err error) *ApiError {
	if message == "" {
		message = "Bad Request"
	}
	return New(http.StatusBadRequest, message, err)
}

func NotFound(message string, err error) *ApiError {
	if message == "" {
		message = "Requested resource not found"
	}
	return New(http.StatusNotFound, message, err)
}

func Unauthorized(message string, err error) *ApiError {
	if message == "" {
		message = "Unauthorized"
	}
	return New(http.StatusUnauthorized, message, err)
}

func Forbidden(message string, err error) *ApiError {
	if message == "" {
		message = "Forbidden"
	}
	return New(http.StatusForbidden, message, err)
}

func Conflict(message string, err error) *ApiError {
	if message == "" {
		message = "Resource conflict"
	}
	return New(http.StatusConflict, message, err)
}

func Internal(message string, err error) *ApiError {
	if message == "" {
		message = "Internal Server Error"
	}
	return New(http.StatusInternalServerError, message, err)
}

func FromError(err error) (int, string) {
	var apiErr *ApiError

	if errors.As(err, &apiErr) {
		return apiErr.Code, apiErr.Message
	}

	return http.StatusInternalServerError, "An unexpected error occurred"
}
