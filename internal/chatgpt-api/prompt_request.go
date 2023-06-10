package chatgptapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (c *Client) PromptChatGPT(userFields UserFields) (string, error) {
	endpoint := BaseURL + "/completions"

	fmt.Println("User fields:", userFields)
	fmt.Println("Endpoint:", endpoint)

	floatTemperature, err := strconv.ParseFloat(userFields.Temperature, 64)
	if err != nil {
		return "", err
	}

	payload := map[string]interface{}{
		"model":       "text-davinci-003",
		"max_tokens":  2000,
		"temperature": floatTemperature,
		"prompt":      userFields.Prompt,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	fmt.Println("Req:", req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("Resp:", resp)

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Resp:", string(responseBody))

	response := ChatGPTResponse{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", err
	}


	fmt.Println("Resp:", response.Choices[0].Text)

	return response.Choices[0].Text, nil
}
