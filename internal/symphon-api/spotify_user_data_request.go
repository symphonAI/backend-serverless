package symphonapi

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func (c *Client) GetUserDataFromSpotify(spotifyAccessToken string) (SpotifyUserData, error) {
	endpoint := SPOTIFY_BASE_URL + "/me"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return SpotifyUserData{}, err
	}

	req.Header.Add("Authorization", "Bearer "+spotifyAccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SpotifyUserData{}, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return SpotifyUserData{}, err
	}


	userData := SpotifyUserData{}
	err = json.Unmarshal(responseBody, &userData)
	if err != nil {
		return SpotifyUserData{}, err
	}

	return userData, nil
}
