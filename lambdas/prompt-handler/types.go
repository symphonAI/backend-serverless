package main

type PromptRequestBody struct {
	Prompt      string   `json:"prompt"`
	Temperature string   `json:"temperature"`
	Options     []string `json:"options`
}

type Track struct {
	Artist string `json:"artist"`
	Track  string `json:"track"`
}
