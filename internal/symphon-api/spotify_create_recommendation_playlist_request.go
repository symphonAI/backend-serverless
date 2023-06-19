package symphonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (c *Client) CreateRecommendationPlaylist(spotifyAccessToken string, userId string, trackIDs []string, prompt string, options []string) (string, error) {
	// create playlist
	playlistURI, playlistID, err := c.createPlaylist(spotifyAccessToken, userId, prompt, options)
	if err != nil {
		return "", err
	}
	fmt.Println(playlistURI)

	// add tracks to playlist
	err = c.addTracksToPlaylist(spotifyAccessToken, playlistID, trackIDs)
	if err != nil {
		return "", err
	}

	return playlistURI, nil
}

func (c *Client) createPlaylist(spotifyAccessToken string, userId string, prompt string, options []string) (string, string, error) {

	endpoint := SPOTIFY_BASE_URL + "/me/" + userId + "/playlists"

	playlistName := prompt + " - " + strings.Join(options, ", ")
	playlistDescription := "Created by Symphon.ai"
	playlistPublic := true
	playlistCollaborative := false

	payload := map[string]interface{}{
		"name":          playlistName,
		"description":   playlistDescription,
		"public":        playlistPublic,
		"collaborative": playlistCollaborative,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+spotifyAccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	response := CreatePlaylistResponse{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", "", err
	}

	fmt.Println("Playlist response:", response)

	playlistID := response.ID
	playlistURI := response.URI

	return playlistURI, playlistID, nil
}

func (c *Client) addTracksToPlaylist(spotifyAccessToken string, playlistID string, trackIDs []string) error {
	endpoint := SPOTIFY_BASE_URL + "/playlists/" + playlistID + "/tracks"

	payload := map[string]interface{}{
		"uris":     trackIDs,
		"position": 0,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+spotifyAccessToken)

	_, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
