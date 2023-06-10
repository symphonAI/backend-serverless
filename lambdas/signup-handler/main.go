// main.go
package main

import (
	"context"
	"encoding/json"
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
		response := events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "JSON key: 'code' missing from response",
		}
		return response, nil 
	}

	code := requestBody.Code

	fmt.Println("Received auth code:", code)

	refresh_token, access_token, err := exchangeCodeForAuthTokens(code)
	if err != nil {
		response := events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to exchange auth code with refresh token",
		}
		return response, nil // TODO should return error instead?
	}
	fmt.Println("Successfully generated refresh token and access token.")

	// Get User Email
	id, email, err := getUserIdentifiers(access_token)
	if err != nil {
		response := events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to get Spotify user identifiers",
		}
		return response, nil 

	}
	// Save user in Cognito User Pool, retain User ID
	err = saveUserToCognito(id, email)
	if err != nil {
		response := events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to save user in Cognito user pool",
		}
		return response, nil 
	}
	
	// Save User ID, Refresh token against this user in DB
	saveUserAndRefreshTokenToDb(id, email, refresh_token)
	if err != nil {
		response := events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to save user to DB",
		}
		return response, nil 
	}

	jwToken, err := GenerateJWT("ap-southeast-2", email)
	if err != nil {
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError, 
			Body: "Error creating JWT Token",
		}
		return response, nil
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
