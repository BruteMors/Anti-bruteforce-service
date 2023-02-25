package service

import (
	"golang.org/x/time/rate"
	"time"
)

type RateLimiterWithLastEventTime struct {
	rate      *rate.Limiter
	LastEvent time.Time
}

func NewLimiter(r rate.Limit, b int) *RateLimiterWithLastEventTime {
	limiter := rate.NewLimiter(r, b)
	return &RateLimiterWithLastEventTime{rate: limiter}
}

func (t *RateLimiterWithLastEventTime) Allow() bool {
	t.LastEvent = time.Now()
	allow := t.rate.Allow()
	return allow
}
