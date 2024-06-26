package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {

	token, _, err := getTokenAndStringFromHeader(r, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token invalid/expired")
		return
	}

	issuer, authorId, err := parseToken(token)
	if issuer != "chirpy-access" {
		respondWithError(w, http.StatusUnauthorized, "Bad token")
	}

	type parameters struct {
		// these tags indicate how the keys in the JSON should be mapped to the struct fields
		// the struct fields must be exported (start with a capital letter) if you want them parsed
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	msg := params.Body
	if len(msg) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	msg = cleanChirp(msg)

	chirp, err := cfg.chirpsDB.StoreChirp(msg, authorId)

	respondWithJSON(w, http.StatusCreated, chirp)

}
