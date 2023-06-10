// main.go
package main

import (
	"fmt"
	"net/http"

	chatgptapi "backend-serverless/internal/chatgpt-api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type config struct {
	chatgptClient chatgptapi.Client
}

func handlePrompt(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Prompt called")
	cfg := config{
		chatgptClient: chatgptapi.NewClient(),
	}
	fmt.Println(request)

	// call chatGPT api
	userFields := chatgptapi.UserFields{
		Prompt:      "Hello. This is Squid.",
		Temperature: "0.9",
	}
	

	chatgptResponse, err := cfg.chatgptClient.PromptChatGPT(userFields)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	fmt.Println(chatgptResponse)

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}

	return response, nil
}

func main() {

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handlePrompt)
}
