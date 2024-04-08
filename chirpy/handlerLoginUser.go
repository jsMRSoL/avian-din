package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		// these tags indicate how the keys in the JSON should be mapped to the struct fields
		// the struct fields must be exported (start with a capital letter) if you want them parsed
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding json: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	user, err := cfg.userDB.AuthenticateUser(params.Email, params.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessTokenString, err := createSignedString(
		user.Id,
		"chirpy-access",
		1*time.Hour,
		cfg.secret,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	refreshTokenString, err := createSignedString(
		user.Id,
		"chirpy-refresh",
		60*24*time.Hour,
		cfg.secret,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		user.ToSignedUser(accessTokenString, refreshTokenString),
	)
}
