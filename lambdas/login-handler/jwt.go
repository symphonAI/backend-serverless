package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(region, email string) (string, error) {
    fmt.Println("Generating JWT...")

	issuer_id := os.Getenv("ISSUER_ID")
	aud_id := os.Getenv("AUDIENCE_ID")

    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    claims["aud"] = aud_id
    claims["iss"] = issuer_id
    claims["user"] = email
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time

    tokenString, err := token.SignedString([]byte(issuer_id))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
