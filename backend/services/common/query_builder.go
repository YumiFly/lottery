package common

import "gorm.io/gorm"

// QueryOption 定义查询条件的函数式选项
type QueryOption func(*gorm.DB) *gorm.DB

// ApplyQueryOptions 应用查询选项
func ApplyQueryOptions(query *gorm.DB, options []QueryOption) *gorm.DB {
	for _, opt := range options {
		query = opt(query)
	}
	return query
}
