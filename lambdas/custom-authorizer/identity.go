package main

import "fmt"

type IdentityProvider interface {
	getAccessCredentials(u *User) (a AccessCredentials, e error)
}


func getIdentityProvider(u *User) (i IdentityProvider, e error) {
	fmt.Println("Fetching identity provider based on user...")
	switch (u.IDProvider){
	case Spotify:
		return SpotifyIdentityProvider{}, nil
	default:
		return nil, fmt.Errorf("Could not find Identity Provider with name: %v", u.IDProvider)
	}
}