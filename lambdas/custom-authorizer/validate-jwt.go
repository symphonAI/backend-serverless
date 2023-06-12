package main

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func validateJWT(tokenString string) (jwt.MapClaims, error) {
	issuerID := os.Getenv("ISSUER_ID")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		return []byte(issuerID), nil
	})

	if err != nil {
		fmt.Println("Error while parsing token:", err)
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		fmt.Println("Token validity is false")
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Failed to extract claims.")
		return nil, fmt.Errorf("Failed to extract claims")
	}
	return claims, nil
}
