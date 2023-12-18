package main

import (
	"net/http"

	"github.com/juju/ratelimit"
)

func rateLimitHandler(h http.Handler) http.Handler {

	//Create a bucket with a capacity of 1000 tokens and a replenishement rate of 30 per second
	limiter := ratelimit.NewBucketWithRate(30, 1000)

	// Define a token cost function that evaluates the cost of each request
	tokenCostFunc := func(r *http.Request) int64 {
		// Calculate the token cost based on the request properties
		// Return the token cost for the request
		if r.URL.Path == "/database" {
			return 500 // Higher token cost for expensive requests
		}
		return 20 // Default token cost for other requests
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check if the client has a valid cached version of the data
		if r.Header.Get("If-None-Match") != "" || r.Header.Get("If-Modified-Since") != "" {
			// If so, skip the rate limit check and proceed to the next handler
			h.ServeHTTP(w, r)
			return
		}

		// Calculate the token cost for the request
		tokenCost := tokenCostFunc(r)

		// Check if the client has enough tokens
		if limiter.TakeAvailable(tokenCost) < tokenCost {
			http.Error(w, "You have exceeded your request rate limit, try again later", http.StatusTooManyRequests)
			return
		}

		// Serve the request
		h.ServeHTTP(w, r)
	})
}
