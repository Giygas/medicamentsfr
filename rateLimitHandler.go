package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/juju/ratelimit"
)

// Per-client rate limiting
type RateLimiter struct {
	clients map[string]*ratelimit.Bucket
	mu      sync.RWMutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*ratelimit.Bucket),
	}
}

func (rl *RateLimiter) getBucket(clientIP string) *ratelimit.Bucket {
	rl.mu.RLock()
	bucket, exists := rl.clients[clientIP]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check pattern
		if bucket, exists = rl.clients[clientIP]; !exists {
			// Create bucket: 30 tokens per second, max 1000 tokens
			bucket = ratelimit.NewBucketWithRate(30, 1000)
			rl.clients[clientIP] = bucket
		}
		rl.mu.Unlock()
	}

	return bucket
}

// Clean up old clients periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			// Remove clients that haven't been used recently
			// TODO: better cleanup
			for ip, bucket := range rl.clients {
				if bucket.Available() == bucket.Capacity() {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
}

var globalRateLimiter = NewRateLimiter()

func init() {
	globalRateLimiter.cleanup()
}

func getTokenCost(r *http.Request) int64 {
	switch r.URL.Path {
	case "/database":
		return 500 // Higher cost for full database
	// TODO: better bucket costs
	// case "/database/":
	// 	return 100 // Medium cost for paged results
	// case "/medicament/":
	// 	return 50 // Medium cost for search
	default:
		return 20 // Default cost for specific lookups
	}
}

func rateLimitHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the client has a valid cached version of the data
		if r.Header.Get("If-None-Match") != "" || r.Header.Get("If-Modified-Since") != "" {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		// Get client IP (consider X-Forwarded-For if behind proxy)
		// TODO: ask if this is good if behind nginx server
		clientIP := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		}

		bucket := globalRateLimiter.getBucket(clientIP)

		// Calculate the token cost for the request
		tokenCost := getTokenCost(r)

		// Check if the client has enough tokens
		if bucket.TakeAvailable(tokenCost) < tokenCost {
			w.Header().Set("X-RateLimit-Limit", "30")
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		// Add rate limit headers
		w.Header().Set("X-RateLimit-Limit", "30")
		w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(bucket.Available(), 10))

		// Serve the request
		h.ServeHTTP(w, r)
	})
}
