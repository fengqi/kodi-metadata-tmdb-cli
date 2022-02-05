package kodi

import (
	"sync"
	"time"
)

type Limiter struct {
	lock   *sync.Map
	second int
}

func NewLimiter(second int) *Limiter {
	return &Limiter{
		lock:   &sync.Map{},
		second: second,
	}
}

func (l *Limiter) take() bool {
	var t time.Time
	value, ok := l.lock.LoadOrStore("time", time.Time{})
	if !ok {
		t = time.Now().Add(time.Second * time.Duration(-l.second))
	} else {
		t = value.(time.Time)
	}

	if t.Add(time.Second * time.Duration(l.second)).Before(time.Now()) {
		l.lock.Store("time", time.Now())
		return true
	}

	return false
}
