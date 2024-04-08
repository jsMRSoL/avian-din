package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	err := cfg.userDB.AddRevokedToken(tokenString)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}
