package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/giygas/medicamentsfr/medicamentsparser"
	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
	"github.com/juju/ratelimit"
)

var medicos []entities.Medicament

func getDatabase(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	val, err := json.MarshalIndent(medicos, "", "  ")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		io.WriteString(w, string(err.Error()))
	}
	io.WriteString(w, string(val))
}

func main() {
	medicos = medicamentsparser.ParseAllMedicaments()
	
	//Create a bucket with a capacity of 1000 tokens and a replenishement rate of 30 per second
	limiter := ratelimit.NewBucketWithRate(30, 1000)
	
	// Define a token cost function that evaluates the cost of each request
	tokenCostFunc := func(r *http.Request) int64 {
		// Calculate the token cost based on the request properties
		// Return the token cost for the request
		if r.URL.Path == "/Medicaments.json" {
			return 1000 // Higher token cost for expensive requests
		}
		return 15 // Default token cost for other requests
	}
	
	// Handler function that enforces rate limiting
	rateLimitHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Calculate the token cost for the request
			tokenCost := tokenCostFunc(r)

			// Check if the client has enough tokens
			if limiter.TakeAvailable(tokenCost) < tokenCost {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// Serve the request
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