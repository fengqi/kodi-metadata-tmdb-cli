package memcache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Cache *cache.Cache

func InitCache() {
	Cache = cache.New(time.Minute*5, time.Minute*20)
}
