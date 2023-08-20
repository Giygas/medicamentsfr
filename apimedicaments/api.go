package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/giygas/medicamentsfr/medicamentsparser"
	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var medicaments []entities.Medicament

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

	medicaments = medicamentsparser.ParseAllMedicaments()

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
	// Search medicaments by elementPharmaceutique or cis
	router.Get("/medicament/{cis}", findMedicament)
	// Searh medicaments by generiques libelle or generiques group
	router.Get("/generiques/{libelle}", findGeneriques)

	fmt.Printf("Starting server at PORT: %v\n", portString)
	err = server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
