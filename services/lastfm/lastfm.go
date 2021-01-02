package lastfm

import (
	"fmt"
	"os"

	lastfm_api "github.com/shkh/lastfm-go/lastfm"
)

// Login runs the CLI procedure for logging in on Last.fm.
func Login(oauthCode string) {
	apiClient := createAPIClient()

	if len(oauthCode) == 0 {
		fmt.Println("Last.fm authentication URL: " + createAuthURL(apiClient))
		return
	}
}

func createAPIClient() *lastfm_api.Api {
	clientID := os.Getenv("LASTFM_CLIENT_ID")
	clientSecret := os.Getenv("LASTFM_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		panic("Please set LASTFM_CLIENT_ID and LASTFM_CLIENT_SECRET environment variables.")
	}

	return lastfm_api.New(clientID, clientSecret)
}

func createAuthURL(apiClient *lastfm_api.Api) string {
	// Not an actual web server (yet).
	redirectURL := "https://admirer.test"

	return apiClient.GetAuthRequestUrl(redirectURL)
}
