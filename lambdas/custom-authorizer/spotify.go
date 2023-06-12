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

type SpotifyIdentityProvider struct {}

func (su SpotifyIdentityProvider) getAccessCredentials(u *User) (a AccessCredentials, e error){
	accessToken, err := exchangeRefreshTokenForAuthToken(u.RefreshToken)
	if err != nil {
		return nil, err
	}
	a = make(AccessCredentials)
	a["accessToken"] = accessToken
	return a, nil
}

func exchangeRefreshTokenForAuthToken(refreshToken string) (string, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	tokenURL := "https://accounts.spotify.com/api/token"
	authData := url.Values{
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	}

	authBody := bytes.NewBufferString(authData.Encode())

	req, err := http.NewRequest("POST", tokenURL, authBody)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error calling Spotify for token:", err)
		return "", err
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
		return "", err
	}

	return tokenResponse.AccessToken, nil
}