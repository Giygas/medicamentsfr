package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/giygas/medicamentsfr/medicamentsparser"
	"github.com/juju/ratelimit"
)

func main() {
	
	medicamentsparser.ParseAllMedicaments()
	
	//Create a bucket with a capacity of 1000 tokens and a replenishement rate of 30 per second
	limiter := ratelimit.NewBucketWithRate(30, 1000)
	
	// Define a token cost function that evaluates the cost of each request
	tokenCostFunc := func(r *http.Request) int64 {
		// Calculate the token cost based on the request properties
		// Return the token cost for the request
		if r.URL.Path == "/Medicaments.json" {
			return 500 // Higher token cost for expensive requests
		}
		return 15 // Default token cost for other requests
	}
	
	// Handler function that enforces rate limiting
	rateLimitHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			// Check the If-Modified-Since header, so if the user has the json cached, there's no need to
			// discount tokens from bucket
			ifModifiedSince := r.Header.Get("If-Modified-Since")
			if ifModifiedSince != "" {
				// Parse the If-Modified-Since date
				ifModifiedSinceTime, err := http.ParseTime(ifModifiedSince)
				if err == nil {
					// Get the modification time of the file
					fileInfo, err := os.Stat("./src/Medicaments.json")
					if err == nil && fileInfo.ModTime().Before(ifModifiedSinceTime.Add(1*time.Second)) {
						// If the file has not been modified since the If-Modified-Since date, return a 304 Not Modified response
						w.WriteHeader(http.StatusNotModified)
						return
					}
				}
			}
			
			// Calculate the token cost for the request
			tokenCost := tokenCostFunc(r)

			// Check if the client has enough tokens
			if limiter.TakeAvailable(tokenCost) < tokenCost {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// Serve the request
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Cache-Control", "public, max-age=3600")
			w.Header().Set("Expires", time.Now().Add(time.Hour).Format(http.TimeFormat))
			w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
			h.ServeHTTP(w, r)
		})
	}	
	
	// File server handler
	fs := http.FileServer(http.Dir("./src"))

	// Apply the rate limit handler to the file server handler
	http.Handle("/Medicaments.json", rateLimitHandler(fs))
	
	fmt.Printf("Starting server at localhost:1337\n")
	
	err := http.ListenAndServe(":1337", nil)
	
  if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}