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

	spotify.Authenticate(oauthCode)

	fmt.Println("Logged in on Spotify as " + spotify.GetUsername())
}

// Spotify is the external Spotify service implementation.
type Spotify struct {
	authenticator *spotify.Authenticator
	client        *spotify.Client
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
func (s *Spotify) Authenticate(code string) {
	token, err := s.authenticator.Exchange(code)
	if err != nil {
		panic("Failed to parse Spotify token.")
	}

	client := s.authenticator.NewClient(token)
	s.client = &client
}

// GetUsername requests and returns the username of the logged in user.
func (s *Spotify) GetUsername() string {
	user, err := s.client.CurrentUser()
	if err != nil {
		panic("Failed to read Spotify profile data.")
	}

	return user.DisplayName
}
