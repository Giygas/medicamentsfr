package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var medicaments = make([]entities.Medicament, 0)
var generiques = make([]entities.GeneriqueList, 0)
var medicamentsMap = make(map[int]entities.Medicament)
var generiquesMap = make(map[int]entities.Generique)

func init() {
	// Load the env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go scheduleMedicaments(&medicaments, &medicamentsMap, &generiques, &generiquesMap)
}

func main() {

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the evironment")
	}

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
	router.Get("/medicament/{element}", findMedicament)
	router.Get("/medicament/id/{cis}", findMedicamentById)
	router.Get("/generiques/{libelle}", findGeneriques)
	router.Get("/generiques/group/{groupId}", findGeneriquesByGroupId)

	fmt.Printf("Starting server at PORT: %v\n", portString)
	err := server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
