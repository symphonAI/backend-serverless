// main.go
package main

// hi sunny :^)

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handlePrompt(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)  {
	// Retrieve the request body from the event
	requestBody := SignupRequestBody{}
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

	fmt.Println("Received auth code:", code)

	refresh_token, access_token, err := exchangeCodeForAuthTokens(code)
	if err != nil {
		errorString := fmt.Sprintf("unable to exchange auth code with refresh token: %s", err.Error())
		response := events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:      errorString,
		}
		return response, errors.New(errorString) // TODO should return error instead?
	}
	fmt.Println("Successfully generated refresh token and access token.")

	// TEMP LOG
	fmt.Println("[TEMP] Access token:", access_token)

	// Get User Email
	id, email, err := getUserIdentifiers(access_token)
	if err != nil {
		errorString := fmt.Sprintf("unable to get user identifiers: %s", err.Error())
		response := events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       errorString,
		}
		return response, errors.New(errorString) 

	}
	// Save user in Cognito User Pool, retain User ID
	err = saveUserToCognito(id, email)
	if err != nil {
		errorString := fmt.Sprintf("unable to save user in Cognito user pool: %s", err.Error())
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

	fmt.Println("Successfully generated JWT. Returning JWT in response...")
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK, 
		Body: jwToken,
	}

	return response, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handlePrompt)
}
