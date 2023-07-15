package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func exchangeCodeForAuthTokens(code string, redirectURI string) (string, string, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	fmt.Println("Redirect URI:", redirectURI)

	authURL := "https://accounts.spotify.com/api/token"
	authData := url.Values{
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	}

	authBody := bytes.NewBufferString(authData.Encode())

	req, err := http.NewRequest("POST", authURL, authBody)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error calling Spotify for tokens:", err)
		return "", "", err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// TODO error handling 
		fmt.Println("Error parsing:", err)
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(responseData, &tokenResponse)
	if err != nil {
		return "", "", err
	}
	if tokenResponse.AccessToken == "" {
		return "", "", fmt.Errorf("invalid response from Spotify API: %v", string(responseData[:]))
	}


	return tokenResponse.RefreshToken, tokenResponse.AccessToken, nil
}

func getUserIdentifiers(accessToken string) (string, string, error){
	fmt.Println("Getting user data from spotify...", accessToken)

	url := "https://api.spotify.com/v1/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error requesting Spotify:", err)
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Spotify response:", err)
		return "", "", err
	}

	var data SpotifyResponse
	err = json.Unmarshal(body, &data)

	fmt.Println("Response as string:", string(body[:]))

	if err != nil {
		fmt.Println("Error parsing Spotify response:", err)
		return "", "", err
	}

	fmt.Println("Successfully retrieved Spotify user data:", data.Email)
	if data.Email == "" {
		return "", "", fmt.Errorf("could not obtain email from Spotify")
	}
	return data.ID, data.Email, nil

}