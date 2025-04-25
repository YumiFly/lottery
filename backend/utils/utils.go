package utils

import (
	"fmt"
	"math/big"
	"strings"
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

// parseBetContent 解析投注内容
func ParseBetContent(content string) []*big.Int {
	parts := strings.Split(content, ",")
	result := make([]*big.Int, 0, len(parts))
	for _, part := range parts {
		num, ok := new(big.Int).SetString(strings.TrimSpace(part), 10)
		if ok {
			result = append(result, num)
		}
	}
	return result
}

// parseBetContentV2bigIntSlice 解析投注内容（重命名版本）
func ParseBetContentV2bigIntSlice(content string) []*big.Int {
	return ParseBetContent(content)
}

func RemoveDuplicateString(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
