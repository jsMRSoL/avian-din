package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	// "time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {

	// Authorization string `json:"Authorization"`
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	type MyCustomClaims struct {
		Email string `json:"email"`
		jwt.RegisteredClaims
	}

	/// Get env variable
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token invalid/expired")
		return
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Token claims could not be accessed",
		)
		return
	}

	idString, err := claims.GetSubject()
	if !ok {
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
