// main.go
package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func authorize(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	fmt.Println("Custom authorizer called")

	cookieStr := event.AuthorizationToken

	tokenString := extractJwtFromCookie(cookieStr)

	if tokenString == nil {
		fmt.Println("Missing JWT from cookie")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
	}
	jwtClaims, err := validateJWT(*tokenString)
	if err != nil {
		fmt.Println("JWT is not valid")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
	}

	// Generate the policy document for the user
	// TODO I think there will be some BS here to deal with
	policyDocument := generatePolicy(event.MethodArn)

	// Generate the authorizer response
	response := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID:  jwtClaims["user"].(string),
		PolicyDocument: policyDocument,
	}

	return response, nil
}


func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(authorize)
}