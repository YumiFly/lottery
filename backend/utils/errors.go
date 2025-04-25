// utils/errors.go
package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

const (
	ErrCodeInvalidInput     = 1001
	ErrCodeValidationFailed = 1002
	ErrCodeUnauthorized     = 1003
	ErrCodeInternalServer   = 1004
	ErrCodeForbidden        = 1005
	ErrCodeBadRequest       = 1006
)

func ErrorResponse(code int, message string, data interface{}) Response {
	return Response{
		Message: message,
		Code:    code,
		Data:    data,
	}
}

func NewErrorResponse(err error) Response {
	if customErr, ok := err.(*Error); ok {
		return Response{Message: customErr.Message, Code: customErr.Code, Data: nil}
	}
	return Response{Message: err.Error(), Code: http.StatusInternalServerError, Data: nil}
}

func NewValidator() *validator.Validate {
	return validator.New()
}

// Error 自定义错误
type Error struct {
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func NewBadRequestError(message string, err error) *Error {
	return &Error{Code: http.StatusBadRequest, Message: message, Err: err}
}

func NewInternalError(message string, err error) *Error {
	return &Error{Code: http.StatusInternalServerError, Message: message, Err: err}
}
