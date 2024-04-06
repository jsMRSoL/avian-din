package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) addUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		// these tags indicate how the keys in the JSON should be mapped to the struct fields
		// the struct fields must be exported (start with a capital letter) if you want them parsed
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	email := params.Email
	user, err := cfg.userDB.AddUser(email)

	respondWithJSON(w, http.StatusCreated, user)

}
