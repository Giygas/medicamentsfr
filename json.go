package main

import (
	"compress/gzip"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(code)
	gz := gzip.NewWriter(w)
	defer gz.Close()
	gz.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}
