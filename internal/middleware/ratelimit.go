package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]time.Time
	interval time.Duration
}

func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]time.Time),
		interval: interval,
	}
}

func (rl *RateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		rl.mu.Lock()
		last, exists := rl.visitors[ip]
		if exists && time.Since(last) < rl.interval {
			rl.mu.Unlock()
			http.Error(w, "Bạn đang gửi request quá nhanh, vui lòng từ từ!", http.StatusTooManyRequests)
			return
		}
		rl.visitors[ip] = time.Now()
		rl.mu.Unlock()
		next(w, r)
	}
}
