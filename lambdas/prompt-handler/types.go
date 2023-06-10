package main

type PromptRequestBody struct {
	Prompt string `json:"prompt"`
	Temperature string `json:"temperature"`
}