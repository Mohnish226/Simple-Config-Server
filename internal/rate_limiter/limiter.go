package rate_limiter

import (
	"sync"

	"golang.org/x/time/rate"
)

var (
	rateLimiters = make(map[string]*rate.Limiter)
	mu           sync.RWMutex
)

func GetRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	if limiter, exists := rateLimiters[ip]; exists {
		return limiter
	}
	limiter := rate.NewLimiter(5, 5) // Allow 5 requests per second per IP
	rateLimiters[ip] = limiter
	return limiter
}
