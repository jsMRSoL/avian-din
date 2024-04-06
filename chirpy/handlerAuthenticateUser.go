package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) authenticateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		// these tags indicate how the keys in the JSON should be mapped to the struct fields
		// the struct fields must be exported (start with a capital letter) if you want them parsed
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	user, err := cfg.userDB.AuthenticateUser(params.Email, params.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
	return
}
