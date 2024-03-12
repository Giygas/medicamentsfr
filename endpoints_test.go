package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/giygas/medicamentsfr/medicamentsparser"
	"github.com/go-chi/chi/v5"
)

func initRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(rateLimitHandler)

	router.Get("/database/{pageNumber}", servePagedMedicaments)
	router.Get("/database", serveAllMedicaments)
	router.Get("/medicament/{element}", findMedicament)
	router.Get("/medicament/id/{cis}", findMedicamentById)
	router.Get("/generiques/{libelle}", findGeneriques)
	router.Get("/generiques/group/{groupId}", findGeneriquesByGroupId)

	return router
}

func TestMain(m *testing.M) {
	medicaments = medicamentsparser.ParseAllMedicaments()

	// Create a map of all medicaments to reduce algorithm complexity
	for i := range medicaments {
		medicamentsMap[(medicaments)[i].Cis] = (medicaments)[i]
	}

	generiques, generiquesMap = medicamentsparser.GeneriquesParser(&medicaments, &medicamentsMap)
	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestEndpoints(t *testing.T) {

	testCases := []struct {
		name     string
		endpoint string
		expected int
	}{

		{"Test database", "/database", http.StatusOK},
		{"Test generiques/paracetamol", "/generiques/paracetamol", http.StatusOK},
		{"Test generiques/group/1643", "/generiques/group/1643", http.StatusOK},
		{"Test medicament/doli", "/medicament/doli", http.StatusOK},
		{"Test database with a", "/database/a", http.StatusNotFound},
		{"Test database with 1", "/database/1", http.StatusOK},
		{"Test database with 0", "/database/0", http.StatusNotFound},
		{"Test database with -1", "/database/-1", http.StatusNotFound},
		{"Test generiques", "/generiques", http.StatusNotFound},
		{"Test generiques/aaaaaaaaaaa", "/generiques/aaaaaaaaaaa", http.StatusNotFound},
		{"Test medicament", "/medicament", http.StatusNotFound},
		{"Test medicament/1000000000000000", "/medicament/100000000000000", http.StatusNotFound},
		{"Test medicament/id/1", "/medicament/id/1", http.StatusNotFound},
		{"Test generiques/group/a", "/generiques/group/a", http.StatusBadRequest},
	}

	router := initRouter()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.endpoint, nil)

			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			if status := rr.Code; status != tt.expected {
				t.Errorf("%v returned wrong status code: got %v want %v", tt.endpoint, status, tt.expected)
			}
		})
	}
}
