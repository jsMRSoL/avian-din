package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	// "time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {

	// Authorization string `json:"Authorization"`
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.secret), nil
		},
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token invalid/expired")
		return
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bad token")
		return
	}

	if issuer == "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Bad token")
		return
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bad token")
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Printf("String conversion error: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Bad token: no id")
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	user, err := cfg.userDB.UpdateUser(id, params.Email, params.Password)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update db")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
	return
}
