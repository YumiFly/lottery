// utils/cache.go
package utils

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var Cache *cache.Cache

func InitCache() {
	// 初始化缓存，默认 TTL 为 5 分钟，清理间隔 10 分钟
	Cache = cache.New(5*time.Minute, 10*time.Minute)
}
