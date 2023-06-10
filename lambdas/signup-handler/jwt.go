package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(region, email string) (string, error) {
    fmt.Println("Generating JWT...")

	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	userPoolClientID := os.Getenv("COGNITO_USER_POOL_CLIENT_ID")

    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    claims["aud"] = userPoolClientID
    claims["cognito:username"] = email
    claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Token expiration time

    tokenString, err := token.SignedString([]byte(userPoolID))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
