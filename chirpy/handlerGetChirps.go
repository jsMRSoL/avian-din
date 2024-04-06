package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, _ *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Println("Could not retrieve chirps from database")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
