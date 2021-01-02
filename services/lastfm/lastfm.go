package lastfm

import (
	"os"

	lastfm_api "github.com/shkh/lastfm-go/lastfm"
)

// Login runs the CLI procedure for logging in on Last.fm.
func Login(oauthCode string) {
	apiClient := createAPIClient()
}

func createAPIClient() *lastfm_api.Api {
	clientID := os.Getenv("LASTFM_CLIENT_ID")
	clientSecret := os.Getenv("LASTFM_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		panic("Please set LASTFM_CLIENT_ID and LASTFM_CLIENT_SECRET environment variables.")
	}

	return lastfm_api.New(clientID, clientSecret)
}
