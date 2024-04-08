package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {

	token, _, err := getTokenAndStringFromHeader(r, cfg.secret)
	issuer, authorId, err := parseToken(token)
	if issuer != "chirpy-access" {
		respondWithError(w, http.StatusUnauthorized, "Bad token")
		return
	}

	path := r.PathValue("ID")
	chirpId, err := strconv.Atoi(path)
	if err != nil {
		log.Printf("Error: ID %s could not be converted to integer", path)
		return
	}

	// is the deleter the author of the chirp?
	chirp, err := cfg.chirpsDB.GetChirp(chirpId)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			fmt.Sprintf("Chirp ID:%d was not found.", chirpId),
		)
		return
	}

	if chirp.AuthorId != authorId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = cfg.chirpsDB.DeleteChirp(chirpId)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			fmt.Sprintf("Chirp ID:%d was not deleted.", chirpId),
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}
