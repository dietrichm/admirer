package spotify

import (
	"errors"
	"os"

	"github.com/zmb3/spotify"
)

// Spotify is the external Spotify service implementation.
type Spotify struct {
	authenticator *spotify.Authenticator
	client        *spotify.Client
}

// NewSpotify creates a Spotify instance.
func NewSpotify() (*Spotify, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		return nil, errors.New("please set SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET environment variables")
	}

	// Not an actual web server (yet).
	redirectURL := "https://admirer.test"
	authenticator := spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate)
	authenticator.SetAuthInfo(clientID, clientSecret)

	return &Spotify{
		authenticator: &authenticator,
	}, nil
}

// Name returns the human readable service name.
func (s *Spotify) Name() string {
	return "Spotify"
}

// CreateAuthURL returns an authorization URL to authorize the integration.
func (s *Spotify) CreateAuthURL() string {
	return s.authenticator.AuthURL("")
}

// Authenticate takes an authorization code and authenticates the user.
func (s *Spotify) Authenticate(code string) error {
	token, err := s.authenticator.Exchange(code)
	if err != nil {
		return errors.New("failed to parse Spotify token")
	}

	client := s.authenticator.NewClient(token)
	s.client = &client
	return nil
}

// GetUsername requests and returns the username of the logged in user.
func (s *Spotify) GetUsername() (string, error) {
	user, err := s.client.CurrentUser()
	if err != nil {
		return "", errors.New("failed to read Spotify profile data")
	}

	return user.DisplayName, nil
}
