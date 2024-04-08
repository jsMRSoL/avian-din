package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getTokenAndStringFromHeader(
	r *http.Request,
	secret string,
) (*jwt.Token, string, error) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	return token, tokenString, err
}

func parseToken(token *jwt.Token) (issuer string, authorId int, err error) {
	issuer, err = token.Claims.GetIssuer()
	if err != nil {
		log.Printf("Could not get issuer from access token: %s", err)
		return "", 0, err
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("Could not get subject from access token: %s", err)
		return "", 0, err
	}

	authorId, err = strconv.Atoi(idString)
	if err != nil {
		log.Printf("Could not convert token id string to int: %s", err)
		return "", 0, err
	}

	return
}

func (cfg *apiConfig) refreshAccessToken(w http.ResponseWriter, r *http.Request) {

	token, tokenString, err := getTokenAndStringFromHeader(r, cfg.secret)
	// Authorization string `json:"Authorization"`
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token invalid/expired")
		return
	}

	issuer, authorId, err := parseToken(token)

	if issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Bad token")
		return
	}

	isRevoked, err := cfg.userDB.IsRevoked(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bad token")
		return
	}
	if isRevoked {
		respondWithError(w, http.StatusUnauthorized, "Bad token")
		return
	}

	accessTokenString, err := createSignedString(
		authorId,
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

	type accessToken struct {
		Token string `json:"token"`
	}
	tkn := accessToken{
		Token: accessTokenString,
	}

	respondWithJSON(
		w,
		http.StatusOK,
		tkn,
	)
}
