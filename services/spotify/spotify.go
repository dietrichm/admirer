//go:generate mockgen -source ../../services/spotify/spotify.go -destination ../../services/spotify/spotify_mock.go -package spotify

package spotify

import (
	"errors"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// Authenticator is our interface for a Spotify authenticator.
type Authenticator interface {
	SetAuthInfo(clientID, secretKey string)
	AuthURL(state string) string
	Exchange(code string) (*oauth2.Token, error)
	NewClient(token *oauth2.Token) spotify.Client
}

// Client is our interface for a Spotify client.
type Client interface {
	CurrentUser() (*spotify.PrivateUser, error)
}

// Spotify is the external Spotify service implementation.
type Spotify struct {
	authenticator Authenticator
	client        Client
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
