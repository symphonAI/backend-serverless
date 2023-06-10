package chatgptapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (c *Client) PromptChatGPT(userFields UserFields) (string, error) {
	endpoint := BaseURL + "/chat/completions"

	fmt.Println("User fields:", userFields)
	fmt.Println("Endpoint:", endpoint)

	body := []byte(`{
		"model":` + os.Getenv("OPENAI_MODEL") + `,
		"prompt": "` + userFields.Prompt + `",
		"max_tokens": ` + os.Getenv("OPENAI_MAX_TOKENS") + `,
		"temperature":` + userFields.Temperature + `,
		`)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	fmt.Printf("Req:\n %v", req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code %d", resp.StatusCode)
		return "", err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	

	response := ""
	err = json.Unmarshal(data, &response)
	if err != nil {
		return "", err
	}

	fmt.Printf("Resp:\n %v", response)


	return response, nil
}
