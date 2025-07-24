package utils

import (
	"time"
)

// CacheExpire 判断缓存到期，1年以上的永久缓存，一年到半年的一礼拜，半年内的3天
// TODO 可配置
func CacheExpire(modTime, airTime time.Time) bool {
	airSub := time.Since(airTime)
	if airSub.Hours() >= 24*365 {
		return false
	}

	cacheSub := time.Since(modTime)
	if airSub.Hours() >= 24*180 && cacheSub.Hours() > 24*7 {
		return false
	}

	return cacheSub.Hours() > 24*3
}
