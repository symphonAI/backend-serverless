package symphonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) GetRecommendedTracks(userFields UserFields) ([]Track, error) {
	endpoint := c.openAIModel.GetUrl()

	payload, err := c.openAIModel.GeneratePayload(userFields)
	if err != nil {
		fmt.Println("An error occurred while generating the payload for the ChatGPT API:", err)
		return []Track{}, err
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("An error occurred while marshalling the payload for the ChatGPT API:", err)
		return []Track{}, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("An error occurred while preparing a request to the ChatGPT API:", err)
		return []Track{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("An error occurred while making a request to the ChatGPT API:", err)
		return []Track{}, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Track{}, err
	}

	tracks, err := c.openAIModel.ParseRecommendedTracksFromResponse(responseBody)
	if err != nil {
		fmt.Println("An error occurred while parsing the recommended tracks from the ChatGPT API response:", err)
		return []Track{}, err
	}
	
	return tracks, nil
}