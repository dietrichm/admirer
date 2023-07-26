//go:generate mockgen -source lastfm.go -destination lastfm_mock.go -package lastfm

package lastfm

import (
	"errors"
	"fmt"
	"os"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/config"
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
	GetLovedTracks(args map[string]interface{}) (result lastfm.UserGetLovedTracks, err error)
}

// TrackAPI is our interface for a Last.fm track API.
type TrackAPI interface {
	Love(args map[string]interface{}) (err error)
}

// Lastfm is the external Lastfm service implementation.
type Lastfm struct {
	api      API
	userAPI  UserAPI
	trackAPI TrackAPI
	secrets  config.Config
}

// NewLastfm creates a Lastfm instance.
func NewLastfm(secrets config.Config) (*Lastfm, error) {
	clientID := os.Getenv("LASTFM_CLIENT_ID")
	clientSecret := os.Getenv("LASTFM_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		return nil, errors.New("please set LASTFM_CLIENT_ID and LASTFM_CLIENT_SECRET environment variables")
	}

	api := lastfm.New(clientID, clientSecret)
	api.SetSession(secrets.GetString("session_key"))

	return &Lastfm{
		api:      api,
		userAPI:  api.User,
		trackAPI: api.Track,
		secrets:  secrets,
	}, nil
}

// Name returns the human readable service name.
func (l *Lastfm) Name() string {
	return "Last.fm"
}

// Authenticated returns whether the service is logged in.
func (l *Lastfm) Authenticated() bool {
	return l.api.GetSessionKey() != ""
}

// CreateAuthURL returns an authorization URL to authorize the integration.
func (l *Lastfm) CreateAuthURL(redirectURL string) string {
	return l.api.GetAuthRequestUrl(redirectURL)
}

// CodeParam is the query parameter name used in the authentication callback.
func (l *Lastfm) CodeParam() string {
	return "token"
}

// Authenticate takes an authorization code and authenticates the user.
func (l *Lastfm) Authenticate(oauthCode string, redirectURL string) error {
	if err := l.api.LoginWithToken(oauthCode); err != nil {
		return fmt.Errorf("failed to authenticate on Last.fm: %w", err)
	}

	l.secrets.Set("session_key", l.api.GetSessionKey())

	if err := l.secrets.Save(); err != nil {
		return fmt.Errorf("failed to save Last.fm secrets: %w", err)
	}

	return nil
}

// GetUsername requests and returns the username of the logged in user.
func (l *Lastfm) GetUsername() (string, error) {
	user, err := l.userAPI.GetInfo(lastfm.P{})
	if err != nil {
		return "", fmt.Errorf("failed to read Last.fm profile data: %w", err)
	}

	return user.Name, nil
}

// GetLovedTracks returns loved tracks from the external service.
func (l *Lastfm) GetLovedTracks(limit int) (tracks []domain.Track, err error) {
	username, err := l.GetUsername()
	if err != nil {
		return
	}

	result, err := l.userAPI.GetLovedTracks(lastfm.P{
		"user":  username,
		"limit": limit,
	})
	if err != nil {
		return tracks, fmt.Errorf("failed to read Last.fm loved tracks: %w", err)
	}

	for _, resultTrack := range result.Tracks {
		track := domain.Track{
			Artist: resultTrack.Artist.Name,
			Name:   resultTrack.Name,
		}
		tracks = append(tracks, track)
	}
	return
}

// LoveTrack marks a track as loved on the external service.
func (l *Lastfm) LoveTrack(track domain.Track) error {
	if err := l.trackAPI.Love(lastfm.P{
		"track":  track.Name,
		"artist": track.Artist,
	}); err != nil {
		return fmt.Errorf("failed to mark track as loved on Last.fm: %w", err)
	}

	return nil
}

// Close persists any state before quitting the application.
func (l *Lastfm) Close() error {
	return nil
}
