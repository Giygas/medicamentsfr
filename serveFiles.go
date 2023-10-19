package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
	"github.com/go-chi/chi"
)

func serveAllMedicaments(w http.ResponseWriter, r *http.Request) {
	meds := &medicaments

	parsedJson, err := json.Marshal(meds)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
	}
	meds = nil
	// Write the headers for the json response
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Cache-Control", "public, max-age=43200") // caches for half a day
	w.WriteHeader(200)
	w.Write(parsedJson)
}

func findMedicament(w http.ResponseWriter, r *http.Request) {
	matchingMedicaments := make([]entities.Medicament, 0)

	userPattern := chi.URLParam(r, "element")
	pattern, compileErr := regexp.Compile(`(?i).*` + regexp.QuoteMeta(userPattern) + `.*`)

	if compileErr != nil {
		log.Panic("An error has ocurred with the search parameter", compileErr)
	} else {
		for _, med := range medicaments {
			// Search the value in medicament denomination, composition
			// denomination and generique libelle

			// For now I'll just use the medicament denomination
			medOk := pattern.MatchString(med.Denomination)
			if medOk {
				matchingMedicaments = append(matchingMedicaments, med)
			}
		}
	}
	if len(matchingMedicaments) > 0 {
		respondWithJSON(w, 200, matchingMedicaments)
	} else {
		respondWithError(w, 404, "No medicaments found that matches the word")
	}
}

func findMedicamentById(w http.ResponseWriter, r *http.Request) {
	cis, err := strconv.Atoi(chi.URLParam(r, "cis"))
	if err != nil {
		log.Fatal("An error ocurred when getting the medicament cis", cis)
	}

	medicament, ok := medicamentsMap[cis]
	if ok {
		respondWithJSON(w, 200, medicament)
	} else {
		respondWithError(w, 404, "No medicaments found with this cis")
	}
}

func findGeneriques(w http.ResponseWriter, r *http.Request) {

}

func findGeneriquesByGroupId(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(chi.URLParam(r, "groupId"))
	if err != nil {
		log.Fatal("An error ocurred when getting the medicament cis", groupId)
	}

	generique, ok := generiquesMap[groupId]
	if ok {
		respondWithJSON(w, 200, generique)
	} else {
		respondWithError(w, 404, "There are no medicaments in this group")
	}
}
