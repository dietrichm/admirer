package spotify

import (
	"fmt"
	"os"

	"github.com/zmb3/spotify"
)

// Login runs the CLI procedure for logging in on Spotify.
func Login(oauthCode string) {
	spotify := NewSpotify()

	if len(oauthCode) == 0 {
		fmt.Println("Spotify authentication URL: " + spotify.CreateAuthURL())
		return
	}

	client := spotify.Authenticate(oauthCode)

	user, err := client.CurrentUser()
	if err != nil {
		panic("Failed to read Spotify profile data.")
	}

	fmt.Println("Logged in on Spotify as " + user.DisplayName)
}

// Spotify is the external Spotify service implementation.
type Spotify struct {
	authenticator *spotify.Authenticator
}

// NewSpotify creates a Spotify instance.
func NewSpotify() *Spotify {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		panic("Please set SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET environment variables.")
	}

	// Not an actual web server (yet).
	redirectURL := "https://admirer.test"
	authenticator := spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate)
	authenticator.SetAuthInfo(clientID, clientSecret)

	return &Spotify{
		authenticator: &authenticator,
	}
}

// CreateAuthURL returns an authorization URL to authorize the integration.
func (s *Spotify) CreateAuthURL() string {
	return s.authenticator.AuthURL("")
}

// Authenticate takes an authorization code and authenticates the user.
func (s *Spotify) Authenticate(code string) spotify.Client {
	token, err := s.authenticator.Exchange(code)
	if err != nil {
		panic("Failed to parse Spotify token.")
	}

	return s.authenticator.NewClient(token)
}
