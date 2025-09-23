package main

import (
	"net"
	"net/http"
	"strings"
)

func getRealClientIP(r *http.Request) string {
	// Check X-Forwarded-For first (most common)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs: "client, proxy1, proxy2"
		// Take the first one (original client)
		ips := strings.Split(xff, ",")
		clientIP := strings.TrimSpace(ips[0])

		// Validate it's a real IP
		if net.ParseIP(clientIP) != nil {
			return clientIP
		}
	}

	// Fallback to X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if net.ParseIP(xri) != nil {
			return xri
		}
	}

	// Last resort: extract from RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // Return as-is if can't parse
	}
	return host
}

// Middleware to get real IP from nginx
func realIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the real client IP from nginx headers
		realIP := getRealClientIP(r)

		// Replace the RemoteAddr with the real IP
		r.RemoteAddr = realIP + ":0" // Port doesn't matter for logging

		next.ServeHTTP(w, r)
	})
}

// Middleware to block direct access to the app
func blockDirectAccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request is coming through nginx
		if r.Header.Get("X-Forwarded-For") == "" && r.Header.Get("X-Real-IP") == "" {
			// Direct access without nginx headers
			host, _, _ := net.SplitHostPort(r.RemoteAddr)

			// Allow localhost connections (for health checks, etc.)
			if host != "127.0.0.1" && host != "::1" && host != "localhost" {
				http.Error(w, "Direct access not allowed", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
