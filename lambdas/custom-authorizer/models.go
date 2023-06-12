package main

type User struct {
	PartitionKey string
	SortKey string
	RefreshToken string
	Username string
	IDProvider IDProvider
	UpdatedOn int64
	CreatedOn int64
}


type IDProvider string

const (
	Spotify IDProvider = "SPOTIFY"
)