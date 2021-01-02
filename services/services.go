package services

import (
	"github.com/dietrichm/admirer/services/lastfm"
	"github.com/dietrichm/admirer/services/spotify"
)

// Service is the external service interface.
type Service interface {
	Name() string
	CreateAuthURL() string
	Authenticate(code string)
	GetUsername() string
}

// ForName returns service instance for service name.
func ForName(serviceName string) Service {
	switch serviceName {
	case "spotify":
		return spotify.NewSpotify()
	case "lastfm":
		return lastfm.NewLastfm()
	default:
		panic("Unknown service " + serviceName)
	}
}
