package utils

import (
	"fmt"
)

// ServiceError 自定义错误类型
type ServiceError struct {
	Msg string
	Err error
}

func (e *ServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

func NewServiceError(msg string, err error) *ServiceError {
	return &ServiceError{Msg: msg, Err: err}
}
