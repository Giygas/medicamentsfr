package main

import (
	// "log"
	"encoding/json"
	"log"
	"net/http"

	// "os"
	"time"
)

func serveAllMedicaments(w http.ResponseWriter, r *http.Request) {
	meds := &medicaments

	// medicaments, err := os.ReadFile("./src/Medicaments.json")
	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(500)
	// }

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

func findGeneriques(w http.ResponseWriter, r *http.Request) {

}
