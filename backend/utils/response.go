// utils/response.go
package utils

import "net/http"

// Response 统一响应结构体
type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

// SuccessResponse 返回成功响应
func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Message: message,
		Code:    http.StatusOK,
		Data:    data,
	}
}
