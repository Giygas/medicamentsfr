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
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("Cache-Control", "public, max-age=3600")
		w.Header().Add("Expires", time.Now().Add(time.Hour).Format(http.TimeFormat))
		w.Header().Add("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		response, err := json.Marshal(medicament)
		if err != nil {
			w.WriteHeader(500)
			log.Fatal("An error has ocurred when marshalling the medicament", err)
		} else {
			w.WriteHeader(200)
			w.Write(response)
		}
	}
}

func findGeneriques(w http.ResponseWriter, r *http.Request) {

}

func findGeneriquesById(w http.ResponseWriter, r *http.Request) {

}
