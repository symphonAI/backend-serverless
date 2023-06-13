package symphonapi

import (
	"net/http"
	"os"
	"time"
)

const BaseURL = "https://api.openai.com/v1"

// Client is a client for the ChatGPT API.
type Client struct {
	httpClient http.Client
	apiKey     string
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Minute,
		},
		apiKey: os.Getenv("OPENAI_API_KEY"),
	}
}
