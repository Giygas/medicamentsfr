package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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
	w.Header().Add("Cache-Control", "public, max-age=3600")
	w.Header().Add("Expires", time.Now().Add(time.Hour).Format(http.TimeFormat))
	w.Header().Add("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.WriteHeader(200)
	w.Write(parsedJson)
}

func findMedicament(w http.ResponseWriter, r *http.Request) {

}

func findMedicamentById(w http.ResponseWriter, r *http.Request) {
	cis, err := strconv.Atoi(chi.URLParam(r, "cis"))
	if err != nil {
		log.Fatal("An error ocurred when getting the medicament cis", cis)
	}

	medicament, ok := medicamentsMap[cis]
	if ok {
		respondWithJSON(w, 200, medicament)
	}
	respondWithError(w, 404, "No medicaments found with this cis")
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
