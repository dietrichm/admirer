//go:generate mockgen -source lastfm.go -destination lastfm_mock.go -package lastfm

package lastfm

import (
	"errors"
	"os"

	"github.com/shkh/lastfm-go/lastfm"
)

// API is our interface for a Last.fm API.
type API interface {
	GetAuthRequestUrl(callback string) (uri string)
	LoginWithToken(token string) (err error)
	GetSessionKey() (sk string)
}

// UserAPI is our interface for a Last.fm user API.
type UserAPI interface {
	GetInfo(args map[string]interface{}) (result lastfm.UserGetInfo, err error)
}

// Lastfm is the external Lastfm service implementation.
type Lastfm struct {
	api     API
	userAPI UserAPI
}

// NewLastfm creates a Lastfm instance.
func NewLastfm(sessionKey string) (*Lastfm, error) {
	clientID := os.Getenv("LASTFM_CLIENT_ID")
	clientSecret := os.Getenv("LASTFM_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		return nil, errors.New("please set LASTFM_CLIENT_ID and LASTFM_CLIENT_SECRET environment variables")
	}

	api := lastfm.New(clientID, clientSecret)
	api.SetSession(sessionKey)

	return &Lastfm{
		api:     api,
		userAPI: api.User,
	}, nil
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
func (l *Lastfm) Authenticate(oauthCode string) error {
	if err := l.api.LoginWithToken(oauthCode); err != nil {
		return errors.New("failed to parse Last.fm token")
	}

	return nil
}

// GetUsername requests and returns the username of the logged in user.
func (l *Lastfm) GetUsername() (string, error) {
	user, err := l.userAPI.GetInfo(lastfm.P{})
	if err != nil {
		return "", errors.New("failed to read Last.fm profile data")
	}

	return user.Name, nil
}

// AccessToken returns the API access token to persist.
func (l *Lastfm) AccessToken() string {
	return l.api.GetSessionKey()
}
