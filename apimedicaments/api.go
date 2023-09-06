// TODO: Make a map of generiques with the format
// [group]: []{cis, libelle} for searching by group
// and for the libelle searching will just iterate
// over this map of libelles -> include regexp

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	// "github.com/giygas/medicamentsfr/medicamentsparser"
	"github.com/giygas/medicamentsfr/medicamentsparser"
	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var medicaments []entities.Medicament
var generiques []entities.GeneriqueList

func checkMedicaments(medicaments *[]entities.Medicament) {
	if len(*medicaments) == 0 {
		fmt.Println("medicaments slice is empty")
		return
	}

	for i, medicament := range *medicaments {
		if medicament.Cis == 0 {
			fmt.Printf("medicament at index %d has Cis set to 0\n", i)
		}
	}
}
func init() {
	// Load the env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	// Create the initial medicaments parsing
	medicaments = medicamentsparser.ParseAllMedicaments()
	checkMedicaments(&medicaments)
	generiques = medicamentsparser.GeneriquesParser(&medicaments)
	// go scheduleMedicaments()
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
	// Search medicaments by elementPharmaceutique or cis
	router.Get("/medicament/{cis}", findMedicament)
	// Searh medicaments by generiques libelle or generiques group
	router.Get("/generiques/{libelle}", findGeneriques)

	fmt.Printf("Starting server at PORT: %v\n", portString)
	err := server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
