package main

type LoginRequestBody struct {
	Code string `json:"code"`
	RedirectURI string `json:"redirectUri"`
}

type TokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type SpotifyResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

