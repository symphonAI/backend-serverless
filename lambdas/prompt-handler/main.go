// main.go
package main

import (
	"context"
	"net/http"

	chatgptapi "backend-serverless/internal/chatgpt-api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type config struct {
	chatgptClient chatgptapi.Client
}

func handlePrompt(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg := config{
		chatgptClient: chatgptapi.NewClient(),
	}
	// call chatGPT api
	userFields := chatgptapi.UserFields {
		Prompt: "Hello. This is Squid.",
		Temperature: "0.9",
	}

	chatgptResponse, err := cfg.chatgptClient.PromptChatGPT(userFields)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}



	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK, 
		Body:       chatgptResponse,
	}

	return response, nil
}

func main() {

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handlePrompt)
}
