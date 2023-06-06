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

func exchangeCodeForAuthTokens(code string) (string, string, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")

	fmt.Println("Client ID:", clientID)
	fmt.Println("Client Secret:", clientSecret)
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
		return "", "", err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// TODO error handling 
	}
	
	var tokenResponse TokenResponse
	err = json.Unmarshal(responseData, &tokenResponse)
	if err != nil {
		return "", "", err
	}

	return tokenResponse.AccessToken, tokenResponse.RefreshToken, nil
}