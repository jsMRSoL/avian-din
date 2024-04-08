package main

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createSignedString(
	id int,
	issuer string,
	expirationDuration time.Duration,
	jwtSecret string,
) (string, error) {

	expiresAt := jwt.NewNumericDate(time.Now().Add(expirationDuration))
	claims := jwt.RegisteredClaims{
		ExpiresAt: expiresAt,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    issuer,
		Subject:   strconv.Itoa(id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSecret))
	return ss, err
}
