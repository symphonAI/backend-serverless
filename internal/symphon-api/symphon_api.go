package symphonapi

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const OPENAI_BASE_URL = "https://api.openai.com/v1"
const SPOTIFY_BASE_URL = "https://api.spotify.com/v1"

// Client is a client for the ChatGPT API.
type Client struct {
	httpClient http.Client
	apiKey     string
	openAIModel OpenAIModel 
}

func NewClient() Client {
	modelIdentifier := os.Getenv("OPENAI_MODEL")
	model, err := getModel(modelIdentifier)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Instantiating API client with model: %v\n", modelIdentifier)
	return Client{
		httpClient: http.Client{
			Timeout: time.Minute,
		},
		apiKey: os.Getenv("OPENAI_API_KEY"),
		openAIModel: model,
	}
}

func getModel(model string) (OpenAIModel, error) {
	switch (model){
	case "davinci":
		return &DaVinciModel{}, nil
	case "gpt-3.5-turbo":
		return &ChatCompletionModel{
			chatGptApiLLMModel: "gpt-3.5-turbo-0613",
		}, nil
	case "gpt-4":
		return &ChatCompletionModel{
			chatGptApiLLMModel: "gpt-4-0613",
		}, nil
	default:
		return nil, fmt.Errorf("could not find supported model with name: %v", model)
	}
}
