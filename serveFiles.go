package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
	"github.com/go-chi/chi/v5"
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Cache-Control", "public, max-age=43200") // caches for half a day
	w.WriteHeader(200)
	gz := gzip.NewWriter(w)
	defer gz.Close()
	gz.Write(parsedJson)
}

func servePagedMedicaments(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(chi.URLParam(r, "pageNumber"))

	if err != nil || page < 1 {
		respondWithError(w, 400, "Invalid page number")

		return
	}
	// Get the maximal page possible for the medicaments
	maxPagePossible := len(medicamentsMap) / 10

	if len(medicamentsMap)%10 != 0 {
		maxPagePossible++
	}

	medicamentsUpper := page * 10

	// If the upper slice of the medicaments is bigger than the maximum, give an error message and exit
	if medicamentsUpper > maxPagePossible*10 {
		respondWithError(w, 404, "Maximum page possible is: "+strconv.Itoa(maxPagePossible))
		return
	}

	medicamentsLower := medicamentsUpper - 10

	if page < maxPagePossible {
		fmt.Println(medicamentsUpper)
		respondWithJSON(w, 200, medicaments[medicamentsLower:medicamentsUpper])
	} else {
		respondWithJSON(w, 200, medicaments[medicamentsLower:])
	}
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
		respondWithError(w, 400, "Bad Request: Not a Number")
		return
	}

	medicament, ok := medicamentsMap[cis]
	if ok {
		respondWithJSON(w, 200, medicament)
	} else {
		respondWithError(w, 404, "No medicaments found with this cis")
	}
}

func findGeneriques(w http.ResponseWriter, r *http.Request) {
	matchingGeneriques := make([]entities.GeneriqueList, 0)

	userPattern := chi.URLParam(r, "libelle")
	pattern, compileErr := regexp.Compile(`(?i).*` + regexp.QuoteMeta(userPattern) + `.*`)

	if compileErr != nil {
		log.Panic("An error has ocurred with the search parameter", compileErr)
	} else {
		for _, gen := range generiques {
			// Search for the generique name in generique libelle

			medOk := pattern.MatchString(gen.Libelle)
			if medOk {
				matchingGeneriques = append(matchingGeneriques, gen)
			}
		}
	}
	if len(matchingGeneriques) > 0 {
		respondWithJSON(w, 200, matchingGeneriques)
	} else {
		respondWithError(w, 404, "No generiques found that matches the word")
	}
}

func findGeneriquesByGroupId(w http.ResponseWriter, r *http.Request) {

	groupId, err := strconv.Atoi(chi.URLParam(r, "groupId"))
	if err != nil {
		respondWithError(w, 400, "Bad Request: Not a Number")
		return
	}

	generique, ok := generiquesMap[groupId]
	if ok {
		respondWithJSON(w, 200, generique)
	} else {
		respondWithError(w, 404, "There are no medicaments in this generique group")
	}
}
