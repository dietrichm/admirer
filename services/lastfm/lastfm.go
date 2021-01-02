package lastfm

import (
	"fmt"
	"os"

	lastfm_api "github.com/shkh/lastfm-go/lastfm"
)

// Login runs the CLI procedure for logging in on Last.fm.
func Login(oauthCode string) {
	lastfm := NewLastfm()
	apiClient := lastfm.api

	if len(oauthCode) == 0 {
		fmt.Println("Last.fm authentication URL: " + lastfm.CreateAuthURL())
		return
	}

	callback(apiClient, oauthCode)

	user, err := apiClient.User.GetInfo(lastfm_api.P{})
	if err != nil {
		panic("Failed to read Last.fm profile data.")
	}

	fmt.Println("Logged in on Last.fm as " + user.Name)
}

// Lastfm is the external Lastfm service implementation.
type Lastfm struct {
	api *lastfm_api.Api
}

// NewLastfm creates a Lastfm instance.
func NewLastfm() *Lastfm {
	clientID := os.Getenv("LASTFM_CLIENT_ID")
	clientSecret := os.Getenv("LASTFM_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		panic("Please set LASTFM_CLIENT_ID and LASTFM_CLIENT_SECRET environment variables.")
	}

	return &Lastfm{
		api: lastfm_api.New(clientID, clientSecret),
	}
}

// CreateAuthURL returns an authorization URL to authorize the integration.
func (l *Lastfm) CreateAuthURL() string {
	// Not an actual web server (yet).
	redirectURL := "https://admirer.test"

	return l.api.GetAuthRequestUrl(redirectURL)
}

func callback(apiClient *lastfm_api.Api, oauthCode string) {
	err := apiClient.LoginWithToken(oauthCode)
	if err != nil {
		panic("Failed to parse Last.fm token.")
	}
}
