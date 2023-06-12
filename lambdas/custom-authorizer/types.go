package main

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
}

type AccessCredentials map[string]interface{}