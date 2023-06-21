// main.go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)


func logout(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Logging out by sending expired cookie...")
	// Send expired cookie to clear client's
	// existing valid cookie
	cookie := &http.Cookie{
        Name:     "jwt",
		MaxAge:   -1,
        HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
        Secure:   true,
		Domain: request.RequestContext.DomainName,
    }
	headers := make(map[string]string)
    headers["Set-Cookie"] = cookie.String()

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK, 
		Body: "Successfully logged out",
		Headers: headers,
	}
	return response, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(logout)
}
