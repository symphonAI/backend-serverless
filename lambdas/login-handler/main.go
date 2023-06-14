// main.go
package main

// hi sunny :^)

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handlePrompt(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)  {
	// Retrieve the request body from the event
	requestBody := LoginRequestBody{}
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		errorString := fmt.Sprintf("unable to unmarshal request body: %s", err.Error())
		response := events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       errorString,
		}
		return response, errors.New(errorString)
	}

	code := requestBody.Code
	redirectURI := requestBody.RedirectURI

	refresh_token, access_token, err := exchangeCodeForAuthTokens(code, redirectURI)
	if err != nil {
		errorString := fmt.Sprintf("unable to exchange auth code with refresh token: %s", err.Error())
		response := events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:      errorString,
		}
		return response, errors.New(errorString) // TODO should return error instead?
	}
	fmt.Println("Successfully generated refresh token and access token:", access_token)

	// Get User Email
	id, email, err := getUserIdentifiers(access_token)
	fmt.Println("User email:", email)
	if err != nil {
		errorString := fmt.Sprintf("unable to get user identifiers: %s", err.Error())
		response := events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       errorString,
		}
		return response, errors.New(errorString) 

	}
	
	// Save User ID, Refresh token against this user in DB
	saveUserAndRefreshTokenToDb(id, email, refresh_token)
	if err != nil {
		errorString := fmt.Sprintf("unable to save user and refresh token to DB: %s", err.Error())
		response := events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       errorString,
		}
		return response, errors.New(errorString) 
	}

	jwToken, err := GenerateJWT("ap-southeast-2", email)
	if err != nil {
		errorString := fmt.Sprintf("unable to generate JWT: %s", err.Error())
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError, 
			Body: errorString,
		}
		return response, errors.New(errorString)
	}

	cookie := &http.Cookie{
        Name:     "jwt",
        Value:    jwToken,
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
        Secure:   true,
		Domain: request.RequestContext.DomainName,
    }
	headers := make(map[string]string)
    headers["Set-Cookie"] = cookie.String()

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK, 
		Body: "Successfully authenticated",
		Headers: headers,
	}

	return response, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handlePrompt)
}
