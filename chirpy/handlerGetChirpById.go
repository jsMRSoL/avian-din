package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) getChirpByID(w http.ResponseWriter, r *http.Request) {
	// get path value
	path := r.PathValue("ID")
	id, err := strconv.Atoi(path)
	if err != nil {
		log.Printf("Error: ID %s could not be converted to integer", path)
		return
	}
	chirp, err := cfg.chirpsDB.GetChirp(id)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			fmt.Sprintf("Chirp ID:%d was not found.", id),
		)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
