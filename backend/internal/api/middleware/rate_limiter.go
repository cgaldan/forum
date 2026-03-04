package middleware

import (
	"net/http"
	"real-time-forum/internal/config"
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	RateLimiter := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			RateLimiter.cleanup()
		}
	}()

	return RateLimiter
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, timestamps := range rl.requests {
		var validTimestamps []time.Time
		for _, timestamp := range timestamps {
			if now.Sub(timestamp) < rl.window {
				validTimestamps = append(validTimestamps, timestamp)
			}
		}
		if len(validTimestamps) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = validTimestamps
		}
	}
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	timestamps := rl.requests[ip]

	var validTimestamps []time.Time
	for _, timestamp := range timestamps {
		if now.Sub(timestamp) < rl.window {
			validTimestamps = append(validTimestamps, timestamp)
		}
	}

	if len(validTimestamps) >= rl.limit {
		return false
	}

	validTimestamps = append(validTimestamps, now)
	rl.requests[ip] = validTimestamps
	return true
}

func RateLimiterMiddleware(config *config.Config) func(http.Handler) http.Handler {
	if !config.RateLimit.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	rateLimiter := NewRateLimiter(config.RateLimit.RequestsPerMinute, time.Minute)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)

			if !rateLimiter.allow(ip) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
