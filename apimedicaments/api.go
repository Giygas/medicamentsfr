package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/giygas/medicamentsfr/medicamentsparser"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	// Load the env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the evironment")
	}

	medicamentsparser.ParseAllMedicaments()

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(rateLimitHandler)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	router.Get("/database", serveAllMedicaments)

	fmt.Printf("Starting server at PORT: %v\n", portString)
	err = server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
