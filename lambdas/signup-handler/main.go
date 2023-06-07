// main.go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handlePrompt(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)  {
	code := request.QueryStringParameters["code"]
	fmt.Println("Received auth code:", code)

	refresh_token, access_token, err := exchangeCodeForAuthTokens(code)
	if err != nil {
		response := events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to exchange auth code with refresh token",
		}
		return response, nil // TODO should return error instead?
	}

	// Get User Email
	id, email, err := getUserIdentifiers(access_token)

	// Save user in Cognito User Pool, retain User ID
	err = SaveUserToCognito(id, email)

	// TODO if err != nil etc.....
	
	// Save User ID, Refresh token against this user in DB
	saveUserAndRefreshTokenToDb(id, email, refresh_token)

	// TODO if err != nil etc.....

	jwToken, err := GenerateJWT("ap-southeast-2", email)
	if err != nil {
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError, 
			Body: "Error creating JWT Token",
		}
		return response, nil
	}

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
