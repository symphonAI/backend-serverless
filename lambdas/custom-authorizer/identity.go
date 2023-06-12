package main

import "fmt"

type IdentityProvider interface {
	getAccessCredentials(u *User) (a AccessCredentials, e error)
}


func getIdentityProvider(u *User) (i IdentityProvider, e error) {
	switch (u.IDProvider){
	case Spotify:
		return SpotifyIdentityProvider{}, nil
	default:
		return nil, fmt.Errorf("Could not find Ideneity Provider with name: %v", u.IDProvider)
	}
}