package symphonapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (c *Client) GetRecommendedTracks(userFields UserFields) ([]Track, error) {
	endpoint := c.openAIModel.GetUrl()

	payload, err := c.openAIModel.GeneratePayload(userFields)
	if err != nil {
		return []Track{}, err
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return []Track{}, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return []Track{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return []Track{}, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Track{}, err
	}

	tracks, err := c.openAIModel.ParseRecommendedTracksFromResponse(responseBody)
	if err != nil {
		return []Track{}, err
	}
	
	return tracks, nil
}