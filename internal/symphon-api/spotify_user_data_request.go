package symphonapi

import (
	"context"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func (c *Client) GetUserDataFromSpotify(ctx context.Context) (SpotifyUserData, error) {
	endpoint := SPOTIFY_BASE_URL + "/me"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return SpotifyUserData{}, err
	}

	spotifyAccessToken := ctx.Value("spotify_access_token")

	req.Header.Add("Authorization", "Bearer "+spotifyAccessToken.(string))

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
