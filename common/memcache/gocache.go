package memcache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var Cache *cache.Cache

func InitCache() {
	Cache = cache.New(time.Minute*5, time.Minute*20)
}
