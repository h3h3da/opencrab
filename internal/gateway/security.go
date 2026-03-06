package gateway

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

// checkOriginLoopback only allows loopback origins (security).
func checkOriginLoopback(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true // Same-origin or no Origin header
	}
	// Only allow loopback origins
	return strings.HasPrefix(origin, "http://127.0.0.1") ||
		strings.HasPrefix(origin, "http://localhost") ||
		strings.HasPrefix(origin, "http://[::1]")
}

// constantTimeCompare prevents timing attacks on auth tokens.
func constantTimeCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}

// rateLimiter limits requests per IP.
type rateLimiter struct {
	limiters map[string]*rate.Limiter
	rate     rate.Limit
	burst    int
}

func newRateLimiter(reqPerMin float64, burst int) *rateLimiter {
	r := rate.Limit(reqPerMin / 60)
	return &rateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    burst,
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	// Extract IP (strip port if present)
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	// Simple in-memory limiter; for production consider redis
	lim, ok := rl.limiters[ip]
	if !ok {
		lim = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = lim
		// TODO: periodic cleanup of old entries
	}
	return lim.Allow()
}

// CleanupStaleLimiters removes old rate limiter entries (call periodically).
func (rl *rateLimiter) CleanupStaleLimiters(olderThan time.Duration) {
	// Placeholder for periodic cleanup
	_ = olderThan
}
