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
		service, err := spotify.NewSpotify()
		if err != nil {
			panic(err)
		}
		return service
	case "lastfm":
		service, err := lastfm.NewLastfm()
		if err != nil {
			panic(err)
		}
		return service
	default:
		panic("Unknown service " + serviceName)
	}
}
