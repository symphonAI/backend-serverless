package chatgptapi

type UserFields struct {
	Prompt      string
	Temperature string
	Options     *string
}

type ChatGPTResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}