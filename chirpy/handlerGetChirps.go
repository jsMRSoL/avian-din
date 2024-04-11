package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {

	s := r.URL.Query().Get("author_id")
	order := r.URL.Query().Get("sort")

	desc := false
	if order == "desc" {
		desc = true
	}

	if s == "" {
		cfg.allChirps(w, desc)
		return
	}

	authorID, err := strconv.Atoi(s)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid author_id: %s", s),
		)
		return
	}

	cfg.chirpsByAuthorID(w, authorID, desc)
	return
}

func (cfg *apiConfig) chirpsByAuthorID(
	w http.ResponseWriter,
	authorID int,
	desc bool,
) {
	chirps, err := cfg.chirpsDB.ChirpsByAuthorID(authorID, desc)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			fmt.Sprintf(
				"Chirps with author_id %d could not be retrieved",
				authorID,
			),
		)
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) allChirps(w http.ResponseWriter, desc bool) {
	chirps, err := cfg.chirpsDB.GetChirps(desc)
	if err != nil {
		log.Println("Could not retrieve chirps from database")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
