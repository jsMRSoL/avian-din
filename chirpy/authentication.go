package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func createSignedString(
	id int,
	// email string,
	expiresInSecs int,
) (string, error) {

	log.Println("expiresInSecs: ", expiresInSecs)
	expirationDuration := 24 * time.Hour
	if expiresInSecs != 0 {
		log.Println("Got to here")
		expirationDuration = time.Duration(expiresInSecs) * time.Second
	}

	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	// type MyCustomClaims struct {
	// 	Email string `json:"email"`
	// 	jwt.RegisteredClaims
	// }
	//
	expiresAt := jwt.NewNumericDate(time.Now().Add(expirationDuration))
	log.Printf("Token expires at %v", expiresAt.Time)
	claims := jwt.RegisteredClaims{
		ExpiresAt: expiresAt,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "Chirpy",
		Subject:   strconv.Itoa(id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSecret))
	return ss, err
}
