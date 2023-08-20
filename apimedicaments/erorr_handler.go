package main

import "net/http"

func errorHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Something went wrong")
}
