// main.go
package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)


func handlePrompt(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	accessToken := ctx.Value("spotifyAccessToken")

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body: accessToken.(string),
	}
	return response, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handlePrompt)
}
