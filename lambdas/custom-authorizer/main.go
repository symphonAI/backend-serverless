// main.go
package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func authorize(ctx context.Context, event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	fmt.Println("Custom authorizer called:", event)

	cookieStr := event.Headers["cookie"]

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

	user, err := getUserFromDB(jwtClaims["user"].(string))
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	ip, err := getIdentityProvider(user)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}
	accessCredentials, err := ip.getAccessCredentials(user)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}
	// Generate the policy document for the user
	// TODO I think there will be some BS here to deal with
	policyDocument := generatePolicy(event.MethodArn)

	// Generate the authorizer response
	response := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID:  user.SortKey,
		PolicyDocument: policyDocument,
		Context: accessCredentials,
	}

	return response, nil
}


func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(authorize)
}
