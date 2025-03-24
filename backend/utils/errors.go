// utils/errors.go
package utils

import "github.com/go-playground/validator/v10"

const (
	ErrCodeInvalidInput     = 1001
	ErrCodeValidationFailed = 1002
	ErrCodeUnauthorized     = 1003
	ErrCodeInternalServer   = 1004
	ErrCodeForbidden        = 1005
)

func ErrorResponse(code int, message string, data interface{}) Response {
	return Response{
		Message: message,
		Code:    code,
		Data:    data,
	}
}

func NewValidator() *validator.Validate {
	return validator.New()
}
