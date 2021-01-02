package lastfm

import (
	"os"

	"github.com/shkh/lastfm-go/lastfm"
)

// Lastfm is the external Lastfm service implementation.
type Lastfm struct {
	api *lastfm.Api
}

// NewLastfm creates a Lastfm instance.
func NewLastfm() *Lastfm {
	clientID := os.Getenv("LASTFM_CLIENT_ID")
	clientSecret := os.Getenv("LASTFM_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		panic("Please set LASTFM_CLIENT_ID and LASTFM_CLIENT_SECRET environment variables.")
	}

	return &Lastfm{
		api: lastfm.New(clientID, clientSecret),
	}
}

// Name returns the human readable service name.
func (l *Lastfm) Name() string {
	return "Last.fm"
}

// CreateAuthURL returns an authorization URL to authorize the integration.
func (l *Lastfm) CreateAuthURL() string {
	// Not an actual web server (yet).
	redirectURL := "https://admirer.test"

	return l.api.GetAuthRequestUrl(redirectURL)
}

// Authenticate takes an authorization code and authenticates the user.
func (l *Lastfm) Authenticate(oauthCode string) {
	err := l.api.LoginWithToken(oauthCode)
	if err != nil {
		panic("Failed to parse Last.fm token.")
	}
}

// GetUsername requests and returns the username of the logged in user.
func (l *Lastfm) GetUsername() string {
	user, err := l.api.User.GetInfo(lastfm.P{})
	if err != nil {
		panic("Failed to read Last.fm profile data.")
	}

	return user.Name
}
