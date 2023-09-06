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
var generiques []entities.GeneriqueList
var medicamentsMap = make(map[int]entities.Medicament)

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
	// Create a map of all medicaments to reduce algorithm complexity
	for i := range medicaments {
		medicamentsMap[(medicaments)[i].Cis] = (medicaments)[i]
	}
	generiques = medicamentsparser.GeneriquesParser(&medicaments, &medicamentsMap)
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
	router.Get("/medicament/{element}", findMedicament)
	router.Get("/medicament/id/{cis}", findMedicamentById)
	router.Get("/generiques/{libelle}", findGeneriques)
	router.Get("/generiques/group/{groupId}", findGeneriquesById)

	fmt.Printf("Starting server at PORT: %v\n", portString)
	err := server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
